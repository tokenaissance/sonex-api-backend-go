package global

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/id_generater"
	"sonex_api_backend_go/internal/lib/quoter"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var cronExecuter *cron.Cron
var idGenerater *id_generater.IdGenerater
var tokenCachers map[int64]*cache.TokenCacher
var tokenCachersLock sync.RWMutex

func init() {
	cronExecuter = cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))
	tokenCachers = map[int64]*cache.TokenCacher{}

	cronExecuter.Start()
}

func AddCronJob(cron string, job func()) (cron.EntryID, error) {
	return cronExecuter.AddFunc(cron, job)
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

func GetToken(chainId int64, address string) (cache.Token, error) {
	tokenCachersLock.RLock()
	defer tokenCachersLock.RUnlock()
	cacher, ok := tokenCachers[chainId]
	if !ok {
		return cache.Token{}, fmt.Errorf("%s", "Miss token cacher")
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

func GetStablecoin(chainId int64) []cache.Token {
	tokenCachersLock.RLock()
	defer tokenCachersLock.RUnlock()
	cacher, ok := tokenCachers[chainId]
	if !ok {
		return nil
	}
	return cacher.GetStablecoins()
}

func GetPrice(chainId int64, token string) ([]quoter.Pool, *big.Int, error) {
	return nil, big.NewInt(0), nil
}
