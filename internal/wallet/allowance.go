package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maxipaz/wallet/internal/common"
	"math/big"
)

const (
	SetAction      = "set"
	GetAction      = "get"
	IncreaseAction = "increase"
	ReduceAction   = "reduce"
)

type Allowance struct {
	privateKey      string
	contractAddress string
}

// NewAllowanceRunner returns a new runner instance
func NewAllowanceRunner(privateKey string, contractAddress string) *Allowance {
	return &Allowance{
		privateKey:      privateKey,
		contractAddress: contractAddress,
	}
}

// GetAllowance get allowance value for a given address
func (r *Allowance) GetAllowance(ctx context.Context, client *ethclient.Client, beneficiaryAddress string) (int64, error) {
	contract, err := common.GetContract(ctx, client, r.contractAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to get contract: %w", err)
	}

	address := ethcommon.HexToAddress(beneficiaryAddress)

	amount, err := contract.Allowance(&bind.CallOpts{Pending: false, Context: ctx}, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get allowance: %w", err)
	}

	return common.WeiToEther(amount).Int64(), nil
}

// ChangeAllowance change the Allowance value for a given address
func (r *Allowance) ChangeAllowance(ctx context.Context, client *ethclient.Client, action string, target string, amount int64) error {
	contract, err := common.GetContract(ctx, client, r.contractAddress)
	if err != nil {
		return fmt.Errorf("failed to get contract: %w", err)
	}

	signer, err := common.GetSigner(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to get signer: %w", err)
	}

	var tx *types.Transaction
	var txErr error
	targetAddress := ethcommon.HexToAddress(target)

	var operation string
	switch action {
	case SetAction:
		{
			tx, txErr = contract.SetAllowance(signer, targetAddress, common.EtherToWei(big.NewInt(amount)))
			operation = "set_allowance"
		}
	case IncreaseAction:
		{
			tx, txErr = contract.IncreaseAllowance(signer, targetAddress, common.EtherToWei(big.NewInt(amount)))
			operation = "increase_allowance"
		}
	case ReduceAction:
		{
			tx, txErr = contract.ReduceAllowance(signer, targetAddress, common.EtherToWei(big.NewInt(amount)))
			operation = "reduce_allowance"
		}
	}
	if txErr != nil {
		return txErr
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("receipt status unsuccessful")
	}

	processTransaction(ctx, tx, operation)

	return nil
}
