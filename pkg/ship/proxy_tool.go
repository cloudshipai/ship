package ship

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// ProxyTool represents a tool that proxies calls to an external MCP server
type ProxyTool struct {
	name        string
	description string
	parameters  []Parameter
	client      client.MCPClient
	toolName    string // The actual tool name on the external server
}

// NewProxyTool creates a new proxy tool that forwards calls to an external MCP server
func NewProxyTool(name, description string, parameters []Parameter, mcpClient client.MCPClient, externalToolName string) *ProxyTool {
	return &ProxyTool{
		name:        name,
		description: description,
		parameters:  parameters,
		client:      mcpClient,
		toolName:    externalToolName,
	}
}

// Name returns the tool name
func (p *ProxyTool) Name() string {
	return p.name
}

// Description returns the tool description
func (p *ProxyTool) Description() string {
	return p.description
}

// Parameters returns the tool parameters
func (p *ProxyTool) Parameters() []Parameter {
	return p.parameters
}

// Execute forwards the tool call to the external MCP server
func (p *ProxyTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
	// Convert parameters to the format expected by mcp-go
	arguments := make(map[string]interface{})
	for key, value := range params {
		arguments[key] = value
	}
	
	// Create the tool call request
	toolCallRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      p.toolName,
			Arguments: arguments,
		},
	}
	
	// Forward the call to the external MCP server
	response, err := p.client.CallTool(ctx, toolCallRequest)
	if err != nil {
		return &ToolResult{
			Error: err,
		}, fmt.Errorf("failed to call external tool %s: %w", p.toolName, err)
	}
	
	// Convert the response back to our ToolResult format
	var content string
	if len(response.Content) > 0 {
		// Handle different content types from the external server
		for _, item := range response.Content {
			switch contentItem := item.(type) {
			case mcp.TextContent:
				content += contentItem.Text
			case mcp.ImageContent:
				// Handle image content
				if jsonData, err := json.Marshal(contentItem); err == nil {
					content += string(jsonData)
				}
			default:
				// Handle other content types
				if jsonData, err := json.Marshal(item); err == nil {
					content += string(jsonData)
				}
			}
		}
	}
	
	// Check for errors in the response
	if response.IsError {
		return &ToolResult{
			Error: fmt.Errorf("external tool error: %s", content),
		}, nil
	}
	
	return &ToolResult{
		Content: content,
	}, nil
}

// Variable represents a configurable environment variable for MCP servers
type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Secret      bool   `json:"secret,omitempty"` // If true, value should be treated as sensitive
}

// MCPServerConfig represents configuration for connecting to an external MCP server
type MCPServerConfig struct {
	Name        string            `json:"name"`
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	BaseURL     string            `json:"baseUrl,omitempty"`
	Transport   string            `json:"transport"` // "stdio", "http", "sse"
	Disabled    bool              `json:"disabled,omitempty"`
	AutoApprove []string          `json:"autoApprove,omitempty"`
	Variables   []Variable        `json:"variables,omitempty"` // Framework-defined variables
}

// MCPProxy manages connections to external MCP servers and creates proxy tools
type MCPProxy struct {
	config MCPServerConfig
	client client.MCPClient
}

// NewMCPProxy creates a new MCP proxy for an external server
func NewMCPProxy(config MCPServerConfig) *MCPProxy {
	return &MCPProxy{
		config: config,
	}
}

// Connect establishes connection to the external MCP server
func (p *MCPProxy) Connect(ctx context.Context) error {
	if p.config.Disabled {
		return fmt.Errorf("MCP server %s is disabled", p.config.Name)
	}
	
	var mcpClient client.MCPClient
	var err error
	
	// Create client based on transport type
	switch p.config.Transport {
	case "stdio":
		if p.config.Command == "" {
			return fmt.Errorf("command is required for stdio transport")
		}
		// Convert env map to slice format
		var envSlice []string
		for key, value := range p.config.Env {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", key, value))
		}
		mcpClient, err = client.NewStdioMCPClient(p.config.Command, envSlice, p.config.Args...)
	case "http":
		if p.config.BaseURL == "" {
			return fmt.Errorf("baseUrl is required for http transport")
		}
		mcpClient, err = client.NewStreamableHttpClient(p.config.BaseURL, nil)
	case "sse":
		if p.config.BaseURL == "" {
			return fmt.Errorf("baseUrl is required for sse transport")
		}
		mcpClient, err = client.NewSSEMCPClient(p.config.BaseURL, nil)
	default:
		return fmt.Errorf("unsupported transport type: %s", p.config.Transport)
	}
	
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}
	
	// Note: stdio clients are automatically started by the constructor
	// HTTP and SSE clients may also be automatically started depending on implementation
	
	// Initialize the connection
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Ship MCP Proxy",
		Version: "1.0.0",
	}
	
	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		mcpClient.Close()
		return fmt.Errorf("failed to initialize MCP connection: %w", err)
	}
	
	p.client = mcpClient
	return nil
}

// DiscoverTools queries the external server for available tools and returns proxy tools
func (p *MCPProxy) DiscoverTools(ctx context.Context) ([]Tool, error) {
	if p.client == nil {
		return nil, fmt.Errorf("not connected to MCP server")
	}
	
	// List tools from the external server
	listToolsRequest := mcp.ListToolsRequest{}
	toolsResponse, err := p.client.ListTools(ctx, listToolsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools from external server: %w", err)
	}
	
	var proxyTools []Tool
	
	// Create proxy tools for each discovered tool
	for _, tool := range toolsResponse.Tools {
		// Convert MCP tool parameters to Ship parameters
		var shipParams []Parameter
		if tool.InputSchema.Properties != nil {
			for propName, propSchema := range tool.InputSchema.Properties {
				param := Parameter{
					Name:        propName,
					Type:        "string", // Default to string, could be enhanced to parse schema
					Description: fmt.Sprintf("Parameter for %s", propName),
					Required:    false, // Could parse required from schema
				}
				
				// Try to extract description from schema
				if schemaMap, ok := propSchema.(map[string]interface{}); ok {
					if desc, ok := schemaMap["description"].(string); ok {
						param.Description = desc
					}
					if typeStr, ok := schemaMap["type"].(string); ok {
						param.Type = typeStr
					}
				}
				
				shipParams = append(shipParams, param)
			}
		}
		
		// Create proxy tool with namespace prefix
		proxyToolName := fmt.Sprintf("%s_%s", p.config.Name, tool.Name)
		description := tool.Description
		if description == "" {
			description = fmt.Sprintf("Proxied tool %s from %s", tool.Name, p.config.Name)
		}
		
		proxyTool := NewProxyTool(
			proxyToolName,
			description,
			shipParams,
			p.client,
			tool.Name, // Original tool name for the external server
		)
		
		proxyTools = append(proxyTools, proxyTool)
	}
	
	return proxyTools, nil
}

// Close closes the connection to the external MCP server
func (p *MCPProxy) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}