package aws
import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/computeoptimizer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
)

// AWSProvider implements the VendorProvider interface for Amazon Web Services
// Based on OpenOps patterns for multi-account, region-based discovery
type AWSProvider struct {
	// AWS service clients
	computeOptimizerClient *computeoptimizer.Client
	costExplorerClient     *costexplorer.Client
	pricingClient         *pricing.Client
	ec2Client             *ec2.Client
	
	// Configuration
	config      AWSConfig
	credentials AWSCredentials
}

// AWSConfig contains AWS-specific configuration
type AWSConfig struct {
	Region           string            `json:"region"`
	Regions          []string          `json:"regions"`
	AccountIDs       []string          `json:"account_ids,omitempty"`
	RoleARN          string            `json:"role_arn,omitempty"`
	ExternalID       string            `json:"external_id,omitempty"`
	SessionName      string            `json:"session_name,omitempty"`
	Tags             map[string]string `json:"default_tags,omitempty"`
	
	// Service-specific configuration
	ComputeOptimizer ComputeOptimizerConfig `json:"compute_optimizer"`
	CostExplorer     CostExplorerConfig     `json:"cost_explorer"`
}

// ComputeOptimizerConfig configures AWS Compute Optimizer integration
type ComputeOptimizerConfig struct {
	Enabled      bool     `json:"enabled"`
	Findings     []string `json:"findings"`      // OVER_PROVISIONED, UNDER_PROVISIONED, OPTIMIZED
	ResourceTypes []string `json:"resource_types"` // EC2, EBS, Lambda, etc.
}

// CostExplorerConfig configures AWS Cost Explorer integration  
type CostExplorerConfig struct {
	Enabled bool     `json:"enabled"`
	Metrics []string `json:"metrics"` // BlendedCost, UnblendedCost, etc.
}

// AWSCredentials contains AWS authentication information
type AWSCredentials struct {
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	SessionToken    string `json:"session_token,omitempty"`
	Profile         string `json:"profile,omitempty"`
	RoleARN         string `json:"role_arn,omitempty"`
}

// NewProvider creates a new AWS provider with default configuration  
func NewProvider() (interfaces.VendorProvider, error) {
	// Default configuration
	cfg := AWSConfig{
		Region: "us-east-1", // Default region
		ComputeOptimizer: ComputeOptimizerConfig{
			Enabled:       true,
			Findings:      []string{"OVER_PROVISIONED", "UNDER_PROVISIONED", "OPTIMIZED"},
			ResourceTypes: []string{"EC2", "EBS"},
		},
		CostExplorer: CostExplorerConfig{
			Enabled: true,
			Metrics: []string{"BlendedCost"},
		},
	}
	
	return NewAWSProvider(cfg)
}

// NewAWSProvider creates a new AWS provider instance
func NewAWSProvider(cfg AWSConfig) (*AWSProvider, error) {
	provider := &AWSProvider{
		config: cfg,
	}
	
	// Initialize AWS clients
	if err := provider.initializeClients(); err != nil {
		return nil, fmt.Errorf("failed to initialize AWS clients: %w", err)
	}
	
	return provider, nil
}

// initializeClients sets up AWS service clients based on configuration
func (p *AWSProvider) initializeClients() error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(p.config.Region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	
	// Initialize service clients
	p.computeOptimizerClient = computeoptimizer.NewFromConfig(cfg)
	p.costExplorerClient = costexplorer.NewFromConfig(cfg)
	p.pricingClient = pricing.NewFromConfig(cfg)
	p.ec2Client = ec2.NewFromConfig(cfg)
	
	return nil
}

// Name returns the provider name
func (p *AWSProvider) Name() string {
	return "AWS"
}

// Type returns the vendor type
func (p *AWSProvider) Type() interfaces.VendorType {
	return interfaces.VendorAWS
}

// GetCapabilities returns AWS provider capabilities
func (p *AWSProvider) GetCapabilities() interfaces.VendorCapabilities {
	return interfaces.VendorCapabilities{
		SupportsRecommendations: true,  // AWS Compute Optimizer
		SupportsCostForecasting: true,  // AWS Cost Explorer forecasting
		SupportsRightsizing:     true,
		ResourceTypes: []string{
			"ec2-instance", "rds-database", "ebs-volume",
			"lambda-function", "s3-bucket", "elasticache-cluster",
		},
		Regions: []string{
			"us-east-1", "us-west-2", "us-west-1", "us-east-2",
			"eu-west-1", "eu-central-1", "ap-southeast-1", "ap-northeast-1",
		},
	}
}

// SetCredentials configures AWS authentication credentials
func (p *AWSProvider) SetCredentials(creds interface{}) error {
	awsCreds, ok := creds.(AWSCredentials)
	if !ok {
		return fmt.Errorf("invalid credentials type for AWS provider")
	}
	
	p.credentials = awsCreds
	
	// Re-initialize clients with new credentials
	return p.initializeClients()
}

// ValidateConfig validates AWS-specific configuration
func (p *AWSProvider) ValidateConfig(config map[string]interface{}) error {
	// Check required fields
	if p.config.Region == "" {
		return fmt.Errorf("region is required for AWS provider")
	}
	
	// Validate regions
	validRegions := map[string]bool{
		"us-east-1": true, "us-west-2": true, "us-west-1": true, "us-east-2": true,
		"eu-west-1": true, "eu-central-1": true, "ap-southeast-1": true, "ap-northeast-1": true,
	}
	
	if !validRegions[p.config.Region] {
		return fmt.Errorf("invalid AWS region: %s", p.config.Region)
	}
	
	return nil
}

// DiscoverResources discovers AWS resources based on OpenOps patterns
func (p *AWSProvider) DiscoverResources(ctx context.Context, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
	var resources []interfaces.Resource
	
	// Multi-region discovery if regions specified
	regions := p.config.Regions
	if opts.Region != "" {
		regions = []string{opts.Region}
	}
	
	for _, region := range regions {
		// Create region-specific clients
		regionResources, err := p.discoverResourcesInRegion(ctx, region, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to discover resources in region %s: %w", region, err)
		}
		resources = append(resources, regionResources...)
	}
	
	return resources, nil
}

// discoverResourcesInRegion discovers resources in a specific AWS region
func (p *AWSProvider) discoverResourcesInRegion(ctx context.Context, region string, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
	var resources []interfaces.Resource
	
	// Filter by resource types or default to all supported types
	resourceTypes := opts.ResourceTypes
	if len(resourceTypes) == 0 {
		resourceTypes = []string{"compute", "storage", "database"}
	}
	
	for _, resourceType := range resourceTypes {
		switch resourceType {
		case "compute":
			ec2Resources, err := p.discoverEC2Instances(ctx, region, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to discover EC2 instances: %w", err)
			}
			resources = append(resources, ec2Resources...)
			
		case "storage":
			ebsResources, err := p.discoverEBSVolumes(ctx, region, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to discover EBS volumes: %w", err)
			}
			resources = append(resources, ebsResources...)
			
		case "database":
			rdsResources, err := p.discoverRDSInstances(ctx, region, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to discover RDS instances: %w", err)
			}
			resources = append(resources, rdsResources...)
		}
	}
	
	// Apply ARN filtering if specified (OpenOps pattern)
	if len(opts.ARNs) > 0 {
		resources = p.filterResourcesByARNs(resources, opts.ARNs)
	}
	
	// Apply tag filtering if specified
	if len(opts.Tags) > 0 {
		resources = p.filterResourcesByTags(resources, opts.Tags)
	}
	
	return resources, nil
}

// GetCostData retrieves cost data from AWS Cost Explorer
func (p *AWSProvider) GetCostData(ctx context.Context, opts interfaces.CostOptions) ([]interfaces.CostRecord, error) {
	costExplorer := NewCostExplorerClient(p.costExplorerClient)
	return costExplorer.GetCostData(ctx, opts)
}

// GetRecommendations retrieves optimization recommendations from AWS services
func (p *AWSProvider) GetRecommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
	var recommendations []interfaces.Recommendation
	
	for _, findingType := range opts.FindingTypes {
		switch findingType {
		case "rightsizing":
			// Follow OpenOps EC2 recommendations pattern
			recs, err := p.getEC2Recommendations(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to get EC2 recommendations: %w", err)
			}
			recommendations = append(recommendations, recs...)
			
		case "reserved_instances":
			// Follow OpenOps reserved instance patterns
			recs, err := p.getReservedInstanceRecommendations(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to get reserved instance recommendations: %w", err)
			}
			recommendations = append(recommendations, recs...)
			
		case "spot_instances":
			recs, err := p.getSpotInstanceRecommendations(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to get spot instance recommendations: %w", err)
			}
			recommendations = append(recommendations, recs...)
		}
	}
	
	return recommendations, nil
}

// Helper methods for resource discovery and filtering
func (p *AWSProvider) filterResourcesByARNs(resources []interfaces.Resource, arns []string) []interfaces.Resource {
	arnSet := make(map[string]bool)
	for _, arn := range arns {
		arnSet[arn] = true
	}
	
	var filtered []interfaces.Resource
	for _, resource := range resources {
		if arnSet[resource.ID] { // Assuming ID contains ARN
			filtered = append(filtered, resource)
		}
	}
	
	return filtered
}

func (p *AWSProvider) filterResourcesByTags(resources []interfaces.Resource, tags map[string]string) []interfaces.Resource {
	var filtered []interfaces.Resource
	for _, resource := range resources {
		matches := true
		for key, value := range tags {
			if resource.Tags[key] != value {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, resource)
		}
	}
	
	return filtered
}