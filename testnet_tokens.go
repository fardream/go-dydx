package dydx

import (
	"context"
	"net/http"
)

type RequestTestnetTokensResponse struct {
	*Transfer `json:"transfer,omitempty"`
}

// RequestTestnetTokens implements https://docs.dydx.exchange/#request-testnet-tokens
// Obtain testnet tokens (USDC)
func (c *Client) RequestTestnetTokens(ctx context.Context) (*RequestTestnetTokensResponse, error) {
	return doRequest[RequestTestnetTokensResponse](ctx, c, http.MethodPost, "testnet/tokens", "", nil, false)
}
