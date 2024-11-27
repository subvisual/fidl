package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/subvisual/fidl/cli"
)

func newRetrievalCommand() *cobra.Command {
	opts := cli.RetrievalOptions{}
	retrievalCmd := &cobra.Command{
		Use:   "retrieval",
		Short: "To retrieval a file from a storage provider.",
		Long:  `This command retrieval a file from a storage provider. To do so, the client must have some FIL in a bank trusted by the storage provider.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg := cli.LoadConfiguration(cfgPath)

			err := cli.Retrieval(cfg.Route.Retrieval, opts)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}
	retrievalCmd.Flags().StringVarP(&opts.Piece, "id", "i", "", "The piece CID to be retrieved")
	retrievalCmd.Flags().StringVarP(&opts.ProxyAddress, "proxy", "p", "", "The proxy address")
	retrievalCmd.Flags().StringVarP(&opts.BankAddress, "bank", "b", "", "The bank address")
	retrievalCmd.Flags().StringVarP(&opts.Authorization, "authorization", "a", "", "The bank address")
	cobra.CheckErr(retrievalCmd.MarkFlagRequired("id"))
	cobra.CheckErr(retrievalCmd.MarkFlagRequired("proxy"))
	cobra.CheckErr(retrievalCmd.MarkFlagRequired("bank"))
	cobra.CheckErr(retrievalCmd.MarkFlagRequired("authorization"))

	return retrievalCmd
}
