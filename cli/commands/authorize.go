package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newAuthorizeCommand(cfg cli.Config) *cobra.Command {
	opts := cli.AuthorizeOptions{}
	authorizeCmd := &cobra.Command{
		Use:   "authorize",
		Short: "To authorize a storage provider to spend a specific amount of FIL to make reedems.",
		Long: `This command allocates funds from your wallet to be spent by a storage provider to reedem files to you. For example:
	`,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cli.Authorize(cfg, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}
	authorizeCmd.Flags().StringVarP(&opts.Amount, "amount", "a", "", "The amount of funds to transfer")
	authorizeCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	cobra.CheckErr(authorizeCmd.MarkFlagRequired("amount"))
	cobra.CheckErr(authorizeCmd.MarkFlagRequired("bank"))

	return authorizeCmd
}
