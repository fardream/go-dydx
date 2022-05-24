# go-dydx

golang client for [dydx.exchange](https://dydx.exchange)

## Prior Art

This is based on the work from [go-numb](https://github.com/go-numb) at [here](https://github.com/go-numb/go-dydx) with some go idiomatic modifications.

There is also another version from [verichenn](https://github.com/verichenn) [here](https://github.com/verichenn/dydx-v3-go).

## Examples

### Subscriptions

Below is an example running the subscription for 15 seconds before unsubscribe and shutdown.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-dydx"
)

func main() {
	const ethAddress = "<eth address>"
	client, err := dydx.NewClient(
		dydx.NewStarkKey(
			ethAddress,
			"<stark pubkey>",
			"<stark pubkey y coordinate>",
			"<stark private key>"),
		dydx.NewApiKey(
            ethAddress,
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
```

### Place order

```go
package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-dydx"
)

func GetOrPanic[T any](input T, err error) T {
	if err != nil {
		panic(err)
	}
	return input
}

func main() {
	const ethAddress = "<eth address>"
	client, err := dydx.NewClient(
		dydx.NewStarkKey(
			ethAddress,
			"<stark pubkey>",
			"<stark pubkey y coordinate>",
			"<stark private key>"),
		dydx.NewApiKey(
            ethAddress,
            "<api key>",
            "<api passphrase>",
            "<api secret>"),
		ethAddress,
		false)
	if err != nil {
		panic(err)
	}

	// generate a big id
	id, _ := rand.Int(rand.Reader, big.NewInt(10000))
	id.Add(id, big.NewInt(1000000))

	// expiration
	expiration := dydx.GetIsoDateStr(time.Now().Add(5 * time.Minute))

	// create a new order
	order := dydx.NewCreateOrderRequest("BTC-USD", dydx.OrderSideSell, dydx.OrderTypeLimit, "0.001", "35000", id.String(), "", expiration, "0.125", false)
	// place the order
	r := GetOrPanic(client.NewOrder(context.Background(), order, 62681))
	spew.Dump(r)
	time.Sleep(5 * time.Second)

	// get the order
	spew.Dump(GetOrPanic(
		client.GetOrderById(context.Background(), r.Order.ID)).Order)
	time.Sleep(2 * time.Second)

	// cancel all active orders
	spew.Dump(client.CancelActiveOrders(context.Background(), &dydx.CancelActiveOrdersParam{Market: "BTC-USD"}))
	time.Sleep(2 * time.Second)

	// get your positions
	spew.Dump(GetOrPanic(client.GetPositions(context.Background(), &dydx.PositionParams{Market: "BTC-USD"})))
	time.Sleep(2 * time.Second)

	// get the active orders
	spew.Dump(GetOrPanic(client.GetActiveOrders(context.Background(), &dydx.QueryActiveOrdersParam{Market: "BTC-USD"})))
}
```
