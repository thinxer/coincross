package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"code.google.com/p/go-commander"

	s "github.com/thinxer/coincross"
)

func init() {
	cmd := newCmd("balance", "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		b, err := client.Balance()
		check(err)
		for k, v := range b {
			fmt.Printf("%v:%v\n", k, v)
		}
	}
}

func trade(c s.Client, tradeType s.TradeType, args []string) {
	price := must(strconv.ParseFloat(args[0], 64)).(float64)
	amount := must(strconv.ParseFloat(args[1], 64)).(float64)
	id, err := client.Trade(tradeType, flagPair, price, amount)
	check(err)
	fmt.Println(id)
}

func init() {
	cmd := newCmd("buy", "price amount")
	cmd.Run = func(cmd *commander.Command, args []string) {
		trade(client, s.Buy, args)
	}
}

func init() {
	cmd := newCmd("sell", "price amount")
	cmd.Run = func(cmd *commander.Command, args []string) {
		trade(client, s.Sell, args)
	}
}

func init() {
	cmd := newCmd("orders", "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		orders, err := client.Orders()
		check(err)
		for _, o := range orders {
			fmt.Println(o)
		}
	}
}

func init() {
	cmd := newCmd("cancel", "orderid")
	cmd.Run = func(cmd *commander.Command, args []string) {
		orderId := must(strconv.ParseInt(args[0], 10, 64)).(int64)
		ok, err := client.Cancel(orderId)
		check(err)
		fmt.Println(ok)
	}
}

func init() {
	cmd := newCmd("transactions", "[-limit 50]")
	limit := (&cmd.Flag).Int("limit", 50, "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		tr, err := client.Transactions(*limit)
		check(err)
		for _, t := range tr {
			fmt.Println(t)
		}
	}
}

func init() {
	cmd := newCmd("history", "[-since=-1]")
	since := (&cmd.Flag).Int64("since", -1, "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		trades, _, err := client.History(flagPair, *since)
		check(err)
		for _, t := range trades {
			fmt.Println(t)
		}

	}
}

func init() {
	cmd := newCmd("orderbook", "[-limit 50]")
	limit := flag.Int("limit", 50, "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		orders, err := client.Orderbook(flagPair, *limit)
		check(err)
		fmt.Println("Amount\t\tAsks\t\tBids\t\tAmount")
		min := len(orders.Asks)
		if len(orders.Bids) < min {
			min = len(orders.Bids)
		}
		for i := 0; i < min; i++ {
			fmt.Printf("%-16.12g%-16.12g%-16.12g%-16.12g\n", orders.Asks[i].Amount, orders.Asks[i].Price, orders.Bids[i].Price, orders.Bids[i].Amount)
		}
	}
}

func init() {
	cmd := newCmd("watch", "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		streamer := client.Stream(flagPair, -1)
		for t := range streamer.C {
			fmt.Println(t)
		}
	}
}

func init() {
	cmd := newCmd("ticker", "")
	cmd.Run = func(cmd *commander.Command, args []string) {
		ticker, err := client.Ticker(flagPair)
		check(err)
		fmt.Printf("%+v\n", ticker)
	}
}

func check(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "Error: %v [%s:%d]\n", err, file, line)
		os.Exit(2)
	}
}

func must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}
