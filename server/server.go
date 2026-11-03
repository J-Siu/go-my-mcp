package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/J-Siu/go-my-mcp/global"
	"github.com/J-Siu/go-my-mcp/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeServer struct {
	MyType string `json:"MyType"`

	ListenAddress string
	Server        *mcp.Server
	Auth          TypeAuth
}

// New - listenAddress is where to bind; apiKey, when not empty, is required
// from every client
func (t *TypeServer) New(listenAddress, apiKey string) *TypeServer {
	t.MyType = "TypeServer"
	// prefix := t.MyType + ".New"
	// the SDK logs session and protocol activity through this logger
	t.Server = mcp.NewServer(&mcp.Implementation{
		Name:    "go-my-mcp",
		Version: global.Version(),
	}, &mcp.ServerOptions{Logger: slog.Default()})
	t.ListenAddress = listenAddress
	t.Auth.Key = apiKey
	return t
}

// Register tool
func (t *TypeServer) Register(tool tool.Tool) *TypeServer {
	tool.Register(t.Server)
	return t
}

func (t *TypeServer) Run() error {
	prefix := t.MyType + ".Run"
	// MCP server using mcp go-sdk
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return t.Server
	}, nil)
	slog.Debug(prefix, "listen", t.ListenAddress, "auth", t.Auth.Key != "")
	if err := http.ListenAndServe(t.ListenAddress, t.Auth.Middleware(handler)); err != nil {
		return errors.New(prefix + ": " + err.Error())
	}
	return nil
}
