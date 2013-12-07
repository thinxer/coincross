package btcchina

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thinxer/gocoins"
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

// Actually the secret should be bytes. string is just more convenient here.
func MakeClient(apikey, secret string) *BTCChina {
	return &BTCChina{apikey, []byte(secret), &http.Client{}}
}

func (bc *BTCChina) request(method string, params []interface{}, reply interface{}) (err error) {
	tonce := time.Now().UnixNano() / 1000
	data := map[string]interface{}{
		"id":            fmt.Sprintf("%d", tonce),
		"tonce":         tonce,
		"accesskey":     bc.apikey,
		"requestmethod": "post",
		"method":        method,
		"params":        params,
	}

	var message bytes.Buffer
	fields := strings.Split("tonce accesskey requestmethod id method", " ")
	for _, field := range fields {
		message.WriteString(fmt.Sprintf("%s=%v&", field, data[field]))
	}
	message.WriteString(fmt.Sprintf("params=%s", php_implode(params)))
	h := hmac.New(sha1.New, bc.secret)
	h.Write(message.Bytes())
	digest := hex.EncodeToString(h.Sum(nil))

	data_json, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", ENDPOINT, bytes.NewReader(data_json))
	req.SetBasicAuth(bc.apikey, digest)
	req.Header.Set("Json-Rpc-Tonce", fmt.Sprintf("%d", tonce))
	if r, err := bc.client.Do(req); err == nil {
		decoder := json.NewDecoder(r.Body)
		var response struct {
			Result interface{}
			Id     string
		}
		response.Result = reply
		err = decoder.Decode(&response)
		r.Body.Close()
	}
	return
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
		BtcDepositAddress    string `json:"btc_deposit_address"`
		BtcWithdrawalAddress string `json:"btc_withdrawal_address"`
		OtpEnabled           bool   `json:"otp_enabled"`
		TradeFee             string `json:"trade_fee"`
		TradePasswordEnabled bool   `json:"trade_password_enabled"`
		DailyBtcLimit        int    `json:"daily_btc_limit"`
	}
}

func (bc *BTCChina) AccountInfo() (info *AccountInfo, err error) {
	info = new(AccountInfo)
	err = bc.request("getAccountInfo", []interface{}{}, info)
	return
}

func (bc *BTCChina) Balance() (balance *gocoins.Balance, err error) {
	if rai, err := bc.AccountInfo(); err == nil {
		balance = gocoins.MakeBalance()
		balance.Available[gocoins.CNY], _ = strconv.ParseFloat(rai.Balance["cny"].Amount, 64)
		balance.Available[gocoins.BTC], _ = strconv.ParseFloat(rai.Balance["btc"].Amount, 64)
		balance.Frozen[gocoins.CNY], _ = strconv.ParseFloat(rai.Frozen["cny"].Amount, 64)
		balance.Frozen[gocoins.BTC], _ = strconv.ParseFloat(rai.Frozen["btc"].Amount, 64)
	}
	return
}

func (bc *BTCChina) Buy(_ gocoins.Pair, price, amount float64) (success bool, err error) {
	err = bc.request("buyOrder", []interface{}{price, amount}, &success)
	return
}

func (bc *BTCChina) Sell(_ gocoins.Pair, price, amount float64) (success bool, err error) {
	err = bc.request("sellOrder", []interface{}{price, amount}, &success)
	return
}

func (bc *BTCChina) Orderbook(_ gocoins.Pair, limit int32) (orderbook *gocoins.Orderbook, err error) {
	var response struct {
		MarketDepth gocoins.Orderbook `json:"market_depth"`
	}
	err = bc.request("getMarketDepth2", []interface{}{limit}, &response)
	orderbook = &response.MarketDepth
	return
}

func (bc *BTCChina) History(_ gocoins.Pair, since int64) (trades []gocoins.Trade, err error) {
	url := HISTORY
	if since >= 0 {
		url = fmt.Sprintf("%s?since=%d", url, since)
	}

	var ts []struct {
		Tid, Type, Date string
		Amount, Price   float64
	}
	if err = getjson(url, &ts); err != nil {
		return
	}

	trades = make([]gocoins.Trade, len(ts))
	for idx, tx := range ts {
		trades[idx].Id, _ = strconv.ParseInt(tx.Tid, 10, 64)
		trades[idx].Date, _ = strconv.ParseInt(tx.Date, 10, 64)
		trades[idx].Price = tx.Price
		trades[idx].Amount = tx.Amount
		if tx.Type == "buy" {
			trades[idx].Type = gocoins.Buy
		} else if tx.Type == "sell" {
			trades[idx].Type = gocoins.Sell
		} else {
			trades[idx].Type = gocoins.Noop
		}
	}
	return
}

func (bc *BTCChina) Ticker(_ gocoins.Pair) (t *gocoins.Ticker, err error) {
	var v map[string]struct {
		Buy, Sell, Last, Vol, High, Low string
	}
	if err = getjson(TICKER, &v); err != nil {
		return
	}
	ticker := v["ticker"]
	t = new(gocoins.Ticker)
	t.Buy, _ = strconv.ParseFloat(ticker.Buy, 64)
	t.Sell, _ = strconv.ParseFloat(ticker.Sell, 64)
	t.Last, _ = strconv.ParseFloat(ticker.Last, 64)
	t.Volume, _ = strconv.ParseFloat(ticker.Vol, 64)
	t.High, _ = strconv.ParseFloat(ticker.High, 64)
	t.Low, _ = strconv.ParseFloat(ticker.Low, 64)
	return
}

func getjson(url string, v interface{}) (err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	res.Body.Close()
	return
}

func php_float(v interface{}) string {
	s := fmt.Sprintf("%f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
func php_implode(values []interface{}) string {
	parts := make([]string, 0)
	for _, v := range values {
		switch v := v.(type) {
		case bool:
			if v {
				parts = append(parts, "1")
			} else {
				parts = append(parts, "")
			}
		case float32:
			parts = append(parts, php_float(v))
		case float64:
			parts = append(parts, php_float(v))
		case string:
			parts = append(parts, v)
		default:
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(parts, ",")
}
