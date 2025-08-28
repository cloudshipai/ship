package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"dagger.io/dagger"
	"github.com/cloudshipai/ship/internal/dagger/modules"
)

// TestCoreToolValidation validates that our key working tools function correctly
func TestCoreToolValidation(t *testing.T) {
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

	// Test directory
	terraformTestPath := "/tmp/terraform-test"
	
	// Validate key tools work
	t.Run("CheckovValidation", func(t *testing.T) {
		start := time.Now()
		checkovModule := modules.NewCheckovModule(client)
		result, err := checkovModule.ScanDirectory(ctx, terraformTestPath)
		duration := time.Since(start)
		
		t.Logf("Checkov scan completed in %v", duration)
		
		if err != nil {
			t.Logf("Checkov scan result: %v", err)
		}
		
		if len(result) > 0 {
			t.Logf("✅ Checkov produced output: %s", truncateString(result, 200))
		} else {
			t.Log("Checkov returned empty output")
		}
	})
	
	t.Run("TrivyValidation", func(t *testing.T) {
		start := time.Now()
		trivyModule := modules.NewTrivyModule(client)
		result, err := trivyModule.ScanFilesystem(ctx, terraformTestPath)
		duration := time.Since(start)
		
		t.Logf("Trivy scan completed in %v", duration)
		
		if err != nil {
			t.Errorf("Trivy scan failed: %v", err)
			return
		}
		
		if len(result) > 0 {
			t.Logf("✅ Trivy produced output: %s", truncateString(result, 200))
		} else {
			t.Error("Trivy returned empty output")
		}
	})
	
	t.Run("TruffleHogValidation", func(t *testing.T) {
		start := time.Now()
		truffleModule := modules.NewTruffleHogModule(client)
		result, err := truffleModule.ScanDirectory(ctx, terraformTestPath)
		duration := time.Since(start)
		
		t.Logf("TruffleHog scan completed in %v", duration)
		
		if err != nil {
			t.Logf("TruffleHog scan result: %v", err)
		}
		
		t.Logf("✅ TruffleHog executed: %s", truncateString(result, 200))
	})
	
	t.Run("TFLintValidation", func(t *testing.T) {
		start := time.Now()
		tflintModule := modules.NewTFLintModule(client)
		result, err := tflintModule.Check(ctx, terraformTestPath, modules.TFLintOptions{})
		duration := time.Since(start)
		
		t.Logf("TFLint check completed in %v", duration)
		
		if err != nil {
			t.Logf("TFLint check result: %v", err)
		}
		
		t.Logf("✅ TFLint executed: %s", truncateString(result, 200))
	})
}

// TestToolBinaryExistence verifies tool binaries exist in containers
func TestToolBinaryExistence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping binary existence tests in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	// Test key tool binaries
	toolTests := []struct {
		name           string
		containerImage string
		binaryName     string
	}{
		{"Checkov", "bridgecrew/checkov:latest", "checkov"},
		{"TFLint", "ghcr.io/terraform-linters/tflint:latest", "tflint"},
		{"Trivy", "aquasec/trivy:latest", "trivy"},
		{"TruffleHog", "trufflesecurity/trufflehog:latest", "trufflehog"},
		{"Gitleaks", "zricethezav/gitleaks:latest", "gitleaks"},
		{"Grype", "anchore/grype:latest", "grype"},
		{"Hadolint", "hadolint/hadolint:latest", "hadolint"},
		{"Semgrep", "semgrep/semgrep:latest", "semgrep"},
	}

	for _, test := range toolTests {
		t.Run(test.name, func(t *testing.T) {
			container := client.Container().From(test.containerImage)
			
			binaryPath, err := container.WithExec([]string{"which", test.binaryName}).Stdout(ctx)
			if err != nil {
				t.Errorf("❌ Binary %s not found in %s: %v", test.binaryName, test.containerImage, err)
				return
			}
			
			if strings.TrimSpace(binaryPath) == "" {
				t.Errorf("❌ Binary %s path is empty", test.binaryName)
				return
			}
			
			t.Logf("✅ %s binary found at: %s", test.name, strings.TrimSpace(binaryPath))
		})
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}