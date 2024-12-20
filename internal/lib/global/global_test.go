package global

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestGlobalGetToken(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	for _, conf := range c.Chains {
		InitTokenCacher(conf)
	}
	token := GetToken(1946, "0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Printf("%+v\n", token)
}
