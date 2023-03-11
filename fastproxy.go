package main

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var client0 *fasthttp.Client

var encoding = []byte("gzip, br")

func init() {
	dial0 := &fasthttp.TCPDialer{
		Concurrency:      512,
		DNSCacheDuration: time.Hour,
	}
	client0 = &fasthttp.Client{
		Name:                          "Go-http-client/1.1",
		ReadTimeout:                   30 * time.Second,
		WriteTimeout:                  30 * time.Second,
		MaxIdleConnDuration:           15 * time.Minute,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial:      dial0.Dial,
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
	req.SetRequestURIBytes(url_)
	req.Header.SetBytesV(fasthttp.HeaderAcceptEncoding, encoding)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := client0.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		fasthttp.ReleaseResponse(resp)
		log.Printf("couldn't read, %v", err)
		return
	}

	ctx.SetContentTypeBytes(resp.Header.ContentType())
	ctx.SetStatusCode(resp.StatusCode())
	ctx.Response.Header.SetBytesV(fasthttp.HeaderContentEncoding, resp.Header.ContentEncoding())
	body := resp.Body()
	ctx.Response.SetBodyRaw(body)
	log.Printf("<- %vb", len(body))
	fasthttp.ReleaseResponse(resp)
}
