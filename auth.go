package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (p *Proxy) Login() error {
	// Get Cookies first
	resp, err := DefaultClient.Get("https://vpns.jlu.edu.cn/login")
	if err != nil {
		return err
	}
	p.Cookies = resp.Header.Get("Set-Cookie")
	p.Cookies = strings.Split(p.Cookies, ";")[0]
	log.Println(p.Cookies)

	// Auth
	values := url.Values{}
	values.Set("auth_type", "local")
	values.Set("username", p.Username)
	values.Set("password", p.Password)
	values.Set("sms_code", "")
	values.Set("remember_cookie", "on")
	req, err := http.NewRequest("POST", "https://vpns.jlu.edu.cn/do-login?fromUrl=", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", p.Cookies)
	resp, err = DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile("logoutOtherToken = '([0-9a-f]+)'")
	if r.Match(body) {
		// Split current active token from html
		p.Cookies = "wengine_vpn_ticket_ecit=" + string(r.FindSubmatch(body)[1])
		log.Println(p.Cookies)
	}

	// TODO: Check whether logon successfully
	return nil
}
