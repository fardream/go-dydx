package dydx

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/apd/v3"
)

type Decimal struct {
	apd.Decimal
}

func NewDecimalFromString(s string) (*Decimal, error) {
	d, _, err := apd.NewFromString(s)
	return &Decimal{Decimal: *d}, err
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &d.Decimal)
	if err != nil {
		if _, _, err1 := d.Decimal.SetString(string(data)); err1 != nil {
			return fmt.Errorf("failed to parse decimal %s - both as str %v or int %v", string(data), err, err1)
		}
	}

	return nil
}

func (d *Decimal) Set(s string) error {
	_, _, err := d.Decimal.SetString(s)
	return err
}

func (d *Decimal) Type() string {
	return "dydx.Decimal (wrapping github.com/cockroachdb/apd/v3 Decimal)"
}
