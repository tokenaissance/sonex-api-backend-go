package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SonexPoolTvlRecordModel = (*customSonexPoolTvlRecordModel)(nil)

type (
	// SonexPoolTvlRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexPoolTvlRecordModel.
	SonexPoolTvlRecordModel interface {
		sonexPoolTvlRecordModel
		withSession(session sqlx.Session) SonexPoolTvlRecordModel
		UpdateOrInsertTvlByPoolAddress(tvlUsd string, poolAddress string, usdDecimals, chainId int64, valueOfType int64) (int64, error)
	}

	customSonexPoolTvlRecordModel struct {
		*defaultSonexPoolTvlRecordModel
	}
)

// 不同的 tvl 代表的值
var ValueOfTypeTVL = int64(1000)
var ValueOfTypeYield = int64(1001)
var ValueOfTypeVolume = int64(1002)

// NewSonexPoolTvlRecordModel returns a model for the database table.
func NewSonexPoolTvlRecordModel(conn sqlx.SqlConn) SonexPoolTvlRecordModel {
	return &customSonexPoolTvlRecordModel{
		defaultSonexPoolTvlRecordModel: newSonexPoolTvlRecordModel(conn),
	}
}

func (m *customSonexPoolTvlRecordModel) withSession(session sqlx.Session) SonexPoolTvlRecordModel {
	return NewSonexPoolTvlRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSonexPoolTvlRecordModel) UpdateOrInsertTvlByPoolAddress(tvlUsd string, poolAddress string, usdDecimals, chainId int64, valueOfType int64) (int64, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`
	SELECT * from %s WHERE pool_address = $1 AND value_of_type = $2 limit 1
	`, m.table)

	var record SonexPoolTvlRecord
	err := m.conn.QueryRowCtx(ctx, &record, query, poolAddress, valueOfType)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logx.Errorf("QueryRowCtx:%v  ", err)
		return 0, err
	} else if errors.Is(err, sqlx.ErrNotFound) {
		data := &SonexPoolTvlRecord{
			PoolAddress: poolAddress,
			Tvl:         tvlUsd,
			Created:     NowUTC(),
			UsdDecimals: usdDecimals,
			ChainId:     chainId,
			Updated:     NowUTC(),
		}

		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6,$7)", m.table, sonexPoolTvlRecordRowsExpectAutoSet)
		ret, err := m.conn.ExecCtx(ctx, query, data.PoolAddress, data.Tvl, data.Created, data.UsdDecimals, data.ChainId, data.Updated, valueOfType)

		if err != nil {
			logx.Errorf("Insert   error :%v ", err)
		} else {
			return ret.RowsAffected()
		}
	} else if err == nil {
		query = fmt.Sprintf(`
UPDATE % s 
SET 
tvl =  $1,
updated = $2 
WHERE
	ID = $3
		`, m.table)
		retModified, err := m.conn.ExecCtx(ctx, query, tvlUsd, NowUTC(), record.Id)
		if err != nil {
			logx.Errorf("UpdateTvlByPoolAddress,%v ", err)
			return 0, err
		} else {
			return retModified.RowsAffected()
		}

	}
	return 0, nil
}
