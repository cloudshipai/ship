# External MCP Server Integration

Ship can proxy external MCP servers, allowing you to integrate existing MCP tools seamlessly into your Ship-based infrastructure. This document covers how to use external MCP servers with Ship's CLI and framework.

## Overview

Ship acts as both an MCP server (to AI assistants) and an MCP client (to external MCP servers), providing:

- **Dynamic Tool Discovery**: Tools from external MCP servers are discovered at runtime
- **Namespace Management**: External tools are prefixed with server names to avoid conflicts
- **Variable Support**: Pass environment variables and configuration to external servers
- **Seamless Integration**: External tools appear alongside Ship's built-in tools

## Available External MCP Servers

Ship comes with pre-configured external MCP servers:

### Filesystem Server
- **Name**: `filesystem`
- **Description**: File and directory operations
- **Tools**: `list_directory`, `read_file`, `write_file`, `create_directory`, `get_file_info`, etc.
- **Variables**:
  - `FILESYSTEM_ROOT` (optional): Root directory for operations (default: `/tmp`)

### Memory Server
- **Name**: `memory`
- **Description**: Persistent knowledge storage and retrieval
- **Tools**: `store_memory`, `search_memory`, `list_memories`, etc.
- **Variables**:
  - `MEMORY_STORAGE_PATH` (optional): Storage location (default: `/tmp/mcp-memory`)
  - `MEMORY_MAX_SIZE` (optional): Maximum storage size (default: `50MB`)

### Brave Search Server
- **Name**: `brave-search`
- **Description**: Web search capabilities using Brave Search API
- **Tools**: `brave_web_search`, `search_summarize`, etc.
- **Variables**:
  - `BRAVE_API_KEY` (required): Brave Search API key
  - `BRAVE_SEARCH_COUNT` (optional): Number of results (default: `10`)

### Grafana Server
- **Name**: `grafana`
- **Description**: Grafana monitoring and visualization platform integration
- **Docker Image**: `mcp/grafana:latest`
- **Tools**: Dashboard management, alerting, metrics querying, data source operations
- **Variables**:
  - `GRAFANA_URL` (required): Grafana server URL (e.g., `http://localhost:3000` or `https://myinstance.grafana.net`)
  - `GRAFANA_API_KEY` (required): Grafana service account token
  - `GRAFANA_USERNAME` (optional): Username for basic auth (alternative to API key)
  - `GRAFANA_PASSWORD` (optional): Password for basic auth (alternative to API key)

### Development & CI/CD Integration

#### Bitbucket Server
- **Name**: `bitbucket`
- **Description**: Atlassian Bitbucket Cloud integration for repository management and CI/CD
- **Tools**: Repository operations, pull request management, pipeline control
- **Variables**:
  - `BITBUCKET_USERNAME` (required): Atlassian Bitbucket username
  - `BITBUCKET_APP_PASSWORD` (required): Bitbucket app password (from account settings)
  - `BITBUCKET_WORKSPACE` (required): Bitbucket workspace name

#### Trello Server
- **Name**: `trello`  
- **Description**: Trello project management integration
- **Tools**: Board management, card operations, team collaboration
- **Variables**:
  - `TRELLO_API_KEY` (required): Trello API key (get from https://trello.com/app-key)
  - `TRELLO_TOKEN` (required): Trello API token (authorize from the API key page)

### Browser Automation & Testing

#### Playwright Server
- **Name**: `playwright`
- **Description**: Browser automation capabilities using Playwright
- **Tools**: Web automation, testing, screenshot capture, PDF generation
- **Variables**:
  - `PLAYWRIGHT_BROWSER` (optional): Browser to use (chromium, firefox, webkit) (default: chromium)
  - `PLAYWRIGHT_HEADLESS` (optional): Run browser in headless mode (true/false) (default: true)

### Database Integration

#### Supabase Server
- **Name**: `supabase`
- **Description**: Supabase backend-as-a-service integration (read-only mode)
- **Tools**: Database queries, authentication management, storage operations
- **Variables**:
  - `SUPABASE_ACCESS_TOKEN` (required): Supabase personal access token
  - `SUPABASE_PROJECT_REF` (required): Supabase project reference ID

#### PostgreSQL Server
- **Name**: `postgresql`
- **Description**: Direct PostgreSQL database operations and queries
- **Tools**: SQL queries, table management, data analysis
- **Variables**:
  - `POSTGRES_CONNECTION_STRING` (required): PostgreSQL connection string

## CLI Usage

### Basic Usage

Start an external MCP server proxy:

```bash
# Start filesystem server
ship mcp filesystem

# Start memory server
ship mcp memory

# Start Brave search (requires API key)
ship mcp brave-search --var BRAVE_API_KEY=your_api_key

# Start Grafana server (requires URL and API key)
ship mcp grafana --var GRAFANA_URL=http://localhost:3000 --var GRAFANA_API_KEY=your_service_account_token

# Start Bitbucket server (requires credentials and workspace)
ship mcp bitbucket --var BITBUCKET_USERNAME=your_username --var BITBUCKET_APP_PASSWORD=your_app_password --var BITBUCKET_WORKSPACE=your_workspace

# Start Trello server (requires API credentials)
ship mcp trello --var TRELLO_API_KEY=your_api_key --var TRELLO_TOKEN=your_token

# Start Playwright browser automation
ship mcp playwright --var PLAYWRIGHT_BROWSER=chromium --var PLAYWRIGHT_HEADLESS=false

# Start Supabase integration (requires project info)
ship mcp supabase --var SUPABASE_ACCESS_TOKEN=your_token --var SUPABASE_PROJECT_REF=your_project_ref

# Start PostgreSQL server (requires connection string)
ship mcp postgresql --var POSTGRES_CONNECTION_STRING=postgresql://user:password@localhost:5432/dbname
```

### Using Variables

Pass environment variables using the `--var` flag:

```bash
# Filesystem with custom root
ship mcp filesystem --var FILESYSTEM_ROOT=/workspace

# Memory with custom storage
ship mcp memory --var MEMORY_STORAGE_PATH=/data --var MEMORY_MAX_SIZE=100MB

# Brave search with custom result count
ship mcp brave-search --var BRAVE_API_KEY=your_key --var BRAVE_SEARCH_COUNT=20

# Grafana with API key authentication
ship mcp grafana --var GRAFANA_URL=https://myinstance.grafana.net --var GRAFANA_API_KEY=glsa_xyz123

# Grafana with username/password authentication
ship mcp grafana --var GRAFANA_URL=http://localhost:3000 --var GRAFANA_USERNAME=admin --var GRAFANA_PASSWORD=admin
```

### Variable Discovery

List available variables for external MCP servers:

```bash
ship modules info filesystem
ship modules info memory
ship modules info brave-search
ship modules info grafana
```

Example output:
```
External MCP Server: brave-search
Type: mcp-external
Description: Brave search MCP server for web search capabilities
Transport: stdio
Command: npx -y @modelcontextprotocol/server-brave-search
Source: hardcoded
Trusted: true

Variables:
  BRAVE_API_KEY (required) (secret)
    Brave Search API key for search functionality
  BRAVE_SEARCH_COUNT [default: 10]
    Number of search results to return (default: 10)

Usage:
  ship mcp brave-search    # Start MCP server proxy

Examples with variables:
  ship mcp brave-search --var BRAVE_API_KEY=<your_brave_api_key>
  ship mcp brave-search --var BRAVE_SEARCH_COUNT=20
```

## AI Assistant Integration

### Claude Code Configuration

Configure external MCP servers in Claude Code's MCP settings:

```json
{
  "mcpServers": {
    "ship-filesystem": {
      "command": "ship",
      "args": ["mcp", "filesystem"],
      "env": {
        "FILESYSTEM_ROOT": "/workspace"
      }
    },
    "ship-memory": {
      "command": "ship",
      "args": ["mcp", "memory", "--var", "MEMORY_STORAGE_PATH=/data"]
    },
    "ship-search": {
      "command": "ship",
      "args": ["mcp", "brave-search", "--var", "BRAVE_API_KEY=your_api_key"]
    }
  }
}
```

### Cursor Configuration

Configure in Cursor's MCP settings:

```json
{
  "mcpServers": {
    "ship-filesystem": {
      "command": "ship",
      "args": ["mcp", "filesystem"],
      "env": {
        "FILESYSTEM_ROOT": "/workspace"
      }
    }
  }
}
```

## Ship Framework Integration

Use external MCP servers in your Ship framework applications:

### Basic Integration

```go
package main

import (
    "context"
    "log"
    "github.com/cloudshipai/ship/pkg/ship"
)

func main() {
    ctx := context.Background()
    
    // Create Ship server
    server := ship.NewServer("my-app", "1.0.0")
    
    // Configure external MCP server
    filesystemConfig := ship.MCPServerConfig{
        Name:      "filesystem",
        Command:   "npx",
        Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/workspace"},
        Transport: "stdio",
        Env: map[string]string{
            "FILESYSTEM_ROOT": "/workspace",
        },
    }
    
    // Create and connect proxy
    proxy := ship.NewMCPProxy(filesystemConfig)
    if err := proxy.Connect(ctx); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer proxy.Close()
    
    // Discover tools dynamically
    tools, err := proxy.DiscoverTools(ctx)
    if err != nil {
        log.Fatalf("Failed to discover tools: %v", err)
    }
    
    // Add tools to server
    for _, tool := range tools {
        server.AddTool(tool)
    }
    
    // Start server
    mcpServer := server.Build()
    if err := mcpServer.ServeStdio(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### Multiple External Servers

```go
func addExternalMCPServers(server *ship.ServerBuilder, ctx context.Context) error {
    externalServers := []ship.MCPServerConfig{
        {
            Name:      "filesystem",
            Command:   "npx",
            Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/workspace"},
            Transport: "stdio",
            Env:       map[string]string{"FILESYSTEM_ROOT": "/workspace"},
        },
        {
            Name:      "memory",
            Command:   "npx",
            Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
            Transport: "stdio",
            Env:       map[string]string{"MEMORY_STORAGE_PATH": "/data"},
        },
    }
    
    for _, config := range externalServers {
        proxy := ship.NewMCPProxy(config)
        if err := proxy.Connect(ctx); err != nil {
            return fmt.Errorf("failed to connect to %s: %w", config.Name, err)
        }
        defer proxy.Close()
        
        tools, err := proxy.DiscoverTools(ctx)
        if err != nil {
            return fmt.Errorf("failed to discover tools from %s: %w", config.Name, err)
        }
        
        for _, tool := range tools {
            server.AddTool(tool)
        }
    }
    
    return nil
}
```

## Variable Configuration

### Framework-Defined Variables

Each external MCP server defines its own variables with:

- **Name**: Environment variable name
- **Description**: Human-readable description
- **Required**: Whether the variable is mandatory
- **Default**: Default value if not provided
- **Secret**: Whether the variable contains sensitive data

### Variable Validation

Ship validates variables before starting external MCP servers:

1. **Required Check**: Ensures all required variables are provided
2. **Default Assignment**: Sets default values for optional variables
3. **Environment Merge**: Combines user-provided and default values

### Custom Variable Support

You can pass any environment variable to external MCP servers:

```bash
# Custom environment variables
ship mcp filesystem --var CUSTOM_VAR=value --var DEBUG=true

# Override npm configuration
ship mcp memory --var NPM_CONFIG_REGISTRY=https://custom-registry.com
```

## Architecture

### MCP-to-MCP Proxying

```
AI Assistant  <--MCP-->  Ship CLI  <--MCP-->  External MCP Server
   (Claude)                (Proxy)              (filesystem)
```

1. **AI Assistant** calls Ship MCP server
2. **Ship** forwards calls to external MCP server
3. **External Server** processes the request
4. **Ship** returns response to AI assistant

### Tool Namespacing

External tools are prefixed with server names to avoid conflicts:

- `filesystem.list_directory`
- `filesystem.read_file`
- `memory.store_memory`
- `brave-search.web_search`

### Dynamic Discovery

Tools are discovered at runtime using the MCP `ListTools` method:

1. Ship connects to external MCP server
2. Calls `ListTools()` to get available tools
3. Creates proxy tools that forward calls
4. Registers tools with Ship's MCP server

## Troubleshooting

### Connection Issues

```bash
# Check if npm package is available
npx -y @modelcontextprotocol/server-filesystem --help

# Test connection manually
ship mcp filesystem --var DEBUG=true
```

### Variable Issues

```bash
# List required variables
ship modules info brave-search

# Check environment
echo $BRAVE_API_KEY

# Validate variable names
ship mcp brave-search --var BRAVE_API_KEY=test
```

### Tool Discovery Issues

External MCP servers must:
1. Implement MCP protocol correctly
2. Respond to `ListTools` requests
3. Support `CallTool` requests
4. Use stdio transport (for Ship's current implementation)

## Adding New External MCP Servers

To add new external MCP servers to Ship:

1. **Add Configuration**: Update `hardcodedMCPServers` in `internal/cli/mcp_cmd.go`
2. **Define Variables**: Specify required and optional variables
3. **Update Modules**: Add to module list in `internal/cli/modules_cmd.go`
4. **Test Integration**: Ensure tools are discovered correctly
5. **Update Documentation**: Add to this guide and README

Example configuration:
```go
"new-server": {
    Name:      "new-server",
    Command:   "npx",
    Args:      []string{"-y", "@example/mcp-server"},
    Transport: "stdio",
    Env:       map[string]string{},
    Variables: []ship.Variable{
        {
            Name:        "API_KEY",
            Description: "API key for external service",
            Required:    true,
            Secret:      true,
        },
    },
},
```

## Best Practices

1. **Use Variables**: Always use `--var` for configuration instead of hardcoding
2. **Handle Secrets**: Mark sensitive variables as `Secret: true`
3. **Provide Defaults**: Set sensible defaults for optional variables
4. **Test Thoroughly**: Verify tool discovery and execution
5. **Document Variables**: Provide clear descriptions for all variables
6. **Error Handling**: Handle connection failures gracefully
7. **Resource Cleanup**: Always close proxy connections

## Examples

See the [examples/](../examples/) directory for complete working examples of external MCP server integration.