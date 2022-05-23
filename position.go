package dydx

import (
	"context"
	"net/http"
	"time"
)

type PositionResponse struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	Market        string      `json:"market"`
	Status        string      `json:"status"`
	Side          string      `json:"side"`
	Size          string      `json:"size"`
	MaxSize       string      `json:"maxSize"`
	EntryPrice    string      `json:"entryPrice"`
	ExitPrice     interface{} `json:"exitPrice"`
	UnrealizedPnl string      `json:"unrealizedPnl"`
	RealizedPnl   string      `json:"realizedPnl"`
	CreatedAt     time.Time   `json:"createdAt"`
	ClosedAt      interface{} `json:"closedAt"`
	NetFunding    string      `json:"netFunding"`
	SumOpen       string      `json:"sumOpen"`
	SumClose      string      `json:"sumClose"`
}

const (
	PositionStatusOpen       = "OPEN"
	PositionStatusClosed     = "CLOSED"
	PositionStatusLiquidated = "LIQUIDATED"
)

type PositionParams struct {
	Market            string    `url:"market,omitempty"`
	Status            string    `url:"status,omitempty"`
	Limit             int       `url:"limit,omitempty"`
	CreatedBeforeOrAt time.Time `url:"createdBeforeOrAt,omitempty"`
}

func (c *Client) GetPositions(ctx context.Context, params *PositionParams) (*PositionResponse, error) {
	return doRequest[PositionResponse](ctx, c, http.MethodGet, "positions", params, nil, false)
}
