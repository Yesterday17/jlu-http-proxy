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
	"sync"
	"time"
)

// Certificates are cached locally to avoid unnecessary regeneration
const certCacheMaxSize = 1000

var (
	certCache      = make(map[*tls.Certificate]map[string]*tls.Certificate)
	certCacheMutex sync.RWMutex
)

func getCert(ca *tls.Certificate, host string) (*tls.Certificate, error) {
	if c := getCachedCert(ca, host); c != nil {
		return c, nil
	}
	cert, err := GenerateCert(ca, host)
	if err != nil {
		return nil, err
	}
	cacheCert(ca, host, cert)
	return cert, nil
}

func getCachedCert(ca *tls.Certificate, host string) *tls.Certificate {
	certCacheMutex.RLock()
	defer certCacheMutex.RUnlock()

	if certCache[ca] == nil {
		return nil
	}
	cert := certCache[ca][host]
	if cert == nil || cert.Leaf.NotAfter.Before(time.Now()) {
		return nil
	}
	return cert
}

func cacheCert(ca *tls.Certificate, host string, cert *tls.Certificate) {
	certCacheMutex.Lock()
	defer certCacheMutex.Unlock()

	if certCache[ca] == nil || len(certCache[ca]) > certCacheMaxSize {
		certCache[ca] = make(map[string]*tls.Certificate)
	}
	certCache[ca][host] = cert
}
