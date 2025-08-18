package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AddTfstateReaderTools adds Terraform state analysis MCP tool implementations using real CLI commands
func AddTfstateReaderTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
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
		args := []string{"terraform", "state", "list"}
		
		if state := request.GetString("state", ""); state != "" {
			args = append(args, "-state="+state)
		}
		if id := request.GetString("id", ""); id != "" {
			args = append(args, id)
		}
		
		return executeShipCommand(args)
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
		resourceAddress := request.GetString("resource_address", "")
		args := []string{"terraform", "state", "show"}
		
		if state := request.GetString("state", ""); state != "" {
			args = append(args, "-state="+state)
		}
		args = append(args, resourceAddress)
		
		return executeShipCommand(args)
	})

	// Terraform state pull tool
	statePullTool := mcp.NewTool("terraform_state_pull",
		mcp.WithDescription("Download and output the state from remote backend using real terraform CLI"),
	)
	s.AddTool(statePullTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform", "state", "pull"}
		return executeShipCommand(args)
	})

	// Terraform show JSON state tool
	showJsonTool := mcp.NewTool("terraform_show_json",
		mcp.WithDescription("Show state in JSON format using real terraform CLI"),
		mcp.WithString("state_path",
			mcp.Description("Path to state file (optional)"),
		),
	)
	s.AddTool(showJsonTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"terraform", "show", "-json"}
		
		if statePath := request.GetString("state_path", ""); statePath != "" {
			args = append(args, statePath)
		}
		
		return executeShipCommand(args)
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
		args := []string{"tfstate-lookup"}
		
		if state := request.GetString("state_file", ""); state != "" {
			args = append(args, "-s", state)
		}
		if request.GetBool("interactive", false) {
			args = append(args, "-i")
		}
		if request.GetBool("dump", false) {
			args = append(args, "-dump")
		}
		if request.GetBool("jid", false) {
			args = append(args, "-j")
		}
		
		if !request.GetBool("interactive", false) && !request.GetBool("dump", false) {
			resourceAddress := request.GetString("resource_address", "")
			if resourceAddress != "" {
				args = append(args, resourceAddress)
			}
		}
		
		return executeShipCommand(args)
	})

	// tfstate-lookup dump all tool
	tfstateDumpTool := mcp.NewTool("tfstate_dump_all",
		mcp.WithDescription("Dump all resources from tfstate using real tfstate-lookup CLI"),
		mcp.WithString("state_file",
			mcp.Description("Path to state file or URL"),
		),
	)
	s.AddTool(tfstateDumpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tfstate-lookup", "-dump"}
		
		if state := request.GetString("state_file", ""); state != "" {
			args = append(args, "-s", state)
		}
		
		return executeShipCommand(args)
	})

	// tfstate-lookup interactive tool
	tfstateInteractiveTool := mcp.NewTool("tfstate_interactive",
		mcp.WithDescription("Browse tfstate interactively using real tfstate-lookup CLI"),
		mcp.WithString("state_file",
			mcp.Description("Path to state file or URL"),
		),
	)
	s.AddTool(tfstateInteractiveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"tfstate-lookup", "-i"}
		
		if state := request.GetString("state_file", ""); state != "" {
			args = append(args, "-s", state)
		}
		
		return executeShipCommand(args)
	})
}