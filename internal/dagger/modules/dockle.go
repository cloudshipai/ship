package modules

import (
	"context"
	"fmt"

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

// GetVersion returns the version of Dockle
func (m *DockleModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14").
		WithExec([]string{"dockle", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to get dockle version: no output received")
}

// ScanImage scans a container image for security issues using Dockle
func (m *DockleModule) ScanImage(ctx context.Context, imageRef string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		Format:    "json",
		ExitLevel: "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14")

	args := []string{"dockle"}

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

	return container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}), nil
}

// ScanTarball scans a container image tarball
func (m *DockleModule) ScanTarball(ctx context.Context, tarballPath string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		Format:    "json",
		ExitLevel: "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14")

	// Mount tarball file
	if tarballPath != "" {
		container = container.WithMountedFile("/workspace/image.tar", m.client.Host().File(tarballPath))
	}

	args := []string{"dockle"}

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

	return container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}), nil
}

// ScanDockerfile scans a Dockerfile for best practices
func (m *DockleModule) ScanDockerfile(ctx context.Context, dockerfilePath string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		Format:    "json",
		ExitLevel: "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14")

	// Mount Dockerfile
	if dockerfilePath != "" {
		container = container.WithMountedFile("/workspace/Dockerfile", m.client.Host().File(dockerfilePath))
	}

	args := []string{"dockle"}

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

	return container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}), nil
}

// ListChecks lists all available Dockle security checks
func (m *DockleModule) ListChecks(ctx context.Context) (*dagger.Container, error) {
	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14")

	return container.WithExec([]string{"dockle", "--help"}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}), nil
}

// ScanWithPolicy scans using a custom policy file
func (m *DockleModule) ScanWithPolicy(ctx context.Context, imageRef string, policyPath string, opts ...DockleOption) (*dagger.Container, error) {
	config := &DockleConfig{
		Format:    "json",
		ExitLevel: "warn",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14")

	// Mount policy file
	if policyPath != "" {
		container = container.WithMountedFile("/workspace/.dockleignore", m.client.Host().File(policyPath))
	}

	args := []string{"dockle"}

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

	return container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}), nil
}

// ScanImageString scans a container image and returns string output (MCP compatible)
func (m *DockleModule) ScanImageString(ctx context.Context, imageRef string) (string, error) {
	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14").
		WithExec([]string{"dockle", imageRef}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to scan image: no output received")
}

// ScanTarballString scans a container image tarball and returns string output (MCP compatible)
func (m *DockleModule) ScanTarballString(ctx context.Context, tarballPath string) (string, error) {
	tarballFile := m.client.Host().File(tarballPath)
	
	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14").
		WithFile("/workspace/image.tar", tarballFile).
		WithExec([]string{"dockle", "--input", "/workspace/image.tar"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to scan tarball: no output received")
}

// ScanImageJSON scans a container image and returns JSON output (MCP compatible)
func (m *DockleModule) ScanImageJSON(ctx context.Context, imageRef string, outputFile string) (string, error) {
	args := []string{"dockle", "-f", "json"}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}
	args = append(args, imageRef)

	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to scan image with JSON output: no output received")
}

// ScanTarballJSON scans a container image tarball and returns JSON output (MCP compatible)
func (m *DockleModule) ScanTarballJSON(ctx context.Context, tarballPath string, outputFile string) (string, error) {
	tarballFile := m.client.Host().File(tarballPath)
	
	args := []string{"dockle", "-f", "json", "--input", "/workspace/image.tar"}
	
	if outputFile != "" {
		args = append(args, "-o", outputFile)
	}

	container := m.client.Container().
		From("goodwithtech/dockle:v0.4.14").
		WithFile("/workspace/image.tar", tarballFile).
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to scan tarball with JSON output: no output received")
}

type DockleConfig struct {
	Format     string
	Output     string
	ExitLevel  string
	AcceptKey  []string
	AcceptFile []string
	Ignore     []string
}

type DockleOption func(*DockleConfig)

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