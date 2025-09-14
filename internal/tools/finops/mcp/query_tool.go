package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudshipai/ship/internal/tools/finops/cloudshipai"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/providers/stub"
)

// FinOpsQueryTool implements an agent-driven query tool for flexible finops operations
type FinOpsQueryTool struct {
	lighthouseClient interfaces.LighthouseReporter
	providers        map[interfaces.VendorType]interfaces.VendorProvider
}

// NewFinOpsQueryTool creates a new finops-query tool instance
func NewFinOpsQueryTool(config interfaces.LighthouseConfig) (*FinOpsQueryTool, error) {
	lighthouseClient := cloudshipai.NewClient(config)
	
	// Initialize providers (using stub for testing)
	providers := make(map[interfaces.VendorType]interfaces.VendorProvider)
	
	// Initialize stub providers
	providers[interfaces.VendorAWS] = stub.NewStubProvider(interfaces.VendorAWS)
	providers[interfaces.VendorGCP] = stub.NewStubProvider(interfaces.VendorGCP)
	providers[interfaces.VendorAzure] = stub.NewStubProvider(interfaces.VendorAzure)
	providers[interfaces.VendorKubernetes] = stub.NewStubProvider(interfaces.VendorKubernetes)
	
	return &FinOpsQueryTool{
		lighthouseClient: lighthouseClient,
		providers:        providers,
	}, nil
}

// QueryRequest represents a flexible query request
type QueryRequest struct {
	Query             string                 `json:"query"`                     // Natural language or structured query
	Provider          string                 `json:"provider"`                  // aws, gcp, azure, kubernetes
	Context           map[string]interface{} `json:"context,omitempty"`         // Additional context for the query
	Operations        []string               `json:"operations,omitempty"`      // ["discover", "recommend", "analyze"]
	Filters           QueryFilters           `json:"filters,omitempty"`         // Query-specific filters
	EnableCloudshipAI bool                   `json:"enable_cloudshipai,omitempty"` // enable Station API reporting
}

// QueryFilters contains various filter options
type QueryFilters struct {
	Regions       []string          `json:"regions,omitempty"`
	ResourceTypes []string          `json:"resource_types,omitempty"`
	TimeWindow    string            `json:"time_window,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	MinSavings    float64           `json:"min_savings,omitempty"`
	AccountIDs    []string          `json:"account_ids,omitempty"`
}

// QueryResponse represents the query result
type QueryResponse struct {
	Query           string                      `json:"query"`
	ParsedIntent    QueryIntent                 `json:"parsed_intent"`
	Results         QueryResults                `json:"results"`
	Summary         string                      `json:"summary"`
	Suggestions     []string                    `json:"suggestions,omitempty"`
	Lighthouse      LighthouseStatus            `json:"lighthouse"`
}

// QueryIntent represents what the agent is trying to accomplish
type QueryIntent struct {
	Operation   string                 `json:"operation"`   // discover, recommend, analyze, or mixed
	Provider    string                 `json:"provider"`
	Intent      string                 `json:"intent"`      // human-readable intent
	Parameters  map[string]interface{} `json:"parameters"`
}

// QueryResults contains the actual query results
type QueryResults struct {
	Resources       []interfaces.Resource       `json:"resources,omitempty"`
	Recommendations []interfaces.Recommendation `json:"recommendations,omitempty"`
	Opportunities   []interfaces.Opportunity    `json:"opportunities,omitempty"`
	CostData        []interfaces.CostRecord     `json:"cost_data,omitempty"`
	Insights        []interfaces.Insight        `json:"insights,omitempty"`
	TotalSavings    float64                     `json:"total_savings,omitempty"`
}

// Execute implements the MCP tool execution for agent-driven queries
func (t *FinOpsQueryTool) Execute(ctx context.Context, request json.RawMessage) (interface{}, error) {
	var req QueryRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query request: %w", err)
	}
	
	// Parse the query intent
	intent, err := t.parseQueryIntent(req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query intent: %w", err)
	}
	
	// Get the appropriate provider
	vendorType := interfaces.VendorType(intent.Provider)
	provider, exists := t.providers[vendorType]
	if !exists {
		return nil, fmt.Errorf("unsupported provider: %s", intent.Provider)
	}
	
	// Execute the query based on parsed intent
	results, err := t.executeQuery(ctx, provider, intent, req.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	
	// Generate summary and suggestions
	summary := t.generateSummary(intent, results)
	suggestions := t.generateSuggestions(intent, results)
	
	// Report results to lighthouse if applicable
	lighthouseStatus := t.reportToLighthouse(ctx, results, req.EnableCloudshipAI)
	
	response := QueryResponse{
		Query:        req.Query,
		ParsedIntent: intent,
		Results:      results,
		Summary:      summary,
		Suggestions:  suggestions,
		Lighthouse:   lighthouseStatus,
	}
	
	return response, nil
}

// parseQueryIntent analyzes the query and determines what operation(s) to perform
func (t *FinOpsQueryTool) parseQueryIntent(req QueryRequest) (QueryIntent, error) {
	query := strings.ToLower(req.Query)
	
	intent := QueryIntent{
		Provider:   req.Provider,
		Parameters: make(map[string]interface{}),
	}
	
	// Simple intent parsing (in production, this could use NLP/LLM)
	switch {
	case strings.Contains(query, "discover") || strings.Contains(query, "find") || strings.Contains(query, "list"):
		intent.Operation = "discover"
		intent.Intent = "Discover cloud resources"
		
	case strings.Contains(query, "recommend") || strings.Contains(query, "optimize") || strings.Contains(query, "savings"):
		intent.Operation = "recommend"
		intent.Intent = "Generate optimization recommendations"
		
	case strings.Contains(query, "analyze") || strings.Contains(query, "cost") || strings.Contains(query, "spend"):
		intent.Operation = "analyze"
		intent.Intent = "Analyze cost data and trends"
		
	case strings.Contains(query, "expensive") || strings.Contains(query, "highest cost"):
		intent.Operation = "mixed"
		intent.Intent = "Find most expensive resources and optimization opportunities"
		
	default:
		// Default to discovery for ambiguous queries
		intent.Operation = "discover"
		intent.Intent = "Discover cloud resources (default operation)"
	}
	
	// Extract parameters from query
	if strings.Contains(query, "ec2") || strings.Contains(query, "instance") {
		intent.Parameters["resource_types"] = []string{"compute"}
	}
	if strings.Contains(query, "storage") || strings.Contains(query, "ebs") {
		intent.Parameters["resource_types"] = []string{"storage"}
	}
	if strings.Contains(query, "production") {
		intent.Parameters["tags"] = map[string]string{"Environment": "production"}
	}
	
	return intent, nil
}

// executeQuery performs the actual query operation
func (t *FinOpsQueryTool) executeQuery(ctx context.Context, provider interfaces.VendorProvider, intent QueryIntent, filters QueryFilters) (QueryResults, error) {
	var results QueryResults
	
	switch intent.Operation {
	case "discover":
		opts := interfaces.DiscoveryOptions{
			ResourceTypes: filters.ResourceTypes,
			Tags:          filters.Tags,
			AccountIDs:    filters.AccountIDs,
		}
		if len(filters.Regions) > 0 {
			opts.Region = filters.Regions[0] // Use first region for now
		}
		
		resources, err := provider.DiscoverResources(ctx, opts)
		if err != nil {
			return results, err
		}
		results.Resources = resources
		
	case "recommend":
		opts := interfaces.RecommendationOptions{
			FindingTypes: []string{"rightsizing"},
			Regions:      filters.Regions,
			AccountIDs:   filters.AccountIDs,
		}
		
		recommendations, err := provider.GetRecommendations(ctx, opts)
		if err != nil {
			return results, err
		}
		results.Recommendations = recommendations
		
		// Transform to opportunities
		opportunities := t.transformToOpportunities(recommendations)
		results.Opportunities = opportunities
		results.TotalSavings = t.calculateTotalSavings(opportunities)
		
	case "analyze":
		opts := interfaces.CostOptions{
			TimeWindow:  filters.TimeWindow,
			Granularity: "daily",
			Currency:    "USD",
		}
		if opts.TimeWindow == "" {
			opts.TimeWindow = "30d"
		}
		
		costData, err := provider.GetCostData(ctx, opts)
		if err != nil {
			return results, err
		}
		results.CostData = costData
		
	case "mixed":
		// Execute multiple operations
		// First discover resources
		discoverOpts := interfaces.DiscoveryOptions{
			ResourceTypes: filters.ResourceTypes,
			Tags:          filters.Tags,
		}
		resources, _ := provider.DiscoverResources(ctx, discoverOpts)
		results.Resources = resources
		
		// Then get recommendations
		recommendOpts := interfaces.RecommendationOptions{
			FindingTypes: []string{"rightsizing"},
			Regions:      filters.Regions,
		}
		recommendations, _ := provider.GetRecommendations(ctx, recommendOpts)
		results.Recommendations = recommendations
		results.Opportunities = t.transformToOpportunities(recommendations)
		results.TotalSavings = t.calculateTotalSavings(results.Opportunities)
	}
	
	return results, nil
}

// Helper methods (reuse from other tools)
func (t *FinOpsQueryTool) transformToOpportunities(recommendations []interfaces.Recommendation) []interfaces.Opportunity {
	// Implementation similar to recommend_tool.go
	var opportunities []interfaces.Opportunity
	for _, rec := range recommendations {
		opp := interfaces.Opportunity{
			ID:               rec.ID,
			ResourceID:       rec.ResourceID,
			ResourceARN:      rec.ResourceARN,
			Type:             t.mapRecommendationTypeToOpportunityType(rec.RecommendationType),
			Vendor:           rec.Provider,
			Title:            rec.Title,
			Description:      rec.Description,
			EstimatedSavings: rec.EstimatedSavings.Monthly,
			Currency:         rec.EstimatedSavings.Currency,
			CreatedAt:        rec.CreatedAt,
		}
		opportunities = append(opportunities, opp)
	}
	return opportunities
}

func (t *FinOpsQueryTool) mapRecommendationTypeToOpportunityType(recType string) interfaces.OpportunityType {
	switch recType {
	case "RightSizeInstance", "UpgradeInstanceGeneration":
		return interfaces.OpportunityRightsizing
	case "TerminateInstance":
		return interfaces.OpportunityIdleResource
	default:
		return interfaces.OpportunityRightsizing
	}
}

func (t *FinOpsQueryTool) calculateTotalSavings(opportunities []interfaces.Opportunity) float64 {
	var total float64
	for _, opp := range opportunities {
		total += opp.EstimatedSavings
	}
	return total
}

func (t *FinOpsQueryTool) generateSummary(intent QueryIntent, results QueryResults) string {
	switch intent.Operation {
	case "discover":
		return fmt.Sprintf("Found %d resources in %s", len(results.Resources), intent.Provider)
	case "recommend":
		return fmt.Sprintf("Generated %d recommendations with potential savings of $%.2f/month", 
			len(results.Recommendations), results.TotalSavings)
	case "analyze":
		return fmt.Sprintf("Analyzed cost data with %d records", len(results.CostData))
	default:
		return fmt.Sprintf("Executed %s operation on %s", intent.Operation, intent.Provider)
	}
}

func (t *FinOpsQueryTool) generateSuggestions(intent QueryIntent, results QueryResults) []string {
	var suggestions []string
	
	if len(results.Recommendations) > 0 {
		suggestions = append(suggestions, "Consider implementing the top recommendations for maximum savings")
	}
	if len(results.Resources) > 10 {
		suggestions = append(suggestions, "Large number of resources found - consider filtering by tags or regions")
	}
	if results.TotalSavings > 1000 {
		suggestions = append(suggestions, "Significant savings potential identified - prioritize high-impact optimizations")
	}
	
	return suggestions
}

func (t *FinOpsQueryTool) reportToLighthouse(ctx context.Context, results QueryResults, enableCloudshipAI bool) LighthouseStatus {
	status := LighthouseStatus{Reported: false}
	
	if enableCloudshipAI && len(results.Opportunities) > 0 {
		if err := t.lighthouseClient.ReportOpportunities(ctx, results.Opportunities); err != nil {
			status.Error = err.Error()
		} else {
			status.Reported = true
		}
	} else if !enableCloudshipAI {
		status.Disabled = true
	}
	
	return status
}

// GetName returns the tool name
func (t *FinOpsQueryTool) GetName() string {
	return "finops-query"
}

// GetDescription returns the tool description
func (t *FinOpsQueryTool) GetDescription() string {
	return "Agent-driven query tool for flexible finops operations with natural language support"
}

// GetInputSchema returns the JSON schema for the tool input
func (t *FinOpsQueryTool) GetInputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Natural language or structured query (e.g., 'find overprovisioned EC2 instances')",
			},
			"provider": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"aws", "gcp", "azure", "kubernetes"},
				"description": "Cloud provider to query",
			},
			"context": map[string]interface{}{
				"type":        "object",
				"description": "Additional context for the query",
			},
			"operations": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
					"enum": []string{"discover", "recommend", "analyze"},
				},
				"description": "Specific operations to perform",
			},
			"filters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"regions": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{"type": "string"},
					},
					"resource_types": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{"type": "string"},
					},
					"time_window": map[string]interface{}{
						"type": "string",
					},
					"tags": map[string]interface{}{
						"type": "object",
					},
					"min_savings": map[string]interface{}{
						"type": "number",
					},
					"account_ids": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{"type": "string"},
					},
				},
			},
			"enable_cloudshipai": map[string]interface{}{
				"type":        "boolean",
				"description": "Enable reporting results to CloudshipAI Station API",
				"default":     false,
			},
		},
		"required": []string{"query", "provider"},
	}
}