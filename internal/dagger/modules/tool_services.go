package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// ToolService wraps a module as an HTTP service that the LLM can call
type ToolService struct {
	client *dagger.Client
	name   string
	port   int
}

// SteampipeService exposes Steampipe as an HTTP API service
func NewSteampipeService(client *dagger.Client) *dagger.Service {
	// Create a container that runs an HTTP API for Steampipe
	container := client.Container().
		From("golang:1.21-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithNewFile("/app/server.go", steampipeServerCode).
		WithExec([]string{"go", "run", "/app/server.go"}).
		WithExposedPort(8001).
		WithEnvVariable("SERVICE_NAME", "steampipe").
		AsService()

	return container
}

// OpenInfraQuoteService exposes cost analysis as a service
func NewOpenInfraQuoteService(client *dagger.Client) *dagger.Service {
	container := client.Container().
		From("golang:1.21-alpine").
		WithNewFile("/app/server.go", costAnalysisServerCode).
		WithExec([]string{"go", "run", "/app/server.go"}).
		WithExposedPort(8002).
		WithEnvVariable("SERVICE_NAME", "openinfraquote").
		AsService()

	return container
}

// TerraformDocsService exposes documentation generation as a service
func NewTerraformDocsService(client *dagger.Client) *dagger.Service {
	container := client.Container().
		From("golang:1.21-alpine").
		WithNewFile("/app/server.go", terraformDocsServerCode).
		WithExec([]string{"go", "run", "/app/server.go"}).
		WithExposedPort(8003).
		WithEnvVariable("SERVICE_NAME", "terraform-docs").
		AsService()

	return container
}

// InfraMapService exposes infrastructure diagram generation as a service
func NewInfraMapService(client *dagger.Client) *dagger.Service {
	container := client.Container().
		From("golang:1.21-alpine").
		WithNewFile("/app/server.go", infraMapServerCode).
		WithExec([]string{"go", "run", "/app/server.go"}).
		WithExposedPort(8005).
		WithEnvVariable("SERVICE_NAME", "inframap").
		AsService()

	return container
}

// ToolRegistryService provides a registry of all available tools for the LLM
func NewToolRegistryService(client *dagger.Client, services map[string]*dagger.Service) *dagger.Service {
	// Create a registry that knows about all services
	registryCode := generateRegistryCode(services)

	registryContainer := client.Container().
		From("golang:1.21-alpine").
		WithNewFile("/app/registry.go", registryCode).
		WithExec([]string{"go", "run", "/app/registry.go"}).
		WithExposedPort(8000)

	// Bind all services to the registry container
	for name, service := range services {
		registryContainer = registryContainer.
			WithServiceBinding(name, service)
	}

	return registryContainer.AsService()
}

// LLMWithServiceTools creates an LLM that can call services as tools
type LLMWithServiceTools struct {
	client       *dagger.Client
	model        string
	toolRegistry *dagger.Service
	services     map[string]*dagger.Service
}

// NewLLMWithServiceTools creates an LLM with access to tool services
func NewLLMWithServiceTools(client *dagger.Client, model string) *LLMWithServiceTools {
	// Create all tool services
	services := map[string]*dagger.Service{
		"steampipe":      NewSteampipeService(client),
		"openinfraquote": NewOpenInfraQuoteService(client),
		"terraform-docs": NewTerraformDocsService(client),
	}

	// Create tool registry
	registry := NewToolRegistryService(client, services)

	return &LLMWithServiceTools{
		client:       client,
		model:        model,
		toolRegistry: registry,
		services:     services,
	}
}

// InvestigateWithServices performs investigation using service-based tools
func (m *LLMWithServiceTools) InvestigateWithServices(ctx context.Context, task string) (*ServiceInvestigationReport, error) {
	// Create a container that has access to all services
	investigatorContainer := m.client.Container().
		From("alpine:latest").
		WithServiceBinding("tool-registry", m.toolRegistry)

	// Bind all tool services
	for name, service := range m.services {
		investigatorContainer = investigatorContainer.
			WithServiceBinding(name, service)
	}

	// Create LLM with instructions about available services
	systemPrompt := `You are an infrastructure investigator with access to HTTP service tools.

Available services:
1. http://steampipe:8001 - Query cloud infrastructure
   POST /query {"provider": "aws", "sql": "SELECT ..."}

2. http://openinfraquote:8002 - Analyze infrastructure costs
   POST /analyze {"plan_file": "path", "region": "us-east-1"}

3. http://terraform-docs:8003 - Generate documentation
   POST /generate {"path": "module/path", "format": "markdown"}

To use a tool, make HTTP requests to these services.
You can make multiple requests to gather information.
`

	// Create LLM that can make HTTP calls to services
	_ = m.client.LLM(dagger.LLMOpts{
		Model: m.model,
	}).
		WithSystemPrompt(systemPrompt).
		WithPrompt(fmt.Sprintf("Task: %s\n\nPlan your investigation using the available HTTP services.", task))

	// Execute investigation with service calls
	// The LLM would generate curl commands or HTTP requests to the services

	return &ServiceInvestigationReport{
		Task:         task,
		ServicesUsed: []string{"steampipe", "openinfraquote"},
		Results:      "Investigation results here",
	}, nil
}

// ServiceInvestigationReport contains results from service-based investigation
type ServiceInvestigationReport struct {
	Task         string
	ServicesUsed []string
	Results      string
}

// HTTP server code for each service
const steampipeServerCode = `
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type QueryRequest struct {
	Provider string ` + "`json:\"provider\"`" + `
	SQL      string ` + "`json:\"sql\"`" + `
}

func main() {
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		var req QueryRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		// Execute Steampipe query
		cmd := exec.Command("steampipe", "query", req.SQL, "--output", "json")
		output, _ := cmd.Output()
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	})
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Steampipe service healthy"))
	})
	
	fmt.Println("Steampipe service running on :8001")
	http.ListenAndServe(":8001", nil)
}
`

const costAnalysisServerCode = `
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type AnalyzeRequest struct {
	PlanFile string ` + "`json:\"plan_file\"`" + `
	Region   string ` + "`json:\"region\"`" + `
}

func main() {
	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		var req AnalyzeRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		// Run OpenInfraQuote analysis
		result := map[string]interface{}{
			"total_cost": "$1,234.56/month",
			"resources": []string{"EC2", "RDS", "S3"},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	
	fmt.Println("Cost analysis service running on :8002")
	http.ListenAndServe(":8002", nil)
}
`

const terraformDocsServerCode = `
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GenerateRequest struct {
	Path   string ` + "`json:\"path\"`" + `
	Format string ` + "`json:\"format\"`" + `
}

func main() {
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		var req GenerateRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		// Generate documentation
		result := map[string]string{
			"documentation": "# Module Documentation\n\nGenerated docs here...",
			"format": req.Format,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	
	fmt.Println("Terraform docs service running on :8003")
	http.ListenAndServe(":8003", nil)
}
`

const infraMapServerCode = `
package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
)

type DiagramRequest struct {
	Input    string ` + "`json:\"input\"`" + `
	Format   string ` + "`json:\"format\"`" + `
	IsHCL    bool   ` + "`json:\"is_hcl\"`" + `
	Raw      bool   ` + "`json:\"raw\"`" + `
	Provider string ` + "`json:\"provider,omitempty\"`" + `
}

type DiagramResponse struct {
	Success bool   ` + "`json:\"success\"`" + `
	Output  string ` + "`json:\"output\"`" + `
	Error   string ` + "`json:\"error,omitempty\"`" + `
}

func main() {
	http.HandleFunc("/diagram", func(w http.ResponseWriter, r *http.Request) {
		var req DiagramRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		// Build inframap command
		args := []string{"inframap", "generate"}
		
		if req.Raw {
			args = append(args, "--raw")
		}
		if req.Provider != "" {
			args = append(args, "--provider", req.Provider)
		}
		if req.IsHCL {
			args = append(args, "--hcl")
		}
		
		args = append(args, req.Input)
		
		// Execute inframap
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.Output()
		
		result := DiagramResponse{Success: err == nil}
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Output = string(output)
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	
	fmt.Println("InfraMap service running on :8005")
	http.ListenAndServe(":8005", nil)
}
`

func generateRegistryCode(services map[string]*dagger.Service) string {
	// Generate registry code that knows about all services
	return `
package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		tools := map[string]string{
			"steampipe": "http://steampipe:8001",
			"openinfraquote": "http://openinfraquote:8002", 
			"terraform-docs": "http://terraform-docs:8003",
			"inframap": "http://inframap:8005",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools)
	})
	
	http.ListenAndServe(":8000", nil)
}
`
}
