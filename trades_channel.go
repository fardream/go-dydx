package dydx

import "context"

type tradesChannelRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`

	ID             string `json:"id"`
	IncludeOffsets *bool  `json:"includeOffsets,omitempty"`
}

type tradesChannelResponseContents = TradesResponse

func newTradesChannelRequest(market string) *tradesChannelRequest {
	return &tradesChannelRequest{
		Type:    subscribeChannelRequestType,
		Channel: TradesChannel,
		ID:      market,
	}
}

type TradesChannelResponse = ChannelResponse[tradesChannelResponseContents]

func (c *Client) SubscribeTrades(ctx context.Context, market string, outputChan chan<- *TradesChannelResponse) error {
	return subscribeForType(ctx, c.wsUrl, newTradesChannelRequest(market), newUnsubscribeRequest(TradesChannel, market), outputChan)
}
