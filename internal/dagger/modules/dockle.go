package modules

import (
	"context"

	"dagger.io/dagger"
)

type DockleModule struct {
	client *dagger.Client
}

func NewDockleModule(client *dagger.Client) *DockleModule {
	return &DockleModule{
		client: client,
	}
}

// ScanImage scans a container image for security issues using Dockle
func (m *DockleModule) ScanImage(ctx context.Context, imageRef string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		DockleVersion: "v0.4.14",
		Format:        "json",
		ExitLevel:     "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:" + config.DockleVersion).
		WithWorkdir("/workspace")

	args := []string{}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add exit level
	if config.ExitLevel != "" {
		args = append(args, "--exit-level", config.ExitLevel)
	}

	// Add accept keys for ignoring specific issues
	for _, key := range config.AcceptKey {
		args = append(args, "--accept-key", key)
	}

	// Add accept files for ignoring file-based issues
	for _, file := range config.AcceptFile {
		args = append(args, "--accept-file", file)
	}

	// Add ignore rules
	for _, ignore := range config.Ignore {
		args = append(args, "--ignore", ignore)
	}

	// Add image reference
	args = append(args, imageRef)

	return container.WithExec(args), nil
}

// ScanTarball scans a container image tarball
func (m *DockleModule) ScanTarball(ctx context.Context, tarballPath string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		DockleVersion: "v0.4.14",
		Format:        "json",
		ExitLevel:     "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:" + config.DockleVersion).
		WithWorkdir("/workspace")

	// Mount tarball file
	if tarballPath != "" {
		container = container.WithMountedFile("/workspace/image.tar", m.client.Host().File(tarballPath))
	}

	args := []string{}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add exit level
	if config.ExitLevel != "" {
		args = append(args, "--exit-level", config.ExitLevel)
	}

	// Add accept keys for ignoring specific issues
	for _, key := range config.AcceptKey {
		args = append(args, "--accept-key", key)
	}

	// Add ignore rules
	for _, ignore := range config.Ignore {
		args = append(args, "--ignore", ignore)
	}

	// Add input flag for tarball
	args = append(args, "--input", "/workspace/image.tar")

	return container.WithExec(args), nil
}

// ScanDockerfile scans a Dockerfile for best practices
func (m *DockleModule) ScanDockerfile(ctx context.Context, dockerfilePath string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		DockleVersion: "v0.4.14",
		Format:        "json",
		ExitLevel:     "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:" + config.DockleVersion).
		WithWorkdir("/workspace")

	// Mount Dockerfile
	if dockerfilePath != "" {
		container = container.WithMountedFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath))
	}

	args := []string{}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add exit level
	if config.ExitLevel != "" {
		args = append(args, "--exit-level", config.ExitLevel)
	}

	// Add accept keys for ignoring specific issues
	for _, key := range config.AcceptKey {
		args = append(args, "--accept-key", key)
	}

	// Add ignore rules
	for _, ignore := range config.Ignore {
		args = append(args, "--ignore", ignore)
	}

	// Add Dockerfile flag
	args = append(args, "--input", "/workspace/Dockerfile")

	return container.WithExec(args), nil
}

// GenerateConfig generates a Dockle configuration file
func (m *DockleModule) GenerateConfig(ctx context.Context, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "yq"}).
		WithWorkdir("/workspace")

	// Create a basic Dockle configuration
	dockleConfig := `# Dockle Configuration
# Accept keys for specific CIS Docker Benchmark checks
accept-key:
  - CIS-DI-0001  # Create a user for the container
  - CIS-DI-0005  # Enable content trust for Docker
  - CIS-DI-0006  # Add HEALTHCHECK instruction to the container image

# Accept file patterns that are safe to ignore
accept-file:
  - /usr/share/man/*
  - /tmp/*
  - /var/tmp/*

# Ignore specific findings
ignore:
  - DKL-DI-0006  # Avoid latest tag
  - DKL-LI-0003  # Only one EXPOSE instruction is available

# Exit level configuration
exit-level: warn
`

	container = container.
		WithNewFile("/workspace/.dockleignore", dockleConfig).
		WithExec([]string{"cat", "/workspace/.dockleignore"})

	return container, nil
}

// ListChecks lists all available Dockle security checks
func (m *DockleModule) ListChecks(ctx context.Context, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		DockleVersion: "v0.4.14",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:" + config.DockleVersion)

	args := []string{"--help"}

	return container.WithExec(args), nil
}

// ScanWithPolicy scans using a custom policy file
func (m *DockleModule) ScanWithPolicy(ctx context.Context, imageRef string, policyPath string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		DockleVersion: "v0.4.14",
		Format:        "json",
		ExitLevel:     "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:" + config.DockleVersion).
		WithWorkdir("/workspace")

	// Mount policy file
	if policyPath != "" {
		container = container.WithMountedFile("/workspace/.dockleignore", m.client.Host().File(policyPath))
	}

	args := []string{}

	// Add format
	if config.Format != "" {
		args = append(args, "--format", config.Format)
	}

	// Add output file
	if config.Output != "" {
		args = append(args, "--output", config.Output)
	}

	// Add exit level
	if config.ExitLevel != "" {
		args = append(args, "--exit-level", config.ExitLevel)
	}

	// Add image reference
	args = append(args, imageRef)

	return container.WithExec(args), nil
}

type DockleConfig struct {
	DockleVersion string
	Format        string
	Output        string
	ExitLevel     string
	AcceptKey     []string
	AcceptFile    []string
	Ignore        []string
}

type DockleOption func(*DockleConfig)

func WithDockleVersion(version string) DockleOption {
	return func(c *DockleConfig) {
		c.DockleVersion = version
	}
}

func WithDockleFormat(format string) DockleOption {
	return func(c *DockleConfig) {
		c.Format = format
	}
}

func WithDockleOutput(output string) DockleOption {
	return func(c *DockleConfig) {
		c.Output = output
	}
}

func WithExitLevel(level string) DockleOption {
	return func(c *DockleConfig) {
		c.ExitLevel = level
	}
}

func WithAcceptKey(keys []string) DockleOption {
	return func(c *DockleConfig) {
		c.AcceptKey = keys
	}
}

func WithAcceptFile(files []string) DockleOption {
	return func(c *DockleConfig) {
		c.AcceptFile = files
	}
}

func WithDockleIgnore(ignores []string) DockleOption {
	return func(c *DockleConfig) {
		c.Ignore = ignores
	}
}
