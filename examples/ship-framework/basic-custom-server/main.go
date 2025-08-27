package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudshipai/ship/pkg/dagger"
	"github.com/cloudshipai/ship/pkg/ship"
)

// EchoTool is a simple example tool
type EchoTool struct{}

func (t *EchoTool) Name() string {
	return "echo"
}

func (t *EchoTool) Description() string {
	return "Echoes the input message back to you"
}

func (t *EchoTool) Parameters() []ship.Parameter {
	return []ship.Parameter{
		{
			Name:        "message",
			Type:        "string",
			Description: "The message to echo back",
			Required:    true,
		},
		{
			Name:        "uppercase",
			Type:        "boolean",
			Description: "Whether to return the message in uppercase",
			Required:    false,
		},
	}
}

func (t *EchoTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	message, ok := params["message"].(string)
	if !ok {
		return &ship.ToolResult{
			Error: fmt.Errorf("message parameter is required"),
		}, fmt.Errorf("message parameter is required")
	}

	uppercase, _ := params["uppercase"].(bool)

	result := message
	if uppercase {
		result = strings.ToUpper(message)
	}

	return &ship.ToolResult{
		Content: fmt.Sprintf("Echo: %s", result),
		Metadata: map[string]interface{}{
			"original_message": message,
			"was_uppercase":    uppercase,
		},
	}, nil
}

// FileListTool demonstrates a containerized tool using Alpine Linux
type FileListTool struct{}

func (t *FileListTool) Name() string {
	return "file-list"
}

func (t *FileListTool) Description() string {
	return "Lists files in the current directory using a container"
}

func (t *FileListTool) Parameters() []ship.Parameter {
	return []ship.Parameter{
		{
			Name:        "path",
			Type:        "string",
			Description: "Path to list files from (default: current directory)",
			Required:    false,
		},
		{
			Name:        "show_hidden",
			Type:        "boolean",
			Description: "Show hidden files",
			Required:    false,
		},
	}
}

func (t *FileListTool) Execute(ctx context.Context, params map[string]interface{}, engine *dagger.Engine) (*ship.ToolResult, error) {
	path := "."
	if p, ok := params["path"].(string); ok && p != "" {
		path = p
	}

	args := []string{"ls", "-la"}
	showHidden, _ := params["show_hidden"].(bool)
	if !showHidden {
		args = []string{"ls", "-l"}
	}
	args = append(args, path)

	if engine == nil {
		return &ship.ToolResult{
			Error: fmt.Errorf("dagger engine not available"),
		}, fmt.Errorf("dagger engine not available")
	}

	// Mount current directory and run ls using Ship's Dagger wrapper
	container := engine.Container().
		From("alpine:latest").
		WithMountedDirectory("/workspace", ".").
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return &ship.ToolResult{
			Error: fmt.Errorf("failed to execute container command: %w", err),
		}, fmt.Errorf("failed to execute container command: %w", err)
	}

	return &ship.ToolResult{
		Content: fmt.Sprintf("File listing for %s:\n%s", path, output),
		Metadata: map[string]interface{}{
			"path":        path,
			"show_hidden": showHidden,
			"command":     strings.Join(args, " "),
		},
	}, nil
}

func main() {
	// Build MCP server with Ship framework
	server := ship.NewServer("basic-custom-server", "1.0.0").
		AddTool(&EchoTool{}).
		AddTool(&FileListTool{}).
		Build()

	fmt.Fprintf(log.Writer(), "Starting basic custom MCP server...\n")
	fmt.Fprintf(log.Writer(), "Available tools: echo, file-list\n")

	// Start the MCP server
	if err := server.ServeStdio(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}