# Basic Custom Server Example

This example demonstrates how to create a custom MCP server using the Ship Framework with two simple tools.

## Tools Included

### 1. Echo Tool
- **Purpose**: Simple text processing example
- **Parameters**: 
  - `message` (required): Text to echo back
  - `uppercase` (optional): Whether to convert to uppercase
- **Usage**: Demonstrates basic parameter handling and string processing

### 2. File List Tool  
- **Purpose**: Containerized file system operation example
- **Parameters**:
  - `path` (optional): Directory path to list (defaults to current)
  - `show_hidden` (optional): Whether to show hidden files
- **Usage**: Shows how to use Dagger containers for file system operations

## Running the Example

```bash
# Install dependencies
go mod tidy

# Start the MCP server
go run main.go
```

The server will start and wait for MCP protocol messages on stdin/stdout.

## Testing with Claude Desktop

Add this to your Claude Desktop MCP configuration:

```json
{
  "mcpServers": {
    "basic-custom-server": {
      "command": "go",
      "args": ["run", "main.go"],
      "cwd": "/path/to/ship/examples/ship-framework/basic-custom-server"
    }
  }
}
```

Then you can ask Claude to:
- "Echo 'Hello World' in uppercase"
- "List files in the current directory"
- "Show hidden files in /tmp"

## Key Learning Points

1. **Custom Tool Implementation**: Both tools implement the `ship.Tool` interface
2. **Parameter Handling**: Shows string and boolean parameter types
3. **Container Integration**: File List tool uses Alpine Linux container via Dagger
4. **Error Handling**: Proper error handling for missing parameters and container failures
5. **Metadata**: Tools can return structured metadata alongside results

## Architecture

```
Ship Framework
├── Tool Interface (ship.Tool)
├── Server Builder (ship.NewServer)
├── Dagger Engine Integration
└── MCP Protocol Handling
```

This example shows the **actual working** Ship Framework API, not the incorrect examples in the README.