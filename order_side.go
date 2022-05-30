package dydx

import "encoding/json"

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

var orderSides []string = []string{string(OrderSideBuy), string(OrderSideSell)}

func GetOrderSide(input string) (OrderSide, error) {
	return getProperStringEnum[OrderSide](input, orderSides, "OrderSide")
}

var (
	_ json.Marshaler   = (*OrderSide)(nil)
	_ json.Unmarshaler = (*OrderSide)(nil)
)

func (os *OrderSide) UnmarshalJSON(input []byte) error {
	return unmarshalJsonForStringEnum(input, "OrderSide", os, orderSides)
}

func (os OrderSide) MarshalJSON() ([]byte, error) {
	return marshalJsonForStringEnum(os, "OrderSide", orderSides)
}
