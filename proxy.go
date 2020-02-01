package main

import (
	"encoding/json"
	"io/ioutil"
)

type Proxy struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`

	Directory string `json:"directory"`

	Cookies string
}

func LoadConfig(file string) *Proxy {
	var p Proxy
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &p)
	if err != nil {
		panic(err)
	}
	return &p
}
