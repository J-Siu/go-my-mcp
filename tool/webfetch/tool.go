package webfetch

import (
	"fj.vms.local/jsds/go-helper-mini/basestruct"
	"fj.vms.local/jsds/go-helper-mini/ezlog"
	"fj.vms.local/jsds/go-helper-mini/str"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	bodySizeMax = 2 << 20
)

type TypeWebFetch struct {
	basestruct.Embedded
	BodySizeMax int64
}

func (t *TypeWebFetch) New() (searXng *TypeWebFetch, err error) {
	t.MyType = "TypeWebFetch"
	prefix := t.MyType + ".New"
	t.BodySizeMax = bodySizeMax
	t.Initialized = true
	if t.Debug {
		ezlog.L.Println(str.Struct2String(prefix, t))
	}
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
				Description: "fetch a web page by URL"},
			t.Fetch,
		)
		if t.Debug {
			ezlog.L.Println(prefix + ": registered")
		}
	}
}
