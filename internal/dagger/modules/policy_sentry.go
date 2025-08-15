package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// PolicySentryModule runs Policy Sentry for AWS IAM policy generation
type PolicySentryModule struct {
	client *dagger.Client
	name   string
}

// NewPolicySentryModule creates a new Policy Sentry module
func NewPolicySentryModule(client *dagger.Client) *PolicySentryModule {
	return &PolicySentryModule{
		client: client,
		name:   "policy-sentry",
	}
}

// CreateTemplate creates a policy template
func (m *PolicySentryModule) CreateTemplate(ctx context.Context, templateType string, outputFile string) (string, error) {
	args := []string{"policy-sentry", "create-template", "--template-type", templateType}
	
	if outputFile != "" {
		args = append(args, "--output-file", outputFile)
	}

	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to create template: %w", err)
	}

	return output, nil
}

// WritePolicy writes an IAM policy from a YAML template
func (m *PolicySentryModule) WritePolicy(ctx context.Context, inputFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithFile("/workspace/input.yml", m.client.Host().File(inputFile)).
		WithWorkdir("/workspace").
		WithExec([]string{"policy-sentry", "write-policy", "--input-file", "input.yml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to write policy: %w", err)
	}

	return output, nil
}

// WritePolicyFromTemplate writes a policy from an inline template
func (m *PolicySentryModule) WritePolicyFromTemplate(ctx context.Context, templateYAML string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithNewFile("/workspace/template.yml", templateYAML).
		WithWorkdir("/workspace").
		WithExec([]string{"policy-sentry", "write-policy", "--input-file", "template.yml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to write policy from template: %w", err)
	}

	return output, nil
}

// WritePolicyWithActions writes a policy based on specific actions
func (m *PolicySentryModule) WritePolicyWithActions(ctx context.Context, actions []string, resourceArns []string) (string, error) {
	// Create a simple actions-based template
	template := `mode: actions
name: 'MyPolicy'
actions:
`
	for _, action := range actions {
		template += fmt.Sprintf("  - '%s'\n", action)
	}

	if len(resourceArns) > 0 {
		template += "conditions:\n"
		for _, arn := range resourceArns {
			template += fmt.Sprintf("  - '%s'\n", arn)
		}
	}

	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithNewFile("/workspace/actions.yml", template).
		WithWorkdir("/workspace").
		WithExec([]string{"policy-sentry", "write-policy", "--input-file", "actions.yml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to write actions-based policy: %w", err)
	}

	return output, nil
}

// WritePolicyWithCRUD writes a policy based on CRUD operations
func (m *PolicySentryModule) WritePolicyWithCRUD(ctx context.Context, resourceArns []string, accessLevels []string) (string, error) {
	// Create a CRUD-based template
	template := `mode: crud
name: 'MyPolicy'
crud:
`
	for _, arn := range resourceArns {
		template += fmt.Sprintf("  '%s':\n", arn)
		for _, level := range accessLevels {
			template += fmt.Sprintf("    - '%s'\n", level)
		}
	}

	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithNewFile("/workspace/crud.yml", template).
		WithWorkdir("/workspace").
		WithExec([]string{"policy-sentry", "write-policy", "--input-file", "crud.yml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to write CRUD-based policy: %w", err)
	}

	return output, nil
}

// QueryActionTable queries the action table for service information
func (m *PolicySentryModule) QueryActionTable(ctx context.Context, service string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithExec([]string{"policy-sentry", "query", "action-table", "--service", service})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to query action table: %w", err)
	}

	return output, nil
}

// QueryConditionTable queries the condition table for service information
func (m *PolicySentryModule) QueryConditionTable(ctx context.Context, service string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/policy-sentry:latest").
		WithExec([]string{"policy-sentry", "query", "condition-table", "--service", service})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to query condition table: %w", err)
	}

	return output, nil
}