// Client implementation for BTC-E.
package btce

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	s "github.com/thinxer/coincross"
)

const (
	PRIVATE_API = "https://btc-e.com/tapi"
	PUBLIC_API  = "https://btc-e.com/api"
)

type BTCE struct {
	key    string
	secret []byte
	client *http.Client
}

func NewClient(apikey, secret string, transport *http.Transport) s.Client {
	return &BTCE{apikey, []byte(secret), &http.Client{
		Transport: transport,
	}}
}

func (b *BTCE) request(method string, params map[string]interface{}, v interface{}) (err error) {
	params["method"] = method
	params["nonce"] = time.Now().Unix()
	form := url.Values{}
	for key, value := range params {
		form.Set(key, fmt.Sprintf("%v", value))
	}
	data := []byte(form.Encode())

	h := hmac.New(sha512.New, b.secret)
	h.Write(data)
	sign := hex.EncodeToString(h.Sum(nil))

	request, _ := http.NewRequest("POST", PRIVATE_API, bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Key", b.key)
	request.Header.Set("Sign", sign)
	response, err := b.client.Do(request)
	if err == nil {
		var body struct {
			Success int
			Return  interface{}
			Error   string
		}
		body.Return = v
		err = decode(response.Body, &body)
		if err == nil {
			if body.Success == 0 {
				return fmt.Errorf("BTC-E Error: %v", body.Error)
			}
		}
	}
	return
}

type Funds map[string]float64
type AccountInfo struct {
	Funds  Funds
	Rights struct {
		Info     int
		Trade    int
		Withdraw int
	}
	TransactionCount int   `json:"transaction_count"`
	OpenOrders       int   `json:"open_orders"`
	ServerTime       int64 `json:"server_time"`
}

func (b *BTCE) AccountInfo() (info *AccountInfo, err error) {
	info = new(AccountInfo)
	err = b.request("getInfo", map[string]interface{}{}, info)
	return
}

func (b *BTCE) Balance() (balance map[s.Symbol]float64, err error) {
	if info, err := b.AccountInfo(); err == nil {
		balance = make(map[s.Symbol]float64)
		for symbol, amount := range info.Funds {
			balance[s.Symbol(strings.ToUpper(symbol))] = amount
		}
	}
	return
}

func (b *BTCE) Trade(tradeType s.TradeType, pair s.Pair, price, amount float64) (orderId int64, err error) {
	var reply struct {
		Received float64
		Remains  float64
		OrderId  int64 `json:"order_id"`
		Funds    Funds
	}
	err = b.request("Trade", map[string]interface{}{
		"pair":   pair.LowerString(),
		"type":   strings.ToLower(tradeType.String()),
		"rate":   price,
		"amount": amount,
	}, &reply)
	if err == nil {
		orderId = reply.OrderId
	}
	return
}

func (b *BTCE) Cancel(orderId int64) (success bool, err error) {
	var reply struct {
		OrderId int64 `json:"order_id"`
		Funds   Funds
	}
	err = b.request("CancelOrder", map[string]interface{}{"order_id": orderId}, &reply)
	success = reply.OrderId == orderId
	return
}

// Transactions returns your transactions,
// including trades, deposits, withdraws, placed/cancelled orders etc.
func (b *BTCE) Transactions(limit int) (transactions []s.Transaction, err error) {
	var reply map[string]struct {
		Type      int
		Amount    float64
		Currency  string
		Desc      string
		Status    int
		Timestamp int64
	}
	if err = b.request("TransHistory", map[string]interface{}{"count": limit, "order": "DESC"}, &reply); err != nil {
		return
	}
	for id, tr := range reply {
		var t s.Transaction
		t.Id, _ = strconv.ParseInt(id, 10, 64)
		t.Timestamp = tr.Timestamp
		// TODO parse DESC and fill amounts better
		t.Amounts = map[s.Symbol]float64{
			s.Symbol(strings.ToUpper(tr.Currency)): tr.Amount,
		}
		t.Descritpion = tr.Desc
		transactions = append(transactions, t)
	}
	return
}

// Orders will return your active orders for all pairs.
func (b *BTCE) Orders() (orders []s.Order, err error) {
	var reply map[string]struct {
		Pair             s.Pair
		Type             s.TradeType
		Amount           float64
		Rate             float64
		TimestampCreated int64 `json:"timestamp_created"`
		Status           int
	}
	err = b.request("ActiveOrders", map[string]interface{}{}, &reply)
	for id, order := range reply {
		var o s.Order
		o.Id, _ = strconv.ParseInt(id, 10, 64)
		o.Pair = order.Pair
		o.Type = order.Type
		o.Price = order.Rate
		o.Remain = order.Amount
		o.Amount = order.Amount
		o.Timestamp = order.TimestampCreated
		orders = append(orders, o)
	}
	return
}

// TradeHistory returns your past trade transactions.
func (b *BTCE) TradeHistory(pair s.Pair, since int64) (trades []s.Trade, err error) {
	var reply map[string]struct {
		Pair        s.Pair
		Type        s.TradeType
		Amount      float64
		Rate        float64
		OrderId     int64 `json:"order_id"`
		IsYourOrder int   `json:"is_your_order"`
		Timestamp   int64
	}
	params := map[string]interface{}{"from_id": since, "order": "DESC"}
	if pair != s.ALL {
		params["pair"] = pair.LowerString()
	}
	if err = b.request("TradeHistory", params, &reply); err != nil {
		return
	}
	for id, trade := range reply {
		var t s.Trade
		t.Id, _ = strconv.ParseInt(id, 10, 64)
		t.Price = trade.Rate
		t.Amount = trade.Amount
		t.Timestamp = trade.Timestamp
		t.Type = trade.Type
		t.Pair = trade.Pair
		trades = append(trades, t)
	}
	return
}

func (b *BTCE) Orderbook(pair s.Pair, limit int) (orderbook *s.Orderbook, err error) {
	url := fmt.Sprintf("%s/3/depth/%s", PUBLIC_API, pair.LowerString())
	var reply map[string]struct {
		Asks, Bids [][]float64
	}
	transform := func(trades [][]float64) (r []struct{ Price, Amount float64 }) {
		for _, p := range trades {
			r = append(r, struct{ Price, Amount float64 }{p[0], p[1]})
		}
		return
	}
	if err = getjson(b.client, url, &reply); err != nil {
		return
	}
	orderbook = new(s.Orderbook)
	reply_orderbook, _ := reply[pair.LowerString()]
	orderbook.Asks = transform(reply_orderbook.Asks)
	orderbook.Bids = transform(reply_orderbook.Bids)
	return
}

// Note that BTC-E use `Timestamp` field for the `since` parameter
func (b *BTCE) History(pair s.Pair, since int64) (trades []s.Trade, next int64, err error) {
	next = since
	url := fmt.Sprintf("%s/3/trades/%s", PUBLIC_API, pair.LowerString())
	if since > 0 {
		url = fmt.Sprintf("%s?since=%d", url, since)
	}
	var reply map[string][]struct {
		Tid       int64
		Price     float64
		Amount    float64
		Type      s.TradeType
		Timestamp int64
	}
	if err = getjson(b.client, url, &reply); err != nil {
		return
	}
	reply_trades := reply[pair.LowerString()]
	for i := len(reply_trades) - 1; i >= 0; i-- {
		trade := reply_trades[i]
		var t s.Trade
		t.Id = trade.Tid
		t.Timestamp = trade.Timestamp
		t.Price = trade.Price
		t.Amount = trade.Amount
		t.Type = trade.Type
		t.Pair = pair
		trades = append(trades, t)
		next = t.Timestamp
	}
	return
}

func (b *BTCE) Ticker(pair s.Pair) (t *s.Ticker, err error) {
	url := fmt.Sprintf("%s/3/ticker/%s", PUBLIC_API, pair.LowerString())
	var reply map[string]struct {
		High, Low, Avg, Vol, Last, Buy, Sell float64
		Vol_Cur                              float64 `json:"vol_cur"`
		Updated                              int64
	}
	if err = getjson(b.client, url, &reply); err != nil {
		return
	}
	tt := reply[pair.LowerString()]
	t = &s.Ticker{tt.Buy, tt.Sell, tt.High, tt.Low, tt.Last, tt.Vol_Cur}
	return
}

func (b *BTCE) Stream(pair s.Pair, since int64, out chan s.Trade) error {
	s.Tail(b, pair, since, 2*time.Second, out)
	return nil
}

type Info struct {
	ServerTime int64 `json:"server_time"`
	Pairs      map[string]struct {
		DecimalPlaces int     `json:"decimal_places"`
		MinPrice      float64 `json:"min_price"`
		MaxPrice      float64 `json:"max_price"`
		MinAmount     float64 `json:"min_amount"`
		Hidden        int
		Fee           float64
	}
}

func (b *BTCE) Info() (info *Info, err error) {
	url := fmt.Sprintf("%s/3/info", PUBLIC_API)
	info = new(Info)
	err = getjson(b.client, url, info)
	return
}

func getjson(client *http.Client, url string, v interface{}) (err error) {
	res, err := client.Get(url)
	if err == nil {
		return decode(res.Body, v)
	}
	return
}

func decode(body io.ReadCloser, v interface{}) error {
	content, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return fmt.Errorf("Unmarshal failed: %s", string(content))
	}
	return nil
}

func init() {
	s.Register("btce", NewClient)
}
