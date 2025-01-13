package pre_deposit

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestSign(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	pointsServiceConf := c.PointsService
	req := struct {
		User      string            `json:"user_address"`
		Points    int64             `json:"points"`
		Source    string            `json:"source"`
		SubTitle  string            `json:"sub_title"`
		ExtraId   string            `json:"extra_id"`
		Extra     map[string]string `json:"extra"`
		Signature string            `json:"signature"`
	}{
		User:     "event.User.Hex()",
		Points:   100,
		Source:   pointsServiceConf.Source,
		SubTitle: "preDepositConf.SubTitle",
		ExtraId:  fmt.Sprintf("%s-%d", "event.Raw.TxHash.Hex()", 1),
		Extra:    map[string]string{},
	}
	signature, err := Sign(&req, pointsServiceConf.PrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(signature)
}
