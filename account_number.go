package dydx

import (
	"encoding/json"
	"strconv"
)

// AccountNumber is an int, but can be json-unmarshaled from either int or string, marshalled into string.
type AccountNumber int

var (
	_ json.Marshaler   = (*AccountNumber)(nil)
	_ json.Unmarshaler = (*AccountNumber)(nil)
)

func (an AccountNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(an)))
}

func (an *AccountNumber) UnmarshalJSON(input []byte) error {
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		var v int
		if err := json.Unmarshal(input, &v); err != nil {
			return err
		}
		*an = AccountNumber(v)
		return nil
	} else {
		if v, err := strconv.ParseInt(s, 10, 64); err != nil {
			return err
		} else {
			*an = AccountNumber(v)
		}
		return nil
	}
}
