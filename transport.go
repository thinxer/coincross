package coincross

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

func timeoutDialer(connectTimeout time.Duration, totalTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
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

// TimeoutTransport returns a http transport with given connect and total timeout.
func TimeoutTransport(connectTimeout time.Duration, totalTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Dial: timeoutDialer(connectTimeout, totalTimeout),
	}
}

// ProxyTransport can setup proxy for transport
func ProxyTransport(transport *http.Transport, proxy_addr string) *http.Transport {
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(proxy_addr)
	transport.Proxy = http.ProxyURL(url_proxy)
	return transport
}

//SSLTransport can setup SSL for transport
func SSLTransport(transport *http.Transport, flag bool) *http.Transport {
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: flag}
	return transport
}
