package dagger

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/cloudship/ship/internal/dagger/modules"
)

// Engine manages the Dagger client and provides methods to run Dagger pipelines
type Engine struct {
	client *dagger.Client
	ctx    context.Context
}

// NewEngine creates a new Dagger engine instance
func NewEngine(ctx context.Context) (*Engine, error) {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
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

// RunSteampipeQuery runs a Steampipe query in a container
func (e *Engine) RunSteampipeQuery(provider, query string) (string, error) {
	// For now, use the official Steampipe image
	// In production, we'd use a custom image with pre-installed plugins
	container := e.client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "query", query, "--output", "json"})

	output, err := container.Stdout(e.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run steampipe query: %w", err)
	}

	return output, nil
}

// NewOpenInfraQuoteModule creates a new OpenInfraQuote module
func (e *Engine) NewOpenInfraQuoteModule() *modules.OpenInfraQuoteModule {
	return modules.NewOpenInfraQuoteModule(e.client)
}

// NewInfraScanModule creates a new InfraScan module
func (e *Engine) NewInfraScanModule() *modules.InfraScanModule {
	return modules.NewInfraScanModule(e.client)
}

// NewTerraformDocsModule creates a new terraform-docs module
func (e *Engine) NewTerraformDocsModule() *modules.TerraformDocsModule {
	return modules.NewTerraformDocsModule(e.client)
}

// NewTFLintModule creates a new TFLint module
func (e *Engine) NewTFLintModule() *modules.TFLintModule {
	return modules.NewTFLintModule(e.client)
}

// NewCheckovModule creates a new Checkov module
func (e *Engine) NewCheckovModule() *modules.CheckovModule {
	return modules.NewCheckovModule(e.client)
}

// NewInfracostModule creates a new Infracost module
func (e *Engine) NewInfracostModule() *modules.InfracostModule {
	return modules.NewInfracostModule(e.client)
}

// NewSteampipeModule creates a new Steampipe module
func (e *Engine) NewSteampipeModule() *modules.SteampipeModule {
	return modules.NewSteampipeModule(e.client)
}

// NewLLMModule creates a new LLM module
func (e *Engine) NewLLMModule(provider, model string) *modules.LLMModule {
	return modules.NewLLMModule(e.client, provider, model)
}

// NewInfraMapModule creates a new InfraMap module
func (e *Engine) NewInfraMapModule() *modules.InfraMapModule {
	return modules.NewInfraMapModule(e.client)
}
