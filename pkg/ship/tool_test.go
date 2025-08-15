package ship

import (
	"context"
	"testing"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/stretchr/testify/assert"
)

func TestParameter(t *testing.T) {
	t.Run("basic parameter creation", func(t *testing.T) {
		param := Parameter{
			Name:        "directory",
			Type:        "string",
			Description: "Directory to scan",
			Required:    true,
			Enum:        []string{".", "/tmp"},
		}

		assert.Equal(t, "directory", param.Name)
		assert.Equal(t, "string", param.Type)
		assert.Equal(t, "Directory to scan", param.Description)
		assert.True(t, param.Required)
		assert.Equal(t, []string{".", "/tmp"}, param.Enum)
	})
}

func TestToolResult(t *testing.T) {
	t.Run("successful result", func(t *testing.T) {
		result := &ToolResult{
			Content: "scan completed",
			Metadata: map[string]interface{}{
				"duration":      "5s",
				"files_scanned": 42,
			},
		}

		assert.Equal(t, "scan completed", result.Content)
		assert.Equal(t, "5s", result.Metadata["duration"])
		assert.Equal(t, 42, result.Metadata["files_scanned"])
		assert.Nil(t, result.Error)
	})

	t.Run("error result", func(t *testing.T) {
		result := &ToolResult{
			Content: "",
			Error:   ErrInvalidParameter,
		}

		assert.Empty(t, result.Content)
		assert.Equal(t, ErrInvalidParameter, result.Error)
	})
}

func TestContainerTool(t *testing.T) {
	t.Run("tool creation and basic properties", func(t *testing.T) {
		params := []Parameter{
			{Name: "directory", Type: "string", Required: true},
			{Name: "format", Type: "string", Enum: []string{"json", "text"}},
		}

		config := ContainerToolConfig{
			Description: "Test scanner tool",
			Image:       "test/scanner:latest",
			Command:     []string{"scanner", "--scan"},
			Parameters:  params,
		}

		tool := NewContainerTool("test-scanner", config)

		assert.Equal(t, "test-scanner", tool.Name())
		assert.Equal(t, "Test scanner tool", tool.Description())
		assert.Equal(t, params, tool.Parameters())
	})

	t.Run("tool execution without executor", func(t *testing.T) {
		config := ContainerToolConfig{
			Description: "Test tool",
			Image:       "test:latest",
		}

		tool := NewContainerTool("test-tool", config)

		result, err := tool.Execute(context.Background(), map[string]interface{}{}, nil)

		assert.Error(t, err)
		assert.Equal(t, ErrExecutorNotSet, err)
		assert.Equal(t, ErrExecutorNotSet, result.Error)
	})

	t.Run("tool execution with executor", func(t *testing.T) {
		expectedResult := &ToolResult{
			Content:  "test output",
			Metadata: map[string]interface{}{"test": true},
		}

		executor := func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			// Verify parameters are passed correctly
			assert.Equal(t, "test-value", params["test-param"])
			return expectedResult, nil
		}

		config := ContainerToolConfig{
			Description: "Test tool",
			Image:       "test:latest",
			Execute:     executor,
		}

		tool := NewContainerTool("test-tool", config)

		params := map[string]interface{}{
			"test-param": "test-value",
		}

		result, err := tool.Execute(context.Background(), params, nil)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("tool execution with executor error", func(t *testing.T) {
		executor := func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			return &ToolResult{
				Content: "",
				Error:   ErrInvalidParameter,
			}, ErrInvalidParameter
		}

		config := ContainerToolConfig{
			Description: "Test tool",
			Execute:     executor,
		}

		tool := NewContainerTool("test-tool", config)

		result, err := tool.Execute(context.Background(), map[string]interface{}{}, nil)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidParameter, err)
		assert.Equal(t, ErrInvalidParameter, result.Error)
	})
}

func TestContainerToolConfig(t *testing.T) {
	t.Run("complete configuration", func(t *testing.T) {
		params := []Parameter{
			{Name: "input", Type: "string", Required: true},
			{Name: "output", Type: "string", Required: false},
		}

		executor := func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			return &ToolResult{Content: "success"}, nil
		}

		config := ContainerToolConfig{
			Description: "Complete test tool",
			Image:       "test/tool:v1.0.0",
			Command:     []string{"tool", "--verbose"},
			Parameters:  params,
			Execute:     executor,
		}

		// Verify all fields are set correctly
		assert.Equal(t, "Complete test tool", config.Description)
		assert.Equal(t, "test/tool:v1.0.0", config.Image)
		assert.Equal(t, []string{"tool", "--verbose"}, config.Command)
		assert.Equal(t, params, config.Parameters)
		assert.NotNil(t, config.Execute)
	})
}
