package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddAllstarTools adds informational tools about Allstar (NOTE: Allstar is a GitHub App, not a CLI tool)
func AddAllstarTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addAllstarToolsDirect(s)
}

// addAllstarToolsDirect implements direct Dagger calls for allstar tools
func addAllstarToolsDirect(s *server.MCPServer) {
	// Allstar info tool - explains that Allstar is not a CLI tool
	allstarInfoTool := mcp.NewTool("allstar_info",
		mcp.WithDescription("Get information about Allstar (Note: Allstar is a GitHub App service, not a CLI tool)"),
	)
	s.AddTool(allstarInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create allstar module and get info
		allstarModule := modules.NewAllstarModule(client)
		result, err := allstarModule.GetInfo(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("allstar info failed: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})
}