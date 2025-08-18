package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// Binary paths for license detection tools
const (
	licenseeBinary         = "/usr/local/bin/licensee"
	askalonoBinary         = "/usr/local/bin/askalono"
	goLicenseDetectorBinary = "/usr/local/bin/license-detector"
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

// DetectLicenses detects licenses in a directory using multiple tools for comprehensive analysis
func (m *LicenseDetectorModule) DetectLicenses(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("ruby:alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "build-base"}).
		WithExec([]string{"gem", "install", "licensee"}).
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			licenseeBinary, "detect",
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

// AskalonoIdentify identifies license in a file using askalono
func (m *LicenseDetectorModule) AskalonoIdentify(ctx context.Context, filePath string, optimize bool) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/amzn/askalono/releases/latest/download/askalono-Linux.tar.gz | tar xz && mv askalono /usr/local/bin/"}).
		WithFile("/workspace/license.txt", m.client.Host().File(filePath))

	args := []string{askalonoBinary, "id", "/workspace/license.txt"}
	if optimize {
		args = []string{askalonoBinary, "id", "--optimize", "/workspace/license.txt"}
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to identify license with askalono: %w", err)
	}

	return output, nil
}

// AskalonoCrawl crawls directory for license files using askalono
func (m *LicenseDetectorModule) AskalonoCrawl(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/amzn/askalono/releases/latest/download/askalono-Linux.tar.gz | tar xz && mv askalono /usr/local/bin/"}).
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{askalonoBinary, "crawl", "."})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to crawl for licenses with askalono: %w", err)
	}

	return output, nil
}

// LicenseScannerFile scans specific file using CycloneDX license-scanner
func (m *LicenseDetectorModule) LicenseScannerFile(ctx context.Context, filePath string, showCopyrights bool, showHash bool, showKeywords bool, debug bool) (string, error) {
	container := m.client.Container().
		From("python:alpine").
		WithExec([]string{"pip", "install", "cyclone-scanner"}).
		WithFile("/workspace/file.txt", m.client.Host().File(filePath))

	args := []string{"license-scanner", "--file", "/workspace/file.txt"}
	if showCopyrights {
		args = append(args, "--copyrights")
	}
	if showHash {
		args = append(args, "--hash")
	}
	if showKeywords {
		args = append(args, "--keywords")
	}
	if debug {
		args = append(args, "--debug")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan file with license-scanner: %w", err)
	}

	return output, nil
}

// LicenseScannerDirectory scans directory using CycloneDX license-scanner
func (m *LicenseDetectorModule) LicenseScannerDirectory(ctx context.Context, dir string, showCopyrights bool, showHash bool, quiet bool) (string, error) {
	container := m.client.Container().
		From("python:alpine").
		WithExec([]string{"pip", "install", "cyclone-scanner"}).
		WithDirectory("/workspace", m.client.Host().Directory(dir))

	args := []string{"license-scanner", "--dir", "/workspace"}
	if showCopyrights {
		args = append(args, "--copyrights")
	}
	if showHash {
		args = append(args, "--hash")
	}
	if quiet {
		args = append(args, "--quiet")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan directory with license-scanner: %w", err)
	}

	return output, nil
}

// LicenseFinderReport generates license report using LicenseFinder
func (m *LicenseDetectorModule) LicenseFinderReport(ctx context.Context, projectPath string, format string) (string, error) {
	container := m.client.Container().
		From("licensefinder/license_finder:latest")

	if projectPath != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(projectPath)).WithWorkdir("/workspace")
	}

	args := []string{"license_finder", "report"}
	if format != "" {
		args = append(args, "--format", format)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate license finder report: %w", err)
	}

	return output, nil
}

// GoLicenseDetector detects project license using go-license-detector
func (m *LicenseDetectorModule) GoLicenseDetector(ctx context.Context, projectPath string) (string, error) {
	container := m.client.Container().
		From("golang:alpine").
		WithExec([]string{"go", "install", "github.com/go-enry/go-license-detector/v4/cmd/license-detector@latest"})

	if projectPath != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(projectPath))
		container = container.WithExec([]string{"license-detector", "/workspace"})
	} else {
		container = container.WithExec([]string{"license-detector", "."})
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to detect license with go-license-detector: %w", err)
	}

	return output, nil
}
