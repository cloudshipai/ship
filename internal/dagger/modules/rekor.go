package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// RekorModule runs Rekor for transparency log operations
type RekorModule struct {
	client *dagger.Client
	name   string
}

// NewRekorModule creates a new Rekor module
func NewRekorModule(client *dagger.Client) *RekorModule {
	return &RekorModule{
		client: client,
		name:   "rekor",
	}
}

// Upload uploads an artifact to the transparency log
func (m *RekorModule) Upload(ctx context.Context, artifactPath string, signaturePath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithFile("/artifact", m.client.Host().File(artifactPath)).
		WithFile("/signature", m.client.Host().File(signaturePath)).
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"upload",
			"--artifact", "/artifact",
			"--signature", "/signature",
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upload to rekor: %w", err)
	}

	return output, nil
}

// Search searches the transparency log
func (m *RekorModule) Search(ctx context.Context, query string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"search",
			"--artifact", query,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to search rekor: %w", err)
	}

	return output, nil
}

// Get retrieves an entry from the log
func (m *RekorModule) Get(ctx context.Context, logIndex string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"get",
			"--log-index", logIndex,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get rekor entry: %w", err)
	}

	return output, nil
}

// Verify verifies an entry in the log
func (m *RekorModule) Verify(ctx context.Context, artifactPath string, signaturePath string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithFile("/artifact", m.client.Host().File(artifactPath)).
		WithFile("/signature", m.client.Host().File(signaturePath)).
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"verify",
			"--artifact", "/artifact",
			"--signature", "/signature",
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify rekor entry: %w", err)
	}

	return output, nil
}

// GetByUUID gets a log entry by UUID
func (m *RekorModule) GetByUUID(ctx context.Context, uuid string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"get",
			"--uuid", uuid,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get rekor entry by UUID: %w", err)
	}

	return output, nil
}

// VerifyByUUID verifies an entry by UUID
func (m *RekorModule) VerifyByUUID(ctx context.Context, uuid string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"verify",
			"--uuid", uuid,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify rekor entry by UUID: %w", err)
	}

	return output, nil
}

// VerifyByIndex verifies an entry by log index
func (m *RekorModule) VerifyByIndex(ctx context.Context, logIndex string) (string, error) {
	container := m.client.Container().
		From("gcr.io/projectsigstore/rekor-cli:latest").
		WithExec([]string{
			"/usr/local/bin/rekor-cli",
			"verify",
			"--log-index", logIndex,
			"--format", "json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify rekor entry by index: %w", err)
	}

	return output, nil
}
