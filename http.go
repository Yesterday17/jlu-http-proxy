package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		// HTTPS
		NewHTTPSRequest(w, r)
	} else {
		// HTTP
		NewHTTPRequest(w, r)
	}
}

func NewHTTPRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host != "vpns.jlu.edu.cn" && strings.Contains(r.URL.Path, "wengine-vpn") {
		// wdnmd
		w.WriteHeader(200)
		return
	}

	var url string
	reg := regexp.MustCompile("/(https?)(?:-\\d*)?/([0-9a-f]+)(.+)")
	if reg.MatchString(r.URL.Path) {
		// Keep vpns
		if r.URL.Host == "vpns.jlu.edu.cn" {
			url = r.URL.String()
		} else {
			ret := reg.FindStringSubmatch(r.URL.Path)
			url = ret[1] + "://" + Decrypt(ret[2]) + ret[3]
		}
	} else if r.URL.Host == "vpns.jlu.edu.cn" {
		url = "https://vpns.jlu.edu.cn" + r.URL.Path
	} else {
		// Without VPN
		url = "https://vpns.jlu.edu.cn/http-" + r.URL.Port() + "/" + Encrypt(r.URL.Hostname()) + r.URL.Path
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

	// Replace https
	re := regexp.MustCompile("https://")
	result = re.ReplaceAll(result, []byte("http://"))

	// Restore vpns url to original
	re = regexp.MustCompile("vpns.jlu.edu.cn/https?(?:-[0-9]+)?/([0-9a-f]+)([^\")]+)")
	result = re.ReplaceAllFunc(result, func(bytes []byte) []byte {
		ret := re.FindSubmatch(bytes)
		return []byte(Decrypt(string(ret[1])) + string(ret[2]))
	})

	// Un VPN ify
	re = regexp.MustCompile("vpn_eval\\(\\(function\\(\\){\r?\n")
	result = re.ReplaceAll(result, []byte(""))

	re = regexp.MustCompile("}\r?\n\\).toString\\(\\).slice\\(12, -2\\),\"\"\\);")
	result = re.ReplaceAll(result, []byte(""))

	re = regexp.MustCompile("vpn-\\d+&")
	result = re.ReplaceAll(result, []byte(""))

	re = regexp.MustCompile("\\?vpn-\\d+")
	result = re.ReplaceAll(result, []byte(""))

	// Remove added tags
	re = regexp.MustCompile("<script>(?:\r?\nvar __vpn_[^ ]+ = [^;]+;)+\r?\n</script>")
	result = re.ReplaceAll(result, []byte(""))
	re = regexp.MustCompile("<script src=\"/wengine-vpn/js/main.js[^>]+></script>")
	result = re.ReplaceAll(result, []byte(""))

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

func NewHTTPSRequest(w http.ResponseWriter, r *http.Request) {
	// TODO
}
