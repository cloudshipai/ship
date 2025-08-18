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

	// Slack MCP Server
	"slack": {
		Name:      "slack",
		Command:   "sh",
		Args:      []string{"-c", "curl -sL https://github.com/korotovsky/slack-mcp-server/releases/latest/download/slack-mcp-server-darwin-amd64 -o /tmp/slack-mcp-server && chmod +x /tmp/slack-mcp-server && /tmp/slack-mcp-server"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "SLACK_MCP_XOXC_TOKEN",
				Description: "Slack browser token (xoxc-...) - required for stealth mode",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_XOXD_TOKEN",
				Description: "Slack browser cookie d (xoxd-...) - required for stealth mode",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_XOXP_TOKEN",
				Description: "User OAuth token (xoxp-...) - alternative to xoxc/xoxd for OAuth mode",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_PORT",
				Description: "Port for the MCP server to listen on (default: 13080)",
				Required:    false,
				Default:     "13080",
			},
			{
				Name:        "SLACK_MCP_HOST",
				Description: "Host for the MCP server to listen on (default: 127.0.0.1)",
				Required:    false,
				Default:     "127.0.0.1",
			},
			{
				Name:        "SLACK_MCP_SSE_API_KEY",
				Description: "Bearer token for SSE transport",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_PROXY",
				Description: "Proxy URL for outgoing requests",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_USER_AGENT",
				Description: "Custom User-Agent (for Enterprise Slack environments)",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_ADD_MESSAGE_TOOL",
				Description: "Enable message posting (true for all channels, comma-separated list for specific channels, or !channelID to exclude)",
				Required:    false,
				Default:     "true",
			},
			{
				Name:        "SLACK_MCP_ADD_MESSAGE_MARK",
				Description: "Automatically mark posted messages as read when enabled",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_ADD_MESSAGE_UNFURLING",
				Description: "Enable link unfurling (true for all domains, comma-separated list for specific domains)",
				Required:    false,
			},
			{
				Name:        "SLACK_MCP_USERS_CACHE",
				Description: "Path to users cache file (default: .users_cache.json)",
				Required:    false,
				Default:     ".users_cache.json",
			},
			{
				Name:        "SLACK_MCP_CHANNELS_CACHE",
				Description: "Path to channels cache file (default: .channels_cache_v2.json)",
				Required:    false,
				Default:     ".channels_cache_v2.json",
			},
			{
				Name:        "SLACK_MCP_LOG_LEVEL",
				Description: "Log level (debug, info, warn, error, panic, fatal) (default: info)",
				Required:    false,
				Default:     "info",
			},
		},
	},

	// GitHub MCP Server
	"github": {
		Name:      "github",
		Command:   "docker",
		Args:      []string{"run", "-i", "--rm", "ghcr.io/github/github-mcp-server"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "GITHUB_PERSONAL_ACCESS_TOKEN",
				Description: "GitHub Personal Access Token (required for authentication)",
				Required:    true,
				Secret:      true,
			},
			{
				Name:        "GITHUB_HOST",
				Description: "GitHub host (default: github.com, use for GitHub Enterprise)",
				Required:    false,
				Default:     "github.com",
			},
			{
				Name:        "GITHUB_TOOLSETS",
				Description: "Comma-separated list of toolsets to enable (e.g., repos,issues,pull_requests,actions,code_security)",
				Required:    false,
				Default:     "all",
			},
			{
				Name:        "GITHUB_READ_ONLY",
				Description: "Run in read-only mode (1 for true, 0 for false)",
				Required:    false,
				Default:     "0",
			},
			{
				Name:        "GITHUB_DYNAMIC_TOOLSETS",
				Description: "Enable dynamic toolset discovery (1 for true, 0 for false)",
				Required:    false,
				Default:     "0",
			},
		},
	},

	// DesktopCommander MCP Server
	"desktop-commander": {
		Name:      "desktop-commander",
		Command:   "npx",
		Args:      []string{"-y", "@wonderwhy-er/desktop-commander-mcp"},
		Transport: "stdio",
		Env:       map[string]string{},
		Variables: []ship.Variable{
			{
				Name:        "DESKTOP_COMMANDER_ROOT",
				Description: "Root directory for desktop operations (default: current working directory)",
				Required:    false,
				Default:     ".",
			},
			{
				Name:        "DESKTOP_COMMANDER_SAFE_MODE",
				Description: "Enable safe mode to prevent destructive operations (true/false)",
				Required:    false,
				Default:     "true",
			},
			{
				Name:        "DESKTOP_COMMANDER_LOG_LEVEL",
				Description: "Log level for desktop commander operations (debug, info, warn, error)",
				Required:    false,
				Default:     "info",
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
