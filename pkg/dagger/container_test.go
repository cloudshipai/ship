package dagger

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerBuilder_From(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("set base image", func(t *testing.T) {
		container := engine.Container().From("alpine:latest")
		
		assert.NotNil(t, container)
		assert.Equal(t, "alpine:latest", container.image)
		assert.Nil(t, container.err)
	})
	
	t.Run("empty image", func(t *testing.T) {
		container := engine.Container().From("")
		
		assert.NotNil(t, container)
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "image cannot be empty")
	})
	
	t.Run("chain after error", func(t *testing.T) {
		container := engine.Container().From("").From("alpine:latest")
		
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "image cannot be empty")
	})
}

func TestContainerBuilder_WithMountedDirectory(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("mount directory", func(t *testing.T) {
		container := engine.Container().WithMountedDirectory("/workspace", "/host/path")
		
		assert.NotNil(t, container)
		assert.Equal(t, "/host/path", container.mounts["/workspace"])
		assert.Nil(t, container.err)
	})
	
	t.Run("empty container path", func(t *testing.T) {
		container := engine.Container().WithMountedDirectory("", "/host/path")
		
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "container path cannot be empty")
	})
	
	t.Run("empty host path", func(t *testing.T) {
		container := engine.Container().WithMountedDirectory("/workspace", "")
		
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "host path cannot be empty")
	})
	
	t.Run("multiple mounts", func(t *testing.T) {
		container := engine.Container().
			WithMountedDirectory("/workspace", "/host/path1").
			WithMountedDirectory("/data", "/host/path2")
		
		assert.Nil(t, container.err)
		assert.Equal(t, "/host/path1", container.mounts["/workspace"])
		assert.Equal(t, "/host/path2", container.mounts["/data"])
	})
}

func TestContainerBuilder_WithEnvVariable(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("set environment variable", func(t *testing.T) {
		container := engine.Container().WithEnvVariable("TEST_VAR", "test_value")
		
		assert.NotNil(t, container)
		assert.Equal(t, "test_value", container.env["TEST_VAR"])
		assert.Nil(t, container.err)
	})
	
	t.Run("empty variable name", func(t *testing.T) {
		container := engine.Container().WithEnvVariable("", "value")
		
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "environment variable name cannot be empty")
	})
	
	t.Run("multiple environment variables", func(t *testing.T) {
		container := engine.Container().
			WithEnvVariable("VAR1", "value1").
			WithEnvVariable("VAR2", "value2")
		
		assert.Nil(t, container.err)
		assert.Equal(t, "value1", container.env["VAR1"])
		assert.Equal(t, "value2", container.env["VAR2"])
	})
}

func TestContainerBuilder_WithWorkdir(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("set working directory", func(t *testing.T) {
		container := engine.Container().WithWorkdir("/app")
		
		assert.NotNil(t, container)
		assert.Equal(t, "/app", container.workdir)
		assert.Nil(t, container.err)
	})
	
	t.Run("empty workdir", func(t *testing.T) {
		container := engine.Container().WithWorkdir("")
		
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "workdir cannot be empty")
	})
}

func TestContainerBuilder_WithExec(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("add command", func(t *testing.T) {
		container := engine.Container().WithExec([]string{"echo", "hello"})
		
		assert.NotNil(t, container)
		assert.Len(t, container.commands, 1)
		assert.Equal(t, []string{"echo", "hello"}, container.commands[0])
		assert.Nil(t, container.err)
	})
	
	t.Run("empty command", func(t *testing.T) {
		container := engine.Container().WithExec([]string{})
		
		assert.Error(t, container.err)
		assert.Contains(t, container.err.Error(), "exec args cannot be empty")
	})
	
	t.Run("multiple commands", func(t *testing.T) {
		container := engine.Container().
			WithExec([]string{"echo", "hello"}).
			WithExec([]string{"ls", "-la"})
		
		assert.Nil(t, container.err)
		assert.Len(t, container.commands, 2)
		assert.Equal(t, []string{"echo", "hello"}, container.commands[0])
		assert.Equal(t, []string{"ls", "-la"}, container.commands[1])
	})
	
	t.Run("command mutation safety", func(t *testing.T) {
		originalCmd := []string{"echo", "hello"}
		container := engine.Container().WithExec(originalCmd)
		
		// Modify original slice
		originalCmd[1] = "modified"
		
		// Container should have original values
		assert.Equal(t, "hello", container.commands[0][1])
	})
}

func TestContainerBuilder_Stdout(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("successful execution", func(t *testing.T) {
		container := engine.Container().
			From("alpine:latest").
			WithWorkdir("/app").
			WithEnvVariable("TEST_VAR", "test_value").
			WithMountedDirectory("/workspace", "/host/path").
			WithExec([]string{"echo", "hello"})
		
		output, err := container.Stdout(ctx)
		
		assert.NoError(t, err)
		assert.Contains(t, output, "Mock execution of image: alpine:latest")
		assert.Contains(t, output, "Working directory: /app")
		assert.Contains(t, output, "TEST_VAR=test_value")
		assert.Contains(t, output, "/host/path -> /workspace")
		assert.Contains(t, output, "echo hello")
	})
	
	t.Run("execution with builder error", func(t *testing.T) {
		container := engine.Container().From("")
		
		output, err := container.Stdout(ctx)
		
		assert.Error(t, err)
		assert.Empty(t, output)
		assert.Contains(t, err.Error(), "image cannot be empty")
	})
	
	t.Run("execution on closed engine", func(t *testing.T) {
		container := engine.Container().From("alpine:latest")
		engine.Close()
		
		output, err := container.Stdout(ctx)
		
		assert.Error(t, err)
		assert.Empty(t, output)
		assert.Contains(t, err.Error(), "engine is closed")
	})
	
	t.Run("execution without base image", func(t *testing.T) {
		newEngine, _ := NewEngine(ctx)
		defer newEngine.Close()
		
		container := newEngine.Container().WithExec([]string{"echo", "hello"})
		
		output, err := container.Stdout(ctx)
		
		assert.Error(t, err)
		assert.Empty(t, output)
		assert.Contains(t, err.Error(), "no base image specified")
	})
}

func TestContainerBuilder_Stderr(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("stderr execution", func(t *testing.T) {
		container := engine.Container().From("alpine:latest")
		
		stderr, err := container.Stderr(ctx)
		
		assert.NoError(t, err)
		assert.Empty(t, stderr) // Mock implementation returns empty stderr
	})
	
	t.Run("stderr with builder error", func(t *testing.T) {
		container := engine.Container().From("")
		
		stderr, err := container.Stderr(ctx)
		
		assert.Error(t, err)
		assert.Empty(t, stderr)
	})
}

func TestContainerBuilder_CombinedOutput(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("combined output", func(t *testing.T) {
		container := engine.Container().From("alpine:latest").WithExec([]string{"echo", "test"})
		
		output, err := container.CombinedOutput(ctx)
		
		assert.NoError(t, err)
		assert.Contains(t, output, "Mock execution of image: alpine:latest")
		assert.Contains(t, output, "echo test")
	})
}

func TestHostManager_Directory(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("create directory reference", func(t *testing.T) {
		host := engine.Host()
		dirRef := host.Directory("/test/path")
		
		assert.NotNil(t, dirRef)
		assert.Equal(t, "/test/path", dirRef.Path())
		assert.Equal(t, engine, dirRef.engine)
	})
}

func TestContainerBuilder_ChainedOperations(t *testing.T) {
	ctx := context.Background()
	engine, err := NewEngine(ctx)
	require.NoError(t, err)
	defer engine.Close()
	
	t.Run("complex chained container build", func(t *testing.T) {
		output, err := engine.Container().
			From("ubuntu:20.04").
			WithEnvVariable("DEBIAN_FRONTEND", "noninteractive").
			WithWorkdir("/workspace").
			WithMountedDirectory("/workspace", "/host/project").
			WithExec([]string{"apt", "update"}).
			WithExec([]string{"apt", "install", "-y", "curl"}).
			WithExec([]string{"curl", "--version"}).
			Stdout(ctx)
		
		assert.NoError(t, err)
		assert.Contains(t, output, "ubuntu:20.04")
		assert.Contains(t, output, "DEBIAN_FRONTEND=noninteractive")
		assert.Contains(t, output, "Working directory: /workspace")
		assert.Contains(t, output, "/host/project -> /workspace")
		
		// Check all commands are present
		lines := strings.Split(output, "\n")
		cmdSection := false
		cmdCount := 0
		for _, line := range lines {
			if strings.Contains(line, "Commands executed:") {
				cmdSection = true
				continue
			}
			if cmdSection && strings.Contains(line, ":") {
				cmdCount++
			}
		}
		assert.Equal(t, 3, cmdCount) // Three WithExec calls
	})
}