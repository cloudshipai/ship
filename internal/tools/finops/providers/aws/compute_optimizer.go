package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/computeoptimizer"
	"github.com/aws/aws-sdk-go-v2/service/computeoptimizer/types"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
)

// EC2RecommendationsBuilder follows OpenOps pattern for building recommendations
type EC2RecommendationsBuilder struct {
	client           *computeoptimizer.Client
	findingType      types.Finding
	recommendationType string
}

// NewEC2RecommendationsBuilder creates a new builder instance
func NewEC2RecommendationsBuilder(client *computeoptimizer.Client, findingType types.Finding) *EC2RecommendationsBuilder {
	return &EC2RecommendationsBuilder{
		client:           client,
		findingType:      findingType,
		recommendationType: getRecommendationType(findingType),
	}
}

// GetRecommendations gets EC2 rightsizing recommendations for a region
// Based on OpenOps getEC2RecommendationsForRegions pattern
func (b *EC2RecommendationsBuilder) GetRecommendations(ctx context.Context, region string, arns []string) ([]interfaces.Recommendation, error) {
	// Prepare input based on whether ARNs are specified (OpenOps pattern)
	input := &computeoptimizer.GetEC2InstanceRecommendationsInput{
		Filters: []types.Filter{
			{
				Name:   types.FilterNameFinding,
				Values: []string{string(b.findingType)},
			},
		},
	}
	
	// Add ARN filtering if specified (OpenOps pattern)
	if len(arns) > 0 {
		input.InstanceArns = arns
	}
	
	// Call AWS Compute Optimizer API
	resp, err := b.client.GetEC2InstanceRecommendations(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get EC2 recommendations: %w", err)
	}
	
	var recommendations []interfaces.Recommendation
	for _, awsRec := range resp.InstanceRecommendations {
		// Transform AWS recommendation to our standard format
		rec, err := b.transformRecommendation(awsRec, region)
		if err != nil {
			return nil, fmt.Errorf("failed to transform recommendation: %w", err)
		}
		recommendations = append(recommendations, rec)
	}
	
	return recommendations, nil
}

// transformRecommendation converts AWS Compute Optimizer recommendation to our format
// Based on OpenOps ec2-recommendations-builder transformation pattern
func (b *EC2RecommendationsBuilder) transformRecommendation(awsRec types.InstanceRecommendation, region string) (interfaces.Recommendation, error) {
	// Filter out options without savings (OpenOps pattern)
	validOptions := b.filterOptionsWithSavings(awsRec.RecommendationOptions)
	
	// Sort options by rank (OpenOps pattern)
	sortedOptions := b.sortOptionsByRank(validOptions)
	
	// Get the best option for primary recommendation
	var bestOption *interfaces.RecommendationOption
	var totalSavings float64
	
	if len(sortedOptions) > 0 {
		bestAWSOption := sortedOptions[0]
		bestOption = &interfaces.RecommendationOption{
			OptionType:      "instance_type_change",
			Description:     fmt.Sprintf("Change from %s to %s", *awsRec.CurrentInstanceType, *bestAWSOption.InstanceType),
			ExpectedSavings: b.getSavingsValue(bestAWSOption.SavingsOpportunity),
			Parameters: map[string]interface{}{
				"current_instance_type":    *awsRec.CurrentInstanceType,
				"recommended_instance_type": *bestAWSOption.InstanceType,
				"migration_effort":         string(bestAWSOption.MigrationEffort),
				"performance_risk":         bestAWSOption.PerformanceRisk,
			},
		}
		totalSavings = bestOption.ExpectedSavings
	}
	
	// Build recommendation following our standard format
	rec := interfaces.Recommendation{
		ID:                 fmt.Sprintf("ec2-%s", extractInstanceID(*awsRec.InstanceArn)),
		ResourceID:         extractInstanceID(*awsRec.InstanceArn),
		ResourceARN:        *awsRec.InstanceArn,
		ResourceType:       interfaces.ResourceTypeCompute,
		Provider:           interfaces.VendorAWS,
		RecommendationType: b.recommendationType,
		Title:              b.generateRecommendationTitle(awsRec, bestOption),
		Description:        b.generateRecommendationDescription(awsRec, bestOption),
		EstimatedSavings: interfaces.EstimatedSavings{
			Currency: "USD",
			Monthly:  totalSavings,
		},
		PerformanceRisk: b.getPerformanceRisk(bestOption),
		MigrationEffort: b.getMigrationEffort(bestOption),
		Confidence:      b.calculateConfidence(&awsRec.Finding),
		Justification:   b.generateJustification(awsRec, bestOption),
		CreatedAt:       *awsRec.LastRefreshTimestamp,
	}
	
	// Add all options
	for _, option := range sortedOptions {
		rec.Options = append(rec.Options, interfaces.RecommendationOption{
			OptionType:      "instance_type_change",
			Description:     fmt.Sprintf("Change to %s", *option.InstanceType),
			ExpectedSavings: b.getSavingsValue(option.SavingsOpportunity),
			Parameters: map[string]interface{}{
				"instance_type":     *option.InstanceType,
				"migration_effort":  string(option.MigrationEffort),
				"performance_risk":  option.PerformanceRisk,
				"rank":             option.Rank,
			},
		})
	}
	
	return rec, nil
}

// filterOptionsWithSavings filters out recommendation options without savings (OpenOps pattern)
func (b *EC2RecommendationsBuilder) filterOptionsWithSavings(options []types.InstanceRecommendationOption) []types.InstanceRecommendationOption {
	var filtered []types.InstanceRecommendationOption
	for _, option := range options {
		if option.SavingsOpportunity != nil && 
		   option.SavingsOpportunity.EstimatedMonthlySavings != nil &&
		   option.SavingsOpportunity.EstimatedMonthlySavings.Value > 0 {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

// sortOptionsByRank sorts recommendation options by rank (OpenOps pattern)  
func (b *EC2RecommendationsBuilder) sortOptionsByRank(options []types.InstanceRecommendationOption) []types.InstanceRecommendationOption {
	// AWS provides options pre-sorted by rank, but we can add custom sorting if needed
	return options
}

// Helper functions based on OpenOps patterns

func getRecommendationType(findingType types.Finding) string {
	switch findingType {
	case types.FindingOptimized:
		return "UpgradeInstanceGeneration"
	case types.FindingOverProvisioned:
		return "RightSizeInstance" 
	case types.FindingUnderProvisioned:
		return "RightSizeInstance"
	default:
		return "RightSizeInstance"
	}
}

func (b *EC2RecommendationsBuilder) getSavingsValue(savings *types.SavingsOpportunity) float64 {
	if savings == nil || savings.EstimatedMonthlySavings == nil {
		return 0
	}
	return savings.EstimatedMonthlySavings.Value
}

func (b *EC2RecommendationsBuilder) getPerformanceRisk(option *interfaces.RecommendationOption) string {
	if option == nil {
		return "Medium"
	}
	if risk, exists := option.Parameters["performance_risk"].(string); exists {
		return risk
	}
	return "Medium"
}

func (b *EC2RecommendationsBuilder) getMigrationEffort(option *interfaces.RecommendationOption) string {
	if option == nil {
		return "Medium"
	}
	if effort, exists := option.Parameters["migration_effort"].(string); exists {
		return effort
	}
	return "Medium"
}

func (b *EC2RecommendationsBuilder) calculateConfidence(finding *types.Finding) float64 {
	if finding == nil {
		return 0.5
	}
	
	switch *finding {
	case types.FindingOverProvisioned:
		return 0.85 // High confidence for overprovisioned
	case types.FindingUnderProvisioned:
		return 0.75 // Good confidence for underprovisioned
	case types.FindingOptimized:
		return 0.70 // Medium confidence for generation upgrades
	default:
		return 0.5
	}
}

func (b *EC2RecommendationsBuilder) generateRecommendationTitle(awsRec types.InstanceRecommendation, bestOption *interfaces.RecommendationOption) string {
	instanceID := extractInstanceID(*awsRec.InstanceArn)
	
	if bestOption != nil {
		if newType, exists := bestOption.Parameters["recommended_instance_type"].(string); exists {
			return fmt.Sprintf("Right-size %s from %s to %s", instanceID, *awsRec.CurrentInstanceType, newType)
		}
	}
	
	return fmt.Sprintf("Optimize EC2 instance %s", instanceID)
}

func (b *EC2RecommendationsBuilder) generateRecommendationDescription(awsRec types.InstanceRecommendation, bestOption *interfaces.RecommendationOption) string {
	instanceID := extractInstanceID(*awsRec.InstanceArn)
	
	var description strings.Builder
	description.WriteString(fmt.Sprintf("AWS Compute Optimizer recommends optimizing EC2 instance %s", instanceID))
	
	if bestOption != nil && bestOption.ExpectedSavings > 0 {
		description.WriteString(fmt.Sprintf(" to save approximately $%.2f per month", bestOption.ExpectedSavings))
	}
	
	return description.String()
}

func (b *EC2RecommendationsBuilder) generateJustification(awsRec types.InstanceRecommendation, bestOption *interfaces.RecommendationOption) string {
	var justification strings.Builder
	
	justification.WriteString(fmt.Sprintf("AWS Compute Optimizer analyzed the workload patterns of instance %s", extractInstanceID(*awsRec.InstanceArn)))
	
	switch awsRec.Finding {
	case types.FindingOverProvisioned:
		justification.WriteString(" and determined it is over-provisioned based on CPU, memory, and network utilization patterns.")
	case types.FindingUnderProvisioned:
		justification.WriteString(" and determined it is under-provisioned based on resource utilization patterns.")
	case types.FindingOptimized:
		justification.WriteString(" and recommends upgrading to a newer instance generation for better price-performance.")
	}
	
	return justification.String()
}

// extractInstanceID extracts the instance ID from an ARN
func extractInstanceID(arn string) string {
	// Extract instance ID from ARN: arn:aws:ec2:region:account:instance/i-1234567890abcdef0
	parts := strings.Split(arn, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return arn
}