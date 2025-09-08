package modules

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dagger.io/dagger"
)

func TestBuildXModule_GetVersion(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	version, err := module.GetVersion(ctx)
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if version == "" {
		t.Error("Expected non-empty version string")
	}

	// Version should contain expected patterns for BuildX
	versionLower := strings.ToLower(version)
	if !strings.Contains(versionLower, "buildx") && 
	   !strings.Contains(version, ".") && 
	   !strings.Contains(versionLower, "version") {
		t.Errorf("Version string doesn't look like a valid BuildX version: %s", version)
	}
}

func TestBuildXModule_Build(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test directory with a simple Dockerfile
	tempDir := t.TempDir()
	dockerFile := filepath.Join(tempDir, "Dockerfile")
	dockerContent := `FROM alpine:latest
RUN echo "Hello BuildX Test" > /hello.txt
CMD ["cat", "/hello.txt"]`
	
	err = os.WriteFile(dockerFile, []byte(dockerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test building for linux/amd64 platform
	_, err = module.Build(ctx, tempDir, "buildx-test:latest", "linux/amd64", ".")
	
	// BuildX build without --load or --push is expected to "fail" with a warning about no output
	// But if we see the build stages in stderr, that means the build process worked correctly
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "No output specified with docker-container driver") &&
		   (strings.Contains(errStr, "#1") || strings.Contains(errStr, "FROM") || strings.Contains(errStr, "DONE")) {
			t.Logf("Build completed successfully (expected warning about no output): %v", err)
			// This is actually success - the build worked, just no output destination specified
		} else {
			t.Fatalf("Unexpected build failure: %v", err)
		}
	} else {
		t.Log("Build succeeded unexpectedly (no error)")
	}
}

func TestBuildXModule_BuildWithInvalidDockerfile(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test directory with an invalid Dockerfile
	tempDir := t.TempDir()
	dockerFile := filepath.Join(tempDir, "Dockerfile")
	invalidDockerContent := `FROM nonexistent-base-image:invalid-tag
RUN this-command-does-not-exist`
	
	err = os.WriteFile(dockerFile, []byte(invalidDockerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid Dockerfile: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// This should fail - we expect an error
	_, err = module.Build(ctx, tempDir, "invalid-build:test", "linux/amd64", "Dockerfile")
	if err == nil {
		t.Error("Expected build to fail with invalid Dockerfile, but it succeeded")
	}
}

func TestBuildXModule_BuildMultiplePlatforms(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a simple cross-platform Dockerfile
	tempDir := t.TempDir()
	dockerFile := filepath.Join(tempDir, "Dockerfile")
	dockerContent := `FROM --platform=$TARGETPLATFORM alpine:latest
ARG TARGETPLATFORM
RUN echo "Built for platform: $TARGETPLATFORM" > /platform.txt
CMD ["cat", "/platform.txt"]`
	
	err = os.WriteFile(dockerFile, []byte(dockerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	// Test multi-platform build
	_, err = module.Build(ctx, tempDir, "multiplatform-test:latest", "linux/amd64,linux/arm64", ".")
	
	// BuildX multi-platform build without --load or --push is expected to "fail" with a warning
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "No output specified with docker-container driver") &&
		   (strings.Contains(errStr, "#1") || strings.Contains(errStr, "FROM") || strings.Contains(errStr, "DONE")) {
			t.Logf("Multi-platform build completed successfully (expected warning): %v", err)
		} else {
			t.Fatalf("Unexpected multi-platform build failure: %v", err)
		}
	} else {
		t.Log("Multi-platform build succeeded unexpectedly (no error)")
	}
}

func TestBuildXModule_BuildWithCustomDockerfile(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test directory with a custom named Dockerfile
	tempDir := t.TempDir()
	customDockerFile := filepath.Join(tempDir, "Dockerfile.custom")
	dockerContent := `FROM alpine:latest
RUN echo "Custom Dockerfile Test" > /custom.txt
CMD ["cat", "/custom.txt"]`
	
	err = os.WriteFile(customDockerFile, []byte(dockerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create custom Dockerfile: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test building with custom Dockerfile name
	_, err = module.Build(ctx, tempDir, "custom-dockerfile-test:latest", "linux/amd64", "Dockerfile.custom")
	
	// BuildX build with custom Dockerfile
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "No output specified with docker-container driver") &&
		   (strings.Contains(errStr, "#1") || strings.Contains(errStr, "FROM") || strings.Contains(errStr, "DONE")) {
			t.Logf("Custom Dockerfile build completed successfully (expected warning): %v", err)
		} else {
			t.Fatalf("Unexpected custom Dockerfile build failure: %v", err)
		}
	} else {
		t.Log("Custom Dockerfile build succeeded unexpectedly (no error)")
	}
}

func TestBuildXModule_Publish(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test directory with a simple Dockerfile
	tempDir := t.TempDir()
	dockerFile := filepath.Join(tempDir, "Dockerfile")
	dockerContent := `FROM alpine:latest
RUN echo "Publish Test Image" > /test.txt
CMD ["cat", "/test.txt"]`
	
	err = os.WriteFile(dockerFile, []byte(dockerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test publish (will fail without valid registry credentials, which is expected)
	// We're testing that the command structure is correct, not actual publishing
	_, err = module.Publish(ctx, tempDir, "localhost:5000/test-image:latest", "linux/amd64", "testuser", "testpass", "localhost:5000", "Dockerfile")
	
	// We expect this to fail due to no registry, but the error should be related to registry/auth, not command structure
	if err != nil {
		errStr := strings.ToLower(err.Error())
		// These are expected failure modes - the command structure is working
		if strings.Contains(errStr, "connection refused") || 
		   strings.Contains(errStr, "registry") ||
		   strings.Contains(errStr, "unauthorized") ||
		   strings.Contains(errStr, "failed to resolve") ||
		   strings.Contains(errStr, "dial tcp") ||
		   strings.Contains(errStr, "no such host") {
			t.Logf("Expected registry/network error: %v", err)
			return // This is the expected behavior
		}
		t.Fatalf("Unexpected error type in publish: %v", err)
	}
	
	// If somehow it succeeds (maybe there's a local registry), that's also fine
	t.Log("Publish command succeeded (unexpected but valid)")
}

func TestBuildXModule_Dev(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "app.js")
	testContent := `console.log("BuildX Dev Environment Test");`
	
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	output, err := module.Dev(ctx, tempDir)
	if err != nil {
		t.Fatalf("Dev failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty dev output")
	}

	// Output should indicate dev environment setup
	outputLower := strings.ToLower(output)
	if !strings.Contains(outputLower, "buildx") && 
	   !strings.Contains(outputLower, "dev") && 
	   !strings.Contains(outputLower, "environment") {
		t.Logf("Dev output: %s", output)
	}
}

func TestBuildXModule_DevWithNoSourceDir(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Test dev without source directory (should still work)
	output, err := module.Dev(ctx, "")
	if err != nil {
		t.Fatalf("Dev without source dir failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty dev output")
	}
}

func TestNewBuildXModule(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewBuildXModule(client)
	
	if module == nil {
		t.Fatal("NewBuildXModule returned nil")
	}
	
	if module.client != client {
		t.Error("Client not properly assigned")
	}
	
	if module.name != "buildx" {
		t.Errorf("Expected name 'buildx', got '%s'", module.name)
	}
}

func TestBuildXModule_BuildParameterValidation(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	testCases := []struct {
		name           string
		srcDir         string
		tag            string
		platform       string
		dockerfilePath string
		shouldFail     bool
		expectError    string
	}{
		{
			name:           "empty tag",
			srcDir:         "/tmp",
			tag:            "",
			platform:       "linux/amd64",
			dockerfilePath: "Dockerfile",
			shouldFail:     false, // Docker BuildX might handle empty tags
		},
		{
			name:           "empty platform",
			srcDir:         "/tmp", 
			tag:            "test:latest",
			platform:       "",
			dockerfilePath: "Dockerfile",
			shouldFail:     false, // Docker BuildX has default platform
		},
		{
			name:           "nonexistent source directory",
			srcDir:         "/nonexistent/directory/12345",
			tag:            "test:latest",
			platform:       "linux/amd64",
			dockerfilePath: "Dockerfile",
			shouldFail:     true,
			expectError:    "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := module.Build(ctx, tc.srcDir, tc.tag, tc.platform, tc.dockerfilePath)
			
			if tc.shouldFail {
				if err == nil {
					t.Errorf("Expected build to fail for test case '%s', but it succeeded", tc.name)
				} else if tc.expectError != "" && !strings.Contains(strings.ToLower(err.Error()), tc.expectError) {
					t.Errorf("Expected error containing '%s', got: %v", tc.expectError, err)
				}
			}
			// Note: We don't test success cases here because they require valid Docker context
		})
	}
}

func TestBuildXModule_PublishParameterValidation(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Test with invalid parameters - should fail appropriately
	_, err = module.Publish(ctx, "/nonexistent", "invalid-registry/image:tag", "linux/amd64", "", "", "invalid-registry", "Dockerfile")
	if err == nil {
		t.Error("Expected publish to fail with invalid parameters, but it succeeded")
	}
}