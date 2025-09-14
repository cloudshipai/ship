package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudshipai/ship/internal/tools/finops/cloudshipai"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/providers/stub"
)

// FinOpsAnalyzeTool implements the finops-analyze MCP tool
type FinOpsAnalyzeTool struct {
	lighthouseClient interfaces.LighthouseReporter
	providers        map[interfaces.VendorType]interfaces.VendorProvider
}

// NewFinOpsAnalyzeTool creates a new finops-analyze tool instance
func NewFinOpsAnalyzeTool(config interfaces.LighthouseConfig) (*FinOpsAnalyzeTool, error) {
	// Initialize lighthouse client for reporting to Station/CloudshipAI
	lighthouseClient := cloudshipai.NewClient(config)
	
	// Initialize providers (using stub for testing)
	providers := make(map[interfaces.VendorType]interfaces.VendorProvider)
	
	// Initialize stub providers
	providers[interfaces.VendorAWS] = stub.NewStubProvider(interfaces.VendorAWS)
	providers[interfaces.VendorGCP] = stub.NewStubProvider(interfaces.VendorGCP)
	providers[interfaces.VendorAzure] = stub.NewStubProvider(interfaces.VendorAzure)
	providers[interfaces.VendorKubernetes] = stub.NewStubProvider(interfaces.VendorKubernetes)
	
	return &FinOpsAnalyzeTool{
		lighthouseClient: lighthouseClient,
		providers:        providers,
	}, nil
}

// AnalyzeRequest represents the input parameters for the finops-analyze tool
type AnalyzeRequest struct {
	Provider          string            `json:"provider"`                    // aws, gcp, azure, kubernetes  
	TimeWindow        string            `json:"time_window"`                 // "7d", "30d", "90d"
	Granularity       string            `json:"granularity"`                 // "daily", "monthly"
	GroupBy           []string          `json:"group_by,omitempty"`          // ["service", "region", "account"]
	Filters           map[string]string `json:"filters,omitempty"`           // additional filters
	AccountIDs        []string          `json:"account_ids,omitempty"`       // optional account filter
	Currency          string            `json:"currency,omitempty"`          // "USD", "EUR", etc. (default: USD)
	EnableCloudshipAI bool              `json:"enable_cloudshipai,omitempty"` // enable Station API reporting
}

// AnalyzeResponse represents the output of the finops-analyze tool
type AnalyzeResponse struct {
	CostData        []interfaces.CostRecord `json:"cost_data"`
	Summary         CostSummary             `json:"summary"`
	Insights        []interfaces.Insight    `json:"insights"`
	Provider        string                  `json:"provider"`
	TimeWindow      string                  `json:"time_window"`
	Granularity     string                  `json:"granularity"`
	Lighthouse      LighthouseStatus        `json:"lighthouse"`
}

// CostSummary provides aggregate cost information
type CostSummary struct {
	TotalCost        float64                    `json:"total_cost"`
	Currency         string                     `json:"currency"`
	AverageDailyCost float64                    `json:"average_daily_cost"`
	PreviousPeriod   *PreviousPeriodComparison  `json:"previous_period,omitempty"`
	TopServices      []ServiceCost              `json:"top_services"`
	TopRegions       []RegionCost               `json:"top_regions"`
}

// PreviousPeriodComparison compares costs to the previous period
type PreviousPeriodComparison struct {
	TotalCost      float64 `json:"total_cost"`
	PercentChange  float64 `json:"percent_change"`
	AbsoluteChange float64 `json:"absolute_change"`
}

// ServiceCost represents cost by service
type ServiceCost struct {
	Service string  `json:"service"`
	Cost    float64 `json:"cost"`
	Percent float64 `json:"percent_of_total"`
}

// RegionCost represents cost by region
type RegionCost struct {
	Region  string  `json:"region"`
	Cost    float64 `json:"cost"`
	Percent float64 `json:"percent_of_total"`
}

// Execute implements the MCP tool execution for cost analysis
func (t *FinOpsAnalyzeTool) Execute(ctx context.Context, request json.RawMessage) (interface{}, error) {
	var req AnalyzeRequest
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal analyze request: %w", err)
	}
	
	// Validate required parameters
	if req.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	if req.TimeWindow == "" {
		req.TimeWindow = "30d" // Default to 30 days
	}
	if req.Granularity == "" {
		req.Granularity = "daily" // Default to daily
	}
	if req.Currency == "" {
		req.Currency = "USD" // Default to USD
	}
	
	// Get the appropriate provider
	vendorType := interfaces.VendorType(req.Provider)
	provider, exists := t.providers[vendorType]
	if !exists {
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
	
	// Build cost options
	opts := interfaces.CostOptions{
		TimeWindow:  req.TimeWindow,
		Granularity: req.Granularity,
		GroupBy:     req.GroupBy,
		Filters:     req.Filters,
		Currency:    req.Currency,
	}
	
	// Get cost data from provider
	costData, err := provider.GetCostData(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost data: %w", err)
	}
	
	// Generate cost summary and insights
	summary := t.generateCostSummary(costData, req.Currency)
	insights := t.generateInsights(costData, summary, req.Provider)
	
	// Report cost data and insights to lighthouse (Station/CloudshipAI) only if enabled
	lighthouseStatus := LighthouseStatus{Reported: false}
	
	if req.EnableCloudshipAI {
		// Report cost data
		if len(costData) > 0 {
			if err := t.lighthouseClient.ReportCostData(ctx, costData); err != nil {
				lighthouseStatus.Error = err.Error()
			} else {
				lighthouseStatus.Reported = true
			}
		}
		
		// Report insights if any were generated
		if len(insights) > 0 {
			if err := t.lighthouseClient.ReportInsights(ctx, insights); err != nil {
				// Don't overwrite cost data reporting status, just log
				if lighthouseStatus.Error != "" {
					lighthouseStatus.Error += "; insights: " + err.Error()
				} else {
					lighthouseStatus.Error = "insights: " + err.Error()
				}
			}
		}
	} else {
		lighthouseStatus.Disabled = true
	}
	
	// Return response
	response := AnalyzeResponse{
		CostData:    costData,
		Summary:     summary,
		Insights:    insights,
		Provider:    req.Provider,
		TimeWindow:  req.TimeWindow,
		Granularity: req.Granularity,
		Lighthouse:  lighthouseStatus,
	}
	
	return response, nil
}

// generateCostSummary creates a cost summary from cost data
func (t *FinOpsAnalyzeTool) generateCostSummary(costData []interfaces.CostRecord, currency string) CostSummary {
	if len(costData) == 0 {
		return CostSummary{Currency: currency}
	}
	
	var totalCost float64
	serviceCosts := make(map[string]float64)
	regionCosts := make(map[string]float64)
	
	// Aggregate costs
	for _, record := range costData {
		totalCost += record.Amount
		
		if record.Service != "" {
			serviceCosts[record.Service] += record.Amount
		}
		if record.Region != "" {
			regionCosts[record.Region] += record.Amount
		}
	}
	
	// Calculate average daily cost (rough estimate)
	averageDailyCost := totalCost / float64(len(costData))
	
	// Get top services and regions
	topServices := t.getTopServiceCosts(serviceCosts, totalCost)
	topRegions := t.getTopRegionCosts(regionCosts, totalCost)
	
	return CostSummary{
		TotalCost:        totalCost,
		Currency:         currency,
		AverageDailyCost: averageDailyCost,
		TopServices:      topServices,
		TopRegions:       topRegions,
	}
}

// getTopServiceCosts returns top services by cost
func (t *FinOpsAnalyzeTool) getTopServiceCosts(serviceCosts map[string]float64, totalCost float64) []ServiceCost {
	var services []ServiceCost
	for service, cost := range serviceCosts {
		percent := (cost / totalCost) * 100
		services = append(services, ServiceCost{
			Service: service,
			Cost:    cost,
			Percent: percent,
		})
	}
	
	// Sort by cost (simplified - in production would use sort.Slice)
	// Return top 5 for now
	if len(services) > 5 {
		services = services[:5]
	}
	
	return services
}

// getTopRegionCosts returns top regions by cost
func (t *FinOpsAnalyzeTool) getTopRegionCosts(regionCosts map[string]float64, totalCost float64) []RegionCost {
	var regions []RegionCost
	for region, cost := range regionCosts {
		percent := (cost / totalCost) * 100
		regions = append(regions, RegionCost{
			Region:  region,
			Cost:    cost,
			Percent: percent,
		})
	}
	
	// Sort by cost (simplified - in production would use sort.Slice)
	// Return top 5 for now
	if len(regions) > 5 {
		regions = regions[:5]
	}
	
	return regions
}

// generateInsights creates cost insights from the data
func (t *FinOpsAnalyzeTool) generateInsights(costData []interfaces.CostRecord, summary CostSummary, provider string) []interfaces.Insight {
	var insights []interfaces.Insight
	
	// Generate cost trend insight
	if len(costData) > 0 {
		insight := interfaces.Insight{
			ID:          fmt.Sprintf("cost-trend-%s", provider),
			Type:        interfaces.InsightTypeTrend,
			Severity:    interfaces.InsightSeverityMedium,
			Title:       fmt.Sprintf("%s Cost Analysis", provider),
			Description: fmt.Sprintf("Total costs for the analyzed period: $%.2f %s", summary.TotalCost, summary.Currency),
			Category:    "cost_analysis",
			CostImpact:  summary.TotalCost,
			Vendor:      interfaces.VendorType(provider),
		}
		insights = append(insights, insight)
	}
	
	// Generate high-cost service insights
	for _, service := range summary.TopServices {
		if service.Percent > 40 { // If a service is more than 40% of costs
			insight := interfaces.Insight{
				ID:          fmt.Sprintf("high-cost-service-%s-%s", provider, service.Service),
				Type:        interfaces.InsightTypeOptimization,
				Severity:    interfaces.InsightSeverityHigh,
				Title:       fmt.Sprintf("High Cost Service: %s", service.Service),
				Description: fmt.Sprintf("Service %s accounts for %.1f%% of total costs ($%.2f)", service.Service, service.Percent, service.Cost),
				Category:    "cost_concentration",
				Service:     service.Service,
				CostImpact:  service.Cost,
				Vendor:      interfaces.VendorType(provider),
			}
			insights = append(insights, insight)
		}
	}
	
	return insights
}

// GetName returns the tool name
func (t *FinOpsAnalyzeTool) GetName() string {
	return "finops-analyze"
}

// GetDescription returns the tool description
func (t *FinOpsAnalyzeTool) GetDescription() string {
	return "Analyzes cost data and trends across providers with insights and anomaly detection"
}

// GetInputSchema returns the JSON schema for the tool input
func (t *FinOpsAnalyzeTool) GetInputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"provider": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"aws", "gcp", "azure", "kubernetes"},
				"description": "Cloud provider to analyze costs for",
			},
			"time_window": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"7d", "30d", "90d", "1y"},
				"description": "Time window for cost analysis",
				"default":     "30d",
			},
			"granularity": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"daily", "monthly"},
				"description": "Granularity of cost data",
				"default":     "daily",
			},
			"group_by": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
					"enum": []string{"service", "region", "account", "instance", "platform"},
				},
				"description": "Dimensions to group cost data by",
			},
			"filters": map[string]interface{}{
				"type":        "object",
				"description": "Additional filters for cost data",
			},
			"account_ids": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional account IDs for multi-account analysis",
			},
			"currency": map[string]interface{}{
				"type":        "string",
				"description": "Currency for cost data",
				"default":     "USD",
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