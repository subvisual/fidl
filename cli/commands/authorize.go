package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
	"github.com/subvisual/fidl/types"
)

func newAuthorizeCommand(cl cli.CLI) *cobra.Command {
	opts := cli.AuthorizeOptions{}
	authorizeCmd := &cobra.Command{
		Use:   "authorize",
		Short: "To authorize a storage provider to spend a specific amount of FIL to make reedems.",
		Long:  `This command allocates funds from your wallet to be spent by a storage provider to reedem files to you.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := cl.Validate.Struct(opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg := cli.LoadConfiguration(cfgPath)

			ki, err := types.ReadWallet(cfg.Wallet)
			if err != nil {
				return fmt.Errorf("failed to read wallet: %w", err)
			}

			_, err = cli.Authorize(ki, cfg.Wallet.Address, cfg.Route.Authorize, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	authorizeCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	authorizeCmd.Flags().StringVarP(&opts.Proxy, "proxy", "p", "", "The proxy wallet address")
	cobra.CheckErr(authorizeCmd.MarkFlagRequired("bank"))
	cobra.CheckErr(authorizeCmd.MarkFlagRequired("proxy"))

	return authorizeCmd
}
