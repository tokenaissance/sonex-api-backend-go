package repository

import (
	"context"
	"fmt"
	"math/big"

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
		GetLatestHistoryPrices(tokenAddresses []string) (map[string]*big.Int, error)
		GetCurrentPrice(tokenAddress string) (float64, error)
		GetCurrentPrices(tokenAddresses []string) (map[string]*TokenPriceRecord, error)
		GetLastHourPrice(tokenAddress string) (float64, error)
		GetYesterdayPrice(tokenAddress string) (float64, error)
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

	query := fmt.Sprintf("select %s from %s where token_address = $1 and  created_at >= NOW() - INTERVAL '%s' ORDER BY created_at ASC", sonexTokenPriceRecordRows, m.table, getIntervalDuration(interval))

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
}

func (m *customSonexTokenPriceRecordModel) GetLatestHistoryPrices(tokenAddresses []string) (map[string]*big.Int, error) {
	ctx := context.Background()

	type PriceRecord struct {
		TokenAddress string `db:"token_address"`
		Price        string `db:"price"`
	}

	var records []PriceRecord
	query := `
		SELECT DISTINCT ON (token_address) token_address, price
		FROM sonex_token_price_record
		WHERE token_address = ANY($1)
		ORDER BY token_address, created_at DESC;
	`

	err := m.conn.QueryRowsCtx(ctx, &records, query, tokenAddresses)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*big.Int, len(records))
	for _, record := range records {
		value, _ := new(big.Int).SetString(record.Price, 10)
		result[record.TokenAddress] = value
	}

	return result, nil
}

func (m *customSonexTokenPriceRecordModel) GetCurrentPrice(tokenAddress string) (float64, error) {
	ctx := context.Background()
	var price float64
	query := `
		SELECT price
		FROM sonex_token_price_record
		WHERE token_address = $1
		ORDER BY created_at DESC
		LIMIT 1;
	`
	err := m.conn.QueryRowCtx(ctx, &price, query, tokenAddress)
	return price, err
}

func (m *customSonexTokenPriceRecordModel) GetCurrentPrices(tokenAddresses []string) (map[string]*TokenPriceRecord, error) {
	ctx := context.Background()

	var records []*TokenPriceRecord
	query := fmt.Sprintf("SELECT DISTINCT ON (token_address) %s FROM %s WHERE token_address = ANY($1) ORDER BY token_address, created_at DESC", sonexTokenPriceRecordRows, m.table)

	err := m.conn.QueryRowsCtx(ctx, &records, query, tokenAddresses)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*TokenPriceRecord, len(records))
	for _, record := range records {
		result[record.TokenAddress] = record
	}

	return result, nil
}

func (m *customSonexTokenPriceRecordModel) GetLastHourPrice(tokenAddress string) (float64, error) {
	ctx := context.Background()
	var price float64
	query := `
		SELECT price
		FROM sonex_token_price_record
		WHERE token_address = $1
		  AND created_at <= NOW() - INTERVAL '1 hour'
		ORDER BY created_at DESC
		LIMIT 1;
	`
	err := m.conn.QueryRowCtx(ctx, &price, query, tokenAddress)
	return price, err
}

func (m *customSonexTokenPriceRecordModel) GetYesterdayPrice(tokenAddress string) (float64, error) {
	ctx := context.Background()
	var price float64
	query := `
		SELECT price
		FROM sonex_token_price_record
		WHERE token_address = $1
		  AND created_at <= NOW() - INTERVAL '1 day'
		ORDER BY created_at DESC
		LIMIT 1;
	`
	err := m.conn.QueryRowCtx(ctx, &price, query, tokenAddress)
	return price, err
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
