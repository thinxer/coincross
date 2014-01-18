// Client implementation for BTCChina.
package btcchina

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	s "github.com/thinxer/coincross"
)

const (
	ENDPOINT  = "https://api.btcchina.com/api_trade_v1.php"
	HISTORY   = "https://data.btcchina.com/data/historydata"
	ORDERBOOK = "https://data.btcchina.com/data/orderbook"
	TICKER    = "https://data.btcchina.com/data/ticker"
)

type BTCChina struct {
	apikey string
	secret []byte
	client *http.Client
}

func New(apikey, secret string, transport *http.Transport) *BTCChina {
	return &BTCChina{apikey, []byte(secret), &http.Client{
		Transport: transport,
	}}
}

type AccountInfo struct {
	Balance, Frozen map[string]struct {
		Amount           string
		AmountInteger    string `json:"amount_integer"`
		Currency, Symbol string
		AmountDecimal    int `json:"amount_decimal"`
	}
	Profile struct {
		Username             string
		BtcDepositAddress    string  `json:"btc_deposit_address"`
		BtcWithdrawalAddress string  `json:"btc_withdrawal_address"`
		OtpEnabled           bool    `json:"otp_enabled"`
		TradeFee             float64 `json:"trade_fee"`
		TradePasswordEnabled bool    `json:"trade_password_enabled"`
		DailyBtcLimit        int     `json:"daily_btc_limit"`
	}
}

func (bc *BTCChina) AccountInfo() (info *AccountInfo, err error) {
	info = new(AccountInfo)
	err = bc.request("getAccountInfo", []interface{}{}, info)
	return
}

func (bc *BTCChina) Balance() (balance map[s.Symbol]float64, err error) {
	rai, err := bc.AccountInfo()
	if err == nil {
		balance = make(map[s.Symbol]float64)
		balance[s.CNY], _ = strconv.ParseFloat(rai.Balance["cny"].Amount, 64)
		balance[s.BTC], _ = strconv.ParseFloat(rai.Balance["btc"].Amount, 64)
	}
	return
}

func (bc *BTCChina) Trade(tradeType s.TradeType, _ s.Pair, price, amount float64) (orderId int64, err error) {
	var success bool
	switch tradeType {
	case s.Sell:
		err = bc.request("sellOrder", []interface{}{price, amount}, &success)
	case s.Buy:
		err = bc.request("buyOrder", []interface{}{price, amount}, &success)
	}
	if err == nil && !success {
		err = s.TradeError(fmt.Errorf("place order failed"))
	}
	// TODO
	orderId = -1
	return
}

func (bc *BTCChina) Cancel(orderId int64) (success bool, err error) {
	err = bc.request("cancelOrder", []interface{}{orderId}, &success)
	return
}

func (bc *BTCChina) Transactions(limit int) (transactions []s.Transaction, err error) {
	var response struct {
		Transaction []struct {
			Id        int64
			Type      string
			BtcAmount string `json:"btc_amount"`
			CnyAmount string `json:"cny_amount"`
			Date      int64
		}
	}
	if err = bc.request("getTransactions", []interface{}{"all", limit}, &response); err == nil {
		for _, tr := range response.Transaction {
			var t s.Transaction
			t.Id = tr.Id
			t.Timestamp = tr.Date
			t.Amounts = make(map[s.Symbol]float64)
			t.Amounts[s.BTC], _ = strconv.ParseFloat(tr.BtcAmount, 64)
			t.Amounts[s.CNY], _ = strconv.ParseFloat(tr.CnyAmount, 64)
			t.Descritpion = tr.Type
			transactions = append(transactions, t)
		}
	}
	return
}

func (bc *BTCChina) Orders() (orders []s.Order, err error) {
	var response struct {
		Order []struct {
			Id             int64
			Type           s.TradeType
			Price          string
			Currency       string
			Amount         string
			AmountOriginal string `json:"amount_original"`
			Date           int64
			Status         string
		}
	}
	if err = bc.request("getOrders", []interface{}{}, &response); err == nil {
		for _, order := range response.Order {
			var o s.Order
			o.Id = order.Id
			o.Type = order.Type
			o.Price, _ = strconv.ParseFloat(order.Price, 64)
			o.Amount, _ = strconv.ParseFloat(order.AmountOriginal, 64)
			o.Remain, _ = strconv.ParseFloat(order.Amount, 64)
			o.Pair = s.BTC_CNY
			o.Timestamp = order.Date
			orders = append(orders, o)
		}
	}
	return
}

func (bc *BTCChina) Orderbook(_ s.Pair, limit int) (orderbook *s.Orderbook, err error) {
	var response struct {
		MarketDepth struct {
			Ask, Bid []struct {
				Price, Amount float64
			}
		} `json:"market_depth"`
	}
	err = bc.request("getMarketDepth2", []interface{}{limit}, &response)
	orderbook = &s.Orderbook{response.MarketDepth.Ask, response.MarketDepth.Bid}
	return
}

func (bc *BTCChina) History(_ s.Pair, since int64) (trades []s.Trade, next int64, err error) {
	next = since
	url := HISTORY
	if since >= 0 {
		url = fmt.Sprintf("%s?since=%d", url, since)
	}

	var ts []struct {
		Tid, Date     string
		Type          s.TradeType
		Amount, Price float64
	}
	if err = getjson(bc.client, url, &ts); err != nil {
		return
	}

	var t s.Trade
	for _, tx := range ts {
		t.Id, _ = strconv.ParseInt(tx.Tid, 10, 64)
		t.Timestamp, _ = strconv.ParseInt(tx.Date, 10, 64)
		t.Price = tx.Price
		t.Amount = tx.Amount
		t.Type = tx.Type
		t.Pair = s.BTC_CNY
		trades = append(trades, t)
		next = t.Id
	}
	return
}

func (bc *BTCChina) Ticker(_ s.Pair) (t *s.Ticker, err error) {
	var v map[string]struct {
		Buy, Sell, Last, Vol, High, Low string
	}
	if err = getjson(bc.client, TICKER, &v); err != nil {
		return
	}
	ticker := v["ticker"]
	t = new(s.Ticker)
	t.Buy, _ = strconv.ParseFloat(ticker.Buy, 64)
	t.Sell, _ = strconv.ParseFloat(ticker.Sell, 64)
	t.Last, _ = strconv.ParseFloat(ticker.Last, 64)
	t.Volume, _ = strconv.ParseFloat(ticker.Vol, 64)
	t.High, _ = strconv.ParseFloat(ticker.High, 64)
	t.Low, _ = strconv.ParseFloat(ticker.Low, 64)
	return
}

func (bc *BTCChina) Stream(pair s.Pair, since int64) *s.Streamer {
	return s.Tail(bc, pair, since, time.Second)
}

func init() {
	s.Register("btcchina", func(apikey, secret string, transport *http.Transport) s.Client {
		return New(apikey, secret, transport)
	})
}
