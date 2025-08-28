package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// IacPlanModule runs Terraform/OpenTofu plan operations
type IacPlanModule struct {
	client *dagger.Client
	name   string
}

const (
	terraformBinary = "/bin/terraform"
	tofuBinary = "tofu"
)

func getBinaryPath(tool string) string {
	switch tool {
	case "tofu":
		return tofuBinary
	default:
		return terraformBinary
	}
}

// NewIacPlanModule creates a new IaC plan module
func NewIacPlanModule(client *dagger.Client) *IacPlanModule {
	return &IacPlanModule{
		client: client,
		name:   "iac-plan",
	}
}

// GeneratePlan generates Terraform/OpenTofu plan and exports to JSON
func (m *IacPlanModule) GeneratePlan(ctx context.Context, workdir string, tool string, varFiles []string, destroy bool) (string, error) {
	var image string
	switch tool {
	case "tofu":
		image = "opentofu/opentofu:latest"
	default:
		image = "hashicorp/terraform:latest"
	}

	container := m.client.Container().
		From(image).
		WithDirectory("/workspace", m.client.Host().Directory(workdir)).
		WithWorkdir("/workspace")

	// Initialize if needed
	container = container.WithExec([]string{tool, "init"})

	// Generate plan
	planCmd := []string{tool, "plan", "-out=tfplan"}
	if destroy {
		planCmd = append(planCmd, "-destroy")
	}

	for _, varFile := range varFiles {
		planCmd = append(planCmd, "-var-file="+varFile)
	}

	container = container.WithExec(planCmd)

	// Convert to JSON
	container = container.WithExec([]string{tool, "show", "-json", "tfplan"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate %s plan: %w", tool, err)
	}

	return output, nil
}

// ValidateConfiguration validates Terraform/OpenTofu configuration
func (m *IacPlanModule) ValidateConfiguration(ctx context.Context, workdir string, tool string) (string, error) {
	var image string
	switch tool {
	case "tofu":
		image = "opentofu/opentofu:latest"
	default:
		image = "hashicorp/terraform:latest"
	}

	container := m.client.Container().
		From(image).
		WithDirectory("/workspace", m.client.Host().Directory(workdir)).
		WithWorkdir("/workspace").
		WithExec([]string{getBinaryPath(tool), "init"}).
		WithExec([]string{getBinaryPath(tool), "validate", "-json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate %s configuration: %w", tool, err)
	}

	return output, nil
}

// FormatConfiguration formats Terraform/OpenTofu configuration files
func (m *IacPlanModule) FormatConfiguration(ctx context.Context, workdir string, tool string, check bool) (string, error) {
	var image string
	switch tool {
	case "tofu":
		image = "opentofu/opentofu:latest"
	default:
		image = "hashicorp/terraform:latest"
	}

	container := m.client.Container().
		From(image).
		WithDirectory("/workspace", m.client.Host().Directory(workdir)).
		WithWorkdir("/workspace")

	fmtCmd := []string{tool, "fmt"}
	if check {
		fmtCmd = append(fmtCmd, "-check")
	}

	container = container.WithExec(fmtCmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to format %s configuration: %w", tool, err)
	}

	return output, nil
}

// AnalyzePlan analyzes plan JSON for security and compliance insights
func (m *IacPlanModule) AnalyzePlan(ctx context.Context, planJsonContent string, analysisTypes []string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithNewFile("/plan.json", planJsonContent).
		WithExec([]string{
			"sh", "-c",
			`
				echo "=== IaC Plan Analysis ==="
				echo "Resource changes:"
				jq '.resource_changes | length' /plan.json
				echo "Resources by action:"
				jq '.resource_changes | group_by(.change.actions[0]) | map({action: .[0].change.actions[0], count: length})' /plan.json
				echo "Security findings:"
				jq '.resource_changes[] | select(.type == "aws_s3_bucket" and (.change.after.acl == "public-read" or .change.after.acl == "public-read-write")) | {address, type, public_acl: .change.after.acl}' /plan.json
			`,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to analyze plan: %w", err)
	}

	return output, nil
}

// ComparePlans compares two plan JSON files to show differences
func (m *IacPlanModule) ComparePlans(ctx context.Context, baselinePlan string, currentPlan string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/baseline.json", m.client.Host().File(baselinePlan)).
		WithFile("/current.json", m.client.Host().File(currentPlan)).
		WithExec([]string{
			"sh", "-c",
			`
				echo "=== Plan Comparison ==="
				echo "Baseline resources: $(jq '.resource_changes | length' /baseline.json)"
				echo "Current resources: $(jq '.resource_changes | length' /current.json)"
				echo "Resource diff analysis:"
				jq -n --slurpfile baseline /baseline.json --slurpfile current /current.json '
					{
						baseline_count: ($baseline[0].resource_changes | length),
						current_count: ($current[0].resource_changes | length),
						difference: (($current[0].resource_changes | length) - ($baseline[0].resource_changes | length))
					}
				'
			`,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to compare plans: %w", err)
	}

	return output, nil
}

// ManageWorkspace manages Terraform workspaces
func (m *IacPlanModule) ManageWorkspace(ctx context.Context, workdir string, tool string, operation string, workspaceName string) (string, error) {
	var image string
	switch tool {
	case "tofu":
		image = "opentofu/opentofu:latest"
	default:
		image = "hashicorp/terraform:latest"
	}

	container := m.client.Container().
		From(image).
		WithDirectory("/workspace", m.client.Host().Directory(workdir)).
		WithWorkdir("/workspace").
		WithExec([]string{getBinaryPath(tool), "init"})

	cmd := []string{tool, "workspace", operation}
	if workspaceName != "" {
		cmd = append(cmd, workspaceName)
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to manage %s workspace: %w", tool, err)
	}

	return output, nil
}

// GenerateGraph generates dependency graph from Terraform configuration
func (m *IacPlanModule) GenerateGraph(ctx context.Context, workdir string, tool string, graphType string) (string, error) {
	var image string
	switch tool {
	case "tofu":
		image = "opentofu/opentofu:latest"
	default:
		image = "hashicorp/terraform:latest"
	}

	container := m.client.Container().
		From(image).
		WithDirectory("/workspace", m.client.Host().Directory(workdir)).
		WithWorkdir("/workspace").
		WithExec([]string{getBinaryPath(tool), "init"})

	cmd := []string{tool, "graph"}
	if graphType == "plan" {
		// Generate plan first for plan graph
		container = container.WithExec([]string{tool, "plan", "-out=tfplan"})
		cmd = append(cmd, "tfplan")
	}

	container = container.WithExec(cmd)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate %s graph: %w", tool, err)
	}

	return output, nil
}
