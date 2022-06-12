package dydx_test

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-dydx"
)

var _ = "keep"

// An example with subscriptions
func Example() {
	const ethAddress = "<eth address>"
	client, err := dydx.NewClient(
		dydx.NewStarkKey(
			"<stark pubkey>",
			"<stark pubkey y coordinate>",
			"<stark private key>"),
		dydx.NewApiKey(
			"<api key>",
			"<api passphrase>",
			"<api secret>"),
		ethAddress,
		false)
	if err != nil {
		panic(err)
	}

	orderbook := make(chan *dydx.OrderbookChannelResponse)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	go func() {
		if err := client.SubscribeOrderbook(ctx, "BTC-USD", orderbook); err != nil {
			fmt.Printf("sub order error: %#v\n", err)
		}
		close(orderbook)
	}()

	for v := range orderbook {
		spew.Printf("%#v\n", v)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	trades := make(chan *dydx.TradesChannelResponse)
	go func() {
		if err := client.SubscribeTrades(ctx, "BTC-USD", trades); err != nil {
			fmt.Printf("sub trades error: %v\n", err)
		}
		close(trades)
	}()

	for v := range trades {
		spew.Printf("%v\n", v)
	}

	markets := make(chan *dydx.MarketsChannelResponse)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	go func() {
		if err := client.SubscribeMarkets(ctx, markets); err != nil {
			fmt.Printf("sub order error: %#v\n", err)
		}
		close(markets)
	}()

	for v := range markets {
		spew.Printf("%#v\n", v)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	accounts := make(chan *dydx.AccountChannelResponse)
	go func() {
		if err := client.SubscribeAccount(ctx, 0, accounts); err != nil {
			fmt.Printf("sub trades error: %v\n", err)
		}
		close(accounts)
	}()

	for v := range accounts {
		spew.Printf("%v\n", v)
	}
}
