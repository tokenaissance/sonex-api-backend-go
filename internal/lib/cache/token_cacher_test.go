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
	token, _ := cacher.Get("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Println(token)
}

func TestTokenCacherSet(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	err := cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)
	if err != nil {
		t.Error(err)
	}
	token, _ := cacher.Get("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Println(token)
}

func TestTokenCacherGetStablecoins(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	err := cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)
	if err != nil {
		t.Error(err)
	}
	tokens := cacher.GetStablecoins()
	fmt.Println(tokens)
}

func TestTokenCacherGetAll(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	err := cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)
	if err != nil {
		t.Error(err)
	}
	err = cacher.Set("0xA5D6513082EF1F157A33A066293309E74A8aF6Df", false)
	if err != nil {
		t.Error(err)
	}
	tokens := cacher.GetAll()
	fmt.Println(tokens)
}
