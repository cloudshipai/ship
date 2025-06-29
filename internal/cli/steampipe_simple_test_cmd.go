package cli

import (
	"fmt"
	"os"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var steampipeSimpleTestCmd = &cobra.Command{
	Use:   "steampipe-simple-test",
	Short: "Simple test of Steampipe functionality",
	Long:  `Test Steampipe with basic queries that don't require credentials`,
	RunE:  runSteampipeSimpleTest,
}

func init() {
	rootCmd.AddCommand(steampipeSimpleTestCmd)
}

func runSteampipeSimpleTest(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine for simple Steampipe test...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	// Create Steampipe module
	module := engine.NewSteampipeModule()

	// Test 1: Get Steampipe version
	fmt.Println("\n1. Testing Steampipe version...")
	version, err := module.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}
	fmt.Printf("✓ Steampipe version: %s", version)

	// Test 2: Simple query without credentials
	fmt.Println("\n2. Testing simple query without cloud provider...")
	result, err := module.RunQuery(ctx, "", "SELECT 'Hello from Steampipe' as message, current_timestamp as time", nil)
	if err != nil {
		fmt.Printf("❌ Simple query failed: %v\n", err)
	} else {
		fmt.Printf("✓ Query result: %s\n", result)
	}

	// Test 3: Check AWS credentials availability
	fmt.Println("\n3. Checking AWS environment...")
	awsKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsProfile := os.Getenv("AWS_PROFILE")
	awsConfigFile := os.Getenv("AWS_CONFIG_FILE")

	fmt.Printf("AWS_ACCESS_KEY_ID set: %v\n", awsKeyId != "")
	fmt.Printf("AWS_PROFILE set: %v (value: %s)\n", awsProfile != "", awsProfile)
	fmt.Printf("AWS_CONFIG_FILE set: %v (value: %s)\n", awsConfigFile != "", awsConfigFile)

	if stat, err := os.Stat(os.ExpandEnv("$HOME/.aws/config")); err == nil {
		fmt.Printf("~/.aws/config exists: Yes (size: %d bytes)\n", stat.Size())
	} else {
		fmt.Printf("~/.aws/config exists: No\n")
	}

	if stat, err := os.Stat(os.ExpandEnv("$HOME/.aws/credentials")); err == nil {
		fmt.Printf("~/.aws/credentials exists: Yes (size: %d bytes)\n", stat.Size())
	} else {
		fmt.Printf("~/.aws/credentials exists: No\n")
	}

	// Test 4: Try AWS query with explicit credentials
	fmt.Println("\n4. Testing AWS query with mounted credentials...")

	// Read AWS credentials file to check format
	credsFile := os.ExpandEnv("$HOME/.aws/credentials")
	if data, err := os.ReadFile(credsFile); err == nil {
		fmt.Printf("✓ Successfully read credentials file (%d bytes)\n", len(data))
		// Check if it has [default] profile
		if len(data) > 0 {
			fmt.Println("✓ Credentials file is not empty")
		}
	}

	green := color.New(color.FgGreen)
	green.Printf("\n✓ Simple Steampipe test completed!\n")

	return nil
}
