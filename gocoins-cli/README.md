#gocoins-cli

A simple CLI for gocoins.
It is recommended that you create a wrapper script for convenience.

	#!/bin/sh
	export EXCHANGE=btcchina
	export APIKEY=your_api_key
	export SECRET=your_secret
	./gocoins-cli $@
