package global_storage

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/lib/id_generater"
	"sync"
)

var idGenerater *id_generater.IdGenerater

type Token struct {
	Address      string
	Name         string
	Symbol       string
	Decimals     int64
	IsStablecoin bool
}

var tokens map[int]map[string]Token
var tokensLock sync.RWMutex

var prices map[int]map[string]*big.Int
var pricesLock sync.RWMutex

func init() {
	tokens = map[int]map[string]Token{}
	prices = map[int]map[string]*big.Int{}
}

func NewIdGenerater(workerId int64) error {
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

func SetToken(chainId int64, address string, token Token) {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	_, ok := tokens[int(chainId)]
	if !ok {
		tokens[int(chainId)] = map[string]Token{}
	}
	tokens[int(chainId)][address] = token
}

func GetToken(chainId int64, address string) (*Token, error) {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	_, ok := tokens[int(chainId)]
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token infomation")
	}
	token, ok := tokens[int(chainId)][address]
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token infomation")
	}
	return &token, nil
}

func AllTokens(chainId int64) []Token {
	tokensLock.RLock()
	defer tokensLock.RUnlock()
	_, ok := tokens[int(chainId)]
	if !ok {
		tokens[int(chainId)] = map[string]Token{}
		return []Token{}
	}
	result := []Token{}
	for _, value := range tokens[int(chainId)] {
		result = append(result, value)
	}
	return result
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
