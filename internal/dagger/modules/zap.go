package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// ZapModule runs OWASP ZAP for web application security testing
type ZapModule struct {
	client *dagger.Client
	name   string
}

// NewZapModule creates a new ZAP module
func NewZapModule(client *dagger.Client) *ZapModule {
	return &ZapModule{
		client: client,
		name:   "zap",
	}
}

// BaselineScan performs a baseline scan
func (m *ZapModule) BaselineScan(ctx context.Context, target string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithExec([]string{
			"zap-baseline.py",
			"-t", target,
			"-J", "/zap/wrk/baseline-report.json",
			"-r", "/zap/wrk/baseline-report.html",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run ZAP baseline scan: %w", err)
	}

	return output, nil
}

// FullScan performs a full scan
func (m *ZapModule) FullScan(ctx context.Context, target string, maxDuration int) (string, error) {
	container := m.client.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithExec([]string{
			"zap-full-scan.py",
			"-t", target,
			"-J", "/zap/wrk/full-report.json",
			"-r", "/zap/wrk/full-report.html",
			"-m", fmt.Sprintf("%d", maxDuration),
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run ZAP full scan: %w", err)
	}

	return output, nil
}

// ApiScan performs an API scan using OpenAPI/Swagger spec
func (m *ZapModule) ApiScan(ctx context.Context, target string, apiSpecPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithFile("/zap/wrk/api-spec.json", m.client.Host().File(apiSpecPath)).
		WithExec([]string{
			"zap-api-scan.py",
			"-t", target,
			"-f", "openapi",
			"-d", "/zap/wrk/api-spec.json",
			"-J", "/zap/wrk/api-report.json",
			"-r", "/zap/wrk/api-report.html",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run ZAP API scan: %w", err)
	}

	return output, nil
}

// ScanWithContext performs a scan with context file
func (m *ZapModule) ScanWithContext(ctx context.Context, target string, contextPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithFile("/zap/wrk/context.context", m.client.Host().File(contextPath)).
		WithExec([]string{
			"zap-baseline.py",
			"-t", target,
			"-n", "/zap/wrk/context.context",
			"-J", "/zap/wrk/context-report.json",
			"-r", "/zap/wrk/context-report.html",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run ZAP context scan: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of ZAP
func (m *ZapModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithExec([]string{"zap.sh", "-version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get ZAP version: %w", err)
	}

	return output, nil
}
