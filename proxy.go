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

	client = http.Client{Timeout: 40 * time.Second, Jar: jar, Transport: ct}
}

const (
	ContentEncoding = "Content-Encoding"
	ContentType     = "Content-Type"
	GET             = "GET"
	ips             = "f" + "ips"
)

func proxy() func(w http.ResponseWriter, r *http.Request) {
	ua := envOr("UA", "")
	setCookie := envOr("SET_COOKIE", "")
	noDDG := envOr("NO_DDG", "") != ""
	log.Printf("PROXY: UA=%v | Cookie=%v | NO_DDG=%v", ua, setCookie, noDDG)

	return func(w http.ResponseWriter, r *http.Request) {
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

		log.Printf("-> [%v]", url_)
		req, err := http.NewRequest(GET, url_, nil)
		if err != nil {
			log.Printf("bad url: %v", url_)
			return
		}

		req.Header.Add("Accept-Encoding", "gzip, br")

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

		hh := w.Header()
		hh.Set(ContentType, resp.Header.Get(ContentType))
		encoding := resp.Header.Get(ContentEncoding)
		if encoding != "" {
			hh.Set(ContentEncoding, encoding)
		}
		_, _ = w.Write(body)
		log.Printf("<- %v (%v)", len(body), encoding)
	}
}
