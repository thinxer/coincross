package gocoins

import (
	"log"
	"time"
)

// This function is used for tailing trading data of a Client
// by repeatly calling History function.
// Trades are returned to the chan Trade t.
func Tail(c Client, pair Pair, interval time.Duration, t chan Trade) {
	var (
		tid, since int64 = -1, -1
		trades     []Trade
		err        error
	)
	fib := newFibonacci()
	for {
		start := time.Now()
		trades, since, err = c.History(pair, since)
		dur := time.Now().Sub(start)
		if err == nil {
			for _, tx := range trades {
				if tx.Id > tid {
					tid = tx.Id
					t <- tx
				}
			}
			fib.Prev()
		} else {
			backoff := fib.Next()
			log.Printf("Error getting history: %s", err.Error())
			log.Printf("Waiting for %ds before retrying...", backoff)
			time.Sleep(time.Duration(backoff) * time.Second)
		}
		time.Sleep(interval - dur)
	}
}
