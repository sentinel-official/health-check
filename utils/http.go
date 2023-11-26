package utils

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func NewTransport(proxyAddr string) (*http.Transport, error) {
	var (
		tlsClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		transport = &http.Transport{
			TLSClientConfig: tlsClientConfig,
		}
	)

	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			return nil, err
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return nil, err
		}

		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			TLSClientConfig: tlsClientConfig,
		}
	}

	return transport, nil
}
