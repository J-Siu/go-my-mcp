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
	"net/url"

	"log/slog"

	"github.com/J-Siu/go-my-mcp/global"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeSearXng struct {
	Err         error  `json:"Err,omitempty"`
	Initialized bool   `json:"Initialized,omitempty"`
	MyType      string `json:"MyType"`

	SearXNG_URL string
	BodySizeMax int64
	url_parsed  *url.URL // Parsed URL object of SearXNG_URL
	TlsCheck    bool
}

func (t *TypeSearXng) New(searXngUrl string, bodySizeMax int64) (searXng *TypeSearXng, err error) {
	t.MyType = "TypeSearXng"
	prefix := t.MyType + ".New"
	t.BodySizeMax = bodySizeMax
	t.SearXNG_URL = searXngUrl
	t.url_parsed, t.Err = global.UrlParse(searXngUrl)
	if t.Err == nil {
		t.Initialized = true
		slog.Debug(prefix, "config", slog.GroupValue(
			slog.String("SearXngUrl", t.SearXNG_URL),
			slog.Int64("BodySizeMax", t.BodySizeMax),
		))
	} else {
		slog.Error("%s: %s", prefix, t.Err.Error())
	}
	return t, t.Err
}

// Register SearXNG with mcpServer
func (t *TypeSearXng) Register(mcpServer *mcp.Server) {
	prefix := t.MyType + ".Register"
	if t.Err != nil {
		return
	}
	if !t.Initialized {
		slog.Error("%s: %s", prefix, "not initialized")
		return
	}
	mcp.AddTool(
		mcpServer,
		&mcp.Tool{
			Name:        "searxng_search",
			Description: "web search"},
		t.Search,
	)
	slog.Debug("%s: %s", prefix, "registered")
}
