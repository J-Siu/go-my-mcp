package webfetch

import (
	"log/slog"

	"fj.vms.local/jsds/go-my-mcp/lib"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeWebFetch struct {
	lib.Embedded
	BodySizeMax int64
	RedirectMax int
}

func (t *TypeWebFetch) New(bodySizeMax int64, redirectMax int) (webFetch *TypeWebFetch, err error) {
	t.MyType = "TypeWebFetch"
	prefix := t.MyType + ".New"
	t.BodySizeMax = bodySizeMax
	t.RedirectMax = redirectMax
	t.Initialized = true
	slog.Debug(prefix, "config", slog.GroupValue(
		slog.Int64("BodySizeMax", t.BodySizeMax),
		slog.Int("RedirectMax", t.RedirectMax),
	))
	return t, t.Err
}

// Register SearXNG with mcpServer
func (t *TypeWebFetch) Register(mcpServer *mcp.Server) {
	prefix := t.MyType + ".Register"
	if t.IsInitAndNoErr(prefix) {
		mcp.AddTool(
			mcpServer,
			&mcp.Tool{
				Name:        "web_fetch",
				Description: "fetch a web page by URL, following redirects"},
			t.Fetch,
		)
		slog.Debug(prefix + ": registered")
	}
}
