package main

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const jsCheckURL = "https://check.ddos-guard.net/check.js"
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var m = regexp.MustCompile(`new Image\(\).src = '(.+?)';`)
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

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

	res := m.FindStringSubmatch(text)

	if len(res) < 2 {
		return errors.New("no dgg URL")
	}

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

func ddosGuardTokenized(u *url.URL, client *http.Client) error {
	cookie := &http.Cookie{
		Domain:   domain(u, true),
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		MaxAge:   int((time.Hour * 24 * 365).Seconds()),
		Name:     "__ddg2_",
		HttpOnly: true,
		Path:     "/",
		Value:    token(16),
	}
	client.Jar.SetCookies(u, []*http.Cookie{cookie})
	return nil
}

func token(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
