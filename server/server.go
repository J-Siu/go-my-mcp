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

package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/J-Siu/go-my-mcp/helper"
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
	// the SDK logs session and protocol activity through this logger
	t.Server = mcp.NewServer(&mcp.Implementation{
		Name:    "go-my-mcp",
		Version: helper.Version(),
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
	return fmt.Errorf("%s: %w", prefix, http.ListenAndServe(t.ListenAddress, t.Auth.Middleware(handler)))
}
