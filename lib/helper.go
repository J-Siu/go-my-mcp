package lib

import (
	"fmt"
	"net/url"
	"strings"
)

// Check for http, https and empty host
func UrlParse(url_str string) (url_parsed *url.URL, err error) {
	url_parsed, err = url.Parse(strings.TrimSpace(url_str))
	if err != nil || (url_parsed.Scheme != "http" && url_parsed.Scheme != "https") ||
		url_parsed.Hostname() == "" || url_parsed.User != nil {
		err = fmt.Errorf("\"%s\" not a valid URL.", url_str)
	}
	return url_parsed, err
}
