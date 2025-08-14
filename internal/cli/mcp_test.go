package cli

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServerIntegration(t *testing.T) {
	// Skip if in CI or if ship binary doesn't exist
	if os.Getenv("CI") != "" {
		t.Skip("Skipping MCP integration test in CI")
	}

	// Check if ship binary exists (when running from project root)
	shipPath := "./ship"
	if _, err := os.Stat(shipPath); err != nil {
		// Try absolute path
		shipPath = "/home/epuerta/projects/ship/ship"
	}
	if _, err := os.Stat(shipPath); err != nil {
		t.Skip("Ship binary not found, run 'make build' first")
	}
	t.Logf("Using ship binary at: %s", shipPath)

	t.Run("MCP Server Lists Tools", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create MCP client with stdio transport
		clientTransport := transport.NewStdio(shipPath, nil, "mcp", "lint")
		mcpClient := client.NewClient(clientTransport)

		// Start the client
		err := mcpClient.Start(ctx)
		require.NoError(t, err)
		defer mcpClient.Close()

		// Initialize MCP session
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "ship-test-client",
			Version: "1.0.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		initResult, err := mcpClient.Initialize(ctx, initRequest)
		require.NoError(t, err)
		assert.Equal(t, mcp.LATEST_PROTOCOL_VERSION, initResult.ProtocolVersion)

		// List available tools
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
		require.NoError(t, err)
		
		// Verify we get the lint tool
		require.Len(t, toolsResult.Tools, 1)
		tool := toolsResult.Tools[0]
		assert.Equal(t, "lint", tool.Name)
		assert.Contains(t, tool.Description, "TFLint")
	})

	t.Run("MCP Server All Tools", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create MCP client for all tools
		clientTransport := transport.NewStdio(shipPath, nil, "mcp", "all")
		mcpClient := client.NewClient(clientTransport)

		// Start the client
		err := mcpClient.Start(ctx)
		require.NoError(t, err)
		defer mcpClient.Close()

		// Initialize MCP session
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "ship-test-client",
			Version: "1.0.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = mcpClient.Initialize(ctx, initRequest)
		require.NoError(t, err)

		// List all tools
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
		require.NoError(t, err)
		
		// Verify we get all 6 tools (lint, checkov, trivy, cost, docs, diagram)
		assert.Len(t, toolsResult.Tools, 6)
		
		// Check that tool names are clean (no prefixes)
		toolNames := make([]string, len(toolsResult.Tools))
		for i, tool := range toolsResult.Tools {
			toolNames[i] = tool.Name
		}
		
		expectedTools := []string{"checkov", "cost", "diagram", "docs", "lint", "trivy"}
		for _, expectedTool := range expectedTools {
			assert.Contains(t, toolNames, expectedTool, "Missing expected tool: %s", expectedTool)
		}

		// Verify no tools have old prefixed names
		for _, tool := range toolsResult.Tools {
			assert.NotContains(t, tool.Name, "terraform_")
			assert.NotContains(t, tool.Name, "terraform-")
		}

		t.Logf("All tools found: %v", toolNames)
	})

	t.Run("MCP Server Resources", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create MCP client for all tools (which includes resources)
		clientTransport := transport.NewStdio(shipPath, nil, "mcp", "all")
		mcpClient := client.NewClient(clientTransport)

		// Start the client
		err := mcpClient.Start(ctx)
		require.NoError(t, err)
		defer mcpClient.Close()

		// Initialize MCP session
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "ship-test-client",
			Version: "1.0.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		initResult, err := mcpClient.Initialize(ctx, initRequest)
		require.NoError(t, err)

		// Check if server supports resources
		if initResult.Capabilities.Resources != nil {
			// List resources
			resourcesRequest := mcp.ListResourcesRequest{}
			resourcesResult, err := mcpClient.ListResources(ctx, resourcesRequest)
			require.NoError(t, err)
			
			// Should have help and tools resources
			assert.Len(t, resourcesResult.Resources, 2)
			
			resourceNames := make([]string, len(resourcesResult.Resources))
			for i, resource := range resourcesResult.Resources {
				resourceNames[i] = resource.URI
			}
			
			assert.Contains(t, resourceNames, "ship://help")
			assert.Contains(t, resourceNames, "ship://tools")
			
			t.Logf("Resources found: %v", resourceNames)
		} else {
			t.Log("Server does not support resources")
		}
	})

	t.Run("MCP Tool Call", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Create a temporary directory with some Terraform files for testing
		tempDir := t.TempDir()
		terraformFile := `
resource "aws_instance" "example" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t2.micro"
}
`
		err := os.WriteFile(tempDir+"/main.tf", []byte(terraformFile), 0644)
		require.NoError(t, err)

		// Create MCP client for lint tool
		clientTransport := transport.NewStdio(shipPath, nil, "mcp", "lint")
		mcpClient := client.NewClient(clientTransport)

		// Start the client
		err = mcpClient.Start(ctx)
		require.NoError(t, err)
		defer mcpClient.Close()

		// Initialize MCP session
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "ship-test-client",
			Version: "1.0.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = mcpClient.Initialize(ctx, initRequest)
		require.NoError(t, err)

		// Call the lint tool
		callRequest := mcp.CallToolRequest{}
		callRequest.Params.Name = "lint"
		callRequest.Params.Arguments = map[string]interface{}{
			"directory": tempDir,
			"format":    "json",
		}

		callResult, err := mcpClient.CallTool(ctx, callRequest)
		require.NoError(t, err)
		
		// Verify we got a result
		assert.NotEmpty(t, callResult.Content)
		
		// The result should contain some text content
		if len(callResult.Content) > 0 {
			if textContent, ok := callResult.Content[0].(mcp.TextContent); ok {
				assert.NotEmpty(t, textContent.Text)
				t.Logf("Lint result: %s", textContent.Text)
			}
		}
	})
}