package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TrivyGoldenModule runs enhanced Trivy workflows for golden image scanning
type TrivyGoldenModule struct {
	client *dagger.Client
	name   string
}

const trivyGoldenBinary = "/usr/local/bin/trivy"

// NewTrivyGoldenModule creates a new Trivy Golden module
func NewTrivyGoldenModule(client *dagger.Client) *TrivyGoldenModule {
	return &TrivyGoldenModule{
		client: client,
		name:   "trivy-golden",
	}
}

// ScanGoldenImage performs comprehensive golden image scanning
func (m *TrivyGoldenModule) ScanGoldenImage(ctx context.Context, imageName string, maxCritical int, maxHigh int) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			trivyGoldenBinary,
			"image",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			"--scanners", "vuln,secret,config",
			"--exit-code", "1",
			imageName,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Check if it's a vulnerability failure (expected for golden image validation)
		stderr, _ := container.Stderr(ctx)
		return fmt.Sprintf(`{"scan_result": %s, "stderr": "%s", "validation": "failed"}`, output, stderr), nil
	}

	return fmt.Sprintf(`{"scan_result": %s, "validation": "passed"}`, output), nil
}

// CompareImages compares two images for golden image validation
func (m *TrivyGoldenModule) CompareImages(ctx context.Context, baseImage string, candidateImage string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(`
				trivy image --format json %s > /base.json
				trivy image --format json %s > /candidate.json
				echo '{"base_scan": '$(cat /base.json)', "candidate_scan": '$(cat /candidate.json)'}'
			`, baseImage, candidateImage),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to compare images: %w", err)
	}

	return output, nil
}

// ValidateImagePolicy validates image against policy
func (m *TrivyGoldenModule) ValidateImagePolicy(ctx context.Context, imageName string, policyPath string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithFile("/policy.rego", m.client.Host().File(policyPath)).
		WithExec([]string{
			trivyGoldenBinary,
			"image",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			"--ignore-policy", "/policy.rego",
			imageName,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate image policy: %w", err)
	}

	return output, nil
}

// GenerateImageAttestation generates SLSA attestation for image
func (m *TrivyGoldenModule) GenerateImageAttestation(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			trivyGoldenBinary,
			"image",
			"--format", "cosign-vuln",
			imageName,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate attestation: %w", err)
	}

	return output, nil
}

// ScanImageBasic performs basic image scanning with customizable parameters
func (m *TrivyGoldenModule) ScanImageBasic(ctx context.Context, imageName string, severity string, format string, scanners string, exitCode bool, outputFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	args := []string{trivyGoldenBinary, "image"}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if format != "" {
		args = append(args, "--format", format)
	}
	if scanners != "" {
		args = append(args, "--scanners", scanners)
	}
	if exitCode {
		args = append(args, "--exit-code", "1")
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	args = append(args, imageName)

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan image: %w", err)
	}

	return output, nil
}

// ScanFilesystemBasic performs basic filesystem scanning with customizable parameters
func (m *TrivyGoldenModule) ScanFilesystemBasic(ctx context.Context, path string, scanners string, severity string, format string, outputFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path)).
		WithWorkdir("/workspace")

	args := []string{trivyGoldenBinary, "fs"}
	if scanners != "" {
		args = append(args, "--scanners", scanners)
	}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if format != "" {
		args = append(args, "--format", format)
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	args = append(args, ".")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan filesystem: %w", err)
	}

	return output, nil
}

// ScanConfigBasic performs basic configuration scanning with customizable parameters
func (m *TrivyGoldenModule) ScanConfigBasic(ctx context.Context, path string, severity string, format string, policy string, outputFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path)).
		WithWorkdir("/workspace")

	if policy != "" {
		container = container.WithFile("/policy.rego", m.client.Host().File(policy))
	}

	args := []string{trivyGoldenBinary, "config"}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if format != "" {
		args = append(args, "--format", format)
	}
	if policy != "" {
		args = append(args, "--policy", "/policy.rego")
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	args = append(args, ".")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan config: %w", err)
	}

	return output, nil
}

// GenerateSBOMBasic generates SBOM with customizable parameters
func (m *TrivyGoldenModule) GenerateSBOMBasic(ctx context.Context, imageName string, format string, outputFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	args := []string{trivyGoldenBinary, "image", "--format", format}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	args = append(args, imageName)

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM: %w", err)
	}

	return output, nil
}

// ScanSecretsBasic performs secret scanning with customizable parameters
func (m *TrivyGoldenModule) ScanSecretsBasic(ctx context.Context, target string, targetType string, format string, outputFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	if targetType == "fs" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}

	args := []string{trivyGoldenBinary, targetType, "--scanners", "secret"}
	if format != "" {
		args = append(args, "--format", format)
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	args = append(args, target)

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan secrets: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Trivy
func (m *TrivyGoldenModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{trivyGoldenBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get trivy version: %w", err)
	}

	return output, nil
}
