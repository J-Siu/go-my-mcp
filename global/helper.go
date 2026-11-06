/*
Copyright © 2026 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package global

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Expand Linux `~` and environment variable in string
func TildeEnvExpand(strIn string) (strOut string) {
	if strIn == "~" {
		strOut = "$HOME"
	} else if strings.HasPrefix(strIn, "~/") {
		strOut = "$HOME" + strIn[1:]
	} else {
		strOut = strIn
	}
	return os.ExpandEnv(strOut)
}

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
