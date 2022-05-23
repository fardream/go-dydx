package dydx

import (
	"context"
	"fmt"
	"net/http"
)

// CancelOrder
func (c *Client) CancelOrder(ctx context.Context, id string) (*CancelOrderResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("order id is empty")
	}

	return doRequest[CancelOrderResponse](ctx, c, http.MethodDelete, urlJoin("orders", id), "", nil, false)
}

// CancelOrders cancels multiple orders at the same time.
// It implements https://docs.dydx.exchange/?json#cancel-orders
func (c *Client) CancelOrders(ctx context.Context, params *CancelOrdersParam) (*CancelOrdersResponse, error) {
	return doRequest[CancelOrdersResponse](ctx, c, http.MethodDelete, "orders", params, nil, false)
}

type CancelOrdersParam struct {
	Market string `url:"market,omitempty"`
}
