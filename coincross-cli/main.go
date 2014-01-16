// A simple tool for trading on CLI.
package main

import (
	"flag"
	"os"
	"time"

	commander "code.google.com/p/go-commander"
	s "github.com/thinxer/coincross"
	_ "github.com/thinxer/coincross/all"
)

// Populate this slice with newCmd, and it will be used in the commander.
var cmds []*commander.Command

func newCmd(name, short string) *commander.Command {
	cmd := &commander.Command{
		UsageLine: name + " " + short,
		Flag:      *flag.NewFlagSet(name, flag.ExitOnError),
		Short:     short,
	}
	cmds = append(cmds, cmd)
	return cmd
}

var (
	flagPair    = s.Pair{s.CNY, s.BTC}
	flagTimeout time.Duration
	client      s.Client
)

func main() {
	// Construct the commander
	cmd := commander.Commander{
		Name:     os.Args[0],
		Commands: cmds,
		Flag:     flag.NewFlagSet("cli", flag.ExitOnError),
	}
	cmd.Flag.Var(&flagPair, "pair", "pair to operate on")
	cmd.Flag.DurationVar(&flagTimeout, "timeout", 10*time.Second, "timeout for connections")
	if err := cmd.Flag.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	// Construct the client
	exchange := os.Getenv("EXCHANGE")
	apikey := os.Getenv("APIKEY")
	secret := os.Getenv("SECRET")
	client = s.New(exchange, apikey, secret, s.TimeoutTransport(flagTimeout, flagTimeout))
	if client == nil {
		panic("create client failed")
	}

	// Actually run the commands
	if err := cmd.Run(cmd.Flag.Args()); err != nil {
		return
	}
}
