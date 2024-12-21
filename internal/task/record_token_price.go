package task

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"time"
)

func HandleRecordTokenPrice(serverCtx *svc.ServiceContext) {
	conn := repository.NewConn(serverCtx.Config.PostgreSQL.DateSource())
	//model := repository.NewSonexTokenModel(conn)
	//results, total, err := model.GetTokenList(context.Background(), req)

	query := fmt.Sprintf("select * from %s", "sonex_token")
	var tokens []repository.TokenData
	_ = conn.QueryRowsCtx(context.Background(), &tokens, query)

	blk, err := sqlx.NewBulkInserter(conn, "insert into sonex_token_prices (token_symbol, token_address,price, created_at) values (?, ?,?,?)")
	if err != nil {
		panic(err)
	}
	defer blk.Flush()
	for _, token := range tokens {
		err = blk.Insert(token.Symbol, token.Address, 0, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			panic(err)
		}
	}
}
