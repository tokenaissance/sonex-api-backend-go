package cache

import (
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
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

func (c *TokenCacher) Get(address string) (Token, error) {
	var result Token
	if token, ok := c.tokens[address]; ok {
		result = token
	} else {
		client, err := ethclient.Dial(c.chain.GetNode())
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
		c.tokensLock.Lock()
		defer c.tokensLock.Unlock()
		c.tokens[address] = Token{
			Address:      address,
			Name:         name,
			Symbol:       symbol,
			Decimals:     int64(decimals),
			IsStablecoin: false,
		}
		result = c.tokens[address]
	}
	return result, nil
}

func (c *TokenCacher) Set(address string, isStablecoin bool) error {
	token, err := c.Get(address)
	if err != nil {
		return err
	}
	token.IsStablecoin = isStablecoin
	c.tokensLock.Lock()
	defer c.tokensLock.Unlock()
	c.tokens[address] = token
	return nil
}

func (c *TokenCacher) GetStablecoins() []Token {
	c.tokensLock.RLock()
	defer c.tokensLock.RUnlock()
	var result []Token
	linq.From(c.tokens).
		WhereT(func(item linq.KeyValue) bool {
			token, _ := item.Value.(Token)
			return token.IsStablecoin
		}).
		SelectT(func(item linq.KeyValue) Token { return item.Value.(Token) }).
		ToSlice(&result)
	return result
}

func (c *TokenCacher) GetAll() []Token {
	c.tokensLock.RLock()
	defer c.tokensLock.RUnlock()
	var result []Token
	linq.From(c.tokens).
		SelectT(func(item linq.KeyValue) Token { return item.Value.(Token) }).
		ToSlice(&result)
	return result
}

func (c *TokenCacher) GetPrice(address string) (*big.Int, error) {
	c.pricesLock.RLock()
	defer c.pricesLock.RUnlock()
	return nil, nil
}
