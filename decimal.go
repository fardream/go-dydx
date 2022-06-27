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

func DecimalToString(d *Decimal) string {
	return d.Decimal.String()
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &d.Decimal)
	if err != nil {
		var v int64
		if err1 := json.Unmarshal(data, &v); err1 != nil {
			return fmt.Errorf("failed to parse decimal %s - both as str %v or int %v", string(data), err, err1)
		}
		d.Decimal.SetInt64(v)
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
