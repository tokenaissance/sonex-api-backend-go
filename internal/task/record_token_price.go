package task

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func HandleRecordTokenPrice(config *config.Config) {
	logx.Info("Token history price updating")
	conn := repository.NewConn(config.PostgreSQL.DateSource())

	for _, chainConf := range config.Chains {
		err := Update(chainConf, conn)
		if err != nil {
			logx.Errorf("Token history price update failed: %d %s", chainConf.ChainId, err)
		}
	}
	logx.Info("Token history price updated")
}

func Update(chainConf config.ChainConfig, conn sqlx.SqlConn) error {
	tokenModel := repository.NewTokenModel(conn)

	tokens, err := tokenModel.QueryByChainId(context.Background(), chainConf.ChainId)
	if err != nil {
		return err
	}

	blk, err := sqlx.NewBulkInserter(conn, "insert into sonex_token_price_record (token_symbol, token_address,price, created_at) values (?,?,?,?)")
	if err != nil {
		return err
	}
	defer blk.Flush()

	var wg sync.WaitGroup
	var asyncErr error

	for _, token := range tokens {
		wg.Add(1)

		go func(item *repository.Token) {
			defer wg.Done()

			bestPrice, _, err := global.GetPrices(context.Background(), nil, item.ChainId, item.Address, false, 2, &chainConf.USD)
			if err != nil {
				asyncErr = err
				return
			}
			if bestPrice == nil {
				asyncErr = fmt.Errorf("%s", fmt.Sprintf("Can not get %s[%s] price", item.Address, item.Symbol))
				return
			}
			err = blk.Insert(item.Symbol, item.Address, bestPrice.Value.String(), time.Now().UTC().Format("2006-01-02 15:04:05"))
			if err != nil {
				asyncErr = err
				return
			}
		}(&token)
	}
	wg.Wait()
	if asyncErr != nil {
		return asyncErr
	}

	return nil
}
