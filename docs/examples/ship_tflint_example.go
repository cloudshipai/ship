//go:build ignore
// +build ignore

// Package examples demonstrates how to use the Ship MCP SDK
//
// This example shows how to create an MCP server using the Ship SDK
// with the TFLint tool integrated. This demonstrates the three usage patterns:
// 1. Pure Ship SDK usage (custom tools only)
// 2. Ship SDK with selected Ship tools (cherry-pick)
// 3. Ship SDK with all Ship tools (everything plus)

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudshipai/ship/pkg/ship"
	"github.com/cloudshipai/ship/pkg/tools"
	"github.com/cloudshipai/ship/pkg/dagger"
)

func main() {
	// Example 1: Pure Ship SDK usage with custom tools only
	if err := runPureFrameworkExample(); err != nil {
		log.Fatalf("Pure Ship SDK example failed: %v", err)
	}

	// Example 2: Ship SDK with selected Ship tools (cherry-pick)
	if err := runCherryPickExample(); err != nil {
		log.Fatalf("Cherry-pick example failed: %v", err)
	}

	// Example 3: Ship SDK with all Ship tools (everything plus)
	if err := runEverythingPlusExample(); err != nil {
		log.Fatalf("Everything plus example failed: %v", err)
	}

	// Bonus: Demonstrate stdio usage
	if err := demonstrateStdioUsage(); err != nil {
		log.Fatalf("Stdio demo failed: %v", err)
	}
}

// Example 1: Pure Ship SDK usage with custom tools only
func runPureFrameworkExample() error {
	fmt.Println("=== Pure Ship SDK Example ===")
	
	// Create custom tool using Ship SDK
	customTool := ship.NewContainerTool("echo-tool", ship.ContainerToolConfig{
		Description: "Simple echo tool for testing",
		Image:       "alpine:latest",
		Parameters: []ship.Parameter{
			{
				Name:        "message",
				Type:        "string",
				Description: "Message to echo",
				Required:    true,
			},
		},
		Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
			message := params["message"].(string)
			
			result, err := engine.Container().
				From("alpine:latest").
				WithExec([]string{"echo", message}).
				Stdout(ctx)
			
			if err != nil {
				return &ship.ToolResult{Error: err}, err
			}
			
			return &ship.ToolResult{
				Content: result,
				Metadata: map[string]interface{}{
					"tool": "echo-tool",
					"message": message,
				},
			}, nil
		},
	})

	// Build server with only custom tools
	server := ship.NewServer("pure-framework-server", "1.0.0").
		AddTool(customTool).
		Build()

	fmt.Printf("✓ Pure Ship SDK server created with %d tools\n", server.GetRegistry().ToolCount())
	return nil
}

// Example 2: Ship SDK with selected Ship tools (cherry-pick)
func runCherryPickExample() error {
	fmt.Println("\n=== Cherry-Pick Example ===")
	
	// Create server with selected Ship tools
	server := ship.NewServer("cherry-pick-server", "1.0.0").
		AddTool(tools.NewTFLintTool()).  // Add TFLint from Ship tools
		AddContainerTool("custom-validator", ship.ContainerToolConfig{
			Description: "Custom validation tool",
			Image:       "alpine:latest",
			Parameters: []ship.Parameter{
				{
					Name:        "file",
					Type:        "string",
					Description: "File to validate",
					Required:    true,
				},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
				file := params["file"].(string)
				
				// Simple file existence check
				_, err := engine.Container().
					From("alpine:latest").
					WithMountedDirectory("/workspace", ".").
					WithExec([]string{"test", "-f", fmt.Sprintf("/workspace/%s", file)}).
					CombinedOutput(ctx)
				
				if err != nil {
					return &ship.ToolResult{
						Content: fmt.Sprintf("File %s does not exist", file),
						Metadata: map[string]interface{}{"valid": false},
					}, nil
				}
				
				return &ship.ToolResult{
					Content: fmt.Sprintf("File %s exists and is valid", file),
					Metadata: map[string]interface{}{"valid": true},
				}, nil
			},
		}).
		Build()

	fmt.Printf("✓ Cherry-pick server created with %d tools\n", server.GetRegistry().ToolCount())
	
	// Test TFLint tool
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer server.Close()

	// Get TFLint tool and test execution
	tflintTool, err := server.GetRegistry().GetTool("tflint")
	if err != nil {
		return fmt.Errorf("failed to get tflint tool: %w", err)
	}

	result, err := tflintTool.Execute(ctx, map[string]interface{}{
		"directory": ".",
		"format":    "json",
	}, server.GetEngine())
	
	if err != nil {
		return fmt.Errorf("failed to execute tflint: %w", err)
	}

	fmt.Printf("✓ TFLint executed successfully, output length: %d characters\n", len(result.Content))
	return nil
}

// Example 3: Ship SDK with all Ship tools (everything plus)
func runEverythingPlusExample() error {
	fmt.Println("\n=== Everything Plus Example ===")
	
	// Create registry with all Ship tools
	shipRegistry := ship.NewRegistry()
	shipRegistry.RegisterTool(tools.NewTFLintTool())
	// In a real implementation, you would add all Ship tools here:
	// shipRegistry.RegisterTool(tools.NewCheckovTool())
	// shipRegistry.RegisterTool(tools.NewCostAnalysisTool())
	// etc.

	// Create server with all Ship tools plus custom extensions
	everythingServer := ship.NewServer("everything-plus-server", "1.0.0").
		ImportRegistry(shipRegistry).  // Import all Ship tools
		AddContainerTool("deployment-validator", ship.ContainerToolConfig{
			Description: "Advanced deployment validation",
			Image:       "alpine:latest",
			Parameters: []ship.Parameter{
				{
					Name:        "environment",
					Type:        "string",
					Description: "Target environment",
					Required:    true,
					Enum:        []string{"dev", "staging", "prod"},
				},
				{
					Name:        "strict",
					Type:        "boolean",
					Description: "Enable strict validation mode",
					Required:    false,
				},
			},
			Execute: func(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
				env := params["environment"].(string)
				strict := false
				if s, ok := params["strict"].(bool); ok {
					strict = s
				}
				
				// Simulate deployment validation
				validation := fmt.Sprintf("Deployment validation for %s environment", env)
				if strict {
					validation += " (strict mode enabled)"
				}
				
				return &ship.ToolResult{
					Content: validation,
					Metadata: map[string]interface{}{
						"environment": env,
						"strict":      strict,
						"valid":       true,
					},
				}, nil
			},
		}).
		Build()

	fmt.Printf("✓ Everything plus server created with %d tools\n", everythingServer.GetRegistry().ToolCount())
	
	// Show that we can list all tools
	tools := everythingServer.GetRegistry().ListTools()
	fmt.Printf("✓ Available tools: %v\n", tools)
	
	return nil
}

// Bonus: Example of using the server in stdio mode
func demonstrateStdioUsage() error {
	fmt.Println("\n=== Stdio Server Example ===")
	
	stdioServer := ship.NewServer("stdio-demo", "1.0.0").
		AddTool(tools.NewTFLintTool()).
		Build()

	fmt.Printf("✓ Server ready for stdio mode with %d tools\n", stdioServer.GetRegistry().ToolCount())
	fmt.Println("  To run in stdio mode, call: stdioServer.ServeStdio()")
	fmt.Println("  This would start the MCP server and communicate via stdin/stdout")
	
	// Note: We don't actually call ServeStdio() here as it would block
	// In a real MCP server binary, you would call:
	// return stdioServer.ServeStdio()
	
	return nil
}