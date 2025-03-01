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
	conn := repository.NewConn(config.PostgreSQL.DataSource())

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

	blk, err := sqlx.NewBulkInserter(conn, "insert into sonex_token_price_record (token_symbol, token_address, price, created_at, usd_decimals) values (?,?,?,?,?)")
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

			prices := []*cache.Price{}
			{
				amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals), nil))
				successed, _, err := global.QuoteExactInput(context.Background(), nil, item.ChainId, item.Address, chainConf.USD, amountIn, 3)
				if err != nil {
					logx.Errorf("Get %s[%s] price failed: %s", item.Address, item.Symbol, err)
					return
				}
				if len(successed) > 0 {
					var ps []*cache.Price
					linq.From(successed).
						WhereT(func(item cache.Price) bool {
							return math.Abs(item.PriceImpact) < 0.1
						}).
						SelectT(func(item cache.Price) *cache.Price {
							return &item
						}).
						ToSlice(&ps)
					prices = append(prices, ps...)
				}
			}

			{
				accuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals/2), nil)
				amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals-usdToken.Decimals/2), nil))
				successed, _, err := global.QuoteExactInput(context.Background(), nil, item.ChainId, item.Address, chainConf.USD, amountIn, 3)
				if err != nil {
					logx.Errorf("Get %s[%s] price failed: %s", item.Address, item.Symbol, err)
					return
				}
				if len(successed) > 0 {
					var ps []*cache.Price
					linq.From(successed).
						WhereT(func(item cache.Price) bool {
							return math.Abs(item.PriceImpact) < 0.1
						}).
						SelectT(func(item cache.Price) *cache.Price {
							item.Value = new(big.Int).Mul(item.Value, accuracy)
							return &item
						}).
						ToSlice(&ps)
					prices = append(prices, ps...)
				}
			}

			{
				accuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil)
				amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals-usdToken.Decimals), nil))
				successed, _, err := global.QuoteExactInput(context.Background(), nil, item.ChainId, item.Address, chainConf.USD, amountIn, 3)
				if err != nil {
					logx.Errorf("Get %s[%s] price failed: %s", item.Address, item.Symbol, err)
					return
				}
				if len(successed) > 0 {
					var ps []*cache.Price
					linq.From(successed).
						WhereT(func(item cache.Price) bool {
							return math.Abs(item.PriceImpact) < 0.1
						}).
						SelectT(func(item cache.Price) *cache.Price {
							item.Value = new(big.Int).Mul(item.Value, accuracy)
							return &item
						}).
						ToSlice(&ps)
					prices = append(prices, ps...)
				}
			}

			average := linq.From(prices).SelectT(func(item *cache.Price) float64 {
				value, _ := new(big.Float).SetInt(item.Value).Float64()
				return value
			}).Average()
			priceA := linq.From(prices).
				WhereT(func(item *cache.Price) bool {
					value, _ := new(big.Float).SetInt(item.Value).Float64()
					return math.Abs(item.PriceImpact) < 0.1 && value >= average
				}).
				OrderByT(func(item *cache.Price) float64 {
					return math.Abs(item.PriceImpact)
				}).First()
			var bestPrice *cache.Price
			if price, ok := priceA.(*cache.Price); ok {
				bestPrice = price
			}

			if bestPrice == nil {
				logx.Errorf("Get %s[%s] best price failed, Success count: %d", item.Address, item.Symbol, len(prices))
				return
			}

			var paths []string
			linq.From(bestPrice.Path.Tokens).SelectT(func(item common.Address) string { return item.Hex() }).ToSlice(&paths)
			logx.Infof("Best price %s %s %s %s", item.Symbol, item.Address, bestPrice.Value.String(), strings.Join(paths, "->"))

			err = blk.Insert(item.Symbol, item.Address, bestPrice.Value.String(), time.Now().UTC().Format("2006-01-02 15:04:05+00"), usdToken.Decimals)
			if err != nil {
				logx.Errorf("Insert %s[%s] price failed: %s", item.Address, item.Symbol, err)
				return
			}
		}(&token)
	}
	wg.Wait()

	return nil
}
