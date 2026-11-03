package global

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/charlievieth/strcase"
)

// Source - https://stackoverflow.com/a/77775482
// Posted by Eli Davis
// Retrieved 2026-07-19, License - CC BY-SA 4.0
func LoadJSON[T any](filename string) (T, error) {
	var data T
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}
	return data, json.Unmarshal(fileData, &data)
}

// Check if string contains any substring of an array
//   - result (bool)
//   - if result is true, `resultVal` == matching substring
//   - if result is false, `resultVal` == ""
func ContainsAnySubStrings(str string, subStrArrP *[]string, caseSensitive bool) (result bool, resultVal string) {
	// prefix := "matchList"
	if subStrArrP != nil {
		for _, subStr := range *subStrArrP {
			if caseSensitive && strcase.Contains(str, subStr) || strcase.Contains(str, subStr) {
				result = true
				resultVal = subStr
				break
			}
		}
	}
	return result, resultVal
}

// Return only bool from [ContainsAnySubStrings]
func ContainsAnySubStringsBool(str string, subStrArrP *[]string, caseSensitive bool) (result bool) {
	result, _ = ContainsAnySubStrings(str, subStrArrP, caseSensitive)
	return result
}

// Expand Linux `~` and environment variable in string
func TildeEnvExpand(strIn string) (strOut string) {
	if strIn == "~" {
		strOut = "$HOME"
	} else {
		re := regexp.MustCompile(`^~/`)
		strOut = re.ReplaceAllString(strIn, "$$HOME/")
	}
	return os.ExpandEnv(strOut)
}

// Check workPath is directory
func IsDir(path string) (result bool) {
	info, e := os.Stat(path)
	if e == nil {
		result = info.IsDir()
	}
	return result
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
