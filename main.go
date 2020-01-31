package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Yesterday17/jlu-http-proxy/mitm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	hostname, _ = os.Hostname()

	dir      = path.Join(os.Getenv("HOME"), ".jlu-http-proxy")
	keyFile  = path.Join(dir, "ca-key.pem")
	certFile = path.Join(dir, "ca-cert.pem")
)

func main() {
	DefaultClient = NewClient()
	p := LoadConfig("config.json")

	ca, err := loadCA()
	if err != nil {
		log.Fatal(err)
	}

	err = p.Login()
	if err != nil {
		panic(err)
	}

	proxy := &mitm.Proxy{
		Wrap: func(upstream http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				HandleRequest(w, r)
			})
		},
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		//SkipRequest: func(request *http.Request) bool {
		//	return request.URL.Host == "vpns.jlu.edu.cn"
		//},
	}

	fmt.Println("Start server on port " + p.Port)
	err = http.ListenAndServe(":"+p.Port, proxy)
	if err != nil {
		panic(err)
	}
}

func loadCA() (cert tls.Certificate, err error) {
	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if os.IsNotExist(err) {
		cert, err = genCA()
	}
	if err == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	}
	return
}

func genCA() (cert tls.Certificate, err error) {
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}
	certPEM, keyPEM, err := mitm.GenerateCA(hostname)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = ioutil.WriteFile(certFile, certPEM, 0400)
	if err == nil {
		err = ioutil.WriteFile(keyFile, keyPEM, 0400)
	}
	return cert, err
}
