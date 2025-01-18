package repository

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/module/sonex"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/zeromicro/go-zero/core/logx"
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
		QueryLatestHistoryPrices() (map[string]*big.Int, error)
		QueryLatestHistoryPricesByAddresses(addresses []string) (map[string]*big.Int, error)
		QueryLatestHistoryPricesByAddressesInterval(addresses []string, interval string) (map[string]*big.Int, error)
		GetCurrentPrice(tokenAddress string) (float64, error)
		GetCurrentPrices(tokenAddresses []string) (map[string]*TokenPriceRecord, error)
		GetCurrentPricesAndTokenDecimals(tokenAddresses []string) ([]*TokenPriceRecordWithTokenDecimals, error)
		GetLastHourPrice(ctx context.Context, tokenAddress string) (map[string]*sonex.TokenPrice, error)
		GetYesterdayPrice(ctx context.Context, tokenAddress string) (map[string]*sonex.TokenPrice, error)
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
	UsdDecimals  int    `db:"usd_decimals"`
}

type MiniTokenPriceRecord struct {
	TokenAddress string `db:"token_address"`
	Price        string `db:"price"`
	CreatedAt    string `db:"created_at"`
}

// downsampleRecords performs time-based downsampling of price records
func (m *customSonexTokenPriceRecordModel) downsampleRecords(records []*TokenPriceRecord, interval string) []*TokenPriceRecord {
	if len(records) == 0 {
		return records
	}

	// Define desired number of K-lines
	var targetKLines int
	switch interval {
	case "1h":
		targetKLines = 60 // 60 K-lines per hour
	case "1d":
		targetKLines = 144 // 144 K-lines per day
	case "1w":
		targetKLines = 168 // 168 K-lines per week
	case "1m":
		targetKLines = 780 // 780 K-lines per month
	case "1y":
		targetKLines = 365 // 365 K-lines per year
	default:
		return records
	}

	// If we have fewer records than target K-lines, return all records
	if len(records) <= targetKLines {
		return records
	}

	// Calculate time range
	firstTime, _ := gtime.StrToTime(records[0].CreatedAt)
	lastTime, _ := gtime.StrToTime(records[len(records)-1].CreatedAt)
	totalDuration := lastTime.Sub(firstTime)

	// Calculate ideal time interval between samples
	timeInterval := totalDuration / time.Duration(targetKLines-1)

	var result []*TokenPriceRecord
	result = append(result, records[0]) // Always include first record

	nextTargetTime := firstTime.Add(timeInterval)
	lastAddedTime := firstTime

	// Find closest records to each target time
	for i := 1; i < len(records)-1; i++ {
		currentTime, err := gtime.StrToTime(records[i].CreatedAt)
		if err != nil {
			logx.Errorf("[downsampleRecords]: time convert err: %v", err)
			continue
		}

		if currentTime.After(nextTargetTime) || currentTime.Equal(nextTargetTime) {
			// Find closest record to target time
			prev := gtime.New(records[i-1].CreatedAt)
			curr := currentTime
			if curr.Sub(nextTargetTime) > nextTargetTime.Sub(prev) {
				result = append(result, records[i-1])
				lastAddedTime = prev
			} else {
				result = append(result, records[i])
				lastAddedTime = curr
			}
			nextTargetTime = firstTime.Add(timeInterval * time.Duration(len(result)))
		}
	}

	// Always include last record if it's not too close to the last added record
	lastRecord := records[len(records)-1]
	lastRecordTime, _ := gtime.StrToTime(lastRecord.CreatedAt)
	if lastRecordTime.Sub(lastAddedTime) > timeInterval/2 {
		result = append(result, lastRecord)
	}

	return result
}

// minutesDifference calculates the difference in minutes between two times
func minutesDifference(t1, t2 *gtime.Time) int {
	minutes := t2.Sub(t1) / time.Minute
	return int(minutes)
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
		// downsample
		return m.downsampleRecords(resp, interval), nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customSonexTokenPriceRecordModel) QueryLatestHistoryPrices() (map[string]*big.Int, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`
		SELECT DISTINCT ON (token_address) token_address, price, created_at
		FROM %s
		ORDER BY token_address, created_at DESC;
	`, m.table)

	var records []MiniTokenPriceRecord
	err := m.conn.QueryRowsCtx(ctx, &records, query)
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

func (m *customSonexTokenPriceRecordModel) QueryLatestHistoryPricesByAddresses(tokenAddresses []string) (map[string]*big.Int, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`
		SELECT DISTINCT ON (token_address) token_address, price, created_at
		FROM %s
		WHERE token_address = ANY($1)
		ORDER BY token_address, created_at DESC;
	`, m.table)

	var records []MiniTokenPriceRecord
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

func (m *customSonexTokenPriceRecordModel) QueryLatestHistoryPricesByAddressesInterval(tokenAddresses []string, interval string) (map[string]*big.Int, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`
		SELECT DISTINCT ON (token_address) token_address, price, created_at
		FROM %s
		WHERE 
			token_address = ANY($1)
			AND created_at <= NOW() - INTERVAL '%s'
	`, m.table, getIntervalDuration(interval))

	var records []MiniTokenPriceRecord
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
	if err != nil {
		return 0, err
	}
	return price, nil
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

type TokenPriceRecordWithTokenDecimals struct {
	Id            int64  `db:"id"`
	TokenSymbol   string `db:"token_symbol"`
	TokenAddress  string `db:"token_address"`
	Price         string `db:"price"`
	CreatedAt     string `db:"created_at"`
	UsdDecimals   int    `db:"usd_decimals"`
	TokenDecimals int    `db:"token_decimals"`
}

// 获取指定地址对应的token的decimals
// 方便后续计算
func (m *customSonexTokenPriceRecordModel) GetCurrentPricesAndTokenDecimals(tokenAddresses []string) ([]*TokenPriceRecordWithTokenDecimals, error) {
	ctx := context.Background()

	var records []*TokenPriceRecordWithTokenDecimals
	query := `
	SELECT DISTINCT ON
	(token_address) s1.*,
	s2.decimals AS token_decimals 
FROM
	sonex_token_price_record s1
	JOIN sonex_token s2 ON s2.address = s1.token_address 
WHERE
	s1.token_address = ANY ($1) 
ORDER BY
	s1.token_address,
	s1.created_at DESC
	`

	err := m.conn.QueryRowsCtx(ctx, &records, query, tokenAddresses)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (m *customSonexTokenPriceRecordModel) GetLastHourPrice(ctx context.Context, tokenAddress string) (map[string]*sonex.TokenPrice, error) {
	var price string
	query := `
		SELECT price
		FROM sonex_token_price_record
		WHERE token_address = $1
		  AND created_at <= NOW() - INTERVAL '1 hour'
		ORDER BY created_at DESC
		LIMIT 1;
	`
	err := m.conn.QueryRowCtx(ctx, &price, query, tokenAddress)
	if err != nil {
		return nil, err
	}
	tp, err := sonex.NewTokenPrice(price, 18)
	if err != nil {
		return nil, err
	}
	return map[string]*sonex.TokenPrice{tokenAddress: tp}, nil
}

func (m *customSonexTokenPriceRecordModel) GetYesterdayPrice(ctx context.Context, tokenAddress string) (map[string]*sonex.TokenPrice, error) {
	var price string
	query := `
		SELECT price
		FROM sonex_token_price_record
		WHERE token_address = $1
		  AND created_at <= NOW() - INTERVAL '1 day'
		ORDER BY created_at DESC
		LIMIT 1;
	`
	err := m.conn.QueryRowCtx(ctx, &price, query, tokenAddress)
	if err != nil {
		return nil, err
	}
	tp, err := sonex.NewTokenPrice(price, 18)
	if err != nil {
		return nil, err
	}
	return map[string]*sonex.TokenPrice{tokenAddress: tp}, nil
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
