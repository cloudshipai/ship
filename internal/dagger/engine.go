package dagger

import (
	"context"
	"fmt"
	"io"
	"os"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
)

// Engine manages the Dagger client and provides methods to run Dagger pipelines
type Engine struct {
	client *dagger.Client
	ctx    context.Context
}

// NewEngine creates a new Dagger engine instance
func NewEngine(ctx context.Context, logLevel ...string) (*Engine, error) {
	var logOutput io.Writer = io.Discard
	if len(logLevel) > 0 && logLevel[0] == "debug" {
		logOutput = os.Stderr
	}

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logOutput))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dagger: %w", err)
	}

	return &Engine{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close closes the Dagger client connection
func (e *Engine) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

// GetClient returns the underlying Dagger client
func (e *Engine) GetClient() *dagger.Client {
	return e.client
}

// RunContainer runs a container with the specified image and returns the output
func (e *Engine) RunContainer(image string, args []string) (string, error) {
	container := e.client.Container().
		From(image).
		WithExec(args)

	output, err := container.Stdout(e.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run container: %w", err)
	}

	return output, nil
}

// BuildContainer builds a container from a directory with a Dockerfile
func (e *Engine) BuildContainer(contextDir string, dockerfile string) (*dagger.Container, error) {
	// Get the build context directory
	src := e.client.Host().Directory(contextDir)

	// Build the container
	container := e.client.Container().Build(src, dagger.ContainerBuildOpts{
		Dockerfile: dockerfile,
	})

	return container, nil
}

// NewOpenInfraQuoteModule creates a new OpenInfraQuote module
func (e *Engine) NewOpenInfraQuoteModule() *modules.OpenInfraQuoteModule {
	return modules.NewOpenInfraQuoteModule(e.client)
}

// NewTerraformDocsModule creates a new terraform-docs module
func (e *Engine) NewTerraformDocsModule() *modules.TerraformDocsModule {
	return modules.NewTerraformDocsModule(e.client)
}

// NewTFLintModule creates a new TFLint module
func (e *Engine) NewTFLintModule() *modules.TFLintModule {
	return modules.NewTFLintModule(e.client)
}

// NewGitSecretsModule creates a new git-secrets module
func (e *Engine) NewGitSecretsModule() *modules.GitSecretsModule {
	return modules.NewGitSecretsModule(e.client)
}

// NewGitleaksModule creates a new Gitleaks module
func (e *Engine) NewGitleaksModule() *modules.GitleaksModule {
	return modules.NewGitleaksModule(e.client)
}

// NewInfraScanModule creates a new InfraScan module
func (e *Engine) NewInfraScanModule() *modules.InfraScanModule {
	return modules.NewInfraScanModule(e.client)
}

// NewCheckovModule creates a new Checkov module
func (e *Engine) NewCheckovModule() *modules.CheckovModule {
	return modules.NewCheckovModule(e.client)
}


// NewLLMModule has been removed - use the new Eino agent system instead
// See internal/agent package for the new AI-powered investigation capabilities

// NewInfraMapModule creates a new InfraMap module
func (e *Engine) NewInfraMapModule() *modules.InfraMapModule {
	return modules.NewInfraMapModule(e.client)
}

// OpenCode creates a new OpenCode module
func (e *Engine) OpenCode() *modules.OpenCodeModule {
	return modules.NewOpenCodeModule(e.client)
}
