// +build !linux

package main

import "crypto/tls"

func dialOpt(network, addr string, cfg *tls.Config) (net.Conn, error) {
	return tls.Dial(network, addr, cfg)
}
