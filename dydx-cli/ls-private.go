package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type lsPrivateCmd struct {
	*cobra.Command
	*commonFields
	market string

	orderIds // for orders

	subaccount bool // for account
	sublength  duration

	orders       *cobra.Command
	activeOrders *cobra.Command
	accounts     *cobra.Command
	positions    *cobra.Command
	fills        *cobra.Command
}

func newLsPrivateCmd() *lsPrivateCmd {
	c := &lsPrivateCmd{
		Command: &cobra.Command{
			Use:   "ls",
			Short: "list private information/subscribe to accounts",
		},
		commonFields: &commonFields{},

		orders: &cobra.Command{
			Use:   "orders",
			Short: "list orders",
		},
		fills: &cobra.Command{
			Use:   "fills",
			Short: "list fills",
		},
		positions: &cobra.Command{
			Use:   "positions",
			Short: "list positions",
		},
		accounts: &cobra.Command{
			Use:   "accounts",
			Short: "list/subscribe accounts",
		},
		activeOrders: &cobra.Command{
			Use:   "active-orders",
			Short: "list active orders",
		},
	}
	c.setupCommonFields(c.Command)

	c.PersistentFlags().StringVarP(&c.market, "market", "m", "", "market")

	c.accounts.Flags().BoolVarP(&c.subaccount, "sub", "s", false, "subscribe to the account feed")
	c.sublength = duration(time.Hour * 24)
	c.accounts.Flags().Var(&c.sublength, "sub-length", "subscribe length")
	c.accounts.Run = c.doAccounts

	c.setupOrderIds(c.orders)
	c.orders.Run = c.doOrders

	c.positions.Run = c.doPositions

	c.fills.Run = c.doFills

	c.activeOrders.Run = c.doActiveOrders

	c.AddCommand(c.orders, c.fills, c.accounts, c.positions, c.activeOrders)

	return c
}

func (c *lsPrivateCmd) doOrders(*cobra.Command, []string) {
	client := getOrPanic(c.getDydxClient())
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

func (c *lsPrivateCmd) doActiveOrders(*cobra.Command, []string) {
	client := getOrPanic(c.getDydxClient())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()
	if c.market == "" {
		orPanic(fmt.Errorf("market is required to query active orders"))
	}
	printOrPanic(getOrPanic(client.GetActiveOrders(ctx, &dydx.QueryActiveOrdersParam{Market: c.market})))
}

func (c *lsPrivateCmd) doAccounts(*cobra.Command, []string) {
	client := getOrPanic(c.getDydxClient())
	if !c.subaccount {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
		defer cancel()
		printOrPanic(getOrPanic(client.GetAccounts(ctx)).Accounts)
	} else {
		runLoop(func(ctx context.Context, outputs chan<- *dydx.AccountChannelResponse) error {
			return client.SubscribeAccount(ctx, 0, outputs)
		}, time.Duration(c.sublength), defaultLoopPrinter[dydx.AccountChannelResponseContents])
	}
}

func (c *lsPrivateCmd) doPositions(*cobra.Command, []string) {
	client := getOrPanic(c.getDydxClient())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()

	printOrPanic(getOrPanic(client.GetPositions(ctx, &dydx.PositionParams{Market: c.market})).Positions)
}

func (c *lsPrivateCmd) doFills(*cobra.Command, []string) {
	client := getOrPanic(c.getDydxClient())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()
	printOrPanic(getOrPanic(client.GetFills(ctx, &dydx.FillsParam{Market: c.market})).Fills)
}
