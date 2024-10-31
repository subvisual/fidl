package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newDepositCommand(cfg cli.Config) *cobra.Command {
	opts := cli.DepositOptions{}
	depositCmd := &cobra.Command{
		Use:   "deposit",
		Short: "To deposit FIL into the client's bank account.",
		Long: `This command transfers a specified amount of FIL from the client's wallet to their
	account in the bank's system, securely updating the client's balance within the service. For example:
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cli.Deposit(cfg, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}
	depositCmd.Flags().StringVarP(&opts.Amount, "amount", "a", "", "The amount of funds to transfer")
	cobra.CheckErr(depositCmd.MarkFlagRequired("amount"))

	return depositCmd
}
