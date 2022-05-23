package dydx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	OrderStatusPending     = "PENDING"
	OrderStatusOpen        = "OPEN"
	OrderStatusFilled      = "FILLED"
	OrderStatusCanceled    = "CANCELED"
	OrderStatusUntriggered = "UNTRIGGERED"
)

// Order is the information returned from dydx
type Order struct {
	ID              string    `json:"id"`
	ClientID        string    `json:"clientId"`
	AccountID       string    `json:"accountId"`
	Market          string    `json:"market"`
	Side            string    `json:"side"`
	Price           string    `json:"price"`
	TriggerPrice    string    `json:"triggerPrice"`
	TrailingPercent string    `json:"trailingPercent"`
	Size            string    `json:"size"`
	RemainingSize   string    `json:"remainingSize"`
	Type            string    `json:"type"`
	UnfillableAt    string    `json:"unfillableAt"`
	Status          string    `json:"status"`
	TimeInForce     string    `json:"timeInForce"`
	CancelReason    string    `json:"cancelReason"`
	PostOnly        bool      `json:"postOnly"`
	CreatedAt       time.Time `json:"createdAt"`
	ExpiresAt       time.Time `json:"expiresAt"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

type OrderResponse struct {
	Order Order `json:"order"`
}

type CancelOrderResponse struct {
	CancelOrder Order `json:"cancelOrder"`
}

type CancelOrdersResponse struct {
	CancelOrders []Order `json:"cancelOrders"`
}

type OrderQueryParam struct {
	Limit              int    `url:"limit,omitempty"`
	Market             string `url:"market,omitempty"`
	Status             string `url:"status,omitempty"`
	Type               string `url:"type,omitempty"`
	Side               string `url:"side,omitempty"`
	CreatedBeforeOrAt  string `url:"createdAt,omitempty"`
	ReturnLatestOrders string `url:"returnLatestOrders,omitempty"`
}

func (c *Client) GetOrders(ctx context.Context, params *OrderQueryParam) (*OrdersResponse, error) {
	return doRequest[OrdersResponse](ctx, c, http.MethodGet, "orders", params, nil, false)
}

func (c *Client) GetOrderById(ctx context.Context, id string) (*OrderResponse, error) {
	return doRequest[OrderResponse](ctx, c, http.MethodGet, fmt.Sprintf("orders/%s", id), "", nil, false)
}

func (c *Client) GetOrderByClientId(ctx context.Context, clientId string) (*OrderResponse, error) {
	return doRequest[OrderResponse](ctx, c, http.MethodGet, fmt.Sprintf("orders/client/%s", clientId), "", nil, false)
}
