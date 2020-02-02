package main

import (
	"io/ioutil"
	"regexp"
	"strings"
)

func (p *Proxy) Login() error {
	// Get Cookies first
	resp, err := LoginClient.Get("https://vpns.jlu.edu.cn/login")
	if err != nil {
		return err
	}
	p.Cookies = resp.Header.Get("Set-Cookie")
	p.Cookies = strings.Split(p.Cookies, ";")[0]

	// Auth
	resp, err = p.SimpleFetchLogin("POST", "/do-login?fromUrl=", map[string][]string{
		"auth_type":       {"local"},
		"username":        {p.Username},
		"password":        {p.Password},
		"sms_code":        {""},
		"remember_cookie": {"on"},
	})
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile("logoutOtherToken = '([0-9a-f]+)'")
	if r.Match(body) {
		// Split current active token from html
		p.Cookies = "wengine_vpn_ticket_ecit=" + string(r.FindSubmatch(body)[1])
	}

	// TODO: Check whether logon successfully
	return nil
}
