package modules

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dagger.io/dagger"
)

// TestBuildXModule_Integration tests the complete BuildX workflow
func TestBuildXModule_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a realistic test project
	tempDir := t.TempDir()
	
	// Create a simple Go application
	appFile := filepath.Join(tempDir, "main.go")
	appContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello from BuildX integration test!")
}`
	
	err = os.WriteFile(appFile, []byte(appContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create app file: %v", err)
	}

	// Create a multi-stage Dockerfile
	dockerFile := filepath.Join(tempDir, "Dockerfile")
	dockerContent := `# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go build -o hello main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/hello .
CMD ["./hello"]`
	
	err = os.WriteFile(dockerFile, []byte(dockerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Test 1: Get version
	version, err := module.GetVersion(ctx)
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}
	t.Logf("BuildX version: %s", version)

	// Test 2: Build single platform
	t.Log("Testing single platform build...")
	buildOutput, err := module.Build(ctx, tempDir, "buildx-integration-test:latest", "linux/amd64", "Dockerfile")
	if err != nil {
		t.Fatalf("Single platform build failed: %v", err)
	}
	t.Logf("Build output: %s", buildOutput)

	// Test 3: Dev environment (quick test)
	t.Log("Testing dev environment...")
	devOutput, err := module.Dev(ctx, tempDir)
	if err != nil {
		t.Fatalf("Dev environment test failed: %v", err)
	}
	t.Logf("Dev output: %s", devOutput)
}

// TestBuildXModule_RealWorldScenarios tests realistic usage patterns
func TestBuildXModule_RealWorldScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real-world scenario test in short mode")
	}

	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	scenarios := []struct {
		name           string
		dockerfile     string
		tag            string
		platform       string
		shouldSucceed  bool
	}{
		{
			name: "Simple Node.js app",
			dockerfile: `FROM node:18-alpine
WORKDIR /app
RUN echo '{"name":"test","version":"1.0.0"}' > package.json
RUN echo 'console.log("Hello Node!");' > index.js
CMD ["node", "index.js"]`,
			tag:           "node-test:latest",
			platform:      "linux/amd64",
			shouldSucceed: true,
		},
		{
			name: "Python app with requirements",
			dockerfile: `FROM python:3.11-slim
WORKDIR /app
RUN echo 'print("Hello Python!")' > app.py
CMD ["python", "app.py"]`,
			tag:           "python-test:latest", 
			platform:      "linux/amd64",
			shouldSucceed: true,
		},
		{
			name: "Multi-platform Alpine",
			dockerfile: `FROM alpine:latest
RUN echo "Multi-platform test" > /test.txt
CMD ["cat", "/test.txt"]`,
			tag:           "alpine-multiplatform:latest",
			platform:      "linux/amd64,linux/arm64",
			shouldSucceed: true,
		},
	}

	module := NewBuildXModule(client)

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Create temporary directory for this scenario
			tempDir := t.TempDir()
			dockerFile := filepath.Join(tempDir, "Dockerfile")
			
			err := os.WriteFile(dockerFile, []byte(scenario.dockerfile), 0644)
			if err != nil {
				t.Fatalf("Failed to create Dockerfile for scenario %s: %v", scenario.name, err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
			defer cancel()

			output, err := module.Build(ctx, tempDir, scenario.tag, scenario.platform, "Dockerfile")
			
			if scenario.shouldSucceed {
				if err != nil {
					t.Errorf("Scenario %s should have succeeded but failed: %v", scenario.name, err)
				} else {
					t.Logf("Scenario %s succeeded with output: %s", scenario.name, output)
				}
			} else {
				if err == nil {
					t.Errorf("Scenario %s should have failed but succeeded", scenario.name)
				} else {
					t.Logf("Scenario %s failed as expected: %v", scenario.name, err)
				}
			}
		})
	}
}

// TestBuildXModule_ErrorHandling tests various error conditions
func TestBuildXModule_ErrorHandling(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewBuildXModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	errorTests := []struct {
		name     string
		testFunc func() error
		expectError bool
	}{
		{
			name: "Missing Dockerfile",
			testFunc: func() error {
				tempDir := t.TempDir()
				// Don't create a Dockerfile
				_, err := module.Build(ctx, tempDir, "test:latest", "linux/amd64", "Dockerfile")
				return err
			},
			expectError: true,
		},
		{
			name: "Invalid source directory",
			testFunc: func() error {
				_, err := module.Build(ctx, "/this/directory/does/not/exist", "test:latest", "linux/amd64", "Dockerfile")
				return err
			},
			expectError: true,
		},
		{
			name: "Invalid platform format",
			testFunc: func() error {
				tempDir := t.TempDir()
				dockerfile := filepath.Join(tempDir, "Dockerfile")
				os.WriteFile(dockerfile, []byte("FROM alpine:latest"), 0644)
				_, err := module.Build(ctx, tempDir, "test:latest", "invalid/platform/format/too/many/parts", "Dockerfile")
				return err
			},
			expectError: true,
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			err := test.testFunc()
			
			if test.expectError && err == nil {
				t.Errorf("Test %s expected an error but got none", test.name)
			} else if !test.expectError && err != nil {
				t.Errorf("Test %s expected no error but got: %v", test.name, err)
			} else if err != nil {
				t.Logf("Test %s got expected error: %v", test.name, err)
			}
		})
	}
}