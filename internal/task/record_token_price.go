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
	fmt.Println("开始执行token价格记录任务")

	ctx := context.Background()

	conn := repository.NewConn(config.PostgreSQL.DateSource())
	//model := repository.NewSonexTokenModel(conn)
	//results, total, err := model.GetTokenList(context.Background(), req)

	query := fmt.Sprintf("select * from %s", "sonex_token")
	var tokens []repository.TokenData
	_ = conn.QueryRowsCtx(context.Background(), &tokens, query)

	blk, err := sqlx.NewBulkInserter(conn, "insert into sonex_token_price_record (token_symbol, token_address,price, created_at) values (?, ?,?,?)")
	if err != nil {
		panic(err)
	}
	defer blk.Flush()

	var wg sync.WaitGroup

	for _, token := range tokens {
		wg.Add(1)

		go func(item *repository.TokenData) {
			defer wg.Done()

			// todo: 2024-12-29 chainID 固定，目前 sonex_token_price_record 中并没有包含 chainId 字段
			chainId := int64(1946)
			usdAddress := config.GetChainConfig(chainId).USD
			bestPrice, _, err := global.GetPrices(ctx, nil, chainId, item.Address, true, 2, &usdAddress)
			if err != nil {
				logx.Errorf("[HandleRecordTokenPrice]: get token price error: %s", err)
				return
			}
			if bestPrice == nil {
				logx.Errorf("[HandleRecordTokenPrice]: bestPrice == nil, get token price error")
				return
			}
			err = blk.Insert(item.Symbol, item.Address, bestPrice.Value.String(), time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				logx.Errorf("[HandleRecordTokenPrice]: insert token price record error: %s", err)
				return
			}
		}(&token)
	}

	wg.Wait()
}
