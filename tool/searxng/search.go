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
	var (
		body         []byte
		endpoint_str string
		httpClient   *http.Client
		httpReq      *http.Request
		httpRes      *http.Response
		httpUrl      *url.URL
		result       SearchOutput
	)

	t.Err = nil

	if strings.TrimSpace(input.Query) == "" {
		t.Err = fmt.Errorf("search query must not be empty")
	}
	// prepare search endpoint
	slog.Debug(prefix, "query", input.Query)
	if t.Err == nil {
		endpoint_str, t.Err = url.JoinPath(t.SearXNG_URL, "search")
	}
	if t.Err == nil {
		httpUrl, t.Err = url.Parse(endpoint_str)
	}
	// prepare http request
	if t.Err == nil {
		query := httpUrl.Query()
		query.Set("q", input.Query)
		query.Set("format", "json")
		httpUrl.RawQuery = query.Encode()
		slog.Debug(prefix, "url", httpUrl.String())
		httpReq, t.Err = http.NewRequestWithContext(ctx, http.MethodGet, httpUrl.String(), nil)
	}
	// http header and request
	if t.Err == nil {
		httpReq.Header.Set("Accept", "application/json")

		if t.TlsCheck {
			httpClient = http.DefaultClient
		} else {
			httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		httpRes, t.Err = httpClient.Do(httpReq)
		if t.Err != nil {
			t.Err = fmt.Errorf("httpClient.Do: %s", t.Err.Error())
		}
	}
	if t.Err == nil {
		defer httpRes.Body.Close()
		if httpRes.StatusCode < http.StatusOK || httpRes.StatusCode >= http.StatusMultipleChoices {
			t.Err = fmt.Errorf("SearXNG HTTP %s", httpRes.Status)
		}
	}
	if t.Err == nil {
		body, t.Err = io.ReadAll(io.LimitReader(httpRes.Body, t.BodySizeMax))
		if t.Err != nil {
			t.Err = fmt.Errorf("SearXNG HTTP %s: %s", httpRes.Status, strings.TrimSpace(string(body)))
		}
	}
	if t.Err == nil {
		t.Err = json.Unmarshal(body, &result)
		if t.Err != nil {
			t.Err = fmt.Errorf("json.Unmarshal: %s", t.Err.Error())
		}
	}
	if t.Err != nil {
		// the caller only gets the error text, so keep the details here
		slog.Error(prefix, "err", t.Err.Error())
	} else if slog.Default().Enabled(ctx, slog.LevelDebug) {
		slog.Debug(prefix, "results", len(result.Results))
	}
	return nil, result, t.Err
}
