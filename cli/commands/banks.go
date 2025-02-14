package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newBanksCommand(cl cli.CLI) *cobra.Command {
	opts := cli.BanksOptions{}
	banksCmd := &cobra.Command{
		Use:   "banks",
		Short: "To list the banks that a given proxy trusts.",
		Long:  `This command lists all the banks that a given proxy trusts and accepts to process payment for a retrieval.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := cl.Validate.Struct(opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			cfgPath, _ := cmd.Flags().GetString("config")
			cfg := cli.LoadConfiguration(cfgPath)

			_, err = cli.Banks(cfg.Route.Banks, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	banksCmd.Flags().StringVarP(&opts.ProxyAddress, "proxy", "p", "", "The proxy address")
	cobra.CheckErr(banksCmd.MarkFlagRequired("proxy"))

	return banksCmd
}
