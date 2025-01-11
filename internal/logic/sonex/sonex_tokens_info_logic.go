package sonex

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokensInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokensInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokensInfoLogic {
	return &SonexTokensInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokensInfoLogic) SonexTokensInfo(req *types.TokensInfoRequest) (resp *types.TokenInfoResponse, err error) {
	_tokenAddress := req.Address
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	ctx := context.Background()
	tokenModel := repository.NewTokenModel(conn)
	_tokenInfo, err := tokenModel.QueryByChainIdAndOneAddress(ctx, req.ChainId, _tokenAddress)
	if err != nil {
		logx.Infof("SonexTokensInfo ERR:%v ", err)
		return &types.TokenInfoResponse{StatusCode: 500}, errors.New("internal error ")
	}

	if !strings.EqualFold(_tokenInfo.Address, "") {

		var token0Additional types.PoolsResponse_TokenAdditional
		_ = json.Unmarshal([]byte(_tokenInfo.Additional.String), &token0Additional)

		return &types.TokenInfoResponse{
			StatusCode: 200,
			Token: types.PoolsResponse_Token{
				Address:    _tokenInfo.Address,
				Symbol:     _tokenInfo.Symbol,
				Name:       _tokenInfo.Name,
				Decimals:   _tokenInfo.Decimals,
				Logo:       _tokenInfo.Logo,
				Additional: token0Additional,
			},
		}, nil
	} else {
		return &types.TokenInfoResponse{StatusCode: 500}, errors.New("Internal  error")
	}

}
