package dydx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// https://docs.dydx.exchange/?json#get-candles-for-market
type CandlesResponse struct {
	Candles []Candle `json:"candles"`
}

type Candle struct {
	Market               string    `json:"market"`
	Resolution           string    `json:"resolution"`
	Low                  string    `json:"low"`
	High                 string    `json:"high"`
	Open                 string    `json:"open"`
	Close                string    `json:"close"`
	BaseTokenVolume      string    `json:"baseTokenVolume"`
	Trades               string    `json:"trades"`
	UsdVolume            string    `json:"usdVolume"`
	StartingOpenInterest string    `json:"startingOpenInterest"`
	StartedAt            time.Time `json:"startedAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type CandlesParam struct {
	Market     string `url:"-"`
	Resolution string `url:"resolution,omitempty"`
	FromISO    string `url:"fromISO,omitempty"`
	ToISO      string `url:"toISO,omitempty"`
	// Max:100
	Limit int `url:"limit,omitempty"`
}

func (c *Client) GetCandles(ctx context.Context, params *CandlesParam) (*CandlesResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil for candles request, market must be provided")
	}
	if params.Market == "" {
		return nil, fmt.Errorf("market cannot be empty for candles request")
	}
	return doRequest[CandlesResponse](ctx, c, http.MethodGet, urlJoin("candles", params.Market), "", nil, true)
}
