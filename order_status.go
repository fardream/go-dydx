package dydx

import "encoding/json"

type OrderStatus string

const (
	OrderStatusPending     OrderStatus = "PENDING"
	OrderStatusOpen        OrderStatus = "OPEN"
	OrderStatusFilled      OrderStatus = "FILLED"
	OrderStatusCanceled    OrderStatus = "CANCELED"
	OrderStatusUntriggered OrderStatus = "UNTRIGGERED"
)

var orderstatuses []string = []string{string(OrderStatusPending), string(OrderStatusOpen), string(OrderStatusFilled), string(OrderStatusCanceled), string(OrderStatusUntriggered)}

func GetOrderStatus(input string) (OrderStatus, error) {
	return getProperStringEnum[OrderStatus](input, orderstatuses, "OrderStatus")
}

var (
	_ json.Marshaler   = (*OrderStatus)(nil)
	_ json.Unmarshaler = (*OrderStatus)(nil)
)

func (ot *OrderStatus) UnmarshalJSON(input []byte) error {
	return unmarshalJsonForStringEnum(input, "OrderStatus", ot, orderstatuses)
}

func (ot OrderStatus) MarshalJSON() ([]byte, error) {
	return marshalJsonForStringEnum(ot, "OrderStatus", orderstatuses)
}
