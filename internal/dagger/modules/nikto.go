package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// NiktoModule runs Nikto for web vulnerability scanning
type NiktoModule struct {
	client *dagger.Client
	name   string
}

// NewNiktoModule creates a new Nikto module
func NewNiktoModule(client *dagger.Client) *NiktoModule {
	return &NiktoModule{
		client: client,
		name:   "nikto",
	}
}

// ScanHost scans a web host for vulnerabilities
func (m *NiktoModule) ScanHost(ctx context.Context, host string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{
			"nikto.pl",
			"-h", host,
			"-Format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run nikto scan: %w", err)
	}

	return output, nil
}

// ScanWithSSL scans a host with SSL/TLS analysis
func (m *NiktoModule) ScanWithSSL(ctx context.Context, host string, port int) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{
			"nikto.pl",
			"-h", host,
			"-p", fmt.Sprintf("%d", port),
			"-ssl",
			"-Format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run nikto SSL scan: %w", err)
	}

	return output, nil
}

// ScanWithTuning scans with specific tuning options
func (m *NiktoModule) ScanWithTuning(ctx context.Context, host string, tuning string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{
			"nikto.pl",
			"-h", host,
			"-Tuning", tuning,
			"-Format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run nikto scan with tuning: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Nikto
func (m *NiktoModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{"nikto.pl", "-Version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get nikto version: %w", err)
	}

	return output, nil
}
