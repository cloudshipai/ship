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
			"trivy",
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
			"trivy",
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
			"trivy",
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
