package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// CertManagerModule runs cert-manager for certificate management
type CertManagerModule struct {
	client *dagger.Client
	name   string
}

// NewCertManagerModule creates a new cert-manager module
func NewCertManagerModule(client *dagger.Client) *CertManagerModule {
	return &CertManagerModule{
		client: client,
		name:   "cert-manager",
	}
}

// GetCertificates lists certificates
func (m *CertManagerModule) GetCertificates(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"cmctl",
		"status", "certificate",
		"--namespace", namespace,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get certificates: %w", err)
	}

	return output, nil
}

// CheckCertificate checks certificate status
func (m *CertManagerModule) CheckCertificate(ctx context.Context, name string, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"cmctl",
		"status", "certificate", name,
		"--namespace", namespace,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check certificate: %w", err)
	}

	return output, nil
}

// RenewCertificate renews a certificate
func (m *CertManagerModule) RenewCertificate(ctx context.Context, name string, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"cmctl",
		"renew", name,
		"--namespace", namespace,
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to renew certificate: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of cert-manager
func (m *CertManagerModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest").
		WithExec([]string{"cmctl", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cert-manager version: %w", err)
	}

	return output, nil
}
