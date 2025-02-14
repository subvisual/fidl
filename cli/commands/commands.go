package commands

import (
	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func Parse(cl cli.CLI) *cobra.Command {
	var cfgPath string
	rootCmd := &cobra.Command{
		Use:   "fidl",
		Short: "FIDL is a CLI tool to retrieve files from Filecoin storage providers.",
		Long:  "FIDL is a CLI tool to retrieve files from Filecoin storage providers via HTTP. To get the data, the client must deposit funds in a bank. Then, the storages providers can perform payments in banks they trust.",
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "./cli.ini", "Path to the configuration file")

	rootCmd.AddCommand(newDepositCommand(cl))
	rootCmd.AddCommand(newWithdrawCommand(cl))
	rootCmd.AddCommand(newBalanceCommand(cl))
	rootCmd.AddCommand(newBanksCommand(cl))
	rootCmd.AddCommand(newAuthorizeCommand(cl))
	rootCmd.AddCommand(newRefundCommand(cl))
	rootCmd.AddCommand(newRetrievalCommand(cl))

	return rootCmd
}
