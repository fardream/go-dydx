package dydx

import (
	"context"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type PositionResponse struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	Market        string           `json:"market,omitempty"`
	Status        string           `json:"status,omitempty"`
	Side          string           `json:"side,omitempty"`
	Size          decimal.Decimal  `json:"size,omitempty"`
	MaxSize       decimal.Decimal  `json:"maxSize,omitempty"`
	EntryPrice    decimal.Decimal  `json:"entryPrice,omitempty"`
	ExitPrice     *decimal.Decimal `json:"exitPrice,omitempty"`
	UnrealizedPnl decimal.Decimal  `json:"unrealizedPnl,omitempty"`
	RealizedPnl   decimal.Decimal  `json:"realizedPnl,omitempty"`
	CreatedAt     time.Time        `json:"createdAt,omitempty"`
	ClosedAt      *time.Time       `json:"closedAt,omitempty"`
	NetFunding    decimal.Decimal  `json:"netFunding,omitempty"`
	SumOpen       decimal.Decimal  `json:"sumOpen,omitempty"`
	SumClose      decimal.Decimal  `json:"sumClose,omitempty"`
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
