package webfetch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

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
		slog.Debug(prefix, "url", input.URL)
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
				slog.Debug(prefix, "redirect", redirects, "to", req.URL.String())
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
	if t.Err != nil {
		// the caller only gets the error text, so keep the details here
		slog.Error(prefix, "status", result.StatusCode, "content_len", len(body),
			"content_type", result.ContentType, "err", t.Err.Error())
	} else if slog.Default().Enabled(ctx, slog.LevelDebug) {
		attrs := []any{"status", result.StatusCode, "content_len", len(body)}
		if result.ContentType != "" {
			attrs = append(attrs, "content_type", result.ContentType)
		}
		slog.Debug(prefix, attrs...)
	}
	result.Content = string(body)
	return nil, result, t.Err
}
