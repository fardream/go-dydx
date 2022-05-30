package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"time"

	"github.com/fardream/go-dydx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

type sendCmd struct {
	*cobra.Command
	commonFields
	duration
	size       decimalValue
	orderType  string
	price      decimalValue
	clientId   string
	market     string
	side       string
	tif        string
	limitfee   decimalValue
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
	c.Flags().VarP(&c.size, "size", "s", "order size")
	c.MarkFlagRequired("size")
	c.Flags().VarP(&c.price, "price", "p", "price for the order")
	c.MarkFlagRequired("price")
	c.Flags().StringVar(&c.orderType, "order-type", "MARKET", "order type")
	c.Flags().StringVar(&c.clientId, "client-id", "", "set an optional client order id. if unset, will be automatically generated")
	c.Flags().StringVar(&c.market, "market", "m", "market for this order")
	c.MarkFlagRequired("market")
	c.limitfee.Set("0.125")
	c.Flags().Var(&c.limitfee, "limit-fee", "limit fee for this order")
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

	order := dydx.NewCreateOrderRequest(
		c.market,
		getOrPanic(dydx.GetOrderSide(c.side)),
		getOrPanic(dydx.GetOrderType(c.orderType)),
		decimal.Decimal(c.size),
		decimal.Decimal(c.price),
		c.clientId, getOrPanic(dydx.GetTimeInForce(c.tif)),
		now.Add((time.Duration)(c.duration)),
		decimal.Decimal(c.limitfee),
		c.postonly)

	printOrPanic(order)
	result := getOrPanic(client.NewOrder(ctx, order, c.positionId)).Order
	printOrPanic(result)
	if c.outputFile != "" {
		os.WriteFile(c.outputFile, getOrPanic(json.MarshalIndent(result, "", "  ")), 0o666)
	}
}
