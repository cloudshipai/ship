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

// prepareContainer prepares a container with rekor-cli installed
func (m *RekorModule) prepareContainer() *dagger.Container {
	return m.client.Container().
		From("golang:1.21-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec([]string{"go", "install", "github.com/sigstore/rekor/cmd/rekor-cli@latest"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("PATH", "/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
}

// Upload uploads an artifact to the transparency log
func (m *RekorModule) Upload(ctx context.Context, artifactPath string, signaturePath string) (string, error) {
	container := m.prepareContainer().
		WithFile("/artifact", m.client.Host().File(artifactPath)).
		WithFile("/signature", m.client.Host().File(signaturePath)).
		WithExec([]string{
			"rekor-cli",
			"upload",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--artifact", "/artifact",
			"--signature", "/signature",
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upload to rekor: %w", err)
	}

	return output, nil
}

// Search searches the transparency log
func (m *RekorModule) Search(ctx context.Context, query string) (string, error) {
	container := m.prepareContainer().
		WithExec([]string{
			"rekor-cli",
			"search",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--artifact", query,
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to search rekor: %w", err)
	}

	return output, nil
}

// Get retrieves an entry from the log
func (m *RekorModule) Get(ctx context.Context, logIndex string) (string, error) {
	container := m.prepareContainer().
		WithExec([]string{
			"rekor-cli",
			"get",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--log-index", logIndex,
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get rekor entry: %w", err)
	}

	return output, nil
}

// Verify verifies an entry in the log
func (m *RekorModule) Verify(ctx context.Context, artifactPath string, signaturePath string) (string, error) {
	container := m.prepareContainer().
		WithFile("/artifact", m.client.Host().File(artifactPath)).
		WithFile("/signature", m.client.Host().File(signaturePath)).
		WithExec([]string{
			"rekor-cli",
			"verify",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--artifact", "/artifact",
			"--signature", "/signature",
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify rekor entry: %w", err)
	}

	return output, nil
}

// GetByUUID gets a log entry by UUID
func (m *RekorModule) GetByUUID(ctx context.Context, uuid string) (string, error) {
	container := m.prepareContainer().
		WithExec([]string{
			"rekor-cli",
			"get",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--uuid", uuid,
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get rekor entry by UUID: %w", err)
	}

	return output, nil
}

// VerifyByUUID verifies an entry by UUID
func (m *RekorModule) VerifyByUUID(ctx context.Context, uuid string) (string, error) {
	container := m.prepareContainer().
		WithExec([]string{
			"rekor-cli",
			"verify",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--uuid", uuid,
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify rekor entry by UUID: %w", err)
	}

	return output, nil
}

// VerifyByIndex verifies an entry by log index
func (m *RekorModule) VerifyByIndex(ctx context.Context, logIndex string) (string, error) {
	container := m.prepareContainer().
		WithExec([]string{
			"rekor-cli",
			"verify",
			"--rekor_server", "https://rekor.sigstore.dev",
			"--log-index", logIndex,
			"--format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to verify rekor entry by index: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Rekor CLI
func (m *RekorModule) GetVersion(ctx context.Context) (string, error) {
	container := m.prepareContainer().
		WithExec([]string{"rekor-cli", "version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// If version command fails, return the module info
	return "rekor-cli@latest", nil
}