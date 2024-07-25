package common

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/maxipaz/wallet/config"
	contracts "github.com/maxipaz/wallet/contracts/interfaces"
	errs "github.com/maxipaz/wallet/internal/errors"
	"math/big"
	"regexp"
)

// GetSigner get the signer for sign transactions
func GetSigner(ctx context.Context, client *ethclient.Client) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(config.App.Blockchain.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hex to ECDSA: %w", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errs.ErrInvalidKey
	}

	address := crypto.PubkeyToAddress(*publicKey)

	nonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce at: %s: %w", address, err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested gas price: %w", err)
	}

	signer.Nonce = big.NewInt(int64(nonce))
	signer.Value = big.NewInt(config.App.Contract.DefaultWeiFounds)
	signer.GasLimit = 0 // automatically estimates gas limit
	signer.GasPrice = gasPrice

	return signer, nil
}

// GetContract get an instance of the deployed contract
func GetContract(ctx context.Context, client *ethclient.Client, contractAddress string) (*contracts.Contract, error) {
	err := ValidateContractAddress(ctx, client, contractAddress)
	if err != nil {
		return nil, err
	}
	contract, err := contracts.NewContract(common.HexToAddress(contractAddress), client)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

// ValidateContractAddress validate the contract address checking if the contract is deployed
func ValidateContractAddress(ctx context.Context, client *ethclient.Client, address string) error {
	if err := ValidateAddress(address); err != nil {
		return err
	}

	contractAddress := common.HexToAddress(address)
	bytecode, err := client.CodeAt(ctx, contractAddress, nil)
	if err != nil {
		return err
	}

	if len(bytecode) == 0 {
		return errs.ErrInvalidContractAddress
	}

	return nil
}

// ValidateAddress validate address format
func ValidateAddress(address string) error {
	regex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if ok := regex.MatchString(address); !ok {
		return errs.ErrInvalidAddress
	}
	return nil
}

// EtherToWei convert Ether to Wei
func EtherToWei(eth *big.Int) *big.Int {
	return new(big.Int).Mul(eth, big.NewInt(params.Ether))
}

// WeiToEther convert Wei to Ether
func WeiToEther(wei *big.Int) *big.Int {
	return new(big.Int).Div(wei, big.NewInt(params.Ether))
}
