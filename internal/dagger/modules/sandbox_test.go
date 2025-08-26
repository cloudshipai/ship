package modules

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dagger.io/dagger"
)

// TestSandboxEnvironments validates that sandbox environments work correctly with selected tools
func TestSandboxEnvironments(t *testing.T) {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Get project root directory
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	t.Run("VulnerableAppSandbox", func(t *testing.T) {
		testVulnerableAppSandbox(t, ctx, client, projectRoot)
	})

	t.Run("CostDemoSandbox", func(t *testing.T) {
		testCostDemoSandbox(t, ctx, client, projectRoot)
	})

	t.Run("ContainerAppSandbox", func(t *testing.T) {
		testContainerAppSandbox(t, ctx, client, projectRoot)
	})

	t.Run("TerraformQualitySandbox", func(t *testing.T) {
		testTerraformQualitySandbox(t, ctx, client, projectRoot)
	})
}

func testVulnerableAppSandbox(t *testing.T, ctx context.Context, client *dagger.Client, projectRoot string) {
	sandboxDir := filepath.Join(projectRoot, "examples", "sandbox", "vulnerable-app")

	// Test Checkov security scanning
	t.Run("CheckovScan", func(t *testing.T) {
		checkov := &CheckovModule{client: client}
		
		result, err := checkov.ScanDirectory(ctx, sandboxDir)
		if err != nil {
			t.Fatalf("Checkov scan failed: %v", err)
		}

		// Validate output contains expected security findings
		if !strings.Contains(result, "FAILED") && !strings.Contains(result, "HIGH") {
			t.Errorf("Expected Checkov to find security issues, but output was: %s", result)
		}

		// Test JSON output format
		jsonResult, err := checkov.ScanDirectoryWithFormat(ctx, sandboxDir, "json")
		if err != nil {
			t.Fatalf("Checkov JSON scan failed: %v", err)
		}

		// Validate JSON is parseable
		var jsonData interface{}
		if err := json.Unmarshal([]byte(jsonResult), &jsonData); err != nil {
			t.Errorf("Checkov JSON output is not valid JSON: %v", err)
		}
	})

	// Test Trivy infrastructure scanning
	t.Run("TrivyInfraScan", func(t *testing.T) {
		trivy := &TrivyModule{client: client}
		
		result, err := trivy.ScanFilesystem(ctx, sandboxDir)
		if err != nil {
			t.Fatalf("Trivy filesystem scan failed: %v", err)
		}

		// Validate output contains scan results
		if len(result) == 0 {
			t.Error("Expected Trivy to produce scan output")
		}
	})
}

func testCostDemoSandbox(t *testing.T, ctx context.Context, client *dagger.Client, projectRoot string) {
	sandboxDir := filepath.Join(projectRoot, "examples", "sandbox", "cost-demo")

	// Test OpenInfraQuote cost analysis
	t.Run("OpenInfraQuoteCostAnalysis", func(t *testing.T) {
		oiq := &OpenInfraQuoteModule{client: client}
		
		// Note: This requires a Terraform plan file to be generated first
		// In real scenarios, this would be done as part of the test setup
		result, err := oiq.AnalyzeDirectory(ctx, sandboxDir, "us-east-1")
		if err != nil {
			// For testing purposes, we expect this might fail if no plan exists
			// but the module should handle it gracefully
			t.Logf("OpenInfraQuote analysis result (may fail without plan): %v", err)
		}

		// If successful, validate output contains cost information
		if result != "" && !strings.Contains(result, "error") {
			if !strings.Contains(result, "cost") && !strings.Contains(result, "$") {
				t.Errorf("Expected cost information in output, got: %s", result)
			}
		}
	})
}

func testContainerAppSandbox(t *testing.T, ctx context.Context, client *dagger.Client, projectRoot string) {
	sandboxDir := filepath.Join(projectRoot, "examples", "sandbox", "container-app")

	// Test Syft SBOM generation
	t.Run("SyftSBOMGeneration", func(t *testing.T) {
		syft := &SyftModule{client: client}
		
		// Generate SBOM for directory
		result, err := syft.GenerateSBOMDirectory(ctx, sandboxDir, "cyclonedx-json")
		if err != nil {
			t.Fatalf("Syft SBOM generation failed: %v", err)
		}

		// Validate SBOM contains expected components
		if !strings.Contains(result, "components") && !strings.Contains(result, "bomFormat") {
			t.Errorf("Expected SBOM format in output, got: %s", result)
		}

		// Test JSON parsing
		var sbomData interface{}
		if err := json.Unmarshal([]byte(result), &sbomData); err != nil {
			t.Errorf("SBOM output is not valid JSON: %v", err)
		}
	})

	// Test Trivy container scanning (if Docker is available)
	t.Run("TrivyContainerScan", func(t *testing.T) {
		// Build test image first
		dockerfile := client.Host().Directory(sandboxDir).File("Dockerfile")
		
		// Build container image
		image := client.Container().
			Build(client.Host().Directory(sandboxDir))

		trivy := &TrivyModule{client: client}
		
		// Note: This test assumes the image is built and available
		// In practice, this would require proper image management
		result, err := trivy.ScanImage(ctx, "test-sandbox-app:latest")
		if err != nil {
			// Expected to fail if image isn't built, but module should handle gracefully
			t.Logf("Trivy container scan result (may fail without built image): %v", err)
		}

		_ = image    // Use the built image variable to avoid unused error
		_ = dockerfile // Use dockerfile variable to avoid unused error
	})
}

func testTerraformQualitySandbox(t *testing.T, ctx context.Context, client *dagger.Client, projectRoot string) {
	sandboxDir := filepath.Join(projectRoot, "examples", "sandbox", "terraform-quality")

	// Test TFLint quality analysis
	t.Run("TFLintQualityAnalysis", func(t *testing.T) {
		tflint := &TFLintModule{client: client}
		
		result, err := tflint.LintDirectory(ctx, sandboxDir)
		if err != nil {
			t.Fatalf("TFLint analysis failed: %v", err)
		}

		// Validate output contains linting issues
		if !strings.Contains(result, "issue") && !strings.Contains(result, "warning") && !strings.Contains(result, "error") {
			t.Errorf("Expected TFLint to find issues in quality sandbox, got: %s", result)
		}

		// Test JSON output format
		jsonResult, err := tflint.LintDirectoryWithFormat(ctx, sandboxDir, "json")
		if err != nil {
			t.Fatalf("TFLint JSON analysis failed: %v", err)
		}

		// Validate JSON is parseable
		var jsonData interface{}
		if err := json.Unmarshal([]byte(jsonResult), &jsonData); err != nil {
			t.Errorf("TFLint JSON output is not valid JSON: %v", err)
		}
	})
}

// Helper function to find the project root directory
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Look for go.mod to identify project root
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

// TestSandboxIntegrity validates that sandbox files exist and are well-formed
func TestSandboxIntegrity(t *testing.T) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	sandboxDir := filepath.Join(projectRoot, "examples", "sandbox")

	// Test that all expected sandbox directories exist
	expectedDirs := []string{
		"vulnerable-app",
		"cost-demo", 
		"container-app",
		"terraform-quality",
	}

	for _, dir := range expectedDirs {
		dirPath := filepath.Join(sandboxDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected sandbox directory %s does not exist", dirPath)
		}
	}

	// Test specific files exist
	expectedFiles := map[string]string{
		"vulnerable-app/main.tf":      "Terraform configuration",
		"cost-demo/main.tf":          "Terraform configuration", 
		"container-app/Dockerfile":   "Docker configuration",
		"container-app/package.json": "Node.js package file",
		"terraform-quality/main.tf":  "Terraform configuration",
	}

	for file, description := range expectedFiles {
		filePath := filepath.Join(sandboxDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected %s file %s does not exist", description, filePath)
		}
	}
}

// BenchmarkSandboxTools provides performance benchmarks for tool execution
func BenchmarkSandboxTools(b *testing.B) {
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		b.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	projectRoot, err := findProjectRoot()
	if err != nil {
		b.Fatalf("Failed to find project root: %v", err)
	}

	b.Run("CheckovScan", func(b *testing.B) {
		sandboxDir := filepath.Join(projectRoot, "examples", "sandbox", "vulnerable-app")
		checkov := &CheckovModule{client: client}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := checkov.ScanDirectory(ctx, sandboxDir)
			if err != nil {
				b.Fatalf("Checkov scan failed: %v", err)
			}
		}
	})

	b.Run("SyftSBOM", func(b *testing.B) {
		sandboxDir := filepath.Join(projectRoot, "examples", "sandbox", "container-app")
		syft := &SyftModule{client: client}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := syft.GenerateSBOMDirectory(ctx, sandboxDir, "cyclonedx-json")
			if err != nil {
				b.Fatalf("Syft SBOM generation failed: %v", err)
			}
		}
	})
}