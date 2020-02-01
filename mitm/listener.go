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
	"net"
)

type listener struct {
	net.Listener
	ca   *tls.Certificate
	conf *tls.Config
}

// NewListener returns a net.Listener that generates a new cert from ca for
// each new Accept. It uses SNI to generate the cert, and therefore only
// works with clients that send SNI headers.
//
// This is useful for building transparent MITM proxies.
func NewListener(inner net.Listener, ca *tls.Certificate, conf *tls.Config) net.Listener {
	return &listener{inner, ca, conf}
}

func (l *listener) Accept() (net.Conn, error) {
	cn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	sc := Server(cn, ServerParam{
		CA:        l.ca,
		TLSConfig: l.conf,
	})
	return sc, nil
}
