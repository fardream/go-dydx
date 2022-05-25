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
			Long:  `cli for dydx.exchange`,
		},
	}

	c.Run = c.do

	send := newSendCmd()
	getCmd := newLsPrivateCmd()
	cancelCmd := newCancelCmd()
	subCmd := newLsPublicCmd()
	c.AddCommand(
		send.Command,
		getCmd.Command,
		cancelCmd.Command,
		subCmd.Command)

	return c
}

func (c *rootCmd) do(*cobra.Command, []string) {
	c.Command.Help()
}
