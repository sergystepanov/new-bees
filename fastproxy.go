package main

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var client0 *fasthttp.Client

func init() {
	readTimeout, _ := time.ParseDuration("30s")
	writeTimeout, _ := time.ParseDuration("30s")
	maxIdleConnDuration, _ := time.ParseDuration("15m")
	client0 = &fasthttp.Client{
		Name:                          "Go-http-client/1.1",
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func fastProxy(ctx *fasthttp.RequestCtx) {
	url_ := ctx.FormValue("_url")
	if url_ == nil {
		log.Printf("bad url!")
		return
	}

	log.Printf("[%s] -> [%s]", ctx.Host(), url_)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(string(url_))
	req.Header.Set(fasthttp.HeaderAcceptEncoding, "gzip, br")
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := client0.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		fasthttp.ReleaseResponse(resp)
		log.Printf("couldn't read, %v", err)
		return
	}
	body := resp.Body()
	contentType := resp.Header.ContentEncoding()
	contentEncoding := resp.Header.ContentEncoding()
	fasthttp.ReleaseResponse(resp)

	ctx.SetContentType(string(contentType))
	if contentEncoding != nil {
		ctx.Response.Header.Set(fasthttp.HeaderContentEncoding, string(contentEncoding))
	}
	ctx.SetBody(body)

	log.Printf("<- %v (%v)", len(body), string(contentEncoding))
}
