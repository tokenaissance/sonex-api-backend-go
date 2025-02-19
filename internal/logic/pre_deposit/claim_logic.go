package pre_deposit

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get claim signature
func NewClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimLogic {
	return &ClaimLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimLogic) Claim(req *types.ClaimRequest) (resp *types.ClaimResponse, err error) {
	if req.UserAddress == "" || !common.IsHexAddress(req.UserAddress) {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	userAddress := common.HexToAddress(req.UserAddress)

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	preDepositConfig := l.svcCtx.Config.GetPreDepositConfig(req.ChainId)
	if preDepositConfig == nil {
		return nil, fmt.Errorf("%s", "Miss pre deposit configure")
	}
	rewardTokenAddress := common.HexToAddress(preDepositConfig.RewardToken)
	rewardClaimerAddress := common.HexToAddress(preDepositConfig.RewardClaimer)

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	preDepositRewardClaimRecordModel := repository.NewPreDepositRewardClaimRecordModel(conn)
	preDepositRewardSnapshotModel := repository.NewPreDepositRewardSnapshotModel(conn)
	eventPreDepositClaimerRewardClaimedModel := repository.NewEventPreDepositClaimerRewardClaimedModel(conn)

	claimRecord, err := preDepositRewardClaimRecordModel.FindLatestByChainIdUser(l.ctx, req.ChainId, userAddress.Hex())
	var nonce int64
	switch err {
	case repository.ErrNotFound:
		nonce = 1
	case nil:
		// There is a unclaimed record, return latest signature
		if !claimRecord.TxHash.Valid {
			return &types.ClaimResponse{
				Token:     rewardTokenAddress.Hex(),
				Amount:    claimRecord.Amount,
				Nonce:     claimRecord.Nonce,
				Signature: claimRecord.Signature,
			}, nil
		} else {
			nonce = claimRecord.Nonce + 1
		}
	default:
		return nil, err
	}

	// Total reward
	totalReward, err := preDepositRewardSnapshotModel.SummarizeTotalRewardByChainIdUserToken(l.ctx, req.ChainId, userAddress.String(), rewardTokenAddress.String())
	if err != nil {
		return nil, err
	}
	totalClaimed, err := eventPreDepositClaimerRewardClaimedModel.SummarizeTotalClaimedByChainIdUserToken(l.ctx, req.ChainId, userAddress.String(), rewardTokenAddress.String())
	if err != nil {
		return nil, err
	}
	reward := new(big.Int).Sub(totalReward, totalClaimed)
	if reward.Cmp(common.Big0) <= 0 {
		return &types.ClaimResponse{
			Token:     rewardTokenAddress.Hex(),
			Amount:    common.Big0.String(),
			Nonce:     0,
			Signature: "",
		}, nil
	}

	domainSeparator := crypto.Keccak256([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	domainSeparator = append(domainSeparator, crypto.Keccak256([]byte("RewardClaimer"))...)
	domainSeparator = append(domainSeparator, crypto.Keccak256([]byte("1.0"))...)
	domainSeparator = append(domainSeparator, common.LeftPadBytes(big.NewInt(req.ChainId).Bytes(), 32)...)
	domainSeparator = append(domainSeparator, common.LeftPadBytes(rewardClaimerAddress.Bytes(), 32)...)
	domainSeparator = crypto.Keccak256(domainSeparator)

	structHash := crypto.Keccak256([]byte("Claim(address user,address token,uint256 amount,uint256 nonce)"))
	structHash = append(structHash, common.LeftPadBytes(userAddress.Bytes(), 32)...)
	structHash = append(structHash, common.LeftPadBytes(rewardTokenAddress.Bytes(), 32)...)
	structHash = append(structHash, common.LeftPadBytes(reward.Bytes(), 32)...)
	structHash = append(structHash, common.LeftPadBytes(big.NewInt(nonce).Bytes(), 32)...)
	structHash = crypto.Keccak256(structHash)

	digest := crypto.Keccak256(append(append([]byte("\x19\x01"), domainSeparator...), structHash...))

	privateKey, err := crypto.HexToECDSA(preDepositConfig.SignerPrivateKey)
	if err != nil {
		return nil, err
	}
	signatureBytes, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return nil, err
	}
	if signatureBytes[64] < 27 {
		signatureBytes[64] += 27
	}
	signature := fmt.Sprintf("0x%s", common.Bytes2Hex(signatureBytes))

	id, err := global.NewId()
	if err != nil {
		return nil, err
	}
	err = preDepositRewardClaimRecordModel.New(l.ctx, id, req.ChainId, userAddress.Hex(), rewardTokenAddress.Hex(), reward.String(), nonce, signature)
	if err != nil {
		return nil, err
	}

	return &types.ClaimResponse{
		Token:     rewardTokenAddress.Hex(),
		Amount:    reward.String(),
		Nonce:     nonce,
		Signature: signature,
	}, nil
}
