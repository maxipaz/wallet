package command

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maxipaz/wallet/config"
	deploy2 "github.com/maxipaz/wallet/internal/deploy"
	"github.com/spf13/cobra"
	"log/slog"
)

// NewDeployCommand creates the deploy command
func NewDeployCommand(ctx context.Context) *cobra.Command {
	deployCommand := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy contract to blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(ctx)
		},
	}

	return deployCommand
}

func deploy(ctx context.Context) error {
	slog.DebugContext(ctx, "deploying contract")

	ctx, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.App.Blockchain.Address)
	if err != nil {
		return err
	}

	deployer := deploy2.NewDeployer()

	if err := deployer.Deploy(ctx, client); err != nil {
		return err
	}

	slog.DebugContext(ctx, "contract deployed at address %s", deployer.ContractAddress())
	return nil
}
