package tools

import (
	"context"
	"testing"

	"github.com/cloudshipai/ship/internal/ship"
	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTFLintTool(t *testing.T) {
	tool := NewTFLintTool()

	assert.Equal(t, "tflint", tool.Name())
	assert.Equal(t, "Run TFLint on Terraform code to check for syntax errors and best practices", tool.Description())

	params := tool.Parameters()
	assert.Len(t, params, 6)

	// Verify parameter names
	paramNames := make(map[string]bool)
	for _, param := range params {
		paramNames[param.Name] = true
	}

	expectedParams := []string{"directory", "format", "config", "enable_rules", "disable_rules", "init"}
	for _, expected := range expectedParams {
		assert.True(t, paramNames[expected], "Expected parameter %s not found", expected)
	}
}

func TestTFLintParameterTypes(t *testing.T) {
	tool := NewTFLintTool()
	params := tool.Parameters()

	paramMap := make(map[string]ship.Parameter)
	for _, param := range params {
		paramMap[param.Name] = param
	}

	t.Run("string parameters", func(t *testing.T) {
		stringParams := []string{"directory", "format", "config", "enable_rules", "disable_rules"}
		for _, name := range stringParams {
			param, exists := paramMap[name]
			assert.True(t, exists, "Parameter %s should exist", name)
			assert.Equal(t, "string", param.Type, "Parameter %s should be string type", name)
		}
	})

	t.Run("boolean parameters", func(t *testing.T) {
		param, exists := paramMap["init"]
		assert.True(t, exists)
		assert.Equal(t, "boolean", param.Type)
		assert.False(t, param.Required)
	})

	t.Run("format enum", func(t *testing.T) {
		param, exists := paramMap["format"]
		assert.True(t, exists)
		assert.Equal(t, []string{"default", "json", "compact"}, param.Enum)
	})
}

func TestTFLintExecution(t *testing.T) {
	ctx := context.Background()
	engine, err := dagger.NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()

	tool := NewTFLintTool()

	t.Run("basic execution", func(t *testing.T) {
		params := map[string]interface{}{
			"directory": ".",
			"format":    "json",
		}

		result, err := tool.Execute(ctx, params, engine)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Content)

		// Verify metadata
		assert.Equal(t, "json", result.Metadata["format"])
		assert.Equal(t, ".", result.Metadata["directory"])
		assert.Equal(t, "tflint", result.Metadata["tool"])
		assert.Equal(t, TFLintImage, result.Metadata["image"])
	})

	t.Run("execution with rules", func(t *testing.T) {
		params := map[string]interface{}{
			"directory":     ".",
			"format":        "json",
			"enable_rules":  "rule1,rule2",
			"disable_rules": "rule3,rule4",
			"init":          false,
		}

		result, err := tool.Execute(ctx, params, engine)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Content)

		// Verify metadata
		assert.Equal(t, false, result.Metadata["initialized"])
	})

	t.Run("execution with config file", func(t *testing.T) {
		params := map[string]interface{}{
			"directory": ".",
			"config":    ".tflint.hcl",
			"format":    "compact",
		}

		result, err := tool.Execute(ctx, params, engine)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Verify metadata
		assert.Equal(t, ".tflint.hcl", result.Metadata["config_file"])
		assert.Equal(t, "compact", result.Metadata["format"])
	})
}

func TestTFLintParameterHelpers(t *testing.T) {
	t.Run("getStringParam", func(t *testing.T) {
		params := map[string]interface{}{
			"existing": "value",
			"empty":    "",
		}

		assert.Equal(t, "value", getStringParam(params, "existing", "default"))
		assert.Equal(t, "", getStringParam(params, "empty", "default"))
		assert.Equal(t, "default", getStringParam(params, "missing", "default"))
	})

	t.Run("getBoolParam", func(t *testing.T) {
		params := map[string]interface{}{
			"bool_true":    true,
			"bool_false":   false,
			"string_true":  "true",
			"string_false": "false",
			"string_TRUE":  "TRUE",
		}

		assert.True(t, getBoolParam(params, "bool_true", false))
		assert.False(t, getBoolParam(params, "bool_false", true))
		assert.True(t, getBoolParam(params, "string_true", false))
		assert.False(t, getBoolParam(params, "string_false", true))
		assert.True(t, getBoolParam(params, "string_TRUE", false))
		assert.True(t, getBoolParam(params, "missing", true))
		assert.False(t, getBoolParam(params, "missing", false))
	})
}
