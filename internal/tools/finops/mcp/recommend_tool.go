package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudshipai/ship/internal/tools/finops/cloudshipai"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/providers/stub"
)

// FinOpsRecommendTool implements the finops-recommend MCP tool
type FinOpsRecommendTool struct {
	lighthouseClient interfaces.LighthouseReporter
	providers        map[interfaces.VendorType]interfaces.VendorProvider
}

// NewFinOpsRecommendTool creates a new finops-recommend tool instance
func NewFinOpsRecommendTool(config interfaces.LighthouseConfig) (*FinOpsRecommendTool, error) {
	// Initialize lighthouse client for reporting to Station/CloudshipAI
	lighthouseClient := cloudshipai.NewClient(config)
	
	// Initialize providers (using stub for testing)
	providers := make(map[interfaces.VendorType]interfaces.VendorProvider)
	
	// Initialize stub providers
	providers[interfaces.VendorAWS] = stub.NewStubProvider(interfaces.VendorAWS)
	providers[interfaces.VendorGCP] = stub.NewStubProvider(interfaces.VendorGCP)
	providers[interfaces.VendorAzure] = stub.NewStubProvider(interfaces.VendorAzure)
	providers[interfaces.VendorKubernetes] = stub.NewStubProvider(interfaces.VendorKubernetes)
	
	return &FinOpsRecommendTool{
		lighthouseClient: lighthouseClient,
		providers:        providers,
	}, nil
}

// RecommendRequest represents the input parameters for the finops-recommend tool
type RecommendRequest struct {
	Provider          string   `json:"provider"`                    // aws, gcp, azure, kubernetes
	FindingTypes      []string `json:"finding_types"`              // ["rightsizing", "reserved_instances", "spot_instances"]
	Regions           []string `json:"regions,omitempty"`          // optional region filter
	ARNs              []string `json:"arns,omitempty"`             // optional ARN filter (OpenOps pattern)  
	AccountIDs        []string `json:"account_ids,omitempty"`      // optional account filter
	MinSavings        float64  `json:"min_savings,omitempty"`      // minimum monthly savings threshold
	EnableCloudshipAI bool     `json:"enable_cloudshipai,omitempty"` // enable Station API reporting
}

// RecommendResponse represents the output of the finops-recommend tool
type RecommendResponse struct {
	Recommendations  []interfaces.Recommendation `json:"recommendations"`
	Opportunities    []interfaces.Opportunity    `json:"opportunities"`
	TotalCount       int                         `json:"total_count"`
	TotalSavings     float64                     `json:"total_monthly_savings_usd"`
	Provider         string                      `json:"provider"`
	FindingTypes     []string                    `json:"finding_types"`
	Lighthouse       LighthouseStatus            `json:"lighthouse"`
}

// Execute implements the MCP tool execution for generating recommendations
func (t *FinOpsRecommendTool) Execute(ctx context.Context, request json.RawMessage) (interface{}, error) {
	var req RecommendRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recommend request: %w", err)
	}
	
	// Validate required parameters
	if req.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	if len(req.FindingTypes) == 0 {
		req.FindingTypes = []string{"rightsizing"} // Default to rightsizing
	}
	
	// Get the appropriate provider
	vendorType := interfaces.VendorType(req.Provider)
	provider, exists := t.providers[vendorType]
	if !exists {
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
	
	// Build recommendation options based on OpenOps patterns
	opts := interfaces.RecommendationOptions{
		FindingTypes: req.FindingTypes,
		Regions:      req.Regions,
		ARNs:         req.ARNs,
		AccountIDs:   req.AccountIDs,
	}
	
	// Get recommendations from provider
	recommendations, err := provider.GetRecommendations(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}
	
	// Filter by minimum savings if specified
	if req.MinSavings > 0 {
		recommendations = t.filterByMinSavings(recommendations, req.MinSavings)
	}
	
	// Transform recommendations to opportunities for CloudshipAI
	opportunities := t.transformToOpportunities(recommendations)
	
	// Calculate total savings
	totalSavings := t.calculateTotalSavings(opportunities)
	
	// Report opportunities to lighthouse (Station/CloudshipAI) only if enabled
	lighthouseStatus := LighthouseStatus{Reported: false}
	if req.EnableCloudshipAI && len(opportunities) > 0 {
		if err := t.lighthouseClient.ReportOpportunities(ctx, opportunities); err != nil {
			lighthouseStatus.Error = err.Error()
		} else {
			lighthouseStatus.Reported = true
		}
	} else if !req.EnableCloudshipAI {
		lighthouseStatus.Disabled = true
	}
	
	// Return response
	response := RecommendResponse{
		Recommendations: recommendations,
		Opportunities:   opportunities,
		TotalCount:      len(recommendations),
		TotalSavings:    totalSavings,
		Provider:        req.Provider,
		FindingTypes:    req.FindingTypes,
		Lighthouse:      lighthouseStatus,
	}
	
	return response, nil
}

// filterByMinSavings filters recommendations by minimum monthly savings
func (t *FinOpsRecommendTool) filterByMinSavings(recommendations []interfaces.Recommendation, minSavings float64) []interfaces.Recommendation {
	var filtered []interfaces.Recommendation
	for _, rec := range recommendations {
		if rec.EstimatedSavings.Monthly >= minSavings {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

// transformToOpportunities converts recommendations to CloudshipAI opportunity format
func (t *FinOpsRecommendTool) transformToOpportunities(recommendations []interfaces.Recommendation) []interfaces.Opportunity {
	var opportunities []interfaces.Opportunity
	
	for _, rec := range recommendations {
		opportunity := interfaces.Opportunity{
			ID:                   rec.ID,
			ResourceID:           rec.ResourceID,
			ResourceARN:          rec.ResourceARN,
			Type:                 t.mapRecommendationTypeToOpportunityType(rec.RecommendationType),
			Vendor:               rec.Provider,
			Title:                rec.Title,
			Description:          rec.Description,
			EstimatedSavings:     rec.EstimatedSavings.Monthly,
			Currency:             rec.EstimatedSavings.Currency,
			Complexity:           t.mapEffortToComplexity(rec.MigrationEffort),
			RiskLevel:            rec.PerformanceRisk,
			ImplementationEffort: rec.MigrationEffort,
			CreatedAt:            rec.CreatedAt,
		}
		
		// Add recommendation details
		if len(rec.Options) > 0 {
			opportunity.RecommendationDetails = &interfaces.RecommendationDetails{
				CurrentConfiguration:     map[string]interface{}{},
				RecommendedConfiguration: map[string]interface{}{},
				ImplementationSteps:      []string{},
			}
			
			// Extract details from the first (best) option
			bestOption := rec.Options[0]
			if bestOption.Parameters != nil {
				// Copy parameters directly
				for key, value := range bestOption.Parameters {
					opportunity.RecommendationDetails.CurrentConfiguration[key] = value
					opportunity.RecommendationDetails.RecommendedConfiguration[key] = value
				}
			}
		}
		
		opportunities = append(opportunities, opportunity)
	}
	
	return opportunities
}

// mapRecommendationTypeToOpportunityType maps recommendation types to opportunity types
func (t *FinOpsRecommendTool) mapRecommendationTypeToOpportunityType(recType string) interfaces.OpportunityType {
	switch recType {
	case "RightSizeInstance", "UpgradeInstanceGeneration":
		return interfaces.OpportunityRightsizing
	case "TerminateInstance":
		return interfaces.OpportunityIdleResource
	case "ReservedInstance":
		return interfaces.OpportunityCommitment
	case "StorageOptimization":
		return interfaces.OpportunityStorageOptimization
	case "SpotInstance":
		return interfaces.OpportunitySpotInstances
	default:
		return interfaces.OpportunityRightsizing
	}
}

// mapEffortToComplexity maps migration effort to complexity scale
func (t *FinOpsRecommendTool) mapEffortToComplexity(effort string) string {
	switch effort {
	case "VeryLow", "Low":
		return "XS"
	case "Medium":
		return "M"
	case "High":
		return "L"
	case "VeryHigh":
		return "XL"
	default:
		return "S"
	}
}

// calculateTotalSavings calculates total monthly savings from opportunities
func (t *FinOpsRecommendTool) calculateTotalSavings(opportunities []interfaces.Opportunity) float64 {
	var total float64
	for _, opp := range opportunities {
		total += opp.EstimatedSavings
	}
	return total
}

// GetName returns the tool name
func (t *FinOpsRecommendTool) GetName() string {
	return "finops-recommend"
}

// GetDescription returns the tool description  
func (t *FinOpsRecommendTool) GetDescription() string {
	return "Generates cost optimization recommendations across providers using vendor-specific recommendation engines"
}

// GetInputSchema returns the JSON schema for the tool input
func (t *FinOpsRecommendTool) GetInputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"provider": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"aws", "gcp", "azure", "kubernetes"},
				"description": "Cloud provider to generate recommendations for",
			},
			"finding_types": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
					"enum": []string{"rightsizing", "reserved_instances", "spot_instances", "storage_optimization"},
				},
				"description": "Types of optimization opportunities to find",
				"default":     []string{"rightsizing"},
			},
			"regions": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional regions to analyze",
			},
			"arns": map[string]interface{}{
				"type": "array", 
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional resource ARNs to analyze (AWS only)",
			},
			"account_ids": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional account IDs for multi-account analysis",
			},
			"min_savings": map[string]interface{}{
				"type":        "number",
				"description": "Minimum monthly savings threshold in USD",
				"minimum":     0,
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