# Container Tools Example

This example demonstrates how to build practical containerized tools using the Ship Framework. Each tool runs in an isolated Docker container, providing consistent behavior across environments.

## Tools Included

### 1. Terraform Validate Tool
- **Container**: `hashicorp/terraform:latest`
- **Purpose**: Validates Terraform configuration syntax
- **Parameters**: 
  - `directory` (optional): Terraform project directory
  - `json_output` (optional): Return JSON formatted results
- **Features**: Automatic `terraform init` before validation

### 2. Dockerfile Lint Tool  
- **Container**: `hadolint/hadolint:latest`
- **Purpose**: Lints Dockerfiles for best practices and security
- **Parameters**:
  - `dockerfile_path` (optional): Path to Dockerfile
  - `format` (optional): Output format (json, tty, gcc, etc.)
- **Features**: Multiple output formats for different CI/CD integrations

### 3. YAML Lint Tool
- **Container**: `cytopia/yamllint:latest`  
- **Purpose**: Validates YAML syntax and formatting
- **Parameters**:
  - `files` (optional): Space-separated list of files/directories
  - `strict` (optional): Treat warnings as errors
- **Features**: Configurable strictness levels

## Running the Example

```bash
# Install dependencies
go mod tidy

# Start the MCP server
go run main.go
```

## Testing with Sample Files

Create test files in the same directory:

**Dockerfile** (with intentional issues):
```dockerfile
from ubuntu:18.04
RUN apt-get update
RUN apt-get install -y curl
ADD . /app
CMD /app/start.sh
```

**main.tf** (simple Terraform):
```hcl
resource "aws_s3_bucket" "example" {
  bucket = "my-tf-test-bucket-12345"
}

output "bucket_name" {
  value = aws_s3_bucket.example.bucket
}
```

**config.yaml** (with YAML issues):
```yaml
name: test
version: 1.0
environment:
	production: true  # Tab instead of spaces
```

## Claude Desktop Configuration

```json
{
  "mcpServers": {
    "container-tools": {
      "command": "go",
      "args": ["run", "main.go"],
      "cwd": "/path/to/ship/examples/ship-framework/container-tools",
      "env": {
        "DOCKER_HOST": "unix:///var/run/docker.sock"
      }
    }
  }
}
```

## Example Interactions

Ask Claude to:
- "Validate the Terraform configuration in the current directory"
- "Lint the Dockerfile and show results in JSON format"  
- "Check all YAML files with strict mode enabled"
- "Validate Terraform in the ./infrastructure directory with JSON output"

## Key Architecture Patterns

### 1. Container Isolation
Each tool runs in its own container, ensuring:
- Consistent tool versions
- No local dependency conflicts  
- Secure execution environment

### 2. Directory Mounting
```go
container := client.Container().
    From("tool-image:latest").
    WithDirectory("/workspace", host.Directory(directory)).
    WithWorkdir("/workspace")
```

### 3. Error Handling Strategy
```go
stdout, err := container.Stdout(ctx)
if err != nil {
    stderr, _ := container.Stderr(ctx)
    // Handle both stdout and stderr appropriately
}
```

### 4. Flexible Parameter Handling
Tools support optional parameters with sensible defaults and enum validation for constrained values.

## Prerequisites

- **Docker**: Required for Dagger container execution
- **Go 1.23+**: For building and running the server
- **Network Access**: To pull container images

## Production Considerations

1. **Image Caching**: Dagger caches container images for faster execution
2. **Security**: Containers run with limited permissions
3. **Resource Management**: Consider memory/CPU limits for production use
4. **Error Handling**: Distinguish between tool errors and execution failures

This example shows **realistic containerized tool patterns** that you can adapt for your specific use cases.