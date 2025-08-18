package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddInfraMapTools adds InfraMap (infrastructure diagram generator) MCP tool implementations using direct Dagger calls
func AddInfraMapTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addInfraMapToolsDirect(s)
}

// addInfraMapToolsDirect adds InfraMap tools using direct Dagger module calls
func addInfraMapToolsDirect(s *server.MCPServer) {
	// InfraMap generate tool
	generateTool := mcp.NewTool("inframap_generate",
		mcp.WithDescription("Generate infrastructure graph using inframap generate"),
		mcp.WithString("input",
			mcp.Description("Path to Terraform state file, HCL file, or directory"),
			mcp.Required(),
		),
		mcp.WithBoolean("hcl",
			mcp.Description("Force HCL input type"),
		),
		mcp.WithBoolean("tfstate",
			mcp.Description("Force Terraform state input type"),
		),
		mcp.WithBoolean("connections",
			mcp.Description("Enable connections in graph (default: true)"),
		),
		mcp.WithBoolean("raw",
			mcp.Description("Show configuration without InfraMap processing"),
		),
		mcp.WithBoolean("clean",
			mcp.Description("Remove unconnected nodes (default: true)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: dot, png, svg, pdf (default: dot)"),
			mcp.Enum("dot", "png", "svg", "pdf"),
		),
		mcp.WithString("provider",
			mcp.Description("Filter by provider (aws, google, azurerm, etc.)"),
		),
	)
	s.AddTool(generateTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfraMapModule(client)

		// Get parameters
		input := request.GetString("input", "")
		if input == "" {
			return mcp.NewToolResultError("input is required"), nil
		}

		// Determine if input is HCL or tfstate
		isHCL := request.GetBool("hcl", false)
		isTFState := request.GetBool("tfstate", false)
		format := request.GetString("format", "dot")

		// Check if we need to use GenerateWithOptions for advanced features
		raw := request.GetBool("raw", false)
		clean := request.GetBool("clean", true)
		provider := request.GetString("provider", "")
		
		// If we have options specified, use GenerateWithOptions
		if raw || !clean || provider != "" {
			opts := modules.InfraMapOptions{
				Raw:      raw,
				Clean:    clean,
				Provider: provider,
				Format:   format,
			}
			
			output, err := module.GenerateWithOptions(ctx, input, opts)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to generate diagram: %v", err)), nil
			}
			return mcp.NewToolResultText(output), nil
		}

		// Otherwise use the specific generate functions
		var output string
		if isHCL || (!isTFState && strings.HasSuffix(input, "/")) {
			// If HCL is specified or input is a directory, use GenerateFromHCL
			output, err = module.GenerateFromHCL(ctx, input, format)
		} else {
			// Otherwise assume it's a state file
			output, err = module.GenerateFromState(ctx, input, format)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to generate diagram: %v", err)), nil
		}

		// Add note about connections if disabled
		if !request.GetBool("connections", true) {
			output = "Note: Connections disabled in graph\n\n" + output
		}

		return mcp.NewToolResultText(output), nil
	})

	// InfraMap prune tool
	pruneTool := mcp.NewTool("inframap_prune",
		mcp.WithDescription("Remove unnecessary information from Terraform state/HCL using inframap prune"),
		mcp.WithString("input",
			mcp.Description("Path to Terraform state file or HCL file to prune"),
			mcp.Required(),
		),
	)
	s.AddTool(pruneTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewInfraMapModule(client)

		// Get input
		input := request.GetString("input", "")
		if input == "" {
			return mcp.NewToolResultError("input is required"), nil
		}

		// Prune the state/HCL
		output, err := module.PruneState(ctx, input)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to prune: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}