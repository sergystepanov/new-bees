package main

import (
	"net/url"
	"testing"
)

func TestFindCharset(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "empty", text: "", want: defaultCharset},
		{name: "normal", text: "<meta charset=\"cp1251\">", want: "cp1251"},
		{name: "case", text: "<MeTA ChaRseT=\"CP1251\">", want: "CP1251"},
		{
			name: "multiline",
			text: `
<html>
<meta charset="cp1251">
<meta charset="cp2020">
</html>
`,
			want: "cp1251",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findCharset(tt.text); got != tt.want {
				t.Errorf("findCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_domain(t *testing.T) {
	toURL := func(s string) url.URL { r, _ := url.Parse(s); return *r }
	type args struct {
		url    url.URL
		dotted bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "a", args: args{url: toURL("https://www.site.com"), dotted: false}, want: "site.com"},
		{name: "b", args: args{url: toURL("https://www.site.com"), dotted: true}, want: ".site.com"},
		{name: "c", args: args{url: toURL("site.com"), dotted: true}, want: ""},
		{name: "d", args: args{url: toURL("https://www1.site.com"), dotted: true}, want: ".site.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domain(&tt.args.url, tt.args.dotted); got != tt.want {
				t.Errorf("domain() = %v, want %v", got, tt.want)
			}
		})
	}
}
