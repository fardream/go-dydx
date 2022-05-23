package dydx

import (
	"context"
	"net/http"
	"time"
)

type FundingPaymentsResponse struct {
	FundingPayments []FundingPayment `json:"fundingPayments"`
}

type FundingPayment struct {
	Market       string    `json:"market"`
	Payment      string    `json:"payment"`
	Rate         string    `json:"rate"`
	PositionSize string    `json:"positionSize"`
	Price        string    `json:"price"`
	EffectiveAt  time.Time `json:"effectiveAt"`
}

type FundingPaymentsParam struct {
	Market              string `url:"market,omitempty"`
	Limit               string `url:"limit,omitempty"`
	EffectiveBeforeOrAt string `url:"effectiveBeforeOrAt,omitempty"`
}

func (c *Client) GetFundingPayments(ctx context.Context, params *FundingPaymentsParam) (*FundingPaymentsResponse, error) {
	return doRequest[FundingPaymentsResponse](ctx, c, http.MethodGet, "funding", params, nil, false)
}
