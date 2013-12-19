package gocoins

import (
	"errors"
)

type TradeError error

func NewTradeError(text string) TradeError {
	return TradeError(errors.New(text))
}

var (
	ErrInvalidCredential      = NewTradeError("Invalid Credential")
	ErrInsufficientPermission = NewTradeError("Insufficient Permissions")
	ErrInsufficientBalance    = NewTradeError("Insufficient Balance")
)
