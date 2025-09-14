package finops

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFinOpsWorkflowIntegration demonstrates the complete FinOps workflow
// This replaces the need for a demo binary by showing all functionality through tests
func TestFinOpsWorkflowIntegration(t *testing.T) {
	t.Log("🚀 FinOps Tools Integration Test")
	t.Log("==========================================")

	// Configuration
	config := interfaces.LighthouseConfig{
		Timeout:       30 * time.Second,
		EnableTracing: true,
	}

	ctx := context.Background()

	// Test 1: Resource Discovery
	t.Log("\n📍 Test 1: Resource Discovery")
	t.Log("------------------------------")
	
	discoverTool, err := mcp.NewFinOpsDiscoverTool(config)
	require.NoError(t, err)

	discoverReq := mcp.DiscoverRequest{
		Provider:      "aws",
		Region:        "us-east-1",
		ResourceTypes: []string{"compute", "storage"},
		Tags:          map[string]string{"Environment": "production"},
	}

	reqJSON, err := json.Marshal(discoverReq)
	require.NoError(t, err)

	result, err := discoverTool.Execute(ctx, reqJSON)
	require.NoError(t, err)

	discoverResponse, ok := result.(mcp.DiscoverResponse)
	require.True(t, ok)
	
	assert.Equal(t, "aws", discoverResponse.Provider)
	assert.Greater(t, discoverResponse.TotalCount, 0)
	assert.NotEmpty(t, discoverResponse.Resources)
	t.Logf("✅ Discovered %d resources", discoverResponse.TotalCount)

	// Test 2: Cost Optimization Recommendations
	t.Log("\n💡 Test 2: Cost Optimization Recommendations")
	t.Log("---------------------------------------------")
	
	recommendTool, err := mcp.NewFinOpsRecommendTool(config)
	require.NoError(t, err)

	recommendReq := mcp.RecommendRequest{
		Provider:     "aws",
		FindingTypes: []string{"rightsizing"},
		Regions:      []string{"us-east-1"},
		MinSavings:   0,
	}

	reqJSON, err = json.Marshal(recommendReq)
	require.NoError(t, err)

	result, err = recommendTool.Execute(ctx, reqJSON)
	require.NoError(t, err)

	recommendResponse, ok := result.(mcp.RecommendResponse)
	require.True(t, ok)
	
	assert.Equal(t, "aws", recommendResponse.Provider)
	assert.Equal(t, []string{"rightsizing"}, recommendResponse.FindingTypes)
	assert.GreaterOrEqual(t, recommendResponse.TotalCount, 0)
	assert.GreaterOrEqual(t, recommendResponse.TotalSavings, float64(0))
	t.Logf("✅ Found %d recommendations with potential savings of $%.2f/month", 
		recommendResponse.TotalCount, recommendResponse.TotalSavings)

	// Test 3: Cost Analysis
	t.Log("\n📊 Test 3: Cost Analysis")
	t.Log("-------------------------")
	
	analyzeTool, err := mcp.NewFinOpsAnalyzeTool(config)
	require.NoError(t, err)

	analyzeReq := mcp.AnalyzeRequest{
		Provider:    "aws",
		TimeWindow:  "30d",
		Granularity: "daily",
		GroupBy:     []string{"service", "region"},
		Currency:    "USD",
	}

	reqJSON, err = json.Marshal(analyzeReq)
	require.NoError(t, err)

	result, err = analyzeTool.Execute(ctx, reqJSON)
	require.NoError(t, err)

	analyzeResponse, ok := result.(mcp.AnalyzeResponse)
	require.True(t, ok)
	
	assert.Equal(t, "aws", analyzeResponse.Provider)
	assert.Equal(t, "USD", analyzeResponse.Summary.Currency)
	assert.NotEmpty(t, analyzeResponse.CostData)
	t.Logf("✅ Analyzed cost data with %d records totaling $%.2f", 
		len(analyzeResponse.CostData), analyzeResponse.Summary.TotalCost)

	// Test 4: Agent-Driven Query
	t.Log("\n🤖 Test 4: Agent-Driven Query")
	t.Log("------------------------------")
	
	queryTool, err := mcp.NewFinOpsQueryTool(config)
	require.NoError(t, err)

	queryReq := mcp.QueryRequest{
		Query:    "find overprovisioned EC2 instances in production",
		Provider: "aws",
		Filters: mcp.QueryFilters{
			Regions:       []string{"us-east-1"},
			ResourceTypes: []string{"compute"},
			MinSavings:    25.0,
		},
	}

	reqJSON, err = json.Marshal(queryReq)
	require.NoError(t, err)

	result, err = queryTool.Execute(ctx, reqJSON)
	require.NoError(t, err)

	queryResponse, ok := result.(mcp.QueryResponse)
	require.True(t, ok)
	
	assert.NotEmpty(t, queryResponse.Query)
	assert.NotNil(t, queryResponse.Results)
	t.Logf("✅ Query executed successfully: %s", queryResponse.Summary)

	t.Log("\n✅ Integration test completed successfully!")
	
	// Show usage examples
	t.Log("\n📚 CLI Usage Examples:")
	t.Log("  ship finops discover --provider=aws --region=us-east-1")
	t.Log("  ship finops recommend --provider=aws --finding-types=rightsizing --cloudshipai")
	t.Log("  ship finops analyze --provider=aws --time-window=30d --cloudshipai")
	t.Log("  ship finops query --query='find expensive resources' --provider=aws")
}

// TestAllProvidersSupported verifies all provider types are properly supported
func TestAllProvidersSupported(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}

	providers := []string{"aws", "gcp", "azure", "kubernetes"}
	
	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			// Test discover tool
			discoverTool, err := mcp.NewFinOpsDiscoverTool(config)
			require.NoError(t, err)

			req := mcp.DiscoverRequest{
				Provider:      provider,
				Region:        "us-east-1",
				ResourceTypes: []string{"compute"},
			}

			reqJSON, err := json.Marshal(req)
			require.NoError(t, err)

			result, err := discoverTool.Execute(context.Background(), reqJSON)
			require.NoError(t, err)
			
			response, ok := result.(mcp.DiscoverResponse)
			require.True(t, ok)
			assert.Equal(t, provider, response.Provider)
			
			t.Logf("✅ Provider %s is properly supported", provider)
		})
	}
}

// TestLighthouseReporting verifies CloudshipAI lighthouse integration
func TestLighthouseReporting(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}

	// Test that lighthouse status is properly tracked
	recommendTool, err := mcp.NewFinOpsRecommendTool(config)
	require.NoError(t, err)

	req := mcp.RecommendRequest{
		Provider:     "aws",
		FindingTypes: []string{"rightsizing"},
		Regions:      []string{"us-east-1"},
	}

	reqJSON, err := json.Marshal(req)
	require.NoError(t, err)

	result, err := recommendTool.Execute(context.Background(), reqJSON)
	require.NoError(t, err)

	response, ok := result.(mcp.RecommendResponse)
	require.True(t, ok)
	
	// Lighthouse status should be tracked
	assert.NotEmpty(t, response.Lighthouse)
	
	// If there are opportunities, lighthouse reporting should have been attempted
	if len(response.Opportunities) > 0 {
		// In our test setup, this will fail due to no real endpoint, 
		// but the attempt should be tracked
		t.Logf("Lighthouse reporting attempted: reported=%v, error=%s", 
			response.Lighthouse.Reported, response.Lighthouse.Error)
	}
}

// TestSchemaValidation ensures all input/output schemas are properly structured
func TestSchemaValidation(t *testing.T) {
	config := interfaces.LighthouseConfig{
		Timeout: 30 * time.Second,
	}

	tools := []struct {
		name string
		tool interface{
			GetName() string
			GetDescription() string
			GetInputSchema() interface{}
		}
	}{
		{"discover", mustCreateDiscoverTool(t, config)},
		{"recommend", mustCreateRecommendTool(t, config)},
		{"analyze", mustCreateAnalyzeTool(t, config)},
		{"query", mustCreateQueryTool(t, config)},
	}

	for _, tt := range tools {
		t.Run(tt.name, func(t *testing.T) {
			// Validate tool metadata
			assert.NotEmpty(t, tt.tool.GetName())
			assert.NotEmpty(t, tt.tool.GetDescription())
			
			// Validate schema structure
			schema := tt.tool.GetInputSchema()
			require.NotNil(t, schema)
			
			schemaMap, ok := schema.(map[string]interface{})
			require.True(t, ok, "Schema must be a map")
			assert.Equal(t, "object", schemaMap["type"], "Schema must be an object type")
			
			properties, exists := schemaMap["properties"].(map[string]interface{})
			require.True(t, exists, "Schema must have properties")
			assert.Contains(t, properties, "provider", "Schema must include provider property")
			
			t.Logf("✅ Tool %s has valid schema", tt.tool.GetName())
		})
	}
}

// Helper functions for test setup
func mustCreateDiscoverTool(t *testing.T, config interfaces.LighthouseConfig) *mcp.FinOpsDiscoverTool {
	tool, err := mcp.NewFinOpsDiscoverTool(config)
	require.NoError(t, err)
	return tool
}

func mustCreateRecommendTool(t *testing.T, config interfaces.LighthouseConfig) *mcp.FinOpsRecommendTool {
	tool, err := mcp.NewFinOpsRecommendTool(config)
	require.NoError(t, err)
	return tool
}

func mustCreateAnalyzeTool(t *testing.T, config interfaces.LighthouseConfig) *mcp.FinOpsAnalyzeTool {
	tool, err := mcp.NewFinOpsAnalyzeTool(config)
	require.NoError(t, err)
	return tool
}

func mustCreateQueryTool(t *testing.T, config interfaces.LighthouseConfig) *mcp.FinOpsQueryTool {
	tool, err := mcp.NewFinOpsQueryTool(config)
	require.NoError(t, err)
	return tool
}