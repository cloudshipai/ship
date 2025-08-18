package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// OpenSCAPModule runs OpenSCAP for security compliance scanning
type OpenSCAPModule struct {
	client *dagger.Client
	name   string
}

const oscapBinary = "/usr/bin/oscap"

// NewOpenSCAPModule creates a new OpenSCAP module
func NewOpenSCAPModule(client *dagger.Client) *OpenSCAPModule {
	return &OpenSCAPModule{
		client: client,
		name:   "openscap",
	}
}

// EvaluateProfile evaluates a system against SCAP content
func (m *OpenSCAPModule) EvaluateProfile(ctx context.Context, contentPath string, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/content.xml", m.client.Host().File(contentPath)).
		WithExec([]string{
			oscapBinary,
			"xccdf", "eval",
			"--profile", profile,
			"--results", "/results.xml",
			"--report", "/report.html",
			"/content.xml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate profile: %w", err)
	}

	return output, nil
}

// ScanImage scans a container image for compliance
func (m *OpenSCAPModule) ScanImage(ctx context.Context, imageName string, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithExec([]string{
			"oscap-podman",
			imageName,
			"xccdf", "eval",
			"--profile", profile,
			"--report", "/report.html",
			"--results", "/results.xml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan image: %w", err)
	}

	return output, nil
}

// GenerateReport generates compliance report
func (m *OpenSCAPModule) GenerateReport(ctx context.Context, resultsPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/results.xml", m.client.Host().File(resultsPath)).
		WithExec([]string{
			oscapBinary,
			"xccdf", "generate", "report",
			"/results.xml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	return output, nil
}

// OvalEvaluate evaluates OVAL definitions
func (m *OpenSCAPModule) OvalEvaluate(ctx context.Context, ovalFile string, resultsFile string, variablesFile string, definitionId string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/oval.xml", m.client.Host().File(ovalFile))

	args := []string{oscapBinary, "oval", "eval"}
	if resultsFile != "" {
		args = append(args, "--results", "/results.xml")
	}
	if variablesFile != "" {
		container = container.WithFile("/variables.xml", m.client.Host().File(variablesFile))
		args = append(args, "--variables", "/variables.xml")
	}
	if definitionId != "" {
		args = append(args, "--id", definitionId)
	}
	args = append(args, "/oval.xml")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate OVAL: %w", err)
	}

	return output, nil
}

// GenerateGuide generates HTML guide from XCCDF content
func (m *OpenSCAPModule) GenerateGuide(ctx context.Context, xccdfFile string, profile string, outputFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/xccdf.xml", m.client.Host().File(xccdfFile))

	args := []string{oscapBinary, "xccdf", "generate", "guide"}
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	args = append(args, "/xccdf.xml")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate guide: %w", err)
	}

	return output, nil
}

// ValidateDataStream validates Source DataStream file
func (m *OpenSCAPModule) ValidateDataStream(ctx context.Context, datastreamFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/datastream.xml", m.client.Host().File(datastreamFile)).
		WithExec([]string{oscapBinary, "ds", "sds-validate", "/datastream.xml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate datastream: %w", err)
	}

	return output, nil
}

// ValidateContent validates SCAP content
func (m *OpenSCAPModule) ValidateContent(ctx context.Context, contentFile string, contentType string, schematron bool) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/content.xml", m.client.Host().File(contentFile))

	var args []string
	if contentType != "" {
		args = []string{oscapBinary, contentType, "validate"}
		if contentType == "oval" && schematron {
			args = append(args, "--schematron")
		}
	} else {
		args = []string{oscapBinary, "info", "/content.xml"}
	}
	args = append(args, "/content.xml")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate content: %w", err)
	}

	return output, nil
}

// GetInfo displays information about SCAP content
func (m *OpenSCAPModule) GetInfo(ctx context.Context, contentFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/content.xml", m.client.Host().File(contentFile)).
		WithExec([]string{oscapBinary, "info", "/content.xml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get content info: %w", err)
	}

	return output, nil
}

// RemediateXCCDF applies remediation based on XCCDF results
func (m *OpenSCAPModule) RemediateXCCDF(ctx context.Context, resultsFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/results.xml", m.client.Host().File(resultsFile)).
		WithExec([]string{oscapBinary, "xccdf", "remediate", "/results.xml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to remediate: %w", err)
	}

	return output, nil
}

// GenerateOvalReport generates report from OVAL results
func (m *OpenSCAPModule) GenerateOvalReport(ctx context.Context, ovalResultsFile string, outputFile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/oval_results.xml", m.client.Host().File(ovalResultsFile)).
		WithExec([]string{oscapBinary, "oval", "generate", "report", "/oval_results.xml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate oval report: %w", err)
	}

	return output, nil
}

// SplitDataStream splits DataStream into component files
func (m *OpenSCAPModule) SplitDataStream(ctx context.Context, datastreamFile string, outputDir string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/datastream.xml", m.client.Host().File(datastreamFile))

	args := []string{oscapBinary, "ds", "sds-split"}
	if outputDir != "" {
		args = append(args, "--output-dir", "/output")
	}
	args = append(args, "/datastream.xml")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to split datastream: %w", err)
	}

	return output, nil
}
