package pre_deposit

import (
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"
)

func StakedCallback(conf config.Config, chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked) {
	conn := repository.NewConn(conf.PostgreSQL.DateSource())
	_ = conn

	// basePoints, _ := new(big.Int).SetString(pointsServiceConf.Points, 10)
	// points := new(big.Int).Div(new(big.Int).Mul(basePoints, big.NewInt(int64(pointsServiceConf.PointsRate*100))), big.NewInt(100))
}
