package sonex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SonexSwapTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexSwapTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexSwapTokensLogic {
	return &SonexSwapTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var ORFERBY_ALL = "ALL"
var ORFERBY_TOP = "TOP"           // TOPFIVE
var ORFERBY_OFFICIAL = "OFFICIAL" // OFFICIAL

type SwapTokenInfoCacheWarpper struct {
	Data       []repository.TokenInfoTmp
	LastUpdate time.Time
}

func getTop5VolumeIn24Hours(cacher cache.ICache, eventSonexPoolSwapModel repository.EventSonexPoolSwapModel, chainId int64) ([]repository.TokenInfoTmp, error) {

	var volumes []repository.TokenInfoTmp
	volumeCacheKey := fmt.Sprintf("sonex:swap_page_tokens_volume:%d", chainId)
	// Check cache
	var cacheWarpper SwapTokenInfoCacheWarpper
	lastUpdate, err := cacher.Get(volumeCacheKey)
	if err == nil && !strings.EqualFold(lastUpdate, "") {
		err := json.Unmarshal([]byte(lastUpdate), &cacheWarpper)
		if err == nil {
			volumes = cacheWarpper.Data
			if time.Now().UTC().After(cacheWarpper.LastUpdate.Add(time.Hour)) {
				volumes = nil
			}
		}
	}

	if len(volumes) == 0 {
		volumes, err = eventSonexPoolSwapModel.SummarizeVolumesUsdByChainId(context.Background(), chainId, 5)
		if err != nil {
			logx.Errorf("Summarize volumes error: %s", err)
			return []repository.TokenInfoTmp{}, err
		}
		if len(volumes) == 0 {
			logx.Error("No volume  data")
			return []repository.TokenInfoTmp{}, errors.New("no volume data ")
		}

		respBytes, err := json.Marshal(SwapTokenInfoCacheWarpper{Data: volumes, LastUpdate: time.Now().UTC()})
		if err != nil {
			logx.Errorf("Marshal response error: %v", err)
			return []repository.TokenInfoTmp{}, err
		}
		err = cacher.Set(volumeCacheKey, respBytes, time.Hour*1)
		if err != nil {
			logx.Errorf("Set cache error: %v", err)
			return []repository.TokenInfoTmp{}, err
		}
	}

	return volumes, nil
}

// for tokenList
func getTokenBalanceMultipleCall(chain config.ChainConfig, tokenList []types.SwapTokenInfo, userAddress string) ([]types.SwapTokenInfo, error) {
	if len(tokenList) == 0 {
		return tokenList, errors.New("no token list")
	}
	if strings.EqualFold(userAddress, "") {
		logx.Info("useraddress is nil")
		return tokenList, nil
	}

	var wg sync.WaitGroup

	poolChan := make(chan types.SwapTokenInfo, len(tokenList))
	pools := make([]types.SwapTokenInfo, 0, len(tokenList))

	for _, v := range tokenList {
		wg.Add(1)
		go func(swapTokenInfo types.SwapTokenInfo) {
			defer wg.Done()
			balance := big.NewInt(0)
			Multicall(
				context.Background(),
				chain,
				Call{
					Contract: common.HexToAddress(swapTokenInfo.Address),
					Method:   "balanceOf",
					Args:     []interface{}{common.HexToAddress(userAddress)},
					Result:   &balance,
				},
			)
			// balance 大于 0
			if balance.Cmp(big.NewInt(0)) > 0 {
				swapTokenInfo.Balance = balance.String()
				userBalanceReadableUsd, userBalanceReadable := global.CalucateOneTokenUsdInPool(balance, swapTokenInfo.Price, int64(6), swapTokenInfo.Decimals)
				swapTokenInfo.UsdValue = userBalanceReadableUsd.String()
				swapTokenInfo.Balance = userBalanceReadable.String()
				logx.Infof("swapTokenInfo.Address: %v  balance: %v", swapTokenInfo.Address, balance)
			}

			poolChan <- swapTokenInfo
		}(v)
	}

	go func() {
		wg.Wait()
		close(poolChan)
	}()

	for pool := range poolChan {
		pools = append(pools, pool)
	}
	return pools, nil

}
func getUserTokenDetailStructFromScout(chain config.ChainConfig, tokenList []types.SwapTokenInfo, userAddressWithCheckSum string) ([]types.SwapTokenInfo, error) {

	if !strings.EqualFold(userAddressWithCheckSum, "") {
		if !common.IsHexAddress(userAddressWithCheckSum) {
			logx.Infof("invalid user address:%v ", userAddressWithCheckSum)
			return tokenList, nil
		}
	}

	if !strings.EqualFold(userAddressWithCheckSum, "") {

		tokenWithBalance, err := getTokenBalanceMultipleCall(chain, tokenList, userAddressWithCheckSum)
		if err != nil {
			logx.Infof("getTokenBalanceMultipleCall ERR:%v ", err)
			return tokenList, err
		}
		return tokenWithBalance, nil

	}

	return tokenList, nil
}

// // ALL / OFFICAL
func getTokenList(conn sqlx.SqlConn, chainId int64, orderBy string) ([]types.SwapTokenInfo, error) {
	var tokenList []types.SwapTokenInfo
	var err error
	tokenModel := repository.NewTokenModel(conn)

	// OFFICAL : sonex_token
	if strings.EqualFold(orderBy, ORFERBY_OFFICIAL) {

		tokenList, err = tokenModel.QueryByChainIdWithPrice(context.Background(), chainId, true)
		if err != nil {
			logx.Errorf(" QueryByChainIdWithPrice ERR:%v ", err)
			return []types.SwapTokenInfo{}, err
		}

		return tokenList, nil

	}

	// 选择 ALL : token
	if strings.EqualFold(orderBy, ORFERBY_ALL) {
		tokenList, err = tokenModel.QueryByChainIdWithPrice(context.Background(), chainId, false)
		if err != nil {
			logx.Errorf(" QueryByChainIdWithPrice ERR:%v ", err)
			return []types.SwapTokenInfo{}, err
		}

		return tokenList, nil

	}
	return tokenList, nil
}
func (l *SonexSwapTokensLogic) SonexSwapTokens(req *types.SwapTokensRequest) (*types.SwapTokensResponse, error) {
	if req.ChainId == 0 {
		return &types.SwapTokensResponse{}, errors.New("invalid chainid")
	}

	if !strings.EqualFold(req.Address, "") {
		if !common.IsHexAddress(req.Address) {
			logx.Errorf("invalid user address:%v ", req.Address)
			return &types.SwapTokensResponse{}, errors.New("invalid address")
		}
	}

	// 判断conf
	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return &types.SwapTokensResponse{}, errors.New("miss chain configure")
	}

	// checksum address
	userAddressWithCheckSum := common.HexToAddress(req.Address).Hex()

	// 默认用ALL字段
	if !strings.EqualFold(req.OrderBy, ORFERBY_ALL) &&
		!strings.EqualFold(req.OrderBy, ORFERBY_OFFICIAL) &&
		!strings.EqualFold(req.OrderBy, ORFERBY_TOP) {
		req.OrderBy = ORFERBY_ALL
	}

	// 设置 10分钟 缓存

	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	// var volumes []types.SwapTokenInfo
	var err error
	// 选择 TOP  FIVE
	if strings.EqualFold(req.OrderBy, ORFERBY_TOP) {

		volumes, err := getTop5VolumeIn24Hours(cacher, eventSonexPoolSwapModel, req.ChainId)

		if err != nil {
			logx.Errorf("getTop5VolumeIn24Hours ERR:%v ", err)
			return &types.SwapTokensResponse{}, err
		}
		var lists []types.SwapTokenInfo
		for _, v := range volumes {
			isOffical := int64(0)
			if v.IsOffical {
				isOffical = 1
			}
			list := types.SwapTokenInfo{
				Address:    v.Address,  // string `json:"address"`
				Name:       v.Name,     //, string `json:"name"`
				Symbol:     v.Symbol,   //  string `json:"symbol"`
				Volume:     v.Volume,   // string `json:"volume,optional"` // TOP
				Price:      v.Price,    //  string `json:"price,optional"`
				Decimals:   v.Decimals, // int64  `json:"decimals"`
				Logo:       v.Logo,     //  string `json:"logo,optional"`
				IsOfficial: isOffical,  // int64  `json:"isOffical,optional"` // 1:官方代币,
			}
			lists = append(lists, list)
		}

		return &types.SwapTokensResponse{List: lists}, nil

	}

	var tokenList []types.SwapTokenInfo

	tokenList, err = getTokenList(conn, req.ChainId, req.OrderBy) // ALL / OFFICAL
	if err != nil {
		logx.Errorf("getTokenList ERR:%v ", err)
		return &types.SwapTokensResponse{}, err
	}

	// 没有数据返回
	if len(tokenList) == 0 {
		logx.Error("no token list")
		return &types.SwapTokensResponse{}, errors.New("no token data ")
	}

	// get user balance from block scout
	// 获取该user address 的所有代币
	tokenListWithBalance, err := getUserTokenDetailStructFromScout(*chainConf, tokenList, userAddressWithCheckSum)
	if err != nil {
		logx.Errorf("getUserTokenDetailStructFromScout:%v ", err)
		// 默认给个值
		return &types.SwapTokensResponse{List: tokenList}, nil
	}

	logx.Infof("tokenListWithBalance:%#v ", tokenListWithBalance)

	sort.SliceStable(tokenListWithBalance, func(i, j int) bool {
		usd_i, _ := big.NewFloat(0).SetString(tokenList[i].UsdValue)
		usd_j, _ := big.NewFloat(0).SetString(tokenList[j].UsdValue)
		if usd_i.Cmp(usd_j) == 0 {
			// 按照 symbol 排序序列
			return tokenListWithBalance[i].Symbol < tokenListWithBalance[i].Symbol
		} else {
			return usd_i.Cmp(usd_j) > 0
		}
	})

	return &types.SwapTokensResponse{List: tokenListWithBalance}, nil

}

/*

type TokenDetailStruct struct {
	Address              string `json:"address"`
	CirculatingMarketCap string `json:"circulating_market_cap"`
	Decimals             string `json:"decimals"`
	ExchangeRate         string `json:"exchange_rate"`
	Holders              string `json:"holders"`
	IconUrl              string `json:"icon_url"`
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	TotalSupply          string `json:"total_supply"`
	TypeString           string `json:"type"`
	Volume24h            string `json:"volume_24h"`
}

type TokenUserStruct struct {
	TokenDetail   TokenDetailStruct `json:"token"`
	TokenId       string            `json:"token_id"`
	TokenInstance string            `json:"token_instance"`
	Value         string            `json:"value"`
}

func roundRequestIds(RequestIds []string) func() string {
	var n uint64
	return func() string {
		tn := atomic.LoadUint64(&n)
		atomic.AddUint64(&n, 1)
		return RequestIds[int(tn)%len(RequestIds)]
	}
}

// 从浏览器来的userbalance 数据
func getUserBalanceFromBlockScout(url, x_request_id string) ([]*TokenUserStruct, error) {
	var tokenDetails []*TokenUserStruct

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	logx.Infof("url:%v ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logx.Errorf("Get getUserBalanceFromBlockScout NewRequest:%v ", err)
		return tokenDetails, err
	}

	req.Header.Add("x-request-id", x_request_id)

	response, err := client.Do(req)
	if err != nil {
		logx.Errorf("getUserBalanceFromBlockScout DO error:%v ", err)
		return tokenDetails, err
	}

	defer response.Body.Close()

	// Body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logx.Errorf("getUserBalanceFromBlockScout ReadAll error:%v ", err)
		return tokenDetails, err
	}
	logx.Infof("sssss:%v ", string(body))
	// token details
	err = json.Unmarshal(body, &tokenDetails)
	if err != nil {
		logx.Errorf("Unmarshal error:%v ", err)
		return tokenDetails, err
	}

	return tokenDetails, nil
}
*/
