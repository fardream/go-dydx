package dydx

import "context"

type orderbookChannelRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`

	ID             string `json:"id"`
	IncludeOffsets *bool  `json:"includeOffsets,omitempty"`
}

func newOrderbookChannelRequest(market string) *orderbookChannelRequest {
	r := &orderbookChannelRequest{Type: subscribeChannelRequestType, Channel: OrderbookChannel}
	r.ID = market
	b := true
	r.IncludeOffsets = &b
	return r
}

type (
	OrderbookChannelResponseContents = OrderbookResponse
	OrderbookChannelResponse         = ChannelResponse[OrderbookChannelResponseContents]
)

func (c *Client) SubscribeOrderbook(ctx context.Context, market string, outputChan chan<- *OrderbookChannelResponse) error {
	return subscribeForType(ctx, c.wsUrl, newOrderbookChannelRequest(market), newUnsubscribeRequest(OrderbookChannel, market), outputChan)
}
