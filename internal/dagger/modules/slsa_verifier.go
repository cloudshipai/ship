package modules

import (
	"context"

	"dagger.io/dagger"
)

type SLSAVerifierModule struct {
	client *dagger.Client
}

func NewSLSAVerifierModule(client *dagger.Client) *SLSAVerifierModule {
	return &SLSAVerifierModule{
		client: client,
	}
}

// VerifyProvenance verifies SLSA provenance for artifacts
func (m *SLSAVerifierModule) VerifyProvenance(ctx context.Context, artifactPath, provenancePath string, opts ...SLSAVerifierOption) (*dagger.Container, error) {
	config := &SLSAVerifierConfig{
		PrintProvenance: false,
		VerifierVersion: "v2.6.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:" + config.VerifierVersion).
		WithWorkdir("/workspace")

	// Mount artifact and provenance files
	if artifactPath != "" {
		container = container.WithMountedFile("/workspace/artifact", m.client.Host().File(artifactPath))
	}
	if provenancePath != "" {
		container = container.WithMountedFile("/workspace/provenance", m.client.Host().File(provenancePath))
	}

	args := []string{"verify-artifact"}

	if artifactPath != "" {
		args = append(args, "/workspace/artifact")
	}

	if provenancePath != "" {
		args = append(args, "--provenance-path", "/workspace/provenance")
	}

	if config.SourceURI != "" {
		args = append(args, "--source-uri", config.SourceURI)
	}

	if config.SourceTag != "" {
		args = append(args, "--source-tag", config.SourceTag)
	}

	if config.BuilderID != "" {
		args = append(args, "--builder-id", config.BuilderID)
	}

	if config.PrintProvenance {
		args = append(args, "--print-provenance")
	}

	return container.WithExec(args), nil
}

// VerifyImage verifies SLSA provenance for container images
func (m *SLSAVerifierModule) VerifyImage(ctx context.Context, imageRef string, opts ...SLSAVerifierOption) (*dagger.Container, error) {
	config := &SLSAVerifierConfig{
		PrintProvenance: false,
		VerifierVersion: "v2.6.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:" + config.VerifierVersion)

	args := []string{"verify-image", imageRef}

	if config.SourceURI != "" {
		args = append(args, "--source-uri", config.SourceURI)
	}

	if config.SourceTag != "" {
		args = append(args, "--source-tag", config.SourceTag)
	}

	if config.BuilderID != "" {
		args = append(args, "--builder-id", config.BuilderID)
	}

	if config.PrintProvenance {
		args = append(args, "--print-provenance")
	}

	return container.WithExec(args), nil
}

// GeneratePolicy generates SLSA policy configuration
func (m *SLSAVerifierModule) GeneratePolicy(ctx context.Context, opts ...SLSAVerifierOption) (*dagger.Container, error) {
	config := &SLSAVerifierConfig{
		VerifierVersion: "v2.6.0",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("ghcr.io/slsa-framework/slsa-verifier:" + config.VerifierVersion).
		WithWorkdir("/workspace")

	args := []string{"policy", "generate"}

	if config.SourceURI != "" {
		args = append(args, "--source-uri", config.SourceURI)
	}

	if config.BuilderID != "" {
		args = append(args, "--builder-id", config.BuilderID)
	}

	return container.WithExec(args), nil
}

type SLSAVerifierConfig struct {
	SourceURI       string
	SourceTag       string
	BuilderID       string
	PrintProvenance bool
	VerifierVersion string
}

type SLSAVerifierOption func(*SLSAVerifierConfig)

func WithSourceURI(uri string) SLSAVerifierOption {
	return func(c *SLSAVerifierConfig) {
		c.SourceURI = uri
	}
}

func WithSourceTag(tag string) SLSAVerifierOption {
	return func(c *SLSAVerifierConfig) {
		c.SourceTag = tag
	}
}

func WithBuilderID(id string) SLSAVerifierOption {
	return func(c *SLSAVerifierConfig) {
		c.BuilderID = id
	}
}

func WithPrintProvenance(print bool) SLSAVerifierOption {
	return func(c *SLSAVerifierConfig) {
		c.PrintProvenance = print
	}
}

func WithVerifierVersion(version string) SLSAVerifierOption {
	return func(c *SLSAVerifierConfig) {
		c.VerifierVersion = version
	}
}
