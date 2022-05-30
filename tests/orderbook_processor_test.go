package dydx_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/fardream/go-dydx"
	"github.com/fardream/go-dydx/heap"
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
		if !heap.IsHeap[*dydx.Bids, *dydx.OrderbookOrder](&ob.Bids) {
			t.Fatalf("bids not heap: %s\n", ob.Bids.PrintBook())
		}
		if !heap.IsHeap[*dydx.Asks, *dydx.OrderbookOrder](&ob.Asks) {
			t.Fatalf("asks not heap: %s\n", ob.Bids.PrintBook())
		}
	}
}
