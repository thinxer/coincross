package gocoins

import (
	"sync"
)

// A stat interface for streaming trades
type StreamingStat interface {
	// A unique name for this stat, like EMA_0.8
	Name() string
	// Accumulate a trade
	Stat(Trade)
	// Return the lastest statistics
	Last() []float64
}

// A set of analytics for easier tracking
type Analytics struct {
	stats map[string]StreamingStat
	mu    sync.Mutex
}

func MakeAnalytics() *Analytics {
	return &Analytics{stats: make(map[string]StreamingStat)}
}

func (a *Analytics) Add(s StreamingStat) {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.stats[s.Name()]
	if !ok {
		a.stats[s.Name()] = s
	}
}

func (a *Analytics) Get(sstat StreamingStat) []float64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	if stat, ok := a.stats[sstat.Name()]; ok {
		return stat.Last()
	}
	return nil
}

func (a *Analytics) GetAll() (m map[string][]float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	m = make(map[string][]float64)
	for name, stat := range a.stats {
		m[name] = stat.Last()
	}
	return
}

func (a *Analytics) Stat(t Trade) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, stat := range a.stats {
		stat.Stat(t)
	}
}
