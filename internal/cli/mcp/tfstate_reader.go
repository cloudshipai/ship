package mcp

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"dagger.io/dagger"
)

// AddTfstateReaderTools adds Terraform state analysis MCP tool implementations using direct Dagger calls
func AddTfstateReaderTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// Ignore executeShipCommand - we use direct Dagger calls
	addTfstateReaderToolsDirect(s)
}

// addTfstateReaderToolsDirect adds Terraform state reader tools using direct Dagger module calls
func addTfstateReaderToolsDirect(s *server.MCPServer) {
	// Terraform state list resources tool
	stateListTool := mcp.NewTool("terraform_state_list",
		mcp.WithDescription("List all resources in Terraform state using real terraform CLI"),
		mcp.WithString("state",
			mcp.Description("Path to state file or remote state config"),
		),
		mcp.WithString("id",
			mcp.Description("Resource ID pattern to filter results"),
		),
	)
	s.AddTool(stateListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Get parameters
		statePath := request.GetString("state", "")
		resourceId := request.GetString("id", "")

		if statePath == "" {
			return mcp.NewToolResultError("state file path is required"), nil
		}

		// List state resources
		output, err := module.StateListResources(ctx, statePath, resourceId)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform state list failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform state show resource tool
	stateShowTool := mcp.NewTool("terraform_state_show",
		mcp.WithDescription("Show attributes of a resource in Terraform state using real terraform CLI"),
		mcp.WithString("resource_address",
			mcp.Description("Address of the resource to show"),
			mcp.Required(),
		),
		mcp.WithString("state",
			mcp.Description("Path to state file or remote state config"),
		),
	)
	s.AddTool(stateShowTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Get parameters
		resourceAddress := request.GetString("resource_address", "")
		statePath := request.GetString("state", "")

		if statePath == "" {
			return mcp.NewToolResultError("state file path is required"), nil
		}

		// Show state resource
		output, err := module.StateShowResource(ctx, statePath, resourceAddress)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform state show failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform state pull tool
	statePullTool := mcp.NewTool("terraform_state_pull",
		mcp.WithDescription("Download and output the state from remote backend using real terraform CLI"),
	)
	s.AddTool(statePullTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Pull remote state - requires a workspace directory (use current working directory)
		output, err := module.PullRemoteState(ctx, ".")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform state pull failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// Terraform show JSON state tool
	showJsonTool := mcp.NewTool("terraform_show_json",
		mcp.WithDescription("Show state in JSON format using real terraform CLI"),
		mcp.WithString("state_path",
			mcp.Description("Path to state file (optional)"),
		),
	)
	s.AddTool(showJsonTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Get parameters
		statePath := request.GetString("state_path", "")

		if statePath == "" {
			return mcp.NewToolResultError("state file path is required"), nil
		}

		// Show state in JSON format
		output, err := module.ShowState(ctx, statePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Terraform show JSON failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// tfstate-lookup query tool
	tfstateLookupTool := mcp.NewTool("tfstate_lookup_resource",
		mcp.WithDescription("Look up resource attributes in tfstate using real tfstate-lookup CLI"),
		mcp.WithString("resource_address",
			mcp.Description("Resource address to look up"),
			mcp.Required(),
		),
		mcp.WithString("state_file",
			mcp.Description("Path to state file or URL"),
		),
		mcp.WithBoolean("interactive",
			mcp.Description("Enable interactive mode"),
		),
		mcp.WithBoolean("dump",
			mcp.Description("Dump all resources, outputs, and data sources"),
		),
		mcp.WithBoolean("jid",
			mcp.Description("Run jid after selecting an item"),
		),
	)
	s.AddTool(tfstateLookupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Get parameters
		resourceAddress := request.GetString("resource_address", "")
		statePath := request.GetString("state_file", "")
		interactive := request.GetBool("interactive", false)
		dump := request.GetBool("dump", false)

		if statePath == "" {
			return mcp.NewToolResultError("state file path is required"), nil
		}

		// Note: jid parameter not supported with direct Dagger calls
		if request.GetBool("jid", false) {
			return mcp.NewToolResultError("Warning: jid parameter is not supported with direct Dagger calls"), nil
		}

		// Handle different modes
		var output string

		if dump {
			output, err = module.DumpAllResources(ctx, statePath)
		} else if interactive {
			output, err = module.InteractiveExplorer(ctx, statePath)
		} else {
			output, err = module.LookupResource(ctx, statePath, resourceAddress)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Tfstate lookup failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// tfstate-lookup dump all tool
	tfstateDumpTool := mcp.NewTool("tfstate_dump_all",
		mcp.WithDescription("Dump all resources from tfstate using real tfstate-lookup CLI"),
		mcp.WithString("state_file",
			mcp.Description("Path to state file or URL"),
		),
	)
	s.AddTool(tfstateDumpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Get parameters
		statePath := request.GetString("state_file", "")

		if statePath == "" {
			return mcp.NewToolResultError("state file path is required"), nil
		}

		// Dump all resources
		output, err := module.DumpAllResources(ctx, statePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Tfstate dump all failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})

	// tfstate-lookup interactive tool
	tfstateInteractiveTool := mcp.NewTool("tfstate_interactive",
		mcp.WithDescription("Browse tfstate interactively using real tfstate-lookup CLI"),
		mcp.WithString("state_file",
			mcp.Description("Path to state file or URL"),
		),
	)
	s.AddTool(tfstateInteractiveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Create Dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create Dagger client: %v", err)), nil
		}
		defer client.Close()

		// Create module instance
		module := modules.NewTfstateReaderModule(client)

		// Get parameters
		statePath := request.GetString("state_file", "")

		if statePath == "" {
			return mcp.NewToolResultError("state file path is required"), nil
		}

		// Interactive exploration
		output, err := module.InteractiveExplorer(ctx, statePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Tfstate interactive failed: %v", err)), nil
		}

		return mcp.NewToolResultText(output), nil
	})
}