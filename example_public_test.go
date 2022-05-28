package dydx_test

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-dydx"
)

// Get the orderbook for BTC-USD
func ExampleClient_GetOrderbook() {
	// No private key necesary
	client, _ := dydx.NewClient(nil, nil, "", false)
	// Set 15 minutes timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	spew.Dump(client.GetOrderbook(ctx, "BTC-USD"))
}
