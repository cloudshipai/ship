package cli

import (
	"fmt"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var daggerTestCmd = &cobra.Command{
	Use:   "dagger-test",
	Short: "Test Dagger engine integration",
	Long:  `Test that Dagger engine can be initialized and run simple containers`,
	RunE:  runDaggerTest,
}

func init() {
	rootCmd.AddCommand(daggerTestCmd)
}

func runDaggerTest(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Run a simple container
	fmt.Println("Running test container...")
	output, err := engine.RunContainer("alpine:latest", []string{"echo", "Hello from Dagger!"})
	if err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	fmt.Printf("Container output: %s", output)

	// Test building a container with Go
	fmt.Println("\nTesting Go build container...")
	goOutput, err := engine.RunContainer("golang:1.23-alpine", []string{"go", "version"})
	if err != nil {
		return fmt.Errorf("failed to run go container: %w", err)
	}

	fmt.Printf("Go version: %s", goOutput)

	green := color.New(color.FgGreen)
	green.Printf("\n✓ Dagger engine is working correctly!\n")

	return nil
}
