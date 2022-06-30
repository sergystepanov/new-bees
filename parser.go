package main

import (
	"net/url"
	"regexp"
	"strings"
)

const defaultCharset = "utf-8"

var metaRx = regexp.MustCompile("(?is)<meta charset=\"(.*?)\"")

func findCharset(text string) string {
	match := metaRx.FindStringSubmatch(text)
	if len(match) == 0 {
		return defaultCharset
	}
	return match[1]
}

func domain(url *url.URL, dotted bool) string {
	parts := strings.Split(url.Hostname(), ".")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		return ""
	}
	rez := ""
	if len(parts) > 2 && dotted {
		rez = "."
	}
	rez += parts[len(parts)-2] + "." + parts[len(parts)-1]
	return rez
}
