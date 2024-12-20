package global

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/id_generater"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var cronExecuter *cron.Cron
var idGenerater *id_generater.IdGenerater
var tokenCachers map[int64]*cache.TokenCacher
var tokenCachersLock sync.RWMutex

var prices map[int]map[string]*big.Int
var pricesLock sync.RWMutex

func init() {
	cronExecuter = cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))
	tokenCachers = map[int64]*cache.TokenCacher{}
	prices = map[int]map[string]*big.Int{}

	cronExecuter.Start()
}

func InitIdGenerater(workerId int64) error {
	generater, err := id_generater.NewIdGenerater(workerId)
	if err != nil {
		return err
	}
	idGenerater = generater
	return nil
}

func NewId() (int64, error) {
	return idGenerater.Generate()
}

func InitTokenCacher(chain config.ChainConfig) {
	tokenCachersLock.Lock()
	defer tokenCachersLock.Unlock()
	tokenCachers[chain.ChainId] = cache.NewTokenCacher(chain)
}

func GetToken(chainId int64, address string) *cache.Token {
	tokenCachersLock.RLock()
	defer tokenCachersLock.RUnlock()
	cacher, ok := tokenCachers[chainId]
	if !ok {
		return nil
	}
	return cacher.Get(address)
}

func SetToken(chainId int64, address string, isStablecoin bool) error {
	tokenCachersLock.RLock()
	defer tokenCachersLock.RUnlock()
	cacher, ok := tokenCachers[chainId]
	if !ok {
		return fmt.Errorf("%s", "Miss token cacher")
	}
	err := cacher.Set(address, isStablecoin)
	if err != nil {
		return err
	}
	return nil
}

func SetPrice(chainId int64, token string, price *big.Int) {
	pricesLock.Lock()
	defer pricesLock.Unlock()
	_, ok := prices[int(chainId)]
	if !ok {
		prices[int(chainId)] = map[string]*big.Int{}
	}
	prices[int(chainId)][token] = price
}

func GetPrice(chainId int64, token string) (*big.Int, error) {
	pricesLock.RLock()
	defer pricesLock.RUnlock()
	_, ok := prices[int(chainId)]
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token price")
	}
	price, ok := prices[int(chainId)][token]
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token price")
	}
	return price, nil
}
