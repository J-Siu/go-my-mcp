package tool

import "github.com/modelcontextprotocol/go-sdk/mcp"

type Tool interface {
	Register(mcpServer *mcp.Server)
}
