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
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
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
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
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
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
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
		WithExec([]string{"velero", "version", "--client-only"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// If version command fails, return the image tag
	return "velero/velero:latest", nil
}

// Install installs Velero with storage provider configuration
func (m *VeleroModule) Install(ctx context.Context, provider string, bucket string, region string, secretFile string, useNodeAgent bool, namespace string, noSecret bool, plugins string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if secretFile != "" {
		container = container.WithFile("/secret", m.client.Host().File(secretFile))
	}

	args := []string{"velero", "install"}
	if provider != "" {
		args = append(args, "--provider", provider)
	}
	if bucket != "" {
		args = append(args, "--bucket", bucket)
	}
	if region != "" {
		args = append(args, "--backup-location-config", "region="+region)
	}
	if secretFile != "" {
		args = append(args, "--secret-file", "/secret")
	}
	if useNodeAgent {
		args = append(args, "--use-node-agent")
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	if noSecret {
		args = append(args, "--no-secret")
	}
	if plugins != "" {
		args = append(args, "--plugins", plugins)
	}
	args = append(args, "--wait")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install velero: %w", err)
	}

	return output, nil
}

// CreateSchedule creates a backup schedule
func (m *VeleroModule) CreateSchedule(ctx context.Context, name string, schedule string, includeNamespaces string, excludeNamespaces string, ttl string, labels string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "schedule", "create", name, "--schedule", schedule}
	if ttl != "" {
		args = append(args, "--ttl", ttl)
	}
	if includeNamespaces != "" {
		args = append(args, "--include-namespaces", includeNamespaces)
	}
	if excludeNamespaces != "" {
		args = append(args, "--exclude-namespaces", excludeNamespaces)
	}
	if labels != "" {
		args = append(args, "--labels", labels)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create schedule: %w", err)
	}

	return output, nil
}

// CreateBackupAdvanced creates an on-demand backup with advanced options
func (m *VeleroModule) CreateBackupAdvanced(ctx context.Context, name string, fromSchedule string, includeNamespaces string, excludeNamespaces string, labels string, wait bool, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "backup", "create", name}
	if fromSchedule != "" {
		args = append(args, "--from-schedule", fromSchedule)
	}
	if includeNamespaces != "" {
		args = append(args, "--include-namespaces", includeNamespaces)
	}
	if excludeNamespaces != "" {
		args = append(args, "--exclude-namespaces", excludeNamespaces)
	}
	if labels != "" {
		args = append(args, "--labels", labels)
	}
	if wait {
		args = append(args, "--wait")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create advanced backup: %w", err)
	}

	return output, nil
}

// CreateRestoreAdvanced creates a restore from backup with advanced options
func (m *VeleroModule) CreateRestoreAdvanced(ctx context.Context, name string, fromBackup string, namespaceMappings string, includeNamespaces string, excludeNamespaces string, restorePVs bool, wait bool, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "restore", "create", name, "--from-backup", fromBackup}
	if namespaceMappings != "" {
		args = append(args, "--namespace-mappings", namespaceMappings)
	}
	if includeNamespaces != "" {
		args = append(args, "--include-namespaces", includeNamespaces)
	}
	if excludeNamespaces != "" {
		args = append(args, "--exclude-namespaces", excludeNamespaces)
	}
	if restorePVs {
		args = append(args, "--restore-volumes=true")
	}
	if wait {
		args = append(args, "--wait")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create advanced restore: %w", err)
	}

	return output, nil
}

// GetBackups gets list of backups
func (m *VeleroModule) GetBackups(ctx context.Context, output string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "backup", "get"}
	if output != "" {
		args = append(args, "-o", output)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get backups: %w", err)
	}

	return output_result, nil
}

// DescribeBackup describes a specific backup
func (m *VeleroModule) DescribeBackup(ctx context.Context, name string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"velero", "backup", "describe", name}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe backup: %w", err)
	}

	return output, nil
}

// GetRestores gets list of restores
func (m *VeleroModule) GetRestores(ctx context.Context, output string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "restore", "get"}
	if output != "" {
		args = append(args, "-o", output)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get restores: %w", err)
	}

	return output_result, nil
}

// DescribeRestore describes a specific restore
func (m *VeleroModule) DescribeRestore(ctx context.Context, name string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"velero", "restore", "describe", name}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe restore: %w", err)
	}

	return output, nil
}

// GetSchedules gets list of backup schedules
func (m *VeleroModule) GetSchedules(ctx context.Context, output string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "schedule", "get"}
	if output != "" {
		args = append(args, "-o", output)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get schedules: %w", err)
	}

	return output_result, nil
}

// DeleteBackup deletes a backup
func (m *VeleroModule) DeleteBackup(ctx context.Context, name string, confirm bool, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "backup", "delete", name}
	if confirm {
		args = append(args, "--confirm")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to delete backup: %w", err)
	}

	return output, nil
}

// DeleteSchedule deletes a backup schedule
func (m *VeleroModule) DeleteSchedule(ctx context.Context, name string, confirm bool, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "schedule", "delete", name}
	if confirm {
		args = append(args, "--confirm")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to delete schedule: %w", err)
	}

	return output, nil
}

// CreateBackupLocation creates backup storage location
func (m *VeleroModule) CreateBackupLocation(ctx context.Context, name string, provider string, bucket string, config string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "backup-location", "create", name, "--provider", provider, "--bucket", bucket}
	if config != "" {
		args = append(args, "--config", config)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create backup location: %w", err)
	}

	return output, nil
}

// GetBackupLocations gets backup storage locations
func (m *VeleroModule) GetBackupLocations(ctx context.Context, output string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"velero", "backup-location", "get"}
	if output != "" {
		args = append(args, "-o", output)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get backup locations: %w", err)
	}

	return output_result, nil
}

// GetBackupLogs gets backup logs
func (m *VeleroModule) GetBackupLogs(ctx context.Context, name string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"velero", "backup", "logs", name}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get backup logs: %w", err)
	}

	return output, nil
}

// GetRestoreLogs gets restore logs
func (m *VeleroModule) GetRestoreLogs(ctx context.Context, name string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("velero/velero:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{"velero", "restore", "logs", name}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get restore logs: %w", err)
	}

	return output, nil
}
