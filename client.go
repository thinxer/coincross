/*
Coincross is a collection of cryptocurrency trading APIs.

This is the base package that defines the unified interface, as well as some
common data structures. API implementations should live in their own package.

However, as the APIs vary a lot from exchange to exchange, it is difficult
to include all the functionalities of the APIs. This package should aim to
provide the minimalist interface for common exchanges, while fixing some broken
APIs in the hidden layer. For example, some API does not return the order id
when you place an order (i.e. BTCChina), implementations should fix this.
Well, to their best.

As more and more exchanges are being supported, this interface may subject to
changes.
*/
package coincross

// Symbol represents a currency, such as USD, CNY, BTC etc.
type Symbol string

// Pair is for a trading pair, such as {USD, BTC}
type Pair struct {
	Base, Target Symbol
}

// Some predefined Symbols. You are not limited to use these.
const (
	BTC Symbol = "BTC"
	LTC        = "LTC"
	CNY        = "CNY"
	USD        = "USD"
)

// Some predefined pairs.
var (
	ALL     = Pair{"", ""}
	BTC_CNY = Pair{CNY, BTC}
	BTC_USD = Pair{USD, BTC}
	LTC_CNY = Pair{CNY, LTC}
	LTC_USD = Pair{USD, LTC}
	LTC_BTC = Pair{BTC, LTC}
)

// TradeType is the direction of a trade.
type TradeType int

const (
	_ TradeType = iota
	Sell
	Buy
)

// Ticker represents for the result of Ticker APIs.
type Ticker struct {
	Buy, Sell, High, Low, Last, Volume float64
}

// A historical trade instance.
// As you can see Trade is a special case of Transaction.
// It's just special enough to make it a unique type.
type Trade struct {
	Id        int64
	Timestamp int64
	Type      TradeType
	Price     float64
	Amount    float64
	Pair      Pair
}

// An active order.
type Order struct {
	Id             int64
	Timestamp      int64
	Type           TradeType
	Price          float64
	Remain, Amount float64
	Pair           Pair
}

// A transaction is an operation to your account's balance.
// All the historical transactions should add up to your current balance.
type Transaction struct {
	Id          int64
	Timestamp   int64
	Amounts     map[Symbol]float64
	Descritpion string
}

// Well, the order book, or the market depth.
type Orderbook struct {
	Asks, Bids []struct {
		Price, Amount float64
	}
}

type Streamer struct {
	C       <-chan Trade
	Closing chan<- bool
}

// This is the interface that every API implementation should use.
type Client interface {
	// Should return the balance of current account.
	Balance() (map[Symbol]float64, error)
	// Use with caution: this method is for real trading.
	Trade(tradeType TradeType, pair Pair, price, amount float64) (int64, error)
	// Cancel an active order.
	Cancel(orderId int64) (bool, error)
	// Returns active orders.
	Orders() ([]Order, error)
	// Returns the transaction history of current account.
	Transactions(limit int) ([]Transaction, error)

	// Returns the orderbook (or market depth) of the given pair.
	// Usually this is a public API.
	Orderbook(pair Pair, limit int) (*Orderbook, error)
	// Returns the trade history of the given pair, as well as a cursor for next since.
	// Usually this is a public API.
	History(pair Pair, since int64) ([]Trade, int64, error)
	// Returns the ticker of the given pair.
	// Usually this is a public API.
	Ticker(pair Pair) (*Ticker, error)

	// History streaming.
	Stream(pair Pair, since int64) *Streamer
}
