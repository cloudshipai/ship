package mcp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/stretchr/testify/assert"
)

// Simple tests that don't require AWS SDK dependencies
func TestToolNames(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	// Test tool creation and basic properties
	discoverTool, err := NewFinOpsDiscoverTool(config)
	if err == nil { // Only test if creation succeeds (may fail without AWS creds)
		assert.Equal(t, "finops-discover", discoverTool.GetName())
		assert.NotEmpty(t, discoverTool.GetDescription())
		assert.NotNil(t, discoverTool.GetInputSchema())
	}
	
	recommendTool, err := NewFinOpsRecommendTool(config)
	if err == nil {
		assert.Equal(t, "finops-recommend", recommendTool.GetName())
		assert.NotEmpty(t, recommendTool.GetDescription())
		assert.NotNil(t, recommendTool.GetInputSchema())
	}
	
	analyzeTool, err := NewFinOpsAnalyzeTool(config)
	if err == nil {
		assert.Equal(t, "finops-analyze", analyzeTool.GetName())
		assert.NotEmpty(t, analyzeTool.GetDescription())
		assert.NotNil(t, analyzeTool.GetInputSchema())
	}
	
	queryTool, err := NewFinOpsQueryTool(config)
	if err == nil {
		assert.Equal(t, "finops-query", queryTool.GetName())
		assert.NotEmpty(t, queryTool.GetDescription())
		assert.NotNil(t, queryTool.GetInputSchema())
	}
}

func TestRequestSerialization(t *testing.T) {
	// Test that our request structures serialize/deserialize correctly
	
	discoverReq := DiscoverRequest{
		Provider:      "aws",
		Region:        "us-east-1",
		ResourceTypes: []string{"compute", "storage"},
		Tags:          map[string]string{"Environment": "production"},
		AccountIDs:    []string{"123456789012"},
		ARNs:          []string{"arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0"},
	}
	
	data, err := json.Marshal(discoverReq)
	assert.NoError(t, err)
	
	var decoded DiscoverRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, discoverReq.Provider, decoded.Provider)
	assert.Equal(t, discoverReq.Region, decoded.Region)
	assert.Equal(t, discoverReq.ResourceTypes, decoded.ResourceTypes)
	assert.Equal(t, discoverReq.Tags, decoded.Tags)
}

func TestRecommendRequestSerialization(t *testing.T) {
	recommendReq := RecommendRequest{
		Provider:     "aws",
		FindingTypes: []string{"rightsizing", "reserved_instances"},
		Regions:      []string{"us-east-1", "us-west-2"},
		MinSavings:   100.50,
		AccountIDs:   []string{"123456789012"},
	}
	
	data, err := json.Marshal(recommendReq)
	assert.NoError(t, err)
	
	var decoded RecommendRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, recommendReq.Provider, decoded.Provider)
	assert.Equal(t, recommendReq.FindingTypes, decoded.FindingTypes)
	assert.Equal(t, recommendReq.MinSavings, decoded.MinSavings)
}

func TestAnalyzeRequestSerialization(t *testing.T) {
	analyzeReq := AnalyzeRequest{
		Provider:    "aws",
		TimeWindow:  "30d",
		Granularity: "daily",
		GroupBy:     []string{"service", "region"},
		Currency:    "USD",
	}
	
	data, err := json.Marshal(analyzeReq)
	assert.NoError(t, err)
	
	var decoded AnalyzeRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, analyzeReq.Provider, decoded.Provider)
	assert.Equal(t, analyzeReq.TimeWindow, decoded.TimeWindow)
	assert.Equal(t, analyzeReq.Granularity, decoded.Granularity)
}

func TestQueryRequestSerialization(t *testing.T) {
	queryReq := QueryRequest{
		Query:    "find overprovisioned EC2 instances",
		Provider: "aws",
		Context:  map[string]interface{}{"region": "us-east-1"},
		Filters: QueryFilters{
			Regions:       []string{"us-east-1"},
			ResourceTypes: []string{"compute"},
			MinSavings:    50.0,
		},
	}
	
	data, err := json.Marshal(queryReq)
	assert.NoError(t, err)
	
	var decoded QueryRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, queryReq.Query, decoded.Query)
	assert.Equal(t, queryReq.Provider, decoded.Provider)
	assert.Equal(t, queryReq.Filters.MinSavings, decoded.Filters.MinSavings)
}

func TestInputSchemas(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	// Test that input schemas are valid JSON structures
	tools := []struct {
		name string
		tool interface {
			GetInputSchema() interface{}
		}
	}{}
	
	if discoverTool, err := NewFinOpsDiscoverTool(config); err == nil {
		tools = append(tools, struct {
			name string
			tool interface{ GetInputSchema() interface{} }
		}{"discover", discoverTool})
	}
	
	if recommendTool, err := NewFinOpsRecommendTool(config); err == nil {
		tools = append(tools, struct {
			name string
			tool interface{ GetInputSchema() interface{} }
		}{"recommend", recommendTool})
	}
	
	if analyzeTool, err := NewFinOpsAnalyzeTool(config); err == nil {
		tools = append(tools, struct {
			name string
			tool interface{ GetInputSchema() interface{} }
		}{"analyze", analyzeTool})
	}
	
	if queryTool, err := NewFinOpsQueryTool(config); err == nil {
		tools = append(tools, struct {
			name string
			tool interface{ GetInputSchema() interface{} }
		}{"query", queryTool})
	}
	
	for _, tool := range tools {
		t.Run(tool.name, func(t *testing.T) {
			schema := tool.tool.GetInputSchema()
			assert.NotNil(t, schema)
			
			// Verify it can be marshaled to JSON
			_, err := json.Marshal(schema)
			assert.NoError(t, err, "Schema should be JSON serializable")
			
			// Verify it has required structure
			if schemaMap, ok := schema.(map[string]interface{}); ok {
				assert.Equal(t, "object", schemaMap["type"])
				assert.Contains(t, schemaMap, "properties")
			}
		})
	}
}

func TestOpportunityTypeMapping(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	if err != nil {
		t.Skip("Skipping test due to AWS SDK dependency")
		return
	}
	
	// Test recommendation type to opportunity type mapping
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

func TestEffortToComplexityMapping(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}
	
	tool, err := NewFinOpsRecommendTool(config)
	if err != nil {
		t.Skip("Skipping test due to AWS SDK dependency")
		return
	}
	
	// Test migration effort to complexity mapping
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