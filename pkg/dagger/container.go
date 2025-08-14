package dagger

import (
	"context"
	"fmt"
	"strings"
)

// ContainerBuilder provides a fluent API for building containers
type ContainerBuilder struct {
	engine   *Engine
	image    string
	mounts   map[string]string
	env      map[string]string
	workdir  string
	commands [][]string
	err      error
}

// From sets the base image for the container
func (c *ContainerBuilder) From(image string) *ContainerBuilder {
	if c.err != nil {
		return c
	}
	
	if image == "" {
		c.err = fmt.Errorf("image cannot be empty")
		return c
	}
	
	c.image = image
	return c
}

// WithMountedDirectory mounts a directory into the container
func (c *ContainerBuilder) WithMountedDirectory(containerPath, hostPath string) *ContainerBuilder {
	if c.err != nil {
		return c
	}
	
	if containerPath == "" {
		c.err = fmt.Errorf("container path cannot be empty")
		return c
	}
	
	if hostPath == "" {
		c.err = fmt.Errorf("host path cannot be empty")
		return c
	}
	
	c.mounts[containerPath] = hostPath
	return c
}

// WithEnvVariable sets an environment variable
func (c *ContainerBuilder) WithEnvVariable(name, value string) *ContainerBuilder {
	if c.err != nil {
		return c
	}
	
	if name == "" {
		c.err = fmt.Errorf("environment variable name cannot be empty")
		return c
	}
	
	c.env[name] = value
	return c
}

// WithWorkdir sets the working directory
func (c *ContainerBuilder) WithWorkdir(workdir string) *ContainerBuilder {
	if c.err != nil {
		return c
	}
	
	if workdir == "" {
		c.err = fmt.Errorf("workdir cannot be empty")
		return c
	}
	
	c.workdir = workdir
	return c
}

// WithExec adds a command to execute
func (c *ContainerBuilder) WithExec(args []string) *ContainerBuilder {
	if c.err != nil {
		return c
	}
	
	if len(args) == 0 {
		c.err = fmt.Errorf("exec args cannot be empty")
		return c
	}
	
	// Create a copy to avoid mutation
	cmdCopy := make([]string, len(args))
	copy(cmdCopy, args)
	
	c.commands = append(c.commands, cmdCopy)
	return c
}

// Stdout executes the container and returns stdout
func (c *ContainerBuilder) Stdout(ctx context.Context) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	
	if c.engine.closed {
		return "", fmt.Errorf("engine is closed")
	}
	
	if c.image == "" {
		return "", fmt.Errorf("no base image specified")
	}
	
	// For testing purposes, return a mock output based on the configuration
	// In real implementation, this would execute the container via Dagger
	return c.mockExecution()
}

// Stderr executes the container and returns stderr
func (c *ContainerBuilder) Stderr(ctx context.Context) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	
	if c.engine.closed {
		return "", fmt.Errorf("engine is closed")
	}
	
	// Mock implementation
	return "", nil
}

// CombinedOutput executes the container and returns combined stdout/stderr
func (c *ContainerBuilder) CombinedOutput(ctx context.Context) (string, error) {
	stdout, err := c.Stdout(ctx)
	if err != nil {
		return "", err
	}
	
	stderr, _ := c.Stderr(ctx)
	
	if stderr != "" {
		return stdout + "\n" + stderr, nil
	}
	
	return stdout, nil
}

// mockExecution provides a mock implementation for testing
func (c *ContainerBuilder) mockExecution() (string, error) {
	var output strings.Builder
	
	output.WriteString(fmt.Sprintf("Mock execution of image: %s\n", c.image))
	
	if c.workdir != "" {
		output.WriteString(fmt.Sprintf("Working directory: %s\n", c.workdir))
	}
	
	if len(c.env) > 0 {
		output.WriteString("Environment variables:\n")
		for k, v := range c.env {
			output.WriteString(fmt.Sprintf("  %s=%s\n", k, v))
		}
	}
	
	if len(c.mounts) > 0 {
		output.WriteString("Mounted directories:\n")
		for container, host := range c.mounts {
			output.WriteString(fmt.Sprintf("  %s -> %s\n", host, container))
		}
	}
	
	if len(c.commands) > 0 {
		output.WriteString("Commands executed:\n")
		for i, cmd := range c.commands {
			output.WriteString(fmt.Sprintf("  %d: %s\n", i+1, strings.Join(cmd, " ")))
		}
	}
	
	return output.String(), nil
}

// HostManager manages host directory operations
type HostManager struct {
	engine *Engine
}

// Directory returns a host directory reference
func (h *HostManager) Directory(path string) *DirectoryRef {
	return &DirectoryRef{
		path:   path,
		engine: h.engine,
	}
}

// DirectoryRef represents a reference to a host directory
type DirectoryRef struct {
	path   string
	engine *Engine
}

// Path returns the directory path
func (d *DirectoryRef) Path() string {
	return d.path
}