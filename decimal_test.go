package dydx_test

import (
	"encoding/json"
	"testing"

	"github.com/fardream/go-dydx"
)

func TestDecimal_UnmarshalJSON(t *testing.T) {
	strInput := "\"125\""
	var d dydx.Decimal
	if err := json.Unmarshal([]byte(strInput), &d); err != nil {
		t.Fatalf("failed to parse string: %s", strInput)
	}
	intInput := "125"
	if err := json.Unmarshal([]byte(intInput), &d); err != nil {
		t.Fatalf("failed to parse int input: %s", intInput)
	}
}
