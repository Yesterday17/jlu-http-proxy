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
	dir, _      = os.Getwd()
	proxy       *Proxy
)

func main() {
	cfgPath := "config.json"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	proxy = LoadConfig(cfgPath)
	InitClient()

	// Directory
	if proxy.Directory != "" {
		dir = proxy.Directory
	}
	ca := loadCA(dir)

	// Login
	if err := proxy.Login(); err != nil {
		panic(err)
	}

	// Proxy
	p := &mitm.Proxy{
		Handle: func(https bool) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				r.URL.Scheme = "http"
				if https {
					r.URL.Scheme += "s"
				}
				proxy.HandleRequest(w, r)
			}
		},
		CA: &ca,
		TLSServerConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
		TLSClientConfig: &tls.Config{
			// VPNS is insecure
			InsecureSkipVerify: true,
		},
	}

	fmt.Println("Start server on port " + proxy.Port)

	// Listen
	if err := http.ListenAndServe(":"+proxy.Port, p); err != nil {
		panic(err)
	}
}

func loadCA(dir string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(path.Join(dir, "ca-cert.pem"), path.Join(dir, "ca-key.pem"))
	if err != nil {
		if os.IsNotExist(err) {
			cert, err = genCA(dir)
			if err != nil {
				Panic(err)
			}
		} else {
			Panic(err)
		}
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		Panic(err)
	}
	return cert
}

func genCA(dir string) (cert tls.Certificate, err error) {
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}
	certPEM, keyPEM, err := mitm.GenerateCA(hostname)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = ioutil.WriteFile(path.Join(dir, "ca-cert.pem"), certPEM, 0400)
	if err == nil {
		err = ioutil.WriteFile(path.Join(dir, "ca-key.pem"), keyPEM, 0400)
	}
	return cert, err
}

func Panic(v ...interface{}) {
	log.Print(v...)
	os.Exit(23)
}
