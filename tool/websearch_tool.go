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

package tool

import (
	"net/url"

	"log/slog"

	"github.com/J-Siu/go-my-mcp/conf"
	"github.com/J-Siu/go-my-mcp/helper"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TypeWebSearch - using SearXNG
type TypeWebSearch struct {
	UrlParsed *url.URL // Parsed URL object of SearXNG_URL
}

// Register into a mcp.Server
func (t *TypeWebSearch) Register(mcpServer *mcp.Server) {
	prefix := "TypeSearXng.Register"
	var err error
	t.UrlParsed, err = helper.UrlParse(conf.Conf.SearXngUrl)
	if err == nil {
		slog.Debug(prefix, "config", slog.GroupValue(
			slog.String("SearXngUrl", conf.Conf.SearXngUrl),
			slog.Int64("BodySizeMax", conf.Conf.SearchBodySizeMax),
		))
	} else {
		slog.Error(prefix + ": " + err.Error())
	}
	if err != nil {
		return
	}
	mcp.AddTool(
		mcpServer,
		&mcp.Tool{
			Name:        "web_search",
			Description: "web search"},
		t.Search,
	)
}

func init() { Tools.Add(new(TypeWebSearch)) }
