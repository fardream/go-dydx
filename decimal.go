package dydx

import "github.com/shopspring/decimal"

type Decimal = decimal.Decimal

func NewDecimalFromString(s string) (Decimal, error) {
	return decimal.NewFromString(s)
}

func DecimalToString(d *Decimal) string {
	return d.String()
}
