package command

import (
	"context"
	"errors"
	"github.com/maxipaz/wallet/config"
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command
func NewRootCommand(ctx context.Context) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "wallet",
		Short: "run the wallet service",

		PersistentPreRunE: config.Setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("command was not provided, please specify a command: deploy, monitor or run")
		},
	}

	rootCommand.Flags().StringVarP(
		&config.Filename,
		"config",
		"f",
		"config/config.yaml",
		"Relative path to the config file",
	)

	rootCommand.PersistentFlags().StringP("blockchain.pk", "k", "", "Account private key")
	rootCommand.AddCommand(NewDeployCommand(ctx))
	rootCommand.AddCommand(NewMonitorCommand(ctx))
	rootCommand.AddCommand(NewRunnerCommand(ctx))

	return rootCommand
}
