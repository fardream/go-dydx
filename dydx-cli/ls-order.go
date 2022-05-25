package main

import (
	"context"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type lsOrderCmd struct {
	*cobra.Command
	commonFields
	orderIds
}

func newLsOrderCmd() *lsOrderCmd {
	c := &lsOrderCmd{
		Command: &cobra.Command{
			Use:   "ls-order",
			Short: "get orders",
		},
		commonFields: commonFields{},
	}
	c.setupCommonFields(c.Command)
	c.setupOrderIds(c.Command)
	c.Run = c.do

	return c
}

func (c *lsOrderCmd) do(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient((*dydx.StarkKey)(&c.starkKey), (*dydx.ApiKey)(&c.apiKey), c.ethAddress, c.isMainnet))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()
	switch {
	case c.clientId != "":
		printOrPanic(getOrPanic(client.GetOrderByClientId(ctx, c.clientId)).Order)
	case c.orderId != "":
		printOrPanic(getOrPanic(client.GetOrderById(ctx, c.orderId)).Order)
	default:
		printOrPanic(getOrPanic(client.GetOrders(ctx, nil)).Orders)
	}
}
