package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudshipai/ship/internal/tools/finops/cloudshipai"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/providers/aws"
)

// FinOpsDiscoverTool implements the finops-discover MCP tool
type FinOpsDiscoverTool struct {
	lighthouseClient interfaces.LighthouseReporter
	providers        map[interfaces.VendorType]interfaces.VendorProvider
}

// NewFinOpsDiscoverTool creates a new finops-discover tool instance
func NewFinOpsDiscoverTool(config interfaces.LighthouseConfig) (*FinOpsDiscoverTool, error) {
	// Initialize lighthouse client for reporting to Station/CloudshipAI
	lighthouseClient := cloudshipai.NewClient(config)
	
	// Initialize real providers
	providers := make(map[interfaces.VendorType]interfaces.VendorProvider)
	
	// Initialize real AWS provider
	awsProvider, err := aws.NewProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AWS provider: %w", err)
	}
	providers[interfaces.VendorAWS] = awsProvider
	
	// TODO: Add real GCP, Azure, Kubernetes providers when implemented
	// For now we only support AWS
	
	return &FinOpsDiscoverTool{
		lighthouseClient: lighthouseClient,
		providers:        providers,
	}, nil
}

// DiscoverRequest represents the input parameters for the finops-discover tool
type DiscoverRequest struct {
	Provider           string            `json:"provider"`                    // aws, gcp, azure, kubernetes
	Region             string            `json:"region,omitempty"`             // optional region filter  
	ResourceTypes      []string          `json:"resource_types,omitempty"`    // optional: ["compute", "storage"]
	Tags               map[string]string `json:"tags,omitempty"`               // optional tag filters
	AccountIDs         []string          `json:"account_ids,omitempty"`       // optional account filter
	ARNs               []string          `json:"arns,omitempty"`               // optional ARN filter (OpenOps pattern)
	EnableCloudshipAI  bool              `json:"enable_cloudshipai,omitempty"` // enable Station API reporting
}

// DiscoverResponse represents the output of the finops-discover tool
type DiscoverResponse struct {
	Resources         []interfaces.Resource `json:"resources"`
	TotalCount        int                   `json:"total_count"`
	Provider          string                `json:"provider"`
	Region            string                `json:"region,omitempty"`
	Lighthouse        LighthouseStatus      `json:"lighthouse"`
}

// LighthouseStatus indicates whether data was successfully reported to lighthouse
type LighthouseStatus struct {
	Reported bool   `json:"reported"`
	Error    string `json:"error,omitempty"`
	Disabled bool   `json:"disabled,omitempty"` // true when CloudshipAI reporting is disabled
}

// Execute implements the MCP tool execution for resource discovery
func (t *FinOpsDiscoverTool) Execute(ctx context.Context, request json.RawMessage) (interface{}, error) {
	var req DiscoverRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discover request: %w", err)
	}
	
	// Validate required parameters
	if req.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	
	// Get the appropriate provider
	vendorType := interfaces.VendorType(req.Provider)
	provider, exists := t.providers[vendorType]
	if !exists {
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
	
	// Build discovery options
	opts := interfaces.DiscoveryOptions{
		Region:        req.Region,
		ResourceTypes: req.ResourceTypes,
		ARNs:          req.ARNs,
		Tags:          req.Tags,
	}
	
	// Handle multi-account discovery if account IDs specified
	if len(req.AccountIDs) > 0 {
		return t.discoverMultiAccount(ctx, provider, opts, req)
	}
	
	// Discover resources
	resources, err := provider.DiscoverResources(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to discover resources: %w", err)
	}
	
	// Report resources to lighthouse (Station/CloudshipAI) only if enabled
	lighthouseStatus := LighthouseStatus{Reported: false}
	if req.EnableCloudshipAI && len(resources) > 0 {
		if err := t.lighthouseClient.ReportResources(ctx, resources); err != nil {
			lighthouseStatus.Error = err.Error()
		} else {
			lighthouseStatus.Reported = true
		}
	} else if !req.EnableCloudshipAI {
		lighthouseStatus.Disabled = true
	}
	
	// Return response
	response := DiscoverResponse{
		Resources:  resources,
		TotalCount: len(resources),
		Provider:   req.Provider,
		Region:     req.Region,
		Lighthouse: lighthouseStatus,
	}
	
	return response, nil
}

// discoverMultiAccount handles discovery across multiple AWS accounts (OpenOps pattern)
func (t *FinOpsDiscoverTool) discoverMultiAccount(ctx context.Context, provider interfaces.VendorProvider, opts interfaces.DiscoveryOptions, req DiscoverRequest) (interface{}, error) {
	var allResources []interfaces.Resource
	
	// Discover resources for each account
	for _, accountID := range req.AccountIDs {
		// Set account ID for this iteration
		accountOpts := opts
		accountOpts.AccountID = accountID
		
		resources, err := provider.DiscoverResources(ctx, accountOpts)
		if err != nil {
			// Log error but continue with other accounts
			fmt.Printf("Warning: failed to discover resources for account %s: %v\n", accountID, err)
			continue
		}
		
		allResources = append(allResources, resources...)
	}
	
	// Report all resources to lighthouse only if enabled
	lighthouseStatus := LighthouseStatus{Reported: false}
	if req.EnableCloudshipAI && len(allResources) > 0 {
		if err := t.lighthouseClient.ReportResources(ctx, allResources); err != nil {
			lighthouseStatus.Error = err.Error()
		} else {
			lighthouseStatus.Reported = true
		}
	} else if !req.EnableCloudshipAI {
		lighthouseStatus.Disabled = true
	}
	
	response := DiscoverResponse{
		Resources:  allResources,
		TotalCount: len(allResources),
		Provider:   req.Provider,
		Region:     req.Region,
		Lighthouse: lighthouseStatus,
	}
	
	return response, nil
}

// GetName returns the tool name
func (t *FinOpsDiscoverTool) GetName() string {
	return "finops-discover"
}

// GetDescription returns the tool description
func (t *FinOpsDiscoverTool) GetDescription() string {
	return "Discovers cloud resources across providers (AWS, GCP, Azure, Kubernetes) with cost and utilization data"
}

// GetInputSchema returns the JSON schema for the tool input
func (t *FinOpsDiscoverTool) GetInputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"provider": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"aws", "gcp", "azure", "kubernetes"},
				"description": "Cloud provider to discover resources from",
			},
			"region": map[string]interface{}{
				"type":        "string",
				"description": "Optional region filter (provider-specific)",
			},
			"resource_types": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
					"enum": []string{"compute", "storage", "database", "serverless", "network", "cache", "kubernetes"},
				},
				"description": "Optional resource types to discover",
			},
			"tags": map[string]interface{}{
				"type":        "object",
				"description": "Optional tag filters (key-value pairs)",
			},
			"account_ids": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional account IDs for multi-account discovery",
			},
			"arns": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional resource ARNs to filter by (AWS only)",
			},
			"enable_cloudshipai": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable reporting results to CloudshipAI Station API",
				"default":     false,
			},
		},
		"required": []string{"provider"},
	}
}