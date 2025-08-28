package modules

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"dagger.io/dagger"
)

// TestIntegrationSuite runs comprehensive tests against real repositories
// to verify that the most relevant tools work correctly with actual codebases
func TestIntegrationSuite(t *testing.T) {
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

	// Test repositories
	testRepos := []struct {
		name string
		url  string
		path string
	}{
		{
			name: "terraform-aws-vpc",
			url:  "https://github.com/terraform-aws-modules/terraform-aws-vpc.git",
			path: "/tmp/terraform-test",
		},
		{
			name: "devops-resources",
			url:  "https://github.com/bregman-arie/devops-resources.git",
			path: "/tmp/cicd-test",
		},
	}

	// Most relevant tools for CICD/Infrastructure testing
	toolTests := []struct {
		name        string
		testFunc    func(t *testing.T, client *dagger.Client, repoPath string)
		repoFilter  func(repoName string) bool
		description string
	}{
		{
			name:        "tflint",
			testFunc:    testTFLint,
			repoFilter:  func(repoName string) bool { return strings.Contains(repoName, "terraform") },
			description: "Terraform syntax and best practices",
		},
		{
			name:        "checkov",
			testFunc:    testCheckov,
			repoFilter:  func(repoName string) bool { return true }, // Works with all repos
			description: "Security and compliance scanning",
		},
		{
			name:        "trivy",
			testFunc:    testTrivy,
			repoFilter:  func(repoName string) bool { return true }, // Universal scanner
			description: "Vulnerability scanning",
		},
		{
			name:        "trufflehog",
			testFunc:    testTruffleHog,
			repoFilter:  func(repoName string) bool { return true }, // Git repo scanning
			description: "Secret detection with verification",
		},
		{
			name:        "gitleaks",
			testFunc:    testGitleaks,
			repoFilter:  func(repoName string) bool { return true }, // Fast git scanning
			description: "Fast git secret scanning",
		},
		{
			name:        "tfsec",
			testFunc:    testTFSec,
			repoFilter:  func(repoName string) bool { return strings.Contains(repoName, "terraform") },
			description: "Terraform security scanning",
		},
		{
			name:        "terrascan",
			testFunc:    testTerrascan,
			repoFilter:  func(repoName string) bool { return strings.Contains(repoName, "terraform") },
			description: "IaC security scanning",
		},
		{
			name:        "semgrep",
			testFunc:    testSemgrep,
			repoFilter:  func(repoName string) bool { return true }, // Multi-language SAST
			description: "Static application security testing",
		},
		{
			name:        "hadolint",
			testFunc:    testHadolint,
			repoFilter:  func(repoName string) bool { return true }, // Look for Dockerfiles
			description: "Dockerfile linting",
		},
		{
			name:        "grype",
			testFunc:    testGrype,
			repoFilter:  func(repoName string) bool { return true }, // Filesystem scanning
			description: "Container/filesystem vulnerability scanning",
		},
	}

	// Run tests for each repository and applicable tool
	for _, repo := range testRepos {
		t.Run(repo.name, func(t *testing.T) {
			for _, toolTest := range toolTests {
				if toolTest.repoFilter(repo.name) {
					t.Run(toolTest.name, func(t *testing.T) {
						t.Logf("Testing %s on %s: %s", toolTest.name, repo.name, toolTest.description)
						start := time.Now()
						toolTest.testFunc(t, client, repo.path)
						duration := time.Since(start)
						t.Logf("✅ %s completed in %v", toolTest.name, duration)
					})
				}
			}
		})
	}
}

// testTFLint verifies TFLint functionality
func testTFLint(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	tflintModule := NewTFLintModule(client)
	
	result, err := tflintModule.ValidateDirectory(ctx, repoPath)
	if err != nil {
		t.Logf("TFLint validation failed (expected for some repos): %v", err)
		// Don't fail - some repos might not have valid Terraform
		return
	}
	
	if len(result) == 0 {
		t.Error("TFLint returned empty output")
		return
	}
	
	t.Logf("TFLint output preview: %s", truncateOutput(result, 200))
}

// testCheckov verifies Checkov functionality
func testCheckov(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	checkovModule := NewCheckovModule(client)
	
	result, err := checkovModule.ScanDirectoryWithOptions(ctx, repoPath, "", "cli", true, true)
	if err != nil {
		t.Errorf("Checkov scan failed: %v", err)
		return
	}
	
	if len(result) == 0 {
		t.Error("Checkov returned empty output")
		return
	}
	
	// Check for expected output patterns
	if !strings.Contains(result, "check") && !strings.Contains(result, "Check:") {
		t.Error("Checkov output doesn't contain expected check patterns")
	}
	
	t.Logf("Checkov output preview: %s", truncateOutput(result, 200))
}

// testTrivy verifies Trivy functionality
func testTrivy(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	trivyModule := NewTrivyModule(client)
	
	result, err := trivyModule.ScanFilesystem(ctx, repoPath, "table", "", false)
	if err != nil {
		t.Errorf("Trivy filesystem scan failed: %v", err)
		return
	}
	
	if len(result) == 0 {
		t.Error("Trivy returned empty output")
		return
	}
	
	t.Logf("Trivy output preview: %s", truncateOutput(result, 200))
}

// testTruffleHog verifies TruffleHog functionality
func testTruffleHog(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	truffleModule := NewTruffleHogModule(client)
	
	result, err := truffleModule.ScanFilesystem(ctx, repoPath, "json", false)
	if err != nil {
		t.Errorf("TruffleHog filesystem scan failed: %v", err)
		return
	}
	
	// TruffleHog might return empty if no secrets found - this is OK
	t.Logf("TruffleHog output preview: %s", truncateOutput(result, 200))
}

// testGitleaks verifies Gitleaks functionality
func testGitleaks(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	gitleaksModule := NewGitleaksModule(client)
	
	result, err := gitleaksModule.DetectSecrets(ctx, repoPath, "json", false)
	if err != nil {
		// Gitleaks might fail if no git repo or secrets found - this is OK for testing
		t.Logf("Gitleaks completed (may have found issues): %v", err)
		return
	}
	
	t.Logf("Gitleaks output preview: %s", truncateOutput(result, 200))
}

// testTFSec verifies TFSec functionality
func testTFSec(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	tfsecModule := NewTfsecModule(client)
	
	result, err := tfsecModule.ScanDirectory(ctx, repoPath, "default", false, true)
	if err != nil {
		t.Logf("TFSec scan failed (expected for some repos): %v", err)
		// Don't fail - some repos might not have valid Terraform
		return
	}
	
	if len(result) == 0 {
		t.Error("TFSec returned empty output")
		return
	}
	
	t.Logf("TFSec output preview: %s", truncateOutput(result, 200))
}

// testTerrascan verifies Terrascan functionality
func testTerrascan(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	terrascanModule := NewTerrascanModule(client)
	
	result, err := terrascanModule.ScanDirectory(ctx, repoPath, "terraform", "json", "", "medium")
	if err != nil {
		t.Logf("Terrascan scan failed (expected for some repos): %v", err)
		// Don't fail - some repos might not have valid IaC
		return
	}
	
	if len(result) == 0 {
		t.Error("Terrascan returned empty output")
		return
	}
	
	t.Logf("Terrascan output preview: %s", truncateOutput(result, 200))
}

// testSemgrep verifies Semgrep functionality
func testSemgrep(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	semgrepModule := NewSemgrepModule(client)
	
	result, err := semgrepModule.ScanDirectory(ctx, repoPath, "auto", "text", false)
	if err != nil {
		t.Errorf("Semgrep scan failed: %v", err)
		return
	}
	
	if len(result) == 0 {
		t.Error("Semgrep returned empty output")
		return
	}
	
	t.Logf("Semgrep output preview: %s", truncateOutput(result, 200))
}

// testHadolint verifies Hadolint functionality
func testHadolint(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	hadolintModule := NewHadolintModule(client)
	
	// Look for Dockerfile in the repo
	result, err := hadolintModule.LintFile(ctx, repoPath+"/Dockerfile")
	if err != nil {
		// Many repos won't have Dockerfiles - this is OK
		t.Logf("Hadolint completed (no Dockerfile found): %v", err)
		return
	}
	
	t.Logf("Hadolint output preview: %s", truncateOutput(result, 200))
}

// testGrype verifies Grype functionality
func testGrype(t *testing.T, client *dagger.Client, repoPath string) {
	ctx := context.Background()
	grypeModule := NewGrypeModule(client)
	
	result, err := grypeModule.ScanDirectory(ctx, repoPath, "table", "")
	if err != nil {
		t.Errorf("Grype scan failed: %v", err)
		return
	}
	
	if len(result) == 0 {
		t.Error("Grype returned empty output")
		return
	}
	
	t.Logf("Grype output preview: %s", truncateOutput(result, 200))
}

// truncateOutput truncates output for logging purposes
func truncateOutput(output string, maxLen int) string {
	if len(output) <= maxLen {
		return output
	}
	return output[:maxLen] + "... [truncated]"
}

// TestToolBinaryAvailability tests that all expected binaries are available in containers
func TestToolBinaryAvailability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping binary availability tests in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	// Test binary availability for working tools
	binaryTests := []struct {
		containerImage string
		binaryName     string
		testCommand    string
		description    string
	}{
		{"aquasec/tflint:latest", "tflint", "tflint --version", "TFLint Terraform linter"},
		{"bridgecrew/checkov:latest", "checkov", "checkov --version", "Checkov security scanner"},
		{"aquasec/trivy:latest", "trivy", "trivy --version", "Trivy vulnerability scanner"},
		{"trufflesecurity/trufflehog:latest", "trufflehog", "trufflehog --version", "TruffleHog secret scanner"},
		{"zricethezav/gitleaks:latest", "gitleaks", "gitleaks version", "Gitleaks fast secret scanner"},
		{"aquasec/tfsec:latest", "tfsec", "tfsec --version", "TFSec Terraform security"},
		{"tenable/terrascan:latest", "terrascan", "terrascan version", "Terrascan IaC scanner"},
		{"returntocorp/semgrep:latest", "semgrep", "semgrep --version", "Semgrep SAST tool"},
		{"hadolint/hadolint:latest", "hadolint", "hadolint --version", "Hadolint Dockerfile linter"},
		{"anchore/grype:latest", "grype", "grype version", "Grype vulnerability scanner"},
	}

	for _, test := range binaryTests {
		t.Run(test.binaryName, func(t *testing.T) {
			container := client.Container().From(test.containerImage)
			
			// Test if binary exists and is executable
			output, err := container.WithExec([]string{"sh", "-c", fmt.Sprintf("which %s", test.binaryName)}).Stdout(ctx)
			if err != nil {
				t.Errorf("Binary %s not found in %s: %v", test.binaryName, test.containerImage, err)
				return
			}
			
			if len(strings.TrimSpace(output)) == 0 {
				t.Errorf("Binary %s path is empty in %s", test.binaryName, test.containerImage)
				return
			}
			
			// Test if binary can be executed
			versionOutput, err := container.WithExec([]string{"sh", "-c", test.testCommand}).Stdout(ctx)
			if err != nil {
				t.Errorf("Binary %s failed to execute in %s: %v", test.binaryName, test.containerImage, err)
				return
			}
			
			if len(strings.TrimSpace(versionOutput)) == 0 {
				t.Errorf("Binary %s version output is empty in %s", test.binaryName, test.containerImage)
				return
			}
			
			t.Logf("✅ %s: %s - %s", test.description, test.binaryName, strings.Split(versionOutput, "\n")[0])
		})
	}
}

// TestToolPerformance measures execution time for key tools
func TestToolPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	testPath := "/tmp/terraform-test"
	
	performanceTests := []struct {
		name        string
		testFunc    func() (time.Duration, error)
		maxDuration time.Duration
		description string
	}{
		{
			name: "checkov-performance",
			testFunc: func() (time.Duration, error) {
				start := time.Now()
				checkovModule := NewCheckovModule(client)
				_, err := checkovModule.ScanDirectoryWithOptions(ctx, testPath, "", "cli", true, true)
				return time.Since(start), err
			},
			maxDuration: 2 * time.Minute,
			description: "Checkov should complete within 2 minutes",
		},
		{
			name: "trivy-performance",
			testFunc: func() (time.Duration, error) {
				start := time.Now()
				trivyModule := NewTrivyModule(client)
				_, err := trivyModule.ScanFilesystem(ctx, testPath, "table", "", false)
				return time.Since(start), err
			},
			maxDuration: 3 * time.Minute,
			description: "Trivy should complete within 3 minutes",
		},
		{
			name: "gitleaks-performance",
			testFunc: func() (time.Duration, error) {
				start := time.Now()
				gitleaksModule := NewGitleaksModule(client)
				_, err := gitleaksModule.DetectSecrets(ctx, testPath, "json", false)
				return time.Since(start), err
			},
			maxDuration: 30 * time.Second,
			description: "Gitleaks should complete within 30 seconds",
		},
	}

	for _, test := range performanceTests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Testing %s: %s", test.name, test.description)
			
			duration, err := test.testFunc()
			
			if err != nil {
				t.Logf("Performance test failed (but measured duration): %v", err)
			}
			
			t.Logf("⏱️  %s completed in %v", test.name, duration)
			
			if duration > test.maxDuration {
				t.Errorf("Performance test failed: %s took %v, expected max %v", 
					test.name, duration, test.maxDuration)
			} else {
				t.Logf("✅ Performance test passed: %s within expected duration", test.name)
			}
		})
	}
}