package dydx_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/fardream/go-dydx"
)

//go:embed data.json
var jsondata string

func TestParseOrderbook(t *testing.T) {
	var data dydx.OrderbookChannelResponse
	err := json.Unmarshal([]byte(jsondata), &data)
	if err != nil {
		t.Fatalf("failed to parse the data: %#v", err)
	}
	t.Logf("parsed data: %#v", *data.Contents)
}
