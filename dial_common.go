// +build !linux

package main

import (
	"crypto/tls"
	"net"
)

func tlsDialOptWithoutCfg(network, addr string) (net.Conn, error) {
	return tlsDialOptWithCfg(network, addr, tlsConfig)
}

func tlsDialOptWithCfg(network, addr string, cfg *tls.Config) (net.Conn, error) {
	return tls.Dial(network, addr, cfg)
}
