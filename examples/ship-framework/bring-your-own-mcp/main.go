package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/cloudshipai/ship/pkg/ship"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Your existing custom tool
type CustomTool struct{}

func (t *CustomTool) Name() string {
	return "custom-hello"
}

func (t *CustomTool) Description() string {
	return "A custom tool that says hello"
}

func (t *CustomTool) Parameters() []ship.Parameter {
	return []ship.Parameter{
		{
			Name:        "name",
			Type:        "string",
			Description: "Name to greet",
			Required:    true,
		},
	}
}

func (t *CustomTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	name, _ := params["name"].(string)
	return &ship.ToolResult{
		Content: fmt.Sprintf("Hello, %s! This is a custom tool.", name),
	}, nil
}

// Your custom MCP tool (not using Ship framework)
func addCustomMCPTool(s *server.MCPServer) {
	tool := mcp.NewTool("native-tool", 
		mcp.WithDescription("A native MCP tool not using Ship framework"),
		mcp.WithString("message", mcp.Description("Message to process")),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		message := request.GetString("message", "")
		result := fmt.Sprintf("Native MCP tool processed: %s", message)
		return mcp.NewToolResultText(result), nil
	})
}

func main() {
	fmt.Println("Starting MCP server with mixed tools...")

	// Create your own mcp-go server
	mcpServer := server.NewMCPServer("mixed-tools-server", "1.0.0")

	// Add your existing native MCP tools
	addCustomMCPTool(mcpServer)

	// Create Ship adapter to add Ship tools to your existing server
	adapter := ship.NewMCPAdapter().
		AddTool(&CustomTool{})

	// Attach Ship tools to your existing MCP server
	ctx := context.Background()
	if err := adapter.AttachToServer(ctx, mcpServer); err != nil {
		log.Fatalf("Failed to attach Ship tools: %v", err)
	}

	// Clean up adapter when done
	defer adapter.Close()

	fmt.Println("Available tools:")
	fmt.Println("- native-tool (pure mcp-go)")
	fmt.Println("- custom-hello (Ship framework tool)")

	// Start your MCP server with both types of tools
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}