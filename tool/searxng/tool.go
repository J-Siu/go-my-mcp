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
