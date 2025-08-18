package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTerraformerTools adds Terraformer MCP tool implementations using real terraformer CLI commands
func AddTerraformerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Terraformer import tool
	importTool := mcp.NewTool("terraformer_import",
		mcp.WithDescription("Import existing infrastructure to Terraform using real terraformer CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, gcp, azure, google, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("resources",
			mcp.Description("Comma-separated list of resources to import or '*' for all"),
			mcp.Required(),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of regions"),
		),
		mcp.WithString("excludes",
			mcp.Description("Comma-separated list of services to exclude"),
		),
		mcp.WithString("filter",
			mcp.Description("Filter resources by identifiers or attributes"),
		),
		mcp.WithString("path_output",
			mcp.Description("Set output directory (default 'generated')"),
		),
		mcp.WithString("output_format",
			mcp.Description("Output format"),
			mcp.Enum("hcl", "json"),
		),
		mcp.WithBoolean("connect",
			mcp.Description("Connect resources (default true)"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose mode"),
		),
	)
	s.AddTool(importTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		resources := request.GetString("resources", "")
		args := []string{"terraformer", "import", provider, "-r", resources}
		
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "-z", regions)
		}
		if excludes := request.GetString("excludes", ""); excludes != "" {
			args = append(args, "-x", excludes)
		}
		if filter := request.GetString("filter", ""); filter != "" {
			args = append(args, "-f", filter)
		}
		if pathOutput := request.GetString("path_output", ""); pathOutput != "" {
			args = append(args, "-o", pathOutput)
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			args = append(args, "-O", outputFormat)
		}
		if !request.GetBool("connect", true) {
			args = append(args, "-c", "false")
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		
		return executeShipCommand(args)
	})

	// Terraformer list resources tool
	listResourcesTool := mcp.NewTool("terraformer_list_resources",
		mcp.WithDescription("List supported resources for a provider using real terraformer CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, gcp, azure, google, etc.)"),
			mcp.Required(),
		),
	)
	s.AddTool(listResourcesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		args := []string{"terraformer", "import", provider, "list"}
		return executeShipCommand(args)
	})

	// Terraformer plan tool
	planTool := mcp.NewTool("terraformer_plan",
		mcp.WithDescription("Generate planfile for importing resources using real terraformer CLI"),
		mcp.WithString("provider",
			mcp.Description("Cloud provider (aws, gcp, azure, google, etc.)"),
			mcp.Required(),
		),
		mcp.WithString("resources",
			mcp.Description("Comma-separated list of resources to plan or '*' for all"),
			mcp.Required(),
		),
		mcp.WithString("regions",
			mcp.Description("Comma-separated list of regions"),
		),
		mcp.WithString("filter",
			mcp.Description("Filter resources by identifiers or attributes"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Enable verbose mode"),
		),
	)
	s.AddTool(planTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		provider := request.GetString("provider", "")
		resources := request.GetString("resources", "")
		args := []string{"terraformer", "plan", provider, "-r", resources}
		
		if regions := request.GetString("regions", ""); regions != "" {
			args = append(args, "-z", regions)
		}
		if filter := request.GetString("filter", ""); filter != "" {
			args = append(args, "-f", filter)
		}
		if request.GetBool("verbose", false) {
			args = append(args, "-v")
		}
		
		return executeShipCommand(args)
	})

	// Terraformer version tool
	versionTool := mcp.NewTool("terraformer_version",
		mcp.WithDescription("Get Terraformer version information using real terraformer CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraformer", "version"}
		return executeShipCommand(args)
	})
}