package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"dagger.io/dagger"
)

// TestCheckovDaggerPerformance tests Checkov performance via Dagger
func TestCheckovDaggerPerformance(t *testing.T) {
	// Skip if not in performance test mode
	if os.Getenv("CHECKOV_PERFORMANCE_TEST") == "" {
		t.Skip("Set CHECKOV_PERFORMANCE_TEST=1 to run performance tests")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Test directory - use the agents-cicd repo
	testDir := "/home/epuerta/projects/hack/agents-cicd"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Skipf("Test directory %s does not exist", testDir)
	}

	t.Logf("Testing Checkov Dagger module performance on: %s", testDir)

	// Create Checkov module
	checkov := NewCheckovModule(client)

	// Test 1: Basic directory scan with timing
	t.Run("ScanDirectory", func(t *testing.T) {
		start := time.Now()
		
		result, err := checkov.ScanDirectory(ctx, testDir)
		
		elapsed := time.Since(start)
		t.Logf("Dagger scan took: %v", elapsed)
		
		if err != nil {
			t.Fatalf("Failed to scan directory: %v", err)
		}

		// Validate JSON output
		var scanResult interface{}
		if err := json.Unmarshal([]byte(result), &scanResult); err != nil {
			t.Errorf("Failed to parse JSON output: %v", err)
			t.Logf("Raw output (first 500 chars): %s", truncateString(result, 500))
		} else {
			t.Logf("Successfully parsed JSON output")
			// Print summary if it's an array of results
			if results, ok := scanResult.([]interface{}); ok {
				for i, r := range results {
					if resultMap, ok := r.(map[string]interface{}); ok {
						if summary, exists := resultMap["summary"]; exists {
							t.Logf("Result %d summary: %+v", i, summary)
						}
					}
				}
			}
		}

		// Log output size
		t.Logf("Output size: %d bytes", len(result))
	})

	// Test 2: Scan with different output formats
	t.Run("ScanWithOptions_CLI", func(t *testing.T) {
		start := time.Now()
		
		result, err := checkov.ScanDirectoryWithOptions(ctx, testDir, "terraform", "cli", true, false)
		
		elapsed := time.Since(start)
		t.Logf("CLI format scan took: %v", elapsed)
		
		if err != nil {
			t.Fatalf("Failed to scan with CLI format: %v", err)
		}

		t.Logf("CLI output size: %d bytes", len(result))
		t.Logf("CLI output preview:\n%s", truncateString(result, 1000))
	})

	// Test 3: Multi-framework scan
	t.Run("ScanMultiFramework", func(t *testing.T) {
		start := time.Now()
		
		frameworks := []string{"terraform", "dockerfile", "github_configuration"}
		result, err := checkov.ScanMultiFramework(ctx, testDir, frameworks)
		
		elapsed := time.Since(start)
		t.Logf("Multi-framework scan took: %v", elapsed)
		
		if err != nil {
			t.Fatalf("Failed to scan with multiple frameworks: %v", err)
		}

		t.Logf("Multi-framework output size: %d bytes", len(result))
	})

	// Test 4: Secrets scan
	t.Run("ScanSecrets", func(t *testing.T) {
		start := time.Now()
		
		result, err := checkov.ScanSecrets(ctx, testDir, "json")
		
		elapsed := time.Since(start)
		t.Logf("Secrets scan took: %v", elapsed)
		
		if err != nil {
			t.Fatalf("Failed to scan secrets: %v", err)
		}

		t.Logf("Secrets scan output size: %d bytes", len(result))
		
		// Try to parse and show secrets found
		var scanResult interface{}
		if err := json.Unmarshal([]byte(result), &scanResult); err == nil {
			if results, ok := scanResult.([]interface{}); ok {
				for _, r := range results {
					if resultMap, ok := r.(map[string]interface{}); ok {
						if summary, exists := resultMap["summary"]; exists {
							t.Logf("Secrets scan summary: %+v", summary)
						}
					}
				}
			}
		}
	})

	// Test 5: Version check (should be fast)
	t.Run("GetVersion", func(t *testing.T) {
		start := time.Now()
		
		version, err := checkov.GetVersion(ctx)
		
		elapsed := time.Since(start)
		t.Logf("Version check took: %v", elapsed)
		
		if err != nil {
			t.Fatalf("Failed to get version: %v", err)
		}

		t.Logf("Checkov version: %s", version)
	})
}

// TestCheckovDaggerVsDockerComparison runs both approaches for comparison
func TestCheckovDaggerVsDockerComparison(t *testing.T) {
	if os.Getenv("CHECKOV_COMPARISON_TEST") == "" {
		t.Skip("Set CHECKOV_COMPARISON_TEST=1 to run comparison tests")
	}

	ctx := context.Background()
	testDir := "/home/epuerta/projects/hack/agents-cicd"
	
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Skipf("Test directory %s does not exist", testDir)
	}

	// Test Dagger approach
	t.Run("Dagger", func(t *testing.T) {
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		if err != nil {
			t.Fatalf("Failed to connect to Dagger: %v", err)
		}
		defer client.Close()

		checkov := NewCheckovModule(client)
		
		start := time.Now()
		result, err := checkov.ScanDirectoryWithOptions(ctx, testDir, "terraform", "cli", true, false)
		elapsed := time.Since(start)
		
		if err != nil {
			t.Errorf("Dagger scan failed: %v", err)
		} else {
			t.Logf("✅ Dagger scan completed in: %v", elapsed)
			t.Logf("   Output size: %d bytes", len(result))
			
			// Count issues found (rough estimate)
			failedCount := countSubstring(result, "FAILED")
			passedCount := countSubstring(result, "PASSED")
			t.Logf("   Approximate results: %d failed, %d passed", failedCount, passedCount)
		}
	})

	// For comparison, you could also run docker directly here
	t.Run("DirectDocker_Reference", func(t *testing.T) {
		// This would be for reference only - showing the command that was run earlier
		t.Logf("Direct Docker command for comparison:")
		t.Logf("docker run --rm -v \"%s:/tf\" bridgecrew/checkov:latest -d /tf --output cli", testDir)
		t.Logf("(This test just logs the equivalent command - actual execution was done earlier)")
	})
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}

// Helper function to count substring occurrences
func countSubstring(s, substr string) int {
	count := 0
	start := 0
	for {
		pos := indexOf(s[start:], substr)
		if pos == -1 {
			break
		}
		count++
		start += pos + len(substr)
	}
	return count
}

// Simple indexOf implementation
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestCheckovDaggerDetailed provides detailed output for debugging
func TestCheckovDaggerDetailed(t *testing.T) {
	if os.Getenv("CHECKOV_DETAILED_TEST") == "" {
		t.Skip("Set CHECKOV_DETAILED_TEST=1 to run detailed tests")
	}

	ctx := context.Background()
	testDir := "/home/epuerta/projects/hack/agents-cicd"
	
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Skipf("Test directory %s does not exist", testDir)
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	checkov := NewCheckovModule(client)

	t.Run("DetailedScan", func(t *testing.T) {
		t.Log("Starting detailed Checkov scan via Dagger...")
		t.Logf("Target directory: %s", testDir)
		
		start := time.Now()
		result, err := checkov.ScanDirectoryWithOptions(ctx, testDir, "", "json", false, false)
		elapsed := time.Since(start)
		
		if err != nil {
			t.Fatalf("Detailed scan failed: %v", err)
		}

		t.Logf("Scan completed in: %v", elapsed)
		t.Logf("Output size: %d bytes", len(result))

		// Write full output to file for inspection
		outputFile := fmt.Sprintf("/tmp/checkov-dagger-output-%d.json", time.Now().Unix())
		if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
			t.Logf("Failed to write output to file: %v", err)
		} else {
			t.Logf("Full output written to: %s", outputFile)
		}

		// Parse and summarize results
		var scanResults []interface{}
		if err := json.Unmarshal([]byte(result), &scanResults); err != nil {
			t.Logf("Failed to parse JSON: %v", err)
			return
		}

		for i, result := range scanResults {
			if resultMap, ok := result.(map[string]interface{}); ok {
				checkType := resultMap["check_type"]
				if summary, exists := resultMap["summary"]; exists {
					t.Logf("Result %d (%v) summary: %+v", i, checkType, summary)
				}
				
				// Count actual issues
				if results, exists := resultMap["results"]; exists {
					if resultsMap, ok := results.(map[string]interface{}); ok {
						if failed, exists := resultsMap["failed_checks"]; exists {
							if failedSlice, ok := failed.([]interface{}); ok {
								t.Logf("  Failed checks: %d", len(failedSlice))
								
								// Show first few failures for context
								for j, failure := range failedSlice {
									if j >= 3 { // Only show first 3
										t.Logf("  ... and %d more", len(failedSlice)-3)
										break
									}
									if failureMap, ok := failure.(map[string]interface{}); ok {
										checkName := failureMap["check_name"]
										filePath := failureMap["file_path"] 
										t.Logf("    %d. %v in %v", j+1, checkName, filePath)
									}
								}
							}
						}
					}
				}
			}
		}
	})
}