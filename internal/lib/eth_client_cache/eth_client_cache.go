package ethclientcache

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

type ethClientCacheT struct {
	lock    sync.Mutex
	clients map[string]*ethclient.Client
}

var ethClientCache = &ethClientCacheT{
	clients: make(map[string]*ethclient.Client),
}

// Get `*ethclient.Client` from cache if exists, otherwise initialize and cache a new one.
func GetEthClient(nodeUrl string) (*ethclient.Client, error) {
	ethClientCache.lock.Lock()
	defer ethClientCache.lock.Unlock()

	v, ok := ethClientCache.clients[nodeUrl]
	if ok {
		return v, nil
	}

	client, err := ethclient.DialContext(context.Background(), nodeUrl)
	if err != nil {
		return nil, err
	}

	ethClientCache.clients[nodeUrl] = client
	return client, nil
}
