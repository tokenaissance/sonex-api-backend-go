package task

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/repository"
	"time"
)

func HandleRecordTokenPrice(config *config.Config) {
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
	for _, token := range tokens {
		//global.GetPrices()
		err = blk.Insert(token.Symbol, token.Address, grand.N(1, 100), time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			panic(err)
		}
	}
}
