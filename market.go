package dydx

import (
	"context"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type Market struct {
	Market                           string           `json:"market,omitempty"`
	BaseAsset                        string           `json:"baseAsset,omitempty"`
	QuoteAsset                       string           `json:"quoteAsset,omitempty"`
	StepSize                         *decimal.Decimal `json:"stepSize,omitempty"`
	TickSize                         *decimal.Decimal `json:"tickSize,omitempty"`
	IndexPrice                       *decimal.Decimal `json:"indexPrice,omitempty"`
	OraclePrice                      *decimal.Decimal `json:"oraclePrice,omitempty"`
	PriceChange24H                   *decimal.Decimal `json:"priceChange24H,omitempty"`
	NextFundingRate                  *decimal.Decimal `json:"nextFundingRate,omitempty"`
	MinOrderSize                     *decimal.Decimal `json:"minOrderSize,omitempty"`
	Type                             string           `json:"type,omitempty"`
	InitialMarginFraction            *decimal.Decimal `json:"initialMarginFraction,omitempty"`
	MaintenanceMarginFraction        *decimal.Decimal `json:"maintenanceMarginFraction,omitempty"`
	BaselinePositionSize             *decimal.Decimal `json:"baselinePositionSize,omitempty"`
	IncrementalPositionSize          *decimal.Decimal `json:"incrementalPositionSize,omitempty"`
	IncrementalInitialMarginFraction *decimal.Decimal `json:"incrementalInitialMarginFraction,omitempty"`
	Volume24H                        *decimal.Decimal `json:"volume24H,omitempty"`
	Trades24H                        *decimal.Decimal `json:"trades24H,omitempty"`
	OpenInterest                     *decimal.Decimal `json:"openInterest,omitempty"`
	MaxPositionSize                  *decimal.Decimal `json:"maxPositionSize,omitempty"`
	AssetResolution                  string           `json:"assetResolution,omitempty"`
	SyntheticAssetID                 string           `json:"syntheticAssetId,omitempty"`
	Status                           string           `json:"status,omitempty"`
	NextFundingAt                    *time.Time       `json:"nextFundingAt,omitempty"`
}

// Market contains the meta data for a market.
// https://docs.dydx.exchange/#get-markets
type MarketsResponse struct {
	Markets map[string]Market `json:"markets"`
}

func (c *Client) GetMarkets(ctx context.Context) (*MarketsResponse, error) {
	return doRequest[MarketsResponse](ctx, c, http.MethodGet, "markets", "", nil, true)
}
