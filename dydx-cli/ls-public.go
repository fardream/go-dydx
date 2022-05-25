package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type lsPublicCmd struct {
	*cobra.Command
	isMainnet bool
	market    string
	sub       bool
	timeout   duration
	sublength duration

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

	c.orderbook.Run = c.doOrderbook
	c.markets.Run = c.doMarkets
	c.trades.Run = c.doTrades

	c.AddCommand(c.orderbook, c.markets, c.trades)

	return c
}

func runLoop[T any](sub func(context.Context, chan<- *dydx.ChannelResponse[T]) error, length time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), length)
	defer cancel()
	var wg sync.WaitGroup
	defer wg.Wait()

	sig_chan := make(chan os.Signal, 5)
	signal.Notify(sig_chan, syscall.SIGINT)

	wg.Add(1)
	outputs := make(chan *dydx.ChannelResponse[T])
	go func() {
		defer wg.Done()
		defer close(outputs)
		orPanic(sub(ctx, outputs))
	}()

	sigint_called := 0
sigloop:
	for {
		select {
		case <-sig_chan:
			sigint_called++
			cancel()
			if sigint_called >= 5 {
				orPanic(fmt.Errorf("sigint called 5 times, quit"))
			}
		case v, ok := <-outputs:
			if !ok {
				break sigloop
			}
			printOrPanic(v.Contents)
		}
	}
}

func (c *lsPublicCmd) doOrderbook(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient(nil, nil, "", c.isMainnet))
	if !c.sub {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
		defer cancel()
		printOrPanic(getOrPanic(client.GetOrderbook(ctx, c.market)))
	} else {
		runLoop(func(ctx context.Context, outputs chan<- *dydx.OrderbookChannelResponse) error {
			return client.SubscribeOrderbook(ctx, c.market, outputs)
		}, time.Duration(c.sublength))
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
		}, time.Duration(c.sublength))
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
		}, time.Duration(c.sublength))
	}
}
