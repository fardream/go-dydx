package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
)

type testnetTokenCmd struct {
	*cobra.Command
	*commonFields
}

func newTestnetTokenCmd() *testnetTokenCmd {
	r := &testnetTokenCmd{
		Command: &cobra.Command{
			Use:   "airdrop",
			Short: "request testne tokens",
		},
		commonFields: &commonFields{},
	}

	r.Run = r.do

	r.setupCommonFields(r.Command)

	return r
}

func (c *testnetTokenCmd) do(*cobra.Command, []string) {
	client := getOrPanic(c.getDydxClient())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout))
	defer cancel()
	printOrPanic(getOrPanic(client.RequestTestnetTokens(ctx)))
}
