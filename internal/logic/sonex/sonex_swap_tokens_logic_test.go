package sonex

import (
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/repository"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_getUserTokenDetailStructFromScout(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)

	logx.MustSetup(c.Log)
	chainConf := c.GetChainConfig(1946)
	conn := repository.NewConn(c.PostgreSQL.DataSource())
	// tokenList, err := getTokenList(conn, 1946, ORFERBY_ALL)
	tokenList, err := getTokenList(conn, 1946, ORFERBY_OFFICAL)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return
	}

	records, err := getUserTokenDetailStructFromScout(*chainConf, tokenList, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D", c.BlockScoutUrl[0], c.BlockScoutXRequsetId)

	// records, err := getUserTokenDetailStructFromScout(*chainConf, tokenList, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39733D", c.BlockScoutUrl[0], c.BlockScoutXRequsetId)
	// records, err := getUserTokenDetailStructFromScout(*chainConf, tokenList, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39", c.BlockScoutUrl[0], c.BlockScoutXRequsetId)

	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return
	}

	for i, v := range records {
		logx.Infof("records i:%#v v :%#v ", i, v)
	}
}
