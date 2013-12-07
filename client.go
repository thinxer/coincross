package gocoins

type Symbol string
type Pair struct {
	Base, Target Symbol
}

func (p Pair) String() string {
	return string(p.Target + "/" + p.Base)
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

type Ticker struct {
	Buy, Sell, High, Low, Last, Volume float64
}

type TradeType int

const (
	_ TradeType = iota
	Ask
	Bid
)

func (t TradeType) String() string {
	switch t {
	case Ask:
		return "Ask"
	case Bid:
		return "Bid"
	default:
		return ""
	}
}

type Trade struct {
	Id        int64
	Timestamp int64
	Type      TradeType
	Price     float64
	Amount    float64
}

type Orderbook struct {
	Ask, Bid []struct {
		Price, Amount float64
	}
}

type OrderStatus int

const (
	_ OrderStatus = iota
	Open
	Closed
	Cancelled
)

func (s OrderStatus) String() string {
	switch s {
	case Open:
		return "Open"
	case Closed:
		return "Closed"
	case Cancelled:
		return "Cancelled"
	default:
		return ""
	}
}

type Order struct {
	Id             int64
	Type           TradeType
	Price          float64
	Remain, Amount float64
	Pair           Pair
	Timestamp      int64
	Status         OrderStatus
}

type TransactionType int

const (
	_ TransactionType = iota
	Deposit
	Withdrawal
	Bought
	Sold
	TradeFee
)

type Transaction struct {
	Id        int64
	Type      TransactionType
	Amounts   map[Symbol]float64
	Timestamp int64
}

type Client interface {
	Balance() (*Balance, error)
	Trade(tradeType TradeType, pair Pair, price, amount float64) (bool, error)
	Cancel(orderId int64) (bool, error)
	Transactions(limit int) ([]Transaction, error)
	Orders() ([]Order, error)
	Orderbook(pair Pair, limit int) (*Orderbook, error)
	History(pair Pair, since int64) ([]Trade, error)
	Ticker(pair Pair) (*Ticker, error)
}
