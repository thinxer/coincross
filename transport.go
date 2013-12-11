package gocoins

import (
	"net"
	"net/http"
	"time"
)

func TimeoutDialer(connectTimeout time.Duration, readWriteTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(readWriteTimeout))
		return conn, nil
	}
}

func TimeoutTransport(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
	}
}
