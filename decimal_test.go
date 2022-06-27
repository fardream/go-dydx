package dydx_test

import (
	"encoding/json"
	"testing"

	"github.com/cockroachdb/apd/v3"
	"github.com/fardream/go-dydx"
)

func TestDecimal_UnmarshalJSON(t *testing.T) {
	d125_345 := apd.New(125345, -3)
	strInput := "\"125.345\""
	var d dydx.Decimal
	if err := json.Unmarshal([]byte(strInput), &d); err != nil {
		t.Fatalf("failed to parse string: %s", strInput)
	}

	if d125_345.Cmp(&d.Decimal) != 0 {
		t.Fatalf("%#v/%s is not 125.345", d, d.String())
	}
	intInput := "125.345"
	if err := json.Unmarshal([]byte(intInput), &d); err != nil {
		t.Fatalf("failed to parse int input: %s", intInput)
	}
	if d125_345.Cmp(&d.Decimal) != 0 {
		t.Fatalf("%#v/%s is not 125.345", d, d.String())
	}
}
