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

// AddTerraformerTools adds Terraformer MCP tool implementations using direct Dagger calls
func AddTerraformerTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTerraformerToolsDirect(s)
}

// addTerraformerToolsDirect adds Terraformer tools using direct Dagger module calls
func addTerraformerToolsDirect(s *server.MCPServer) {
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformerModule(client)

		// Get parameters
		provider := request.GetString("provider", "")
		resources := request.GetString("resources", "")
		regions := request.GetString("regions", "")

		// Convert comma-separated resources to slice
		var resourceList []string
		if resources != "" && resources != "*" {
			resourceList = strings.Split(resources, ",")
			for i, r := range resourceList {
				resourceList[i] = strings.TrimSpace(r)
			}
		}

		// Build extra arguments map for unsupported parameters
		extraArgs := make(map[string]string)
		if excludes := request.GetString("excludes", ""); excludes != "" {
			extraArgs["excludes"] = excludes
		}
		if filter := request.GetString("filter", ""); filter != "" {
			extraArgs["filter"] = filter
		}
		if pathOutput := request.GetString("path_output", ""); pathOutput != "" {
			extraArgs["output"] = pathOutput
		}
		if outputFormat := request.GetString("output_format", ""); outputFormat != "" {
			extraArgs["output-format"] = outputFormat
		}
		if !request.GetBool("connect", true) {
			extraArgs["connect"] = "false"
		}
		if request.GetBool("verbose", false) {
			extraArgs["verbose"] = "true"
		}

		// Use the generic Import function
		output, err := module.Import(ctx, provider, regions, resourceList, extraArgs)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraformer import failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformerModule(client)

		// Get parameters
		provider := request.GetString("provider", "")

		// List resources for provider
		output, err := module.ListResources(ctx, provider)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraformer list resources failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
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
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformerModule(client)

		// Get parameters
		provider := request.GetString("provider", "")
		resources := request.GetString("resources", "")
		regions := request.GetString("regions", "")

		// Convert comma-separated resources to slice
		var resourceList []string
		if resources != "" && resources != "*" {
			resourceList = strings.Split(resources, ",")
			for i, r := range resourceList {
				resourceList[i] = strings.TrimSpace(r)
			}
		}

		// Note: Plan function doesn't support all the advanced parameters like filter/verbose
		// This is a limitation of the current Dagger module implementation
		if filter := request.GetString("filter", ""); filter != "" {
			return mcp.NewToolResultError("Warning: filter parameter is not supported in plan mode with direct Dagger calls"), nil
		}
		if request.GetBool("verbose", false) {
			return mcp.NewToolResultError("Warning: verbose parameter is not supported in plan mode with direct Dagger calls"), nil
		}

		// Generate plan
		output, err := module.Plan(ctx, provider, regions, resourceList)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraformer plan failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraformer version tool
	versionTool := mcp.NewTool("terraformer_version",
		mcp.WithDescription("Get Terraformer version information using real terraformer CLI"),
	)
	s.AddTool(versionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTerraformerModule(client)

		// Get version
		output, err := module.GetVersion(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraformer get version failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}