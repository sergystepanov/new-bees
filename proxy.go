package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

var client http.Client
var isFirst int32

func init() {
	jar, err := cookiejar.New(nil)

	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	ct := http.DefaultTransport.(*http.Transport).Clone()
	ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client = http.Client{
		Timeout:   30 * time.Second,
		Jar:       jar,
		Transport: ct,
	}
	atomic.StoreInt32(&isFirst, 0)
}

func proxy() func(w http.ResponseWriter, r *http.Request) {
	ua := envOr("UA", "")
	log.Printf("UA: %v", ua)
	genCookie := envOr("GEN_COOKIE", "") != ""
	log.Printf("GEN_COOKIE: %v", genCookie)
	setCookie := envOr("SET_COOKIE", "")
	log.Printf("SET_COOKIE: %v", setCookie)
	noDDG := envOr("NO_DDG", "") != ""
	log.Printf("NO_DDG: %v", noDDG)

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

		if !noDDG && strings.Contains(url_, "f"+"ips") {
			if genCookie {
				err = ddosGuardTokenized(u, &client)
			} else {
				err = ddosGuard(u, &client)
			}

			if err != nil {
				log.Printf("No DDG! %v", err)
			} else {
				log.Printf("Set DDG")
			}
		}

		req, err := http.NewRequest("GET", url_, nil)
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
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("has close error, %v", err)
				return
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("couldn't read, %v", err)
			return
		}

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, "%s", body)
	}
}
