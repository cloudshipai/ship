package mcp

import "github.com/cloudshipai/ship/pkg/ship"

// ExternalMCPServers contains the built-in external MCP server configurations
// These are third-party MCP servers that Ship can proxy to, not Ship's own tools
var ExternalMCPServers = map[string]ship.MCPServerConfig{
	// ModelContextProtocol Official Servers (Node.js based)
	"filesystem": {
		Name:      "filesystem",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "FILESYSTEM_ROOT",
				Description: "Root directory for filesystem operations (overrides /tmp default)",
				Required:    false,
				Default:     "/tmp",
			},
		},
	},
	"memory": {
		Name:      "memory",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "MEMORY_STORAGE_PATH",
				Description: "Path for persistent memory storage",
				Required:    false,
				Default:     "/tmp/mcp-memory",
			},
			{
				Name:        "MEMORY_MAX_SIZE",
				Description: "Maximum memory storage size (e.g., 100MB)",
				Required:    false,
				Default:     "50MB",
			},
		},
	},
	"brave-search": {
		Name:      "brave-search",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-brave-search"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "BRAVE_API_KEY",
				Description: "Brave Search API key for search functionality",
				Required:    true,
				Secret:      true,
			},
			{
				Name:        "BRAVE_SEARCH_COUNT",
				Description: "Number of search results to return (default: 10)",
				Required:    false,
				Default:     "10",
			},
		},
	},

	// AWS Labs Official MCP Servers (Python based, requires 'uv')
	"aws-core": {
		Name:      "aws-core",
		Command:   "uvx",
		Args:      []string{"awslabs.core-mcp-server@latest"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "AWS_PROFILE",
				Description: "AWS profile to use for authentication",
				Required:    false,
			},
			{
				Name:        "AWS_REGION",
				Description: "AWS region for operations",
				Required:    false,
				Default:     "us-east-1",
			},
			{
				Name:        "FASTMCP_LOG_LEVEL",
				Description: "Log level for the MCP server",
				Required:    false,
				Default:     "ERROR",
			},
		},
	},
	"aws-iam": {
		Name:      "aws-iam",
		Command:   "uvx",
		Args:      []string{"awslabs.iam-mcp-server@latest"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "AWS_PROFILE",
				Description: "AWS profile to use for authentication",
				Required:    false,
			},
			{
				Name:        "AWS_REGION",
				Description: "AWS region for operations",
				Required:    false,
				Default:     "us-east-1",
			},
			{
				Name:        "FASTMCP_LOG_LEVEL",
				Description: "Log level for the MCP server",
				Required:    false,
				Default:     "ERROR",
			},
		},
	},
	"aws-pricing": {
		Name:      "aws-pricing",
		Command:   "uvx",
		Args:      []string{"awslabs.pricing-mcp-server@latest"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "AWS_PROFILE",
				Description: "AWS profile to use for authentication",
				Required:    false,
			},
			{
				Name:        "AWS_REGION",
				Description: "AWS region for operations",
				Required:    false,
				Default:     "us-east-1",
			},
			{
				Name:        "FASTMCP_LOG_LEVEL",
				Description: "Log level for the MCP server",
				Required:    false,
				Default:     "ERROR",
			},
		},
	},
	"aws-eks": {
		Name:      "aws-eks",
		Command:   "uvx",
		Args:      []string{"awslabs.eks-mcp-server@latest"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "AWS_PROFILE",
				Description: "AWS profile to use for authentication",
				Required:    false,
			},
			{
				Name:        "AWS_REGION",
				Description: "AWS region for operations",
				Required:    false,
				Default:     "us-east-1",
			},
			{
				Name:        "FASTMCP_LOG_LEVEL",
				Description: "Log level for the MCP server",
				Required:    false,
				Default:     "ERROR",
			},
		},
	},
	"aws-ec2": {
		Name:      "aws-ec2",
		Command:   "uvx",
		Args:      []string{"awslabs.ec2-mcp-server@latest"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "AWS_PROFILE",
				Description: "AWS profile to use for authentication",
				Required:    false,
			},
			{
				Name:        "AWS_REGION",
				Description: "AWS region for operations",
				Required:    false,
				Default:     "us-east-1",
			},
			{
				Name:        "FASTMCP_LOG_LEVEL",
				Description: "Log level for the MCP server",
				Required:    false,
				Default:     "ERROR",
			},
		},
	},
	"aws-s3": {
		Name:      "aws-s3",
		Command:   "uvx",
		Args:      []string{"awslabs.s3-mcp-server@latest"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "AWS_PROFILE",
				Description: "AWS profile to use for authentication",
				Required:    false,
			},
			{
				Name:        "AWS_REGION",
				Description: "AWS region for operations",
				Required:    false,
				Default:     "us-east-1",
			},
			{
				Name:        "FASTMCP_LOG_LEVEL",
				Description: "Log level for the MCP server",
				Required:    false,
				Default:     "ERROR",
			},
		},
	},

	// Third-Party MCP Servers
	"steampipe": {
		Name:      "steampipe",
		Command:   "npx",
		Args:      []string{"-y", "@turbot/steampipe-mcp"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "STEAMPIPE_DATABASE_CONNECTIONS",
				Description: "Database connections configuration for Steampipe",
				Required:    false,
				Default:     "postgres://steampipe@localhost:9193/steampipe",
			},
		},
	},
}

// IsExternalMCPServer checks if the tool name matches an external MCP server
func IsExternalMCPServer(toolName string) bool {
	_, exists := ExternalMCPServers[toolName]
	return exists
}

// GetExternalMCPServer returns the configuration for an external MCP server
func GetExternalMCPServer(toolName string) (ship.MCPServerConfig, bool) {
	config, exists := ExternalMCPServers[toolName]
	return config, exists
}

// ListExternalMCPServers returns a list of all available external MCP server names
func ListExternalMCPServers() []string {
	servers := make([]string, 0, len(ExternalMCPServers))
	for name := range ExternalMCPServers {
		servers = append(servers, name)
	}
	return servers
}