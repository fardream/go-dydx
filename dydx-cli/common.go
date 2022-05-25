package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type (
	starkKey dydx.StarkKey
	apiKey   dydx.ApiKey
)

func (c *starkKey) String() string {
	return "empty"
}

func (c *starkKey) Set(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m, err := dydx.ParseStarkKeyMap(data)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("only one keys is allowed: %s", data)
	}
	for _, v := range m {
		*c = (starkKey)(*v)
	}
	return nil
}

func (c *starkKey) Type() string {
	return "stark-key-map-file"
}

func (c *apiKey) String() string {
	return "empty"
}

func (c *apiKey) Set(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m, err := dydx.ParseApiKeyMap(data)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("only one keys is allowed: %s", data)
	}
	for _, v := range m {
		*c = (apiKey)(*v)
	}
	return nil
}

func (c *apiKey) Type() string {
	return "api-key-map-file"
}

type commonFields struct {
	isMainnet  bool
	starkKey   starkKey
	apiKey     apiKey
	ethAddress string
	timeout    duration
}

func (cmn *commonFields) setupCommonFields(c *cobra.Command) {
	c.Flags().BoolVar(&cmn.isMainnet, "mainnet", false, "turn on mainnet endpoint")
	c.Flags().Var(&cmn.starkKey, "stark", "path to stark key")
	c.Flags().Var(&cmn.apiKey, "api", "path to api key")
	c.Flags().StringVar(&cmn.ethAddress, "eth-address", "", "eth address")
	c.Flags().Var(&cmn.timeout, "time-out", "time out for all requests")
	cmn.timeout = duration(time.Second * 15)
	c.MarkFlagRequired("eth-address")
	c.MarkFlagRequired("stark")
	c.MarkFlagRequired("api")
	c.MarkFlagFilename("stark")
	c.MarkFlagFilename("api")
}

type duration time.Duration

func (d *duration) Type() string {
	return "time.Duration"
}

func (d *duration) Set(s string) error {
	ds, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = (duration)(ds)
	return nil
}

func (d *duration) String() string {
	return (time.Duration)(*d).String()
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
