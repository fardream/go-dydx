package dydx

import (
	"time"
)

type clientOption func(c *Client)

func SetApiKey(apiKey *ApiKey) clientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func SetStarkKey(starkKey *StarkKey) clientOption {
	return func(c *Client) {
		c.starkKey = starkKey
	}
}

func SetEndpoint(isMainnet bool) clientOption {
	rpc := ApiHostRopsten
	ws := WsHostRopsten
	networkId := NetworkIdRopsten
	if isMainnet {
		rpc = ApiHostMainnet
		ws = WsHostRopsten
		networkId = NetworkIdMainnet
	}

	return func(c *Client) {
		c.rpcUrl = rpc
		c.wsUrl = ws
		c.networkId = networkId
	}
}

func SetTimeout(timeout time.Duration) clientOption {
	return func(c *Client) {
		c.timeOut = timeout
	}
}

type Client struct {
	starkKey   *StarkKey
	apiKey     *ApiKey
	ethAddress string

	wsUrl     string
	rpcUrl    string
	networkId int

	timeOut time.Duration
}

func NewClient(starkKey *StarkKey, apiKey *ApiKey, ethAddress string, isMainnet bool, clientOptions ...clientOption) (*Client, error) {
	c := &Client{starkKey: starkKey, apiKey: apiKey, ethAddress: ethAddress, timeOut: time.Second * 15}

	SetEndpoint(isMainnet)(c)

	for _, option := range clientOptions {
		option(c)
	}

	return c, nil
}
