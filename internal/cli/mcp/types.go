package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// ExecuteShipCommandFunc is a function type for executing ship commands
type ExecuteShipCommandFunc func(args []string) (*mcp.CallToolResult, error)
