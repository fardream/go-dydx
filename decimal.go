package dydx

import "github.com/fardream/decimal"

type Decimal = decimal.Decimal

func NewDecimalFromString(s string) (*Decimal, error) {
	return decimal.NewFromString(s)
}
