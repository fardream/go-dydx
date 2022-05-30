package dydx_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/fardream/go-dydx"
)

//go:embed orderbook.json
var orderbook_data string

func TestOrderBookProcessor(t *testing.T) {
	ob := dydx.NewOrderbookProcessor("BTC-USD", false)
	var data []*dydx.OrderbookChannelResponse
	if err := json.Unmarshal([]byte(orderbook_data), &data); err != nil {
		t.Fatalf("failed to parse data")
	}
	for _, a := range data {
		ob.Process(a)
		bid, ask := ob.BookTop()
		t.Logf("%#v :: %#v", bid, ask)
	}
}
