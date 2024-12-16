package repository

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func NewConn(dataSource string) sqlx.SqlConn {
	return sqlx.NewSqlConn("pgx", dataSource)
}

func IsTx(conn sqlx.SqlConn) bool {
	_, err := conn.RawDB()
	return err != nil
}

func NowUTC() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05+00")
}
