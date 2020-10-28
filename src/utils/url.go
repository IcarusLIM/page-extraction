package utils

import (
	urlParsed "net/url"
	"regexp"
)

var baseUrlRe = regexp.MustCompile(`(?i)<base\s[^>]*href\s*=\s*[\"\']\s*([^\"\'\s]+)\s*[\"\']`)

func AbsoluteURL(base *urlParsed.URL, url string) string {
	absUrl, err := base.Parse(url)
	if err != nil {
		return ""
	}
	absUrl.Fragment = ""
	if absUrl.Scheme == "//" {
		absUrl.Scheme = base.Scheme
	}
	return absUrl.String()
}

func GetBaseUrl(url string, html string) (base *urlParsed.URL, err error) {
	base, err = urlParsed.Parse(url)
	if err != nil {
		return
	}
	var l int
	if l = len(html); l > 4096 {
		l = 4096
	}
	if m := baseUrlRe.FindString(html[0:l]); m != "" {
		base, err = base.Parse(m)
	}
	return
}
