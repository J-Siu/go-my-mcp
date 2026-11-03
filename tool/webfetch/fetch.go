package webfetch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/J-Siu/go-my-mcp/global"
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
	start := time.Now()
	// The request state lives in locals: t.Err is shared by every concurrent
	// call on this long-lived tool.
	var (
		body       []byte
		err        error
		httpReq    *http.Request
		httpRes    *http.Response
		redirects  int
		result     FetchOutput
		url_target *url.URL
	)

	url_target, err = global.UrlParse(input.URL)
	if err == nil {
		slog.Debug(prefix, "url", input.URL)
		result.URL = url_target.String()
		httpReq, err = http.NewRequestWithContext(ctx, http.MethodGet, url_target.String(), nil)
	}
	// http header and request
	if err == nil {
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
				if _, err := global.UrlParse(req.URL.String()); err != nil {
					return fmt.Errorf("redirect rejected: %w", err)
				}
				slog.Debug(prefix, "redirect", redirects, "to", req.URL.String())
				return nil
			},
		}
		httpRes, err = client.Do(httpReq)
		if err != nil {
			err = fmt.Errorf("httpClient.Do: %w", err)
		}
	}
	if err == nil {
		defer httpRes.Body.Close()

		result.URL = httpRes.Request.URL.String()
		result.StatusCode = httpRes.StatusCode
		result.ContentType = httpRes.Header.Get("Content-Type")
		result.RedirectCount = redirects

		if httpRes.StatusCode < http.StatusOK || httpRes.StatusCode >= http.StatusMultipleChoices {
			err = fmt.Errorf("HTTP %s", httpRes.Status)
		}
	}
	if err == nil {
		body, err = io.ReadAll(io.LimitReader(httpRes.Body, t.BodySizeMax+1))
	}
	if err == nil && len(body) > int(t.BodySizeMax) {
		err = fmt.Errorf("Exceeds the %d-byte limit", t.BodySizeMax)
	}
	if err != nil {
		// the caller only gets the error text, so keep the details here
		slog.Error(prefix, "status", result.StatusCode, "content_len", len(body),
			"content_type", result.ContentType, "err", err.Error(),
			"elapsed_ms", time.Since(start).Milliseconds())
	} else if slog.Default().Enabled(ctx, slog.LevelDebug) {
		attrs := []any{"status", result.StatusCode, "content_len", len(body),
			"elapsed_ms", time.Since(start).Milliseconds()}
		if result.ContentType != "" {
			attrs = append(attrs, "content_type", result.ContentType)
		}
		slog.Debug(prefix, attrs...)
	}
	result.Content = string(body)
	return nil, result, err
}
