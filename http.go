package main

import (
	"bytes"
	"compress/gzip"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host != "vpns.jlu.edu.cn" && strings.Contains(r.URL.Path, "wengine-vpn") {
		// wdnmd
		w.WriteHeader(200)
		return
	}

	var url string
	if PathMatchRegex.MatchString(r.URL.Path) {
		// Keep vpns
		if r.URL.Host == "vpns.jlu.edu.cn" {
			url = r.URL.String()
		} else {
			ret := PathMatchRegex.FindStringSubmatch(r.URL.Path)
			protocol := ret[1]
			host := ret[2]
			path := ret[3]
			r.URL.Scheme = "https"
			r.URL.Host = "vpns.jlu.edu.cn"
			r.URL.Path = "/" + protocol + "-" + r.URL.Port() + "/" + host + path
			url = r.URL.String()
		}
	} else if r.URL.Host == "vpns.jlu.edu.cn" {
		url = "https://vpns.jlu.edu.cn" + r.URL.Path
	} else {
		protocol := r.URL.Scheme
		host := r.URL.Hostname()
		if r.URL.Scheme == "" {
			// https
			protocol = "https"
			host = r.Host
		}
		// Without VPN
		r.URL.Path = "/" + protocol + "-" + r.URL.Port() + "/" + Encrypt(host) + r.URL.Path
		r.URL.Scheme = "https"
		r.URL.Host = "vpns.jlu.edu.cn"
		url = r.URL.String()
	}

	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		panic(err)
	}
	req.Header = r.Header

	resp, err := DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// Redirect
	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		if location == "/login" {
			// vpn relogin
			w.WriteHeader(resp.StatusCode)
		} else {
			location = RedirectLink.ReplaceAllStringFunc(location, func(s string) string {
				ret := RedirectLink.FindStringSubmatch(s)
				return ret[1] + "://" + Decrypt(ret[2]) + ret[3]
			})
			resp.Header.Set("Location", location)
		}
	}

	for k, v := range resp.Header {
		if k == "Content-Length" {
			continue
		}
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	if r.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Header.Get("Access-Control-Request-Headers") != "" {
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	}

	for _, c := range resp.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}

	w.WriteHeader(resp.StatusCode)

	body := resp.Body
	// Gzip
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
	}
	defer body.Close()

	result, err := ioutil.ReadAll(body)
	if err != nil && err != io.EOF {
		log.Println(err)
	}

	// Restore vpns url to original
	result = VPNsLinkMatch.ReplaceAllFunc(result, func(bytes []byte) []byte {
		ret := VPNsLinkMatch.FindSubmatch(bytes)
		return []byte(string(ret[1]) + "://" + Decrypt(string(ret[2])) + string(ret[3]))
	})

	// Fix URL Escape in links
	result = LinkUnescape.ReplaceAllFunc(result, func(bytes []byte) []byte {
		ret := LinkUnescape.FindSubmatch(bytes)
		return []byte(string(ret[1]) + "=\"" + html.UnescapeString(string(ret[2])) + "\"")
	})

	// Un VPN ify
	result = VPNEvalPrefix.ReplaceAll(result, []byte{})
	result = VPNEvalPostfix.ReplaceAll(result, []byte{})
	result = VPNRewritePrefix.ReplaceAll(result, []byte{})
	result = VPNRewritePostfix.ReplaceAll(result, []byte{})
	result = VPNParamRemoveFirst.ReplaceAll(result, []byte{})
	result = VPNParamRemoveOther.ReplaceAll(result, []byte{})

	// Remove added tags
	result = VPNScriptInfo.ReplaceAll(result, []byte{})
	result = VPNScriptWEngine.ReplaceAll(result, []byte(""))

	// Trim
	result = bytes.Trim(result, "\r\n")

	// Gzip
	if resp.Header.Get("Content-Encoding") == "gzip" {
		var buf bytes.Buffer
		wr := gzip.NewWriter(&buf)
		_, _ = wr.Write(result)
		_ = wr.Close()
		result = buf.Bytes()
	}

	if _, err = w.Write(result); err != nil {
		log.Println(err)
	}
}
