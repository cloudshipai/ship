package ship

import (
	"fmt"
	"sync"
)

// Registry manages a collection of tools, prompts, and resources
type Registry struct {
	mu        sync.RWMutex
	tools     map[string]Tool
	prompts   map[string]Prompt
	resources map[string]Resource
}

// Prompt represents an MCP prompt
type Prompt interface {
	Name() string
	Description() string
	Arguments() []PromptArgument
}

// PromptArgument represents a prompt argument
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// Resource represents an MCP resource
type Resource interface {
	URI() string
	Name() string
	Description() string
	MIMEType() string
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		tools:     make(map[string]Tool),
		prompts:   make(map[string]Prompt),
		resources: make(map[string]Resource),
	}
}

// RegisterTool registers a tool in the registry
func (r *Registry) RegisterTool(tool Tool) error {
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}

	if tool.Name() == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[tool.Name()]; exists {
		return fmt.Errorf("tool with name '%s' already exists", tool.Name())
	}

	r.tools[tool.Name()] = tool
	return nil
}

// GetTool retrieves a tool by name
func (r *Registry) GetTool(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrToolNotFound, name)
	}

	return tool, nil
}

// ListTools returns all registered tool names
func (r *Registry) ListTools() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}

	return names
}

// GetAllTools returns all registered tools
func (r *Registry) GetAllTools() map[string]Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make(map[string]Tool)
	for name, tool := range r.tools {
		tools[name] = tool
	}

	return tools
}

// ImportFrom imports all tools from another registry
func (r *Registry) ImportFrom(other *Registry) error {
	other.mu.RLock()
	defer other.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	for name, tool := range other.tools {
		if _, exists := r.tools[name]; exists {
			return fmt.Errorf("tool with name '%s' already exists", name)
		}
		r.tools[name] = tool
	}

	for name, prompt := range other.prompts {
		if _, exists := r.prompts[name]; exists {
			return fmt.Errorf("prompt with name '%s' already exists", name)
		}
		r.prompts[name] = prompt
	}

	for uri, resource := range other.resources {
		if _, exists := r.resources[uri]; exists {
			return fmt.Errorf("resource with URI '%s' already exists", uri)
		}
		r.resources[uri] = resource
	}

	return nil
}

// RegisterPrompt registers a prompt in the registry
func (r *Registry) RegisterPrompt(prompt Prompt) error {
	if prompt == nil {
		return fmt.Errorf("prompt cannot be nil")
	}

	if prompt.Name() == "" {
		return fmt.Errorf("prompt name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prompts[prompt.Name()]; exists {
		return fmt.Errorf("prompt with name '%s' already exists", prompt.Name())
	}

	r.prompts[prompt.Name()] = prompt
	return nil
}

// GetPrompt retrieves a prompt by name
func (r *Registry) GetPrompt(name string) (Prompt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompt, exists := r.prompts[name]
	if !exists {
		return nil, fmt.Errorf("prompt not found: %s", name)
	}

	return prompt, nil
}

// RegisterResource registers a resource in the registry
func (r *Registry) RegisterResource(resource Resource) error {
	if resource == nil {
		return fmt.Errorf("resource cannot be nil")
	}

	if resource.URI() == "" {
		return fmt.Errorf("resource URI cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resource.URI()]; exists {
		return fmt.Errorf("resource with URI '%s' already exists", resource.URI())
	}

	r.resources[resource.URI()] = resource
	return nil
}

// GetResource retrieves a resource by URI
func (r *Registry) GetResource(uri string) (Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, exists := r.resources[uri]
	if !exists {
		return nil, fmt.Errorf("resource not found: %s", uri)
	}

	return resource, nil
}

// Clear removes all registered items
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tools = make(map[string]Tool)
	r.prompts = make(map[string]Prompt)
	r.resources = make(map[string]Resource)
}

// ToolCount returns the number of registered tools
func (r *Registry) ToolCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tools)
}

// Global default registry used by Ship CLI
var DefaultRegistry = NewRegistry()
