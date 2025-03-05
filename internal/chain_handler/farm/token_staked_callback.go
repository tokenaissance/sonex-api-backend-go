package farm

import (
	"context"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

func TokenStakedCallback(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenStaked) {
	if !repository.BadgeIsInCampaignTime(tx.TxAt) {
		return
	}

	sonexConf := conf.GetSonexConfig(chainId)
	if sonexConf == nil {
		logx.Errorf("SwapCallback error: %s", "Miss sonex configure")
		return
	}

	conn := repository.NewConn(conf.PostgreSQL.DataSource())

	sonexBadgeStakedPosition := repository.NewSonexBadgeStakedPositionModel(conn)
	caller, err := sonex.NewStaticCaller(conf)
	if err != nil {
		logx.Errorf("NewStaticCaller :%v ", err)
		return
	}
	logx.Infof("tokenId :%v ", event.TokenId.Int64())

	positions, err := caller.BatchGetPositions(
		context.Background(),
		sonexConf.ChainId,
		common.HexToAddress(sonexConf.PositionManager),
		[]*big.Int{event.TokenId},
	)
	if err != nil {
		logx.Errorf("BatchGetPositions ERR :%v ", err)
		return
	}
	if len(positions) != 1 {
		logx.Errorf("BatchGetPositions: unable to get token(id=%v)", event.TokenId)
		return
	}
	totalBalanceUSD, _, _, err := global.SummarizePositions(conf, context.Background(), chainId, positions)
	if err != nil {
		logx.Errorf("err:%v ", err)
		return
	}
	logx.Infof("\n tokenid:%v valueusd :%v ", event.TokenId, totalBalanceUSD)

	// 把 tokenId 写入的值写入表 sonex_badge_staked_position
	data := &repository.SonexBadgeStakedPosition{
		PositionId:    event.TokenId.Int64(),
		PositionValue: totalBalanceUSD.String(),
		UserAddress:   event.User.String(),
	}
	retInserted, err := sonexBadgeStakedPosition.Insert(context.Background(), data)
	if err != nil {
		logx.Errorf("Insert staked tokenId :%v  Err:%v ", event.TokenId.Int64(), err)
		return
	}
	inserted, err := retInserted.LastInsertId()
	if err != nil {
		logx.Errorf("Insert staked tokenId :%v  LastInsertId Err:%v ", event.TokenId.Int64(), err)
		return
	}
	logx.Infof("Insert staked tokenId :%v success ret:%v  ", event.TokenId.Int64(), inserted)
}
