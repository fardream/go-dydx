package dydx

import (
	"context"
	"net/http"
	"time"
)

type FillsResponse struct {
	Fills []Fill `json:"fills"`
}

type Fill struct {
	ID        string    `json:"id"`
	Side      string    `json:"side"`
	Liquidity string    `json:"liquidity"`
	Type      string    `json:"type"`
	Market    string    `json:"market"`
	OrderID   string    `json:"orderId"`
	Price     string    `json:"price"`
	Size      string    `json:"size"`
	Fee       string    `json:"fee"`
	CreatedAt time.Time `json:"createdAt"`
}

type FillsParam struct {
	Market            string `json:"market,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
	Limit             string `json:"limit,omitempty"`
	CreatedBeforeOrAt string `json:"createdBeforeOrAt,omitempty"`
}

func (c *Client) GetFills(ctx context.Context, params *FillsParam) (*FillsResponse, error) {
	return doRequest[FillsResponse](ctx, c, http.MethodGet, "fills", params, nil, false)
}
