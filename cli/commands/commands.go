package commands

import (
	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func Parse(cli cli.CLI) *cobra.Command {
	var cfgPath string
	rootCmd := &cobra.Command{
		Use:   "fidl",
		Short: "FIDL is a CLI tool to retrieve files from Filecoin storage providers.",
		Long:  "FIDL is a CLI tool to retrieve files from Filecoin storage providers via HTTP. To get the data, the client must deposit funds in a bank. Then, the storages providers can perform payments in banks they trust.",
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "./cli.ini", "Path to the configuration file")

	rootCmd.AddCommand(newDepositCommand())
	rootCmd.AddCommand(newWithdrawCommand(cli))
	rootCmd.AddCommand(newBalanceCommand())
	rootCmd.AddCommand(newBanksCommand())
	rootCmd.AddCommand(newAuthorizeCommand(cli))
	rootCmd.AddCommand(newRefundCommand())
	rootCmd.AddCommand(newRetrievalCommand())

	return rootCmd
}
