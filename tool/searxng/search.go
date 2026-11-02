package searxng

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"fj.vms.local/jsds/go-helper-mini/ezlog"
	"fj.vms.local/jsds/go-helper-mini/str"

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
	if t.Err == nil {
		if t.Debug {
			ezlog.L.Println(prefix+":", input.Query)
		}
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
		if t.Debug {
			ezlog.L.Println(prefix+": httpUrl: ", httpUrl.String())
		}
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
	if t.Err == nil && t.Debug {
		ezlog.L.Println(str.Struct2String(prefix+": result", result))
	}
	if t.Err != nil {
		ezlog.L.Println(fmt.Errorf("%s: Error: %s", prefix, t.Err.Error()))
	}
	return nil, result, t.Err
}
