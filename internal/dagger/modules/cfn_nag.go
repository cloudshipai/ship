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
			"cfn_nag_scan",
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
			"cfn_nag_scan",
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
			"cfn_nag_scan",
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

// GetVersion returns the version of cfn-nag
func (m *CfnNagModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithExec([]string{"cfn_nag", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get cfn-nag version: %w", err)
	}

	return output, nil
}
