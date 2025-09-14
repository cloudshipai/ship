package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
)

// CostExplorerClient wraps AWS Cost Explorer functionality
type CostExplorerClient struct {
	client *costexplorer.Client
}

// NewCostExplorerClient creates a new Cost Explorer client wrapper
func NewCostExplorerClient(client *costexplorer.Client) *CostExplorerClient {
	return &CostExplorerClient{
		client: client,
	}
}

// GetCostData retrieves cost data from AWS Cost Explorer
func (c *CostExplorerClient) GetCostData(ctx context.Context, opts interfaces.CostOptions) ([]interfaces.CostRecord, error) {
	// Parse time window
	timeRange, err := c.parseTimeWindow(opts.TimeWindow)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time window: %w", err)
	}
	
	// Determine granularity
	granularity := c.getGranularity(opts.Granularity)
	
	// Build dimensions for grouping
	var groupBy []types.GroupDefinition
	for _, groupByField := range opts.GroupBy {
		dimension, err := c.mapToCostExplorerDimension(groupByField)
		if err != nil {
			continue // Skip invalid dimensions
		}
		groupBy = append(groupBy, types.GroupDefinition{
			Type: types.GroupDefinitionTypeDimension,
			Key:  &dimension,
		})
	}
	
	// Prepare Cost Explorer request  
	startStr := timeRange.Start.Format("2006-01-02")
	endStr := timeRange.End.Format("2006-01-02")
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: &startStr,
			End:   &endStr,
		},
		Granularity: granularity,
		GroupBy:     groupBy,
		Metrics:     []string{"BlendedCost", "UnblendedCost", "UsageQuantity"},
	}
	
	// Add filters if specified
	if len(opts.Filters) > 0 {
		input.Filter = c.buildCostFilter(opts.Filters)
	}
	
	// Call Cost Explorer API
	resp, err := c.client.GetCostAndUsage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost and usage data: %w", err)
	}
	
	// Transform response to our standard format
	return c.transformCostResponse(resp, opts.Currency)
}

// GetRightsizingRecommendations retrieves rightsizing recommendations from Cost Explorer
func (c *CostExplorerClient) GetRightsizingRecommendations(ctx context.Context, accountID string, regions []string) ([]interfaces.Recommendation, error) {
	service := "EC2-Instance"
	input := &costexplorer.GetRightsizingRecommendationInput{
		Service: &service,
	}
	
	// Add filters if specified
	if accountID != "" || len(regions) > 0 {
		filter := c.buildRightsizingFilter(accountID, regions)
		if filter != nil {
			input.Filter = filter
		}
	}
	
	resp, err := c.client.GetRightsizingRecommendation(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get rightsizing recommendations: %w", err)
	}
	
	return c.transformRightsizingRecommendations(resp)
}

// parseTimeWindow converts time window string to date range
func (c *CostExplorerClient) parseTimeWindow(window string) (*interfaces.TimePeriod, error) {
	now := time.Now()
	var start, end time.Time
	
	switch window {
	case "7d":
		start = now.AddDate(0, 0, -7)
		end = now
	case "30d":
		start = now.AddDate(0, 0, -30)
		end = now
	case "90d":
		start = now.AddDate(0, 0, -90)
		end = now
	case "1y":
		start = now.AddDate(-1, 0, 0)
		end = now
	default:
		return nil, fmt.Errorf("unsupported time window: %s", window)
	}
	
	return &interfaces.TimePeriod{
		Start: start,
		End:   end,
	}, nil
}

// getGranularity maps granularity string to Cost Explorer granularity
func (c *CostExplorerClient) getGranularity(granularity string) types.Granularity {
	switch granularity {
	case "daily":
		return types.GranularityDaily
	case "monthly":
		return types.GranularityMonthly
	default:
		return types.GranularityDaily
	}
}

// mapToCostExplorerDimension maps our group-by fields to Cost Explorer dimensions
func (c *CostExplorerClient) mapToCostExplorerDimension(field string) (string, error) {
	dimensionMap := map[string]string{
		"service":    "SERVICE",
		"region":     "REGION",
		"account":    "LINKED_ACCOUNT",
		"instance":   "INSTANCE_TYPE",
		"platform":   "PLATFORM",
		"tenancy":    "TENANCY",
		"az":         "AZ",
	}
	
	if dimension, ok := dimensionMap[field]; ok {
		return dimension, nil
	}
	
	return "", fmt.Errorf("unsupported group-by field: %s", field)
}

// buildCostFilter creates Cost Explorer filters from our filter map
func (c *CostExplorerClient) buildCostFilter(filters map[string]string) *types.Expression {
	var expressions []types.Expression
	
	for key, value := range filters {
		dimension, err := c.mapToCostExplorerDimension(key)
		if err != nil {
			continue // Skip invalid filters
		}
		
		expr := types.Expression{
			Dimensions: &types.DimensionValues{
				Key:        types.Dimension(dimension),
				Values:     []string{value},
				MatchOptions: []types.MatchOption{types.MatchOptionEquals},
			},
		}
		expressions = append(expressions, expr)
	}
	
	if len(expressions) == 0 {
		return nil
	}
	
	if len(expressions) == 1 {
		return &expressions[0]
	}
	
	// Combine multiple expressions with AND
	return &types.Expression{
		And: expressions,
	}
}

// buildRightsizingFilter creates filters for rightsizing recommendations
func (c *CostExplorerClient) buildRightsizingFilter(accountID string, regions []string) *types.Expression {
	var expressions []types.Expression
	
	if accountID != "" {
		accountDim := "LINKED_ACCOUNT"
		expressions = append(expressions, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:          types.Dimension(accountDim),
				Values:       []string{accountID},
				MatchOptions: []types.MatchOption{types.MatchOptionEquals},
			},
		})
	}
	
	if len(regions) > 0 {
		regionDim := "REGION"
		expressions = append(expressions, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:          types.Dimension(regionDim),
				Values:       regions,
				MatchOptions: []types.MatchOption{types.MatchOptionEquals},
			},
		})
	}
	
	if len(expressions) == 0 {
		return nil
	}
	
	if len(expressions) == 1 {
		return &expressions[0]
	}
	
	return &types.Expression{
		And: expressions,
	}
}

// transformCostResponse converts Cost Explorer response to our standard format
func (c *CostExplorerClient) transformCostResponse(resp *costexplorer.GetCostAndUsageOutput, currency string) ([]interfaces.CostRecord, error) {
	var records []interfaces.CostRecord
	
	for _, result := range resp.ResultsByTime {
		// Parse time period
		startTime, _ := time.Parse("2006-01-02", *result.TimePeriod.Start)
		endTime, _ := time.Parse("2006-01-02", *result.TimePeriod.End)
		
		period := interfaces.TimePeriod{
			Start: startTime,
			End:   endTime,
		}
		
		// Handle ungrouped data
		if len(result.Groups) == 0 {
			for metricName, metric := range result.Total {
				if metricName == "BlendedCost" || metricName == "UnblendedCost" {
					amount := 0.0
					if metric.Amount != nil {
						fmt.Sscanf(*metric.Amount, "%f", &amount)
					}
					
					record := interfaces.CostRecord{
						ResourceID: "total",
						Period:     period,
						Amount:     amount,
						Currency:   currency,
						Service:    "all",
					}
					records = append(records, record)
				}
			}
		} else {
			// Handle grouped data
			for _, group := range result.Groups {
				for metricName, metric := range group.Metrics {
					if metricName == "BlendedCost" || metricName == "UnblendedCost" {
						amount := 0.0
						if metric.Amount != nil {
							fmt.Sscanf(*metric.Amount, "%f", &amount)
						}
						
						record := interfaces.CostRecord{
							Period:   period,
							Amount:   amount,
							Currency: currency,
						}
						
						// Extract group dimensions
						if len(group.Keys) > 0 {
							record.Service = group.Keys[0]
						}
						if len(group.Keys) > 1 {
							record.Region = group.Keys[1]
						}
						
						records = append(records, record)
					}
				}
			}
		}
	}
	
	return records, nil
}

// transformRightsizingRecommendations converts Cost Explorer rightsizing recommendations
func (c *CostExplorerClient) transformRightsizingRecommendations(resp *costexplorer.GetRightsizingRecommendationOutput) ([]interfaces.Recommendation, error) {
	var recommendations []interfaces.Recommendation
	
	for _, awsRec := range resp.RightsizingRecommendations {
		if awsRec.CurrentInstance == nil || awsRec.CurrentInstance.ResourceId == nil {
			continue
		}
		
		rec := interfaces.Recommendation{
			ID:                 fmt.Sprintf("ce-%s", *awsRec.CurrentInstance.ResourceId),
			ResourceID:         *awsRec.CurrentInstance.ResourceId,
			ResourceType:       interfaces.ResourceTypeCompute,
			Provider:           interfaces.VendorAWS,
			RecommendationType: "RightSizeInstance",
			Title:              fmt.Sprintf("Right-size EC2 instance %s", *awsRec.CurrentInstance.ResourceId),
			Description:        c.generateRightsizingDescription(awsRec),
		}
		
		// Extract savings information from target instances
		totalSavings := 0.0
		if awsRec.ModifyRecommendationDetail != nil {
			for _, targetInstance := range awsRec.ModifyRecommendationDetail.TargetInstances {
				if targetInstance.EstimatedMonthlySavings != nil {
					amount := 0.0
					fmt.Sscanf(*targetInstance.EstimatedMonthlySavings, "%f", &amount)
					if amount > totalSavings {
						totalSavings = amount
					}
				}
			}
		}
		
		rec.EstimatedSavings = interfaces.EstimatedSavings{
			Currency: "USD",
			Monthly:  totalSavings,
		}
		
		// Add recommendation options based on rightsizing type
		switch awsRec.RightsizingType {
		case types.RightsizingTypeTerminate:
			rec.Options = append(rec.Options, interfaces.RecommendationOption{
				OptionType:      "terminate",
				Description:     "Terminate unused instance",
				ExpectedSavings: totalSavings,
			})
		case types.RightsizingTypeModify:
			if awsRec.ModifyRecommendationDetail != nil {
				for _, targetInstance := range awsRec.ModifyRecommendationDetail.TargetInstances {
					if targetInstance.DefaultTargetInstance && targetInstance.ResourceDetails != nil {
						savings := 0.0
						if targetInstance.EstimatedMonthlySavings != nil {
							fmt.Sscanf(*targetInstance.EstimatedMonthlySavings, "%f", &savings)
						}
						
						rec.Options = append(rec.Options, interfaces.RecommendationOption{
							OptionType:      "modify",
							Description:     "Change instance type",
							ExpectedSavings: savings,
							Parameters: map[string]interface{}{
								"default_target": targetInstance.DefaultTargetInstance,
							},
						})
					}
				}
			}
		}
		
		recommendations = append(recommendations, rec)
	}
	
	return recommendations, nil
}

// generateRightsizingDescription creates a description for rightsizing recommendations
func (c *CostExplorerClient) generateRightsizingDescription(awsRec types.RightsizingRecommendation) string {
	if awsRec.CurrentInstance == nil {
		return "AWS Cost Explorer recommends optimizing this EC2 instance."
	}
	
	resourceId := ""
	if awsRec.CurrentInstance.ResourceId != nil {
		resourceId = *awsRec.CurrentInstance.ResourceId
	}
	
	switch awsRec.RightsizingType {
	case types.RightsizingTypeTerminate:
		return fmt.Sprintf("AWS Cost Explorer recommends terminating instance %s as it appears to be underutilized or idle.", resourceId)
	case types.RightsizingTypeModify:
		return fmt.Sprintf("AWS Cost Explorer recommends modifying instance %s to a more cost-effective instance type.", resourceId)
	default:
		return fmt.Sprintf("AWS Cost Explorer recommends optimizing instance %s.", resourceId)
	}
}