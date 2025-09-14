package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/cloudshipai/ship/internal/tools/finops/mcp"
)

// finOpsCmd represents the finops command
var finOpsCmd = &cobra.Command{
	Use:   "finops",
	Short: "FinOps cost optimization tools",
	Long:  `FinOps tools for cloud cost optimization, resource discovery, and recommendations.`,
}

// finOpsDiscoverCmd represents the finops-discover command
var finOpsDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover cloud resources with cost optimization data",
	Long: `Discover cloud resources across providers (AWS, GCP, Azure, Kubernetes) 
with cost and utilization data for optimization analysis.`,
	RunE: runFinOpsDiscover,
}

// finOpsRecommendCmd represents the finops-recommend command
var finOpsRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Generate cost optimization recommendations",
	Long: `Generate cost optimization recommendations using vendor-specific 
recommendation engines like AWS Compute Optimizer.`,
	RunE: runFinOpsRecommend,
}

// finOpsAnalyzeCmd represents the finops-analyze command
var finOpsAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze cost data and trends",
	Long: `Analyze cost data and trends with insights and anomaly detection 
across cloud providers.`,
	RunE: runFinOpsAnalyze,
}

// finOpsQueryCmd represents the finops-query command (agent-driven)
var finOpsQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Agent-driven flexible finops queries",
	Long: `Execute flexible finops queries with natural language support 
for agent-driven cost optimization workflows.`,
	RunE: runFinOpsQuery,
}

func init() {
	rootCmd.AddCommand(finOpsCmd)
	finOpsCmd.AddCommand(finOpsDiscoverCmd)
	finOpsCmd.AddCommand(finOpsRecommendCmd)
	finOpsCmd.AddCommand(finOpsAnalyzeCmd)
	finOpsCmd.AddCommand(finOpsQueryCmd)

	// Common flags for all finops commands
	finOpsDiscoverCmd.Flags().String("provider", "aws", "Cloud provider (aws, gcp, azure, kubernetes)")
	finOpsDiscoverCmd.Flags().String("region", "", "Region filter")
	finOpsDiscoverCmd.Flags().StringSlice("resource-types", []string{}, "Resource types to discover")
	finOpsDiscoverCmd.Flags().StringToString("tags", map[string]string{}, "Tag filters")
	finOpsDiscoverCmd.Flags().StringSlice("account-ids", []string{}, "Account IDs for multi-account discovery")
	finOpsDiscoverCmd.Flags().StringSlice("arns", []string{}, "Resource ARNs to filter by")
	finOpsDiscoverCmd.Flags().Bool("cloudshipai", false, "Send data to CloudshipAI via Station")

	finOpsRecommendCmd.Flags().String("provider", "aws", "Cloud provider")
	finOpsRecommendCmd.Flags().StringSlice("finding-types", []string{"rightsizing"}, "Types of recommendations")
	finOpsRecommendCmd.Flags().StringSlice("regions", []string{}, "Regions to analyze")
	finOpsRecommendCmd.Flags().StringSlice("arns", []string{}, "Resource ARNs to analyze")
	finOpsRecommendCmd.Flags().Float64("min-savings", 0, "Minimum monthly savings threshold")
	finOpsRecommendCmd.Flags().Bool("cloudshipai", false, "Send data to CloudshipAI via Station")

	finOpsAnalyzeCmd.Flags().String("provider", "aws", "Cloud provider")
	finOpsAnalyzeCmd.Flags().String("time-window", "30d", "Time window for analysis")
	finOpsAnalyzeCmd.Flags().String("granularity", "daily", "Data granularity")
	finOpsAnalyzeCmd.Flags().StringSlice("group-by", []string{}, "Dimensions to group by")
	finOpsAnalyzeCmd.Flags().String("currency", "USD", "Currency for cost data")
	finOpsAnalyzeCmd.Flags().Bool("cloudshipai", false, "Send data to CloudshipAI via Station")

	finOpsQueryCmd.Flags().String("query", "", "Natural language query")
	finOpsQueryCmd.Flags().String("provider", "aws", "Cloud provider")
	finOpsQueryCmd.Flags().Bool("cloudshipai", false, "Send data to CloudshipAI via Station")
}

func runFinOpsDiscover(cmd *cobra.Command, args []string) error {
	// Get flags
	provider, _ := cmd.Flags().GetString("provider")
	region, _ := cmd.Flags().GetString("region")
	resourceTypes, _ := cmd.Flags().GetStringSlice("resource-types")
	tags, _ := cmd.Flags().GetStringToString("tags")
	accountIDs, _ := cmd.Flags().GetStringSlice("account-ids")
	arns, _ := cmd.Flags().GetStringSlice("arns")
	cloudshipai, _ := cmd.Flags().GetBool("cloudshipai")

	// Create lighthouse config
	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   cloudshipai,
	}

	// Create the finops discover tool
	tool, err := mcp.NewFinOpsDiscoverTool(lighthouseConfig)
	if err != nil {
		return fmt.Errorf("failed to create finops discover tool: %w", err)
	}

	// Build request
	request := mcp.DiscoverRequest{
		Provider:      provider,
		Region:        region,
		ResourceTypes: resourceTypes,
		Tags:          tags,
		AccountIDs:    accountIDs,
		ARNs:          arns,
	}

	ctx := context.Background()
	requestJSON, _ := json.Marshal(request)
	
	// Execute the tool
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return fmt.Errorf("failed to execute finops discover: %w", err)
	}

	// Print results
	fmt.Printf("FinOps Discovery Results:\n")
	printJSON(result)

	if cloudshipai {
		fmt.Printf("\n✅ Data sent to CloudshipAI via Station\n")
	}

	return nil
}

func runFinOpsRecommend(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	findingTypes, _ := cmd.Flags().GetStringSlice("finding-types")
	regions, _ := cmd.Flags().GetStringSlice("regions")
	arns, _ := cmd.Flags().GetStringSlice("arns")
	minSavings, _ := cmd.Flags().GetFloat64("min-savings")
	cloudshipai, _ := cmd.Flags().GetBool("cloudshipai")

	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   cloudshipai,
	}

	tool, err := mcp.NewFinOpsRecommendTool(lighthouseConfig)
	if err != nil {
		return fmt.Errorf("failed to create finops recommend tool: %w", err)
	}

	request := mcp.RecommendRequest{
		Provider:     provider,
		FindingTypes: findingTypes,
		Regions:      regions,
		ARNs:         arns,
		MinSavings:   minSavings,
	}

	ctx := context.Background()
	requestJSON, _ := json.Marshal(request)
	
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return fmt.Errorf("failed to execute finops recommend: %w", err)
	}

	fmt.Printf("FinOps Recommendations:\n")
	printJSON(result)

	if cloudshipai {
		fmt.Printf("\n✅ Data sent to CloudshipAI via Station\n")
	}

	return nil
}

func runFinOpsAnalyze(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	timeWindow, _ := cmd.Flags().GetString("time-window")
	granularity, _ := cmd.Flags().GetString("granularity")
	groupBy, _ := cmd.Flags().GetStringSlice("group-by")
	currency, _ := cmd.Flags().GetString("currency")
	cloudshipai, _ := cmd.Flags().GetBool("cloudshipai")

	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   cloudshipai,
	}

	tool, err := mcp.NewFinOpsAnalyzeTool(lighthouseConfig)
	if err != nil {
		return fmt.Errorf("failed to create finops analyze tool: %w", err)
	}

	request := mcp.AnalyzeRequest{
		Provider:    provider,
		TimeWindow:  timeWindow,
		Granularity: granularity,
		GroupBy:     groupBy,
		Currency:    currency,
	}

	ctx := context.Background()
	requestJSON, _ := json.Marshal(request)
	
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return fmt.Errorf("failed to execute finops analyze: %w", err)
	}

	fmt.Printf("FinOps Cost Analysis:\n")
	printJSON(result)

	if cloudshipai {
		fmt.Printf("\n✅ Data sent to CloudshipAI via Station\n")
	}

	return nil
}

func runFinOpsQuery(cmd *cobra.Command, args []string) error {
	query, _ := cmd.Flags().GetString("query")
	provider, _ := cmd.Flags().GetString("provider")
	cloudshipai, _ := cmd.Flags().GetBool("cloudshipai")

	if query == "" {
		return fmt.Errorf("query is required")
	}

	lighthouseConfig := interfaces.LighthouseConfig{
		Timeout:         30 * time.Second,
		RetryAttempts:   3,
		BatchSize:       1000,
		ValidateSchema:  true,
		EnableTracing:   cloudshipai,
	}

	tool, err := mcp.NewFinOpsQueryTool(lighthouseConfig)
	if err != nil {
		return fmt.Errorf("failed to create finops query tool: %w", err)
	}

	request := mcp.QueryRequest{
		Query:    query,
		Provider: provider,
	}

	ctx := context.Background()
	requestJSON, _ := json.Marshal(request)
	
	result, err := tool.Execute(ctx, requestJSON)
	if err != nil {
		return fmt.Errorf("failed to execute finops query: %w", err)
	}

	fmt.Printf("FinOps Query Results:\n")
	printJSON(result)

	if cloudshipai {
		fmt.Printf("\n✅ Data sent to CloudshipAI via Station\n")
	}

	return nil
}

// Helper functions
func printJSON(v interface{}) {
	if jsonBytes, err := json.MarshalIndent(v, "", "  "); err == nil {
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%+v\n", v)
	}
}