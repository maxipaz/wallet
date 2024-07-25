package api

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maxipaz/wallet/config"
	errs "github.com/maxipaz/wallet/internal/errors"
	"github.com/maxipaz/wallet/internal/wallet"
	"github.com/spf13/cobra"
	"log/slog"
)

var allowanceActions = map[string]struct{}{
	"set":      {},
	"get":      {},
	"increase": {},
	"reduce":   {},
}

// NewAllowanceCommand creates the allowance command
func NewAllowanceCommand(ctx context.Context) *cobra.Command {
	var (
		action        string
		targetAddress string
		amount        int64
	)
	allowanceCommand := &cobra.Command{
		Use:   "allowance",
		Short: "Change the allowance for a beneficiary",
		Run: func(cmd *cobra.Command, args []string) {

			if err := runAllowance(ctx, action, targetAddress, amount); err != nil {
				slog.ErrorContext(ctx, "failed to run allowance", slog.String("error", err.Error()))
			}
		},
	}

	allowanceCommand.Flags().StringVar(&action, "action", "", "Action to perform: set, get, increase or reduce")
	allowanceCommand.Flags().Int64Var(&amount, "amount", 0, "Amount")
	allowanceCommand.Flags().StringVarP(&targetAddress, "target.address", "t", "", "Target address")
	_ = allowanceCommand.MarkFlagRequired("action")
	_ = allowanceCommand.MarkFlagRequired("target.address")

	return allowanceCommand
}

func runAllowance(ctx context.Context, action string, targetAddress string, amount int64) error {
	if _, ok := allowanceActions[action]; !ok {
		return errs.ErrInvalidAllowanceAction
	}

	ctxCall, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctxCall, config.App.Blockchain.WS)
	if err != nil {
		return err
	}

	defer client.Close()

	runner := wallet.NewAllowanceRunner(config.App.Blockchain.PrivateKey, config.App.Contract.Address)

	switch action {
	case wallet.GetAction:
		allowance, err := runner.GetAllowance(ctx, client, targetAddress)
		if err != nil {
			return fmt.Errorf("failed to get allowance: %w", err)
		}

		fmt.Printf("Current allowance for address %s is %d\n", targetAddress, allowance)
	default:
		if amount <= 0 {
			return errs.ErrInvalidAmountAction
		}

		if err := runner.ChangeAllowance(ctx, client, action, targetAddress, amount); err != nil {
			return fmt.Errorf("failed to change allowance: %w", err)
		}
	}

	return nil
}
