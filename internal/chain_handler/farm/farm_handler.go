package farm

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"

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

	onIncentiveCreated  []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerIncentiveCreated)
	onIncentiveEnded    []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerIncentiveEnded)
	onTokenStaked       []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenStaked)
	onTokenUnstaked     []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenUnstaked)
	onRewardClaimed     []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerRewardClaimed)
	onRewardRateUpdated []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerRewardRateUpdated)
	onTokenDeposited    []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenDeposited)
	onPaused            []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerPaused)
	onUnpaused          []func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerUnpaused)
}

func NewFarmHandler(dataSource string, chain config.ChainConfig, config config.FarmConfig) (*FarmHandler, error) {
	return &FarmHandler{
		dataSource: dataSource,
		chain:      &chain,
		config:     &config,
	}, nil
}

func (handler *FarmHandler) OnIncentiveCreated(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerIncentiveCreated)) {
	handler.onIncentiveCreated = append(handler.onIncentiveCreated, callback)
}

func (handler *FarmHandler) OnIncentiveEnded(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerIncentiveEnded)) {
	handler.onIncentiveEnded = append(handler.onIncentiveEnded, callback)
}

func (handler *FarmHandler) OnTokenStaked(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenStaked)) {
	handler.onTokenStaked = append(handler.onTokenStaked, callback)
}

func (handler *FarmHandler) OnTokenUnstaked(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenUnstaked)) {
	handler.onTokenUnstaked = append(handler.onTokenUnstaked, callback)
}

func (handler *FarmHandler) OnRewardClaimed(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerRewardClaimed)) {
	handler.onRewardClaimed = append(handler.onRewardClaimed, callback)
}

func (handler *FarmHandler) OnRewardRateUpdated(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerRewardRateUpdated)) {
	handler.onRewardRateUpdated = append(handler.onRewardRateUpdated, callback)
}

func (handler *FarmHandler) OnTokenDeposited(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenDeposited)) {
	handler.onTokenDeposited = append(handler.onTokenDeposited, callback)
}

func (handler *FarmHandler) OnPaused(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerPaused)) {
	handler.onPaused = append(handler.onPaused, callback)
}

func (handler *FarmHandler) OnUnpaused(callback func(chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerUnpaused)) {
	handler.onUnpaused = append(handler.onUnpaused, callback)
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

	// IncentiveCreated event
	incentiveCreatedLogs, err := farmContract.FilterIncentiveCreated(opt, nil, nil, nil)
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

	// IncentiveEnded event
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

	// TokenStaked event
	tokenStakedLogs, err := farmContract.FilterTokenStaked(opt, nil, nil, nil)
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

	// TokenUnstaked event
	tokenUnstakedLogs, err := farmContract.FilterTokenUnstaked(opt, nil, nil, nil)
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

	// RewardClaimed event
	rewardClaimedLogs, err := farmContract.FilterRewardClaimed(opt, nil, nil)
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

	// RewardRateUpdated event
	rewardRateUpdatedLogs, err := farmContract.FilterRewardRateUpdated(opt, nil)
	if err != nil {
		return err, false
	}
	rewardRateUpdatedEvents := []*farm.UniswapV3StakerRewardRateUpdated{}
	for rewardRateUpdatedLogs.Next() {
		if err = rewardRateUpdatedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(rewardRateUpdatedLogs.Event.Raw.BlockNumber)] = struct{}{}
		rewardRateUpdatedEvents = append(rewardRateUpdatedEvents, rewardRateUpdatedLogs.Event)
	}

	// TokenDeposited event
	tokenDepositedLogs, err := farmContract.FilterTokenDeposited(opt, nil, nil)
	if err != nil {
		return err, false
	}
	tokenDepositedEvents := []*farm.UniswapV3StakerTokenDeposited{}
	for tokenDepositedLogs.Next() {
		if err = tokenDepositedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(tokenDepositedLogs.Event.Raw.BlockNumber)] = struct{}{}
		tokenDepositedEvents = append(tokenDepositedEvents, tokenDepositedLogs.Event)
	}

	// Paused event
	pausedLogs, err := farmContract.FilterPaused(opt, nil)
	if err != nil {
		return err, false
	}
	pausedEvents := []*farm.UniswapV3StakerPaused{}
	for pausedLogs.Next() {
		if err = pausedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(pausedLogs.Event.Raw.BlockNumber)] = struct{}{}
		pausedEvents = append(pausedEvents, pausedLogs.Event)
	}

	// Unpaused event
	unpausedLogs, err := farmContract.FilterUnpaused(opt, nil)
	if err != nil {
		return err, false
	}
	unpausedEvents := []*farm.UniswapV3StakerUnpaused{}
	for unpausedLogs.Next() {
		if err = unpausedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(unpausedLogs.Event.Raw.BlockNumber)] = struct{}{}
		unpausedEvents = append(unpausedEvents, unpausedLogs.Event)
	}

	var numbers []int64
	linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
	transactions, err := nodeCaller.BatchGetTranscationsByNumber(opt.Context, numbers)
	if err != nil {
		return err, false
	}

	conn := repository.NewConn(handler.dataSource)
	err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
		eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmIncentiveEndedModel := repository.NewEventFarmIncentiveEndedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmTokenUnstakedModel := repository.NewEventFarmTokenUnstakedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmRewardClaimedModel := repository.NewEventFarmRewardClaimedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmRewardRateUpdatedModel := repository.NewEventFarmRewardRateUpdatedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmTokenDepositedModel := repository.NewEventFarmTokenDepositedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmPausedModel := repository.NewEventFarmPausedModel(sqlx.NewSqlConnFromSession(s))
		eventFarmUnpausedModel := repository.NewEventFarmUnpausedModel(sqlx.NewSqlConnFromSession(s))

		callableIncentiveCreatedEvents := []*farm.UniswapV3StakerIncentiveCreated{}
		// IncentiveCreated Event
		for _, event := range incentiveCreatedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventRewardToken := event.RewardToken.String()
			eventPool := event.Pool.String()
			eventStartTime := event.StartTime.Int64()
			eventWeight := event.Weight.Int64()
			eventIncentiveId := event.IncentiveId.Int64()

			// If handled, ignore
			_, err := eventFarmIncentiveCreatedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmIncentiveCreatedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventRewardToken, eventPool, eventStartTime, eventWeight, eventIncentiveId,
			)
			if err != nil {
				return err
			}

			callableIncentiveCreatedEvents = append(callableIncentiveCreatedEvents, event)
		}

		callableIncentiveEndedEvents := []*farm.UniswapV3StakerIncentiveEnded{}
		// IncentiveEnded Event
		for _, event := range incentiveEndedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventIncentiveId := event.IncentiveId.Int64()
			eventRefund := event.Refund.String()

			// If handled, ignore
			_, err := eventFarmIncentiveEndedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmIncentiveEndedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventIncentiveId, eventRefund,
			)
			if err != nil {
				return err
			}

			callableIncentiveEndedEvents = append(callableIncentiveEndedEvents, event)
		}

		callableTokenStakedEvents := []*farm.UniswapV3StakerTokenStaked{}
		// TokenStaked Event
		for _, event := range tokenStakedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventUser := event.User.Hex()
			eventTokenId := event.TokenId.Int64()
			eventIncentiveId := event.IncentiveId.Int64()

			// If handled, ignore
			_, err := eventFarmTokenStakedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmTokenStakedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventUser, eventTokenId, eventIncentiveId,
			)
			if err != nil {
				return err
			}

			callableTokenStakedEvents = append(callableTokenStakedEvents, event)
		}

		callableTokenUnstakedEvents := []*farm.UniswapV3StakerTokenUnstaked{}
		// TokenUnstaked Event
		for _, event := range tokenUnstakedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventUser := event.User.Hex()
			eventTokenId := event.TokenId.Int64()
			eventIncentiveId := event.IncentiveId.Int64()

			// If handled, ignore
			_, err := eventFarmTokenUnstakedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmTokenUnstakedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventUser, eventTokenId, eventIncentiveId,
			)
			if err != nil {
				return err
			}

			callableTokenUnstakedEvents = append(callableTokenUnstakedEvents, event)
		}

		callableRewardClaimedEvents := []*farm.UniswapV3StakerRewardClaimed{}
		// RewardClaimed Event
		for _, event := range rewardClaimedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventUser := event.User.Hex()
			eventIncentiveId := event.IncentiveId.Int64()
			eventAmount := event.Amount.String()

			// If handled, ignore
			_, err := eventFarmRewardClaimedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmRewardClaimedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventUser, eventIncentiveId, eventAmount,
			)
			if err != nil {
				return err
			}

			callableRewardClaimedEvents = append(callableRewardClaimedEvents, event)
		}

		callableRewardRateUpdatedEvents := []*farm.UniswapV3StakerRewardRateUpdated{}
		// RewardRateUpdated Event
		for _, event := range rewardRateUpdatedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventRewardToken := event.RewardToken.Hex()
			eventOldRate := event.OldRate.String()
			eventNewRate := event.NewRate.String()
			eventTimestamp := event.Timestamp.Int64()

			// If handled, ignore
			_, err := eventFarmRewardRateUpdatedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmRewardRateUpdatedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventRewardToken, eventOldRate, eventNewRate, eventTimestamp,
			)
			if err != nil {
				return err
			}

			callableRewardRateUpdatedEvents = append(callableRewardRateUpdatedEvents, event)
		}

		callableTokenDepositedEvents := []*farm.UniswapV3StakerTokenDeposited{}
		// TokenDeposited Event
		for _, event := range tokenDepositedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventTokenId := event.TokenId.Int64()
			eventOwner := event.Owner.Hex()

			// If handled, ignore
			_, err := eventFarmTokenDepositedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmTokenDepositedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventTokenId, eventOwner,
			)
			if err != nil {
				return err
			}

			callableTokenDepositedEvents = append(callableTokenDepositedEvents, event)
		}

		callablePausedEvents := []*farm.UniswapV3StakerPaused{}
		// Paused Event
		for _, event := range pausedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventAccount := event.Account.Hex()

			// If handled, ignore
			_, err := eventFarmPausedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmPausedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventAccount,
			)
			if err != nil {
				return err
			}

			callablePausedEvents = append(callablePausedEvents, event)
		}

		callableUnpausedEvents := []*farm.UniswapV3StakerUnpaused{}
		// Unpaused Event
		for _, event := range unpausedEvents {
			blockNumber := int64(event.Raw.BlockNumber)
			tx := transactions[event.Raw.TxHash.Hex()]
			txIndex := int64(event.Raw.Index)
			txHash := tx.Tx().Hash().Hex()
			txAt := tx.TxAt
			from := tx.From.Hex()
			to := tx.Tx().To().Hex()
			value := tx.Tx().Value().String()
			eventAccount := event.Account.Hex()

			// If handled, ignore
			_, err := eventFarmUnpausedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
			switch err {
			case repository.ErrNotFound:
			case nil:
				continue
			default:
				return err
			}

			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = eventFarmUnpausedModel.New(
				ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value,
				eventAccount,
			)
			if err != nil {
				return err
			}

			callableUnpausedEvents = append(callableUnpausedEvents, event)
		}

		// Callback
		for _, event := range callableIncentiveCreatedEvents {
			for _, callback := range handler.onIncentiveCreated {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableIncentiveEndedEvents {
			for _, callback := range handler.onIncentiveEnded {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableTokenStakedEvents {
			for _, callback := range handler.onTokenStaked {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableTokenUnstakedEvents {
			for _, callback := range handler.onTokenUnstaked {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableRewardClaimedEvents {
			for _, callback := range handler.onRewardClaimed {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableRewardRateUpdatedEvents {
			for _, callback := range handler.onRewardRateUpdated {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableTokenDepositedEvents {
			for _, callback := range handler.onTokenDeposited {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callablePausedEvents {
			for _, callback := range handler.onPaused {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range callableUnpausedEvents {
			for _, callback := range handler.onUnpaused {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}

		return nil
	})
	if err != nil {
		return err, false
	}

	return nil, false
}
