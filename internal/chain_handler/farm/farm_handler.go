package farm

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmHandler struct {
	dataSource string
	chain      *config.ChainConfig
	config     *config.FarmConfig

	onIncentiveCreated   []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventRewardToken string, eventPool string, eventStartTime int64, eventEndTime int64, eventRefundee string, eventReward string, eventIncentiveId string)
	onTokenStaked        []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTokenId int64, eventIncentiveId string, eventLiquidity string)
	onTokenUnstaked      []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTokenId int64, eventIncentiveId string)
	onRewardClaimed      []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTo string, eventReward string)
	onDepositTransferred []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTokenId int64, eventOldOwner string, eventNewOwner string)
	onIncentiveEnded     []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventIncentiveId string, eventRefund string)
}

func NewFarmHandler(dataSource string, chain config.ChainConfig, config config.FarmConfig) (*FarmHandler, error) {
	return &FarmHandler{
		dataSource: dataSource,
		chain:      &chain,
		config:     &config,
	}, nil
}

func (handler *FarmHandler) OnIncentiveCreated(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventRewardToken string, eventPool string, eventStartTime int64, eventEndTime int64, eventRefundee string, eventReward string, eventIncentiveId string)) {
	handler.onIncentiveCreated = append(handler.onIncentiveCreated, callback)
}

func (handler *FarmHandler) OnTokenStaked(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTokenId int64, eventIncentiveId string, eventLiquidity string)) {
	handler.onTokenStaked = append(handler.onTokenStaked, callback)
}

func (handler *FarmHandler) OnTokenUnstaked(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTokenId int64, eventIncentiveId string)) {
	handler.onTokenUnstaked = append(handler.onTokenUnstaked, callback)
}

func (handler *FarmHandler) OnRewardClaimed(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTo string, eventReward string)) {
	handler.onRewardClaimed = append(handler.onRewardClaimed, callback)
}

func (handler *FarmHandler) OnDepositTransferred(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventTokenId int64, eventOldOwner string, eventNewOwner string)) {
	handler.onDepositTransferred = append(handler.onDepositTransferred, callback)
}

func (handler *FarmHandler) OnIncentiveEnded(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventIncentiveId string, eventRefund string)) {
	handler.onIncentiveEnded = append(handler.onIncentiveEnded, callback)
}

func (handler *FarmHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool) {
	nodeCaller, err := node_caller.NewNodeCaller(*handler.chain)
	if err != nil {
		return err, false
	}

	farmContract, err := farm.NewUniswapV3Staker(common.HexToAddress(handler.config.Staker), client)
	if err != nil {
		return err, false
	}

	numberSet := map[int64]struct{}{}

	// 获取IncentiveCreated事件
	incentiveCreatedLogs, err := farmContract.FilterIncentiveCreated(opt, nil, nil)
	if err != nil {
		return err, false
	}
	incentiveCreatedEvents := []*farm.UniswapV3StakerIncentiveCreated{}
	for incentiveCreatedLogs.Next() {
		if err = incentiveCreatedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(incentiveCreatedLogs.Event.Raw.BlockNumber)] = struct{}{}
		incentiveCreatedEvents = append(incentiveCreatedEvents, incentiveCreatedLogs.Event)
	}

	// 获取TokenStaked事件
	tokenStakedLogs, err := farmContract.FilterTokenStaked(opt, nil, nil)
	if err != nil {
		return err, false
	}
	tokenStakedEvents := []*farm.UniswapV3StakerTokenStaked{}
	for tokenStakedLogs.Next() {
		if err = tokenStakedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(tokenStakedLogs.Event.Raw.BlockNumber)] = struct{}{}
		tokenStakedEvents = append(tokenStakedEvents, tokenStakedLogs.Event)
	}

	// 获取TokenUnstaked事件
	tokenUnstakedLogs, err := farmContract.FilterTokenUnstaked(opt, nil, nil)
	if err != nil {
		return err, false
	}
	tokenUnstakedEvents := []*farm.UniswapV3StakerTokenUnstaked{}
	for tokenUnstakedLogs.Next() {
		if err = tokenUnstakedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(tokenUnstakedLogs.Event.Raw.BlockNumber)] = struct{}{}
		tokenUnstakedEvents = append(tokenUnstakedEvents, tokenUnstakedLogs.Event)
	}

	// 获取RewardClaimed事件
	rewardClaimedLogs, err := farmContract.FilterRewardClaimed(opt, nil)
	if err != nil {
		return err, false
	}
	rewardClaimedEvents := []*farm.UniswapV3StakerRewardClaimed{}
	for rewardClaimedLogs.Next() {
		if err = rewardClaimedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(rewardClaimedLogs.Event.Raw.BlockNumber)] = struct{}{}
		rewardClaimedEvents = append(rewardClaimedEvents, rewardClaimedLogs.Event)
	}

	// 获取DepositTransferred事件
	depositTransferredLogs, err := farmContract.FilterDepositTransferred(opt, nil, nil, nil)
	if err != nil {
		return err, false
	}
	depositTransferredEvents := []*farm.UniswapV3StakerDepositTransferred{}
	for depositTransferredLogs.Next() {
		if err = depositTransferredLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(depositTransferredLogs.Event.Raw.BlockNumber)] = struct{}{}
		depositTransferredEvents = append(depositTransferredEvents, depositTransferredLogs.Event)
	}

	// 获取IncentiveEnded事件
	incentiveEndedLogs, err := farmContract.FilterIncentiveEnded(opt, nil)
	if err != nil {
		return err, false
	}
	incentiveEndedEvents := []*farm.UniswapV3StakerIncentiveEnded{}
	for incentiveEndedLogs.Next() {
		if err = incentiveEndedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(incentiveEndedLogs.Event.Raw.BlockNumber)] = struct{}{}
		incentiveEndedEvents = append(incentiveEndedEvents, incentiveEndedLogs.Event)
	}

	// 获取区块交易信息
	var numbers []int64
	linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
	transactions, err := nodeCaller.GetTranscations(opt.Context, numbers)
	if err != nil {
		return err, false
	}

	conn := repository.NewConn(handler.dataSource)
	err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
		eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmTokenUnstakedModel := repository.NewEventFarmTokenUnstakedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmRewardClaimedModel := repository.NewEventFarmRewardClaimedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmDepositTransferredModel := repository.NewEventFarmDepositTransferredModel(sqlx.NewSqlConnFromSession(s))
		eventFarmIncentiveEndedModel := repository.NewEventFarmIncentiveEndedModel(sqlx.NewSqlConnFromSession(s))

		// IncentiveCreated Event
		for _, event := range incentiveCreatedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx.Hash().Hex()
			txAt := tx.Time
			from := tx.From.Hex()
			to := tx.Tx.To().Hex()
			value := tx.Tx.Value().String()

			// If handled, ignore
			_, err := eventFarmIncentiveCreatedModel.FindOneByTxHashChainId(ctx, txHash, handler.chain.ChainId)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			eventData := &repository.EventFarmIncentiveCreated{
				ChainId:          handler.chain.ChainId,
				BlockNumber:      blockNumber,
				TxIndex:          txIndex,
				TxHash:           txHash,
				TxAt:             txAt,
				From:             from,
				To:               to,
				Value:            value,
				EventRewardToken: event.RewardToken.String(),
				EventPool:        event.Pool.String(),
				EventStartTime:   event.StartTime.Int64(),
				EventEndTime:     event.EndTime.Int64(),
				EventRefundee:    event.Refundee.String(),
				EventReward:      event.Reward.String(),
				EventIncentiveId: common.BytesToHash(event.IncentiveId[:]).Hex(),
			}
			_, err = eventFarmIncentiveCreatedModel.Insert(ctx, eventData)
			if err != nil {
				return err
			}

			// Callback
			for _, callback := range handler.onIncentiveCreated {
				go callback(
					handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
					eventData.EventRewardToken, eventData.EventPool, eventData.EventStartTime,
					eventData.EventEndTime, eventData.EventRefundee, eventData.EventReward, eventData.EventIncentiveId,
				)
			}
		}

		// TokenStaked Event
		for _, event := range tokenStakedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx.Hash().Hex()
			txAt := tx.Time
			from := tx.From.Hex()
			to := tx.Tx.To().Hex()
			value := tx.Tx.Value().String()

			// If handled, ignore
			_, err := eventFarmTokenStakedModel.FindOneByTxHashChainId(ctx, txHash, handler.chain.ChainId)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			eventData := &repository.EventFarmTokenStaked{
				ChainId:          handler.chain.ChainId,
				BlockNumber:      blockNumber,
				TxIndex:          txIndex,
				TxHash:           txHash,
				TxAt:             txAt,
				From:             from,
				To:               to,
				Value:            value,
				EventTokenId:     event.TokenId.Int64(),
				EventIncentiveId: common.BytesToHash(event.IncentiveId[:]).Hex(),
				EventLiquidity:   event.Liquidity.String(),
			}
			_, err = eventFarmTokenStakedModel.Insert(ctx, eventData)
			if err != nil {
				return err
			}

			// Callback
			for _, callback := range handler.onTokenStaked {
				go callback(
					handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
					eventData.EventTokenId, eventData.EventIncentiveId, eventData.EventLiquidity,
				)
			}
		}

		// TokenUnstaked Event
		for _, event := range tokenUnstakedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx.Hash().Hex()
			txAt := tx.Time
			from := tx.From.Hex()
			to := tx.Tx.To().Hex()
			value := tx.Tx.Value().String()

			// If handled, ignore
			_, err := eventFarmTokenUnstakedModel.FindOneByTxHashChainId(ctx, txHash, handler.chain.ChainId)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			eventData := &repository.EventFarmTokenUnstaked{
				ChainId:          handler.chain.ChainId,
				BlockNumber:      blockNumber,
				TxIndex:          txIndex,
				TxHash:           txHash,
				TxAt:             txAt,
				From:             from,
				To:               to,
				Value:            value,
				EventTokenId:     event.TokenId.Int64(),
				EventIncentiveId: common.BytesToHash(event.IncentiveId[:]).Hex(),
			}
			_, err = eventFarmTokenUnstakedModel.Insert(ctx, eventData)
			if err != nil {
				return err
			}

			// Callback
			for _, callback := range handler.onTokenUnstaked {
				go callback(
					handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
					eventData.EventTokenId, eventData.EventIncentiveId,
				)
			}
		}

		// RewardClaimed Event
		for _, event := range rewardClaimedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx.Hash().Hex()
			txAt := tx.Time
			from := tx.From.Hex()
			to := tx.Tx.To().Hex()
			value := tx.Tx.Value().String()

			// If handled, ignore
			_, err := eventFarmRewardClaimedModel.FindOneByTxHashChainId(ctx, txHash, handler.chain.ChainId)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			eventData := &repository.EventFarmRewardClaimed{
				ChainId:     handler.chain.ChainId,
				BlockNumber: blockNumber,
				TxIndex:     txIndex,
				TxHash:      txHash,
				TxAt:        txAt,
				From:        from,
				To:          to,
				Value:       value,
				EventTo:     event.To.String(),
				EventReward: event.Reward.String(),
			}
			_, err = eventFarmRewardClaimedModel.Insert(ctx, eventData)
			if err != nil {
				return err
			}

			// Callback
			for _, callback := range handler.onRewardClaimed {
				go callback(
					handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
					eventData.EventTo, eventData.EventReward,
				)
			}
		}

		// DepositTransferred Event
		for _, event := range depositTransferredEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx.Hash().Hex()
			txAt := tx.Time
			from := tx.From.Hex()
			to := tx.Tx.To().Hex()
			value := tx.Tx.Value().String()

			// If handled, ignore
			_, err := eventFarmDepositTransferredModel.FindOneByTxHashChainId(ctx, txHash, handler.chain.ChainId)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			eventData := &repository.EventFarmDepositTransferred{
				ChainId:       handler.chain.ChainId,
				BlockNumber:   blockNumber,
				TxIndex:       txIndex,
				TxHash:        txHash,
				TxAt:          txAt,
				From:          from,
				To:            to,
				Value:         value,
				EventTokenId:  event.TokenId.Int64(),
				EventOldOwner: event.OldOwner.Hex(),
				EventNewOwner: event.NewOwner.Hex(),
			}
			_, err = eventFarmDepositTransferredModel.Insert(ctx, eventData)
			if err != nil {
				return err
			}

			// Callback
			for _, callback := range handler.onDepositTransferred {
				go callback(
					handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
					eventData.EventTokenId, eventData.EventOldOwner, eventData.EventNewOwner,
				)
			}
		}

		// IncentiveEnded Event
		for _, event := range incentiveEndedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx.Hash().Hex()
			txAt := tx.Time
			from := tx.From.Hex()
			to := tx.Tx.To().Hex()
			value := tx.Tx.Value().String()

			// If handled, ignore
			_, err := eventFarmIncentiveEndedModel.FindOneByTxHashChainId(ctx, txHash, handler.chain.ChainId)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			eventData := &repository.EventFarmIncentiveEnded{
				ChainId:          handler.chain.ChainId,
				BlockNumber:      blockNumber,
				TxIndex:          txIndex,
				TxHash:           txHash,
				TxAt:             txAt,
				From:             from,
				To:               to,
				Value:            value,
				EventIncentiveId: common.BytesToHash(event.IncentiveId[:]).Hex(),
				EventRefund:      event.Refund.String(),
			}
			_, err = eventFarmIncentiveEndedModel.Insert(ctx, eventData)
			if err != nil {
				return err
			}

			// Callback
			for _, callback := range handler.onIncentiveEnded {
				go callback(
					handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
					eventData.EventIncentiveId, eventData.EventRefund,
				)
			}
		}

		return nil
	})
	if err != nil {
		return err, false
	}

	return nil, false
}
