package searxng

import (
	"fmt"
	"net/url"

	"log/slog"

	"fj.vms.local/jsds/go-my-mcp/lib"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeSearXng struct {
	lib.Embedded
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
	t.url_parsed, t.Err = lib.UrlParse(searXngUrl)
	if t.Err == nil {
		t.Initialized = true
	} else {
		t.Err = fmt.Errorf("%s: %s", prefix, t.Err.Error())
	}
	slog.Debug(prefix, "config", slog.GroupValue(
		slog.String("SearXngUrl", t.SearXNG_URL),
		slog.Int64("BodySizeMax", t.BodySizeMax),
	))
	return t, t.Err
}

// Register SearXNG with mcpServer
func (t *TypeSearXng) Register(mcpServer *mcp.Server) {
	prefix := t.MyType + ".Register"
	if t.IsInitAndNoErr(prefix) {
		mcp.AddTool(
			mcpServer,
			&mcp.Tool{
				Name:        "searxng_search",
				Description: "web search"},
			t.Search,
		)
		slog.Debug(prefix + ": registered")
	}
}
