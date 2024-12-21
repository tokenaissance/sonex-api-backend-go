package cache

import (
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/quoter"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenCacher struct {
	chain      config.ChainConfig
	tokens     map[string]Token
	tokensLock sync.RWMutex
	prices     map[string]*big.Int
	pricesLock sync.RWMutex
}

type Token struct {
	Address      string
	Name         string
	Symbol       string
	Decimals     int64
	IsStablecoin bool
}

type Price struct {
	Token     string
	Pools     []quoter.Pool
	Path      []byte
	Value     *big.Int
	UpdatedAt time.Time
}

func NewTokenCacher(chain config.ChainConfig) *TokenCacher {
	return &TokenCacher{
		chain:  chain,
		tokens: map[string]Token{},
		prices: map[string]*big.Int{},
	}
}

func (cacher *TokenCacher) Get(address string) (Token, error) {
	cacher.tokensLock.RLock()
	token, ok := cacher.tokens[address]
	cacher.tokensLock.RUnlock()
	if !ok {
		client, err := ethclient.Dial(cacher.chain.GetNode())
		if err != nil {
			return Token{}, err
		}
		erc20, err := contract_common.NewErc20(common.HexToAddress(address), client)
		if err != nil {
			return Token{}, err
		}
		name, err := erc20.Name(nil)
		if err != nil {
			return Token{}, err
		}
		symbol, err := erc20.Symbol(nil)
		if err != nil {
			return Token{}, err
		}
		decimals, err := erc20.Decimals(nil)
		if err != nil {
			return Token{}, err
		}
		token = Token{
			Address:      address,
			Name:         name,
			Symbol:       symbol,
			Decimals:     int64(decimals),
			IsStablecoin: false,
		}
		cacher.tokensLock.Lock()
		cacher.tokens[address] = token
		cacher.tokensLock.Unlock()
	}
	return token, nil
}

func (c *TokenCacher) Set(address string, isStablecoin bool) error {
	token, err := c.Get(address)
	if err != nil {
		return err
	}
	token.IsStablecoin = isStablecoin
	c.tokensLock.Lock()
	c.tokens[address] = token
	c.tokensLock.Unlock()
	return nil
}

func (cacher *TokenCacher) GetStablecoins() []Token {
	cacher.tokensLock.RLock()
	tokens := cacher.tokens
	cacher.tokensLock.RUnlock()
	var result []Token
	linq.From(tokens).
		WhereT(func(item linq.KeyValue) bool {
			token, _ := item.Value.(Token)
			return token.IsStablecoin
		}).
		SelectT(func(item linq.KeyValue) Token { return item.Value.(Token) }).
		ToSlice(&result)
	return result
}

func (cacher *TokenCacher) GetAll() []Token {
	cacher.tokensLock.RLock()
	tokens := cacher.tokens
	cacher.tokensLock.RUnlock()
	var result []Token
	linq.From(tokens).
		SelectT(func(item linq.KeyValue) Token { return item.Value.(Token) }).
		ToSlice(&result)
	return result
}

func (cacher *TokenCacher) GetPrice(address string, isCache bool) ([]quoter.Pool, []byte, *big.Int, error) {
	cacher.pricesLock.RLock()
	token, err := cacher.Get(address)
	if err != nil {
		return nil, nil, nil, err
	}
	cacher.pricesLock.RUnlock()
	_ = token
	return nil, nil, nil, nil
}
