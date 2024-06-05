# Wallet Allowance Manager

The **Wallet Allowance Manager** is a project built on the Ethereum blockchain, utilizing the `go-ethereum` library to manage and control allowances within a distributed budget. 

This project aims to provide a secure and transparent way to allocate, track, and manage funds across different accounts, ensuring efficient budget management within decentralized teams or groups.

## Features

- **Allowance Allocation**: Set up allowances for different accounts, specifying the maximum amount they can spend.
- **Secure Transactions**: Utilize Ethereum's blockchain technology for secure and transparent transactions.
- **Budget Tracking**: Keep track of how allowances are being spent and adjust them as necessary.
- **Easy Integration**: Designed to be easily integrated into existing financial management systems.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.22 or later)
- [Ethereum Node](https://ethereum.org/en/developers/docs/nodes-and-clients/) (Geth or similar)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/maxipaz/wallet.git
```

Navigate to the project directory:

```bash
cd wallet
```

Build the project (ensure Go is properly installed):

```bash
go build
```

### Usage
To start managing allowances, first, ensure your Ethereum node is running and accessible.

#### Set Allowance

```bash
./wallet set-allowance --account=0xACCOUNT_ADDRESS --amount=AMOUNT
```

#### Check Allowance

```bash
./wallet check-allowance --account=0xACCOUNT_ADDRESS
```

#### Send Funds

```bash
./wallet send --from=0xFROM_ACCOUNT --to=0xTO_ACCOUNT --amount=AMOUNT
```

> Replace `0xACCOUNT_ADDRESS`, `0xFROM_ACCOUNT`, `0xTO_ACCOUNT`, and `AMOUNT` with actual account addresses and amount values.
