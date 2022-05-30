package main

import "github.com/shopspring/decimal"

type decimalValue decimal.Decimal

func (d *decimalValue) Type() string {
	return "decimal.Decimal"
}

func (d *decimalValue) Set(s string) error {
	if v, err := decimal.NewFromString(s); err != nil {
		return err
	} else {
		*d = decimalValue(v)
		return nil
	}
}

func (d *decimalValue) String() string {
	return (*decimal.Decimal)(d).String()
}
