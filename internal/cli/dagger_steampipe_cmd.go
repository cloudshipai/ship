package cli

import (
	"fmt"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var daggerSteampipeTestCmd = &cobra.Command{
	Use:   "dagger-steampipe-test",
	Short: "Test Dagger with Steampipe container",
	Long:  `Test that Dagger can run Steampipe containers`,
	RunE:  runDaggerSteampipeTest,
}

func init() {
	rootCmd.AddCommand(daggerSteampipeTestCmd)
}

func runDaggerSteampipeTest(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine for Steampipe test...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Test running Steampipe container
	fmt.Println("\nTesting Steampipe container...")
	client := engine.GetClient()

	// Run Steampipe with version command
	container := client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to run steampipe container: %w", err)
	}

	fmt.Printf("Steampipe version:\n%s\n", output)

	// Test listing available plugins
	fmt.Println("\nListing available Steampipe plugins...")
	pluginContainer := client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "plugin", "list", "--output", "json"})

	pluginOutput, err := pluginContainer.Stdout(ctx)
	if err != nil {
		// This might fail if no plugins are installed, which is OK
		fmt.Printf("Note: No plugins installed (this is expected): %v\n", err)
	} else {
		fmt.Printf("Installed plugins:\n%s\n", pluginOutput)
	}

	green := color.New(color.FgGreen)
	green.Printf("\n✓ Steampipe container is working correctly!\n")
	fmt.Println("\nThis demonstrates that Ship CLI can run Steampipe")
	fmt.Println("in containers without requiring local installation.")

	return nil
}
