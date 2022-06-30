package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"strings"
)

var client http.Client

func init() {
	jar, err := cookiejar.New(nil)

	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	ct := http.DefaultTransport.(*http.Transport).Clone()
	ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client = http.Client{
		Jar:       jar,
		Transport: ct,
	}
}

func proxy() func(w http.ResponseWriter, r *http.Request) {
	ua := envOr("UA", "")
	log.Printf("UA: %v", ua)
	genCookie := envOr("GEN_COOKIE", "") != ""
	log.Printf("GEN_COOKIE: %v", genCookie)

	return func(w http.ResponseWriter, r *http.Request) {
		// parse POST params
		if err := r.ParseForm(); err != nil {
			log.Printf("couldn't parse params, %v", err)
			return
		}

		url := r.Form.Get("_url")
		if url == "" {
			log.Printf("bad url!")
			return
		}
		u, err := url2.Parse(url)
		if err != nil {
			log.Printf("bad url: %v", err)
			return
		}

		if strings.Contains(url, "f"+"ips") {
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

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("bad url: %v", url)
			return
		}

		if ua != "" {
			req.Header.Set("User-Agent", ua)
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
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("couldn't read, %v", err)
			return
		}

		contentType := resp.Header.Get("Content-Type")
		w.Header().Set("Content-Type", contentType)
		_, _ = fmt.Fprintf(w, "%s", body)
	}
}
