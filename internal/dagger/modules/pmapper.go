package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// PMapperModule runs PMapper for AWS IAM privilege mapping
type PMapperModule struct {
	client *dagger.Client
	name   string
}

// NewPMapperModule creates a new PMapper module
func NewPMapperModule(client *dagger.Client) *PMapperModule {
	return &PMapperModule{
		client: client,
		name:   "pmapper",
	}
}

// CreateGraph creates a privilege graph for an AWS account
func (m *PMapperModule) CreateGraph(ctx context.Context, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/pmapper:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{"pmapper", "graph", "create"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to create graph: %w", err)
	}

	return output, nil
}

// QueryAccess queries if a principal can access a specific action/resource
func (m *PMapperModule) QueryAccess(ctx context.Context, profile string, principal string, action string, resource string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/pmapper:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	args := []string{"pmapper", "query", "can", principal, action}
	if resource != "" {
		args = append(args, resource)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to query access: %w", err)
	}

	return output, nil
}

// FindPrivilegeEscalation finds privilege escalation paths
func (m *PMapperModule) FindPrivilegeEscalation(ctx context.Context, profile string, principal string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/pmapper:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{"pmapper", "query", "preset", "privesc", principal})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to find privilege escalation: %w", err)
	}

	return output, nil
}

// VisualizeGraph creates a visual representation of the privilege graph
func (m *PMapperModule) VisualizeGraph(ctx context.Context, profile string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/pmapper:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	args := []string{"pmapper", "visualize"}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to visualize graph: %w", err)
	}

	return output, nil
}

// ListPrincipals lists all principals in the AWS account
func (m *PMapperModule) ListPrincipals(ctx context.Context, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/pmapper:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{"pmapper", "query", "list", "principals"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to list principals: %w", err)
	}

	return output, nil
}

// CheckAdminAccess checks if a principal has admin access
func (m *PMapperModule) CheckAdminAccess(ctx context.Context, profile string, principal string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/pmapper:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{"pmapper", "query", "preset", "admin", principal})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to check admin access: %w", err)
	}

	return output, nil
}