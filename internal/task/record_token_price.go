package task

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"strings"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func HandleRecordTokenPrice(config *config.Config) {
	logx.Info("Token history price updating")
	conn := repository.NewConn(config.PostgreSQL.DateSource())

	for _, chainConf := range config.Chains {
		err := RecordTokenPrice(chainConf, conn)
		if err != nil {
			logx.Errorf("Token history price update failed: %d %s", chainConf.ChainId, err)
			continue
		}
	}
	logx.Info("Token history price updated")
}

func RecordTokenPrice(chainConf config.ChainConfig, conn sqlx.SqlConn) error {
	tokenModel := repository.NewTokenModel(conn)

	tokens, err := tokenModel.QueryByChainId(context.Background(), chainConf.ChainId)
	if err != nil {
		return err
	}
	usdToken, err := global.GetToken(chainConf.ChainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil || usdToken == nil {
		return fmt.Errorf("%s", fmt.Sprintf("Get USD token[%s] failed: %s", chainConf.USD, err))
	}

	blk, err := sqlx.NewBulkInserter(conn, "insert into sonex_token_price_record (token_symbol, token_address,price, created_at, usd_decimals) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer blk.Flush()

	var wg sync.WaitGroup

	for _, token := range tokens {
		wg.Add(1)

		go func(item *repository.Token) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					logx.Errorf("Update %s[%s] price failed(panic): %s", item.Address, item.Symbol, err)
					return
				}
			}()

			var bestPrice *cache.Price

			amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals), nil))
			successed, failed, err := global.QuoteExactInput(context.Background(), nil, item.ChainId, item.Address, chainConf.USD, amountIn, 2)
			if err != nil {
				logx.Errorf("Get %s[%s] price failed: %s", item.Address, item.Symbol, err)
				return
			}
			if len(successed) > 0 {
				priceA := linq.From(successed).
					WhereT(func(item cache.Price) bool {
						return math.Abs(item.PriceImpact) < 0.1
					}).
					OrderByT(func(item cache.Price) float64 {
						return item.PriceImpact
					}).First()
				if price, ok := priceA.(cache.Price); ok {
					bestPrice = &price
				}
			}

			if bestPrice == nil {
				accuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals/2), nil)
				amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals-usdToken.Decimals/2), nil))
				successed, failed, err = global.QuoteExactInput(context.Background(), nil, item.ChainId, item.Address, chainConf.USD, amountIn, 2)
				if err != nil {
					logx.Errorf("Get %s[%s] price failed: %s", item.Address, item.Symbol, err)
					return
				}
				if len(successed) > 0 {
					priceA := linq.From(successed).
						WhereT(func(item cache.Price) bool {
							return math.Abs(item.PriceImpact) < 0.1
						}).
						OrderByT(func(item cache.Price) float64 {
							return item.PriceImpact
						}).First()
					if price, ok := priceA.(cache.Price); ok {
						price.Value = new(big.Int).Mul(price.Value, accuracy)
						bestPrice = &price
					}
				}
			}

			if bestPrice == nil {
				accuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil)
				amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals-usdToken.Decimals), nil))
				successed, failed, err = global.QuoteExactInput(context.Background(), nil, item.ChainId, item.Address, chainConf.USD, amountIn, 2)
				if err != nil {
					logx.Errorf("Get %s[%s] price failed: %s", item.Address, item.Symbol, err)
					return
				}
				if len(successed) > 0 {
					priceA := linq.From(successed).
						WhereT(func(item cache.Price) bool {
							return math.Abs(item.PriceImpact) < 0.1
						}).
						OrderByT(func(item cache.Price) float64 {
							return item.PriceImpact
						}).First()
					if price, ok := priceA.(cache.Price); ok {
						price.Value = new(big.Int).Mul(price.Value, accuracy)
						bestPrice = &price
					}
				}
			}

			if bestPrice == nil {
				failedList := []string{}
				for _, item := range failed {
					pools := []string{}
					for _, pool := range item.Pools {
						pools = append(pools, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
					}
					failedList = append(failedList, fmt.Sprintf("Path: %s, Error: %s", strings.Join(pools, " -> "), item.Err.Error()))
				}
				logx.Errorf("Get %s[%s] best price failed,Success count: %d, Failed list: %v", item.Address, item.Symbol, len(successed), failedList)
				return
			}
			// fmt.Printf("%s %s %s\n", item.Symbol, item.Address, bestPrice.Value.String())
			err = blk.Insert(item.Symbol, item.Address, bestPrice.Value.String(), time.Now().UTC().Format("2006-01-02 15:04:05"), usdToken.Decimals)
			if err != nil {
				logx.Errorf("Insert %s[%s] price failed: %s", item.Address, item.Symbol, err)
				return
			}
		}(&token)
	}
	wg.Wait()

	return nil
}
