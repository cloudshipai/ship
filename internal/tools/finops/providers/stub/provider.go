package stub

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
)

// StubProvider implements a stub provider for testing
type StubProvider struct {
	name         string
	vendorType   interfaces.VendorType
	capabilities interfaces.VendorCapabilities
}

// NewStubProvider creates a new stub provider
func NewStubProvider(vendorType interfaces.VendorType) *StubProvider {
	return &StubProvider{
		name:       string(vendorType),
		vendorType: vendorType,
		capabilities: interfaces.VendorCapabilities{
			SupportsRecommendations: true,
			SupportsCostForecasting: true,
			SupportsRightsizing:     true,
			ResourceTypes:          []string{"compute", "storage", "database"},
			Regions:               []string{"us-east-1", "us-west-2"},
		},
	}
}

func (p *StubProvider) Name() string {
	return p.name
}

func (p *StubProvider) Type() interfaces.VendorType {
	return p.vendorType
}

func (p *StubProvider) GetCapabilities() interfaces.VendorCapabilities {
	return p.capabilities
}

func (p *StubProvider) SetCredentials(creds interface{}) error {
	return nil // Stub implementation
}

func (p *StubProvider) ValidateConfig(config map[string]interface{}) error {
	return nil // Stub implementation
}

func (p *StubProvider) DiscoverResources(ctx context.Context, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
	// Return mock resources
	resources := []interfaces.Resource{
		{
			ID:     "i-1234567890abcdef0",
			Name:   "test-instance",
			Type:   interfaces.ResourceTypeCompute,
			Vendor: p.vendorType,
			Region: "us-east-1",
			Tags:   map[string]string{"Environment": "test"},
			Cost: &interfaces.CostData{
				HourlyRate:   0.10,
				DailyRate:    2.40,
				MonthlyRate:  72.00,
				Currency:     "USD",
				LastUpdated:  time.Now(),
			},
			Specifications: interfaces.ResourceSpecs{
				InstanceType: "t3.micro",
				CPUCount:     1,
				MemoryGB:     1.0,
				StorageGB:    8.0,
			},
			Configuration: map[string]interface{}{
				"instance_type": "t3.micro",
				"state":         "running",
			},
			LastUpdated: time.Now(),
		},
	}
	
	return resources, nil
}

func (p *StubProvider) GetCostData(ctx context.Context, opts interfaces.CostOptions) ([]interfaces.CostRecord, error) {
	// Return mock cost data
	now := time.Now()
	records := []interfaces.CostRecord{
		{
			ResourceID: "total",
			Period: interfaces.TimePeriod{
				Start: now.AddDate(0, 0, -30),
				End:   now,
			},
			Amount:   150.75,
			Currency: opts.Currency,
			Service:  "EC2-Instance",
			Region:   "us-east-1",
		},
	}
	
	return records, nil
}

func (p *StubProvider) GetRecommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
	// Return mock recommendations
	recommendations := []interfaces.Recommendation{
		{
			ID:                 "rec-1234567890abcdef0",
			ResourceID:         "i-1234567890abcdef0",
			ResourceARN:        fmt.Sprintf("arn:%s:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0", p.vendorType),
			ResourceType:       interfaces.ResourceTypeCompute,
			Provider:           p.vendorType,
			RecommendationType: "RightSizeInstance",
			Title:              "Right-size EC2 instance i-1234567890abcdef0",
			Description:        "This instance appears to be over-provisioned based on utilization patterns",
			EstimatedSavings: interfaces.EstimatedSavings{
				Currency: "USD",
				Monthly:  25.50,
			},
			PerformanceRisk: "Low",
			MigrationEffort: "Medium",
			Confidence:      0.85,
			Justification:   "CPU utilization averaged 5% over the last 30 days",
			Options: []interfaces.RecommendationOption{
				{
					OptionType:      "instance_type_change",
					Description:     "Change from t3.medium to t3.small",
					ExpectedSavings: 25.50,
					Parameters: map[string]interface{}{
						"current_type":     "t3.medium",
						"recommended_type": "t3.small",
					},
				},
			},
			CreatedAt: time.Now(),
		},
	}
	
	return recommendations, nil
}