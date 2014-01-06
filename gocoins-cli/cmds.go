package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	s "github.com/thinxer/gocoins"
)

func must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}

func balance(c s.Client) {
	b, err := c.Balance()
	check(err)
	for k, v := range b {
		fmt.Printf("%v:%v\n", k, v)
	}
}

func trade(c s.Client, tradeType s.TradeType) {
	price := must(strconv.ParseFloat(flag.Arg(1), 64)).(float64)
	amount := must(strconv.ParseFloat(flag.Arg(2), 64)).(float64)
	id, err := c.Trade(tradeType, *flagPair, price, amount)
	check(err)
	fmt.Println(id)
}

func buy(c s.Client) {
	trade(c, s.Buy)
}

func sell(c s.Client) {
	trade(c, s.Sell)
}

func orders(c s.Client) {
	orders, err := c.Orders()
	check(err)
	for _, o := range orders {
		fmt.Println(o)
	}

}

func cancel(c s.Client) {
	orderId := must(strconv.ParseInt(flag.Arg(1), 10, 64)).(int64)
	ok, err := c.Cancel(orderId)
	check(err)
	fmt.Println(ok)
}

func transactions(c s.Client) {
	var limit int64 = 50
	if flag.NArg() > 1 {
		limit = must(strconv.ParseInt(flag.Arg(1), 10, 64)).(int64)
	}
	tr, err := c.Transactions(int(limit))
	check(err)
	for _, t := range tr {
		fmt.Println(t)
	}
}

func history(c s.Client) {
	trades, _, err := c.History(*flagPair, -1)
	check(err)
	for _, t := range trades {
		fmt.Println(t)
	}
}

func orderbook(c s.Client) {
	limit := 50
	if flag.NArg() > 1 {
		limit, _ = strconv.Atoi(flag.Arg(1))
	}
	orders, err := c.Orderbook(*flagPair, limit)
	check(err)
	fmt.Println("Asks:")
	for _, o := range orders.Asks {
		fmt.Printf("%v\t%v\n", o.Price, o.Amount)
	}
	fmt.Println("Bids:")
	for _, o := range orders.Bids {
		fmt.Printf("%v\t%v\n", o.Price, o.Amount)
	}
}

func watch(c s.Client) {
	ct := make(chan s.Trade)
	go c.Stream(*flagPair, -1, ct)
	for t := range ct {
		fmt.Println(t)
	}
}

func ticker(c s.Client) {
	ticker, err := c.Ticker(*flagPair)
	check(err)
	fmt.Printf("%+v\n", ticker)
}

func init() {
	cmds["balance"] = balance
	cmds["buy"] = buy
	cmds["sell"] = sell
	cmds["orders"] = orders
	cmds["cancel"] = cancel
	cmds["transactions"] = transactions

	cmds["history"] = history
	cmds["orderbook"] = orderbook
	cmds["watch"] = watch
	cmds["ticker"] = ticker
}

func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "Error: %v [%s:%d]\n", err, file, line)
		os.Exit(2)
	}
}
