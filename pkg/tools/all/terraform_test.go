package all

import (
	"testing"

	"github.com/cloudshipai/ship/pkg/ship"
	"github.com/stretchr/testify/assert"
)

func TestTerraformRegistry(t *testing.T) {
	registry := TerraformRegistry()
	
	assert.NotNil(t, registry)
	assert.Equal(t, 1, registry.ToolCount()) // Currently only TFLint
	
	// Verify TFLint tool is present
	tools := registry.ListTools()
	assert.Contains(t, tools, "tflint")
	
	// Verify we can get the tool
	tflintTool, err := registry.GetTool("tflint")
	assert.NoError(t, err)
	assert.NotNil(t, tflintTool)
	assert.Equal(t, "tflint", tflintTool.Name())
}

func TestSecurityRegistry(t *testing.T) {
	registry := SecurityRegistry()
	
	assert.NotNil(t, registry)
	assert.Equal(t, 0, registry.ToolCount()) // No security tools converted yet
}

func TestDocsRegistry(t *testing.T) {
	registry := DocsRegistry()
	
	assert.NotNil(t, registry)
	assert.Equal(t, 0, registry.ToolCount()) // No docs tools converted yet
}

func TestAllRegistry(t *testing.T) {
	registry := AllRegistry()
	
	assert.NotNil(t, registry)
	assert.Equal(t, 1, registry.ToolCount()) // Currently only TFLint from Terraform registry
	
	// Verify all tools from sub-registries are present
	tools := registry.ListTools()
	assert.Contains(t, tools, "tflint")
}

func TestServerBuilderConvenienceFunctions(t *testing.T) {
	t.Run("AddTerraformTools", func(t *testing.T) {
		builder := ship.NewServer("test-server", "1.0.0")
		result := AddTerraformTools(builder)
		
		assert.Equal(t, builder, result) // Should return the same builder for chaining
		
		server := builder.Build()
		assert.Equal(t, 1, server.GetRegistry().ToolCount())
		
		tools := server.GetRegistry().ListTools()
		assert.Contains(t, tools, "tflint")
	})
	
	t.Run("AddSecurityTools", func(t *testing.T) {
		builder := ship.NewServer("test-server", "1.0.0")
		result := AddSecurityTools(builder)
		
		assert.Equal(t, builder, result)
		
		server := builder.Build()
		assert.Equal(t, 0, server.GetRegistry().ToolCount()) // No security tools yet
	})
	
	t.Run("AddDocsTools", func(t *testing.T) {
		builder := ship.NewServer("test-server", "1.0.0")
		result := AddDocsTools(builder)
		
		assert.Equal(t, builder, result)
		
		server := builder.Build()
		assert.Equal(t, 0, server.GetRegistry().ToolCount()) // No docs tools yet
	})
	
	t.Run("AddAllTools", func(t *testing.T) {
		builder := ship.NewServer("test-server", "1.0.0")
		result := AddAllTools(builder)
		
		assert.Equal(t, builder, result)
		
		server := builder.Build()
		assert.Equal(t, 1, server.GetRegistry().ToolCount()) // All tools combined
		
		tools := server.GetRegistry().ListTools()
		assert.Contains(t, tools, "tflint")
	})
	
	t.Run("ChainedUsage", func(t *testing.T) {
		// Test that convenience functions can be chained
		server := ship.NewServer("chained-server", "1.0.0")
		
		// This demonstrates the usage pattern
		result := AddTerraformTools(server.
			AddContainerTool("custom-tool", ship.ContainerToolConfig{
				Description: "Custom tool",
				Image:       "alpine:latest",
			})).
			Build()
		
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.GetRegistry().ToolCount()) // TFLint + custom tool
		
		tools := result.GetRegistry().ListTools()
		assert.Contains(t, tools, "tflint")
		assert.Contains(t, tools, "custom-tool")
	})
}