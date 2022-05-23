package dydx

import (
	"context"
	"net/http"
	"time"
)

type HistoricalPnLResponse struct {
	HistoricalPnLs []HistoricalPnL `json:"historicalPnl"`
}

type HistoricalPnL struct {
	AccountID    string    `json:"accountId"`
	Equity       string    `json:"equity"`
	TotalPnl     string    `json:"totalPnl"`
	NetTransfers string    `json:"netTransfers"`
	CreatedAt    time.Time `json:"createdAt"`
}

type HistoricalPnLParam struct {
	EffectiveBeforeOrAt string `url:"effectiveBeforeOrAt,omitempty"`
	EffectiveAtOrAfter  string `url:"effectiveAtOrAfter,omitempty"`
}

func (c *Client) GetHistoricalPnL(ctx context.Context, params *HistoricalPnLParam) (*HistoricalPnLResponse, error) {
	return doRequest[HistoricalPnLResponse](ctx, c, http.MethodGet, "historical-pnl", params, nil, false)
}
