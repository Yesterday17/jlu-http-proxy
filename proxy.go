package main

import (
	"encoding/json"
	"io/ioutil"
)

var DefaultProxy *Proxy

type Proxy struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`

	Directory string `json:"directory"`
	Mark      int    `json:"mark"`
	Http2     bool   `json:"http2"`

	Cookies string
}

func LoadConfig(file string) *Proxy {
	var p Proxy
	content, err := ioutil.ReadFile(file)
	if err != nil {
		Panic(err)
	}
	err = json.Unmarshal(content, &p)
	if err != nil {
		Panic(err)
	}

	if p.Username == "" || p.Password == "" || p.Directory == "" || p.Port == "" {
		Panic("Empty entries found in config file!")
	}

	DefaultProxy = &p
	return DefaultProxy
}
