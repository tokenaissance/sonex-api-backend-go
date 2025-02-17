package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolPoolWithdrawedModel = (*customEventXpoolPoolWithdrawedModel)(nil)

type (
	// EventXpoolPoolWithdrawedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolPoolWithdrawedModel.
	EventXpoolPoolWithdrawedModel interface {
		eventXpoolPoolWithdrawedModel
		withSession(session sqlx.Session) EventXpoolPoolWithdrawedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, withdrawed string, claimed string, height int64) error
		GetWithdrawedEvents(ctx context.Context, chainId int64, address string) (map[string][]*EventXpoolPoolWithdrawed, error)
	}

	customEventXpoolPoolWithdrawedModel struct {
		*defaultEventXpoolPoolWithdrawedModel
	}
)

// NewEventXpoolPoolWithdrawedModel returns a model for the database table.
func NewEventXpoolPoolWithdrawedModel(conn sqlx.SqlConn) EventXpoolPoolWithdrawedModel {
	return &customEventXpoolPoolWithdrawedModel{
		defaultEventXpoolPoolWithdrawedModel: newEventXpoolPoolWithdrawedModel(conn),
	}
}

func (m *customEventXpoolPoolWithdrawedModel) withSession(session sqlx.Session) EventXpoolPoolWithdrawedModel {
	return NewEventXpoolPoolWithdrawedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolPoolWithdrawedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, withdrawed string, claimed string, height int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventXpoolPoolWithdrawedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, withdrawed, claimed, height)
	return err
}

// GetWithdrawedEvents 获取指定地址的提款事件
func (m *customEventXpoolPoolWithdrawedModel) GetWithdrawedEvents(ctx context.Context, chainId int64, address string) (map[string][]*EventXpoolPoolWithdrawed, error) {
	var resp []*EventXpoolPoolWithdrawed
	// 将用户地址转换为以太坊地址格式
	address = common.HexToAddress(address).Hex()
	query := fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1 AND to_ = $2", m.table)
	// 执行查询
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, address)
	if err != nil {
		return nil, err
	}

	withdrawedEvents := make(map[string][]*EventXpoolPoolWithdrawed)
	for _, event := range resp {
		// 获取池子地址
		poolAddress := common.HexToAddress(event.From).Hex()
		// 将事件按池子地址分组
		poolEvents, ok := withdrawedEvents[poolAddress]
		if !ok {
			poolEvents = make([]*EventXpoolPoolWithdrawed, 0)
		}
		poolEvents = append(poolEvents, event)
		withdrawedEvents[poolAddress] = poolEvents
	}
	return withdrawedEvents, nil
}
