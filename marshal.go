package gocoins

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (t *TradeType) MarshalJSON() ([]byte, error) {
	var s string
	switch *t {
	case Buy:
		s = "buy"
	case Sell:
		s = "sell"
	}
	return json.Marshal(s)
}
func (t *TradeType) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err == nil {
		switch strings.ToLower(s) {
		case "buy", "bid":
			*t = Buy
		case "sell", "ask":
			*t = Sell
		default:
			return fmt.Errorf("Unknown TradeType: %v", *t)
		}
	}
	return
}

func (t TradeType) String() string {
	switch t {
	case Sell:
		return "Sell"
	case Buy:
		return "Buy"
	default:
		return ""
	}
}
