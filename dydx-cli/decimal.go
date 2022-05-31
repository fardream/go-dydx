package main

import (
	"github.com/fardream/go-dydx"
	"github.com/shopspring/decimal"
)

type decimalValue dydx.Decimal

func (d *decimalValue) Type() string {
	return "decimal.Decimal"
}

func (d *decimalValue) Set(s string) error {
	if v, err := dydx.NewDecimalFromString(s); err != nil {
		return err
	} else {
		*d = decimalValue(v)
		return nil
	}
}

func (d *decimalValue) String() string {
	return dydx.DecimalToString((*decimal.Decimal)(d))
}
