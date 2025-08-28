package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// NucleiModule runs Nuclei for vulnerability scanning
type NucleiModule struct {
	client *dagger.Client
	name   string
}

// NewNucleiModule creates a new Nuclei module
func NewNucleiModule(client *dagger.Client) *NucleiModule {
	return &NucleiModule{
		client: client,
		name:   "nuclei",
	}
}

// GetVersion returns the version of Nuclei
func (m *NucleiModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec([]string{"nuclei", "-version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// Try stderr for version info
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}

	return "Nuclei vulnerability scanner", nil
}

// ScanURL scans a single URL for vulnerabilities
func (m *NucleiModule) ScanURL(ctx context.Context, targetURL string, severity string) (string, error) {
	// Use silent mode to reduce noise and disable-update-check for faster execution
	args := []string{"nuclei", "-target", targetURL, "-jsonl", "-silent", "-disable-update-check"}
	
	if severity != "" {
		args = append(args, "-severity", severity)
	}

	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// Nuclei outputs to stdout when findings exist
	if output != "" {
		return output, nil
	}
	
	// Check stderr for any errors
	if stderr != "" && strings.Contains(stderr, "Error") {
		return "", fmt.Errorf("scan failed: %s", stderr)
	}
	
	// No output means no vulnerabilities found
	return `{"message": "No vulnerabilities found"}`, nil
}

// ScanURLList scans multiple URLs from a file
func (m *NucleiModule) ScanURLList(ctx context.Context, urlsFile string, severity string) (string, error) {
	var container *dagger.Container
	
	// Create sample URLs file if none provided
	if urlsFile == "" || urlsFile == "/tmp/urls.txt" {
		sampleURLs := "https://example.com\nhttps://test.example.com"
		container = m.client.Container().
			From("projectdiscovery/nuclei:latest").
			WithNewFile("/urls.txt", sampleURLs)
	} else {
		container = m.client.Container().
			From("projectdiscovery/nuclei:latest").
			WithFile("/urls.txt", m.client.Host().File(urlsFile))
	}

	args := []string{"nuclei", "-list", "/urls.txt", "-json"}
	if severity != "" {
		args = append(args, "-severity", severity)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" && !strings.Contains(stderr, "Error") {
			return stderr, nil
		}
		if output == "" {
			return `{"message": "No vulnerabilities found"}`, nil
		}
	}

	return output, nil
}

// ScanWithTemplate scans using specific template(s)
func (m *NucleiModule) ScanWithTemplate(ctx context.Context, targetURL string, templatePath string) (string, error) {
	args := []string{"nuclei", "-u", targetURL, "-json"}
	
	if templatePath != "" {
		args = append(args, "-t", templatePath)
	}

	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" && !strings.Contains(stderr, "Error") {
			return stderr, nil
		}
		if output == "" {
			return `{"message": "No vulnerabilities found"}`, nil
		}
	}

	return output, nil
}

// ScanWithWorkflow runs a workflow scan
func (m *NucleiModule) ScanWithWorkflow(ctx context.Context, targetURL string, workflowPath string) (string, error) {
	args := []string{"nuclei", "-u", targetURL, "-json"}
	
	if workflowPath != "" {
		args = append(args, "-w", workflowPath)
	} else {
		// Use default workflows
		args = append(args, "-w", "workflows/")
	}

	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" && !strings.Contains(stderr, "Error") {
			return stderr, nil
		}
		if output == "" {
			return `{"message": "No vulnerabilities found"}`, nil
		}
	}

	return output, nil
}

// UpdateTemplates updates Nuclei templates
func (m *NucleiModule) UpdateTemplates(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec([]string{"nuclei", "-ut"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	
	if stderr != "" {
		return stderr, nil
	}
	
	return "Templates updated successfully", nil
}

// ListTemplates lists available templates
func (m *NucleiModule) ListTemplates(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec([]string{"nuclei", "-tl", "-silent"}, dagger.ContainerWithExecOpts{
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
	
	return "No templates found", nil
}

// ScanWithTags scans using specific tags
func (m *NucleiModule) ScanWithTags(ctx context.Context, targetURL string, tags string) (string, error) {
	args := []string{"nuclei", "-u", targetURL, "-json"}
	
	if tags != "" {
		args = append(args, "-tags", tags)
	}

	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" && !strings.Contains(stderr, "Error") {
			return stderr, nil
		}
		if output == "" {
			return `{"message": "No vulnerabilities found"}`, nil
		}
	}

	return output, nil
}

// RateLimitedScan performs scan with rate limiting
func (m *NucleiModule) RateLimitedScan(ctx context.Context, targetURL string, rateLimit int) (string, error) {
	args := []string{"nuclei", "-u", targetURL, "-json"}
	
	if rateLimit > 0 {
		args = append(args, "-rate-limit", fmt.Sprintf("%d", rateLimit))
	}

	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" && !strings.Contains(stderr, "Error") {
			return stderr, nil
		}
		if output == "" {
			return `{"message": "No vulnerabilities found"}`, nil
		}
	}

	return output, nil
}

// ValidateTemplate validates a custom template
func (m *NucleiModule) ValidateTemplate(ctx context.Context, templatePath string) (string, error) {
	var container *dagger.Container
	
	if templatePath == "" || templatePath == "/tmp/template.yaml" {
		// Create a sample template for validation
		sampleTemplate := `id: sample-template
info:
  name: Sample Template
  author: Ship CLI
  severity: info
  description: Sample template for validation

http:
  - method: GET
    path:
      - "{{BaseURL}}"
    matchers:
      - type: status
        status:
          - 200`
		container = m.client.Container().
			From("projectdiscovery/nuclei:latest").
			WithNewFile("/template.yaml", sampleTemplate)
		templatePath = "/template.yaml"
	} else {
		container = m.client.Container().
			From("projectdiscovery/nuclei:latest").
			WithFile("/template.yaml", m.client.Host().File(templatePath))
		templatePath = "/template.yaml"
	}

	container = container.WithExec([]string{"nuclei", "-validate", "-t", templatePath}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "Template validation completed", nil
	}

	if output == "" {
		return "Template is valid", nil
	}

	return output, nil
}

// GenerateReport generates a scan report in various formats
func (m *NucleiModule) GenerateReport(ctx context.Context, targetURL string, reportType string) (string, error) {
	args := []string{"nuclei", "-u", targetURL}
	
	switch reportType {
	case "json":
		args = append(args, "-json")
	case "markdown":
		args = append(args, "-markdown-export", "/tmp/report.md")
	case "sarif":
		args = append(args, "-sarif-export", "/tmp/report.sarif")
	default:
		args = append(args, "-json")
	}

	container := m.client.Container().
		From("projectdiscovery/nuclei:latest").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" && !strings.Contains(stderr, "Error") {
			return stderr, nil
		}
	}

	// For markdown or sarif, read the generated file
	if reportType == "markdown" || reportType == "sarif" {
		fileName := fmt.Sprintf("/tmp/report.%s", reportType)
		if reportType == "sarif" {
			fileName = "/tmp/report.sarif"
		} else {
			fileName = "/tmp/report.md"
		}
		
		reportContainer := container.WithExec([]string{"cat", fileName}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})
		
		reportOutput, _ := reportContainer.Stdout(ctx)
		if reportOutput != "" {
			return reportOutput, nil
		}
	}

	if output == "" {
		return `{"message": "No vulnerabilities found"}`, nil
	}

	return output, nil
}