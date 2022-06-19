package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/fardream/go-dydx"
)

type commonFields struct {
	isMainnet  bool
	starkKey   dydx.StarkKey
	apiKey     dydx.ApiKey
	ethAddress string
	timeout    duration
}

func (cmn *commonFields) setupCommonFields(c *cobra.Command) {
	c.PersistentFlags().BoolVar(&cmn.isMainnet, "mainnet", false, "turn on mainnet endpoint")
	c.PersistentFlags().Var(&cmn.starkKey, "stark", "path to stark key")
	c.PersistentFlags().Var(&cmn.apiKey, "api", "path to api key")
	c.PersistentFlags().StringVar(&cmn.ethAddress, "eth-address", "", "eth address")

	cmn.timeout = duration(time.Second * 15)
	c.PersistentFlags().Var(&cmn.timeout, "time-out", "time out for all requests")

	c.MarkPersistentFlagRequired("eth-address")
	c.MarkPersistentFlagRequired("stark")
	c.MarkPersistentFlagRequired("api")

	c.MarkPersistentFlagFilename("stark")
	c.MarkPersistentFlagFilename("api")
}

func (c *commonFields) getDydxClient() (*dydx.Client, error) {
	return dydx.NewClient((*dydx.StarkKey)(&c.starkKey), (*dydx.ApiKey)(&c.apiKey), c.ethAddress, c.isMainnet)
}

type orderIds struct {
	clientId string
	orderId  string
}

func (o *orderIds) setupOrderIds(c *cobra.Command) {
	c.Flags().StringVar(&o.clientId, "client-id", "", "client id")
	c.Flags().StringVar(&o.orderId, "order-id", "", "order id")
	c.MarkFlagsMutuallyExclusive("client-id", "order-id")
}

func getOrPanic[T any](input T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return input
}

func orPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printOrPanic(input any) {
	fmt.Println(string(getOrPanic(json.MarshalIndent(input, "", "  "))))
}

func defaultLoopPrinter[T any](v *dydx.ChannelResponse[T]) {
	printOrPanic(v.Contents)
}

func runLoop[T any](sub func(context.Context, chan<- *dydx.ChannelResponse[T]) error, length time.Duration, printer func(*dydx.ChannelResponse[T])) {
	timeout_ctx, timeout_cancel := context.WithTimeout(context.Background(), length)
	defer timeout_cancel()

	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := signal.NotifyContext(timeout_ctx, syscall.SIGINT)
	defer cancel()

	wg.Add(1)
	outputs := make(chan *dydx.ChannelResponse[T])
	go func() {
		defer wg.Done()
		defer close(outputs)
		orPanic(sub(ctx, outputs))
	}()

sigloop:
	for v := range outputs {
		if v.Contents == nil {
			continue sigloop
		}
		printer(v)
	}
}
