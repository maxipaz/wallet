package errors

import "errors"

var (
	ErrInvalidKey             = errors.New("invalid key")
	ErrInvalidAddress         = errors.New("invalid address")
	ErrInvalidContractAddress = errors.New("invalid contract address")
	ErrInvalidAllowanceAction = errors.New("invalid allowance action")
	ErrInvalidAmountAction    = errors.New("amount should be a positive value")
	ErrInvalidTransferAction  = errors.New("invalid transfer action")
)
