package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolPoolDepositedModel = (*customEventXpoolPoolDepositedModel)(nil)

type (
	// EventXpoolPoolDepositedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolPoolDepositedModel.
	EventXpoolPoolDepositedModel interface {
		eventXpoolPoolDepositedModel
		withSession(session sqlx.Session) EventXpoolPoolDepositedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) error
		GetDepositedEventsByUserAddress(ctx context.Context, chainId int64, address string) (map[string][]*EventXpoolPoolDeposited, error)
	}

	customEventXpoolPoolDepositedModel struct {
		*defaultEventXpoolPoolDepositedModel
	}
)

// NewEventXpoolPoolDepositedModel returns a model for the database table.
func NewEventXpoolPoolDepositedModel(conn sqlx.SqlConn) EventXpoolPoolDepositedModel {
	return &customEventXpoolPoolDepositedModel{
		defaultEventXpoolPoolDepositedModel: newEventXpoolPoolDepositedModel(conn),
	}
}

func (m *customEventXpoolPoolDepositedModel) withSession(session sqlx.Session) EventXpoolPoolDepositedModel {
	return NewEventXpoolPoolDepositedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolPoolDepositedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXpoolPoolDepositedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
	return err
}

// GetDepositedEventsByUserAddress 获取指定地址的存款事件
func (m *customEventXpoolPoolDepositedModel) GetDepositedEventsByUserAddress(ctx context.Context, chainId int64, address string) (map[string][]*EventXpoolPoolDeposited, error) {
	var resp []*EventXpoolPoolDeposited
	// 将用户地址转换为以太坊地址格式
	address = common.HexToAddress(address).Hex()
	query := fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1 AND from_ = $2", m.table)
	// 执行查询
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, address)
	if err != nil {
		return nil, err
	}

	// 将事件按池子地址分组
	depositedEvents := make(map[string][]*EventXpoolPoolDeposited)
	for _, event := range resp {
		// 获取池子地址
		poolAddress := common.HexToAddress(event.To).Hex()
		// 将事件按池子地址分组
		poolEvents, ok := depositedEvents[poolAddress]
		if !ok {
			poolEvents = make([]*EventXpoolPoolDeposited, 0)
		}
		poolEvents = append(poolEvents, event)
		depositedEvents[poolAddress] = poolEvents
	}
	return depositedEvents, nil
}
