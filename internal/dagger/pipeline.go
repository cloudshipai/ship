package dagger

import (
	"fmt"
	"os"
	"path/filepath"
)

// Pipeline represents a Dagger pipeline for running investigations
type Pipeline struct {
	engine *Engine
	name   string
}

// NewPipeline creates a new pipeline with the given name
func NewPipeline(engine *Engine, name string) *Pipeline {
	return &Pipeline{
		engine: engine,
		name:   name,
	}
}

// InvestigateAWS runs an AWS investigation pipeline
func (p *Pipeline) InvestigateAWS(queries []string) (map[string]string, error) {
	results := make(map[string]string)
	client := p.engine.GetClient()

	// Create a container with AWS CLI and Steampipe
	container := client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "plugin", "install", "aws"})

	// Mount AWS credentials from host
	awsCredsPath := filepath.Join(os.Getenv("HOME"), ".aws")
	if _, err := os.Stat(awsCredsPath); err == nil {
		awsCreds := client.Host().Directory(awsCredsPath)
		container = container.WithDirectory("/root/.aws", awsCreds)
	}

	// Run each query
	for i, query := range queries {
		queryContainer := container.WithExec([]string{
			"steampipe", "query", query, "--output", "json",
		})

		output, err := queryContainer.Stdout(p.engine.ctx)
		if err != nil {
			results[fmt.Sprintf("query_%d_error", i)] = err.Error()
			continue
		}

		results[fmt.Sprintf("query_%d", i)] = output
	}

	return results, nil
}

// InvestigateCloudflare runs a Cloudflare investigation pipeline
func (p *Pipeline) InvestigateCloudflare(queries []string) (map[string]string, error) {
	results := make(map[string]string)
	client := p.engine.GetClient()

	// Create a container with Cloudflare CLI and Steampipe
	container := client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "plugin", "install", "cloudflare"})

	// Set Cloudflare API token if available
	if token := os.Getenv("CLOUDFLARE_API_TOKEN"); token != "" {
		container = container.WithEnvVariable("CLOUDFLARE_API_TOKEN", token)
	}

	// Run each query
	for i, query := range queries {
		queryContainer := container.WithExec([]string{
			"steampipe", "query", query, "--output", "json",
		})

		output, err := queryContainer.Stdout(p.engine.ctx)
		if err != nil {
			results[fmt.Sprintf("query_%d_error", i)] = err.Error()
			continue
		}

		results[fmt.Sprintf("query_%d", i)] = output
	}

	return results, nil
}

// InvestigateHeroku runs a Heroku investigation pipeline
func (p *Pipeline) InvestigateHeroku(queries []string) (map[string]string, error) {
	results := make(map[string]string)
	client := p.engine.GetClient()

	// Create a container with Heroku CLI and Steampipe
	container := client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "plugin", "install", "heroku"})

	// Set Heroku API key if available
	if apiKey := os.Getenv("HEROKU_API_KEY"); apiKey != "" {
		container = container.WithEnvVariable("HEROKU_API_KEY", apiKey)
	}

	// Run each query
	for i, query := range queries {
		queryContainer := container.WithExec([]string{
			"steampipe", "query", query, "--output", "json",
		})

		output, err := queryContainer.Stdout(p.engine.ctx)
		if err != nil {
			results[fmt.Sprintf("query_%d_error", i)] = err.Error()
			continue
		}

		results[fmt.Sprintf("query_%d", i)] = output
	}

	return results, nil
}
