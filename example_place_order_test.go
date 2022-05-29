package dydx_test

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-dydx"
)

func getOrPanic[T any](input T, err error) T {
	if err != nil {
		panic(err)
	}
	return input
}

// Place a new order for BTC at the price of 35000
func ExampleClient_NewOrder() {
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
	r := getOrPanic(client.NewOrder(context.Background(), order, 62681))
	spew.Dump(r)
	time.Sleep(5 * time.Second)

	// get the order
	spew.Dump(getOrPanic(
		client.GetOrderById(context.Background(), r.Order.ID)).Order)
	time.Sleep(2 * time.Second)

	// cancel all active orders
	spew.Dump(client.CancelActiveOrders(context.Background(), &dydx.CancelActiveOrdersParam{Market: "BTC-USD"}))
	time.Sleep(2 * time.Second)

	// get your positions
	spew.Dump(getOrPanic(client.GetPositions(context.Background(), &dydx.PositionParams{Market: "BTC-USD"})))
	time.Sleep(2 * time.Second)

	// get the active orders
	spew.Dump(getOrPanic(client.GetActiveOrders(context.Background(), &dydx.QueryActiveOrdersParam{Market: "BTC-USD"})))
}
