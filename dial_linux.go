// +build linux

package main

import (
	"crypto/tls"
	"log"
	"net"
	"syscall"
)

func tlsDialOptWithoutCfg(network, addr string) (net.Conn, error) {
	return tlsDialOptWithCfg(network, addr, tlsConfig)
}

func tlsDialOptWithCfg(network, addr string, cfg *tls.Config) (net.Conn, error) {
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
