package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	DefaultClient *http.Client
	LoginClient   *http.Client
	tlsConfig     = &tls.Config{InsecureSkipVerify: true}
)

func InitClient() {
	DefaultClient = &http.Client{
		Transport: Transport(),
		Timeout:   time.Second * 60,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	LoginClient = &http.Client{
		Transport: Transport(),
		Timeout:   time.Second * 60,
	}
}

func Transport() http.RoundTripper {
	var proxyUrl *url.URL
	var err error
	if proxy.Proxy != "" {
		proxyUrl, err = url.Parse(DefaultProxy.Proxy)
		if err != nil {
			panic(err)
		}
	}

	if proxy.Http2 {
		return &http2.Transport{
			TLSClientConfig: tlsConfig,
			DialTLS:         tlsDialOptWithCfg,
		}
	} else {
		return &http.Transport{
			TLSClientConfig: tlsConfig,
			DialTLS:         tlsDialOptWithoutCfg,
			Proxy:           http.ProxyURL(proxyUrl),
		}
	}
}

func (p *Proxy) SimpleFetchClient(method, path string, header url.Values, client *http.Client) (*http.Response, error) {
	req, err := http.NewRequest(method, "https://vpns.jlu.edu.cn"+path, strings.NewReader(header.Encode()))
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Cookie", p.Cookies)
	return client.Do(req)
}

func (p *Proxy) SimpleFetch(method, path string, header url.Values) (*http.Response, error) {
	return p.SimpleFetchClient(method, path, header, DefaultClient)
}

func (p *Proxy) SimpleFetchLogin(method, path string, header url.Values) (*http.Response, error) {
	return p.SimpleFetchClient(method, path, header, LoginClient)
}
