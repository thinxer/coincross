package gocoins

import (
	"fmt"
)

type TradeError struct {
	Message string
}

func (t TradeError) Error() string {
	return fmt.Sprintf("Trade Error: %s", t.Message)
}
