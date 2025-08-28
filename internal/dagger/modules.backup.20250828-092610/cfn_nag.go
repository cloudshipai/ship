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

// GetVersion returns the version of cfn-nag
func (m *CfnNagModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithExec([]string{"cfn_nag", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to get cfn-nag version: no output received")
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

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
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

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
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

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with rules: no output received")
}

// ScanWithProfile scans with specific rule profile
func (m *CfnNagModule) ScanWithProfile(ctx context.Context, templatePath string, profilePath string, denyListPath string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath)).
		WithWorkdir("/workspace")

	args := []string{
		"cfn_nag_scan",
		"--input-path", "template.yaml",
		"--output-format", "json",
	}

	if profilePath != "" {
		container = container.WithFile("/workspace/profile.yaml", m.client.Host().File(profilePath))
		args = append(args, "--rule-profile", "profile.yaml")
	}

	if denyListPath != "" {
		container = container.WithFile("/workspace/deny.yaml", m.client.Host().File(denyListPath))
		args = append(args, "--deny-list-path", "deny.yaml")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with profile: no output received")
}

// ListRules lists all available cfn-nag rules
func (m *CfnNagModule) ListRules(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithExec([]string{
			"cfn_nag_rules",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to list cfn-nag rules: no output received")
}

// GenerateWhitelist generates a whitelist template
func (m *CfnNagModule) GenerateWhitelist(ctx context.Context, templatePath string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"cfn_nag_scan",
			"--input-path", "template.yaml",
			"--output-format", "json",
			"--print-suppression",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to generate whitelist: no output received")
}

// ScanWithSuppression scans with rule suppression
func (m *CfnNagModule) ScanWithSuppression(ctx context.Context, templatePath string, suppressRules []string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(templatePath)).
		WithWorkdir("/workspace")

	args := []string{
		"cfn_nag_scan",
		"--input-path", "template.yaml",
		"--output-format", "json",
		"--allow-suppression",
	}

	// If suppress rules provided, create a deny list file
	if len(suppressRules) > 0 {
		denyList := ""
		for _, rule := range suppressRules {
			denyList += rule + "\n"
		}
		container = container.WithNewFile("/workspace/deny.txt", denyList).
			WithExec(append(args, "--deny-list-path", "deny.txt"), dagger.ContainerWithExecOpts{
				Expect: "ANY",
			})
	} else {
		container = container.WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})
	}

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with suppression: no output received")
}

// Scan scans a CloudFormation template with options
func (m *CfnNagModule) Scan(ctx context.Context, inputPath string, outputFormat string, debug bool) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest")

	// Determine if input is file or directory
	args := []string{"cfn_nag_scan", "--input-path"}
	
	// Mount the input file or directory
	if inputPath != "" {
		container = container.WithFile("/workspace/input", m.client.Host().File(inputPath)).
			WithWorkdir("/workspace")
		args = append(args, "input")
	} else {
		args = append(args, ".")
	}

	if outputFormat != "" {
		args = append(args, "--output-format", outputFormat)
	} else {
		args = append(args, "--output-format", "json")
	}

	if debug {
		args = append(args, "--debug")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag scan: no output received")
}

// ScanWithParameters scans with parameter values
func (m *CfnNagModule) ScanWithParameters(ctx context.Context, inputPath string, parameterValuesPath string, conditionValuesPath string, ruleArguments string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(inputPath)).
		WithWorkdir("/workspace")

	args := []string{
		"cfn_nag_scan",
		"--input-path", "template.yaml",
		"--output-format", "json",
	}

	if parameterValuesPath != "" {
		container = container.WithFile("/workspace/params.json", m.client.Host().File(parameterValuesPath))
		args = append(args, "--parameter-values-path", "params.json")
	}

	if conditionValuesPath != "" {
		container = container.WithFile("/workspace/conditions.json", m.client.Host().File(conditionValuesPath))
		args = append(args, "--condition-values-path", "conditions.json")
	}

	if ruleArguments != "" {
		args = append(args, "--rule-arguments", ruleArguments)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to run cfn-nag with parameters: no output received")
}

// SPCMScan runs Stelligent Policy Complexity Metrics scan
func (m *CfnNagModule) SPCMScan(ctx context.Context, inputPath string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("stelligent/cfn_nag:latest").
		WithFile("/workspace/template.yaml", m.client.Host().File(inputPath)).
		WithWorkdir("/workspace")

	// SPCM is typically a metric calculation - using cfn_nag_scan with metrics flag
	args := []string{
		"cfn_nag_scan",
		"--input-path", "template.yaml",
	}

	if outputFormat == "html" {
		args = append(args, "--output-format", "html")
	} else {
		args = append(args, "--output-format", "json")
	}

	// Add policy complexity metrics if available
	args = append(args, "--print-suppression")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "", fmt.Errorf("failed to run SPCM scan: no output received")
}