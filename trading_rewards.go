package dydx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TradingRewardsResponse TradingReward

type TradingReward struct {
	Epoch      int       `json:"epoch"`
	EpochStart time.Time `json:"epochStart"`
	EpochEnd   time.Time `json:"epochEnd"`
	Fees       struct {
		FeesPaid      string `json:"feesPaid"`
		TotalFeesPaid string `json:"totalFeesPaid"`
	} `json:"fees"`
	OpenInterest struct {
		AverageOpenInterest      string `json:"averageOpenInterest"`
		TotalAverageOpenInterest string `json:"totalAverageOpenInterest"`
	} `json:"openInterest"`
	StakedDYDX struct {
		AverageStakedDYDX          string `json:"averageStakedDYDX"`
		AverageStakedDYDXWithFloor string `json:"averageStakedDYDXWithFloor"`
		TotalAverageStakedDYDX     string `json:"totalAverageStakedDYDX"`
	} `json:"stakedDYDX"`
	Weight struct {
		Weight      string `json:"weight"`
		TotalWeight string `json:"totalWeight"`
	} `json:"weight"`
	TotalRewards     string `json:"totalRewards"`
	EstimatedRewards string `json:"estimatedRewards"`
}

func (c *Client) GetTradingRewards(ctx context.Context, epoch int64) (*TradingRewardsResponse, error) {
	params := new(url.Values)
	if epoch > 0 {
		params.Add("epoch", strconv.FormatInt(epoch, 10))
	}
	return doRequest[TradingRewardsResponse](ctx, c, http.MethodGet, "rewards/weight", params, nil, false)
}
