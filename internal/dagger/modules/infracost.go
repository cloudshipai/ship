package modules

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

// InfracostModule runs Infracost for cloud cost estimation
type InfracostModule struct {
	client *dagger.Client
	name   string
}

// NewInfracostModule creates a new Infracost module
func NewInfracostModule(client *dagger.Client) *InfracostModule {
	return &InfracostModule{
		client: client,
		name:   "infracost",
	}
}

// BreakdownDirectory generates cost breakdown for a directory
func (m *InfracostModule) BreakdownDirectory(ctx context.Context, dir string) (string, error) {
	container := m.prepareContainer(dir)

	container = container.WithExec([]string{
		"infracost",
		"breakdown",
		"--path", ".",
		"--format", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run infracost breakdown: %w", err)
	}

	return output, nil
}

// BreakdownPlan generates cost breakdown from a Terraform plan
func (m *InfracostModule) BreakdownPlan(ctx context.Context, planFile string) (string, error) {
	dir := filepath.Dir(planFile)
	filename := filepath.Base(planFile)

	container := m.prepareContainer(dir)

	container = container.
		WithFile("/tmp/plan.json", m.client.Host().File(planFile)).
		WithExec([]string{
			"infracost",
			"breakdown",
			"--path", "/tmp/" + filename,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run infracost on plan: %w", err)
	}

	return output, nil
}

// Diff compares costs between current and planned state
func (m *InfracostModule) Diff(ctx context.Context, dir string) (string, error) {
	container := m.prepareContainer(dir)

	container = container.WithExec([]string{
		"infracost",
		"diff",
		"--path", ".",
		"--format", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run infracost diff: %w", err)
	}

	return output, nil
}

// BreakdownWithConfig runs breakdown using a config file
func (m *InfracostModule) BreakdownWithConfig(ctx context.Context, configFile string) (string, error) {
	dir := filepath.Dir(configFile)

	container := m.prepareContainer(dir)

	container = container.
		WithFile("/tmp/infracost.yml", m.client.Host().File(configFile)).
		WithExec([]string{
			"infracost",
			"breakdown",
			"--config-file", "/tmp/infracost.yml",
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run infracost with config: %w", err)
	}

	return output, nil
}

// GenerateHTMLReport generates an HTML cost report
func (m *InfracostModule) GenerateHTMLReport(ctx context.Context, dir string) (string, error) {
	container := m.prepareContainer(dir)

	container = container.WithExec([]string{
		"infracost",
		"breakdown",
		"--path", ".",
		"--format", "html",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML report: %w", err)
	}

	return output, nil
}

// GenerateTableReport generates a table format cost report
func (m *InfracostModule) GenerateTableReport(ctx context.Context, dir string) (string, error) {
	container := m.prepareContainer(dir)

	container = container.WithExec([]string{
		"infracost",
		"breakdown",
		"--path", ".",
		"--format", "table",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate table report: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Infracost
func (m *InfracostModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("infracost/infracost:latest").
		WithExec([]string{"infracost", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get infracost version: %w", err)
	}

	return output, nil
}

// prepareContainer prepares a container with common setup
func (m *InfracostModule) prepareContainer(dir string) *dagger.Container {
	container := m.client.Container().
		From("infracost/infracost:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace")

	// Set Infracost API key if available
	apiKey := os.Getenv("INFRACOST_API_KEY")
	if apiKey != "" {
		container = container.WithEnvVariable("INFRACOST_API_KEY", apiKey)
	}

	// Mount AWS credentials if they exist
	awsCreds := filepath.Join(os.Getenv("HOME"), ".aws")
	if _, err := os.Stat(awsCreds); err == nil {
		container = container.WithDirectory("/root/.aws", m.client.Host().Directory(awsCreds))
	}

	// Mount Azure credentials if they exist
	azureCreds := filepath.Join(os.Getenv("HOME"), ".azure")
	if _, err := os.Stat(azureCreds); err == nil {
		container = container.WithDirectory("/root/.azure", m.client.Host().Directory(azureCreds))
	}

	// Mount GCP credentials if set
	gcpCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gcpCreds != "" {
		if _, err := os.Stat(gcpCreds); err == nil {
			container = container.
				WithFile("/tmp/gcp-creds.json", m.client.Host().File(gcpCreds)).
				WithEnvVariable("GOOGLE_APPLICATION_CREDENTIALS", "/tmp/gcp-creds.json")
		}
	}

	return container
}