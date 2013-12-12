package gocoins

import (
	"fmt"
)

// The EMA statistics
type EMA struct {
	ratio float64
	last  float64
}

func MakeEMA(ratio float64) *EMA {
	return &EMA{ratio: ratio, last: -1}
}

func (e *EMA) Name() string {
	return fmt.Sprintf("EMA_%f", e.ratio)
}

func (e *EMA) Stat(t Trade) {
	if e.last < 0 {
		e.last = t.Price
	} else {
		e.last = e.last*e.ratio + t.Price*(1-e.ratio)
	}
}

func (e *EMA) Last() []float64 {
	return []float64{e.last}
}
