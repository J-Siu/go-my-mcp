package searxng

import (
	"fmt"
	"net/url"

	"fj.vms.local/jsds/go-helper-mini/basestruct"
	"fj.vms.local/jsds/go-helper-mini/ezlog"
	"fj.vms.local/jsds/go-helper-mini/str"
	"fj.vms.local/jsds/go-my-mcp/lib"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	bodySizeMax = 2 << 20
)

type TypeSearXng struct {
	basestruct.Embedded
	SearXNG_URL string
	BodySizeMax int64
	url_parsed  *url.URL // Parsed URL object of SearXNG_URL
	TlsCheck    bool
}

func (t *TypeSearXng) New(searXngUrl string) (searXng *TypeSearXng, err error) {
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
	if t.Debug {
		ezlog.L.Println(str.Struct2String(prefix, t))
	}
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
		if t.Debug {
			ezlog.L.Println(prefix + ": registered")
		}
	}
}
