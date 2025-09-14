package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/mcp"
	"github.com/mark3labs/mcp-go/server"
	mcpTypes "github.com/mark3labs/mcp-go/mcp"
)

// AddFinOpsDiscoverTools registers the finops-discover MCP tool
func AddFinOpsDiscoverTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	tool := mcpTypes.NewTool("finops-discover",
		mcpTypes.WithDescription("Discover cloud resources across providers (AWS, GCP, Azure, Kubernetes) with cost and utilization data"),
		mcpTypes.WithString("provider",
			mcpTypes.Description("Cloud provider to discover resources from"),
			mcpTypes.Required(),
		),
		mcpTypes.WithString("region",
			mcpTypes.Description("Optional region filter (provider-specific)"),
		),
		mcpTypes.WithArray("resource_types",
			mcpTypes.Description("Optional resource types to discover"),
		),
		mcpTypes.WithObject("tags",
			mcpTypes.Description("Optional tag filters (key-value pairs)"),
		),
		mcpTypes.WithArray("account_ids",
			mcpTypes.Description("Optional account IDs for multi-account discovery"),
		),
		mcpTypes.WithArray("arns",
			mcpTypes.Description("Optional resource ARNs to filter by (AWS only)"),
		),
		mcpTypes.WithBoolean("cloudshipai",
			mcpTypes.Description("Enable reporting results to CloudshipAI Station API"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
		return executeFinOpsDiscoverTool(ctx, request)
	})
}

// AddFinOpsRecommendTools registers the finops-recommend MCP tool
func AddFinOpsRecommendTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	tool := mcpTypes.NewTool("finops-recommend",
		mcpTypes.WithDescription("Generate cost optimization recommendations using vendor-specific recommendation engines like AWS Compute Optimizer"),
		mcpTypes.WithString("provider",
			mcpTypes.Description("Cloud provider to analyze"),
			mcpTypes.Required(),
		),
		mcpTypes.WithArray("finding_types",
			mcpTypes.Description("Types of recommendations to generate"),
		),
		mcpTypes.WithArray("regions",
			mcpTypes.Description("Regions to analyze"),
		),
		mcpTypes.WithArray("arns",
			mcpTypes.Description("Resource ARNs to analyze (AWS only)"),
		),
		mcpTypes.WithNumber("min_savings",
			mcpTypes.Description("Minimum monthly savings threshold in USD"),
		),
		mcpTypes.WithBoolean("cloudshipai",
			mcpTypes.Description("Enable reporting results to CloudshipAI Station API"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
		return executeFinOpsRecommendTool(ctx, request)
	})
}

// AddFinOpsAnalyzeTools registers the finops-analyze MCP tool
func AddFinOpsAnalyzeTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	tool := mcpTypes.NewTool("finops-analyze",
		mcpTypes.WithDescription("Analyze cost data and trends with insights and anomaly detection across cloud providers"),
		mcpTypes.WithString("provider",
			mcpTypes.Description("Cloud provider to analyze"),
			mcpTypes.Required(),
		),
		mcpTypes.WithString("time_window",
			mcpTypes.Description("Time window for analysis"),
		),
		mcpTypes.WithString("granularity",
			mcpTypes.Description("Data granularity"),
		),
		mcpTypes.WithArray("group_by",
			mcpTypes.Description("Dimensions to group by"),
		),
		mcpTypes.WithString("currency",
			mcpTypes.Description("Currency for cost data"),
		),
		mcpTypes.WithBoolean("cloudshipai",
			mcpTypes.Description("Enable reporting results to CloudshipAI Station API"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
		return executeFinOpsAnalyzeTool(ctx, request)
	})
}

// AddFinOpsQueryTools registers the finops-query MCP tool
func AddFinOpsQueryTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
	tool := mcpTypes.NewTool("finops-query",
		mcpTypes.WithDescription("Execute flexible finops queries with natural language support for agent-driven cost optimization workflows"),
		mcpTypes.WithString("query",
			mcpTypes.Description("Natural language query about finops data"),
			mcpTypes.Required(),
		),
		mcpTypes.WithString("provider",
			mcpTypes.Description("Cloud provider context"),
		),
		mcpTypes.WithObject("context",
			mcpTypes.Description("Additional context for the query"),
		),
		mcpTypes.WithObject("filters",
			mcpTypes.Description("Query filters"),
		),
		mcpTypes.WithBoolean("cloudshipai",
			mcpTypes.Description("Enable reporting results to CloudshipAI Station API"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
		return executeFinOpsQueryTool(ctx, request)
	})
}

// Tool execution functions

func executeFinOpsDiscoverTool(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
	// Convert arguments to map for easier processing
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	// Build request from arguments
	req := mcp.DiscoverRequest{
		Provider: getString(args, "provider"),
		Region:   getString(args, "region"),
		EnableCloudshipAI: getBool(args, "cloudshipai"),
	}

	// Handle arrays
	if resourceTypes := getStringArray(args, "resource_types"); resourceTypes != nil {
		req.ResourceTypes = resourceTypes
	}
	if accountIDs := getStringArray(args, "account_ids"); accountIDs != nil {
		req.AccountIDs = accountIDs
	}
	if arns := getStringArray(args, "arns"); arns != nil {
		req.ARNs = arns
	}
	if tags := getStringMap(args, "tags"); tags != nil {
		req.Tags = tags
	}

	// Create lighthouse config
	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   true,
	}

	// Create tool and execute
	tool, err := mcp.NewFinOpsDiscoverTool(lighthouseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create finops discover tool: %w", err)
	}

	requestJSON, _ := json.Marshal(req)
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	// Convert result to MCP response
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return &mcpTypes.CallToolResult{
		Content: []mcpTypes.Content{
			mcpTypes.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

func executeFinOpsRecommendTool(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	req := mcp.RecommendRequest{
		Provider:          getString(args, "provider"),
		EnableCloudshipAI: getBool(args, "cloudshipai"),
		MinSavings:       getFloat64(args, "min_savings"),
	}

	if findingTypes := getStringArray(args, "finding_types"); findingTypes != nil {
		req.FindingTypes = findingTypes
	}
	if regions := getStringArray(args, "regions"); regions != nil {
		req.Regions = regions
	}
	if arns := getStringArray(args, "arns"); arns != nil {
		req.ARNs = arns
	}

	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   true,
	}

	tool, err := mcp.NewFinOpsRecommendTool(lighthouseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create finops recommend tool: %w", err)
	}

	requestJSON, _ := json.Marshal(req)
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return &mcpTypes.CallToolResult{
		Content: []mcpTypes.Content{
			mcpTypes.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

func executeFinOpsAnalyzeTool(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	req := mcp.AnalyzeRequest{
		Provider:          getString(args, "provider"),
		TimeWindow:        getString(args, "time_window"),
		Granularity:       getString(args, "granularity"),
		Currency:          getString(args, "currency"),
		EnableCloudshipAI: getBool(args, "cloudshipai"),
	}

	if groupBy := getStringArray(args, "group_by"); groupBy != nil {
		req.GroupBy = groupBy
	}

	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   true,
	}

	tool, err := mcp.NewFinOpsAnalyzeTool(lighthouseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create finops analyze tool: %w", err)
	}

	requestJSON, _ := json.Marshal(req)
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return &mcpTypes.CallToolResult{
		Content: []mcpTypes.Content{
			mcpTypes.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

func executeFinOpsQueryTool(ctx context.Context, request mcpTypes.CallToolRequest) (*mcpTypes.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	req := mcp.QueryRequest{
		Query:             getString(args, "query"),
		Provider:          getString(args, "provider"),
		EnableCloudshipAI: getBool(args, "cloudshipai"),
	}

	// Handle context and filters as generic interface{}
	if context := args["context"]; context != nil {
		req.Context = context.(map[string]interface{})
	}
	if filters := args["filters"]; filters != nil {
		if filtersMap, ok := filters.(map[string]interface{}); ok {
			req.Filters = mcp.QueryFilters{
				Regions:       getStringArray(filtersMap, "regions"),
				ResourceTypes: getStringArray(filtersMap, "resource_types"),
				MinSavings:    getFloat64(filtersMap, "min_savings"),
			}
		}
	}

	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   true,
	}

	tool, err := mcp.NewFinOpsQueryTool(lighthouseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create finops query tool: %w", err)
	}

	requestJSON, _ := json.Marshal(req)
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return &mcpTypes.CallToolResult{
		Content: []mcpTypes.Content{
			mcpTypes.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// Helper functions for extracting values from MCP arguments
func getString(args map[string]interface{}, key string) string {
	if val, ok := args[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBool(args map[string]interface{}, key string) bool {
	if val, ok := args[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func getFloat64(args map[string]interface{}, key string) float64 {
	if val, ok := args[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
		if i, ok := val.(int); ok {
			return float64(i)
		}
	}
	return 0
}

func getStringArray(args map[string]interface{}, key string) []string {
	if val, ok := args[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

func getStringMap(args map[string]interface{}, key string) map[string]string {
	if val, ok := args[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			result := make(map[string]string)
			for k, v := range m {
				if str, ok := v.(string); ok {
					result[k] = str
				}
			}
			return result
		}
	}
	return nil
}
