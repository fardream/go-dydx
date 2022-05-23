package dydx

import (
	"context"
	"net/http"
	"time"
)

type Market struct {
	Market                           string    `json:"market"`
	BaseAsset                        string    `json:"baseAsset"`
	QuoteAsset                       string    `json:"quoteAsset"`
	StepSize                         string    `json:"stepSize"`
	TickSize                         string    `json:"tickSize"`
	IndexPrice                       string    `json:"indexPrice"`
	OraclePrice                      string    `json:"oraclePrice"`
	PriceChange24H                   string    `json:"priceChange24H"`
	NextFundingRate                  string    `json:"nextFundingRate"`
	MinOrderSize                     string    `json:"minOrderSize"`
	Type                             string    `json:"type"`
	InitialMarginFraction            string    `json:"initialMarginFraction"`
	MaintenanceMarginFraction        string    `json:"maintenanceMarginFraction"`
	BaselinePositionSize             string    `json:"baselinePositionSize"`
	IncrementalPositionSize          string    `json:"incrementalPositionSize"`
	IncrementalInitialMarginFraction string    `json:"incrementalInitialMarginFraction"`
	Volume24H                        string    `json:"volume24H"`
	Trades24H                        string    `json:"trades24H"`
	OpenInterest                     string    `json:"openInterest"`
	MaxPositionSize                  string    `json:"maxPositionSize"`
	AssetResolution                  string    `json:"assetResolution"`
	SyntheticAssetID                 string    `json:"syntheticAssetId"`
	Status                           string    `json:"status"`
	NextFundingAt                    time.Time `json:"nextFundingAt"`
}

// Market contains the meta data for a market.
// https://docs.dydx.exchange/#get-markets
type MarketsResponse struct {
	Markets map[string]Market `json:"markets"`
}

func (c *Client) GetMarkets(ctx context.Context) (*MarketsResponse, error) {
	return doRequest[MarketsResponse](ctx, c, http.MethodGet, "markets", "", nil, true)
}
