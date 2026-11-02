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
	Content     string `json:"content" jsonschema:"the response body"`
	ContentType string `json:"content_type,omitempty" jsonschema:"the response content type"`
	StatusCode  int    `json:"status_code" jsonschema:"the HTTP status code"`
	URL         string `json:"url" jsonschema:"the fetched URL"`
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
		url_target *url.URL
	)

	t.Err = nil

	url_target, t.Err = lib.UrlParse(input.URL)
	if t.Err == nil {
		ezlog.L.Println(prefix+":", input.URL)
		httpReq, t.Err = http.NewRequestWithContext(ctx, http.MethodGet, url_target.String(), nil)
	}
	// http header and request
	if t.Err == nil {
		httpReq.Header.Set("Accept", "text/html, text/plain, application/json, */*")

		client := &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		httpRes, t.Err = client.Do(httpReq)
	}
	if t.Err == nil {
		defer httpRes.Body.Close()

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
	return nil,
		FetchOutput{
			URL:         url_target.String(),
			StatusCode:  httpRes.StatusCode,
			ContentType: httpRes.Header.Get("Content-Type"),
			Content:     string(body),
		},
		t.Err
}
