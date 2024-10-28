package commands

import "github.com/spf13/cobra"

func Parse(bankAddress string) *cobra.Command {
	rootCmd := &cobra.Command{Use: "fidl"}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(newDepositCommand(bankAddress))
	rootCmd.AddCommand(newWithdrawCommand(bankAddress))

	return rootCmd
}
