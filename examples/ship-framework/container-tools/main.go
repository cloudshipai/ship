package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/cloudshipai/ship/pkg/ship"
)

// TerraformValidateTool demonstrates a realistic containerized tool
type TerraformValidateTool struct{}

func (t *TerraformValidateTool) Name() string {
	return "terraform-validate"
}

func (t *TerraformValidateTool) Description() string {
	return "Validates Terraform configuration files using hashicorp/terraform container"
}

func (t *TerraformValidateTool) Parameters() []ship.Parameter {
	return []ship.Parameter{
		{
			Name:        "directory",
			Type:        "string",
			Description: "Directory containing Terraform files (default: current directory)",
			Required:    false,
		},
		{
			Name:        "json_output",
			Type:        "boolean", 
			Description: "Output results in JSON format",
			Required:    false,
		},
	}
}

func (t *TerraformValidateTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	directory := "."
	if d, ok := params["directory"].(string); ok && d != "" {
		directory = d
	}

	jsonOutput, _ := params["json_output"].(bool)

	if engine == nil {
		return &ship.ToolResult{
			Error: fmt.Errorf("dagger engine not available"),
		}, fmt.Errorf("dagger engine not available")
	}

	// Build command arguments
	args := []string{"terraform", "validate"}
	if jsonOutput {
		args = append(args, "-json")
	}

	// Mount the directory and run terraform validate
	container := engine.Container().
		From("hashicorp/terraform:latest").
		WithMountedDirectory("/workspace", directory).
		WithWorkdir("/workspace").
		WithExec([]string{"terraform", "init", "-backend=false"}).  // Initialize without backend
		WithExec(args)

	stdout, err := container.Stdout(ctx)
	if err != nil {
		// Try to get stderr for better error messages
		stderr, _ := container.Stderr(ctx)
		return &ship.ToolResult{
			Content: fmt.Sprintf("Terraform validation failed:\nSTDOUT: %s\nSTDERR: %s", stdout, stderr),
			Error:   fmt.Errorf("terraform validation failed: %w", err),
		}, fmt.Errorf("terraform validation failed: %w", err)
	}

	return &ship.ToolResult{
		Content: fmt.Sprintf("Terraform validation for %s:\n%s", directory, stdout),
		Metadata: map[string]interface{}{
			"directory":   directory,
			"json_output": jsonOutput,
			"tool":        "terraform",
			"command":     strings.Join(args, " "),
		},
	}, nil
}

// DockerLintTool demonstrates another practical containerized tool
type DockerLintTool struct{}

func (t *DockerLintTool) Name() string {
	return "dockerfile-lint"
}

func (t *DockerLintTool) Description() string {
	return "Lints Dockerfile using hadolint container"
}

func (t *DockerLintTool) Parameters() []ship.Parameter {
	return []ship.Parameter{
		{
			Name:        "dockerfile_path",
			Type:        "string",
			Description: "Path to Dockerfile (default: ./Dockerfile)",
			Required:    false,
		},
		{
			Name:        "format",
			Type:        "string",
			Description: "Output format (checkstyle, codeclimate, gcc, json, tty)",
			Required:    false,
			Enum:        []string{"checkstyle", "codeclimate", "gcc", "json", "tty"},
		},
	}
}

func (t *DockerLintTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	dockerfilePath := "./Dockerfile"
	if p, ok := params["dockerfile_path"].(string); ok && p != "" {
		dockerfilePath = p
	}

	format := "tty"
	if f, ok := params["format"].(string); ok && f != "" {
		format = f
	}

	if engine == nil {
		return &ship.ToolResult{
			Error: fmt.Errorf("dagger engine not available"),
		}, fmt.Errorf("dagger engine not available")
	}

	// Build hadolint command
	args := []string{"hadolint", "--format", format, dockerfilePath}

	// Mount directory and run hadolint
	container := engine.Container().
		From("hadolint/hadolint:latest").
		WithMountedDirectory("/workspace", ".").
		WithWorkdir("/workspace").
		WithExec(args)

	stdout, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return &ship.ToolResult{
			Content: fmt.Sprintf("Dockerfile linting completed with issues:\nSTDOUT: %s\nSTDERR: %s", stdout, stderr),
			Metadata: map[string]interface{}{
				"dockerfile_path": dockerfilePath,
				"format":          format,
				"had_issues":      true,
			},
		}, nil // Don't return error for linting issues, just warnings
	}

	return &ship.ToolResult{
		Content: fmt.Sprintf("Dockerfile linting results:\n%s", stdout),
		Metadata: map[string]interface{}{
			"dockerfile_path": dockerfilePath,
			"format":          format,
			"had_issues":      false,
		},
	}, nil
}

// YAMLLintTool shows how to create a flexible container tool
type YAMLLintTool struct{}

func (t *YAMLLintTool) Name() string {
	return "yaml-lint"
}

func (t *YAMLLintTool) Description() string {
	return "Validates YAML files using yamllint"
}

func (t *YAMLLintTool) Parameters() []ship.Parameter {
	return []ship.Parameter{
		{
			Name:        "files",
			Type:        "string",
			Description: "Files or directories to lint (space-separated, default: .)",
			Required:    false,
		},
		{
			Name:        "strict",
			Type:        "boolean",
			Description: "Return non-zero exit code on warnings",
			Required:    false,
		},
	}
}

func (t *YAMLLintTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	files := "."
	if f, ok := params["files"].(string); ok && f != "" {
		files = f
	}

	strict, _ := params["strict"].(bool)

	if engine == nil {
		return &ship.ToolResult{
			Error: fmt.Errorf("dagger engine not available"),
		}, fmt.Errorf("dagger engine not available")
	}

	// Build yamllint command
	args := []string{"yamllint"}
	if strict {
		args = append(args, "--strict")
	}
	
	// Split files and add them
	fileList := strings.Fields(files)
	args = append(args, fileList...)

	container := engine.Container().
		From("cytopia/yamllint:latest").
		WithMountedDirectory("/workspace", ".").
		WithWorkdir("/workspace").
		WithExec(args)

	stdout, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return &ship.ToolResult{
			Content: fmt.Sprintf("YAML linting found issues:\nSTDOUT: %s\nSTDERR: %s", stdout, stderr),
			Metadata: map[string]interface{}{
				"files":      files,
				"strict":     strict,
				"had_issues": true,
			},
		}, nil // Return success even with linting issues
	}

	return &ship.ToolResult{
		Content: fmt.Sprintf("YAML linting results:\n%s", stdout),
		Metadata: map[string]interface{}{
			"files":      files,
			"strict":     strict,
			"had_issues": false,
		},
	}, nil
}

func main() {
	// Build MCP server with containerized tools
	server := ship.NewServer("container-tools-server", "1.0.0").
		AddTool(&TerraformValidateTool{}).
		AddTool(&DockerLintTool{}).
		AddTool(&YAMLLintTool{}).
		Build()

	fmt.Fprintf(log.Writer(), "Starting container tools MCP server...\n")
	fmt.Fprintf(log.Writer(), "Available tools: terraform-validate, dockerfile-lint, yaml-lint\n")
	fmt.Fprintf(log.Writer(), "All tools run in isolated Docker containers\n")

	// Start the MCP server
	if err := server.ServeStdio(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}