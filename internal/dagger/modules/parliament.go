package modules

import (
	"context"

	"dagger.io/dagger"
)

// ParliamentModule runs Parliament for AWS IAM policy linting
type ParliamentModule struct {
	client *dagger.Client
	name   string
}

const parliamentBinary = "/usr/local/bin/parliament"

// NewParliamentModule creates a new Parliament module
func NewParliamentModule(client *dagger.Client) *ParliamentModule {
	return &ParliamentModule{
		client: client,
		name:   parliamentBinary,
	}
}

// LintPolicyFile lints a specific IAM policy file
func (m *ParliamentModule) LintPolicyFile(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath)).
		WithWorkdir("/workspace").
		WithExec([]string{parliamentBinary, "--file", "policy.json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Parliament returns non-zero exit code for policy violations
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintPolicyDirectory lints all policy files in a directory
func (m *ParliamentModule) LintPolicyDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{parliamentBinary, "--directory", "."})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintPolicyString lints a policy provided as a string
func (m *ParliamentModule) LintPolicyString(ctx context.Context, policyJSON string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithNewFile("/workspace/policy.json", policyJSON).
		WithWorkdir("/workspace").
		WithExec([]string{parliamentBinary, "--file", "policy.json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintWithCommunityAuditors lints using community auditors
func (m *ParliamentModule) LintWithCommunityAuditors(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath)).
		WithWorkdir("/workspace").
		WithExec([]string{parliamentBinary, "--file", "policy.json", "--include-community-auditors"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintWithPrivateAuditors lints using private auditors
func (m *ParliamentModule) LintWithPrivateAuditors(ctx context.Context, policyPath string, auditorsPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath)).
		WithDirectory("/workspace/auditors", m.client.Host().Directory(auditorsPath)).
		WithWorkdir("/workspace").
		WithExec([]string{parliamentBinary, "--file", "policy.json", "--private-auditors-dir", "auditors"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintWithSeverityFilter lints and filters by severity level
func (m *ParliamentModule) LintWithSeverityFilter(ctx context.Context, policyPath string, minSeverity string) (string, error) {
	args := []string{parliamentBinary, "--file", "policy.json"}
	
	if minSeverity != "" {
		args = append(args, "--minimum-severity", minSeverity)
	}

	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintAWSManagedPolicies lints AWS managed policies
func (m *ParliamentModule) LintAWSManagedPolicies(ctx context.Context, config string, jsonOutput bool) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest")

	args := []string{parliamentBinary, "--aws-managed-policies"}
	if config != "" {
		container = container.WithFile("/workspace/config.yaml", m.client.Host().File(config))
		args = append(args, "--config", "/workspace/config.yaml")
	}
	if jsonOutput {
		args = append(args, "--json")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// LintAuthDetailsFile lints AWS IAM authorization details file
func (m *ParliamentModule) LintAuthDetailsFile(ctx context.Context, authDetailsFile string, config string, jsonOutput bool) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithFile("/workspace/auth_details.json", m.client.Host().File(authDetailsFile))

	args := []string{parliamentBinary, "--auth-details-file", "/workspace/auth_details.json"}
	if config != "" {
		container = container.WithFile("/workspace/config.yaml", m.client.Host().File(config))
		args = append(args, "--config", "/workspace/config.yaml")
	}
	if jsonOutput {
		args = append(args, "--json")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// ComprehensiveAnalysis performs comprehensive IAM policy analysis with all auditors
func (m *ParliamentModule) ComprehensiveAnalysis(ctx context.Context, policyPath string, privateAuditors string, config string, jsonOutput bool) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath))

	if privateAuditors != "" {
		container = container.WithDirectory("/workspace/auditors", m.client.Host().Directory(privateAuditors))
	}
	if config != "" {
		container = container.WithFile("/workspace/config.yaml", m.client.Host().File(config))
	}

	args := []string{parliamentBinary, "--file", "/workspace/policy.json", "--include-community-auditors"}
	if privateAuditors != "" {
		args = append(args, "--private_auditors", "/workspace/auditors")
	}
	if config != "" {
		args = append(args, "--config", "/workspace/config.yaml")
	}
	if jsonOutput {
		args = append(args, "--json")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}

// BatchDirectoryAnalysis performs batch analysis of multiple policy directories
func (m *ParliamentModule) BatchDirectoryAnalysis(ctx context.Context, baseDirectory string, config string, privateAuditors string, jsonOutput bool, includeExtension string, excludePattern string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/parliament:latest").
		WithDirectory("/workspace", m.client.Host().Directory(baseDirectory))

	if privateAuditors != "" {
		container = container.WithDirectory("/workspace/auditors", m.client.Host().Directory(privateAuditors))
	}
	if config != "" {
		container = container.WithFile("/workspace/config.yaml", m.client.Host().File(config))
	}

	args := []string{parliamentBinary, "--directory", "/workspace", "--include-community-auditors"}
	if privateAuditors != "" {
		args = append(args, "--private_auditors", "/workspace/auditors")
	}
	if config != "" {
		args = append(args, "--config", "/workspace/config.yaml")
	}
	if jsonOutput {
		args = append(args, "--json")
	}
	if includeExtension != "" {
		args = append(args, "--include_policy_extension", includeExtension)
	}
	if excludePattern != "" {
		args = append(args, "--exclude_pattern", excludePattern)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return output, nil
	}

	return output, nil
}