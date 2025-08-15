package ship

import (
	"context"
	"fmt"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPAdapter allows users to bring their own mcp-go server and add Ship tools to it
type MCPAdapter struct {
	engine   *dagger.Engine
	registry *Registry
}

// NewMCPAdapter creates a new adapter for integrating Ship tools into existing MCP servers
func NewMCPAdapter() *MCPAdapter {
	return &MCPAdapter{
		registry: NewRegistry(),
	}
}

// WithEngine sets the Dagger engine for the adapter
func (a *MCPAdapter) WithEngine(engine *dagger.Engine) *MCPAdapter {
	a.engine = engine
	return a
}

// AddTool adds a Ship tool to the adapter registry
func (a *MCPAdapter) AddTool(tool Tool) *MCPAdapter {
	if tool != nil {
		a.registry.RegisterTool(tool)
	}
	return a
}

// AddTools adds multiple Ship tools to the adapter registry
func (a *MCPAdapter) AddTools(tools ...Tool) *MCPAdapter {
	for _, tool := range tools {
		if tool != nil {
			a.registry.RegisterTool(tool)
		}
	}
	return a
}

// AddContainerTool creates and adds a container-based tool to the adapter
func (a *MCPAdapter) AddContainerTool(name string, config ContainerToolConfig) *MCPAdapter {
	tool := NewContainerTool(name, config)
	return a.AddTool(tool)
}

// ImportRegistry imports all tools from a Ship registry
func (a *MCPAdapter) ImportRegistry(registry *Registry) *MCPAdapter {
	if registry != nil {
		a.registry.ImportFrom(registry)
	}
	return a
}

// AttachToServer attaches all Ship tools to an existing mcp-go server
func (a *MCPAdapter) AttachToServer(ctx context.Context, mcpServer *server.MCPServer) error {
	// Initialize Dagger engine if not set
	if a.engine == nil {
		engine, err := dagger.NewEngine(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize dagger engine: %w", err)
		}
		a.engine = engine
	}

	// Register all Ship tools with the existing MCP server
	for toolName, tool := range a.registry.GetAllTools() {
		a.attachShipTool(mcpServer, toolName, tool)
	}

	return nil
}

// GetRegistry returns the adapter's registry
func (a *MCPAdapter) GetRegistry() *Registry {
	return a.registry
}

// GetEngine returns the adapter's Dagger engine
func (a *MCPAdapter) GetEngine() *dagger.Engine {
	return a.engine
}

// Close shuts down the adapter and cleans up resources
func (a *MCPAdapter) Close() error {
	if a.engine != nil {
		return a.engine.Close()
	}
	return nil
}

// attachShipTool registers a Ship tool with an mcp-go server
func (a *MCPAdapter) attachShipTool(mcpServer *server.MCPServer, name string, tool Tool) {
	// Convert framework parameters to MCP parameters
	var mcpOptions []mcp.ToolOption
	mcpOptions = append(mcpOptions, mcp.WithDescription(tool.Description()))

	// Add parameters
	for _, param := range tool.Parameters() {
		switch param.Type {
		case "string":
			if len(param.Enum) > 0 {
				mcpOptions = append(mcpOptions, mcp.WithString(param.Name,
					mcp.Description(param.Description),
					mcp.Enum(param.Enum...),
				))
			} else {
				mcpOptions = append(mcpOptions, mcp.WithString(param.Name,
					mcp.Description(param.Description),
				))
			}
		case "boolean":
			mcpOptions = append(mcpOptions, mcp.WithBoolean(param.Name,
				mcp.Description(param.Description),
			))
		case "number":
			// Handle as string for now, can be converted in the handler
			mcpOptions = append(mcpOptions, mcp.WithString(param.Name,
				mcp.Description(param.Description),
			))
		}
	}

	mcpTool := mcp.NewTool(name, mcpOptions...)

	// Add tool handler
	mcpServer.AddTool(mcpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters from request
		params := make(map[string]interface{})
		for _, param := range tool.Parameters() {
			switch param.Type {
			case "string":
				if value := request.GetString(param.Name, ""); value != "" || !param.Required {
					params[param.Name] = value
				}
			case "boolean":
				params[param.Name] = request.GetBool(param.Name, false)
			case "number":
				if value := request.GetString(param.Name, ""); value != "" {
					params[param.Name] = value
				}
			}
		}

		// Validate required parameters
		for _, param := range tool.Parameters() {
			if param.Required {
				if _, exists := params[param.Name]; !exists {
					return mcp.NewToolResultError(fmt.Sprintf("required parameter '%s' is missing", param.Name)), nil
				}
			}
		}

		// Execute the Ship tool
		result, err := tool.Execute(ctx, params, a.engine)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if result.Error != nil {
			return mcp.NewToolResultError(result.Error.Error()), nil
		}

		return mcp.NewToolResultText(result.Content), nil
	})
}

// ToolRouter provides a way to selectively route tools to different MCP servers
type ToolRouter struct {
	adapters map[string]*MCPAdapter
}

// NewToolRouter creates a new tool router
func NewToolRouter() *ToolRouter {
	return &ToolRouter{
		adapters: make(map[string]*MCPAdapter),
	}
}

// AddRoute adds a route for a specific tool or tool pattern to an adapter
func (r *ToolRouter) AddRoute(pattern string, adapter *MCPAdapter) *ToolRouter {
	r.adapters[pattern] = adapter
	return r
}

// RouteToServer routes tools based on patterns to different MCP servers
func (r *ToolRouter) RouteToServer(ctx context.Context, routes map[string]*server.MCPServer) error {
	for pattern, adapter := range r.adapters {
		if mcpServer, exists := routes[pattern]; exists {
			if err := adapter.AttachToServer(ctx, mcpServer); err != nil {
				return fmt.Errorf("failed to route %s: %w", pattern, err)
			}
		}
	}
	return nil
}
