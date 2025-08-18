package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAllstarTools adds informational tools about Allstar (NOTE: Allstar is a GitHub App, not a CLI tool)
func AddAllstarTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Allstar info tool - explains that Allstar is not a CLI tool
	allstarInfoTool := mcp.NewTool("allstar_info",
		mcp.WithDescription("Get information about Allstar (Note: Allstar is a GitHub App service, not a CLI tool)"),
	)
	s.AddTool(allstarInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Return a simple execution that shows the help message via echo
		args := []string{"echo", "IMPORTANT: Allstar is a GitHub App (not a CLI tool). Install from https://github.com/apps/allstar and configure via .allstar/ repository."}
		return executeShipCommand(args)
	})
}