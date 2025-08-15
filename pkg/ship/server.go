package ship

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer represents an MCP server with tools, prompts, and resources
type MCPServer struct {
	name     string
	version  string
	registry *Registry
	engine   *dagger.Engine
	server   *server.MCPServer
	proxies  map[string]*MCPProxy // For managing external MCP server connections
}

// ServerBuilder provides a fluent API for building MCP servers
type ServerBuilder struct {
	name     string
	version  string
	registry *Registry
}

// NewServer creates a new server builder
func NewServer(name, version string) *ServerBuilder {
	if name == "" {
		name = "unnamed-server"
	}
	if version == "" {
		version = "1.0.0"
	}

	return &ServerBuilder{
		name:     name,
		version:  version,
		registry: NewRegistry(),
	}
}

// AddTool adds a single tool to the server
func (b *ServerBuilder) AddTool(tool Tool) *ServerBuilder {
	if tool != nil {
		b.registry.RegisterTool(tool)
	}
	return b
}

// AddTools adds multiple tools to the server
func (b *ServerBuilder) AddTools(tools ...Tool) *ServerBuilder {
	for _, tool := range tools {
		if tool != nil {
			b.registry.RegisterTool(tool)
		}
	}
	return b
}

// AddContainerTool creates and adds a container-based tool
func (b *ServerBuilder) AddContainerTool(name string, config ContainerToolConfig) *ServerBuilder {
	tool := NewContainerTool(name, config)
	return b.AddTool(tool)
}

// AddPrompt adds a prompt to the server
func (b *ServerBuilder) AddPrompt(prompt Prompt) *ServerBuilder {
	if prompt != nil {
		b.registry.RegisterPrompt(prompt)
	}
	return b
}

// AddResource adds a resource to the server
func (b *ServerBuilder) AddResource(resource Resource) *ServerBuilder {
	if resource != nil {
		b.registry.RegisterResource(resource)
	}
	return b
}

// ImportRegistry imports all items from another registry
func (b *ServerBuilder) ImportRegistry(registry *Registry) *ServerBuilder {
	if registry != nil {
		b.registry.ImportFrom(registry)
	}
	return b
}

// Build creates the MCP server instance
func (b *ServerBuilder) Build() *MCPServer {
	return &MCPServer{
		name:     b.name,
		version:  b.version,
		registry: b.registry,
	}
}

// Start initializes the server with a Dagger engine
func (s *MCPServer) Start(ctx context.Context) error {
	// Initialize Dagger engine
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger engine: %w", err)
	}

	s.engine = engine

	// Create MCP server
	s.server = server.NewMCPServer(s.name, s.version)

	// Register all tools
	for toolName, tool := range s.registry.GetAllTools() {
		s.registerMCPTool(toolName, tool)
	}

	return nil
}

// ServeStdio starts the server using stdio transport
func (s *MCPServer) ServeStdio() error {
	ctx := context.Background()

	if err := s.Start(ctx); err != nil {
		return err
	}

	defer s.Close()

	return server.ServeStdio(s.server)
}

// ServeHTTP starts the server using HTTP transport
func (s *MCPServer) ServeHTTP(host string, port int) error {
	ctx := context.Background()

	if err := s.Start(ctx); err != nil {
		return err
	}

	defer s.Close()

	// Note: HTTP server implementation would go here
	// For now, return not implemented
	return fmt.Errorf("HTTP server not implemented yet")
}

// Close shuts down the server and cleans up resources
func (s *MCPServer) Close() error {
	var lastErr error
	
	// Close all proxy connections
	if s.proxies != nil {
		for name, proxy := range s.proxies {
			if err := proxy.Close(); err != nil {
				lastErr = fmt.Errorf("failed to close proxy %s: %w", name, err)
			}
		}
	}
	
	// Close Dagger engine
	if s.engine != nil {
		if err := s.engine.Close(); err != nil {
			lastErr = err
		}
	}
	
	return lastErr
}

// GetRegistry returns the server's registry
func (s *MCPServer) GetRegistry() *Registry {
	return s.registry
}

// GetEngine returns the server's Dagger engine (may be nil if not started)
func (s *MCPServer) GetEngine() *dagger.Engine {
	return s.engine
}

// GetMCPGoServer returns the underlying mcp-go server instance
func (s *MCPServer) GetMCPGoServer() *server.MCPServer {
	return s.server
}

// registerMCPTool registers a framework tool as an MCP tool
func (s *MCPServer) registerMCPTool(name string, tool Tool) {
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
			// Note: The actual API might be different, this is a placeholder
			// We'll need to check the real mcp-go API
			mcpOptions = append(mcpOptions, mcp.WithString(param.Name,
				mcp.Description(param.Description),
			))
		}
	}

	mcpTool := mcp.NewTool(name, mcpOptions...)

	// Add tool handler
	s.server.AddTool(mcpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
				// For now, treat numbers as strings since the API might not have GetNumber
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

		// Execute the tool
		result, err := tool.Execute(ctx, params, s.engine)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if result.Error != nil {
			return mcp.NewToolResultError(result.Error.Error()), nil
		}

		return mcp.NewToolResultText(result.Content), nil
	})
}

// WriteOutput writes server output to a writer (useful for testing)
func (s *MCPServer) WriteOutput(w io.Writer, message string) error {
	_, err := fmt.Fprintf(w, "[%s] %s\n", s.name, message)
	return err
}

// LogInfo logs an info message to stderr
func (s *MCPServer) LogInfo(message string) {
	fmt.Fprintf(os.Stderr, "[%s] INFO: %s\n", s.name, message)
}

// LogError logs an error message to stderr
func (s *MCPServer) LogError(message string) {
	fmt.Fprintf(os.Stderr, "[%s] ERROR: %s\n", s.name, message)
}
