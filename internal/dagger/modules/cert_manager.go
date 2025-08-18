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

const cmctlBinary = "/usr/local/bin/cmctl"

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
		cmctlBinary,
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
		cmctlBinary,
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
		cmctlBinary,
		"renew", name,
		"--namespace", namespace,
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to renew certificate: %w", err)
	}

	return output, nil
}

// Install installs cert-manager using kubectl
func (m *CertManagerModule) Install(ctx context.Context, version string, dryRun bool) (string, error) {
	if version == "" {
		version = "v1.18.2"
	}
	manifestUrl := fmt.Sprintf("https://github.com/cert-manager/cert-manager/releases/download/%s/cert-manager.yaml", version)
	args := []string{"kubectl", "apply", "-f", manifestUrl}
	if dryRun {
		args = append(args, "--dry-run=client")
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install cert-manager: %w", err)
	}

	return output, nil
}

// CheckInstallation checks cert-manager installation status
func (m *CertManagerModule) CheckInstallation(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	if namespace == "" {
		namespace = "cert-manager"
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"kubectl", "get", "pods", "-n", namespace})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check cert-manager installation: %w", err)
	}

	return output, nil
}

// CreateCertificateRequest creates a CertificateRequest using cmctl
func (m *CertManagerModule) CreateCertificateRequest(ctx context.Context, name string, fromCertificateFile string, fetchCertificate bool, timeout string, kubeconfig string) (string, error) {
	args := []string{cmctlBinary, "create", "certificaterequest", name}
	if fromCertificateFile != "" {
		args = append(args, "--from-certificate-file", "/tmp/cert.yaml")
	}
	if fetchCertificate {
		args = append(args, "--fetch-certificate")
	}
	if timeout != "" {
		args = append(args, "--timeout", timeout)
	}

	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}
	if fromCertificateFile != "" {
		container = container.WithFile("/tmp/cert.yaml", m.client.Host().File(fromCertificateFile))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate request: %w", err)
	}

	return output, nil
}

// ListCertificates lists certificates using kubectl
func (m *CertManagerModule) ListCertificates(ctx context.Context, namespace string, allNamespaces bool, kubeconfig string) (string, error) {
	args := []string{"kubectl", "get", "certificates"}
	if allNamespaces {
		args = append(args, "--all-namespaces")
	} else if namespace != "" {
		args = append(args, "-n", namespace)
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list certificates: %w", err)
	}

	return output, nil
}

// RenewAllCertificates renews all certificates in a namespace
func (m *CertManagerModule) RenewAllCertificates(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	args := []string{cmctlBinary, "renew", "--all"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to renew all certificates: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of cert-manager
func (m *CertManagerModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("quay.io/jetstack/cert-manager-ctl:latest").
		WithExec([]string{cmctlBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cert-manager version: %w", err)
	}

	return output, nil
}
