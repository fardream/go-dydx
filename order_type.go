package dydx

import (
	"encoding/json"
)

// OrderType indicates if the order is market, limit, stop, trailing stop or taking profit.
// See https://docs.dydx.exchange/#order-types
type OrderType string

const (
	OrderTypeMarket       OrderType = "MARKET"
	OrderTypeLimit        OrderType = "LIMIT"
	OrderTypeStop         OrderType = "STOP"
	OrderTypeTrailingStop OrderType = "TRAILING_STOP"
	OrderTypeTakingProfit OrderType = "TAKE_PROFIT"
)

var orderTypes []string = []string{string(OrderTypeMarket), string(OrderTypeLimit), string(OrderTypeStop), string(OrderTypeTrailingStop), string(OrderTypeTakingProfit)}

func GetOrderType(input string) (OrderType, error) {
	return getProperStringEnum[OrderType](input, orderTypes, "OrderType")
}

var (
	_ json.Marshaler   = (*OrderType)(nil)
	_ json.Unmarshaler = (*OrderType)(nil)
)

func (ot *OrderType) UnmarshalJSON(input []byte) error {
	return unmarshalJsonForStringEnum(input, "OrderType", ot, orderTypes)
}

func (ot OrderType) MarshalJSON() ([]byte, error) {
	return marshalJsonForStringEnum(ot, "OrderType", orderTypes)
}
