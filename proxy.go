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
	DefaultProxy = &p
	return DefaultProxy
}
