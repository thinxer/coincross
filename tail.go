package coincross

import (
	"log"
	"time"
)

// Tail follows Client.History.
func Tail(c Client, pair Pair, since int64, interval time.Duration) *Streamer {
	// Advanced Go Concurrency Patterns: http://talks.golang.org/2013/advconc.slide

	var (
		trades  = make(chan Trade, 100)
		closing = make(chan bool)

		tid     = int64(-1)
		timer   = time.NewTimer(0)
		backoff = interval
		fetched = make(chan []Trade)
		pending []Trade
	)

	fetch := func() {
		start := time.Now()
		if history, next, err := c.History(pair, since); err == nil {
			filtered := []Trade{}
			for _, t := range history {
				if t.Id > tid {
					filtered = append(filtered, t)
					tid = t.Id
				}
			}
			fetched <- filtered

			since = next
			if backoff > interval {
				backoff = backoff / 2
			}
		} else {
			backoff *= 2
			log.Printf("Error getting history: %s", err.Error())
			log.Printf("Waiting for %v before retrying...", backoff)
		}

		dur := time.Now().Sub(start)
		timer.Reset(backoff - dur)
	}

	go func() {
		var first Trade
		var updates chan Trade

		for {
			if len(pending) > 0 {
				updates = trades
				first = pending[0]
			} else {
				updates = nil
			}

			select {
			case <-timer.C:
				go fetch()
			case history := <-fetched:
				pending = append(pending, history...)
			case updates <- first:
				pending = pending[1:]
			case <-closing:
				close(trades)
				timer.Stop()
				break
			}
		}
	}()

	return &Streamer{C: trades, Closing: closing}
}
