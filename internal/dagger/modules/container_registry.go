package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// ContainerRegistryModule provides Docker registry operations
type ContainerRegistryModule struct {
	client *dagger.Client
	name   string
}

const dockerBinary = "/usr/local/bin/docker"

// NewContainerRegistryModule creates a new container registry module
func NewContainerRegistryModule(client *dagger.Client) *ContainerRegistryModule {
	return &ContainerRegistryModule{
		client: client,
		name:   "container-registry",
	}
}

// Login to container registry
func (m *ContainerRegistryModule) Login(ctx context.Context, registry, username, password string) (string, error) {
	args := []string{dockerBinary, "login"}
	if registry != "" {
		args = append(args, registry)
	}
	if username != "" {
		args = append(args, "--username", username)
	}
	if password != "" {
		args = append(args, "--password", password)
	}

	container := m.client.Container().
		From("docker:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to login to registry: %w", err)
	}

	return output, nil
}

// PushImage pushes an image to the registry
func (m *ContainerRegistryModule) PushImage(ctx context.Context, image string) (string, error) {
	container := m.client.Container().
		From("docker:latest").
		WithExec([]string{dockerBinary, "push", image})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to push image: %w", err)
	}

	return output, nil
}

// PullImage pulls an image from the registry
func (m *ContainerRegistryModule) PullImage(ctx context.Context, image string) (string, error) {
	container := m.client.Container().
		From("docker:latest").
		WithExec([]string{dockerBinary, "pull", image})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	return output, nil
}

// ListImages lists local Docker images
func (m *ContainerRegistryModule) ListImages(ctx context.Context, repository string, all bool) (string, error) {
	args := []string{dockerBinary, "images"}
	if repository != "" {
		args = append(args, repository)
	}
	if all {
		args = append(args, "--all")
	}

	container := m.client.Container().
		From("docker:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list images: %w", err)
	}

	return output, nil
}

// TagImage creates a tag for an image
func (m *ContainerRegistryModule) TagImage(ctx context.Context, sourceImage, targetImage string) (string, error) {
	container := m.client.Container().
		From("docker:latest").
		WithExec([]string{dockerBinary, "tag", sourceImage, targetImage})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to tag image: %w", err)
	}

	return output, nil
}

// Logout from container registry
func (m *ContainerRegistryModule) Logout(ctx context.Context, registry string) (string, error) {
	args := []string{dockerBinary, "logout"}
	if registry != "" {
		args = append(args, registry)
	}

	container := m.client.Container().
		From("docker:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to logout from registry: %w", err)
	}

	return output, nil
}