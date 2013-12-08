package gocoins

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Returns BTC/USD
func (p Pair) String() string {
	return string(p.Target + "/" + p.Base)
}

// Returns btc_usd
func (p Pair) LowerString() string {
	return strings.ToLower(string(p.Target + "_" + p.Base))
}

// Marshal to "btc_usd"
func (p *Pair) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.LowerString())
}

// Unmarshal from "btc_usd"
func (p *Pair) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err == nil {
		parts := strings.Split(strings.ToUpper(string(b)), "_")
		*p = Pair{Symbol(parts[1]), Symbol(parts[0])}
	}
	return
}

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
