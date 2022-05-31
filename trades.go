package dydx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Trade struct {
	Side      string    `json:"side"`
	Size      Decimal   `json:"size"`
	Price     Decimal   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

// https://docs.dydx.exchange/?json#get-trades
type TradesResponse struct {
	Trades []Trade `json:"trades"`
}

type TradesParam struct {
	MarketID           string `url:"-"`
	Limit              int    `url:"limit,omitempty"`
	StartingBeforeOrAt string `url:"startingBeforeOrAt,omitempty"`
}

func (c *Client) GetTrades(ctx context.Context, params *TradesParam) (*TradesResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is nil (GetTrades requires params to be set with market id)")
	}
	if len(params.MarketID) == 0 {
		return nil, fmt.Errorf("market id in params is empty")
	}

	return doRequest[TradesResponse](ctx, c, http.MethodGet, urlJoin("trades", params.MarketID), params, nil, true)
}
