package main

import (
	"flag"
	"fmt"
	"strconv"

	s "github.com/thinxer/gocoins"
)

func mustInt64(v int64, err error) int64 {
	if err != nil {
		panic(err)
	}
	return v
}

func mustFloat64(v float64, err error) float64 {
	if err != nil {
		panic(err)
	}
	return v
}

func balance(c s.Client) {
	b, err := c.Balance()
	if err == nil {
		for k, v := range b {
			fmt.Printf("%v:%v\n", k, v)
		}
	} else {
		panic(err)
	}
}

func trade(c s.Client, tradeType s.TradeType) {
	price := mustFloat64(strconv.ParseFloat(flag.Arg(1), 64))
	amount := mustFloat64(strconv.ParseFloat(flag.Arg(2), 64))
	id, err := c.Trade(tradeType, *flagPair, price, amount)
	if err == nil {
		fmt.Println(id)
	} else {
		panic(err)
	}
}

func buy(c s.Client) {
	trade(c, s.Buy)
}

func sell(c s.Client) {
	trade(c, s.Sell)
}

func orders(c s.Client) {
	orders, err := c.Orders()
	if err == nil {
		for _, o := range orders {
			fmt.Println(o)
		}
	} else {
		panic(err)
	}

}

func cancel(c s.Client) {
	orderId := mustInt64(strconv.ParseInt(flag.Arg(1), 10, 64))
	ok, err := c.Cancel(orderId)
	if err == nil {
		fmt.Println(ok)
	} else {
		panic(err)
	}
}

func transactions(c s.Client) {
	var limit int64 = 50
	if flag.NArg() > 1 {
		limit = mustInt64(strconv.ParseInt(flag.Arg(1), 10, 64))
	}
	tr, err := c.Transactions(int(limit))
	if err == nil {
		for _, t := range tr {
			fmt.Println(t)
		}
	} else {
		panic(err)
	}
}

func history(c s.Client) {
	trades, _, err := c.History(*flagPair, -1)
	if err == nil {
		for _, t := range trades {
			fmt.Println(t)
		}
	} else {
		panic(err)
	}
}

func orderbook(c s.Client) {
	limit := 50
	if flag.NArg() > 1 {
		limit, _ = strconv.Atoi(flag.Arg(1))
	}
	orders, err := c.Orderbook(*flagPair, limit)
	if err == nil {
		fmt.Println("Asks:")
		for _, o := range orders.Asks {
			fmt.Printf("%v\t%v\n", o.Price, o.Amount)
		}
		fmt.Println("Bids:")
		for _, o := range orders.Bids {
			fmt.Printf("%v\t%v\n", o.Price, o.Amount)
		}
	} else {
		panic(err)
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
	if err == nil {
		fmt.Printf("%+v\n", ticker)
	} else {
		panic(err)
	}
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
