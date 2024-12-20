package cache

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenCacher struct {
	chain      config.ChainConfig
	tokens     map[string]Token
	tokensLock sync.RWMutex
}

func NewTokenCacher(chain config.ChainConfig) *TokenCacher {
	return &TokenCacher{
		chain:  chain,
		tokens: map[string]Token{},
	}
}

func (c *TokenCacher) Get(address string) *Token {
	var result Token
	if token, ok := c.tokens[address]; ok {
		result = token
	} else {
		client, err := ethclient.Dial(c.chain.GetNode())
		if err != nil {
			return nil
		}
		erc20, err := contract_common.NewErc20(common.HexToAddress(address), client)
		if err != nil {
			return nil
		}
		name, err := erc20.Name(nil)
		if err != nil {
			panic(err)
		}
		symbol, err := erc20.Symbol(nil)
		if err != nil {
			panic(err)
		}
		decimals, err := erc20.Decimals(nil)
		if err != nil {
			panic(err)
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
	return &result
}

func (c *TokenCacher) Set(address string, isStablecoin bool) error {
	token := c.Get(address)
	if token == nil {
		return fmt.Errorf("%s", "Get token infomaion failed")
	}
	token.IsStablecoin = isStablecoin
	c.tokensLock.Lock()
	defer c.tokensLock.Unlock()
	c.tokens[address] = *token
	return nil
}
