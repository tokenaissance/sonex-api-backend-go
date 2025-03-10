package task

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/util"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

func HandleXpoolPoints(config *config.Config) {
	logx.Info("HandleXpoolPoints start")
	retrys := 0
	err := retry.Do(
		func() error {
			conn := repository.NewConn(config.PostgreSQL.DataSource())
			err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
				for _, chainConf := range config.Chains {
					farmConf := config.GetFarmConfig(chainConf.ChainId)
					if farmConf == nil {
						continue
					}
					err := xpoolPointsLogic(ctx, sqlx.NewSqlConnFromSession(s), *config, chainConf)
					if err != nil {
						logx.Errorf("[HandleXpoolPoints]: failed: %d %s", chainConf.ChainId, err)
						return err
					}
				}

				return nil
			})
			if err != nil {
				return err
			}

			return nil
		},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.Attempts(5),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("[HandleXpoolPoints]: retry %d times", retrys)
		}),
	)
	if err != nil {
		logx.Errorf("[HandleXpoolPoints]: failed: %s", err)
	}

	logx.Info("[HandleXpoolPoints]: successfully")
}

func xpoolPointsLogic(ctx context.Context, conn sqlx.SqlConn, conf config.Config, chainConf config.ChainConfig) error {
	xpoolConf := conf.GetXpoolConfig(chainConf.ChainId)

	factoryModel := repository.NewEventXpoolFactoryPoolCreatedModel(conn)
	depositModel := repository.NewEventXpoolPoolDepositedModel(conn)
	withdrawModel := repository.NewEventXpoolPoolWithdrawedModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	xpoolPointsRecordModel := repository.NewXpoolPointsRecordModel(conn)

	// get all available && destroied pools
	poolsEvents, err := factoryModel.QueryByChainIdFactory(ctx, chainConf.ChainId, xpoolConf.Factory)
	if err != nil {
		return err
	}

	rewardTokens := make([]string, 0)
	for _, pool := range poolsEvents {
		token := common.HexToAddress(pool.EventRewardToken).Hex()
		rewardTokens = append(rewardTokens, token)
	}
	prices, err := getCurrentPrices(ctx, tokenPriceModel, conf, chainConf.ChainId, rewardTokens)
	if err != nil {
		return err
	}
	// 70 分钟前
	beginTime := time.Now().UTC().Add(-70 * time.Minute)
	// 10 分钟前
	endTime := time.Now().UTC().Add(-10 * time.Minute)

	// 针对每个池子，先找出该池子 7 天前有质押的用户地址
	// 判断用户是否被发过积分，发过就不再走积分统计逻辑
	// 没发过积分的用户，查询前七天的质押和解质押表，统计每一天是否质押 amount > 0
	// 第八天（当天）看，第七天质押代币余额的当前价值 > 100U，则发放积分
	for _, pool := range poolsEvents {

		// 获取池子地址
		poolAddress := pool.EventPool
		// 获取池子中用户积分
		userPointsMap, err := xpoolPointsRecordModel.GetUserPointsByPoolAddress(ctx, chainConf.ChainId, poolAddress)
		if err != nil {
			logx.Errorf("[HandleXpoolPoints]: GetUserPointsByPoolAddress failed: %s", err)
			return err
		}

		// 获取池子中用户质押事件
		depositedEvents, err := depositModel.GetDepositedEventsByPoolAddress(ctx, chainConf.ChainId, poolAddress)
		if err != nil {
			logx.Errorf("[HandleXpoolPoints]: GetDepositedEventsByPoolAddress failed: %s", err)
			return err
		}

		// 如果池子没有用户质押事件，则跳过
		if len(depositedEvents) == 0 {
			logx.Infof("[HandleXpoolPoints]: No deposited events for pool %s", poolAddress)
			continue
		}

		// 过滤出未发放积分的用户地址
		userAddresses := make([]string, 0)
		for userAddress, _ := range depositedEvents {
			if _, ok := userPointsMap[userAddress]; !ok {
				userAddresses = append(userAddresses, userAddress)
			}
		}

		// 获取池子中用户提款事件
		withdrawedEvents, err := withdrawModel.GetWithdrawedEventsByPoolAddress(ctx, chainConf.ChainId, poolAddress)
		if err != nil {
			return err
		}

		// 根据 userAddresses, 针对 depositedEvents 和 withdrawedEvents 的 TxAt 进行按照日期的求和
		// map[userAddress][timeStr][amount]
		depositedAmountMap := make(map[string]map[string]*big.Int)
		for _, userAddress := range userAddresses {
			depositedAmountMap[userAddress] = make(map[string]*big.Int)
		}

		// 统计每个时间点的质押 & 解质押之后的余额
		for _, userAddress := range userAddresses {
			m := depositedAmountMap[userAddress]

			// 质押时间按照时间升序
			sort.Slice(depositedEvents[userAddress], func(i, j int) bool {
				return depositedEvents[userAddress][i].TxAt.Before(depositedEvents[userAddress][j].TxAt)
			})

			// 计提每个时间点的质押金额
			for _, depositedEvent := range depositedEvents[userAddress] {
				depositedAmount, _ := big.NewInt(0).SetString(depositedEvent.EventAmount, 10)
				timeStr := depositedEvent.TxAt.UTC().Format("2006-01-02 15:04:05")
				currentDepositedAmount, ok := m[timeStr]
				if !ok {
					m[timeStr] = big.NewInt(0)
					currentDepositedAmount = big.NewInt(0)
				}

				m[timeStr] = big.NewInt(0).Add(currentDepositedAmount, depositedAmount)
			}

			// 提款时间按照时间升序
			sort.Slice(withdrawedEvents[userAddress], func(i, j int) bool {
				return withdrawedEvents[userAddress][i].TxAt.Before(withdrawedEvents[userAddress][j].TxAt)
			})

			// 计减每个时间点的提款金额
			for _, withdrawedEvent := range withdrawedEvents[userAddress] {
				withdrawedAmount, _ := big.NewInt(0).SetString(withdrawedEvent.EventWithdrawn, 10)
				timeStr := withdrawedEvent.TxAt.UTC().Format("2006-01-02 15:04:05")
				currentWithdrawedAmount, ok := m[timeStr]
				if !ok {
					m[timeStr] = big.NewInt(0)
					currentWithdrawedAmount = big.NewInt(0)
				}

				m[timeStr] = big.NewInt(0).Sub(currentWithdrawedAmount, withdrawedAmount)
			}

			// 时间点排序
			timeStrs := make([]string, 0)
			for timeStr := range m {
				timeStrs = append(timeStrs, timeStr)
			}
			// 按照时间字符串排序（因为使用了标准格式，可以直接用字符串排序）
			sort.Strings(timeStrs)

			// 计算每个时间点的累积余额
			total := big.NewInt(0)
			for _, timeStr := range timeStrs {
				amount := m[timeStr]
				total = big.NewInt(0).Add(total, amount)
				m[timeStr] = big.NewInt(0).Set(total) // 更新为累积值
			}
			depositedAmountMap[userAddress] = m

			// 查询前七天到前一天的余额，判断每一天的余额是否 > 0，且第七天余额 > 100U，并且对满足条件的地址发放积分
			qualified, lastDayAmount, err := checkUserQualification(m, beginTime, endTime)
			if err != nil {
				logx.Errorf("[HandleXpoolPoints]: Check qualification failed: %s", err)
				continue
			}

			// 如果用户不满足条件，则跳过
			if !qualified {
				continue
			}

			rewardTokenPrice, ok := prices[pool.EventRewardToken]
			if !ok {
				logx.Errorf("[HandleXpoolPoints]: Get reward token price failed: %s", pool.EventRewardToken)
				continue
			}

			token, err := global.GetToken(chainConf.ChainId, pool.EventRewardToken)
			if err != nil {
				logx.Errorf("[HandleXpoolPoints]: Get usd value failed: %s", err)
				continue
			}
			tokenValue, err := sonex.NewTokenValue(lastDayAmount.String(), int(token.Decimals), sonex.WithTokenPrice(rewardTokenPrice))
			if err != nil {
				logx.Errorf("[HandleXpoolPoints]: Get usd value failed: %s", err)
				continue
			}

			// 如果第七天余额 <= 100U，则跳过
			if tokenValue.Value().Cmp(big.NewFloat(100)) <= 0 {
				continue
			}

			// 发放积分
			id, err := global.NewId()
			if err != nil {
				return err
			}
			// 创建积分记录
			err = xpoolPointsRecordModel.New(ctx, id, chainConf.ChainId, poolAddress, userAddress, xpoolConf.Points, lastDayAmount.String(), token.Decimals, rewardTokenPrice.Price(), int64(rewardTokenPrice.Decimals()))
			if err != nil {
				logx.Errorf("[HandleXpoolPoints]: New points record failed: %s", err)
				continue
			}

			req := struct {
				User      string            `json:"user_address"`
				Points    int64             `json:"points"`
				Source    string            `json:"source"`
				SubTitle  string            `json:"sub_title"`
				ExtraId   string            `json:"extra_id"`
				Extra     map[string]string `json:"extra"`
				Signature string            `json:"signature"`
			}{
				User:     userAddress,
				Points:   xpoolConf.Points,
				Source:   conf.PointsService.Source,
				SubTitle: xpoolConf.SubTitle,
				ExtraId:  fmt.Sprintf("%s-%s", userAddress, poolAddress), // 幂等，通过 extraId 判断是否发放过积分
				Extra: map[string]string{
					"user":          userAddress,
					"poolAddress":   poolAddress,
					"amount":        lastDayAmount.String(),
					"tokenDecimals": strconv.Itoa(int(token.Decimals)),
					"price":         rewardTokenPrice.Price(),
					"usdDecimals":   strconv.Itoa(int(rewardTokenPrice.Decimals())),
				},
			}
			signature, err := util.Sign(&req, conf.PointsService.PrivateKey)
			if err != nil {
				logx.Errorf("[HandleXpoolPoints]: Sign failed: %s", err)
				continue
			}
			req.Signature = signature

			url := fmt.Sprintf("%s/extra/increase_points", strings.TrimRight(conf.PointsService.Host, "/"))
			resp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
			if err != nil {
				logx.Errorf("[HandleXpoolPoints]: Increase points failed: %s, userAddress: %s, poolAddress: %s", err, userAddress, poolAddress)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logx.Errorf("[HandleXpoolPoints]: Increase points failed: %s, userAddress: %s, poolAddress: %s", resp.Status, userAddress, poolAddress)
				continue
			}
			logx.Infof("HandleXpoolPoints increase points: User %s, PoolAddress %s, Points %d",
				userAddress,
				poolAddress,
				xpoolConf.Points,
			)
		}

	}

	return nil
}

// checkUserQualification 检查用户是否满足获取积分的条件
func checkUserQualification(balanceMap map[string]*big.Int, beginTime, endTime time.Time) (bool, *big.Int, error) {
	hasBalance := false

	// 检查从beginTime到endTime的每个10分钟
	currentTime := beginTime
	for currentTime.Before(endTime) || currentTime.Equal(endTime) {
		// 将时间转换为格式化的字符串作为map的key
		timeKey := currentTime.Format("2006-01-02 15:04:05")
		amount, exists := balanceMap[timeKey]
		if !exists || amount.Cmp(big.NewInt(0)) <= 0 {
			return false, nil, nil
		}
		hasBalance = true

		// 移到下一个10分钟
		currentTime = currentTime.Add(10 * time.Minute)
	}

	// 获取最后一个时间点的余额
	lastTimeKey := endTime.Format("2006-01-02 15:04:05")
	lastAmount, exists := balanceMap[lastTimeKey]
	if !exists {
		return false, nil, nil
	}

	return hasBalance, lastAmount, nil
}

// getCurrentPrices 获取当前价格
func getCurrentPrices(ctx context.Context, tokenPriceModel repository.SonexTokenPriceRecordModel, conf config.Config, chainId int64, addressList []string) (prices map[string]*sonex.TokenPrice, err error) {
	sonexConf := conf.GetSonexConfig(chainId)
	if sonexConf == nil {
		return nil, fmt.Errorf("%s", "Miss sonex configure")
	}
	xTokenAddr := common.HexToAddress(sonexConf.XToken)

	chainConf := conf.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, fmt.Errorf("Miss chain configure")
	}
	usdToken, err := global.GetToken(chainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}

	// unique address list
	filterM := make(map[string]bool)
	uniqueAddr := make([]string, 0)
	for _, addr := range addressList {
		if _, ok := filterM[addr]; ok {
			continue
		}
		filterM[addr] = true
		uniqueAddr = append(uniqueAddr, addr)
	}

	addressList = uniqueAddr

	records, err := tokenPriceModel.GetCurrentPrices(ctx, addressList)
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get current prices error: %s", err)
		return nil, err
	}
	prices = make(map[string]*sonex.TokenPrice, 0)
	for _, record := range records {
		tp, err := sonex.NewTokenPrice(record.Price, record.UsdDecimals)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: Invalid token price: %s", err)
			return nil, err
		}
		prices[record.TokenAddress] = tp
	}

	// set xToken price
	for _, addr := range addressList {
		// 1. db 价格存在，pass
		if _, ok := prices[addr]; ok {
			continue
		}

		// 设置默认价格为 1U
		tp, err := sonex.NewTokenPrice("1", 0)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: Invalid token price: %s", err)
			return nil, err
		}

		// 2. db 价格不存在，当 addr 为 xToken 时，设置价格为默认 1U
		if xTokenAddr.Cmp(common.HexToAddress(addr)) == 0 {
			prices[addr] = tp
			continue
		}

		// 3. db 价格不存在，且不是 xToken，则从链上查询数据，如果获取不到，则将价格设置为 1U
		bestPrice0, _, _, err := global.GetPrices(ctx, nil, chainId, common.HexToAddress(addr).Hex(), true, 2)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: getCurrentPrices GetPrices error: %s", err.Error())
		}
		if err == nil && bestPrice0 != nil {
			tp.SetPrice(bestPrice0.Value.String())
			tp.SetDecimals(int(usdToken.Decimals))
		}
		prices[addr] = tp
	}

	return prices, nil
}
