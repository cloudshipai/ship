package ship

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Run("with name and version", func(t *testing.T) {
		builder := NewServer("test-server", "1.0.0")

		assert.NotNil(t, builder)
		assert.Equal(t, "test-server", builder.name)
		assert.Equal(t, "1.0.0", builder.version)
		assert.NotNil(t, builder.registry)
	})

	t.Run("with empty name", func(t *testing.T) {
		builder := NewServer("", "1.0.0")

		assert.Equal(t, "unnamed-server", builder.name)
		assert.Equal(t, "1.0.0", builder.version)
	})

	t.Run("with empty version", func(t *testing.T) {
		builder := NewServer("test-server", "")

		assert.Equal(t, "test-server", builder.name)
		assert.Equal(t, "1.0.0", builder.version)
	})
}

func TestServerBuilder_AddTool(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")

	t.Run("add valid tool", func(t *testing.T) {
		tool := createTestTool("test-tool", "Test tool")
		result := builder.AddTool(tool)

		assert.Equal(t, builder, result) // Should return self for chaining
		assert.Equal(t, 1, builder.registry.ToolCount())

		registeredTool, err := builder.registry.GetTool("test-tool")
		assert.NoError(t, err)
		assert.Equal(t, tool, registeredTool)
	})

	t.Run("add nil tool", func(t *testing.T) {
		originalCount := builder.registry.ToolCount()
		result := builder.AddTool(nil)

		assert.Equal(t, builder, result)
		assert.Equal(t, originalCount, builder.registry.ToolCount()) // Count should not change
	})
}

func TestServerBuilder_AddTools(t *testing.T) {
	t.Run("add multiple tools", func(t *testing.T) {
		builder := NewServer("test-server", "1.0.0")
		tool1 := createTestTool("tool1", "Tool 1")
		tool2 := createTestTool("tool2", "Tool 2")
		tool3 := createTestTool("tool3", "Tool 3")

		result := builder.AddTools(tool1, tool2, tool3)

		assert.Equal(t, builder, result)
		assert.Equal(t, 3, builder.registry.ToolCount())

		// Verify all tools are registered
		_, err1 := builder.registry.GetTool("tool1")
		_, err2 := builder.registry.GetTool("tool2")
		_, err3 := builder.registry.GetTool("tool3")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
	})

	t.Run("add tools with nil values", func(t *testing.T) {
		builder := NewServer("test-server-2", "1.0.0")
		tool1 := createTestTool("valid-tool", "Valid tool")

		result := builder.AddTools(tool1, nil, nil)

		assert.Equal(t, builder, result)
		assert.Equal(t, 1, builder.registry.ToolCount()) // Only valid tool should be added

		_, err := builder.registry.GetTool("valid-tool")
		assert.NoError(t, err)
	})
}

func TestServerBuilder_AddContainerTool(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")

	t.Run("add container tool", func(t *testing.T) {
		config := ContainerToolConfig{
			Description: "Test container tool",
			Image:       "test:latest",
			Parameters: []Parameter{
				{Name: "input", Type: "string", Required: true},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
				return &ToolResult{Content: "test output"}, nil
			},
		}

		result := builder.AddContainerTool("container-tool", config)

		assert.Equal(t, builder, result)
		assert.Equal(t, 1, builder.registry.ToolCount())

		tool, err := builder.registry.GetTool("container-tool")
		assert.NoError(t, err)
		assert.Equal(t, "container-tool", tool.Name())
		assert.Equal(t, "Test container tool", tool.Description())
	})
}

func TestServerBuilder_Build(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")
	tool := createTestTool("test-tool", "Test tool")
	builder.AddTool(tool)

	server := builder.Build()

	assert.NotNil(t, server)
	assert.Equal(t, "test-server", server.name)
	assert.Equal(t, "1.0.0", server.version)
	assert.Equal(t, builder.registry, server.registry)
	assert.Nil(t, server.engine) // Engine not initialized until Start()
}

func TestServerBuilder_ChainedOperations(t *testing.T) {
	tool1 := createTestTool("tool1", "Tool 1")
	tool2 := createTestTool("tool2", "Tool 2")

	server := NewServer("chained-server", "2.0.0").
		AddTool(tool1).
		AddTools(tool2).
		AddContainerTool("container-tool", ContainerToolConfig{
			Description: "Container tool",
			Image:       "test:latest",
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
				return &ToolResult{Content: "output"}, nil
			},
		}).
		Build()

	assert.NotNil(t, server)
	assert.Equal(t, "chained-server", server.name)
	assert.Equal(t, "2.0.0", server.version)
	assert.Equal(t, 3, server.registry.ToolCount())
}

func TestMCPServer_Start(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")
	tool := createTestTool("test-tool", "Test tool")
	builder.AddTool(tool)

	server := builder.Build()
	ctx := context.Background()

	t.Run("successful start", func(t *testing.T) {
		err := server.Start(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, server.engine)
		assert.NotNil(t, server.server)

		// Cleanup
		server.Close()
	})

	t.Run("start with cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Note: In real implementation with actual Dagger, this might fail
		// For our mock implementation, it should still work
		err := server.Start(cancelledCtx)
		assert.NoError(t, err) // Mock implementation doesn't check context cancellation

		server.Close()
	})
}

func TestMCPServer_Close(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")
	server := builder.Build()

	t.Run("close without starting", func(t *testing.T) {
		err := server.Close()
		assert.NoError(t, err)
	})

	t.Run("close after starting", func(t *testing.T) {
		ctx := context.Background()
		err := server.Start(ctx)
		require.NoError(t, err)

		err = server.Close()
		assert.NoError(t, err)
		assert.True(t, server.engine.IsClosed())
	})
}

func TestMCPServer_GetRegistry(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")
	tool := createTestTool("test-tool", "Test tool")
	builder.AddTool(tool)

	server := builder.Build()
	registry := server.GetRegistry()

	assert.NotNil(t, registry)
	assert.Equal(t, server.registry, registry)
	assert.Equal(t, 1, registry.ToolCount())
}

func TestMCPServer_GetEngine(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")
	server := builder.Build()

	t.Run("engine before start", func(t *testing.T) {
		engine := server.GetEngine()
		assert.Nil(t, engine)
	})

	t.Run("engine after start", func(t *testing.T) {
		ctx := context.Background()
		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		engine := server.GetEngine()
		assert.NotNil(t, engine)
		assert.Equal(t, server.engine, engine)
	})
}

func TestMCPServer_WriteOutput(t *testing.T) {
	builder := NewServer("test-server", "1.0.0")
	server := builder.Build()

	var buf bytes.Buffer
	err := server.WriteOutput(&buf, "test message")

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "[test-server]")
	assert.Contains(t, output, "test message")
}

func TestMCPServer_Logging(t *testing.T) {
	builder := NewServer("log-test-server", "1.0.0")
	server := builder.Build()

	// Note: These methods write to stderr, so we can't easily capture output in tests
	// But we can verify they don't panic
	assert.NotPanics(t, func() {
		server.LogInfo("test info message")
	})

	assert.NotPanics(t, func() {
		server.LogError("test error message")
	})
}

func TestMCPServer_ParameterConversion(t *testing.T) {
	// Test that framework parameters are correctly converted to MCP parameters
	builder := NewServer("param-test", "1.0.0")

	// Create tool with various parameter types
	config := ContainerToolConfig{
		Description: "Parameter test tool",
		Image:       "test:latest",
		Parameters: []Parameter{
			{Name: "string_param", Type: "string", Description: "String parameter", Required: true},
			{Name: "bool_param", Type: "boolean", Description: "Boolean parameter"},
			{Name: "number_param", Type: "number", Description: "Number parameter"},
			{Name: "enum_param", Type: "string", Description: "Enum parameter", Enum: []string{"option1", "option2"}},
		},
		Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			// Verify parameter types
			stringVal, stringOk := params["string_param"].(string)
			boolVal, boolOk := params["bool_param"].(bool)
			numberVal, numberOk := params["number_param"].(float64)
			enumVal, enumOk := params["enum_param"].(string)

			result := fmt.Sprintf("string:%s(%t) bool:%t(%t) number:%f(%t) enum:%s(%t)",
				stringVal, stringOk, boolVal, boolOk, numberVal, numberOk, enumVal, enumOk)

			return &ToolResult{Content: result}, nil
		},
	}

	builder.AddContainerTool("param-tool", config)
	server := builder.Build()

	ctx := context.Background()
	err := server.Start(ctx)
	require.NoError(t, err)
	defer server.Close()

	// Verify the tool was registered
	assert.NotNil(t, server.server)

	// The actual MCP tool registration is tested indirectly through the Start method
	// More detailed testing would require mocking the MCP server
}

// createTestTool helper is defined in registry_test.go
