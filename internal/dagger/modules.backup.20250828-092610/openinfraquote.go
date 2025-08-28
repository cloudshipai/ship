package modules

import (
	"context"

	"dagger.io/dagger"
)

// OpenInfraQuoteModule runs OpenInfraQuote for infrastructure cost estimation
type OpenInfraQuoteModule struct {
	client *dagger.Client
	name   string
}

// OpenInfraQuoteOptions contains options for cost estimation
type OpenInfraQuoteOptions struct {
	OutputFormat      string
	OutputFile        string
	TerraformPlanFile string
	Currency          string
	Region            string
	ShowSkipped       bool
	SyncUsageFile     bool
	UsageFile         string
}

// OpenInfraQuoteDiffOptions contains options for cost diff
type OpenInfraQuoteDiffOptions struct {
	OutputFormat      string
	Currency          string
	ShowAllProjects   bool
}

// NewOpenInfraQuoteModule creates a new OpenInfraQuote module
func NewOpenInfraQuoteModule(client *dagger.Client) *OpenInfraQuoteModule {
	return &OpenInfraQuoteModule{
		client: client,
		name:   "openinfraquote",
	}
}

// Estimate generates cost estimates for Terraform infrastructure
func (m *OpenInfraQuoteModule) Estimate(ctx context.Context, terraformPath string, opts OpenInfraQuoteOptions) (string, error) {
	args := []string{"infracost", "breakdown", "--path", "."}

	// Add estimation options
	if opts.OutputFormat != "" && opts.OutputFormat != "table" {
		args = append(args, "--format", opts.OutputFormat)
	}
	if opts.TerraformPlanFile != "" {
		args = append(args, "--terraform-plan-file", opts.TerraformPlanFile)
	}
	if opts.ShowSkipped {
		args = append(args, "--show-skipped")
	}
	if opts.SyncUsageFile {
		args = append(args, "--sync-usage-file")
	}
	if opts.UsageFile != "" {
		args = append(args, "--usage-file", opts.UsageFile)
	}

	// Set environment variables
	envs := []string{
		"INFRACOST_API_KEY=your-api-key", // User needs to set this
	}
	if opts.Currency != "" {
		envs = append(envs, "INFRACOST_CURRENCY="+opts.Currency)
	}

	container := m.client.Container().
		From(getImageTag("openinfraquote", "infracost/infracost:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(terraformPath)).
		WithWorkdir("/workspace")

	// Add environment variables
	for _, env := range envs {
		container = container.WithEnvVariable(env[:len(env)-len("=your-api-key")], env[len(env)-len("your-api-key"):])
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Cost estimation completed", nil
}

// Diff shows cost difference between two Terraform configurations
func (m *OpenInfraQuoteModule) Diff(ctx context.Context, path1, path2 string, opts OpenInfraQuoteDiffOptions) (string, error) {
	args := []string{"infracost", "diff", "--path1", "/path1", "--path2", "/path2"}

	// Add diff options
	if opts.OutputFormat != "" && opts.OutputFormat != "table" {
		args = append(args, "--format", opts.OutputFormat)
	}
	if opts.ShowAllProjects {
		args = append(args, "--show-all-projects")
	}

	// Set environment variables
	envs := []string{
		"INFRACOST_API_KEY=your-api-key", // User needs to set this
	}
	if opts.Currency != "" {
		envs = append(envs, "INFRACOST_CURRENCY="+opts.Currency)
	}

	container := m.client.Container().
		From(getImageTag("openinfraquote", "infracost/infracost:latest")).
		WithDirectory("/path1", m.client.Host().Directory(path1)).
		WithDirectory("/path2", m.client.Host().Directory(path2)).
		WithWorkdir("/")

	// Add environment variables
	for _, env := range envs {
		container = container.WithEnvVariable(env[:len(env)-len("=your-api-key")], env[len(env)-len("your-api-key"):])
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Cost diff completed", nil
}