package main

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"
)

var (
	DefaultClient = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialTLS:         dialOpt,
		},
		Timeout: time.Second * 60,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	LoginClient = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialTLS:         dialOpt,
		},
		Timeout: time.Second * 60,
	}
)

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

func dialOpt(network, addr string, cfg *tls.Config) (net.Conn, error) {
	if DefaultProxy.Mark == 0 {
		return tls.Dial(network, addr, cfg)
	}

	return tls.DialWithDialer(&net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, DefaultProxy.Mark)
				if err != nil {
					if err.(syscall.Errno) == 1 {
						Panic("Operation not permitted, make sure you're running in root privilege!")
					} else {
						log.Printf("control: %s", err)
					}
					return
				}
			})
		},
	}, network, addr, cfg)
}
