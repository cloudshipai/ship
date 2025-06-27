package modules

import (
	"context"
	"time"
)

// ModuleType represents the type of module
type ModuleType string

const (
	ModuleTypeDocker ModuleType = "docker"
	ModuleTypeDagger ModuleType = "dagger"
)

// Module represents a discoverable Ship CLI module
type Module struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   ModuleMetadata `yaml:"metadata"`
	Spec       ModuleSpec     `yaml:"spec"`
	
	// Runtime fields
	Path       string    `yaml:"-"`
	Source     string    `yaml:"-"`
	LoadedAt   time.Time `yaml:"-"`
	Trusted    bool      `yaml:"-"`
}

// ModuleMetadata contains module identification information
type ModuleMetadata struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	Tags        []string          `yaml:"tags,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// ModuleSpec defines the module's behavior and integration
type ModuleSpec struct {
	Type         ModuleType          `yaml:"type"`
	Docker       *DockerModuleSpec   `yaml:"docker,omitempty"`
	Dagger       *DaggerModuleSpec   `yaml:"dagger,omitempty"`
	Commands     []ModuleCommand     `yaml:"commands"`
	Dependencies []string            `yaml:"dependencies,omitempty"`
	Permissions  []string            `yaml:"permissions,omitempty"`
}

// DockerModuleSpec defines Docker-based module configuration
type DockerModuleSpec struct {
	Image      string            `yaml:"image"`
	Entrypoint []string          `yaml:"entrypoint,omitempty"`
	Env        map[string]string `yaml:"env,omitempty"`
	WorkingDir string            `yaml:"workingDir,omitempty"`
	Volumes    []VolumeMount     `yaml:"volumes,omitempty"`
}

// DaggerModuleSpec defines Dagger-based module configuration
type DaggerModuleSpec struct {
	Module   string `yaml:"module"`
	Function string `yaml:"function"`
}

// VolumeMount represents a volume mount for Docker modules
type VolumeMount struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Type   string `yaml:"type,omitempty"` // bind, volume, tmpfs
}

// ModuleCommand represents a CLI command provided by the module
type ModuleCommand struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Usage       string       `yaml:"usage,omitempty"`
	Flags       []ModuleFlag `yaml:"flags,omitempty"`
	Examples    []string     `yaml:"examples,omitempty"`
}

// ModuleFlag represents a command-line flag
type ModuleFlag struct {
	Name        string      `yaml:"name"`
	Short       string      `yaml:"short,omitempty"`
	Type        string      `yaml:"type"`        // string, int, bool, []string
	Default     interface{} `yaml:"default,omitempty"`
	Required    bool        `yaml:"required,omitempty"`
	Description string      `yaml:"description"`
	Enum        []string    `yaml:"enum,omitempty"`
}

// ModuleExecutor defines the interface for executing modules
type ModuleExecutor interface {
	Execute(ctx context.Context, module *Module, command string, args []string, flags map[string]interface{}) (*ExecutionResult, error)
	CanExecute(module *Module) bool
}

// ExecutionResult represents the result of module execution
type ExecutionResult struct {
	ExitCode int               `json:"exitCode"`
	Stdout   string            `json:"stdout"`
	Stderr   string            `json:"stderr"`
	Duration time.Duration     `json:"duration"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ModuleSource represents where a module was discovered from
type ModuleSource struct {
	Type        string    `yaml:"type"`        // builtin, user, project, git
	Path        string    `yaml:"path"`
	URL         string    `yaml:"url,omitempty"`
	Ref         string    `yaml:"ref,omitempty"`
	LastUpdated time.Time `yaml:"lastUpdated,omitempty"`
	Trusted     bool      `yaml:"trusted"`
}

// ModuleConfig represents module configuration from ship config
type ModuleConfig struct {
	Repositories     []GitRepository `yaml:"repositories,omitempty"`
	Directories      []string        `yaml:"directories,omitempty"`
	AllowUntrusted   bool            `yaml:"allow_untrusted"`
	Sandbox          bool            `yaml:"sandbox"`
	CacheDir         string          `yaml:"cache_dir,omitempty"`
	UpdateInterval   string          `yaml:"update_interval,omitempty"`
}

// GitRepository represents a git-based module source
type GitRepository struct {
	URL     string `yaml:"url"`
	Ref     string `yaml:"ref"`
	Path    string `yaml:"path,omitempty"`
	Trusted bool   `yaml:"trusted"`
}

// Discovery interface for module discovery
type Discovery interface {
	DiscoverModules(ctx context.Context) ([]*Module, error)
	GetSourceType() string
	SetConfig(config ModuleConfig)
}