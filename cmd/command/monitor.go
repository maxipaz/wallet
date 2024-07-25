package command

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maxipaz/wallet/config"
	"github.com/maxipaz/wallet/internal/wallet"
	"github.com/spf13/cobra"
)

// NewMonitorCommand creates the monitor command
func NewMonitorCommand(ctx context.Context) *cobra.Command {
	monitorCommand := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor events in the blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return monitoring(ctx)
		},
	}

	monitorCommand.Flags().StringP("contract.address", "c", "", "Contract address")
	return monitorCommand
}

func monitoring(ctx context.Context) error {
	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}

	return wallet.NewMonitor(config.App.Contract.Address).Start(ctx, client)
}
