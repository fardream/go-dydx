package dydx

import "time"

type clientOption func(c *Client)

func SetClientEndpoint(isMainnet bool) clientOption {
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

func SetClientTimeout(timeout time.Duration) clientOption {
	return func(c *Client) {
		c.timeOut = timeout
	}
}

// Client is a struct holding the information necessary to connect to dydx.
type Client struct {
	starkKey   *StarkKey
	apiKey     *ApiKey
	ethAddress string

	wsUrl     string
	rpcUrl    string
	networkId int

	timeOut time.Duration
}

// NewClient creates a new Client, but doesn't connect to the dydx.exchange yet.
// If only public method is needed, keys and eth addersse can be empty/nil.
func NewClient(starkKey *StarkKey, apiKey *ApiKey, ethAddress string, isMainnet bool, clientOptions ...clientOption) (*Client, error) {
	c := &Client{starkKey: starkKey, apiKey: apiKey, ethAddress: ethAddress, timeOut: time.Second * 15}

	SetClientEndpoint(isMainnet)(c)

	for _, option := range clientOptions {
		option(c)
	}

	return c, nil
}
