package cloudshipai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
)

// Client implements the LighthouseReporter interface for CloudshipAI integration
type Client struct {
	config     interfaces.LighthouseConfig
	httpClient *http.Client
}

// NewClient creates a new CloudshipAI lighthouse client
func NewClient(config interfaces.LighthouseConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = 3
	}
	
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// ReportResources sends resource inventory data to CloudshipAI via Station proxy
func (c *Client) ReportResources(ctx context.Context, resources []interfaces.Resource) error {
	request := interfaces.ResourceInventoryRequest{
		ReportRequest: interfaces.ReportRequest{
			DataType:  "resource_inventory",
			Timestamp: time.Now().UTC(),
			Source:    "ship-finops-discover",
			Version:   "1.0",
		},
		Resources: resources,
	}
	
	return c.sendReport(ctx, request)
}

// ReportOpportunities sends optimization opportunities to CloudshipAI
func (c *Client) ReportOpportunities(ctx context.Context, opportunities []interfaces.Opportunity) error {
	// Validate opportunities before sending
	for i, opp := range opportunities {
		if err := c.validateOpportunity(opp); err != nil {
			return fmt.Errorf("opportunity %d validation failed: %w", i, err)
		}
	}
	
	request := interfaces.OpportunityRequest{
		ReportRequest: interfaces.ReportRequest{
			DataType:  "opportunities",
			Timestamp: time.Now().UTC(),
			Source:    "ship-finops-recommend",
			Version:   "1.0",
		},
		Opportunities: opportunities,
	}
	
	return c.sendReport(ctx, request)
}

// ReportCostData sends cost analysis data to CloudshipAI
func (c *Client) ReportCostData(ctx context.Context, costData []interfaces.CostRecord) error {
	request := interfaces.CostDataRequest{
		ReportRequest: interfaces.ReportRequest{
			DataType:  "cost_data",
			Timestamp: time.Now().UTC(),
			Source:    "ship-finops-analyze",
			Version:   "1.0",
		},
		CostRecords: costData,
	}
	
	return c.sendReport(ctx, request)
}

// ReportInsights sends insights and anomalies to CloudshipAI
func (c *Client) ReportInsights(ctx context.Context, insights []interfaces.Insight) error {
	request := interfaces.InsightRequest{
		ReportRequest: interfaces.ReportRequest{
			DataType:  "insights",
			Timestamp: time.Now().UTC(),
			Source:    "ship-finops-insights",
			Version:   "1.0",
		},
		Insights: insights,
	}
	
	return c.sendReport(ctx, request)
}

// HealthCheck verifies connectivity to CloudshipAI lighthouse endpoint
func (c *Client) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.getEndpointURL()+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}
	
	// No authorization needed - Station is on localhost
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// sendReport is the generic method for sending data to CloudshipAI
func (c *Client) sendReport(ctx context.Context, request interface{}) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	url := c.getEndpointURL() + "/finops"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	if c.config.EnableTracing {
		req.Header.Set("X-Trace-ID", generateTraceID())
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("lighthouse API returned status %d", resp.StatusCode)
	}
	
	// Parse response for any warnings or errors
	var response interfaces.ReportResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// Non-critical error - log but don't fail
		if c.config.EnableTracing {
			fmt.Printf("Warning: failed to parse response: %v\n", err)
		}
		return nil
	}
	
	if !response.Success {
		return fmt.Errorf("lighthouse reported failure: %v", response.Errors)
	}
	
	return nil
}

// validateOpportunity validates an opportunity against CloudshipAI requirements
func (c *Client) validateOpportunity(opp interfaces.Opportunity) error {
	if opp.ResourceID == "" {
		return fmt.Errorf("resource_id is required")
	}
	if opp.Type == "" {
		return fmt.Errorf("opportunity_type is required")
	}
	if opp.Vendor == "" {
		return fmt.Errorf("vendor is required")
	}
	if opp.EstimatedSavings < 0 {
		return fmt.Errorf("estimated_savings_monthly_usd must be >= 0")
	}
	
	// Validate enum values
	validTypes := map[interfaces.OpportunityType]bool{
		interfaces.OpportunityRightsizing:         true,
		interfaces.OpportunityIdleResource:        true,
		interfaces.OpportunityCommitment:          true,
		interfaces.OpportunityStorageOptimization: true,
		interfaces.OpportunitySpotInstances:       true,
	}
	
	if !validTypes[opp.Type] {
		return fmt.Errorf("invalid opportunity type: %s", opp.Type)
	}
	
	validVendors := map[interfaces.VendorType]bool{
		interfaces.VendorAWS:        true,
		interfaces.VendorGCP:        true,
		interfaces.VendorAzure:      true,
		interfaces.VendorKubernetes: true,
	}
	
	if !validVendors[opp.Vendor] {
		return fmt.Errorf("invalid vendor: %s", opp.Vendor)
	}
	
	return nil
}

// getEndpointURL constructs the full endpoint URL (Station runs on localhost)
func (c *Client) getEndpointURL() string {
	return "http://localhost:8585/api/v1"
}

// generateTraceID creates a unique trace ID for request tracking
func generateTraceID() string {
	return fmt.Sprintf("ship-finops-%d", time.Now().UnixNano())
}

// ReporterOptions configures the reporter behavior
type ReporterOptions struct {
	BatchSize    int
	EnableRetry  bool
	ValidateData bool
}

// BatchReporter provides batch reporting capabilities with retry logic
type BatchReporter struct {
	client  *Client
	options ReporterOptions
}

// NewBatchReporter creates a new batch reporter
func NewBatchReporter(client *Client, options ReporterOptions) *BatchReporter {
	if options.BatchSize == 0 {
		options.BatchSize = 100
	}
	
	return &BatchReporter{
		client:  client,
		options: options,
	}
}

// ReportOpportunitiesBatch reports opportunities in batches with retry logic
func (br *BatchReporter) ReportOpportunitiesBatch(ctx context.Context, opportunities []interfaces.Opportunity) error {
	if len(opportunities) == 0 {
		return nil
	}
	
	// Process in batches
	for i := 0; i < len(opportunities); i += br.options.BatchSize {
		end := i + br.options.BatchSize
		if end > len(opportunities) {
			end = len(opportunities)
		}
		
		batch := opportunities[i:end]
		
		// Retry logic
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			if err := br.client.ReportOpportunities(ctx, batch); err != nil {
				lastErr = err
				// Exponential backoff
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			break // Success
		}
		
		if lastErr != nil {
			return fmt.Errorf("failed to send batch after retries: %w", lastErr)
		}
	}
	
	return nil
}