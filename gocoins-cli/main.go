// A simple tool for trading on the command line interface
package main

import (
	"flag"
	"os"
	"time"

	s "github.com/thinxer/gocoins"
	_ "github.com/thinxer/gocoins/all"
)

var (
	flagPair    = &s.Pair{s.CNY, s.BTC}
	flagTimeout = flag.Duration("timeout", 10*time.Second, "timeout for connections")

	cmds = make(map[string]func(c s.Client))
)

func init() {
	flag.Var(flagPair, "pair", "pair to operate on")
}

func main() {
	flag.Parse()
	exchange := os.Getenv("EXCHANGE")
	apikey := os.Getenv("APIKEY")
	secret := os.Getenv("SECRET")
	client := s.New(exchange, apikey, secret, s.TimeoutTransport(*flagTimeout, *flagTimeout))

	cmd := flag.Arg(0)
	cmds[cmd](client)
}
