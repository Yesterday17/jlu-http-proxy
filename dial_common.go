// +build !linux

package main

import (
	"crypto/tls"
	"net"
)

func dialOpt(network, addr string, cfg *tls.Config) (net.Conn, error) {
	return tls.Dial(network, addr, cfg)
}
