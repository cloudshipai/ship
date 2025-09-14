package mcp

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinOpsDiscoverTool_GetName(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsDiscoverTool(config)
	assert.NoError(t, err)
	assert.Equal(t, "finops-discover", tool.GetName())
}

func TestFinOpsDiscoverTool_GetDescription(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsDiscoverTool(config)
	assert.NoError(t, err)
	assert.NotEmpty(t, tool.GetDescription())
}

func TestFinOpsDiscoverTool_GetInputSchema(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsDiscoverTool(config)
	assert.NoError(t, err)
	
	schema := tool.GetInputSchema()
	assert.NotNil(t, schema)
	
	// Validate schema structure
	schemaMap, ok := schema.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "object", schemaMap["type"])
	
	properties, exists := schemaMap["properties"].(map[string]interface{})
	assert.True(t, exists)
	assert.Contains(t, properties, "provider")
}

func TestDiscoverRequest_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"provider": "aws",
		"region": "us-east-1",
		"resource_types": ["compute", "storage"],
		"tags": {"Environment": "production"}
	}`
	
	var req DiscoverRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	assert.NoError(t, err)
	assert.Equal(t, "aws", req.Provider)
	assert.Equal(t, "us-east-1", req.Region)
	assert.Equal(t, []string{"compute", "storage"}, req.ResourceTypes)
	assert.Equal(t, map[string]string{"Environment": "production"}, req.Tags)
}

func TestFinOpsRecommendTool_GetName(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	assert.NoError(t, err)
	assert.Equal(t, "finops-recommend", tool.GetName())
}

func TestFinOpsAnalyzeTool_GetName(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsAnalyzeTool(config)
	assert.NoError(t, err)
	assert.Equal(t, "finops-analyze", tool.GetName())
}

func TestRecommendRequest_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"provider": "aws",
		"finding_types": ["rightsizing", "reserved_instances"],
		"regions": ["us-east-1", "us-west-2"],
		"min_savings": 100.50
	}`
	
	var req RecommendRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	assert.NoError(t, err)
	assert.Equal(t, "aws", req.Provider)
	assert.Equal(t, []string{"rightsizing", "reserved_instances"}, req.FindingTypes)
	assert.Equal(t, []string{"us-east-1", "us-west-2"}, req.Regions)
	assert.Equal(t, 100.50, req.MinSavings)
}

func TestAnalyzeRequest_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"provider": "aws",
		"time_window": "30d",
		"granularity": "daily",
		"group_by": ["service", "region"],
		"currency": "USD"
	}`
	
	var req AnalyzeRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	assert.NoError(t, err)
	assert.Equal(t, "aws", req.Provider)
	assert.Equal(t, "30d", req.TimeWindow)
	assert.Equal(t, "daily", req.Granularity)
	assert.Equal(t, []string{"service", "region"}, req.GroupBy)
	assert.Equal(t, "USD", req.Currency)
}

// Integration tests demonstrating full tool execution with stub providers

func TestMapRecommendationTypeToOpportunityType(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	assert.NoError(t, err)
	
	tests := []struct {
		recType  string
		expected interfaces.OpportunityType
	}{
		{"RightSizeInstance", interfaces.OpportunityRightsizing},
		{"UpgradeInstanceGeneration", interfaces.OpportunityRightsizing},
		{"TerminateInstance", interfaces.OpportunityIdleResource},
		{"ReservedInstance", interfaces.OpportunityCommitment},
		{"StorageOptimization", interfaces.OpportunityStorageOptimization},
		{"SpotInstance", interfaces.OpportunitySpotInstances},
		{"Unknown", interfaces.OpportunityRightsizing}, // default case
	}
	
	for _, tt := range tests {
		result := tool.mapRecommendationTypeToOpportunityType(tt.recType)
		assert.Equal(t, tt.expected, result, "Failed for recommendation type: %s", tt.recType)
	}
}

func TestMapEffortToComplexity(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	assert.NoError(t, err)
	
	tests := []struct {
		effort   string
		expected string
	}{
		{"VeryLow", "XS"},
		{"Low", "XS"},
		{"Medium", "M"},
		{"High", "L"},
		{"VeryHigh", "XL"},
		{"Unknown", "S"}, // default case
	}
	
	for _, tt := range tests {
		result := tool.mapEffortToComplexity(tt.effort)
		assert.Equal(t, tt.expected, result, "Failed for effort: %s", tt.effort)
	}
}

// Integration tests demonstrating complete tool functionality

func TestFinOpsDiscoverTool_ExecuteIntegration(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsDiscoverTool(config)
	assert.NoError(t, err)
	
	// Test AWS discovery
	req := DiscoverRequest{
		Provider:      "aws",
		Region:        "us-east-1",
		ResourceTypes: []string{"compute", "storage"},
		Tags:          map[string]string{"Environment": "production"},
	}
	
	reqJSON, err := json.Marshal(req)
	assert.NoError(t, err)
	
	result, err := tool.Execute(context.Background(), reqJSON)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Validate response structure
	response, ok := result.(DiscoverResponse)
	assert.True(t, ok)
	assert.Equal(t, "aws", response.Provider)
	assert.Greater(t, response.TotalCount, 0)
	assert.NotEmpty(t, response.Resources)
	t.Logf("Discovered %d resources", response.TotalCount)
}

func TestFinOpsRecommendTool_ExecuteIntegration(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	assert.NoError(t, err)
	
	// Test AWS recommendations
	req := RecommendRequest{
		Provider:     "aws",
		FindingTypes: []string{"rightsizing"},
		Regions:      []string{"us-east-1"},
		MinSavings:   0,
	}
	
	reqJSON, err := json.Marshal(req)
	assert.NoError(t, err)
	
	result, err := tool.Execute(context.Background(), reqJSON)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Validate response structure
	response, ok := result.(RecommendResponse)
	assert.True(t, ok)
	assert.Equal(t, "aws", response.Provider)
	assert.Equal(t, []string{"rightsizing"}, response.FindingTypes)
	assert.GreaterOrEqual(t, response.TotalCount, 0)
	assert.GreaterOrEqual(t, response.TotalSavings, float64(0))
	t.Logf("Found %d recommendations with potential savings of $%.2f/month", response.TotalCount, response.TotalSavings)
}

func TestFinOpsAnalyzeTool_ExecuteIntegration(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsAnalyzeTool(config)
	assert.NoError(t, err)
	
	// Test AWS cost analysis
	req := AnalyzeRequest{
		Provider:    "aws",
		TimeWindow:  "30d",
		Granularity: "daily",
		GroupBy:     []string{"service", "region"},
		Currency:    "USD",
	}
	
	reqJSON, err := json.Marshal(req)
	assert.NoError(t, err)
	
	result, err := tool.Execute(context.Background(), reqJSON)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Validate response structure
	response, ok := result.(AnalyzeResponse)
	assert.True(t, ok)
	assert.Equal(t, "aws", response.Provider)
	assert.Equal(t, "USD", response.Summary.Currency)
	assert.NotEmpty(t, response.CostData)
	t.Logf("Analyzed cost data with %d records totaling $%.2f", len(response.CostData), response.Summary.TotalCost)
}

func TestFinOpsQueryTool_ExecuteIntegration(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsQueryTool(config)
	assert.NoError(t, err)
	
	// Test agent-driven query
	req := QueryRequest{
		Query:    "find overprovisioned EC2 instances in production",
		Provider: "aws",
		Filters: QueryFilters{
			Regions:       []string{"us-east-1"},
			ResourceTypes: []string{"compute"},
			MinSavings:    25.0,
		},
	}
	
	reqJSON, err := json.Marshal(req)
	assert.NoError(t, err)
	
	result, err := tool.Execute(context.Background(), reqJSON)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Validate response structure
	response, ok := result.(QueryResponse)
	assert.True(t, ok)
	assert.NotEmpty(t, response.Query)
	assert.NotNil(t, response.Results)
	t.Logf("Query executed successfully: %s", response.Summary)
}

// Test CloudshipAI flag behavior

func TestFinOpsDiscoverTool_CloudshipAIFlagDisabled(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsDiscoverTool(config)
	require.NoError(t, err)

	// Test without EnableCloudshipAI flag
	req := DiscoverRequest{
		Provider:          "aws",
		Region:            "us-east-1",
		ResourceTypes:     []string{"compute"},
		EnableCloudshipAI: false, // Explicitly disabled
	}

	reqJSON, err := json.Marshal(req)
	require.NoError(t, err)

	result, err := tool.Execute(context.Background(), reqJSON)
	require.NoError(t, err)

	response, ok := result.(DiscoverResponse)
	require.True(t, ok)

	// Should indicate lighthouse reporting was disabled
	assert.False(t, response.Lighthouse.Reported)
	assert.True(t, response.Lighthouse.Disabled)
	assert.Empty(t, response.Lighthouse.Error)
	t.Logf("✅ CloudshipAI reporting disabled as expected")
}

func TestFinOpsDiscoverTool_CloudshipAIFlagEnabled(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsDiscoverTool(config)
	require.NoError(t, err)

	// Test with EnableCloudshipAI flag
	req := DiscoverRequest{
		Provider:          "aws",
		Region:            "us-east-1",
		ResourceTypes:     []string{"compute"},
		EnableCloudshipAI: true, // Explicitly enabled
	}

	reqJSON, err := json.Marshal(req)
	require.NoError(t, err)

	result, err := tool.Execute(context.Background(), reqJSON)
	require.NoError(t, err)

	response, ok := result.(DiscoverResponse)
	require.True(t, ok)

	// Should attempt lighthouse reporting (will fail due to no real endpoint)
	assert.False(t, response.Lighthouse.Reported) // Will fail in test
	assert.False(t, response.Lighthouse.Disabled)
	assert.NotEmpty(t, response.Lighthouse.Error) // Should have error from failed API call
	t.Logf("✅ CloudshipAI reporting attempted with error: %s", response.Lighthouse.Error)
}

func TestFinOpsRecommendTool_CloudshipAIFlag(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		enabled  bool
		expected string
	}{
		{"disabled", false, "should be disabled"},
		{"enabled", true, "should attempt reporting"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := RecommendRequest{
				Provider:          "aws",
				FindingTypes:      []string{"rightsizing"},
				EnableCloudshipAI: tc.enabled,
			}

			reqJSON, err := json.Marshal(req)
			require.NoError(t, err)

			result, err := tool.Execute(context.Background(), reqJSON)
			require.NoError(t, err)

			response, ok := result.(RecommendResponse)
			require.True(t, ok)

			if tc.enabled {
				// Should attempt reporting
				assert.False(t, response.Lighthouse.Disabled)
			} else {
				// Should be disabled
				assert.True(t, response.Lighthouse.Disabled)
				assert.False(t, response.Lighthouse.Reported)
			}

			t.Logf("✅ CloudshipAI flag test '%s': %s", tc.name, tc.expected)
		})
	}
}