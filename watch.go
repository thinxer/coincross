package gocoins

import (
	"log"
	"time"
)

func Watch(c Client, pair Pair, interval time.Duration, t chan Trade) {
	var tid int64
	tid = -1
	for {
		start := time.Now()
		trades, err := c.History(pair, tid)
		dur := time.Now().Sub(start)
		if err == nil {
			for _, tx := range trades {
				t <- tx
			}
			if len(trades) > 0 {
				tid = trades[len(trades)-1].Id
			}
		} else {
			log.Print(err.Error())
			close(t)
			break
		}
		time.Sleep(interval - dur)
	}
}
