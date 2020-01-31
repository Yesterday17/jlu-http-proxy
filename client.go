package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var DefaultClient *Client

type Client struct {
	httpClient  *http.Client
	http2Client *http.Client
}

func NewClient() *Client {
	var c Client
	jar, _ := cookiejar.New(nil)
	c.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar:     jar,
		Timeout: time.Second * 60,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	c.http2Client = &http.Client{
		Transport: &http2.Transport{
			AllowHTTP:       true,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar:     jar,
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
	return &c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme == "http" {
		return c.httpClient.Do(req)
	} else {
		return c.http2Client.Do(req)
	}
}
