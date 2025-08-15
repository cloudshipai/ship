package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type InTotoModule struct {
	client *dagger.Client
}

func NewInTotoModule(client *dagger.Client) *InTotoModule {
	return &InTotoModule{
		client: client,
	}
}

// RunStep executes an in-toto step and records metadata
func (m *InTotoModule) RunStep(ctx context.Context, stepName string, command []string, opts ...InTotoOption) (*dagger.Container, error) {
	config := &InTotoConfig{
		KeyPath:     "",
		MaterialDir: "/workspace",
		ProductDir:  "/workspace",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "in-toto"}).
		WithWorkdir("/workspace")

	// Mount directories if specified
	if config.MaterialDir != "" {
		container = container.WithMountedDirectory("/workspace/materials", m.client.Host().Directory(config.MaterialDir))
	}

	args := []string{"in-toto-run", "--step-name", stepName}

	// Add key if provided
	if config.KeyPath != "" {
		container = container.WithMountedFile("/keys/key.pem", m.client.Host().File(config.KeyPath))
		args = append(args, "--key", "/keys/key.pem")
	}

	// Add material and product globs
	if len(config.Materials) > 0 {
		for _, material := range config.Materials {
			args = append(args, "--materials", material)
		}
	}

	if len(config.Products) > 0 {
		for _, product := range config.Products {
			args = append(args, "--products", product)
		}
	}

	// Add the command to execute
	args = append(args, "--")
	args = append(args, command...)

	return container.WithExec(args), nil
}

// VerifySupplyChain verifies the entire supply chain
func (m *InTotoModule) VerifySupplyChain(ctx context.Context, layoutPath string, opts ...InTotoOption) (*dagger.Container, error) {
	config := &InTotoConfig{
		LinkDir: "/workspace/links",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "in-toto"}).
		WithWorkdir("/workspace")

	// Mount layout file
	container = container.WithMountedFile("/workspace/layout.json", m.client.Host().File(layoutPath))

	// Mount links directory
	if config.LinkDir != "" {
		container = container.WithMountedDirectory("/workspace/links", m.client.Host().Directory(config.LinkDir))
	}

	args := []string{"in-toto-verify", "--layout", "/workspace/layout.json"}

	if config.LinkDir != "" {
		args = append(args, "--link-dir", "/workspace/links")
	}

	if len(config.PublicKeys) > 0 {
		for i, keyPath := range config.PublicKeys {
			keyMountPath := fmt.Sprintf("/keys/key%d.pem", i)
			container = container.WithMountedFile(keyMountPath, m.client.Host().File(keyPath))
			args = append(args, "--layout-keys", keyMountPath)
		}
	}

	return container.WithExec(args), nil
}

// GenerateLayout creates an in-toto layout file
func (m *InTotoModule) GenerateLayout(ctx context.Context, opts ...InTotoOption) (*dagger.Container, error) {
	config := &InTotoConfig{}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("alpine:latest").
		WithWorkdir("/workspace")

	// Create basic layout JSON using echo to avoid Go syntax issues
	jsonCmd := `echo '{"_type":"layout","keys":{},"steps":[{"name":"build","expected_command":["make","build"],"expected_materials":["**"],"expected_products":["**"],"pubkeys":[]}],"inspect":[]}' > /workspace/layout.json && echo "Layout generated at /workspace/layout.json"`

	return container.WithExec([]string{"sh", "-c", jsonCmd}), nil
}

// RecordMetadata records step metadata without executing commands
func (m *InTotoModule) RecordMetadata(ctx context.Context, stepName string, opts ...InTotoOption) (*dagger.Container, error) {
	config := &InTotoConfig{
		MaterialDir: "/workspace",
		ProductDir:  "/workspace",
	}

	for _, opt := range opts {
		opt(config)
	}

	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "in-toto"}).
		WithWorkdir("/workspace")

	args := []string{"in-toto-record", "--step-name", stepName}

	// Add key if provided
	if config.KeyPath != "" {
		container = container.WithMountedFile("/keys/key.pem", m.client.Host().File(config.KeyPath))
		args = append(args, "--key", "/keys/key.pem")
	}

	// Add material and product globs
	if len(config.Materials) > 0 {
		for _, material := range config.Materials {
			args = append(args, "--materials", material)
		}
	}

	if len(config.Products) > 0 {
		for _, product := range config.Products {
			args = append(args, "--products", product)
		}
	}

	return container.WithExec(args), nil
}

type InTotoConfig struct {
	KeyPath     string
	MaterialDir string
	ProductDir  string
	LinkDir     string
	Materials   []string
	Products    []string
	PublicKeys  []string
}

type InTotoOption func(*InTotoConfig)

func WithKeyPath(path string) InTotoOption {
	return func(c *InTotoConfig) {
		c.KeyPath = path
	}
}

func WithMaterialDir(dir string) InTotoOption {
	return func(c *InTotoConfig) {
		c.MaterialDir = dir
	}
}

func WithProductDir(dir string) InTotoOption {
	return func(c *InTotoConfig) {
		c.ProductDir = dir
	}
}

func WithLinkDir(dir string) InTotoOption {
	return func(c *InTotoConfig) {
		c.LinkDir = dir
	}
}

func WithMaterials(materials []string) InTotoOption {
	return func(c *InTotoConfig) {
		c.Materials = materials
	}
}

func WithProducts(products []string) InTotoOption {
	return func(c *InTotoConfig) {
		c.Products = products
	}
}

func WithPublicKeys(keys []string) InTotoOption {
	return func(c *InTotoConfig) {
		c.PublicKeys = keys
	}
}
