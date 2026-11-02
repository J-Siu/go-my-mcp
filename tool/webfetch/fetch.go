package webfetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"fj.vms.local/jsds/go-helper-mini/ezlog"
	"fj.vms.local/jsds/go-my-mcp/lib"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FetchInput struct {
	URL string `json:"url" jsonschema:"the HTTP or HTTPS URL to fetch"`
}

type FetchOutput struct {
	Content       string `json:"content" jsonschema:"the response body"`
	ContentType   string `json:"content_type,omitempty" jsonschema:"the response content type"`
	StatusCode    int    `json:"status_code" jsonschema:"the HTTP status code"`
	URL           string `json:"url" jsonschema:"the fetched URL, after following redirects"`
	RedirectCount int    `json:"redirect_count,omitempty" jsonschema:"number of redirects followed"`
}

func (t *TypeWebFetch) Fetch(ctx context.Context, req *mcp.CallToolRequest, input FetchInput) (
	*mcp.CallToolResult,
	FetchOutput,
	error,
) {
	prefix := t.MyType + ".Fetch"
	var (
		body       []byte
		httpReq    *http.Request
		httpRes    *http.Response
		redirects  int
		result     FetchOutput
		url_target *url.URL
	)

	t.Err = nil

	url_target, t.Err = lib.UrlParse(input.URL)
	if t.Err == nil {
		ezlog.L.Println(prefix+":", input.URL)
		result.URL = url_target.String()
		httpReq, t.Err = http.NewRequestWithContext(ctx, http.MethodGet, url_target.String(), nil)
	}
	// http header and request
	if t.Err == nil {
		httpReq.Header.Set("Accept", "text/html, text/plain, application/json, */*")

		client := &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if t.RedirectMax == 0 {
					return http.ErrUseLastResponse
				}
				redirects = len(via)
				if redirects >= t.RedirectMax {
					return fmt.Errorf("stopped after %d redirects", t.RedirectMax)
				}
				// each hop must obey the same rules as the requested URL
				if _, err := lib.UrlParse(req.URL.String()); err != nil {
					return fmt.Errorf("redirect rejected: %s", err.Error())
				}
				if t.Debug {
					ezlog.L.Println(prefix+": redirect", redirects, "->", req.URL.String())
				}
				return nil
			},
		}
		httpRes, t.Err = client.Do(httpReq)
		if t.Err != nil {
			t.Err = fmt.Errorf("httpClient.Do: %s", t.Err.Error())
		}
	}
	if t.Err == nil {
		defer httpRes.Body.Close()

		result.URL = httpRes.Request.URL.String()
		result.StatusCode = httpRes.StatusCode
		result.ContentType = httpRes.Header.Get("Content-Type")
		result.RedirectCount = redirects

		if httpRes.StatusCode < http.StatusOK || httpRes.StatusCode >= http.StatusMultipleChoices {
			t.Err = fmt.Errorf("HTTP %s", httpRes.Status)
		}
	}
	if t.Err == nil {
		body, t.Err = io.ReadAll(io.LimitReader(httpRes.Body, t.BodySizeMax+1))
	}
	if t.Err == nil && len(body) > int(t.BodySizeMax) {
		t.Err = fmt.Errorf("Exceeds the %d-byte limit", t.BodySizeMax)
	}
	if t.Err == nil && t.Debug {
		ezlog.L.Println(string(body))
	}
	if t.Err != nil {
		ezlog.L.Println(fmt.Errorf("%s: Error: %s", prefix, t.Err.Error()))
	}
	result.Content = string(body)
	return nil, result, t.Err
}
