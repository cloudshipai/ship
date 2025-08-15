package ship

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	assert.NotNil(t, registry)
	assert.Equal(t, 0, registry.ToolCount())
	assert.Empty(t, registry.ListTools())
}

func TestRegistry_RegisterTool(t *testing.T) {
	registry := NewRegistry()

	t.Run("successful registration", func(t *testing.T) {
		tool := createTestTool("test-tool", "Test tool")

		err := registry.RegisterTool(tool)

		assert.NoError(t, err)
		assert.Equal(t, 1, registry.ToolCount())
		assert.Contains(t, registry.ListTools(), "test-tool")
	})

	t.Run("register nil tool", func(t *testing.T) {
		err := registry.RegisterTool(nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool cannot be nil")
	})

	t.Run("register tool with empty name", func(t *testing.T) {
		tool := createTestTool("", "Test tool")

		err := registry.RegisterTool(tool)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool name cannot be empty")
	})

	t.Run("register duplicate tool", func(t *testing.T) {
		tool1 := createTestTool("duplicate", "Tool 1")
		tool2 := createTestTool("duplicate", "Tool 2")

		err1 := registry.RegisterTool(tool1)
		err2 := registry.RegisterTool(tool2)

		assert.NoError(t, err1)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "already exists")
	})
}

func TestRegistry_GetTool(t *testing.T) {
	registry := NewRegistry()
	tool := createTestTool("test-tool", "Test tool")
	registry.RegisterTool(tool)

	t.Run("get existing tool", func(t *testing.T) {
		retrieved, err := registry.GetTool("test-tool")

		assert.NoError(t, err)
		assert.Equal(t, tool, retrieved)
	})

	t.Run("get non-existent tool", func(t *testing.T) {
		retrieved, err := registry.GetTool("non-existent")

		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.ErrorIs(t, err, ErrToolNotFound)
	})
}

func TestRegistry_ListTools(t *testing.T) {
	registry := NewRegistry()

	tool1 := createTestTool("tool1", "Tool 1")
	tool2 := createTestTool("tool2", "Tool 2")
	tool3 := createTestTool("tool3", "Tool 3")

	registry.RegisterTool(tool1)
	registry.RegisterTool(tool2)
	registry.RegisterTool(tool3)

	tools := registry.ListTools()

	assert.Len(t, tools, 3)
	assert.Contains(t, tools, "tool1")
	assert.Contains(t, tools, "tool2")
	assert.Contains(t, tools, "tool3")
}

func TestRegistry_GetAllTools(t *testing.T) {
	registry := NewRegistry()

	tool1 := createTestTool("tool1", "Tool 1")
	tool2 := createTestTool("tool2", "Tool 2")

	registry.RegisterTool(tool1)
	registry.RegisterTool(tool2)

	allTools := registry.GetAllTools()

	assert.Len(t, allTools, 2)
	assert.Equal(t, tool1, allTools["tool1"])
	assert.Equal(t, tool2, allTools["tool2"])
}

func TestRegistry_ImportFrom(t *testing.T) {
	source := NewRegistry()
	target := NewRegistry()

	tool1 := createTestTool("tool1", "Tool 1")
	tool2 := createTestTool("tool2", "Tool 2")

	source.RegisterTool(tool1)
	source.RegisterTool(tool2)

	t.Run("successful import", func(t *testing.T) {
		err := target.ImportFrom(source)

		assert.NoError(t, err)
		assert.Equal(t, 2, target.ToolCount())

		retrievedTool1, _ := target.GetTool("tool1")
		retrievedTool2, _ := target.GetTool("tool2")

		assert.Equal(t, tool1, retrievedTool1)
		assert.Equal(t, tool2, retrievedTool2)
	})

	t.Run("import with conflicts", func(t *testing.T) {
		conflict := NewRegistry()
		conflictTool := createTestTool("tool1", "Conflict tool")
		conflict.RegisterTool(conflictTool)

		err := target.ImportFrom(conflict)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestRegistry_Clear(t *testing.T) {
	registry := NewRegistry()

	tool := createTestTool("test-tool", "Test tool")
	registry.RegisterTool(tool)

	assert.Equal(t, 1, registry.ToolCount())

	registry.Clear()

	assert.Equal(t, 0, registry.ToolCount())
	assert.Empty(t, registry.ListTools())
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewRegistry()

	// Test concurrent registration and retrieval
	done := make(chan bool, 10)

	// Register tools concurrently
	for i := 0; i < 5; i++ {
		go func(id int) {
			tool := createTestTool(fmt.Sprintf("tool-%d", id), fmt.Sprintf("Tool %d", id))
			err := registry.RegisterTool(tool)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Read tools concurrently
	for i := 0; i < 5; i++ {
		go func() {
			_ = registry.ListTools()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	assert.Equal(t, 5, registry.ToolCount())
}

func TestDefaultRegistry(t *testing.T) {
	// Test that DefaultRegistry is initialized
	assert.NotNil(t, DefaultRegistry)

	// Test that it's a proper registry
	assert.Equal(t, 0, DefaultRegistry.ToolCount())

	// Test registration works
	tool := createTestTool("default-test", "Default test tool")
	err := DefaultRegistry.RegisterTool(tool)
	assert.NoError(t, err)

	// Clean up for other tests
	DefaultRegistry.Clear()
}

// Helper function to create test tools
func createTestTool(name, description string) Tool {
	config := ContainerToolConfig{
		Description: description,
		Image:       "test:latest",
		Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
			return &ToolResult{Content: "test output"}, nil
		},
	}

	return NewContainerTool(name, config)
}
