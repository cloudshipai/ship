package ship

import (
	"context"

	"github.com/cloudshipai/ship/pkg/dagger"
)

// Tool represents a tool that can be executed via MCP
type Tool interface {
	Name() string
	Description() string
	Parameters() []Parameter
	Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error)
}

// Parameter represents a tool parameter
type Parameter struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Error    error                  `json:"error,omitempty"`
}

// ToolExecutor is a function type for executing container-based tools
type ToolExecutor func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error)

// ContainerToolConfig represents configuration for a container-based tool
type ContainerToolConfig struct {
	Description string
	Image       string
	Command     []string
	Parameters  []Parameter
	Execute     ToolExecutor
}

// ContainerTool implements Tool interface for container-based tools
type ContainerTool struct {
	name        string
	description string
	image       string
	command     []string
	parameters  []Parameter
	executor    ToolExecutor
}

// NewContainerTool creates a new container-based tool
func NewContainerTool(name string, config ContainerToolConfig) Tool {
	return &ContainerTool{
		name:        name,
		description: config.Description,
		image:       config.Image,
		command:     config.Command,
		parameters:  config.Parameters,
		executor:    config.Execute,
	}
}

// Name returns the tool name
func (t *ContainerTool) Name() string {
	return t.name
}

// Description returns the tool description
func (t *ContainerTool) Description() string {
	return t.description
}

// Parameters returns the tool parameters
func (t *ContainerTool) Parameters() []Parameter {
	return t.parameters
}

// Execute runs the container tool
func (t *ContainerTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ToolResult, error) {
	if t.executor == nil {
		return &ToolResult{
			Content: "",
			Error:   ErrExecutorNotSet,
		}, ErrExecutorNotSet
	}

	return t.executor(ctx, params, engine)
}
