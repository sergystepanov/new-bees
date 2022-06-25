package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const jsCheckURL = "https://check.ddos-guard.net/check.js"

func ddosGuard(u *url.URL, client *http.Client) error {
	req, err := http.NewRequest("GET", jsCheckURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	text := string(body)

	m := regexp.MustCompile(`new Image\(\).src = '(.+?)';`)
	res := m.FindStringSubmatch(text)

	if len(res) < 2 {
		return errors.New("no dgg URL")
	}

	log.Printf("Found ddg URL: %v", res[1])

	ddgURL := url.URL{Host: u.Host, Path: res[1], Scheme: u.Scheme}

	log.Printf("Result URL: %v", ddgURL)

	req, err = http.NewRequest("GET", ddgURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}
