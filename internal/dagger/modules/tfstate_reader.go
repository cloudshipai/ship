package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TfstateReaderModule reads and analyzes Terraform state files
type TfstateReaderModule struct {
	client *dagger.Client
	name   string
}

const tfstateReaderBinary = "/usr/local/bin/tfstate-lookup"

// NewTfstateReaderModule creates a new Terraform state reader module
func NewTfstateReaderModule(client *dagger.Client) *TfstateReaderModule {
	return &TfstateReaderModule{
		client: client,
		name:   "tfstate-reader",
	}
}

// AnalyzeState analyzes a Terraform state file
func (m *TfstateReaderModule) AnalyzeState(ctx context.Context, statePath string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"jq",
			"{ version: .version, terraform_version: .terraform_version, serial: .serial, lineage: .lineage, resources: [.resources[] | {type: .type, name: .name, provider: .provider, instances: (.instances | length)}] }",
			"/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to analyze state: %w", err)
	}

	return output, nil
}

// ListResources lists resources in state file
func (m *TfstateReaderModule) ListResources(ctx context.Context, statePath string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"jq",
			"[.resources[] | {type: .type, name: .name, mode: .mode, provider: .provider}]",
			"/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list resources: %w", err)
	}

	return output, nil
}

// GetResourceByType gets resources by type
func (m *TfstateReaderModule) GetResourceByType(ctx context.Context, statePath string, resourceType string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"jq",
			fmt.Sprintf(`[.resources[] | select(.type == "%s")]`, resourceType),
			"/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get resources by type: %w", err)
	}

	return output, nil
}

// ExtractOutputs extracts outputs from state file
func (m *TfstateReaderModule) ExtractOutputs(ctx context.Context, statePath string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"jq",
			".outputs",
			"/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to extract outputs: %w", err)
	}

	return output, nil
}

// ShowState shows state information using terraform show
func (m *TfstateReaderModule) ShowState(ctx context.Context, statePath string) (string, error) {
	container := m.client.Container().
		From("hashicorp/terraform:latest").
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"terraform", "show", "-json", "/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to show state: %w", err)
	}

	return output, nil
}

// PullRemoteState pulls state from remote backend
func (m *TfstateReaderModule) PullRemoteState(ctx context.Context, workdir string) (string, error) {
	container := m.client.Container().
		From("hashicorp/terraform:latest").
		WithDirectory("/workspace", m.client.Host().Directory(workdir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terraform", "init"}).
		WithExec([]string{
			"terraform", "state", "pull",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to pull remote state: %w", err)
	}

	return output, nil
}

// InteractiveExplorer provides interactive state exploration
func (m *TfstateReaderModule) InteractiveExplorer(ctx context.Context, statePath string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"sh", "-c", "echo 'Interactive exploration:' && jq '.' /terraform.tfstate | head -50",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to explore state interactively: %w", err)
	}

	return output, nil
}

// StateListResources lists resources using terraform state list
func (m *TfstateReaderModule) StateListResources(ctx context.Context, statePath string, resourceId string) (string, error) {
	container := m.client.Container().
		From("hashicorp/terraform:latest").
		WithFile("/terraform.tfstate", m.client.Host().File(statePath))

	args := []string{"terraform", "state", "list", "-state=/terraform.tfstate"}
	if resourceId != "" {
		args = append(args, resourceId)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list state resources: %w", err)
	}

	return output, nil
}

// StateShowResource shows a specific resource using terraform state show
func (m *TfstateReaderModule) StateShowResource(ctx context.Context, statePath string, resourceAddress string) (string, error) {
	container := m.client.Container().
		From("hashicorp/terraform:latest").
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"terraform", "state", "show", "-state=/terraform.tfstate", resourceAddress,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to show state resource: %w", err)
	}

	return output, nil
}

// LookupResource looks up resource using tfstate-lookup functionality (simulated with jq)
func (m *TfstateReaderModule) LookupResource(ctx context.Context, statePath string, resourceAddress string) (string, error) {
	// Since we don't have tfstate-lookup in container, simulate with jq
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"jq",
			fmt.Sprintf(`[.resources[] | select(.type + "." + .name == "%s") | .instances[]]`, resourceAddress),
			"/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to lookup resource: %w", err)
	}

	return output, nil
}

// DumpAllResources dumps all resources (simulated tfstate-lookup -dump)
func (m *TfstateReaderModule) DumpAllResources(ctx context.Context, statePath string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithFile("/terraform.tfstate", m.client.Host().File(statePath)).
		WithExec([]string{
			"jq",
			`{
				version: .version,
				terraform_version: .terraform_version,
				serial: .serial,
				lineage: .lineage,
				outputs: .outputs,
				resources: [.resources[] | {
					type: .type,
					name: .name,
					provider: .provider,
					instances: [.instances[] | {
						attributes: .attributes,
						dependencies: .dependencies
					}]
				}]
			}`,
			"/terraform.tfstate",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to dump all resources: %w", err)
	}

	return output, nil
}
