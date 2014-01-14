package gocoins

import (
	"log"
	"time"
)

// This function is used for tailing trading data of a Client
// by repeatly calling History function.
// Trades are returned to the chan Trade t.
func Tail(c Client, pair Pair, since int64, interval time.Duration, t chan Trade) {
	var (
		tid     int64 = -1
		trades  []Trade
		err     error
		backoff = time.Second
	)
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
			if backoff > time.Second {
				backoff = backoff / 2
			}
		} else {
			log.Printf("Error getting history: %s", err.Error())
			log.Printf("Waiting for %v before retrying...", backoff)
			time.Sleep(backoff)
			backoff *= 2
		}
		time.Sleep(interval - dur)
	}
}
