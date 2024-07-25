package deploy

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts "github.com/maxipaz/wallet/contracts/interfaces"
	"github.com/maxipaz/wallet/internal/common"
	"log/slog"
)

type Deployer struct {
	address     ethcommon.Address
	transaction *types.Transaction
	contract    *contracts.Contract
}

// NewDeployer returns a new runner instance
func NewDeployer() *Deployer {
	return new(Deployer)
}

// Deploy deploys a new Ethereum contract
func (d *Deployer) Deploy(ctx context.Context, client *ethclient.Client) error {
	signer, err := common.GetSigner(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to get signer: %w", err)
	}

	address, tx, contract, err := contracts.DeployContract(signer, client)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %w", err)
	}

	slog.DebugContext(ctx, "waiting for contract to be deployed...", slog.String("address", address.Hex()))
	if _, err := bind.WaitDeployed(ctx, client, tx); err != nil {
		return fmt.Errorf("failed to wait deployed: %w", err)
	}

	d.address = address
	d.transaction = tx
	d.contract = contract

	return nil
}

// ContractAddress returns the contract address
func (d *Deployer) ContractAddress() string {
	return d.address.Hex()
}
