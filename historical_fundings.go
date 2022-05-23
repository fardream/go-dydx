package dydx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HistoricalFundingsResponse struct {
	HistoricalFundings []HistoricalFunding `json:"historicalFunding"`
}

type HistoricalFunding struct {
	Market      string    `json:"-"`
	Rate        string    `json:"rate"`
	Price       string    `json:"price"`
	EffectiveAt time.Time `json:"effectiveAt"`
}

type HistoricalFundingsParam struct {
	Market              string `url:"-"`
	EffectiveBeforeOrAt string `url:"effectiveBeforeOrAt,omitempty"`
}

func (c *Client) GetHistoricalFunding(ctx context.Context, params *HistoricalFundingsParam) (*HistoricalFundingsResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil for historical funding. market must be provided")
	}
	if len(params.Market) == 0 {
		return nil, fmt.Errorf("market cannot be empty historical funding")
	}

	return doRequest[HistoricalFundingsResponse](ctx, c, http.MethodGet, urlJoin("historical-funding", params.Market), params, nil, true)
}
