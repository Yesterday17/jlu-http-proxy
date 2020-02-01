package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var DefaultClient = NewClient()

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 60,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !strings.HasPrefix(req.URL.Path, "/login") &&
				!strings.HasPrefix(req.URL.Path, "/do-login") &&
				!strings.HasPrefix(req.URL.Path, "/do-confirm-login") {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

func (p *Proxy) SimpleFetch(method, path string, header url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, "https://vpns.jlu.edu.cn"+path, strings.NewReader(header.Encode()))
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Cookie", p.Cookies)
	return DefaultClient.Do(req)
}
