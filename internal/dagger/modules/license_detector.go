package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// LicenseDetectorModule detects and analyzes software licenses
type LicenseDetectorModule struct {
	client *dagger.Client
	name   string
}

// NewLicenseDetectorModule creates a new license detector module
func NewLicenseDetectorModule(client *dagger.Client) *LicenseDetectorModule {
	return &LicenseDetectorModule{
		client: client,
		name:   "license-detector",
	}
}

// DetectLicenses detects licenses in a directory
func (m *LicenseDetectorModule) DetectLicenses(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("licensee/licensee:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"licensee", "detect",
			"--json",
			".",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to detect licenses: %w", err)
	}

	return output, nil
}

// AnalyzeDependencyLicenses analyzes dependency licenses
func (m *LicenseDetectorModule) AnalyzeDependencyLicenses(ctx context.Context, packageFile string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "npm", "jq"}).
		WithFile("/package.json", m.client.Host().File(packageFile)).
		WithExec([]string{
			"sh", "-c",
			`npm install license-checker-rseidelsohn && npx license-checker-rseidelsohn --json --out /licenses.json && cat /licenses.json`,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to analyze dependency licenses: %w", err)
	}

	return output, nil
}

// ValidateLicenseCompliance validates license compliance
func (m *LicenseDetectorModule) ValidateLicenseCompliance(ctx context.Context, dir string, allowedLicenses []string) (string, error) {
	container := m.client.Container().
		From("licensee/licensee:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace")

	// Create allowed licenses file
	allowedList := ""
	for _, license := range allowedLicenses {
		allowedList += license + "\n"
	}

	container = container.
		WithNewFile("/allowed_licenses.txt", allowedList).
		WithExec([]string{
			"sh", "-c",
			`
				detected=$(licensee detect --json . | jq -r '.licenses[]?.spdx_id // empty')
				allowed=$(cat /allowed_licenses.txt)
				echo "Detected licenses: $detected"
				echo "Allowed licenses: $allowed"
				for license in $detected; do
					if ! echo "$allowed" | grep -q "$license"; then
						echo "ERROR: License $license not in allowed list"
						exit 1
					fi
				done
				echo "All licenses are compliant"
			`,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate license compliance: %w", err)
	}

	return output, nil
}
