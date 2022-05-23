package dydx

import "context"

type marketsChannelRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

func newMarketsChannelRequest() *marketsChannelRequest {
	return &marketsChannelRequest{
		Type:    subscribeChannelRequestType,
		Channel: MarketsChannel,
	}
}

type (
	MarketsChannelResponseContents = map[string]Market
	MarketsChannelResponse         = ChannelResponse[MarketsChannelResponseContents]
)

func (c *Client) SubscribeMarkets(ctx context.Context, outputChan chan<- *MarketsChannelResponse) error {
	return subscribeForType(ctx, c.wsUrl, newMarketsChannelRequest(), newUnsubscribeRequest(MarketsChannel, ""), outputChan)
}
