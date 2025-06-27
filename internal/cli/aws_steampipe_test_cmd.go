package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudship/ship/internal/dagger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var awsSteampipeTestCmd = &cobra.Command{
	Use:   "aws-steampipe-test",
	Short: "Test Steampipe AWS connection with proper configuration",
	Long:  `Test Steampipe with AWS using various credential methods`,
	RunE:  runAWSSteampipeTest,
}

func init() {
	rootCmd.AddCommand(awsSteampipeTestCmd)
	awsSteampipeTestCmd.Flags().String("profile", "default", "AWS profile to use")
	awsSteampipeTestCmd.Flags().String("region", "us-east-1", "AWS region")
}

func runAWSSteampipeTest(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")

	// Initialize Dagger engine
	fmt.Println("Initializing Dagger engine for AWS Steampipe test...")
	engine, err := dagger.NewEngine(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	client := engine.GetClient()

	// Create a temporary directory for Steampipe config
	tempDir, err := ioutil.TempDir("", "steampipe-config-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create Steampipe config directory structure
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	// Create AWS connection config
	awsConfig := fmt.Sprintf(`connection "aws" {
  plugin = "aws"
  profile = "%s"
  regions = ["%s"]
}
`, profile, region)

	configFile := filepath.Join(configDir, "aws.spc")
	if err := os.WriteFile(configFile, []byte(awsConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Created Steampipe config:\n%s\n", awsConfig)

	// Test 1: Container with all configurations
	fmt.Println("\n1. Setting up Steampipe container with AWS plugin...")

	container := client.Container().
		From("turbot/steampipe:latest")

	// Mount AWS credentials
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		awsCredsPath := filepath.Join(homeDir, ".aws")
		if _, err := os.Stat(awsCredsPath); err == nil {
			fmt.Println("✓ Mounting AWS credentials from ~/.aws")
			awsCreds := client.Host().Directory(awsCredsPath)
			container = container.WithDirectory("/home/steampipe/.aws", awsCreds)
		}
	}

	// Mount Steampipe config
	fmt.Println("✓ Mounting Steampipe configuration")
	steampipeConfig := client.Host().Directory(tempDir)
	container = container.WithDirectory("/home/steampipe/.steampipe", steampipeConfig)

	// Set environment variables
	container = container.
		WithEnvVariable("AWS_PROFILE", profile).
		WithEnvVariable("AWS_REGION", region).
		WithEnvVariable("AWS_SDK_LOAD_CONFIG", "1")

	// Install AWS plugin
	fmt.Println("\n2. Installing AWS plugin...")
	container = container.WithExec([]string{"steampipe", "plugin", "install", "aws"})

	output, err := container.Stdout(ctx)
	if err != nil {
		fmt.Printf("⚠️  Plugin installation had issues: %v\n", err)
	} else {
		fmt.Println("✓ AWS plugin installed")
	}

	// Test 2: Check plugin list
	fmt.Println("\n3. Checking installed plugins...")
	pluginContainer := container.WithExec([]string{"steampipe", "plugin", "list"})

	output, err = pluginContainer.Stdout(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to list plugins: %v\n", err)
	} else {
		fmt.Printf("✓ Installed plugins:\n%s\n", output)
	}

	// Test 3: Simple AWS query
	fmt.Println("\n4. Testing simple AWS query (list regions)...")
	queryContainer := container.WithExec([]string{
		"steampipe", "query",
		"SELECT name, opt_in_status FROM aws_region ORDER BY name",
		"--output", "json",
	})

	output, err = queryContainer.Stdout(ctx)
	if err != nil {
		fmt.Printf("❌ Query failed: %v\n", err)

		// Try to get error details
		errContainer := container.WithExec([]string{
			"steampipe", "query",
			"SELECT 1",
			"--output", "json",
		})
		if errOut, _ := errContainer.Stdout(ctx); errOut != "" {
			fmt.Printf("Basic query result: %s\n", errOut)
		}
	} else {
		fmt.Printf("✓ AWS regions query successful!\n")
		// Show first 500 chars of output
		if len(output) > 500 {
			fmt.Printf("Results (truncated):\n%s...\n", output[:500])
		} else {
			fmt.Printf("Results:\n%s\n", output)
		}
	}

	// Test 4: Check connection
	fmt.Println("\n5. Checking AWS connection status...")
	connContainer := container.WithExec([]string{
		"steampipe", "query",
		"SELECT * FROM steampipe_connection WHERE name = 'aws'",
		"--output", "json",
	})

	output, err = connContainer.Stdout(ctx)
	if err != nil {
		fmt.Printf("❌ Connection check failed: %v\n", err)
	} else {
		fmt.Printf("✓ Connection info:\n%s\n", output)
	}

	// Test 5: Try with explicit credentials if available
	if accessKey := os.Getenv("AWS_ACCESS_KEY_ID"); accessKey != "" {
		fmt.Println("\n6. Testing with environment credentials...")

		envContainer := client.Container().
			From("turbot/steampipe:latest").
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", region).
			WithExec([]string{"steampipe", "plugin", "install", "aws"})

		envContainer = envContainer.WithExec([]string{
			"steampipe", "query",
			"SELECT count(*) as region_count FROM aws_region",
			"--output", "json",
		})

		output, err = envContainer.Stdout(ctx)
		if err != nil {
			fmt.Printf("❌ Environment credential test failed: %v\n", err)
		} else {
			fmt.Printf("✓ Environment credentials work! Result: %s\n", output)
		}
	}

	green := color.New(color.FgGreen)
	green.Printf("\n✓ AWS Steampipe test completed!\n")

	fmt.Println("\nRecommendations:")
	fmt.Println("1. Ensure AWS credentials are properly configured in ~/.aws/credentials")
	fmt.Println("2. Verify the AWS profile has appropriate permissions")
	fmt.Println("3. Check that the AWS region is valid and accessible")

	return nil
}
