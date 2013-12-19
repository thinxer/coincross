package gocoins

import (
	"net"
	"net/http"
	"time"
)

func TimeoutDialer(connectTimeout time.Duration, totalTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(start.Add(totalTimeout))
		return conn, nil
	}
}

func TimeoutTransport(connectTimeout time.Duration, totalTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Dial: TimeoutDialer(connectTimeout, totalTimeout),
	}
}
