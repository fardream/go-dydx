package dydx

import "encoding/json"

type TimeInForce string

const (
	TimeInForceGtt TimeInForce = "GTT" // Good til time
	TimeInForceFok TimeInForce = "FOK" // Fill or Kill
	TimeInForceIoc TimeInForce = "IOC" // Immediate or Cancel
)

var tifs []string = []string{string(TimeInForceFok), string(TimeInForceGtt), string(TimeInForceIoc)}

func GetTimeInForce(input string) (TimeInForce, error) {
	return getProperStringEnum[TimeInForce](input, tifs, "TimeInForce")
}

var (
	_ json.Marshaler   = (*TimeInForce)(nil)
	_ json.Unmarshaler = (*TimeInForce)(nil)
)

func (ot *TimeInForce) UnmarshalJSON(input []byte) error {
	return unmarshalJsonForStringEnum(input, "TimeInForce", ot, tifs)
}

func (ot TimeInForce) MarshalJSON() ([]byte, error) {
	return marshalJsonForStringEnum(ot, "TimeInForce", tifs)
}
