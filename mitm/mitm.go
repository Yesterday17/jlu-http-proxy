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
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// ServerParam struct
type ServerParam struct {
	CA        *tls.Certificate // the Root CA for generatng on the fly MITM certificates
	TLSConfig *tls.Config      // a template TLS config for the server.
}

// A ServerConn is a net.Conn that holds its clients SNI header in ServerName
// after the handshake.
type ServerConn struct {
	*tls.Conn

	// ServerName is set during Conn's handshake to the client's requested
	// server name set in the SNI header. It is not safe to access across
	// multiple goroutines while Conn is performing the handshake.
	ServerName string
}

// Server wraps cn with a ServerConn configured with p so that during its
// Handshake, it will generate a new certificate using p.CA. After a successful
// Handshake, its ServerName field will be set to the clients requested
// ServerName in the SNI header.
func Server(cn net.Conn, p ServerParam) *ServerConn {
	conf := new(tls.Config)
	if p.TLSConfig != nil {
		*conf = *p.TLSConfig
	}
	sc := new(ServerConn)
	conf.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		sc.ServerName = hello.ServerName
		return getCert(p.CA, hello.ServerName)
	}
	sc.Conn = tls.Server(cn, conf)
	return sc
}

// Proxy is a forward proxy that substitutes its own certificate
// for incoming TLS connections in place of the upstream server's
// certificate.
type Proxy struct {
	// Handle specifies a function for handling the decrypted HTTP request and response.
	Handle func(https bool) func(w http.ResponseWriter, r *http.Request)

	// CA specifies the root CA for generating leaf certs for each incoming
	// TLS request.
	CA *tls.Certificate

	// TLSServerConfig specifies the tls.Config to use when generating leaf
	// cert using CA.
	TLSServerConfig *tls.Config

	// TLSClientConfig specifies the tls.Config to use when establishing
	// an upstream connection for proxying.
	TLSClientConfig *tls.Config

	// Director is function which modifies the request into a new
	// request to be sent using Transport. See the documentation for
	// httputil.ReverseProxy for more details. For mitm proxies, the
	// director defaults to HTTPDirector, but for transparent TLS
	// proxies it should be set to HTTPSDirector.
	Director func(*http.Request)

	// SkipRequest is a function used to skip some requests from being
	// proxied. If it returns true, the request passes by without being
	// wrapped. If false, it's wrapped and proxied. By default we use
	// SkipNone, which doesn't skip any request
	SkipRequest func(*http.Request) bool
}

var okHeader = "HTTP/1.1 200 OK\r\n\r\n"
var establishedHeader = "HTTP/1.1 200 Connection Established\r\n\r\n"

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if p.SkipRequest == nil {
		p.SkipRequest = SkipNone
	}
	if p.Director == nil {
		p.Director = HTTPDirector
	}

	// Skip some requests
	if p.SkipRequest(req) || isWebSocket(req) {
		p.forwardRequest(w, req)
		return
	}

	// http
	if req.Method != "CONNECT" {
		p.Handle(false)(w, req)
		return
	}

	// https
	cn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Println("Hijack:", err)
		http.Error(w, "No Upstream", 503)
		return
	}
	defer cn.Close()

	_, err = io.WriteString(cn, establishedHeader)
	if err != nil {
		log.Println("Write:", err)
		return
	}

	var conn = cn
	var https = false
	if req.URL.Port() == "443" {
		https = true
	} else if req.URL.Port() != "80" {
		if c, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 5}, "tcp", req.URL.Host, nil); err == nil {
			_ = c.Close()
			https = true
		}
	}

	if https {
		serverConn, ok := cn.(*ServerConn)
		if !ok {
			serverConn = Server(cn, ServerParam{
				CA:        p.CA,
				TLSConfig: p.TLSServerConfig,
			})
			if err := serverConn.Handshake(); err != nil {
				log.Println("Server Handshake:", err)
				return
			}
		}
		conn = serverConn
	}

	p.proxyMITM(conn, https)
}

// SkipNone doesn't skip any request and proxy all of them.
func SkipNone(req *http.Request) bool {
	return false
}

func (p *Proxy) proxyMITM(upstream net.Conn, https bool) {
	ch := make(chan int)
	wc := &onCloseConn{upstream, func() { ch <- 1 }}
	_ = http.Serve(&oneShotListener{wc}, http.HandlerFunc(p.Handle(https)))
	<-ch
}

// HTTPDirector is director designed for use in Proxy for http
// proxies.
func HTTPDirector(r *http.Request) {
	r.URL.Host = r.Host
	r.URL.Scheme = "http"
}

// A oneShotListener implements net.Listener whos Accept only returns a
// net.Conn as specified by c followed by an error for each subsequent Accept.
type oneShotListener struct {
	c net.Conn
}

func (l *oneShotListener) Accept() (net.Conn, error) {
	if l.c == nil {
		return nil, errors.New("closed")
	}
	c := l.c
	l.c = nil
	return c, nil
}

func (l *oneShotListener) Close() error {
	return nil
}

func (l *oneShotListener) Addr() net.Addr {
	return l.c.LocalAddr()
}

// A onCloseConn implements net.Conn and calls its f on Close.
type onCloseConn struct {
	net.Conn
	f func()
}

func (c *onCloseConn) Close() error {
	if c.f != nil {
		c.f()
		c.f = nil
	}
	return c.Conn.Close()
}
