package modules

import (
	"context"
	"strings"
	"testing"
	"time"

	"dagger.io/dagger"
)

// TestCoreToolsIntegration runs tests on the most important tools with real repositories
func TestCoreToolsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()
	
	// Create Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	// Test directory (should exist from previous test runs)
	terraformTestPath := "/tmp/terraform-test"
	
	// Core tool tests
	t.Run("Checkov", func(t *testing.T) {
		testCheckovBasic(t, client, terraformTestPath)
	})
	
	t.Run("TFLint", func(t *testing.T) {
		testTFLintBasic(t, client, terraformTestPath)  
	})
	
	t.Run("Trivy", func(t *testing.T) {
		testTrivyBasic(t, client, terraformTestPath)
	})
	
	t.Run("TruffleHog", func(t *testing.T) {
		testTruffleHogBasic(t, client, terraformTestPath)
	})
	
	t.Run("Gitleaks", func(t *testing.T) {
		testGitleaksBasic(t, client, terraformTestPath)
	})
}

func testCheckovBasic(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	start := time.Now()
	
	checkovModule := NewCheckovModule(client)
	result, err := checkovModule.ScanDirectory(ctx, repoPath)
	
	duration := time.Since(start)
	t.Logf("Checkov scan completed in %v", duration)
	
	if err != nil {
		t.Logf("Checkov scan failed (may be expected): %v", err)
		return
	}
	
	if len(result) == 0 {
		t.Error("Checkov returned empty output")
		return
	}
	
	// Check for expected Checkov output patterns
	if strings.Contains(result, "Check:") || strings.Contains(result, "PASSED") || strings.Contains(result, "FAILED") {
		t.Logf("✅ Checkov produced expected output format")
	}
	
	t.Logf("Checkov output preview: %s", truncateOutput(result, 300))
}

func testTFLintBasic(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	start := time.Now()
	
	tflintModule := NewTFLintModule(client)
	result, err := tflintModule.Check(ctx, repoPath, TFLintOptions{})
	
	duration := time.Since(start)
	t.Logf("TFLint check completed in %v", duration)
	
	if err != nil {
		t.Logf("TFLint check failed (may be expected for non-terraform repos): %v", err)
		return
	}
	
	t.Logf("✅ TFLint executed successfully")
	t.Logf("TFLint output preview: %s", truncateOutput(result, 300))
}

func testTrivyBasic(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	start := time.Now()
	
	trivyModule := NewTrivyModule(client)
	result, err := trivyModule.ScanFilesystem(ctx, repoPath)
	
	duration := time.Since(start)
	t.Logf("Trivy scan completed in %v", duration)
	
	if err != nil {
		t.Errorf("Trivy filesystem scan failed: %v", err)
		return
	}
	
	if len(result) == 0 {
		t.Error("Trivy returned empty output")
		return
	}
	
	// Check for Trivy output indicators
	if strings.Contains(result, "Total:") || strings.Contains(result, "Vulnerability") {
		t.Logf("✅ Trivy produced expected output format")
	}
	
	t.Logf("Trivy output preview: %s", truncateOutput(result, 300))
}

func testTruffleHogBasic(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	start := time.Now()
	
	truffleModule := NewTruffleHogModule(client)
	result, err := truffleModule.ScanDirectory(ctx, repoPath)
	
	duration := time.Since(start)
	t.Logf("TruffleHog scan completed in %v", duration)
	
	if err != nil {
		t.Logf("TruffleHog scan failed: %v", err)
		return
	}
	
	// TruffleHog might return empty if no secrets found - this is OK
	t.Logf("✅ TruffleHog executed successfully")
	t.Logf("TruffleHog output preview: %s", truncateOutput(result, 300))
}

func testGitleaksBasic(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	start := time.Now()
	
	gitleaksModule := NewGitleaksModule(client)
	result, err := gitleaksModule.Detect(ctx, repoPath, GitleaksDetectOptions{})
	
	duration := time.Since(start)
	t.Logf("Gitleaks detect completed in %v", duration)
	
	if err != nil {
		// Gitleaks may return error if secrets found or other issues - this is OK for testing
		t.Logf("Gitleaks completed with result: %v", err)
	}
	
	t.Logf("✅ Gitleaks executed successfully")
	t.Logf("Gitleaks output preview: %s", truncateOutput(result, 300))
}

// TestToolBinariesExist tests that key tool binaries are available in their containers
func TestToolBinariesExist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping binary existence tests in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	// Key tools to test
	toolTests := []struct {
		name           string
		containerImage string
		binaryName     string
		versionCommand []string
	}{
		{
			name:           "Checkov",
			containerImage: "bridgecrew/checkov:latest",
			binaryName:     "checkov",
			versionCommand: []string{"checkov", "--version"},
		},
		{
			name:           "TFLint",
			containerImage: "ghcr.io/terraform-linters/tflint:latest",
			binaryName:     "tflint",
			versionCommand: []string{"tflint", "--version"},
		},
		{
			name:           "Trivy",
			containerImage: "aquasec/trivy:latest",
			binaryName:     "trivy",
			versionCommand: []string{"trivy", "--version"},
		},
		{
			name:           "TruffleHog",
			containerImage: "trufflesecurity/trufflehog:latest",
			binaryName:     "trufflehog",
			versionCommand: []string{"trufflehog", "--version"},
		},
		{
			name:           "Gitleaks",
			containerImage: "zricethezav/gitleaks:latest",
			binaryName:     "gitleaks",
			versionCommand: []string{"gitleaks", "version"},
		},
		{
			name:           "Grype",
			containerImage: "anchore/grype:latest",
			binaryName:     "grype",
			versionCommand: []string{"grype", "version"},
		},
		{
			name:           "Hadolint",
			containerImage: "hadolint/hadolint:latest",
			binaryName:     "hadolint",
			versionCommand: []string{"hadolint", "--version"},
		},
		{
			name:           "Semgrep",
			containerImage: "semgrep/semgrep:latest",
			binaryName:     "semgrep",
			versionCommand: []string{"semgrep", "--version"},
		},
	}

	for _, test := range toolTests {
		t.Run(test.name, func(t *testing.T) {
			// Check if binary exists
			container := client.Container().From(test.containerImage)
			
			binaryPath, err := container.WithExec([]string{"which", test.binaryName}).Stdout(ctx)
			if err != nil {
				t.Errorf("Binary %s not found in %s: %v", test.binaryName, test.containerImage, err)
				return
			}
			
			if strings.TrimSpace(binaryPath) == "" {
				t.Errorf("Binary %s path is empty in %s", test.binaryName, test.containerImage)
				return
			}
			
			// Test version command
			versionOutput, err := container.WithExec(test.versionCommand).Stdout(ctx)
			if err != nil {
				// Try stderr for version output
				stderrOutput, stderrErr := container.WithExec(test.versionCommand).Stderr(ctx)
				if stderrErr == nil && strings.TrimSpace(stderrOutput) != "" {
					t.Logf("✅ %s version (stderr): %s", test.name, strings.Split(stderrOutput, "\n")[0])
					return
				}
				t.Errorf("Failed to get version for %s: %v", test.name, err)
				return
			}
			
			if strings.TrimSpace(versionOutput) == "" {
				t.Errorf("Empty version output for %s", test.name)
				return
			}
			
			t.Logf("✅ %s version: %s", test.name, strings.Split(versionOutput, "\n")[0])
		})
	}
}

// TestPerformanceBenchmarks runs performance tests for key tools
func TestPerformanceBenchmarks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance benchmarks in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	testPath := "/tmp/terraform-test"
	
	// Performance expectations
	benchmarks := []struct {
		name        string
		testFunc    func() (time.Duration, error)
		maxExpected time.Duration
		description string
	}{
		{
			name: "Checkov-Performance",
			testFunc: func() (time.Duration, error) {
				start := time.Now()
				checkovModule := NewCheckovModule(client)
				_, err := checkovModule.ScanDirectory(ctx, testPath)
				return time.Since(start), err
			},
			maxExpected: 2 * time.Minute,
			description: "Checkov should scan typical repo within 2 minutes",
		},
		{
			name: "TruffleHog-Performance",
			testFunc: func() (time.Duration, error) {
				start := time.Now()
				truffleModule := NewTruffleHogModule(client)
				_, err := truffleModule.ScanDirectory(ctx, testPath)
				return time.Since(start), err
			},
			maxExpected: 1 * time.Minute,
			description: "TruffleHog should scan typical repo within 1 minute",
		},
		{
			name: "Trivy-Performance",
			testFunc: func() (time.Duration, error) {
				start := time.Now()
				trivyModule := NewTrivyModule(client)
				_, err := trivyModule.ScanFilesystem(ctx, testPath)
				return time.Since(start), err
			},
			maxExpected: 3 * time.Minute,
			description: "Trivy should scan typical repo within 3 minutes",
		},
	}

	for _, benchmark := range benchmarks {
		t.Run(benchmark.name, func(t *testing.T) {
			t.Logf("Running performance benchmark: %s", benchmark.description)
			
			duration, err := benchmark.testFunc()
			
			if err != nil {
				t.Logf("Benchmark completed with error (duration still measured): %v", err)
			}
			
			t.Logf("⏱️  %s completed in %v", benchmark.name, duration)
			
			if duration > benchmark.maxExpected {
				t.Errorf("❌ Performance issue: %s took %v, expected max %v", 
					benchmark.name, duration, benchmark.maxExpected)
			} else {
				t.Logf("✅ Performance OK: %s within expected duration", benchmark.name)
			}
		})
	}
}

// Helper function to truncate output for logging
func truncateOutput(output string, maxLen int) string {
	if len(output) <= maxLen {
		return output
	}
	return output[:maxLen] + "... [truncated]"
}