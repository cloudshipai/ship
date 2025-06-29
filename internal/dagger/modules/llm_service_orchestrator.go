package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// LLMServiceOrchestrator uses Dagger services as tools via HTTP
type LLMServiceOrchestrator struct {
	client   *dagger.Client
	model    string
	services map[string]*dagger.Service
}

// NewLLMServiceOrchestrator creates an orchestrator with service-based tools
func NewLLMServiceOrchestrator(client *dagger.Client, model string) *LLMServiceOrchestrator {
	return &LLMServiceOrchestrator{
		client:   client,
		model:    model,
		services: make(map[string]*dagger.Service),
	}
}

// StartToolServices starts all tool services and returns their endpoints
func (o *LLMServiceOrchestrator) StartToolServices(ctx context.Context) (map[string]string, error) {
	endpoints := make(map[string]string)

	// Start Steampipe service
	steampipeService := o.createSteampipeService()
	o.services["steampipe"] = steampipeService
	endpoints["steampipe"] = "http://steampipe:8001"

	// Start OpenInfraQuote service
	costService := o.createCostAnalysisService()
	o.services["cost-analysis"] = costService
	endpoints["cost-analysis"] = "http://cost-analysis:8002"

	// Start Terraform Docs service
	docsService := o.createDocsService()
	o.services["terraform-docs"] = docsService
	endpoints["terraform-docs"] = "http://terraform-docs:8003"

	// Start Security Scan service
	securityService := o.createSecurityScanService()
	o.services["security-scan"] = securityService
	endpoints["security-scan"] = "http://security-scan:8004"

	// Start InfraMap service
	infraMapService := o.createInfraMapService()
	o.services["inframap"] = infraMapService
	endpoints["inframap"] = "http://inframap:8005"

	return endpoints, nil
}

// createSteampipeService creates a real Steampipe service
func (o *LLMServiceOrchestrator) createSteampipeService() *dagger.Service {
	// Create actual Steampipe service with HTTP API wrapper
	return o.client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "plugin", "install", "aws"}).
		WithNewFile("/app/api.py", `
import json
import subprocess
from http.server import HTTPServer, BaseHTTPRequestHandler

class SteampipeHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/query':
            content_length = int(self.headers['Content-Length'])
            body = self.rfile.read(content_length)
            data = json.loads(body)
            
            # Execute Steampipe query
            cmd = ['steampipe', 'query', data['sql'], '--output', 'json']
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(result.stdout.encode())

    def do_GET(self):
        if self.path == '/health':
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'Steampipe service healthy')

httpd = HTTPServer(('0.0.0.0', 8001), SteampipeHandler)
print('Steampipe API running on :8001')
httpd.serve_forever()
`).
		WithExec([]string{"python3", "/app/api.py"}).
		WithExposedPort(8001).
		AsService()
}

// createCostAnalysisService creates OpenInfraQuote service
func (o *LLMServiceOrchestrator) createCostAnalysisService() *dagger.Service {
	return o.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithNewFile("/app/api.py", `
import json
import subprocess
from http.server import HTTPServer, BaseHTTPRequestHandler

class CostHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/analyze':
            content_length = int(self.headers['Content-Length'])
            body = self.rfile.read(content_length)
            data = json.loads(body)
            
            # Run cost analysis
            cmd = ['oiq', 'analyze', data.get('plan_file', 'tfplan.json')]
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            
            response = {
                'status': 'success' if result.returncode == 0 else 'error',
                'output': result.stdout,
                'error': result.stderr
            }
            self.wfile.write(json.dumps(response).encode())

httpd = HTTPServer(('0.0.0.0', 8002), CostHandler)
print('Cost Analysis API running on :8002')
httpd.serve_forever()
`).
		WithExec([]string{"python3", "/app/api.py"}).
		WithExposedPort(8002).
		AsService()
}

// createDocsService creates Terraform documentation service
func (o *LLMServiceOrchestrator) createDocsService() *dagger.Service {
	return o.client.Container().
		From("quay.io/terraform-docs/terraform-docs:latest").
		WithNewFile("/app/api.sh", `#!/bin/sh
# Simple HTTP API for terraform-docs
while true; do
  echo -e "HTTP/1.1 200 OK\n\n{\"docs\": \"Generated documentation\"}" | nc -l -p 8003
done
`, dagger.ContainerWithNewFileOpts{
			Permissions: 0755,
		}).
		WithExec([]string{"/app/api.sh"}).
		WithExposedPort(8003).
		AsService()
}

// createSecurityScanService creates security scanning service
func (o *LLMServiceOrchestrator) createSecurityScanService() *dagger.Service {
	return o.client.Container().
		From("bridgecrew/checkov:latest").
		WithNewFile("/app/api.py", `
import json
import subprocess
from http.server import HTTPServer, BaseHTTPRequestHandler

class SecurityHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/scan':
            content_length = int(self.headers['Content-Length'])
            body = self.rfile.read(content_length)
            data = json.loads(body)
            
            # Run security scan
            path = data.get('path', '/workspace')
            cmd = ['checkov', '-d', path, '--output', 'json']
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            self.wfile.write(result.stdout.encode())

httpd = HTTPServer(('0.0.0.0', 8004), SecurityHandler)
print('Security Scan API running on :8004')
httpd.serve_forever()
`).
		WithExec([]string{"python", "/app/api.py"}).
		WithExposedPort(8004).
		AsService()
}

// createInfraMapService creates a real InfraMap service
func (o *LLMServiceOrchestrator) createInfraMapService() *dagger.Service {
	return o.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "graphviz", "font-noto", "python3", "py3-pip"}).
		WithExec([]string{"sh", "-c", "curl -sSL https://github.com/cycloidio/inframap/releases/latest/download/inframap-linux-amd64.tar.gz | tar xz -C /usr/local/bin/"}).
		WithExec([]string{"chmod", "+x", "/usr/local/bin/inframap"}).
		WithNewFile("/app/api.py", `
import json
import subprocess
from http.server import HTTPServer, BaseHTTPRequestHandler

class InfraMapHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/diagram':
            content_length = int(self.headers['Content-Length'])
            body = self.rfile.read(content_length)
            data = json.loads(body)
            
            # Build inframap command
            cmd = ['inframap', 'generate']
            
            if data.get('raw', False):
                cmd.append('--raw')
            if data.get('provider'):
                cmd.extend(['--provider', data['provider']])
            if data.get('is_hcl', False):
                cmd.append('--hcl')
                
            cmd.append(data.get('input', 'terraform.tfstate'))
            
            # Add output format conversion if needed
            format = data.get('format', 'dot')
            if format != 'dot':
                # Pipe to dot for conversion
                result = subprocess.run(cmd, capture_output=True, text=True)
                if result.returncode == 0:
                    dot_cmd = ['dot', f'-T{format}']
                    final_result = subprocess.run(dot_cmd, input=result.stdout, capture_output=True, text=True)
                    output = final_result.stdout
                    error = final_result.stderr
                else:
                    output = result.stdout
                    error = result.stderr
            else:
                result = subprocess.run(cmd, capture_output=True, text=True)
                output = result.stdout
                error = result.stderr
            
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            response = {
                'success': result.returncode == 0,
                'output': output,
                'error': error
            }
            self.wfile.write(json.dumps(response).encode())

    def do_GET(self):
        if self.path == '/health':
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'InfraMap service healthy')

httpd = HTTPServer(('0.0.0.0', 8005), InfraMapHandler)
print('InfraMap API running on :8005')
httpd.serve_forever()
`).
		WithExec([]string{"python3", "/app/api.py"}).
		WithExposedPort(8005).
		AsService()
}

// ExecuteWithServices runs an LLM investigation using services
func (o *LLMServiceOrchestrator) ExecuteWithServices(ctx context.Context, task string) (*ServiceExecutionReport, error) {
	// Start all services
	endpoints, err := o.StartToolServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start services: %w", err)
	}

	// Create a container that can access all services
	orchestratorContainer := o.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq"})

	// Bind all services to the orchestrator
	for name, service := range o.services {
		orchestratorContainer = orchestratorContainer.
			WithServiceBinding(name, service)
	}

	// Create LLM with knowledge of service endpoints
	toolsDescription := o.generateToolsDescription(endpoints)

	llm := o.client.LLM(dagger.LLMOpts{
		Model: o.model,
	}).
		WithSystemPrompt(fmt.Sprintf(`You are an infrastructure investigator with access to HTTP API tools.

%s

To use a tool, generate curl commands that the orchestrator will execute.
Format your tool requests as JSON:
{
  "tool": "service-name",
  "endpoint": "/path",
  "method": "POST",
  "data": {...}
}

After each tool use, analyze the results and decide next steps.`, toolsDescription)).
		WithPrompt(fmt.Sprintf("Task: %s\n\nBegin your investigation.", task))

	// Execute investigation loop
	var toolUses []ToolUse
	for i := 0; i < 5; i++ {
		// Get LLM's next action
		synced, err := llm.Sync(ctx)
		if err != nil {
			return nil, err
		}

		response, err := synced.LastReply(ctx)
		if err != nil {
			return nil, err
		}

		// Check if LLM wants to use a tool
		if strings.Contains(response, `"tool":`) {
			// Parse tool request
			var toolRequest ToolRequest
			// Extract JSON from response
			start := strings.Index(response, "{")
			end := strings.LastIndex(response, "}") + 1
			if start >= 0 && end > start {
				jsonStr := response[start:end]
				if err := json.Unmarshal([]byte(jsonStr), &toolRequest); err == nil {
					// Execute tool via HTTP in the orchestrator container
					result, err := o.executeToolInContainer(ctx, orchestratorContainer, toolRequest, endpoints)

					toolUse := ToolUse{
						Tool:     toolRequest.Tool,
						Endpoint: toolRequest.Endpoint,
						Result:   result,
						Error:    err,
					}
					toolUses = append(toolUses, toolUse)

					// Feed result back to LLM
					llm = llm.WithPrompt(fmt.Sprintf("Tool result:\n%s\n\nContinue investigation.", result))
				}
			}
		} else {
			// LLM provided final analysis
			return &ServiceExecutionReport{
				Task:      task,
				ToolUses:  toolUses,
				Analysis:  response,
				Endpoints: endpoints,
			}, nil
		}
	}

	// Get final analysis
	finalLLM := llm.WithPrompt("Provide your final analysis and recommendations.")
	synced, _ := finalLLM.Sync(ctx)
	analysis, _ := synced.LastReply(ctx)

	return &ServiceExecutionReport{
		Task:      task,
		ToolUses:  toolUses,
		Analysis:  analysis,
		Endpoints: endpoints,
	}, nil
}

// executeToolInContainer executes HTTP request to tool service
func (o *LLMServiceOrchestrator) executeToolInContainer(
	ctx context.Context,
	container *dagger.Container,
	request ToolRequest,
	endpoints map[string]string,
) (string, error) {
	endpoint := endpoints[request.Tool]
	if endpoint == "" {
		return "", fmt.Errorf("unknown tool: %s", request.Tool)
	}

	// Build curl command
	dataJSON, _ := json.Marshal(request.Data)
	curlCmd := []string{
		"curl", "-s", "-X", request.Method,
		"-H", "Content-Type: application/json",
		"-d", string(dataJSON),
		endpoint + request.Endpoint,
	}

	// Execute in container with service bindings
	result, err := container.
		WithExec(curlCmd).
		Stdout(ctx)

	return result, err
}

// generateToolsDescription creates documentation for available tools
func (o *LLMServiceOrchestrator) generateToolsDescription(endpoints map[string]string) string {
	return fmt.Sprintf(`Available HTTP API Tools:

1. Steampipe Query Service (%s)
   POST /query - Execute SQL queries on cloud infrastructure
   Body: {"sql": "SELECT * FROM aws_s3_bucket", "provider": "aws"}

2. Cost Analysis Service (%s)
   POST /analyze - Analyze infrastructure costs
   Body: {"plan_file": "/path/to/tfplan.json", "region": "us-east-1"}

3. Documentation Service (%s)
   POST /generate - Generate Terraform documentation
   Body: {"path": "/module/path", "format": "markdown"}

4. Security Scan Service (%s)
   POST /scan - Scan for security issues
   Body: {"path": "/code/path", "framework": "terraform"}

5. InfraMap Diagram Service (%s)
   POST /diagram - Generate infrastructure diagrams
   Body: {"input": "terraform.tfstate", "format": "png", "is_hcl": false, "raw": false}`,
		endpoints["steampipe"],
		endpoints["cost-analysis"],
		endpoints["terraform-docs"],
		endpoints["security-scan"],
		endpoints["inframap"])
}

// ToolRequest represents an LLM's request to use a tool
type ToolRequest struct {
	Tool     string                 `json:"tool"`
	Endpoint string                 `json:"endpoint"`
	Method   string                 `json:"method"`
	Data     map[string]interface{} `json:"data"`
}

// ToolUse records a tool usage
type ToolUse struct {
	Tool     string
	Endpoint string
	Result   string
	Error    error
}

// ServiceExecutionReport contains the complete execution report
type ServiceExecutionReport struct {
	Task      string
	ToolUses  []ToolUse
	Analysis  string
	Endpoints map[string]string
}
