package main

import "github.com/spf13/cobra"

type rootCmd struct {
	*cobra.Command
}

func newRootCmd() *rootCmd {
	c := &rootCmd{
		Command: &cobra.Command{
			Use:   "dydx-cli",
			Short: "cli for dydx.exchange",
		},
	}

	c.Run = c.do

	send := newSendCmd()
	c.AddCommand(send.Command)
	getCmd := newLsOrderCmd()
	c.AddCommand(getCmd.Command)
	cancelCmd := newCancelCmd()
	c.AddCommand(cancelCmd.Command)

	return c
}

func (c *rootCmd) do(*cobra.Command, []string) {
	c.Command.Help()
}
