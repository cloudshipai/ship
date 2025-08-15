package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// VeleroModule runs Velero for Kubernetes backup and restore
type VeleroModule struct {
	client *dagger.Client
	name   string
}

// NewVeleroModule creates a new Velero module
func NewVeleroModule(client *dagger.Client) *VeleroModule {
	return &VeleroModule{
		client: client,
		name:   "velero",
	}
}

// CreateBackup creates a backup of Kubernetes resources
func (m *VeleroModule) CreateBackup(ctx context.Context, backupName string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"velero",
		"backup", "create", backupName,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create velero backup: %w", err)
	}

	return output, nil
}

// ListBackups lists all backups
func (m *VeleroModule) ListBackups(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"velero",
		"backup", "get",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list velero backups: %w", err)
	}

	return output, nil
}

// RestoreBackup restores from a backup
func (m *VeleroModule) RestoreBackup(ctx context.Context, backupName string, restoreName string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"velero",
		"restore", "create", restoreName,
		"--from-backup", backupName,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to restore velero backup: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Velero
func (m *VeleroModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest").
		WithExec([]string{"velero", "version", "--client-only"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get velero version: %w", err)
	}

	return output, nil
}
