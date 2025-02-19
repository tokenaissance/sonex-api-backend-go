package repository

import (
	"context"
	"errors"
	"sonex_api_backend_go/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SonexVolTvlRecordModel = (*customSonexVolTvlRecordModel)(nil)

type (
	// SonexVolTvlRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexVolTvlRecordModel.
	SonexVolTvlRecordModel interface {
		sonexVolTvlRecordModel
		withSession(session sqlx.Session) SonexVolTvlRecordModel
		UpdateOrInsertTvlByTimeStamp(valueUsd string, created_at int64, typeOfUniswap int64, chainId int64, typeOfValue int64) (int64, error)
		GetBytypeOfValueAndChainId(chainId int64, typeOfUniswap int64, typeOfValue int64, numOfItems int64) ([]types.UsdCreatedRecord, error)
	}

	customSonexVolTvlRecordModel struct {
		*defaultSonexVolTvlRecordModel
	}
)

// TypeoOfValue int64     `db:"typeo_of_value"` // 1: volume value; 2:tvl value
var VolumeOfTypeoOfValue = int64(1) //
var TVLTypeoOfValue = int64(2)

var V3OfUniswap = int64(1)
var V2OfUniswap = int64(2)

// NewSonexVolTvlRecordModel returns a model for the database table.
func NewSonexVolTvlRecordModel(conn sqlx.SqlConn) SonexVolTvlRecordModel {
	return &customSonexVolTvlRecordModel{
		defaultSonexVolTvlRecordModel: newSonexVolTvlRecordModel(conn),
	}
}

func (m *customSonexVolTvlRecordModel) withSession(session sqlx.Session) SonexVolTvlRecordModel {
	return NewSonexVolTvlRecordModel(sqlx.NewSqlConnFromSession(session))
}

// 按照创建时间进行检索，每天一条
func (m *customSonexVolTvlRecordModel) UpdateOrInsertTvlByTimeStamp(valueUsd string, created_at int64, typeofUnsiwap int64, chainId int64, typeOfValue int64) (int64, error) {
	ctx := context.Background()

	query := `
	SELECT * from sonex_vol_tvl_record WHERE createtime = $1 And type_of_uniswap = $2 AND chain_id=$3 AND type_of_value = $4 limit 1
	`

	record := &SonexVolTvlRecord{}

	err := m.conn.QueryRowCtx(ctx, record, query, created_at, typeofUnsiwap, chainId, typeOfValue)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logx.Errorf("tvl_vol_total :%v  ", err)
		return 0, err
	} else if errors.Is(err, sqlx.ErrNotFound) {
		data := &SonexVolTvlRecord{
			UsdtValue:     valueUsd,      //   string `db:"usdt_value"`
			ChainId:       chainId,       //   int64  `db:"chain_id"`
			TypeOfValue:   typeOfValue,   //    int64  `db:"type_of_value"`   // 1: volume value; 2:tvl value
			Createtime:    created_at,    //   int64  `db:"created_at"`      // 创建时间，需要拿着创建时间后续做更新的key
			TypeOfUniswap: typeofUnsiwap, //int64  `db:"type_of_uniswap"` // 1: v3, 2: v2
		}

		inserted, err := m.Insert(ctx, data)
		if err != nil {
			logx.Errorf("Insert   error :%v ", err)
		} else {
			return inserted.RowsAffected()
		}
	} else if err == nil {
		if strings.EqualFold(record.UsdtValue, valueUsd) {
			logx.Infof(" tvl_vol_total not changed ,no need to update record:%v ", record)
			return 0, nil
		}
		logx.Infof("tvl_vol_total:%#v ", record)
		query = `
UPDATE sonex_vol_tvl_record 
SET 
usdt_value = $1
WHERE
id = $2
		`
		retModified, err := m.conn.ExecCtx(ctx, query, valueUsd, record.Id)
		if err != nil {
			logx.Errorf("UpdateTvlByPoolAddress,%v ", err)
			return 0, err
		} else {
			return retModified.RowsAffected()
		}

	}
	return 0, nil
}

func (m *customSonexVolTvlRecordModel) GetBytypeOfValueAndChainId(chainId int64, typeOfUniswap int64, typeOfValue int64, numOfItems int64) ([]types.UsdCreatedRecord, error) {
	records := []types.UsdCreatedRecord{}
	ctx := context.Background()
	query := `
	SELECT usdt_value , createtime FROM  sonex_vol_tvl_record WHERE chain_id=$1 AND type_of_uniswap=$2 AND type_of_value = $3 ORDER BY createtime ASC limit $4 
	`

	err := m.conn.QueryRowsCtx(ctx, &records, query, chainId, typeOfUniswap, typeOfValue, numOfItems)
	if err != nil {
		return []types.UsdCreatedRecord{}, err
	}

	return records, nil
}
