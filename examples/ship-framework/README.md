# Ship Framework Examples

This directory contains working examples of how to use the Ship Framework to build custom MCP servers.

**IMPORTANT**: The Ship Framework provides the infrastructure (server, tool interfaces, registry), but pre-built tools like TFLint, Checkov, etc. are only available through the Ship CLI, not as importable packages.

## What's Available for External Use

The Ship Framework public API (`github.com/cloudshipai/ship/pkg/ship`) provides:
- **Server Builder**: Fluent API for building MCP servers
- **Tool Interface**: For creating custom containerized tools
- **Registry System**: For managing tools, prompts, and resources
- **Dagger Integration**: For containerized tool execution

## Example Projects

1. **basic-custom-server**: Simple MCP server with custom tools
2. **container-tools**: Examples of containerized tool integration
3. **terraform-proxy**: How to proxy existing Ship tools (advanced)

## Quick Start

```bash
# Navigate to an example
cd basic-custom-server

# Run the example
go mod tidy
go run main.go
```

Each example includes a complete `go.mod` file and working code that demonstrates real Ship Framework usage patterns.