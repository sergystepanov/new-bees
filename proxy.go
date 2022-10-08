package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var client http.Client

func init() {
	jar, err := cookiejar.New(nil)

	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	ct := http.DefaultTransport.(*http.Transport).Clone()
	ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client = http.Client{Timeout: 60 * time.Second, Jar: jar, Transport: ct}
}

const (
	ContentType = "Content-Type"
	GET         = "GET"
	ips         = "f" + "ips"
)

func proxy() func(w http.ResponseWriter, r *http.Request) {
	ua := envOr("UA", "")
	setCookie := envOr("SET_COOKIE", "")
	noDDG := envOr("NO_DDG", "") != ""
	log.Printf("PROXY: UA=%v | Cookie=%v | NO_DDG=%v", ua, setCookie, noDDG)

	return func(w http.ResponseWriter, r *http.Request) {
		// parse POST params
		if err := r.ParseForm(); err != nil {
			log.Printf("couldn't parse params, %v", err)
			return
		}

		url_ := r.Form.Get("_url")
		if url_ == "" {
			log.Printf("bad url!")
			return
		}
		u, err := url.Parse(url_)
		if err != nil {
			log.Printf("bad url: %v", err)
			return
		}

		if !noDDG && strings.Contains(url_, ips) {
			err = ddosGuard(u, &client)
			//log.Printf("DDG: err?=%v", err)
		}

		req, err := http.NewRequest(GET, url_, nil)
		if err != nil {
			log.Printf("bad url: %v", url_)
			return
		}

		if ua != "" {
			req.Header.Set("User-Agent", ua)
		}
		if setCookie != "" {
			req.Header.Set("cookie", setCookie)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("couldn't get, %v", err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			log.Printf("couldn't read, %v", err)
			return
		}

		w.Header().Set(ContentType, resp.Header.Get(ContentType))
		_, _ = w.Write(body)
	}
}
