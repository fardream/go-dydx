// dydx-replay-orderbook is a cli to replay orderbook updates from dydx websocket subscription.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fardream/go-dydx"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	*cobra.Command

	printBook   bool
	printUpdate bool
}

func newRootCmd() *rootCmd {
	c := &rootCmd{
		Command: &cobra.Command{
			Use:   "dydx-replay-orderbook",
			Short: "replay dydx orderbook events",
			Args:  cobra.ExactArgs(1),
		},
	}

	c.Flags().BoolVarP(&c.printBook, "book", "b", false, "print book")
	c.Flags().BoolVarP(&c.printUpdate, "update", "u", false, "print update")

	c.Run = c.do
	return c
}

func getOrPanic[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func orPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c *rootCmd) do(cmd *cobra.Command, args []string) {
	log.Printf("reading messages from %s", args[0])
	var data []*dydx.OrderbookChannelResponse
	orPanic(json.Unmarshal(getOrPanic(os.ReadFile(args[0])), &data))
	ob := dydx.NewOrderbookProcessor("BTC-USD", true)
	for i, v := range data {
		if c.printUpdate {
			log.Printf("index %d resp: %s", i, getOrPanic(json.MarshalIndent(v, "", "  ")))
		}
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
		log.Printf("index %d: %s || %s", i, bidstr, askstr)
		if c.printBook {
			log.Printf("index %d bids: %s", i, getOrPanic(json.MarshalIndent(ob.Bids, "", "  ")))
			log.Printf("index %d asks: %s", i, getOrPanic(json.MarshalIndent(ob.Asks, "", "  ")))
		}
	}
}

func main() {
	newRootCmd().Execute()
}
