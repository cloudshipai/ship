package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type MCPRequest struct {
	Tool    string            `json:"tool"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type MCPResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	r := mux.NewRouter()

	// CORS middleware
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	// List available tools
	r.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tools := map[string]interface{}{
			"categories": map[string][]string{
				"security":     {"gitleaks", "trivy", "grype", "syft", "checkov", "terrascan", "semgrep"},
				"terraform":    {"tflint", "terraform-docs", "inframap", "iac-plan"},
				"kubernetes":   {"kube-bench", "kube-hunter", "falco", "kubescape"},
				"aws":          {"cloudsplaining", "parliament", "pmapper", "policy-sentry", "prowler"},
				"supply-chain": {"cosign", "guac", "rekor", "in-toto", "slsa-verifier"},
			},
			"collections": []string{"all", "security", "terraform", "kubernetes", "cloud", "aws", "supply-chain"},
		}

		json.NewEncoder(w).Encode(tools)
	}).Methods("GET")

	// Execute MCP command
	r.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req MCPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// For MCP commands, we need to handle them differently
		if req.Tool == "mcp" || strings.HasPrefix(req.Tool, "mcp") {
			response := MCPResponse{
				Success: false,
				Error:   "MCP commands are long-running servers and cannot be executed via HTTP. Use the MCP Inspector instead.",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// For individual tools, try to get help or version info
		args := []string{}
		if req.Tool != "" {
			args = append(args, req.Tool)
		}
		if req.Command != "" {
			args = append(args, req.Command)
		}
		args = append(args, req.Args...)

		// Add --help if no args provided to avoid hanging
		if len(args) == 1 && len(req.Args) == 0 {
			args = append(args, "--help")
		}

		// Execute ship command with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "ship", args...)

		// Set environment variables
		if req.Env != nil {
			for k, v := range req.Env {
				cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", k, v))
			}
		}

		// Capture output
		output, err := cmd.CombinedOutput()

		if err != nil {
			// Check if it's a timeout error
			if ctx.Err() == context.DeadlineExceeded {
				response := MCPResponse{
					Success: false,
					Error:   fmt.Sprintf("Command timed out after 10 seconds. This tool may be designed to run as a long-running server.\n\nTry using the MCP Inspector at http://localhost:6274 instead.\n\nCommand attempted: ship %s", strings.Join(args, " ")),
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			response := MCPResponse{
				Success: false,
				Error:   fmt.Sprintf("Command failed: %v\nOutput: %s", err, string(output)),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		response := MCPResponse{
			Success: true,
			Data: map[string]string{
				"output":  string(output),
				"command": fmt.Sprintf("ship %s", strings.Join(args, " ")),
			},
		}

		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	// Test MCP server connection
	r.HandleFunc("/test-mcp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Test if we can start an MCP server briefly
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "ship", "mcp", "security", "--help")
		output, err := cmd.CombinedOutput()

		if err != nil {
			response := MCPResponse{
				Success: false,
				Error:   fmt.Sprintf("MCP test failed: %v\nOutput: %s", err, string(output)),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		response := MCPResponse{
			Success: true,
			Data: map[string]string{
				"output":  string(output),
				"message": "MCP server test successful. Use the MCP Inspector to connect to running MCP servers.",
			},
		}

		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	// Serve static files for a simple web interface
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	// Apply CORS middleware
	handler := corsMiddleware(r)

	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	fmt.Printf("🚀 Ship MCP HTTP Server running on http://localhost%s\n", port)
	fmt.Printf("📚 Available endpoints:\n")
	fmt.Printf("   GET  /health     - Health check\n")
	fmt.Printf("   GET  /tools      - List available tools\n")
	fmt.Printf("   POST /execute    - Execute ship command (non-MCP)\n")
	fmt.Printf("   GET  /test-mcp   - Test MCP server connection\n")
	fmt.Printf("   GET  /           - Web interface\n")
	fmt.Printf("\n⚠️  Note: MCP commands are long-running servers.\n")
	fmt.Printf("   Use the MCP Inspector at http://localhost:6274 for MCP testing.\n")

	log.Fatal(http.ListenAndServe(port, handler))
}
