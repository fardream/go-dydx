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
	Market        string     `json:"market,omitempty"`
	Status        string     `json:"status,omitempty"`
	Side          string     `json:"side,omitempty"`
	Size          Decimal    `json:"size,omitempty"`
	MaxSize       Decimal    `json:"maxSize,omitempty"`
	EntryPrice    Decimal    `json:"entryPrice,omitempty"`
	ExitPrice     *Decimal   `json:"exitPrice,omitempty"`
	UnrealizedPnl Decimal    `json:"unrealizedPnl,omitempty"`
	RealizedPnl   Decimal    `json:"realizedPnl,omitempty"`
	CreatedAt     time.Time  `json:"createdAt,omitempty"`
	ClosedAt      *time.Time `json:"closedAt,omitempty"`
	NetFunding    Decimal    `json:"netFunding,omitempty"`
	SumOpen       Decimal    `json:"sumOpen,omitempty"`
	SumClose      Decimal    `json:"sumClose,omitempty"`
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
