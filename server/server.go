package server

import (
	"errors"
	"log/slog"
	"net/http"

	"fj.vms.local/jsds/go-my-mcp/global"
	"fj.vms.local/jsds/go-my-mcp/lib"
	"fj.vms.local/jsds/go-my-mcp/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeServer struct {
	lib.Embedded

	ListenAddress string
	McpServer     *mcp.Server
	Auth          TypeAuth
}

// New - listenAddress is where to bind; apiKey, when not empty, is required
// from every client
func (t *TypeServer) New(listenAddress, apiKey string) *TypeServer {
	t.MyType = "TypeServer"
	// prefix := t.MyType + ".New"
	// the SDK logs session and protocol activity through this logger
	t.McpServer = mcp.NewServer(&mcp.Implementation{
		Name:    "go-my-mcp",
		Version: global.Version(),
	}, &mcp.ServerOptions{Logger: slog.Default()})
	t.ListenAddress = listenAddress
	t.Auth.Key = apiKey
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
	slog.Debug(prefix, "listen", t.ListenAddress, "auth", t.Auth.Key != "")
	t.Err = http.ListenAndServe(t.ListenAddress, t.Auth.Middleware(handler))
	if t.Err != nil {
		return errors.New(prefix + ": " + t.Err.Error())
	}
	return nil
}
