package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newBalanceCommand(cfg cli.Config) *cobra.Command {
	opts := cli.BalanceOptions{}
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "To check the client's account balance at a specified bank",
		Long: `This command checks the client's account balance at a specified bank:
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cli.Balance(cfg, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	balanceCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	cobra.CheckErr(balanceCmd.MarkFlagRequired("bank"))

	return balanceCmd
}
