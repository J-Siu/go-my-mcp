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

	"github.com/J-Siu/go-my-mcp/conf"
	"github.com/J-Siu/go-my-mcp/helper"
	"github.com/J-Siu/go-my-mcp/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeServer struct {
	Server *mcp.Server
	Auth   TypeAuth
}

// Run - register all tools in `tool.Tools` and start mcp server
func (t *TypeServer) Run() (err error) {
	prefix := "TypeServer.Run"

	t.Auth.Key = conf.Conf.ApiKey
	t.Server = mcp.NewServer(&mcp.Implementation{
		Name:    conf.AppName,
		Version: helper.Version(),
	}, &mcp.ServerOptions{Logger: slog.Default()})

	// Register all tools into Server
	slog.Debug("tool.Register", "Tool Number", len(tool.Tools))
	for _, tool := range tool.Tools {
		tool.Register(t.Server)
	}

	// MCP server using mcp go-sdk
	slog.Debug(prefix, "listen", conf.Conf.ListenAddress, "auth", conf.Conf.ApiKey != "")
	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return t.Server },
		nil,
	)
	err = http.ListenAndServe(conf.Conf.ListenAddress, t.Auth.Middleware(handler))
	return fmt.Errorf("%s: %w", prefix, err)
}

var Server TypeServer
