package mcp

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// logOpenCodeInteraction logs OpenCode MCP interactions to /tmp/opencode_mcp.log
func logOpenCodeInteraction(phase string, toolName string, request *mcp.CallToolRequest, result *mcp.CallToolResult, err error, duration time.Duration) {
	logFile := "/tmp/opencode_mcp.log"
	
	// Create log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	
	var logEntry map[string]interface{}
	if phase == "INPUT" {
		// Log input
		requestData, _ := json.MarshalIndent(request, "", "  ")
		logEntry = map[string]interface{}{
			"timestamp": timestamp,
			"phase":     phase,
			"tool":      toolName,
			"request":   json.RawMessage(requestData),
		}
	} else {
		// Log output
		var resultData []byte
		if result != nil {
			resultData, _ = json.MarshalIndent(result, "", "  ")
		}
		
		logEntry = map[string]interface{}{
			"timestamp": timestamp,
			"phase":     phase,
			"tool":      toolName,
			"result":    json.RawMessage(resultData),
			"error":     nil,
			"duration":  duration.String(),
		}
		
		if err != nil {
			logEntry["error"] = err.Error()
		}
	}
	
	// Write to log file
	logData, _ := json.MarshalIndent(logEntry, "", "  ")
	
	file, fileErr := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if fileErr != nil {
		log.Printf("Failed to open OpenCode log file: %v", fileErr)
		return
	}
	defer file.Close()
	
	file.WriteString(string(logData) + "\n---\n")
}

// AddOpenCodeTools adds OpenCode AI coding assistant tools to the MCP server
func AddOpenCodeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	// OpenCode Run - The main tool that handles all OpenCode functionality
	opencodeRunTool := mcp.NewTool("opencode_run",
		mcp.WithDescription("Send messages to OpenCode AI coding assistant. This tool can handle all OpenCode operations including code analysis, generation, review, refactoring, and general coding questions."),
		mcp.WithString("message", mcp.Required(), mcp.Description("Message, question, or instruction to send to OpenCode AI (e.g., 'analyze this code', 'generate a function that...', 'review my code')")),
		mcp.WithString("workdir", mcp.Description("Working directory for OpenCode operations (default: current directory)")),
		mcp.WithBoolean("persist_files", mcp.Description("Whether to persist files created by OpenCode back to the host (default: true)")),
		mcp.WithString("model", mcp.Description("AI model to use (format: provider/model, e.g., 'openai/gpt-4', 'anthropic/claude-3-sonnet')")),
		mcp.WithString("session", mcp.Description("Session ID for continuing conversations")),
		mcp.WithBoolean("continue", mcp.Description("Continue the last session")),
		mcp.WithBoolean("share", mcp.Description("Share the session")),
		mcp.WithString("agent", mcp.Description("Specific agent to use")),
	)

	s.AddTool(opencodeRunTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		startTime := time.Now()
		
		// Log input
		logOpenCodeInteraction("INPUT", "opencode_run", &request, nil, nil, 0)
		message := request.GetString("message", "")
		workdir := request.GetString("workdir", ".")
		persistFiles := request.GetBool("persist_files", true)
		model := request.GetString("model", "")
		sessionID := request.GetString("session", "")
		continueSession := request.GetBool("continue", false)
		// Future: share and agent options can be added when implemented

		args := []string{"opencode", "chat"}
		
		// Add working directory flag
		if workdir != "." {
			args = append(args, "--workdir", workdir)
		}
		
		// Add persistence flag
		if !persistFiles {
			args = append(args, "--ephemeral")
		}
		
		// Add model flag
		if model != "" {
			args = append(args, "--model", model)
		}
		
		// Add session support flags
		if sessionID != "" {
			args = append(args, "--session", sessionID)
		}
		if continueSession {
			args = append(args, "--continue")
		}
		
		// Add the message
		args = append(args, message)

		// Execute the command
		result, err := executeShipCommand(args)
		duration := time.Since(startTime)
		
		// Log output
		logOpenCodeInteraction("OUTPUT", "opencode_run", &request, result, err, duration)
		
		return result, err
	})

	// OpenCode Version - Get version information
	opencodeVersionTool := mcp.NewTool("opencode_version",
		mcp.WithDescription("Get OpenCode AI version information"),
	)

	s.AddTool(opencodeVersionTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"opencode", "version"}
		return executeShipCommand(args)
	})
}