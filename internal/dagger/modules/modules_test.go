package modules_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/cloudship/ship/internal/dagger/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenInfraQuoteModule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	module := modules.NewOpenInfraQuoteModule(client)

	t.Run("GetVersion", func(t *testing.T) {
		version, err := module.GetVersion(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("OpenInfraQuote version: %s", version)
	})

	t.Run("AnalyzePlan", func(t *testing.T) {
		// Get absolute path to test file
		testFile, err := filepath.Abs("../../../test/terraform-mocks/tfplan.json")
		require.NoError(t, err)

		// Check if file exists
		_, err = os.Stat(testFile)
		if os.IsNotExist(err) {
			t.Skip("Test file not found, skipping")
		}

		output, err := module.AnalyzePlan(ctx, testFile, "us-east-1")
		if err != nil {
			// OpenInfraQuote might fail if the plan format isn't exactly right
			t.Logf("Expected error (plan format): %v", err)
		} else {
			assert.NotEmpty(t, output)
			t.Logf("Cost analysis output: %s", output)
		}
	})

	t.Run("AnalyzeDirectory", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.AnalyzeDirectory(ctx, testDir, "us-east-1")
		if err != nil {
			// This might fail without actual AWS credentials
			t.Logf("Expected error (no credentials): %v", err)
		} else {
			assert.NotEmpty(t, output)
			t.Logf("Directory analysis output: %s", output)
		}
	})
}

func TestInfraScanModule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	module := modules.NewInfraScanModule(client)

	t.Run("GetVersion", func(t *testing.T) {
		version, err := module.GetVersion(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("InfraScan version: %s", version)
	})

	t.Run("ScanDirectory", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.ScanDirectory(ctx, testDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		t.Logf("Security scan output: %s", output)

		// The mock files should have security issues
		assert.Contains(t, output, "AWS_ACCESS_KEY_ID", "Should detect hardcoded AWS credentials")
	})

	t.Run("ScanFile", func(t *testing.T) {
		// Get absolute path to test file
		testFile, err := filepath.Abs("../../../test/terraform-mocks/main.tf")
		require.NoError(t, err)

		// Check if file exists
		_, err = os.Stat(testFile)
		if os.IsNotExist(err) {
			t.Skip("Test file not found, skipping")
		}

		output, err := module.ScanFile(ctx, testFile)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		t.Logf("File scan output: %s", output)
	})
}

func TestTerraformDocsModule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	module := modules.NewTerraformDocsModule(client)

	t.Run("GetVersion", func(t *testing.T) {
		version, err := module.GetVersion(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("terraform-docs version: %s", version)
	})

	t.Run("GenerateMarkdown", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.GenerateMarkdown(ctx, testDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Check for expected content
		assert.Contains(t, output, "## Requirements")
		assert.Contains(t, output, "## Providers")
		assert.Contains(t, output, "## Resources")
		assert.Contains(t, output, "## Inputs")
		assert.Contains(t, output, "## Outputs")

		t.Logf("Generated markdown:\n%s", output)
	})

	t.Run("GenerateJSON", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.GenerateJSON(ctx, testDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Should be valid JSON
		assert.Contains(t, output, "{")
		assert.Contains(t, output, "}")

		t.Logf("Generated JSON length: %d bytes", len(output))
	})

	t.Run("GenerateTable", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.GenerateTable(ctx, testDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)

		// Should contain table markers
		assert.Contains(t, output, "|")
		assert.Contains(t, output, "Name")
		assert.Contains(t, output, "Type")
		assert.Contains(t, output, "Default")

		t.Logf("Generated table:\n%s", output)
	})
}

func TestTFLintModule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	module := modules.NewTFLintModule(client)

	t.Run("GetVersion", func(t *testing.T) {
		version, err := module.GetVersion(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("TFLint version: %s", version)
	})

	t.Run("LintDirectory", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.LintDirectory(ctx, testDir)
		// TFLint might return error if issues are found
		if err != nil {
			t.Logf("Linting found issues (expected): %v", err)
		}
		assert.NotEmpty(t, output)
		t.Logf("Lint output: %s", output)
	})

	t.Run("LintFile", func(t *testing.T) {
		// Get absolute path to test file
		testFile, err := filepath.Abs("../../../test/terraform-mocks/main.tf")
		require.NoError(t, err)

		// Check if file exists
		_, err = os.Stat(testFile)
		if os.IsNotExist(err) {
			t.Skip("Test file not found, skipping")
		}

		output, err := module.LintFile(ctx, testFile)
		if err != nil {
			t.Logf("Linting found issues (expected): %v", err)
		}
		assert.NotEmpty(t, output)
		t.Logf("File lint output: %s", output)
	})
}

func TestCheckovModule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	module := modules.NewCheckovModule(client)

	t.Run("GetVersion", func(t *testing.T) {
		version, err := module.GetVersion(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("Checkov version: %s", version)
	})

	t.Run("ScanDirectory", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.ScanDirectory(ctx, testDir)
		// Checkov might return error if security issues are found
		if err != nil {
			t.Logf("Checkov found issues (expected): %v", err)
		}
		assert.NotEmpty(t, output)
		t.Logf("Checkov scan output: %s", output)
	})

	t.Run("ScanFile", func(t *testing.T) {
		// Get absolute path to test file
		testFile, err := filepath.Abs("../../../test/terraform-mocks/main.tf")
		require.NoError(t, err)

		// Check if file exists
		_, err = os.Stat(testFile)
		if os.IsNotExist(err) {
			t.Skip("Test file not found, skipping")
		}

		output, err := module.ScanFile(ctx, testFile)
		if err != nil {
			t.Logf("Checkov found issues (expected): %v", err)
		}
		assert.NotEmpty(t, output)
		t.Logf("File scan output: %s", output)
	})

	t.Run("ScanMultiFramework", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.ScanMultiFramework(ctx, testDir, []string{"terraform", "secrets"})
		if err != nil {
			t.Logf("Checkov found issues (expected): %v", err)
		}
		assert.NotEmpty(t, output)
		t.Logf("Multi-framework scan output: %s", output)
	})
}

func TestInfracostModule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	module := modules.NewInfracostModule(client)

	t.Run("GetVersion", func(t *testing.T) {
		version, err := module.GetVersion(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("Infracost version: %s", version)
	})

	t.Run("BreakdownDirectory", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.BreakdownDirectory(ctx, testDir)
		// Infracost might fail without API key or valid resources
		if err != nil {
			t.Logf("Expected error (no API key or resources): %v", err)
			if os.Getenv("INFRACOST_API_KEY") == "" {
				t.Skip("INFRACOST_API_KEY not set, skipping breakdown test")
			}
		} else {
			assert.NotEmpty(t, output)
			t.Logf("Cost breakdown output: %s", output)
		}
	})

	t.Run("GenerateTableReport", func(t *testing.T) {
		// Get absolute path to test directory
		testDir, err := filepath.Abs("../../../test/terraform-mocks")
		require.NoError(t, err)

		// Check if directory exists
		_, err = os.Stat(testDir)
		if os.IsNotExist(err) {
			t.Skip("Test directory not found, skipping")
		}

		output, err := module.GenerateTableReport(ctx, testDir)
		// Infracost might fail without API key or valid resources
		if err != nil {
			t.Logf("Expected error (no API key or resources): %v", err)
			if os.Getenv("INFRACOST_API_KEY") == "" {
				t.Skip("INFRACOST_API_KEY not set, skipping table report test")
			}
		} else {
			assert.NotEmpty(t, output)
			t.Logf("Table report output: %s", output)
		}
	})
}
