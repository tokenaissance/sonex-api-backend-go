package cache

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestTokenCacherGet(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	token := cacher.Get("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Printf("+v\n", token)
}

func TestTokenCacherSet(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	err := cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)
	if err != nil {
		t.Error(err)
	}
	token := cacher.Get("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Printf("+v\n", token)
}
