package dydx

import (
	"encoding/json"
	"strconv"
)

// JsonInt is an int, but can be json-unmarshaled from either int or string, marshalled into string.
type JsonInt int

var (
	_ json.Marshaler   = (*JsonInt)(nil)
	_ json.Unmarshaler = (*JsonInt)(nil)
)

// MarshalJSON converts JsonInt into string and marshal that.
func (an JsonInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(an)))
}

func (an *JsonInt) UnmarshalJSON(input []byte) error {
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		var v int
		if err := json.Unmarshal(input, &v); err != nil {
			return err
		}
		*an = JsonInt(v)
		return nil
	} else {
		if v, err := strconv.ParseInt(s, 10, 64); err != nil {
			return err
		} else {
			*an = JsonInt(v)
		}
		return nil
	}
}
