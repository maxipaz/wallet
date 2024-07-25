package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts "github.com/maxipaz/wallet/contracts/interfaces"
	"github.com/maxipaz/wallet/internal/common"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"math/big"
	"time"
)

type Monitor struct {
	contractAddress string
}

// AllowanceChangedEvent struct
type AllowanceChangedEvent struct {
	Event       string    `json:"event_type"`
	Sender      string    `json:"sender"`
	Beneficiary string    `json:"beneficiary"`
	PrevAmount  *big.Int  `json:"prev_amount"`
	NewAmount   *big.Int  `json:"new_amount"`
	Timestamp   time.Time `json:"timestamp"`
}

// MoneyReceivedEvent struct
type MoneyReceivedEvent struct {
	Event       string    `json:"event_type"`
	Sender      string    `json:"sender"`
	BlockNumber uint64    `json:"block_number"`
	Amount      *big.Int  `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
}

// MoneySentEvent struct
type MoneySentEvent struct {
	Event       string    `json:"event_type"`
	Beneficiary string    `json:"beneficiary"`
	BlockNumber uint64    `json:"block_number"`
	Amount      *big.Int  `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
}

// OwnershipTransferredEvent struct
type OwnershipTransferredEvent struct {
	Event         string    `json:"event_type"`
	PreviousOwner string    `json:"previous_owner"`
	NewOwner      string    `json:"new_owner"`
	BlockNumber   uint64    `json:"block_number"`
	Timestamp     time.Time `json:"timestamp"`
}

// NewMonitor returns a new runner instance
func NewMonitor(contractAddress string) *Monitor {
	return &Monitor{
		contractAddress: contractAddress,
	}
}

// Start register to listen blockchain events
func (m *Monitor) Start(ctx context.Context, client *ethclient.Client) error {
	slog.DebugContext(ctx, "start monitoring", slog.String("contract_address", m.contractAddress))

	if err := common.ValidateContractAddress(ctx, client, m.contractAddress); err != nil {
		return fmt.Errorf("failed to validate contract address: %w", err)
	}

	contract, err := contracts.NewContract(ethcommon.HexToAddress(m.contractAddress), client)
	if err != nil {
		return fmt.Errorf("failed to create contract instance: %w", err)
	}

	eg := new(errgroup.Group)
	eg.Go(func() error {
		return m.watchAllowanceChanged(ctx, contract)
	})

	eg.Go(func() error {
		return m.watchMoneySent(ctx, contract)
	})

	eg.Go(func() error {
		return m.watchMoneyReceived(ctx, contract)
	})

	eg.Go(func() error {
		return m.watchOwnershipTransferred(ctx, contract)
	})

	return eg.Wait()
}

func (m *Monitor) watchAllowanceChanged(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractAllowanceChanged)

	subscription, err := contract.WatchAllowanceChanged(&bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}, events, nil, nil)
	if err != nil {
		return err
	}

	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, err := json.MarshalIndent(
				AllowanceChangedEvent{
					Event:       "AllowanceChanged",
					Sender:      event.Sender.Hex(),
					Beneficiary: event.Beneficiary.Hex(),
					PrevAmount:  common.WeiToEther(event.PrevAmount),
					NewAmount:   common.WeiToEther(event.NewAmount),
					Timestamp:   time.Now().UTC(),
				},
				"",
				"  ",
			)
			if err != nil {
				slog.ErrorContext(ctx, "error marshaling AllowanceChanged event", slog.String("error", err.Error()))
				continue
			}

			slog.DebugContext(ctx, "AllowanceChanged event received", slog.String("event", string(j)))
		}
	}
}

func (m *Monitor) watchMoneySent(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractMoneySent)

	subscription, err := contract.WatchMoneySent(&bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}, events, nil)
	if err != nil {
		return err
	}

	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, err := json.MarshalIndent(
				MoneySentEvent{
					Event:       "MoneySent",
					Beneficiary: event.Beneficiary.Hex(),
					BlockNumber: event.Raw.BlockNumber,
					Amount:      common.WeiToEther(event.Amount),
					Timestamp:   time.Now().UTC(),
				},
				"",
				"  ",
			)
			if err != nil {
				slog.ErrorContext(ctx, "error marshaling MoneySent event", slog.String("error", err.Error()))
				continue
			}

			slog.DebugContext(ctx, "MoneySent event received", slog.String("event", string(j)))
		}
	}
}

func (m *Monitor) watchMoneyReceived(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractMoneyReceived)

	subscription, err := contract.WatchMoneyReceived(&bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}, events, nil)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, err := json.MarshalIndent(
				MoneyReceivedEvent{
					Event:       "MoneyReceived",
					Sender:      event.From.Hex(),
					BlockNumber: event.Raw.BlockNumber,
					Amount:      common.WeiToEther(event.Amount),
					Timestamp:   time.Now(),
				},
				"",
				"  ",
			)
			if err != nil {
				slog.ErrorContext(ctx, "error marshaling MoneyReceived event", slog.String("error", err.Error()))
				continue
			}

			slog.DebugContext(ctx, "MoneyReceived event received", slog.String("event", string(j)))
		}
	}
}

func (m *Monitor) watchOwnershipTransferred(ctx context.Context, contract *contracts.Contract) error {
	events := make(chan *contracts.ContractOwnershipTransferred)

	subscription, err := contract.WatchOwnershipTransferred(&bind.WatchOpts{
		Start:   nil,
		Context: ctx,
	}, events, nil, nil)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil
		case errChan := <-subscription.Err():
			return errChan
		case event := <-events:
			j, err := json.MarshalIndent(
				OwnershipTransferredEvent{
					Event:         "OwnershipTransferred",
					PreviousOwner: event.PreviousOwner.Hex(),
					NewOwner:      event.NewOwner.Hex(),
					BlockNumber:   event.Raw.BlockNumber,
					Timestamp:     time.Now(),
				},
				"",
				"  ",
			)
			if err != nil {
				slog.ErrorContext(ctx, "error marshaling OwnershipTransferred event", slog.String("error", err.Error()))
				continue
			}

			slog.DebugContext(ctx, "OwnershipTransferred event received", slog.String("event", string(j)))
		}
	}
}
