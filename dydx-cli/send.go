package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type sendCmd struct {
	*cobra.Command
	commonFields
	duration
	size       float64
	orderType  string
	price      float64
	clientId   string
	market     string
	side       string
	tif        string
	limitfee   string
	postonly   bool
	positionId int64
	outputFile string
}

func newSendCmd() *sendCmd {
	c := &sendCmd{
		Command: &cobra.Command{
			Use:   "send",
			Short: "send order",
			Long:  "send an order to dydx.exchange",
		},
		commonFields: commonFields{},
		duration:     duration(time.Minute * 15),
	}

	c.setupCommonFields(c.Command)

	c.Flags().VarP(&c.duration, "duration", "t", "order duration")
	c.Flags().Float64VarP(&c.size, "size", "s", 0, "order size")
	c.MarkFlagRequired("size")
	c.Flags().Float64VarP(&c.price, "price", "p", 0, "price for the order")
	c.MarkFlagRequired("price")
	c.Flags().StringVar(&c.orderType, "order-type", "MARKET", "order type")
	c.Flags().StringVar(&c.clientId, "client-id", "", "set an optional client order id. if unset, will be automatically generated")
	c.Flags().StringVar(&c.market, "market", "m", "market for this order")
	c.MarkFlagRequired("market")
	c.Flags().StringVar(&c.limitfee, "limit-fee", "0.1", "limit fee for this order")
	c.Flags().StringVar(&c.side, "side", "", "side")
	c.MarkFlagRequired("side")
	c.Flags().StringVar(&c.tif, "tif", "GTT", "time-in-force")
	c.Flags().BoolVar(&c.postonly, "post-only", false, "post only")
	c.Flags().Int64Var(&c.positionId, "position-id", 0, "position id")
	c.MarkFlagRequired("position-id")
	c.Flags().StringVarP(&c.outputFile, "output", "o", "", "output result to this file")
	c.MarkFlagFilename("output")
	c.Run = c.do

	return c
}

func (c *sendCmd) do(*cobra.Command, []string) {
	client := getOrPanic(dydx.NewClient((*dydx.StarkKey)(&c.starkKey), (*dydx.ApiKey)(&c.apiKey), c.ethAddress, c.isMainnet))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()

	now := time.Now()
	if c.clientId == "" {
		million := big.NewInt(1_000_000_000)
		id, _ := rand.Int(rand.Reader, million)
		date_yyyymmdd := big.NewInt((int64)(now.Year()*10000 + int(now.Month())*100 + now.Day()))
		c.clientId = date_yyyymmdd.Mul(date_yyyymmdd, million).Add(date_yyyymmdd, id).String()
	}

	expiration := dydx.GetIsoDateStr(now.Add((time.Duration)(c.duration)))

	order := dydx.NewCreateOrderRequest(c.market, c.side, c.orderType, fmt.Sprintf("%f", c.size), fmt.Sprintf("%f", c.price), c.clientId, c.tif, expiration, c.limitfee, c.postonly)

	printOrPanic(order)
	result := getOrPanic(client.NewOrder(ctx, order, c.positionId)).Order
	printOrPanic(result)
	if c.outputFile != "" {
		os.WriteFile(c.outputFile, getOrPanic(json.MarshalIndent(result, "", "  ")), 0o666)
	}
}
