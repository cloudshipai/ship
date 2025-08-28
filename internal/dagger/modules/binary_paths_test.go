package modules

import (
	"context"
	"testing"
	"time"

	"dagger.io/dagger"
)

// TestBinaryPaths tests that the binary constants are correctly defined and work in their containers
func TestBinaryPaths(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	tests := []struct {
		name        string
		image       string
		binary      string
		testCmd     []string
		expectError bool
	}{
		{
			name:    "Checkov",
			image:   "bridgecrew/checkov:latest",
			binary:  checkovBinary,
			testCmd: []string{checkovBinary, "--version"},
		},
		{
			name:    "Trivy",
			image:   "aquasec/trivy:latest", 
			binary:  "trivy",
			testCmd: []string{"trivy", "--version"},
		},
		{
			name:    "TFLint",
			image:   "ghcr.io/terraform-linters/tflint:latest",
			binary:  "tflint",
			testCmd: []string{"tflint", "--version"},
		},
		{
			name:    "Syft",
			image:   "anchore/syft:latest",
			binary:  "syft",
			testCmd: []string{"syft", "--version"},
		},
		{
			name:    "Grype",
			image:   "anchore/grype:latest",
			binary:  "grype",
			testCmd: []string{"grype", "--version"},
		},
		{
			name:    "Actionlint",
			image:   "wolff2023/actionlint:latest",
			binary:  actionlintBinary,
			testCmd: []string{actionlintBinary, "--version"},
		},
		{
			name:    "Semgrep",
			image:   "semgrep/semgrep:latest",
			binary:  "semgrep",
			testCmd: []string{"semgrep", "--version"},
		},
		{
			name:    "Hadolint",
			image:   "hadolint/hadolint:latest", 
			binary:  "hadolint",
			testCmd: []string{"hadolint", "--version"},
		},
		{
			name:    "Gitleaks",
			image:   "zricethezav/gitleaks:latest",
			binary:  "gitleaks",
			testCmd: []string{"gitleaks", "--version"},
		},
		{
			name:    "Trufflehog",
			image:   "trufflesecurity/trufflehog:latest",
			binary:  "trufflehog",
			testCmd: []string{"trufflehog", "--version"},
		},
		{
			name:    "Dockle",
			image:   "goodwithtech/dockle:latest",
			binary:  "dockle",
			testCmd: []string{"dockle", "--version"},
		},
		{
			name:    "Cosign",
			image:   "gcr.io/projectsigstore/cosign:latest",
			binary:  cosignBinary,
			testCmd: []string{cosignBinary, "version"},
		},
		{
			name:    "Nmap",
			image:   "instrumentisto/nmap:latest",
			binary:  "nmap",
			testCmd: []string{"nmap", "--version"},
		},
		{
			name:    "Nuclei",
			image:   "projectdiscovery/nuclei:latest",
			binary:  "nuclei",
			testCmd: []string{"nuclei", "--version"},
		},
		{
			name:    "OSSF Scorecard",
			image:   "gcr.io/openssf/scorecard:stable",
			binary:  scorecardBinary,
			testCmd: []string{scorecardBinary, "version"},
		},
		{
			name:    "Packer",
			image:   "hashicorp/packer:latest",
			binary:  packerBinary,
			testCmd: []string{packerBinary, "version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := client.Container().
				From(tt.image).
				WithExec(tt.testCmd, dagger.ContainerWithExecOpts{
					// Some tools may return non-zero exit codes for --version
					Expect: "ANY",
				})

			output, err := container.Stdout(ctx)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
				return
			}

			if err != nil {
				// Try stderr if stdout fails
				stderr, stderrErr := container.Stderr(ctx)
				if stderrErr != nil {
					t.Errorf("Failed to run %s binary '%s': %v", tt.name, tt.binary, err)
					return
				}
				if stderr != "" {
					t.Logf("%s binary works (output in stderr): %s", tt.name, stderr[:min(100, len(stderr))])
					return
				}
				t.Errorf("Failed to run %s binary '%s': %v", tt.name, tt.binary, err)
				return
			}

			if output == "" {
				t.Errorf("No output from %s binary '%s'", tt.name, tt.binary)
				return
			}

			t.Logf("%s binary works: %s", tt.name, output[:min(100, len(output))])
		})
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestModuleCreation tests that modules can be created without errors
func TestModuleCreation(t *testing.T) {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	// Test key modules can be created
	modules := map[string]func(*dagger.Client) interface{}{
		"Checkov":        func(c *dagger.Client) interface{} { return NewCheckovModule(c) },
		"Actionlint":     func(c *dagger.Client) interface{} { return NewActionlintModule(c) },
		"Cosign":         func(c *dagger.Client) interface{} { return NewCosignModule(c) },
		"OSSF Scorecard": func(c *dagger.Client) interface{} { return NewOSSFScorecardModule(c) },
		"Packer":         func(c *dagger.Client) interface{} { return NewPackerModule(c) },
		"InfraScan":      func(c *dagger.Client) interface{} { return NewInfraScanModule(c) },
	}

	for name, createModule := range modules {
		t.Run(name, func(t *testing.T) {
			module := createModule(client)
			if module == nil {
				t.Errorf("Failed to create %s module", name)
			}
		})
	}
}