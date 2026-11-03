package webfetch

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeWebFetch struct {
	Err         error  `json:"Err,omitempty"`
	Initialized bool   `json:"Initialized,omitempty"`
	MyType      string `json:"MyType"`

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
