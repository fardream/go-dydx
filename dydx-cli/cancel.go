package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type cancelCmd struct {
	*cobra.Command
	commonFields
	orderIds
	cancelAll bool
	market    string
}

func newCancelCmd() *cancelCmd {
	c := &cancelCmd{
		Command: &cobra.Command{
			Use:   "cancel",
			Short: "cancel orders",
		},
	}
	c.setupCommonFields(c.Command)
	c.setupOrderIds(c.Command)
	c.Flags().BoolVarP(&c.cancelAll, "cancel-all", "a", false, "cancel all orders")
	c.MarkFlagsMutuallyExclusive("client-id", "order-id", "cancel-all")
	c.Flags().StringVarP(&c.market, "market", "m", "", "market to cancel the orders on")

	c.Run = c.do

	return c
}

func (c *cancelCmd) do(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient((*dydx.StarkKey)(&c.starkKey), (*dydx.ApiKey)(&c.apiKey), c.ethAddress, c.isMainnet))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()
	switch {
	case c.cancelAll:
		printOrPanic(getOrPanic(client.CancelOrders(ctx, &dydx.CancelOrdersParam{Market: c.market})).CancelOrders)
	case c.clientId != "":
		printOrPanic(getOrPanic(client.CancelOrder(ctx, getOrPanic(client.GetOrderByClientId(ctx, c.clientId)).Order.ID)).CancelOrder)
	case c.orderId != "":
		printOrPanic(getOrPanic(client.CancelOrder(ctx, c.orderId)).CancelOrder)
	default:
		orPanic(fmt.Errorf("missing market, or client id, or order id"))
	}
}
