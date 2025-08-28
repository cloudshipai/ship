package modules

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"dagger.io/dagger"
)

// OpenCodeModule runs OpenCode AI coding agent in containers
type OpenCodeModule struct {
	client    *dagger.Client
	name      string
	projectRoot string
}

// NewOpenCodeModule creates a new OpenCode module
func NewOpenCodeModule(client *dagger.Client) *OpenCodeModule {
	return &OpenCodeModule{
		client:      client,
		name:        "opencode",
		projectRoot: findOpenCodeProjectRoot(),
	}
}

// findOpenCodeProjectRoot finds the project root directory by looking for go.mod
func findOpenCodeProjectRoot() string {
	// Start from current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	
	// Walk up the directory tree looking for go.mod
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}
	
	// If we can't find go.mod, try to find it relative to this source file
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		// This file is in internal/dagger/modules/opencode.go
		// So project root is ../../../
		projectRoot := filepath.Join(filepath.Dir(filename), "..", "..", "..")
		if abs, err := filepath.Abs(projectRoot); err == nil {
			return abs
		}
	}
	
	return "."
}

// addCommonEnvVars adds common environment variables for AI providers
func (m *OpenCodeModule) addCommonEnvVars(container *dagger.Container) *dagger.Container {
	// List of common AI provider environment variables
	envVars := []string{
		"OPENAI_API_KEY",
		"ANTHROPIC_API_KEY",
		"CLAUDE_API_KEY", 
		"GEMINI_API_KEY",
		"GROQ_API_KEY",
		"OPENROUTER_API_KEY",
	}
	
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			container = container.WithEnvVariable(envVar, value)
		}
	}
	
	return container
}

// Chat starts an interactive chat session with OpenCode (files persist by default)
func (m *OpenCodeModule) Chat(ctx context.Context, workDir string, message string) (string, error) {
	return m.ChatWithOptions(ctx, workDir, message, true)
}

// ChatWithOptions starts an interactive chat session with OpenCode with configurable file persistence
func (m *OpenCodeModule) ChatWithOptions(ctx context.Context, workDir string, message string, persistFiles bool) (string, error) {
	return m.ChatWithSession(ctx, workDir, message, persistFiles, "", false)
}

// ChatWithSession starts an interactive chat session with OpenCode with session support
func (m *OpenCodeModule) ChatWithSession(ctx context.Context, workDir string, message string, persistFiles bool, sessionID string, continueSession bool) (string, error) {
	return m.ChatWithSessionAndModel(ctx, workDir, message, persistFiles, sessionID, continueSession, "")
}

// ChatWithSessionAndModel starts an interactive chat session with OpenCode with session support and model selection
func (m *OpenCodeModule) ChatWithSessionAndModel(ctx context.Context, workDir string, message string, persistFiles bool, sessionID string, continueSession bool, model string) (string, error) {
	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace")
	
	// Mount OpenCode session storage from host to enable session persistence
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/home/" + os.Getenv("USER")
	}
	opencodeSessionDir := filepath.Join(homeDir, ".local", "share", "opencode")
	
	// Ensure session directory exists on host
	if err := os.MkdirAll(opencodeSessionDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create session directory: %v\n", err)
	} else {
		// Mount the opencode session directory to enable persistence
		container = container.WithDirectory("/root/.local/share/opencode", 
			m.client.Host().Directory(opencodeSessionDir))
	}
	
	// Add environment variables for AI providers
	container = m.addCommonEnvVars(container)
	
	// Build the opencode run command with session support
	args := []string{"opencode", "run"}
	
	// Add model flag if provided
	if model != "" {
		args = append(args, "--model", model)
	}
	
	// Add session flags if provided
	if sessionID != "" {
		args = append(args, "--session", sessionID)
	}
	if continueSession {
		args = append(args, "--continue")
	}
	
	// Add the message
	args = append(args, message)
	
	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	// Conditionally export files based on persistFiles flag
	if persistFiles {
		// List files in workspace to debug what was created
		lsOutput, _ := container.WithExec([]string{"ls", "-la", "/workspace"}).Stdout(ctx)
		fmt.Printf("Debug: Files in container workspace:\n%s\n", lsOutput)
		
		// Export any files that may have been created back to the host
		_, err := container.Directory("/workspace").Export(ctx, workDir)
		if err != nil {
			// Log but don't fail - file creation might not have happened
			fmt.Printf("Note: Could not export files from container: %v\n", err)
		} else {
			fmt.Printf("Successfully exported workspace directory to %s\n", workDir)
		}
	} else {
		fmt.Printf("Running in ephemeral mode - files will not be persisted to host\n")
	}

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run opencode chat: no output received")
}

// Generate generates code based on a prompt
func (m *OpenCodeModule) Generate(ctx context.Context, workDir string, prompt string, outputFile string) (string, error) {
	message := prompt
	if outputFile != "" {
		message = fmt.Sprintf("%s and save it to %s", prompt, outputFile)
	}

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace")
	
	// Add environment variables for AI providers
	container = m.addCommonEnvVars(container)
	
	container = container.WithExec([]string{"opencode", "run", message}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	// List files in workspace to debug what was created
	lsOutput, _ := container.WithExec([]string{"ls", "-la", "/workspace"}).Stdout(ctx)
	fmt.Printf("Debug: Files in container workspace:\n%s\n", lsOutput)
	
	// Export any files that may have been created back to the host
	_, err := container.Directory("/workspace").Export(ctx, workDir)
	if err != nil {
		// Log but don't fail - file creation might not have happened
		fmt.Printf("Note: Could not export files from container: %v\n", err)
	} else {
		fmt.Printf("Successfully exported workspace directory to %s\n", workDir)
	}

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to generate code: no output received")
}

// AnalyzeFile analyzes a specific file with OpenCode
func (m *OpenCodeModule) AnalyzeFile(ctx context.Context, filePath string, question string) (string, error) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	message := fmt.Sprintf("%s about the file %s", question, filename)

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"opencode",
			"run",
			message,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to analyze file: no output received")
}

// Review performs code review on changes
func (m *OpenCodeModule) Review(ctx context.Context, workDir string, target string) (string, error) {
	args := []string{"opencode", "review"}
	if target != "" {
		args = append(args, "--target", target)
	}

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to review code: no output received")
}

// Refactor performs code refactoring based on instructions
func (m *OpenCodeModule) Refactor(ctx context.Context, workDir string, instructions string, files []string) (string, error) {
	args := []string{"opencode", "refactor", "--instructions", instructions}
	for _, file := range files {
		args = append(args, "--file", file)
	}

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to refactor code: no output received")
}

// Test generates and runs tests for code
func (m *OpenCodeModule) Test(ctx context.Context, workDir string, testType string, coverage bool) (string, error) {
	args := []string{"opencode", "test"}
	if testType != "" {
		args = append(args, "--type", testType)
	}
	if coverage {
		args = append(args, "--coverage")
	}

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run tests: no output received")
}

// Document generates documentation for code
func (m *OpenCodeModule) Document(ctx context.Context, workDir string, format string, outputDir string) (string, error) {
	args := []string{"opencode", "document"}
	if format != "" {
		args = append(args, "--format", format)
	}
	if outputDir != "" {
		args = append(args, "--output-dir", outputDir)
	}

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to generate documentation: no output received")
}

// WithAuth configures OpenCode with authentication credentials
func (m *OpenCodeModule) WithAuth(ctx context.Context, workDir string, provider string, apiKey string) (string, error) {
	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithEnvVariable(fmt.Sprintf("%s_API_KEY", provider), apiKey).
		WithExec([]string{
			"opencode",
			"auth",
			"configure",
			"--provider", provider,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to configure auth: no output received")
}

// GetVersion returns the version of OpenCode
func (m *OpenCodeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithExec([]string{"opencode", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get opencode version: %w", err)
	}

	return output, nil
}

// Interactive starts an interactive OpenCode session
func (m *OpenCodeModule) Interactive(ctx context.Context, workDir string, model string) (string, error) {
	args := []string{"opencode", "interactive"}
	if model != "" {
		args = append(args, "--model", model)
	}

	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to start interactive session: no output received")
}

// BatchProcess processes multiple files with OpenCode
func (m *OpenCodeModule) BatchProcess(ctx context.Context, workDir string, pattern string, operation string) (string, error) {
	container := m.client.Container().
		Build(m.client.Host().Directory(m.projectRoot), dagger.ContainerBuildOpts{
			Dockerfile: "internal/dagger/dockerfiles/opencode.dockerfile",
		}).
		WithDirectory("/workspace", m.client.Host().Directory(workDir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"opencode",
			"batch",
			"--pattern", pattern,
			"--operation", operation,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to batch process: no output received")
}