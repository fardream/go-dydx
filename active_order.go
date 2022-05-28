package dydx

import (
	"context"
	"net/http"
)

type ActiveOrder struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Market        string `json:"market"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	RemainingSize string `json:"remainingSize"`
}

type QueryActiveOrdersParam struct {
	Market string `url:"market,omitempty"`
	Side   string `url:"side,omitempty"`
	Id     string `url:"id,omitempty"`
}

type ActiveOrdersResponse struct {
	Orders []ActiveOrder `json:"orders"`
}

// GetActiveOrders implements https://docs.dydx.exchange/#cancel-active-orders
func (c *Client) GetActiveOrders(ctx context.Context, params *QueryActiveOrdersParam) (*ActiveOrdersResponse, error) {
	return doRequest[ActiveOrdersResponse](ctx, c, http.MethodGet, "active-orders", params, nil, false)
}

// CancelActiveOrders cancels mulitple orders at the same time, however, it utilizes the active order api.
// It implements https://docs.dydx.exchange/?json#cancel-active-orders
func (c *Client) CancelActiveOrders(ctx context.Context, params *CancelActiveOrdersParam) (*CancelActiveOrdersResponse, error) {
	return doRequest[CancelActiveOrdersResponse](ctx, c, http.MethodDelete, "active-orders", params, nil, false)
}

type CancelActiveOrdersResponse struct {
	CancelOrders []ActiveOrder `json:"cancelOrders"`
}

type CancelActiveOrdersParam struct {
	Market string `url:"market,omitempty"`
	Side   string `url:"side,omitempty"`
	Id     string `url:"id,omitempty"`
}
