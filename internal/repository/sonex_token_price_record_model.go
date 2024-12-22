package repository

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SonexTokenPriceRecordModel = (*customSonexTokenPriceRecordModel)(nil)

type (
	// SonexTokenPriceRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexTokenPriceRecordModel.
	SonexTokenPriceRecordModel interface {
		sonexTokenPriceRecordModel
		withSession(session sqlx.Session) SonexTokenPriceRecordModel
		Kline(address string, interval string) ([]*TokenPriceRecord, error)
	}

	customSonexTokenPriceRecordModel struct {
		*defaultSonexTokenPriceRecordModel
	}
)

// NewSonexTokenPriceRecordModel returns a model for the database table.
func NewSonexTokenPriceRecordModel(conn sqlx.SqlConn) SonexTokenPriceRecordModel {
	return &customSonexTokenPriceRecordModel{
		defaultSonexTokenPriceRecordModel: newSonexTokenPriceRecordModel(conn),
	}
}

func (m *customSonexTokenPriceRecordModel) withSession(session sqlx.Session) SonexTokenPriceRecordModel {
	return NewSonexTokenPriceRecordModel(sqlx.NewSqlConnFromSession(session))
}

type TokenPriceRecord struct {
	Id           int64  `db:"id"`
	TokenSymbol  string `db:"token_symbol"`
	TokenAddress string `db:"token_address"`
	Price        string `db:"price"`
	CreatedAt    string `db:"created_at"`
}

func (m *customSonexTokenPriceRecordModel) Kline(address string, interval string) ([]*TokenPriceRecord, error) {
	// 获取时间范围
	//intervalDuration := getIntervalDuration(interval)

	query := fmt.Sprintf("select %s from %s where token_address = $1 and  created_at >= NOW() - INTERVAL '%s'", sonexTokenPriceRecordRows, m.table, getIntervalDuration(interval))

	var resp []*TokenPriceRecord
	err := m.conn.QueryRows(&resp, query, address)
	switch err {
	case nil:
		for _, v := range resp {
			v.CreatedAt = gtime.New(v.CreatedAt).Format("Y-m-d H:i:s")
		}
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
	//// 构造 SQL 查询
	//query := fmt.Sprintf(`
	//	SELECT id, token_symbol, token_address, price, created_at
	//	FROM sonex_token_price_record
	//	WHERE token_address = ? AND created_at >= NOW() - INTERVAL ?
	//	ORDER BY created_at DESC;
	//`)
	//
	//var records []SonexTokenPriceRecord
	//err := m.conn.QueryRows(&records, query, tokenAddress, intervalDuration)
	//if err != nil {
	//	return nil, err
	//}
	//return records, nil
}

// 根据 interval 获取查询条件的时间范围
func getIntervalDuration(interval string) string {
	switch interval {
	case "1h":
		return "1 hour"
	case "1d":
		return "1 day"
	case "1w":
		return "1 week"
	case "1m":
		return "1 month"
	case "1y":
		return "1 year"
	default:
		return "1 day" // 默认按1天查询
	}
}
