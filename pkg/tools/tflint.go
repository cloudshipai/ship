package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/cloudshipai/ship/internal/ship"
)

const (
	TFLintImage = "ghcr.io/terraform-linters/tflint:latest"
)

// NewTFLintTool creates a new TFLint tool using the framework
func NewTFLintTool() ship.Tool {
	config := ship.ContainerToolConfig{
		Description: "Run TFLint on Terraform code to check for syntax errors and best practices",
		Image:       TFLintImage,
		Parameters: []ship.Parameter{
			{
				Name:        "directory",
				Type:        "string",
				Description: "Directory containing Terraform files (default: current directory)",
				Required:    false,
			},
			{
				Name:        "format",
				Type:        "string",
				Description: "Output format: default, json, compact",
				Required:    false,
				Enum:        []string{"default", "json", "compact"},
			},
			{
				Name:        "config",
				Type:        "string",
				Description: "Path to TFLint configuration file",
				Required:    false,
			},
			{
				Name:        "enable_rules",
				Type:        "string",
				Description: "Comma-separated list of rules to enable",
				Required:    false,
			},
			{
				Name:        "disable_rules",
				Type:        "string",
				Description: "Comma-separated list of rules to disable",
				Required:    false,
			},
			{
				Name:        "init",
				Type:        "boolean",
				Description: "Run tflint --init before linting to install plugins",
				Required:    false,
			},
		},
		Execute: executeTFLint,
	}

	return ship.NewContainerTool("tflint", config)
}

func executeTFLint(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	// Get parameters with defaults
	directory := getStringParam(params, "directory", ".")
	format := getStringParam(params, "format", "default")
	configFile := getStringParam(params, "config", "")
	enableRules := getStringParam(params, "enable_rules", "")
	disableRules := getStringParam(params, "disable_rules", "")
	shouldInit := getBoolParam(params, "init", true)

	// Validate directory exists
	if directory == "" {
		directory = "."
	}

	// Build container
	container := engine.Container().
		From(TFLintImage).
		WithMountedDirectory("/workspace", directory).
		WithWorkdir("/workspace")

	// Note: Custom config file mounting would need to be implemented
	// For now, we'll rely on the config file being in the mounted directory
	if configFile != "" {
		// The config file should be in the mounted workspace directory
		// We'll reference it by its relative path in the workspace
	}

	// Initialize TFLint plugins if requested
	if shouldInit {
		container = container.WithExec([]string{"tflint", "--init"})
	}

	// Build TFLint command arguments
	args := []string{"tflint"}

	// Add format
	if format != "default" {
		args = append(args, "--format", format)
	}

	// Add custom config (relative to workspace)
	if configFile != "" {
		args = append(args, "--config", configFile)
	}

	// Add enabled rules
	if enableRules != "" {
		rules := strings.Split(enableRules, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if rule != "" {
				args = append(args, "--enable-rule", rule)
			}
		}
	}

	// Add disabled rules
	if disableRules != "" {
		rules := strings.Split(disableRules, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if rule != "" {
				args = append(args, "--disable-rule", rule)
			}
		}
	}

	// Use shell wrapper to capture output regardless of exit code
	// TFLint returns non-zero exit codes when issues are found
	shellCmd := fmt.Sprintf("%s || true", strings.Join(args, " "))
	container = container.WithExec([]string{"sh", "-c", shellCmd})

	// Execute and get output
	output, err := container.Stdout(ctx)
	if err != nil {
		return &ship.ToolResult{
			Content: "",
			Error:   fmt.Errorf("failed to run tflint: %w", err),
		}, err
	}

	// Get stderr for any warnings/errors
	stderr, _ := container.Stderr(ctx)

	// Combine output and stderr if both exist
	finalOutput := output
	if stderr != "" {
		finalOutput = fmt.Sprintf("STDOUT:\n%s\n\nSTDERR:\n%s", output, stderr)
	}

	return &ship.ToolResult{
		Content: finalOutput,
		Metadata: map[string]interface{}{
			"format":      format,
			"directory":   directory,
			"config_file": configFile,
			"initialized": shouldInit,
			"tool":        "tflint",
			"image":       TFLintImage,
		},
	}, nil
}

// Helper functions for parameter extraction
func getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, exists := params[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, exists := params[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
		// Handle string representation of boolean
		if str, ok := val.(string); ok {
			return strings.ToLower(str) == "true"
		}
	}
	return defaultValue
}
