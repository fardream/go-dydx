package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type lsPublicCmd struct {
	*cobra.Command
	isMainnet    bool
	market       string
	sub          bool
	timeout      duration
	sublength    duration
	orderbookTop bool
	outputFile   string

	orderbook *cobra.Command
	markets   *cobra.Command
	trades    *cobra.Command
}

func newLsPublicCmd() *lsPublicCmd {
	c := &lsPublicCmd{
		Command: &cobra.Command{
			Use:   "ls-pub",
			Short: "ls/sub to a public data feed (markets/trades/orderbook)",
			Long: `list/subscribe to a public data feed for a certain time period.

- Use Ctrl-C to cancel the subscription.
- After the time period has elapsed, the subscription will stop.
- Will **NOT** reconnect.
`,
		},
		orderbook: &cobra.Command{
			Use:   "orderbook",
			Short: "list/subscribe to orderbook",
		},
		markets: &cobra.Command{
			Use:   "markets",
			Short: "list/subscribe to markets",
		},
		trades: &cobra.Command{
			Use:   "trades",
			Short: "list/subscribe to trades",
		},
	}

	c.PersistentFlags().BoolVar(&c.isMainnet, "mainnet", false, "set to use the mainnet")
	c.PersistentFlags().StringVarP(&c.market, "market", "m", "", "market to subscribe to")
	c.PersistentFlags().BoolVar(&c.sub, "sub", false, "get the data once and quit, don't subscribe")

	c.timeout = duration(time.Second * 15)
	c.PersistentFlags().Var(&c.timeout, "time-out", "time out for all requests.")

	c.sublength = duration(time.Hour * 24)
	c.Flags().Var(&c.sublength, "subscribe-length", "how long to subscribe to")

	c.orderbook.Flags().BoolVar(&c.orderbookTop, "top", false, "show order book top instead of the data")
	c.orderbook.Flags().StringVarP(&c.outputFile, "out", "o", "", "dump messages into a directory")
	c.orderbook.MarkFlagFilename("out", "json")

	c.orderbook.Run = c.doOrderbook
	c.markets.Run = c.doMarkets
	c.trades.Run = c.doTrades

	c.AddCommand(c.orderbook, c.markets, c.trades)

	return c
}

func (c *lsPublicCmd) doOrderbook(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient(nil, nil, "", c.isMainnet))
	if !c.sub {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
		defer cancel()
		printOrPanic(getOrPanic(client.GetOrderbook(ctx, c.market)))
	} else {
		printer := defaultLoopPrinter[dydx.OrderbookChannelResponseContents]
		var ob *dydx.OrderbookProcessor
		if c.orderbookTop {
			ob = dydx.NewOrderbookProcessor(c.market, false)
			printer = func(v *dydx.OrderbookChannelResponse) {
				ob.Process(v)
				bid, ask := ob.BookTop()
				bidstr := "bid : empty"
				askstr := "empty : ask"
				if bid != nil {
					bidstr = fmt.Sprintf("bid : $%s %s", bid.Price.String(), bid.Size.String())
				}
				if ask != nil {
					askstr = fmt.Sprintf("$%s %s : ask", ask.Price.String(), ask.Size.String())
				}
				log.Printf("%s || %s", bidstr, askstr)
			}
		}
		runLoop(func(ctx context.Context, outputs chan<- *dydx.OrderbookChannelResponse) error {
			return client.SubscribeOrderbook(ctx, c.market, outputs)
		}, time.Duration(c.sublength), printer)
		if ob != nil && c.outputFile != "" {
			os.WriteFile(c.outputFile, getOrPanic(json.Marshal(ob.Data)), 0o666)
		}
	}
}

func (c *lsPublicCmd) doMarkets(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient(nil, nil, "", c.isMainnet))
	if !c.sub {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
		defer cancel()
		printOrPanic(getOrPanic(client.GetMarkets(ctx)))
	} else {
		runLoop(func(ctx context.Context, outputs chan<- *dydx.MarketsChannelResponse) error {
			return client.SubscribeMarkets(ctx, outputs)
		}, time.Duration(c.sublength), defaultLoopPrinter[dydx.MarketsChannelResponseContents])
	}
}

func (c *lsPublicCmd) doTrades(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient(nil, nil, "", c.isMainnet))
	if !c.sub {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
		defer cancel()
		printOrPanic(getOrPanic(client.GetTrades(ctx, &dydx.TradesParam{MarketID: c.market})))
	} else {
		runLoop(func(ctx context.Context, outputs chan<- *dydx.TradesChannelResponse) error {
			return client.SubscribeTrades(ctx, c.market, outputs)
		}, time.Duration(c.sublength), defaultLoopPrinter[dydx.TradesChannelResponseContents])
	}
}
