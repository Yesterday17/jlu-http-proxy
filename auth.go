package main

import (
	"io/ioutil"
	"regexp"
)

func (p *Proxy) Login() error {
	// Get Cookies first
	resp, err := DefaultClient.http2Client.Get("https://vpns.jlu.edu.cn/login")
	if err != nil {
		return err
	}

	// Auth
	resp, err = DefaultClient.http2Client.PostForm("https://vpns.jlu.edu.cn/do-login?fromUrl=", map[string][]string{
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
		// Need confirm login
		resp, err = DefaultClient.http2Client.PostForm("https://vpns.jlu.edu.cn/do-confirm-login", map[string][]string{
			"username":         {p.Username},
			"logoutOtherToken": {string(r.FindSubmatch(body)[1])},
		})
		if err != nil {
			return err
		}
	}

	// TODO: Check whether logon successfully
	return nil
}
