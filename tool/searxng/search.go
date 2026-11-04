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

package searxng

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchResult struct {
	Content string `json:"content,omitempty" jsonschema:"result summary"`
	Engine  string `json:"engine,omitempty" jsonschema:"search engine"`
	Title   string `json:"title,omitempty" jsonschema:"result title"`
	URL     string `json:"url" jsonschema:"result URL"`
}

type SearchInput struct {
	Query string `json:"query" jsonschema:"the search query"`
}

type SearchOutput struct {
	Results []SearchResult `json:"results" jsonschema:"search results"`
}

func (t *TypeSearXng) Search(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult,
	SearchOutput,
	error,
) {
	prefix := t.MyType + ".Search"
	start := time.Now()
	// The request state lives in locals: t.Err is shared by every concurrent
	// call on this long-lived tool.
	var (
		body         []byte
		endpoint_str string
		err          error
		httpClient   *http.Client
		httpReq      *http.Request
		httpRes      *http.Response
		httpUrl      *url.URL
		result       SearchOutput
	)

	if strings.TrimSpace(input.Query) == "" {
		err = fmt.Errorf("search query must not be empty")
	}
	// prepare search endpoint
	slog.Debug(prefix, "query", input.Query)
	if err == nil {
		endpoint_str, err = url.JoinPath(t.SearXNG_URL, "search")
	}
	if err == nil {
		httpUrl, err = url.Parse(endpoint_str)
	}
	// prepare http request
	if err == nil {
		query := httpUrl.Query()
		query.Set("q", input.Query)
		query.Set("format", "json")
		httpUrl.RawQuery = query.Encode()
		slog.Debug(prefix, "url", httpUrl.String())
		httpReq, err = http.NewRequestWithContext(ctx, http.MethodGet, httpUrl.String(), nil)
	}
	// http header and request
	if err == nil {
		httpReq.Header.Set("Accept", "application/json")

		if t.TlsCheck {
			httpClient = http.DefaultClient
		} else {
			httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		httpRes, err = httpClient.Do(httpReq)
		if err != nil {
			err = fmt.Errorf("httpClient.Do: %w", err)
		}
	}
	if err == nil {
		defer httpRes.Body.Close()
		if httpRes.StatusCode < http.StatusOK || httpRes.StatusCode >= http.StatusMultipleChoices {
			err = fmt.Errorf("SearXNG HTTP %s", httpRes.Status)
		}
	}
	if err == nil {
		body, err = io.ReadAll(io.LimitReader(httpRes.Body, t.BodySizeMax))
		if err != nil {
			err = fmt.Errorf("SearXNG HTTP %s: %s: %w", httpRes.Status, strings.TrimSpace(string(body)), err)
		}
	}
	if err == nil {
		err = json.Unmarshal(body, &result)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal: %w", err)
		}
	}
	if err != nil {
		// the caller only gets the error text, so keep the details here
		slog.Error(prefix, "err", err.Error(), "elapsed_ms", time.Since(start).Milliseconds())
	} else if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug(prefix, "results", len(result.Results),
			"status", httpRes.StatusCode,
			"content_len", len(body),
			"elapsed_ms", time.Since(start).Milliseconds(),
		)
	}
	return nil, result, err
}
