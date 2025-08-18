package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// CloudsplainingModule runs Cloudsplaining for AWS IAM security assessment
type CloudsplainingModule struct {
	client *dagger.Client
	name   string
}

const cloudsplainingBinary = "/usr/local/bin/cloudsplaining"

// NewCloudsplainingModule creates a new Cloudsplaining module
func NewCloudsplainingModule(client *dagger.Client) *CloudsplainingModule {
	return &CloudsplainingModule{
		client: client,
		name:   "cloudsplaining",
	}
}

// ScanAccountAuthorization scans account authorization details
func (m *CloudsplainingModule) ScanAccountAuthorization(ctx context.Context, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	// First download, then scan
	container = container.WithExec([]string{
		"sh", "-c", cloudsplainingBinary + " download --profile " + profile + " && " + cloudsplainingBinary + " scan --input-file default.json --output /workspace/results.json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan account authorization: %w", err)
	}

	return output, nil
}

// ScanPolicyFile scans a specific IAM policy file
func (m *CloudsplainingModule) ScanPolicyFile(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			cloudsplainingBinary, "scan-policy-file",
			"--input-file", "policy.json",
			"--output", "results.json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan policy file: %w", err)
	}

	return output, nil
}

// CreateReportFromResults creates an HTML report from scan results
func (m *CloudsplainingModule) CreateReportFromResults(ctx context.Context, resultsPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithFile("/workspace/results.json", m.client.Host().File(resultsPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			cloudsplainingBinary, "create-report",
			"--input-file", "results.json",
			"--output", "report.html",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to create report: %w", err)
	}

	return output, nil
}

// ScanWithMinimization scans with policy minimization recommendations
func (m *CloudsplainingModule) ScanWithMinimization(ctx context.Context, profile string, minimizeStatementId string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	// Build the command with optional parameters
	cmd := cloudsplainingBinary + " download --profile " + profile + " && " + cloudsplainingBinary + " scan --input-file default.json --output /workspace/results.json"
	if minimizeStatementId != "" {
		cmd = cloudsplainingBinary + " download --profile " + profile + " && " + cloudsplainingBinary + " scan --input-file default.json --minimize-statement-id " + minimizeStatementId + " --output /workspace/results.json"
	}

	container = container.WithExec([]string{"sh", "-c", cmd})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan with minimization: %w", err)
	}

	return output, nil
}

// Download downloads AWS account authorization data
func (m *CloudsplainingModule) Download(ctx context.Context, profile string, includeNonDefaultPolicyVersions bool) (string, error) {
	args := []string{cloudsplainingBinary, "download"}
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	if includeNonDefaultPolicyVersions {
		args = append(args, "--include-non-default-policy-versions")
	}

	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest")

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to download account data: %w", err)
	}

	return output, nil
}

// ScanAccountData scans downloaded account authorization data
func (m *CloudsplainingModule) ScanAccountData(ctx context.Context, inputFile string, exclusionsFile string, outputDir string) (string, error) {
	args := []string{cloudsplainingBinary, "scan", "--input-file", "/workspace/input.json"}
	if exclusionsFile != "" {
		args = append(args, "--exclusions-file", "/workspace/exclusions.yml")
	}
	if outputDir != "" {
		args = append(args, "--output", outputDir)
	}

	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithFile("/workspace/input.json", m.client.Host().File(inputFile)).
		WithWorkdir("/workspace")

	if exclusionsFile != "" {
		container = container.WithFile("/workspace/exclusions.yml", m.client.Host().File(exclusionsFile))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan account data: %w", err)
	}

	return output, nil
}

// CreateExclusionsFile creates exclusions file template
func (m *CloudsplainingModule) CreateExclusionsFile(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithExec([]string{cloudsplainingBinary, "create-exclusions-file"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to create exclusions file: %w", err)
	}

	return output, nil
}

// CreateMultiAccountConfig creates multi-account configuration file
func (m *CloudsplainingModule) CreateMultiAccountConfig(ctx context.Context, outputFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithExec([]string{cloudsplainingBinary, "create-multi-account-config-file", "-o", "/workspace/config.yml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to create multi-account config: %w", err)
	}

	return output, nil
}

// ScanMultiAccount scans multiple AWS accounts
func (m *CloudsplainingModule) ScanMultiAccount(ctx context.Context, configFile string, profile string, roleName string, outputBucket string, outputDirectory string) (string, error) {
	args := []string{cloudsplainingBinary, "scan-multi-account", "-c", "/workspace/config.yml"}
	
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	if roleName != "" {
		args = append(args, "--role-name", roleName)
	}
	if outputBucket != "" {
		args = append(args, "--output-bucket", outputBucket)
	}
	if outputDirectory != "" {
		args = append(args, "--output-directory", outputDirectory)
	}

	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithFile("/workspace/config.yml", m.client.Host().File(configFile)).
		WithWorkdir("/workspace")

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan multi-account: %w", err)
	}

	return output, nil
}