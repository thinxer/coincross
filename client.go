package gocoins

type Symbol string
type Pair struct {
	Base, Target Symbol
}

const (
	BTC Symbol = "BTC"
	LTC        = "LTC"
	CNY        = "CNY"
	USD        = "USD"
)

var (
	BTC_CNY = Pair{CNY, BTC}
	LTC_CNY = Pair{CNY, LTC}
	LTC_BTC = Pair{BTC, LTC}
)

type Balance struct {
	Available map[Symbol]float64
	Frozen    map[Symbol]float64
}

func MakeBalance() *Balance {
	b := Balance{}
	b.Available = make(map[Symbol]float64)
	b.Frozen = make(map[Symbol]float64)
	return &b
}

type Ticker struct {
	Buy, Sell, High, Low, Last, Volume float64
}

type TradeType int

const (
	Noop TradeType = iota
	Buy
	Sell
)

type Trade struct {
	Id     int64
	Date   int64
	Type   TradeType
	Price  float64
	Amount float64
}

type Orderbook struct {
	Ask, Bid []struct {
		Price, Amount float64
	}
}

type Client interface {
	Balance() (*Balance, error)
	Buy(pair Pair, price, amount float64) (bool, error)
	Sell(pair Pair, price, amount float64) (bool, error)
	Orderbook(pair Pair, limit int32) (*Orderbook, error)
	History(pair Pair, since int64) ([]Trade, error)
	Ticker(pair Pair) (*Ticker, error)
}
