// Code edited from https://github.com/rainforestapp/mitm
/**
Copyright (c) 2015 Keith Rarick

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package mitm

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func isWebSocket(req *http.Request) bool {
	return strings.Contains(strings.ToLower(req.Header.Get("Upgrade")), "websocket") &&
		strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade")
}

// fixWebsocketHeaders "un-canonicalizes" websocket headers for a
// request. According to https://tools.ietf.org/html/rfc6455 the correct form is
// Sec-WebSocket-*, which header canonicalization breaks (some servers care).
func fixWebsocketHeaders(req *http.Request) {
	for header, _ := range req.Header {
		if strings.Contains(header, "Sec-Websocket") {
			val := req.Header.Get(header)
			correctHeader := strings.Replace(header, "Sec-Websocket", "Sec-WebSocket", 1)
			req.Header[correctHeader] = []string{val}
			delete(req.Header, header)
		}
	}
}

// forwardRequest forwards a WebSocket connection directly to the
// source, skipping the request wrapper. Code shamelessly stolen from
// https://groups.google.com/forum/#!topic/golang-nuts/KBx9pDlvFOc
func (p *Proxy) forwardRequest(w http.ResponseWriter, req *http.Request) {
	p.Director(req)
	if isWebSocket(req) {
		fixWebsocketHeaders(req)
	}
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		// Assume there is no port and use default
		host = req.URL.Host
		if req.URL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	address := net.JoinHostPort(host, port)
	var d io.ReadWriteCloser
	if req.URL.Scheme == "https" {
		d, err = tls.Dial("tcp", address, p.TLSClientConfig)
	} else {
		d, err = net.Dial("tcp", address)
	}
	if err != nil {
		log.Printf("forwardRequest: error dialing websocket backend %s: %v", address, err)
		http.Error(w, "No Upstream", 503)
		return
	}
	defer d.Close()

	nc, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Printf("forwardRequest: hijack error: %v", err)
		http.Error(w, "No Upstream", 503)
		return
	}
	defer nc.Close()

	err = req.Write(d)
	if err != nil {
		log.Printf("forwardRequest: error copying request to target: %v", err)
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go cp(d, nc)
	go cp(nc, d)
	<-errc
}
