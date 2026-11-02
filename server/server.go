package server

import (
	"errors"
	"log/slog"
	"net/http"

	"fj.vms.local/jsds/go-helper-mini/basestruct"
	"fj.vms.local/jsds/go-my-mcp/global"
	"fj.vms.local/jsds/go-my-mcp/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeServer struct {
	basestruct.Embedded

	ListenAddress string
	McpServer     *mcp.Server
}

func (t *TypeServer) New(listenAddress string) *TypeServer {
	t.MyType = "TypeServer"
	// prefix := t.MyType + ".New"
	t.McpServer = mcp.NewServer(&mcp.Implementation{
		Name:    "go-my-mcp",
		Version: global.Version(),
	}, nil)
	t.ListenAddress = listenAddress
	return t
}

// Register tool
func (t *TypeServer) Register(tool tool.Tool) *TypeServer {
	tool.Register(t.McpServer)
	return t
}

func (t *TypeServer) Run() error {
	prefix := t.MyType + ".Run"
	// MCP server using mcp go-sdk
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return t.McpServer
	}, nil)
	slog.Debug(prefix, "listen", t.ListenAddress)
	t.Err = http.ListenAndServe(t.ListenAddress, handler)
	if t.Err != nil {
		return errors.New(prefix + ": " + t.Err.Error())
	}
	return nil
}
