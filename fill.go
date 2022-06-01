package dydx

import (
	"context"
	"net/http"
	"time"
)

type FillsResponse struct {
	Fills []*Fill `json:"fills"`
}

type Fill struct {
	ID        string    `json:"id"`
	Side      OrderSide `json:"side"`
	Liquidity string    `json:"liquidity"`
	Type      OrderType `json:"type"`
	Market    string    `json:"market"`
	OrderID   string    `json:"orderId"`
	Price     Decimal   `json:"price"`
	Size      Decimal   `json:"size"`
	Fee       Decimal   `json:"fee"`
	CreatedAt time.Time `json:"createdAt"`
}

// FillsParam: https://docs.dydx.exchange/#get-fills
type FillsParam struct {
	Market            string `json:"market,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
	Limit             string `json:"limit,omitempty"`
	CreatedBeforeOrAt string `json:"createdBeforeOrAt,omitempty"`
}

// GetFills implements https://docs.dydx.exchange/#get-fills
func (c *Client) GetFills(ctx context.Context, params *FillsParam) (*FillsResponse, error) {
	return doRequest[FillsResponse](ctx, c, http.MethodGet, "fills", params, nil, false)
}
