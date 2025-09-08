package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// BuildXModule provides Docker BuildX functionality for multi-platform image building
type BuildXModule struct {
	client *dagger.Client
	name   string
}

// NewBuildXModule creates a new BuildX module
func NewBuildXModule(client *dagger.Client) *BuildXModule {
	return &BuildXModule{
		client: client,
		name:   "buildx",
	}
}

// Build builds an OCI image using Docker BuildX
func (m *BuildXModule) Build(ctx context.Context, srcDir string, tag string, platform string, dockerfilePath string) (string, error) {
	// Set defaults
	if platform == "" {
		platform = "linux/amd64"
	}
	if dockerfilePath == "" {
		dockerfilePath = "."
	}

	// Create a container with Docker BuildX installed
	container := m.client.Container().
		From(getImageTag("buildx", "docker:latest")).
		WithUnixSocket("/var/run/docker.sock", m.client.Host().UnixSocket("/var/run/docker.sock")).
		WithDirectory("/build", m.client.Host().Directory(srcDir)).
		WithWorkdir("/build")

	// Set up buildx
	container = container.WithExec([]string{
		"docker", "buildx", "create", "--use", "--name", "multiarch-builder",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	// Build the image with BuildX
	args := []string{
		"docker", "buildx", "build",
		"--platform", platform,
		"-t", tag,
	}
	
	// Handle dockerfile path
	if dockerfilePath != "" && dockerfilePath != "." {
		args = append(args, "-f", dockerfilePath)
	}
	
	args = append(args, ".")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("buildx build failed - stdout error: %v, stderr: %s", err, stderr)
	}

	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	return "", fmt.Errorf("buildx build failed: no output received, stderr: %s", stderr)
}

// Publish builds and pushes an OCI image to a registry
func (m *BuildXModule) Publish(ctx context.Context, srcDir string, tag string, platform string, username string, password string, registry string, dockerfilePath string) (string, error) {
	// Set defaults
	if platform == "" {
		platform = "linux/amd64"
	}
	if registry == "" {
		registry = "docker.io"
	}
	if dockerfilePath == "" {
		dockerfilePath = "."
	}

	// Create a container with Docker BuildX installed
	container := m.client.Container().
		From(getImageTag("buildx", "docker:latest")).
		WithUnixSocket("/var/run/docker.sock", m.client.Host().UnixSocket("/var/run/docker.sock")).
		WithDirectory("/build", m.client.Host().Directory(srcDir)).
		WithWorkdir("/build")

	// Login to registry
	if username != "" && password != "" {
		container = container.WithExec([]string{
			"sh", "-c", fmt.Sprintf("echo '%s' | docker login %s -u %s --password-stdin", password, registry, username),
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})
	}

	// Set up buildx
	container = container.WithExec([]string{
		"docker", "buildx", "create", "--use", "--name", "multiarch-builder",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	// Build and push the image with BuildX
	args := []string{
		"docker", "buildx", "build",
		"--platform", platform,
		"-t", tag,
		"--push",
	}
	
	// Handle dockerfile path
	if dockerfilePath != "" && dockerfilePath != "." {
		args = append(args, "-f", dockerfilePath)
	}
	
	args = append(args, ".")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("buildx publish failed - stdout error: %v, stderr: %s", err, stderr)
	}

	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	return "", fmt.Errorf("buildx publish failed: no output received, stderr: %s", stderr)
}

// Dev returns a development container with Docker BuildX installed
func (m *BuildXModule) Dev(ctx context.Context, srcDir string) (string, error) {
	container := m.client.Container().
		From(getImageTag("buildx", "docker:latest")).
		WithUnixSocket("/var/run/docker.sock", m.client.Host().UnixSocket("/var/run/docker.sock"))

	if srcDir != "" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(srcDir)).
			WithWorkdir("/workspace")
	}

	// Set up buildx
	container = container.WithExec([]string{
		"docker", "buildx", "install",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	container = container.WithExec([]string{
		"docker", "buildx", "create", "--use", "--name", "dev-builder",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("buildx dev setup failed - stdout error: %v, stderr: %s", err, stderr)
	}

	return fmt.Sprintf("BuildX development environment ready. %s", output), nil
}

// GetVersion returns the BuildX version information
func (m *BuildXModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From(getImageTag("buildx", "docker:latest")).
		WithExec([]string{"docker", "buildx", "version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to get buildx version - stdout error: %v, stderr: %s", err, stderr)
	}

	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	return "", fmt.Errorf("failed to get buildx version: no output received, stderr: %s", stderr)
}