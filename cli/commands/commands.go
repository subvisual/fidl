package commands

import (
	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func Parse(cfg cli.Config) *cobra.Command {
	rootCmd := &cobra.Command{Use: "fidl"}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(newDepositCommand(cfg))
	rootCmd.AddCommand(newWithdrawCommand(cfg))
	rootCmd.AddCommand(newBalanceCommand(cfg))

	return rootCmd
}
