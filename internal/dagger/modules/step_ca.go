package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// StepCAModule runs Step CA for certificate authority operations
type StepCAModule struct {
	client *dagger.Client
	name   string
}

const stepBinary = "/usr/local/bin/step"

// NewStepCAModule creates a new Step CA module
func NewStepCAModule(client *dagger.Client) *StepCAModule {
	return &StepCAModule{
		client: client,
		name:   "step-ca",
	}
}

// InitCA initializes a certificate authority
func (m *StepCAModule) InitCA(ctx context.Context, name string, dnsName string) (string, error) {
	container := m.client.Container().
		From("smallstep/step-ca:latest").
		WithExec([]string{
			stepBinary, "ca", "init",
			"--name", name,
			"--dns", dnsName,
			"--provisioner", "admin",
			"--password-file", "/dev/null",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to initialize CA: %w", err)
	}

	return output, nil
}

// CreateCertificate creates a certificate
func (m *StepCAModule) CreateCertificate(ctx context.Context, subject string, caURL string, rootCert string) (string, error) {
	container := m.client.Container().
		From("smallstep/step-ca:latest")

	if rootCert != "" {
		container = container.WithFile("/root.crt", m.client.Host().File(rootCert))
	}

	container = container.WithExec([]string{
		stepBinary, "ca", "certificate",
		subject,
		"/tmp/cert.crt",
		"/tmp/cert.key",
		"--ca-url", caURL,
		"--root", "/root.crt",
		"--output-format", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate: %w", err)
	}

	return output, nil
}

// RenewCertificate renews a certificate
func (m *StepCAModule) RenewCertificate(ctx context.Context, certPath string, keyPath string, caURL string) (string, error) {
	container := m.client.Container().
		From("smallstep/step-ca:latest").
		WithFile("/cert.crt", m.client.Host().File(certPath)).
		WithFile("/cert.key", m.client.Host().File(keyPath)).
		WithExec([]string{
			stepBinary, "ca", "renew",
			"/cert.crt",
			"/cert.key",
			"--ca-url", caURL,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to renew certificate: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Step CA
func (m *StepCAModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("smallstep/step-ca:latest").
		WithExec([]string{stepBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get step-ca version: %w", err)
	}

	return output, nil
}

// AddProvisioner adds a new provisioner to the CA
func (m *StepCAModule) AddProvisioner(ctx context.Context, name string, provisionerType string, caConfig string) (string, error) {
	args := []string{stepBinary, "ca", "provisioner", "add", name, "--type", provisionerType}
	if caConfig != "" {
		args = append(args, "--ca-config", "/config/ca.json")
	}

	container := m.client.Container().
		From("smallstep/step-ca:latest")

	if caConfig != "" {
		container = container.WithFile("/config/ca.json", m.client.Host().File(caConfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to add provisioner: %w", err)
	}

	return output, nil
}

// RevokeCertificate revokes a certificate
func (m *StepCAModule) RevokeCertificate(ctx context.Context, certPath string, keyPath string, caURL string, reason string) (string, error) {
	args := []string{stepBinary, "ca", "revoke"}
	if certPath != "" {
		args = append(args, "--cert", "/cert.pem")
	}
	if keyPath != "" {
		args = append(args, "--key", "/key.pem")
	}
	if caURL != "" {
		args = append(args, "--ca-url", caURL)
	}
	if reason != "" {
		args = append(args, "--reason", reason)
	}

	container := m.client.Container().
		From("smallstep/step-ca:latest")

	if certPath != "" {
		container = container.WithFile("/cert.pem", m.client.Host().File(certPath))
	}
	if keyPath != "" {
		container = container.WithFile("/key.pem", m.client.Host().File(keyPath))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to revoke certificate: %w", err)
	}

	return output, nil
}
