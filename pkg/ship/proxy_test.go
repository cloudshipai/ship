package ship

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPProxyConnection(t *testing.T) {
	// Skip this test if we're in CI or don't have npm
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("FilesystemServer", func(t *testing.T) {
		// Create a filesystem server configuration
		config := MCPServerConfig{
			Name:      "filesystem-test",
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
			Transport: "stdio",
			Env:       map[string]string{},
		}

		// Create proxy
		proxy := NewMCPProxy(config)
		require.NotNil(t, proxy)

		// Set a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Connect to the server
		err := proxy.Connect(ctx)
		require.NoError(t, err, "Failed to connect to filesystem MCP server")
		defer proxy.Close()

		// Discover tools
		tools, err := proxy.DiscoverTools(ctx)
		require.NoError(t, err, "Failed to discover tools")
		require.Greater(t, len(tools), 0, "No tools discovered")

		// Verify some expected tools exist
		toolNames := make([]string, len(tools))
		for i, tool := range tools {
			toolNames[i] = tool.Name()
		}

		// Check for expected filesystem tools
		expectedTools := []string{
			"filesystem-test.read_file",
			"filesystem-test.write_file", 
			"filesystem-test.list_directory",
			"filesystem-test.create_directory",
		}

		for _, expectedTool := range expectedTools {
			assert.Contains(t, toolNames, expectedTool, "Expected tool %s not found", expectedTool)
		}

		t.Logf("Discovered %d tools: %v", len(tools), toolNames)
	})

	t.Run("MemoryServer", func(t *testing.T) {
		// Create a memory server configuration
		config := MCPServerConfig{
			Name:      "memory-test",
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
			Transport: "stdio",
			Env:       map[string]string{},
		}

		// Create proxy
		proxy := NewMCPProxy(config)
		require.NotNil(t, proxy)

		// Set a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Connect to the server
		err := proxy.Connect(ctx)
		require.NoError(t, err, "Failed to connect to memory MCP server")
		defer proxy.Close()

		// Discover tools
		tools, err := proxy.DiscoverTools(ctx)
		require.NoError(t, err, "Failed to discover tools")
		require.Greater(t, len(tools), 0, "No tools discovered")

		// Verify tools are namespaced
		for _, tool := range tools {
			assert.Contains(t, tool.Name(), "memory-test.", "Tool name should be namespaced with server name")
		}

		toolNames := make([]string, len(tools))
		for i, tool := range tools {
			toolNames[i] = tool.Name()
		}

		t.Logf("Discovered %d memory tools: %v", len(tools), toolNames)
	})
}

func TestMCPServerProxyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Note: ProxyServerBuilder functionality has been moved to CLI level
	// External MCP server integration is handled in internal/cli/ not pkg/ship/
}

func TestHardcodedMCPServers(t *testing.T) {
	// Test the hardcoded configurations we have in the CLI
	hardcodedConfigs := map[string]MCPServerConfig{
		"filesystem": {
			Name:      "filesystem",
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
			Transport: "stdio",
			Env:       map[string]string{},
		},
		"memory": {
			Name:      "memory",
			Command:   "npx",
			Args:      []string{"-y", "@modelcontextprotocol/server-memory"},
			Transport: "stdio",
			Env:       map[string]string{},
		},
	}

	for name, config := range hardcodedConfigs {
		t.Run(name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("Skipping integration test in short mode")
			}

			proxy := NewMCPProxy(config)
			require.NotNil(t, proxy)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			err := proxy.Connect(ctx)
			if err != nil {
				t.Logf("Failed to connect to %s server (this may be expected if package isn't available): %v", name, err)
				t.Skip("MCP server package not available")
				return
			}
			defer proxy.Close()

			tools, err := proxy.DiscoverTools(ctx)
			require.NoError(t, err)
			assert.Greater(t, len(tools), 0, "Should discover at least one tool")

			t.Logf("Server %s has %d tools", name, len(tools))
		})
	}
}