// Package examples demonstrates the convenience functions for the Ship MCP SDK
//
// This example shows the three main usage patterns with convenience functions:
// 1. Pure Ship SDK usage (custom tools only)
// 2. Ship SDK with selected Ship tools using convenience functions
// 3. Ship SDK with all Ship tools using convenience functions

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/cloudshipai/ship/pkg/ship"
	"github.com/cloudshipai/ship/pkg/tools/all"
)

func main() {
	// Example 1: Cherry-pick specific tool collections
	if err := runConvenienceExample(); err != nil {
		log.Fatalf("Convenience example failed: %v", err)
	}

	// Example 2: Use all Ship tools with convenience function
	if err := runAllToolsExample(); err != nil {
		log.Fatalf("All tools example failed: %v", err)
	}

	// Example 3: Mix convenience functions with custom tools
	if err := runMixedExample(); err != nil {
		log.Fatalf("Mixed example failed: %v", err)
	}
}

// Example 1: Cherry-pick specific tool collections
func runConvenienceExample() error {
	fmt.Println("=== Convenience Functions Example ===")

	// Use convenience functions to add specific tool collections
	server := ship.NewServer("convenience-server", "1.0.0")

	// Add just Terraform tools
	all.AddTerraformTools(server)

	// Could also add security tools (when available):
	// all.AddSecurityTools(server)

	// Build the server
	built := server.Build()

	fmt.Printf("✓ Server with Terraform tools: %d tools\n", built.GetRegistry().ToolCount())
	tools := built.GetRegistry().ListTools()
	fmt.Printf("✓ Available tools: %v\n", tools)

	return nil
}

// Example 2: Use all Ship tools with convenience function
func runAllToolsExample() error {
	fmt.Println("\n=== All Tools Example ===")

	// Use convenience function to add all Ship tools at once
	server := all.AddAllTools(
		ship.NewServer("all-tools-server", "1.0.0"),
	).Build()

	fmt.Printf("✓ Server with all Ship tools: %d tools\n", server.GetRegistry().ToolCount())

	// Test that the server works
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer server.Close()

	// Test TFLint tool execution
	tflintTool, err := server.GetRegistry().GetTool("tflint")
	if err != nil {
		return fmt.Errorf("failed to get tflint tool: %w", err)
	}

	result, err := tflintTool.Execute(ctx, map[string]interface{}{
		"directory": ".",
		"format":    "json",
		"init":      false,
	}, server.GetEngine())

	if err != nil {
		return fmt.Errorf("failed to execute tflint: %w", err)
	}

	fmt.Printf("✓ TFLint executed successfully, output length: %d characters\n", len(result.Content))

	return nil
}

// Example 3: Mix convenience functions with custom tools
func runMixedExample() error {
	fmt.Println("\n=== Mixed Tools Example ===")

	// Combine convenience functions with custom tools in a fluent chain
	server := ship.NewServer("mixed-server", "1.0.0").
		// Add custom tool first
		AddContainerTool("project-validator", ship.ContainerToolConfig{
			Description: "Validates project structure and dependencies",
			Image:       "alpine:latest",
			Parameters: []ship.Parameter{
				{
					Name:        "project_type",
					Type:        "string",
					Description: "Type of project to validate",
					Required:    true,
					Enum:        []string{"terraform", "golang", "javascript", "python"},
				},
				{
					Name:        "strict_mode",
					Type:        "boolean",
					Description: "Enable strict validation rules",
					Required:    false,
				},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
				projectType := params["project_type"].(string)
				strictMode := false
				if s, ok := params["strict_mode"].(bool); ok {
					strictMode = s
				}

				// Simulate project validation
				var checks []string
				switch projectType {
				case "terraform":
					checks = []string{"*.tf files", "terraform.tfvars", "providers.tf"}
				case "golang":
					checks = []string{"go.mod", "go.sum", "main.go"}
				case "javascript":
					checks = []string{"package.json", "node_modules", "src/"}
				case "python":
					checks = []string{"requirements.txt", "setup.py", "__init__.py"}
				}

				validation := fmt.Sprintf("Project validation for %s:\n", projectType)
				for _, check := range checks {
					validation += fmt.Sprintf("- %s: ✓\n", check)
				}

				if strictMode {
					validation += "- Strict mode: ✓\n"
				}

				return &ship.ToolResult{
					Content: validation,
					Metadata: map[string]interface{}{
						"project_type": projectType,
						"strict_mode":  strictMode,
						"checks":       len(checks),
						"valid":        true,
					},
				}, nil
			},
		})

	// Chain with convenience functions
	built := all.AddTerraformTools(server).Build()

	fmt.Printf("✓ Mixed server created with %d tools\n", built.GetRegistry().ToolCount())

	// Show all available tools
	tools := built.GetRegistry().ListTools()
	fmt.Printf("✓ Available tools: %v\n", tools)

	// Test custom tool
	ctx := context.Background()
	if err := built.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer built.Close()

	// Test project validator
	validatorTool, err := built.GetRegistry().GetTool("project-validator")
	if err != nil {
		return fmt.Errorf("failed to get project-validator tool: %w", err)
	}

	result, err := validatorTool.Execute(ctx, map[string]interface{}{
		"project_type": "terraform",
		"strict_mode":  true,
	}, built.GetEngine())

	if err != nil {
		return fmt.Errorf("failed to execute project-validator: %w", err)
	}

	fmt.Printf("✓ Project validator executed successfully:\n%s", result.Content)

	return nil
}

// Bonus: Demonstrate creating a focused MCP server for specific workflows
func demonstrateFocusedServer() error {
	fmt.Println("\n=== Focused Server Example ===")

	// Create a server focused on Terraform linting and validation
	server := ship.NewServer("terraform-linter", "1.0.0")

	// Add only Terraform tools - no other distractions
	focusedServer := all.AddTerraformTools(server).
		AddContainerTool("terraform-format", ship.ContainerToolConfig{
			Description: "Format Terraform files using terraform fmt",
			Image:       "hashicorp/terraform:latest",
			Parameters: []ship.Parameter{
				{
					Name:        "directory",
					Type:        "string",
					Description: "Directory to format",
					Required:    false,
				},
				{
					Name:        "check_only",
					Type:        "boolean",
					Description: "Only check if files are formatted, don't modify",
					Required:    false,
				},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
				directory := "."
				if d, ok := params["directory"].(string); ok && d != "" {
					directory = d
				}

				checkOnly := false
				if c, ok := params["check_only"].(bool); ok {
					checkOnly = c
				}

				args := []string{"terraform", "fmt"}
				if checkOnly {
					args = append(args, "-check")
				}
				args = append(args, directory)

				output, err := engine.Container().
					From("hashicorp/terraform:latest").
					WithMountedDirectory("/workspace", directory).
					WithWorkdir("/workspace").
					WithExec(args).
					CombinedOutput(ctx)

				if err != nil {
					return &ship.ToolResult{
						Content: fmt.Sprintf("Terraform fmt failed: %s", output),
						Error:   err,
					}, err
				}

				content := "All Terraform files are properly formatted"
				if output != "" {
					content = output
				}

				return &ship.ToolResult{
					Content: content,
					Metadata: map[string]interface{}{
						"directory":  directory,
						"check_only": checkOnly,
						"tool":       "terraform-fmt",
					},
				}, nil
			},
		}).
		Build()

	fmt.Printf("✓ Focused Terraform server created with %d tools\n", focusedServer.GetRegistry().ToolCount())
	tools := focusedServer.GetRegistry().ListTools()
	fmt.Printf("✓ Terraform workflow tools: %v\n", tools)

	return nil
}