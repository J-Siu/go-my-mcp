package lib

import (
	"fmt"
	"net/url"
	"strings"
)

// Check for http, https and empty host
func UrlParse(url_str string) (url_parsed *url.URL, err error) {
	url_str = strings.TrimSpace(url_str)
	url_parsed, err = url.Parse(url_str)
	if err != nil {
		return url_parsed, fmt.Errorf("%q is not a valid URL: %w", url_str, err)
	}
	switch {
	case url_parsed.Scheme != "http" && url_parsed.Scheme != "https":
		err = fmt.Errorf("%q is not a valid URL: scheme must be http or https, got %q", url_str, url_parsed.Scheme)
	case url_parsed.Hostname() == "":
		err = fmt.Errorf("%q is not a valid URL: missing host", url_str)
	case url_parsed.User != nil:
		err = fmt.Errorf("%q is not a valid URL: must not contain credentials", url_str)
	}
	return url_parsed, err
}
