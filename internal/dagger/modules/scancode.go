package modules

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// ScanCodeModule runs ScanCode Toolkit for license detection
type ScanCodeModule struct {
	client *dagger.Client
	name   string
}

// NewScanCodeModule creates a new ScanCode module
func NewScanCodeModule(client *dagger.Client) *ScanCodeModule {
	return &ScanCodeModule{
		client: client,
		name:   "scancode",
	}
}

// LicenseScan performs high-accuracy license detection
func (m *ScanCodeModule) LicenseScan(ctx context.Context, path string, outputPath string, extraFlags []string) (string, error) {
	// Default output path
	if outputPath == "" {
		outputPath = "./scancode.json"
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Build scancode command
	args := []string{"scancode", "-l", "--json-pp", "/work/scancode.json", "/work/src"}
	
	// Add extra flags if provided
	if len(extraFlags) > 0 {
		args = append(args, extraFlags...)
	}

	// Run ScanCode in container
	// Using the official ScanCode Toolkit image
	container := m.client.Container().
		From("beevelop/scancode:latest").
		WithDirectory("/work/src", m.client.Host().Directory(absPath)).
		WithWorkdir("/work").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	// Get output
	output, err := container.Stdout(ctx)
	if err != nil {
		// Even on error, try to get stderr for diagnostics
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("scancode failed: %w\nstderr: %s", err, stderr)
	}

	// Export the JSON file
	_, err = container.File("/work/scancode.json").Export(ctx, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to export scancode results: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of ScanCode
func (m *ScanCodeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("beevelop/scancode:latest").
		WithExec([]string{"scancode", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scancode version: %w", err)
	}

	return output, nil
}

// GetHelp returns help information for ScanCode
func (m *ScanCodeModule) GetHelp(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("beevelop/scancode:latest").
		WithExec([]string{"scancode", "--help"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scancode help: %w", err)
	}

	return output, nil
}