package ship

import (
	"context"
	"testing"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

// MockTool for testing
type MockTool struct {
	name        string
	description string
	parameters  []Parameter
}

func (m *MockTool) Name() string                 { return m.name }
func (m *MockTool) Description() string          { return m.description }
func (m *MockTool) Parameters() []Parameter      { return m.parameters }
func (m *MockTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
	return &ToolResult{Content: "mock result"}, nil
}

func TestMCPAdapter(t *testing.T) {
	t.Run("create new adapter", func(t *testing.T) {
		adapter := NewMCPAdapter()

		assert.NotNil(t, adapter)
		assert.NotNil(t, adapter.GetRegistry())
		assert.Nil(t, adapter.GetEngine()) // Engine not initialized yet
	})

	t.Run("add tools to adapter", func(t *testing.T) {
		adapter := NewMCPAdapter()

		// Create a mock tool
		mockTool := &MockTool{
			name:        "test-tool",
			description: "A test tool",
			parameters:  []Parameter{{Name: "input", Type: "string", Required: true}},
		}

		adapter.AddTool(mockTool)

		registry := adapter.GetRegistry()
		tools := registry.GetAllTools()

		assert.Len(t, tools, 1)
		assert.Contains(t, tools, "test-tool")
	})

	t.Run("attach to existing mcp server", func(t *testing.T) {
		ctx := context.Background()

		// Create MCP server (this would normally be user's existing server)
		mcpServer := server.NewMCPServer("test-server", "1.0.0")

		// Create adapter with mock tool
		adapter := NewMCPAdapter()
		mockTool := &MockTool{
			name:        "container-tool",
			description: "A containerized tool",
			parameters:  []Parameter{{Name: "message", Type: "string", Required: true}},
		}
		adapter.AddTool(mockTool)

		// This would normally initialize Dagger, but we'll skip for this test
		// In a real test, we'd need a proper Dagger engine or mock
		t.Skip("Requires Dagger engine for full integration test")

		// Attach Ship tools to the existing MCP server
		err := adapter.AttachToServer(ctx, mcpServer)
		assert.NoError(t, err)

		// Verify the tool was registered (this would require inspecting the server state)
		// In practice, you'd test by actually calling the tool via MCP protocol
	})

	t.Run("import registry", func(t *testing.T) {
		adapter := NewMCPAdapter()

		// Create a registry with tools
		registry := NewRegistry()
		mockTool1 := &MockTool{name: "tool1", description: "Tool 1"}
		mockTool2 := &MockTool{name: "tool2", description: "Tool 2"}

		registry.RegisterTool(mockTool1)
		registry.RegisterTool(mockTool2)

		// Import the registry
		adapter.ImportRegistry(registry)

		tools := adapter.GetRegistry().GetAllTools()
		assert.Len(t, tools, 2)
		assert.Contains(t, tools, "tool1")
		assert.Contains(t, tools, "tool2")
	})

	t.Run("add container tool", func(t *testing.T) {
		adapter := NewMCPAdapter()

		config := ContainerToolConfig{
			Description: "Test container tool",
			Image:       "alpine:latest",
			Parameters: []Parameter{
				{Name: "command", Type: "string", Description: "Command to run", Required: true},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
				return &ToolResult{Content: "test output"}, nil
			},
		}

		adapter.AddContainerTool("alpine-tool", config)

		tools := adapter.GetRegistry().GetAllTools()
		assert.Len(t, tools, 1)
		assert.Contains(t, tools, "alpine-tool")

		tool := tools["alpine-tool"]
		assert.Equal(t, "Test container tool", tool.Description())
		assert.Len(t, tool.Parameters(), 1)
		assert.Equal(t, "command", tool.Parameters()[0].Name)
	})
}

func TestToolRouter(t *testing.T) {
	t.Run("create new router", func(t *testing.T) {
		router := NewToolRouter()

		assert.NotNil(t, router)
		assert.NotNil(t, router.adapters)
		assert.Len(t, router.adapters, 0)
	})

	t.Run("add routes", func(t *testing.T) {
		router := NewToolRouter()
		adapter1 := NewMCPAdapter()
		adapter2 := NewMCPAdapter()

		router.AddRoute("security", adapter1).
			AddRoute("infrastructure", adapter2)

		assert.Len(t, router.adapters, 2)
		assert.Contains(t, router.adapters, "security")
		assert.Contains(t, router.adapters, "infrastructure")
		assert.Equal(t, adapter1, router.adapters["security"])
		assert.Equal(t, adapter2, router.adapters["infrastructure"])
	})
}

// Integration test pattern for users wanting to test their setup
func ExampleMCPAdapter_attachToServer() {
	ctx := context.Background()

	// User's existing MCP server
	mcpServer := server.NewMCPServer("my-app", "1.0.0")

	// Create Ship adapter with infrastructure tools
	adapter := NewMCPAdapter()

	// In real usage, you'd add actual tools:
	// adapter.AddTool(tools.NewTFLintTool())
	// adapter.AddTool(tools.NewCheckovTool())

	// Attach Ship tools to existing server
	if err := adapter.AttachToServer(ctx, mcpServer); err != nil {
		panic(err)
	}

	// Clean up
	defer adapter.Close()

	// Your server now has both your tools and Ship's containerized tools
	// server.ServeStdio(mcpServer)
}

// Example showing the "bring your own mcp-go" pattern
func ExamplePattern_bringYourOwnMCP() {
	ctx := context.Background()

	// Your existing mcp-go server setup
	mcpServer := server.NewMCPServer("existing-app", "2.0.0")

	// Your existing tools would be added here
	// mcpServer.AddTool(myTool, myHandler)

	// Add Ship's containerized infrastructure tools
	shipAdapter := NewMCPAdapter()

	// In real code, use actual tools:
	// shipAdapter.AddTool(tools.NewTFLintTool())
	// shipAdapter.AddContainerTool("custom", myConfig)

	// Attach to your existing server
	shipAdapter.AttachToServer(ctx, mcpServer)
	defer shipAdapter.Close()

	// Now your server has both your tools AND Ship's containerized tools
	// server.ServeStdio(mcpServer)
}
