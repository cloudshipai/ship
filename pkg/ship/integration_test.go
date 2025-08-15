package ship

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFullFrameworkIntegration tests the complete MCP server lifecycle
func TestFullFrameworkIntegration(t *testing.T) {
	ctx := context.Background()

	// Create a test tool that simulates real container execution
	testTool := NewContainerTool("integration-test-tool", ContainerToolConfig{
		Description: "Integration test tool for framework testing",
		Image:       "alpine:latest",
		Parameters: []Parameter{
			{
				Name:        "operation",
				Type:        "string",
				Description: "Operation to perform",
				Required:    true,
				Enum:        []string{"echo", "env", "pwd"},
			},
			{
				Name:        "value",
				Type:        "string",
				Description: "Value for the operation",
				Required:    false,
			},
		},
		Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			operation := params["operation"].(string)
			value := ""
			if v, ok := params["value"].(string); ok {
				value = v
			}

			var output string
			var err error

			switch operation {
			case "echo":
				if value == "" {
					value = "hello world"
				}
				output, err = engine.Container().
					From("alpine:latest").
					WithExec([]string{"echo", value}).
					Stdout(ctx)
			case "env":
				output, err = engine.Container().
					From("alpine:latest").
					WithEnvVariable("TEST_VAR", "test_value").
					WithExec([]string{"env"}).
					Stdout(ctx)
			case "pwd":
				output, err = engine.Container().
					From("alpine:latest").
					WithWorkdir("/tmp").
					WithExec([]string{"pwd"}).
					Stdout(ctx)
			default:
				return &ToolResult{
					Error: fmt.Errorf("unknown operation: %s", operation),
				}, fmt.Errorf("unknown operation: %s", operation)
			}

			if err != nil {
				return &ToolResult{
					Content: "",
					Error:   err,
				}, err
			}

			return &ToolResult{
				Content: output,
				Metadata: map[string]interface{}{
					"operation": operation,
					"value":     value,
					"tool":      "integration-test-tool",
				},
			}, nil
		},
	})

	// Create server with test tool
	server := NewServer("integration-test-server", "1.0.0").
		AddTool(testTool).
		Build()

	require.NotNil(t, server)

	// Test server lifecycle
	t.Run("ServerLifecycle", func(t *testing.T) {
		// Start server
		err := server.Start(ctx)
		require.NoError(t, err)

		// Verify server state
		assert.NotNil(t, server.GetEngine())
		assert.Equal(t, 1, server.GetRegistry().ToolCount())

		// Close server
		err = server.Close()
		assert.NoError(t, err)
	})

	t.Run("ToolExecution", func(t *testing.T) {
		// Start server for tool execution tests
		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		// Get the test tool
		tool, err := server.GetRegistry().GetTool("integration-test-tool")
		require.NoError(t, err)
		assert.Equal(t, "integration-test-tool", tool.Name())

		// Test echo operation
		result, err := tool.Execute(ctx, map[string]interface{}{
			"operation": "echo",
			"value":     "integration test",
		}, server.GetEngine())

		require.NoError(t, err)
		assert.Contains(t, result.Content, "integration test")
		assert.Equal(t, "echo", result.Metadata["operation"])

		// Test environment operation
		result, err = tool.Execute(ctx, map[string]interface{}{
			"operation": "env",
		}, server.GetEngine())

		require.NoError(t, err)
		assert.Contains(t, result.Content, "TEST_VAR=test_value")
		assert.Equal(t, "env", result.Metadata["operation"])

		// Test pwd operation
		result, err = tool.Execute(ctx, map[string]interface{}{
			"operation": "pwd",
		}, server.GetEngine())

		require.NoError(t, err)
		assert.Contains(t, result.Content, "/tmp")
		assert.Equal(t, "pwd", result.Metadata["operation"])
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		tool, err := server.GetRegistry().GetTool("integration-test-tool")
		require.NoError(t, err)

		// Test invalid operation
		result, err := tool.Execute(ctx, map[string]interface{}{
			"operation": "invalid",
		}, server.GetEngine())

		assert.Error(t, err)
		assert.NotNil(t, result.Error)
		assert.Contains(t, err.Error(), "unknown operation")
	})

	t.Run("ConcurrentExecution", func(t *testing.T) {
		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		tool, err := server.GetRegistry().GetTool("integration-test-tool")
		require.NoError(t, err)

		// Execute multiple operations concurrently
		numGoroutines := 5
		results := make(chan *ToolResult, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				result, err := tool.Execute(ctx, map[string]interface{}{
					"operation": "echo",
					"value":     fmt.Sprintf("concurrent test %d", index),
				}, server.GetEngine())

				results <- result
				errors <- err
			}(i)
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			select {
			case result := <-results:
				assert.NotNil(t, result)
				assert.Contains(t, result.Content, "concurrent test")
			case <-time.After(30 * time.Second):
				t.Fatal("Timeout waiting for concurrent execution")
			}

			select {
			case err := <-errors:
				assert.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout waiting for error result")
			}
		}
	})
}

// TestMultipleToolIntegration tests framework with multiple tools
func TestMultipleToolIntegration(t *testing.T) {
	ctx := context.Background()

	// Create multiple test tools
	echoTool := NewContainerTool("echo-tool", ContainerToolConfig{
		Description: "Simple echo tool",
		Image:       "alpine:latest",
		Parameters: []Parameter{
			{
				Name:        "message",
				Type:        "string",
				Description: "Message to echo",
				Required:    true,
			},
		},
		Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			message := params["message"].(string)
			output, err := engine.Container().
				From("alpine:latest").
				WithExec([]string{"echo", message}).
				Stdout(ctx)

			if err != nil {
				return &ToolResult{Error: err}, err
			}

			return &ToolResult{
				Content: output,
				Metadata: map[string]interface{}{
					"tool":    "echo-tool",
					"message": message,
				},
			}, nil
		},
	})

	lsTool := NewContainerTool("ls-tool", ContainerToolConfig{
		Description: "List directory contents",
		Image:       "alpine:latest",
		Parameters: []Parameter{
			{
				Name:        "directory",
				Type:        "string",
				Description: "Directory to list",
				Required:    false,
			},
		},
		Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			directory := "/"
			if d, ok := params["directory"].(string); ok && d != "" {
				directory = d
			}

			output, err := engine.Container().
				From("alpine:latest").
				WithExec([]string{"ls", "-la", directory}).
				Stdout(ctx)

			if err != nil {
				return &ToolResult{Error: err}, err
			}

			return &ToolResult{
				Content: output,
				Metadata: map[string]interface{}{
					"tool":      "ls-tool",
					"directory": directory,
				},
			}, nil
		},
	})

	// Create server with multiple tools
	server := NewServer("multi-tool-server", "1.0.0").
		AddTool(echoTool).
		AddTool(lsTool).
		Build()

	require.NotNil(t, server)

	t.Run("MultipleToolsRegistration", func(t *testing.T) {
		assert.Equal(t, 2, server.GetRegistry().ToolCount())

		tools := server.GetRegistry().ListTools()
		assert.Contains(t, tools, "echo-tool")
		assert.Contains(t, tools, "ls-tool")
	})

	t.Run("ExecuteAllTools", func(t *testing.T) {
		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		// Test echo tool
		echoTool, err := server.GetRegistry().GetTool("echo-tool")
		require.NoError(t, err)

		result, err := echoTool.Execute(ctx, map[string]interface{}{
			"message": "multi-tool test",
		}, server.GetEngine())

		require.NoError(t, err)
		assert.Contains(t, result.Content, "multi-tool test")

		// Test ls tool
		lsTool, err := server.GetRegistry().GetTool("ls-tool")
		require.NoError(t, err)

		result, err = lsTool.Execute(ctx, map[string]interface{}{
			"directory": "/tmp",
		}, server.GetEngine())

		require.NoError(t, err)
		assert.Contains(t, result.Content, "ls -la /tmp")
		assert.Equal(t, "/tmp", result.Metadata["directory"])
	})
}

// TestRegistryIntegration tests registry import functionality
func TestRegistryIntegration(t *testing.T) {
	ctx := context.Background()

	// Create a registry with tools
	sourceRegistry := NewRegistry()
	sourceRegistry.RegisterTool(NewContainerTool("source-tool-1", ContainerToolConfig{
		Description: "Source tool 1",
		Image:       "alpine:latest",
	}))
	sourceRegistry.RegisterTool(NewContainerTool("source-tool-2", ContainerToolConfig{
		Description: "Source tool 2",
		Image:       "alpine:latest",
	}))

	// Create server and import registry
	server := NewServer("registry-import-server", "1.0.0").
		ImportRegistry(sourceRegistry).
		AddContainerTool("custom-tool", ContainerToolConfig{
			Description: "Custom tool",
			Image:       "alpine:latest",
		}).
		Build()

	require.NotNil(t, server)

	t.Run("RegistryImport", func(t *testing.T) {
		assert.Equal(t, 3, server.GetRegistry().ToolCount())

		tools := server.GetRegistry().ListTools()
		assert.Contains(t, tools, "source-tool-1")
		assert.Contains(t, tools, "source-tool-2")
		assert.Contains(t, tools, "custom-tool")
	})

	t.Run("ImportedToolsExecutable", func(t *testing.T) {
		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		// Verify imported tools are accessible and executable
		for _, toolName := range []string{"source-tool-1", "source-tool-2", "custom-tool"} {
			tool, err := server.GetRegistry().GetTool(toolName)
			require.NoError(t, err)
			assert.Equal(t, toolName, tool.Name())
		}
	})
}

// TestFrameworkErrorConditions tests various error conditions
func TestFrameworkErrorConditions(t *testing.T) {
	ctx := context.Background()

	t.Run("ToolExecutionError", func(t *testing.T) {
		// Create a tool that deliberately returns an error
		failingTool := NewContainerTool("failing-tool", ContainerToolConfig{
			Description: "Tool that deliberately fails",
			Image:       "alpine:latest",
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
				// Simulate a deliberate failure
				err := fmt.Errorf("deliberate test failure")
				return &ToolResult{
					Error:   err,
					Content: "Failed as expected",
				}, err
			},
		})

		server := NewServer("failing-server", "1.0.0").
			AddTool(failingTool).
			Build()

		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		tool, err := server.GetRegistry().GetTool("failing-tool")
		require.NoError(t, err)

		result, err := tool.Execute(ctx, map[string]interface{}{}, server.GetEngine())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deliberate test failure")
		assert.NotNil(t, result)
		assert.NotNil(t, result.Error)
	})

	t.Run("InvalidParameters", func(t *testing.T) {
		// Create a tool that validates parameters
		validatingTool := NewContainerTool("validating-tool", ContainerToolConfig{
			Description: "Tool that validates parameters",
			Image:       "alpine:latest",
			Parameters: []Parameter{
				{
					Name:        "required_param",
					Type:        "string",
					Description: "A required parameter",
					Required:    true,
				},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
				// Check if required parameter is present
				if _, ok := params["required_param"]; !ok {
					err := fmt.Errorf("required parameter 'required_param' is missing")
					return &ToolResult{Error: err}, err
				}

				return &ToolResult{
					Content: "Parameters validated successfully",
				}, nil
			},
		})

		server := NewServer("validation-server", "1.0.0").
			AddTool(validatingTool).
			Build()

		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		tool, err := server.GetRegistry().GetTool("validating-tool")
		require.NoError(t, err)

		// Test with missing required parameter
		result, err := tool.Execute(ctx, map[string]interface{}{}, server.GetEngine())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required parameter")
		assert.NotNil(t, result.Error)

		// Test with valid parameters
		result, err = tool.Execute(ctx, map[string]interface{}{
			"required_param": "test_value",
		}, server.GetEngine())
		assert.NoError(t, err)
		assert.Contains(t, result.Content, "validated successfully")
	})

	t.Run("ContextChecking", func(t *testing.T) {
		// Create a tool that checks context cancellation
		contextTool := NewContainerTool("context-tool", ContainerToolConfig{
			Description: "Tool that checks context",
			Image:       "alpine:latest",
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
				// Check if context is already cancelled
				select {
				case <-ctx.Done():
					err := ctx.Err()
					return &ToolResult{Error: err}, err
				default:
					// Context is not cancelled, proceed normally
				}

				return &ToolResult{
					Content: "Context check passed",
				}, nil
			},
		})

		server := NewServer("context-test-server", "1.0.0").
			AddTool(contextTool).
			Build()

		err := server.Start(ctx)
		require.NoError(t, err)
		defer server.Close()

		tool, err := server.GetRegistry().GetTool("context-tool")
		require.NoError(t, err)

		// Test with normal context
		result, err := tool.Execute(ctx, map[string]interface{}{}, server.GetEngine())
		assert.NoError(t, err)
		assert.Contains(t, result.Content, "Context check passed")

		// Test with cancelled context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		result, err = tool.Execute(cancelCtx, map[string]interface{}{}, server.GetEngine())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}
