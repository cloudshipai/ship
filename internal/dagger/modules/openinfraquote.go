package modules

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// oiqBinary is the path to the oiq binary in the container
const oiqBinary = "/usr/bin/oiq"

// OpenInfraQuoteModule runs OpenInfraQuote for Terraform cost analysis
type OpenInfraQuoteModule struct {
	client *dagger.Client
	name   string
}

// NewOpenInfraQuoteModule creates a new OpenInfraQuote module
func NewOpenInfraQuoteModule(client *dagger.Client) *OpenInfraQuoteModule {
	return &OpenInfraQuoteModule{
		client: client,
		name:   "openinfraquote",
	}
}

// AnalyzePlan analyzes a Terraform plan JSON file for cost estimation
func (m *OpenInfraQuoteModule) AnalyzePlan(ctx context.Context, planFile string, region string) (string, error) {
	// Get the directory containing the plan file
	dir := filepath.Dir(planFile)
	filename := filepath.Base(planFile)

	// First, use a utility container to download the pricing sheet
	utilContainer := m.client.Container().
		From("alpine:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		// Download the pricing sheet
		WithExec([]string{"sh", "-c", "apk add --no-cache curl && curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > prices.csv"})

	// Get the workspace with pricing sheet
	_, err := utilContainer.Directory("/workspace").Export(ctx, dir+"_temp")
	if err != nil {
		return "", fmt.Errorf("failed to prepare workspace: %w", err)
	}
	defer func() {
		// Clean up temp directory
		m.client.Host().Directory(dir + "_temp")
	}()

	// First, get the match output
	matchContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir+"_temp")).
		WithWorkdir("/workspace").
		WithExec([]string{oiqBinary, "match", "--pricesheet", "prices.csv", filename})

	matchOutput, err := matchContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}

	// Now run oiq price with the match output as stdin using Dagger's stdin parameter
	priceContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{oiqBinary, "price", "--region", region, "--format", "json"}, dagger.ContainerWithExecOpts{
			Stdin: matchOutput,
		})

	priceOutput, err := priceContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq price: %w", err)
	}

	return priceOutput, nil
}

// AnalyzeDirectory analyzes all Terraform files in a directory
func (m *OpenInfraQuoteModule) AnalyzeDirectory(ctx context.Context, dir string, region string) (string, error) {
	// First, use a utility container to download the pricing sheet and generate Terraform plan
	terraformContainer := m.client.Container().
		From("hashicorp/terraform:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		// Download the pricing sheet
		WithExec([]string{"sh", "-c", "apk add --no-cache curl && curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > prices.csv"}).
		// Generate Terraform plan as JSON
		WithExec([]string{"sh", "-c", "terraform init && terraform plan -out=tf.plan && terraform show -json tf.plan > tfplan.json"})

	// Get the workspace with generated files
	_, err := terraformContainer.Directory("/workspace").Export(ctx, dir+"_temp")
	if err != nil {
		return "", fmt.Errorf("failed to prepare workspace: %w", err)
	}
	defer func() {
		// Clean up temp directory
		m.client.Host().Directory(dir + "_temp")
	}()

	// First, get the match output
	matchContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir+"_temp")).
		WithWorkdir("/workspace").
		WithExec([]string{oiqBinary, "match", "--pricesheet", "prices.csv", "tfplan.json"})

	matchOutput, err := matchContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}

	// Now run oiq price with the match output as stdin using Dagger's stdin parameter
	priceContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{oiqBinary, "price", "--region", region, "--format", "json"}, dagger.ContainerWithExecOpts{
			Stdin: matchOutput,
		})

	priceOutput, err := priceContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq price: %w", err)
	}

	return priceOutput, nil
}

// GetVersion returns the version of OpenInfraQuote
func (m *OpenInfraQuoteModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{oiqBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get openinfraquote version: %w", err)
	}

	return output, nil
}

// CompareRegions compares costs across multiple regions
func (m *OpenInfraQuoteModule) CompareRegions(ctx context.Context, planFile string, regions []string) (string, error) {
	dir := filepath.Dir(planFile)
	filename := filepath.Base(planFile)

	// Prepare workspace with pricing sheet
	utilContainer := m.client.Container().
		From("alpine:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"sh", "-c", "apk add --no-cache curl && curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > prices.csv"})

	_, err := utilContainer.Directory("/workspace").Export(ctx, dir+"_temp")
	if err != nil {
		return "", fmt.Errorf("failed to prepare workspace: %w", err)
	}
	defer func() {
		m.client.Host().Directory(dir + "_temp")
	}()

	// Get match output once
	matchContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir+"_temp")).
		WithWorkdir("/workspace").
		WithExec([]string{oiqBinary, "match", "--pricesheet", "prices.csv", filename})

	matchOutput, err := matchContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}

	// Compare costs across regions
	results := "{"
	for i, region := range regions {
		priceContainer := m.client.Container().
			From("ghcr.io/terrateamio/openinfraquote:latest").
			WithExec([]string{oiqBinary, "price", "--region", region, "--format", "json"}, dagger.ContainerWithExecOpts{
				Stdin: matchOutput,
			})

		regionOutput, err := priceContainer.Stdout(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get pricing for region %s: %w", region, err)
		}

		if i > 0 {
			results += ","
		}
		results += fmt.Sprintf(`"%s": %s`, region, regionOutput)
	}
	results += "}"

	return results, nil
}

// Match matches Terraform resources to pricing
func (m *OpenInfraQuoteModule) Match(ctx context.Context, pricesheet string, tfplanJson string) (string, error) {
	// Get directory containing files
	dir := filepath.Dir(tfplanJson)
	planFilename := filepath.Base(tfplanJson)
	pricesheetFilename := filepath.Base(pricesheet)

	container := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{oiqBinary, "match", "--pricesheet", pricesheetFilename, planFilename})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}

	return output, nil
}

// Price calculates prices from matched resources
func (m *OpenInfraQuoteModule) Price(ctx context.Context, region string, inputFile string) (string, error) {
	args := []string{oiqBinary, "price"}
	if region != "" {
		args = append(args, "--region", region)
	}
	
	var container *dagger.Container
	if inputFile != "" {
		// Use file input
		dir := filepath.Dir(inputFile)
		filename := filepath.Base(inputFile)
		container = m.client.Container().
			From("ghcr.io/terrateamio/openinfraquote:latest").
			WithDirectory("/workspace", m.client.Host().Directory(dir)).
			WithWorkdir("/workspace")
		args = append(args, filename)
	} else {
		// Use stdin (for pipeline)
		container = m.client.Container().
			From("ghcr.io/terrateamio/openinfraquote:latest")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq price: %w", err)
	}

	return output, nil
}

// DownloadPrices downloads AWS pricing data
func (m *OpenInfraQuoteModule) DownloadPrices(ctx context.Context, outputFile string) (string, error) {
	if outputFile == "" {
		outputFile = "prices.csv"
	}

	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > " + outputFile})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to download prices: %w", err)
	}

	return output, nil
}

// GetHelp returns help information for oiq
func (m *OpenInfraQuoteModule) GetHelp(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{oiqBinary, "--help"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get oiq help: %w", err)
	}

	return output, nil
}

// FullPipeline runs the complete cost estimation pipeline (match | price)
func (m *OpenInfraQuoteModule) FullPipeline(ctx context.Context, tfplanJson string, pricesheet string, region string) (string, error) {
	// Get directory containing files
	dir := filepath.Dir(tfplanJson)
	planFilename := filepath.Base(tfplanJson)
	pricesheetFilename := filepath.Base(pricesheet)

	// First, run match to get the output
	matchContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{oiqBinary, "match", "--pricesheet", pricesheetFilename, planFilename})

	matchOutput, err := matchContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}

	// Then pipe to price with region
	priceContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{oiqBinary, "price", "--region", region}, dagger.ContainerWithExecOpts{
			Stdin: matchOutput,
		})

	priceOutput, err := priceContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq price: %w", err)
	}

	return priceOutput, nil
}
