package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// CfnNagModule runs cfn-nag for CloudFormation security scanning
type CfnNagModule struct {
	client *dagger.Client
	name   string
}

// cfn_nag_scan and cfn_nag_rules are separate binaries from cfn_nag
const cfnNagScanBinary = "/usr/local/bundle/bin/cfn_nag_scan"
const cfnNagRulesBinary = "/usr/local/bundle/bin/cfn_nag_rules"

// NewCfnNagModule creates a new cfn-nag module
func NewCfnNagModule(client *dagger.Client) *CfnNagModule {
	return &CfnNagModule{
		client: client,
		name:   "cfn-nag",
	}
}

// ScanTemplate scans a CloudFormation template
func (m *CfnNagModule) ScanTemplate(ctx context.Context, templatePath string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			cfnNagScanBinary,
			"--input-path", "template.yaml",
			"--output-format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag: no output received")
}

// ScanDirectory scans all CloudFormation templates in a directory
func (m *CfnNagModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			cfnNagScanBinary,
			"--input-path", ".",
			"--output-format", "json",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag on directory: no output received")
}

// ScanWithRules scans with custom rules
func (m *CfnNagModule) ScanWithRules(ctx context.Context, templatePath string, rulesPath string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath)).
		WithDirectory("/workspace/rules", m.client.Host().Directory(rulesPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			cfnNagScanBinary,
			"--input-path", "template.yaml",
			"--output-format", "json",
			"--rule-directory", "rules",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with rules: no output received")
}

// ScanWithProfile scans with specific rule profile
func (m *CfnNagModule) ScanWithProfile(ctx context.Context, templatePath string, profilePath string, denyListPath string) (string, error) {
	args := []string{cfnNagScanBinary, "--input-path", "/workspace/template.yaml"}
	if profilePath != "" {
		args = append(args, "--profile-path", "/workspace/profile.yml")
	}
	if denyListPath != "" {
		args = append(args, "--deny-list-path", "/workspace/denylist.yml")
	}

	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath))

	if profilePath != "" {
		container = container.WithFile("/workspace/profile.yml", m.client.Host().File(profilePath))
	}
	if denyListPath != "" {
		container = container.WithFile("/workspace/denylist.yml", m.client.Host().File(denyListPath))
	}

	container = container.WithWorkdir("/workspace").WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with profile: no output received")
}

// ScanWithParameters scans with parameter values
func (m *CfnNagModule) ScanWithParameters(ctx context.Context, templatePath string, parameterValuesPath string, conditionValuesPath string, ruleArguments string) (string, error) {
	args := []string{cfnNagScanBinary, "--input-path", "/workspace/template.yaml"}
	if parameterValuesPath != "" {
		args = append(args, "--parameter-values-path", "/workspace/parameters.json")
	}
	if conditionValuesPath != "" {
		args = append(args, "--condition-values-path", "/workspace/conditions.json")
	}
	if ruleArguments != "" {
		args = append(args, "--rule-arguments", ruleArguments)
	}

	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath))

	if parameterValuesPath != "" {
		container = container.WithFile("/workspace/parameters.json", m.client.Host().File(parameterValuesPath))
	}
	if conditionValuesPath != "" {
		container = container.WithFile("/workspace/conditions.json", m.client.Host().File(conditionValuesPath))
	}

	container = container.WithWorkdir("/workspace").WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with parameters: no output received")
}

// ListRules lists all available CFN Nag rules
func (m *CfnNagModule) ListRules(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithExec([]string{cfnNagRulesBinary})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list cfn-nag rules: %w", err)
	}

	return output, nil
}

// SPCMScan generates Stelligent Policy Complexity Metrics report
func (m *CfnNagModule) SPCMScan(ctx context.Context, templatePath string, outputFormat string) (string, error) {
	args := []string{"spcm_scan", "--input-path", "/workspace/template.yaml"}
	if outputFormat != "" {
		args = append(args, "--output-format", outputFormat)
	}

	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run spcm scan: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of cfn-nag
func (m *CfnNagModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithExec([]string{cfnNagScanBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cfn-nag version: %w", err)
	}

	return output, nil
}

// Scan scans CloudFormation templates with configurable options
func (m *CfnNagModule) Scan(ctx context.Context, inputPath string, outputFormat string, debug bool) (string, error) {
	args := []string{"cfn_nag_scan", "--input-path", "/workspace/input"}
	if outputFormat == "json" {
		args = append(args, "--output-format", "json")
	}
	if debug {
		args = append(args, "--debug")
	}

	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/input", m.client.Host().File(inputPath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{Expect: "ANY"})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag scan: no output received")
}
