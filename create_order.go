package dydx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fardream/go-dydx/starkex"
)

// CreateOrderRequest is the post payload to create a new order
// https://docs.dydx.exchange/?json#create-a-new-order
type CreateOrderRequest struct {
	Signature       string `json:"signature"`
	Expiration      string `json:"expiration"`
	Market          string `json:"market"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	Size            string `json:"size"`
	Price           string `json:"price"`
	ClientId        string `json:"clientId"`
	TimeInForce     string `json:"timeInForce"`
	LimitFee        string `json:"limitFee"`
	CancelId        string `json:"cancelId,omitempty"`
	TriggerPrice    string `json:"triggerPrice,omitempty"`
	TrailingPercent string `json:"trailingPercent,omitempty"`
	PostOnly        bool   `json:"postOnly"`
}

type CreateOrderResponse struct {
	Order Order `json:"order"`
}

func NewCreateOrderRequest(market, side, order_type, size, price, clientid, tif, expiration, limitfee string, postonly bool) *CreateOrderRequest {
	return &CreateOrderRequest{
		Expiration:  expiration,
		Market:      market,
		Side:        side,
		Size:        size,
		Price:       price,
		ClientId:    clientid,
		TimeInForce: tif,
		PostOnly:    postonly,
		LimitFee:    limitfee,
		Type:        order_type,
	}
}

func (c *Client) NewOrder(ctx context.Context, order *CreateOrderRequest, positionId int64) (*CreateOrderResponse, error) {
	if order == nil {
		return nil, fmt.Errorf("order is null")
	}

	if len(order.Signature) == 0 {
		if c.starkKey == nil {
			return nil, fmt.Errorf("stark key is nil")
		}
		if len(c.starkKey.PrivateKey) == 0 {
			return nil, fmt.Errorf("start key is empty")
		}

		order_sign_params := starkex.OrderSignParam{
			NetworkId:  c.networkId,
			Market:     order.Market,
			Side:       order.Side,
			PositionId: positionId,
			HumanSize:  order.Size,
			HumanPrice: order.Price,
			LimitFee:   order.LimitFee,
			ClientId:   order.ClientId,
			Expiration: order.Expiration,
		}

		log.Debugf("sign order: %#v", order_sign_params)

		sign, err := starkex.OrderSign(c.starkKey.PrivateKey, order_sign_params)
		if err != nil {
			return nil, fmt.Errorf("failed to sign order: %w", err)
		}

		order.Signature = sign
	}

	payload, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order: %#v", err)
	}

	return doRequest[CreateOrderResponse](ctx, c, http.MethodPost, "orders", "", payload, false)
}
