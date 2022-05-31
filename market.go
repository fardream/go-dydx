package dydx

import (
	"context"
	"net/http"
	"time"
)

type Market struct {
	Market                           string     `json:"market,omitempty"`
	BaseAsset                        string     `json:"baseAsset,omitempty"`
	QuoteAsset                       string     `json:"quoteAsset,omitempty"`
	StepSize                         *Decimal   `json:"stepSize,omitempty"`
	TickSize                         *Decimal   `json:"tickSize,omitempty"`
	IndexPrice                       *Decimal   `json:"indexPrice,omitempty"`
	OraclePrice                      *Decimal   `json:"oraclePrice,omitempty"`
	PriceChange24H                   *Decimal   `json:"priceChange24H,omitempty"`
	NextFundingRate                  *Decimal   `json:"nextFundingRate,omitempty"`
	MinOrderSize                     *Decimal   `json:"minOrderSize,omitempty"`
	Type                             string     `json:"type,omitempty"`
	InitialMarginFraction            *Decimal   `json:"initialMarginFraction,omitempty"`
	MaintenanceMarginFraction        *Decimal   `json:"maintenanceMarginFraction,omitempty"`
	BaselinePositionSize             *Decimal   `json:"baselinePositionSize,omitempty"`
	IncrementalPositionSize          *Decimal   `json:"incrementalPositionSize,omitempty"`
	IncrementalInitialMarginFraction *Decimal   `json:"incrementalInitialMarginFraction,omitempty"`
	Volume24H                        *Decimal   `json:"volume24H,omitempty"`
	Trades24H                        *Decimal   `json:"trades24H,omitempty"`
	OpenInterest                     *Decimal   `json:"openInterest,omitempty"`
	MaxPositionSize                  *Decimal   `json:"maxPositionSize,omitempty"`
	AssetResolution                  string     `json:"assetResolution,omitempty"`
	SyntheticAssetID                 string     `json:"syntheticAssetId,omitempty"`
	Status                           string     `json:"status,omitempty"`
	NextFundingAt                    *time.Time `json:"nextFundingAt,omitempty"`
}

// Market contains the meta data for a market.
// https://docs.dydx.exchange/#get-markets
type MarketsResponse struct {
	Markets map[string]Market `json:"markets"`
}

func (c *Client) GetMarkets(ctx context.Context) (*MarketsResponse, error) {
	return doRequest[MarketsResponse](ctx, c, http.MethodGet, "markets", "", nil, true)
}
