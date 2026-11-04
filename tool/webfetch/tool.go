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

package webfetch

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeWebFetch struct {
	Initialized bool   `json:"Initialized,omitempty"`
	MyType      string `json:"MyType"`

	BodySizeMax int64
	RedirectMax int
}

func (t *TypeWebFetch) New(bodySizeMax int64, redirectMax int) (*TypeWebFetch, error) {
	t.MyType = "TypeWebFetch"
	prefix := t.MyType + ".New"
	t.BodySizeMax = bodySizeMax
	t.RedirectMax = redirectMax
	t.Initialized = true
	slog.Debug(prefix, "config", slog.GroupValue(
		slog.Int64("BodySizeMax", t.BodySizeMax),
		slog.Int("RedirectMax", t.RedirectMax),
	))
	return t, nil
}

// Register SearXNG with mcpServer
func (t *TypeWebFetch) Register(mcpServer *mcp.Server) {
	prefix := t.MyType + ".Register"
	if !t.Initialized {
		slog.Error("%s: %s", prefix, "not initialized")
		return
	}
	mcp.AddTool(
		mcpServer,
		&mcp.Tool{
			Name:        "web_fetch",
			Description: "fetch a web page by URL, following redirects"},
		t.Fetch,
	)
	slog.Debug("%s: %s", prefix, "registered")
}
