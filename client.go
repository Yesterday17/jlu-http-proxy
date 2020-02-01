package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	"strings"
	"time"
)

var DefaultClient *http.Client

func NewClient() *http.Client {
	DefaultClient = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !strings.HasPrefix(req.URL.Path, "/login") &&
				!strings.HasPrefix(req.URL.Path, "/do-login") &&
				!strings.HasPrefix(req.URL.Path, "/do-confirm-login") {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	return DefaultClient
}
