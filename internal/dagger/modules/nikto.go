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


// ScanHostsFile scans multiple hosts from a file
func (m *NiktoModule) ScanHostsFile(ctx context.Context, hostsFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithFile("/workspace/hosts.txt", m.client.Host().File(hostsFile)).
		WithExec([]string{"nikto.pl", "-h", "/workspace/hosts.txt"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan hosts file: %w", err)
	}

	return output, nil
}

// ScanWithAuth scans with authentication credentials
func (m *NiktoModule) ScanWithAuth(ctx context.Context, host string, authMethod string, credentials string) (string, error) {
	args := []string{"nikto.pl", "-h", host}
	if authMethod != "" && credentials != "" {
		args = append(args, "-id", credentials)
	}

	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan with auth: %w", err)
	}

	return output, nil
}

// ScanWithProxy scans through a proxy server
func (m *NiktoModule) ScanWithProxy(ctx context.Context, host string, proxyHost string, proxyPort string) (string, error) {
	args := []string{"nikto.pl", "-h", host}
	if proxyHost != "" {
		args = append(args, "-useproxy", proxyHost+":"+proxyPort)
	}

	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan with proxy: %w", err)
	}

	return output, nil
}

// ScanWithEvasion scans with IDS evasion techniques
func (m *NiktoModule) ScanWithEvasion(ctx context.Context, host string, evasionLevel string) (string, error) {
	args := []string{"nikto.pl", "-h", host}
	if evasionLevel != "" {
		args = append(args, "-evasion", evasionLevel)
	}

	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan with evasion: %w", err)
	}

	return output, nil
}

// UpdateDatabase updates Nikto's vulnerability database
func (m *NiktoModule) UpdateDatabase(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{"nikto.pl", "-update"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to update database: %w", err)
	}

	return output, nil
}

// DatabaseCheck checks Nikto database integrity
func (m *NiktoModule) DatabaseCheck(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{"nikto.pl", "-dbcheck"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check database: %w", err)
	}

	return output, nil
}

// FindOnly performs discovery-only scan without vulnerability checks
func (m *NiktoModule) FindOnly(ctx context.Context, host string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/nikto:latest").
		WithExec([]string{"nikto.pl", "-h", host, "-findonly"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run find-only scan: %w", err)
	}

	return output, nil
}
