package main

import (
	"fmt"
	"net/http"
)

func main() {
	DefaultClient = NewClient()
	p := LoadConfig("config.json")

	err := p.Login()
	if err != nil {
		panic(err)
	}

	fmt.Println("Start server on port " + p.Port)
	err = http.ListenAndServe(":"+p.Port, p)
	if err != nil {
		panic(err)
	}
}
