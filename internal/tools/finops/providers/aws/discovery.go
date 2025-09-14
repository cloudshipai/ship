package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/computeoptimizer/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
)

// discoverEC2Instances discovers EC2 instances in a region
func (p *AWSProvider) discoverEC2Instances(ctx context.Context, region string, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
	// Use the existing EC2 client (we'll handle region switching differently)
	client := p.ec2Client
	
	input := &ec2.DescribeInstancesInput{}
	
	// Add filters if specified
	var filters []ec2Types.Filter
	
	// Add tag filters
	for key, value := range opts.Tags {
		filters = append(filters, ec2Types.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", key)),
			Values: []string{value},
		})
	}
	
	// Filter by account ID if specified
	if opts.AccountID != "" {
		filters = append(filters, ec2Types.Filter{
			Name:   aws.String("owner-id"),
			Values: []string{opts.AccountID},
		})
	}
	
	if len(filters) > 0 {
		input.Filters = filters
	}
	
	resp, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}
	
	var resources []interfaces.Resource
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			resource, err := p.transformEC2Instance(instance, region)
			if err != nil {
				continue // Skip instances that can't be transformed
			}
			resources = append(resources, resource)
		}
	}
	
	return resources, nil
}

// discoverEBSVolumes discovers EBS volumes in a region
func (p *AWSProvider) discoverEBSVolumes(ctx context.Context, region string, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
	// Use the existing EC2 client
	client := p.ec2Client
	
	input := &ec2.DescribeVolumesInput{}
	
	// Add filters if specified
	var filters []ec2Types.Filter
	
	// Add tag filters
	for key, value := range opts.Tags {
		filters = append(filters, ec2Types.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", key)),
			Values: []string{value},
		})
	}
	
	// Filter by account ID if specified  
	if opts.AccountID != "" {
		filters = append(filters, ec2Types.Filter{
			Name:   aws.String("owner-id"),
			Values: []string{opts.AccountID},
		})
	}
	
	if len(filters) > 0 {
		input.Filters = filters
	}
	
	resp, err := client.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe volumes: %w", err)
	}
	
	var resources []interfaces.Resource
	for _, volume := range resp.Volumes {
		resource, err := p.transformEBSVolume(volume, region)
		if err != nil {
			continue // Skip volumes that can't be transformed
		}
		resources = append(resources, resource)
	}
	
	return resources, nil
}

// discoverRDSInstances discovers RDS instances in a region
func (p *AWSProvider) discoverRDSInstances(ctx context.Context, region string, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
	// TODO: Implement RDS discovery - requires RDS client
	// For now, return empty slice
	return []interfaces.Resource{}, nil
}

// transformEC2Instance converts AWS EC2 instance to our standard resource format
func (p *AWSProvider) transformEC2Instance(instance ec2Types.Instance, region string) (interfaces.Resource, error) {
	if instance.InstanceId == nil {
		return interfaces.Resource{}, fmt.Errorf("instance ID is nil")
	}
	
	// Extract tags
	tags := make(map[string]string)
	for _, tag := range instance.Tags {
		if tag.Key != nil && tag.Value != nil {
			tags[*tag.Key] = *tag.Value
		}
	}
	
	// Determine instance name
	name := *instance.InstanceId
	if nameTag, exists := tags["Name"]; exists {
		name = nameTag
	}
	
	// Build resource specifications
	specs := interfaces.ResourceSpecs{
		InstanceType: string(instance.InstanceType),
	}
	
	// Extract CPU and memory information if available
	// Note: This would typically require additional API calls or lookups
	// for full instance type specifications
	
	// Extract storage information
	var storageGB float64
	for _, bdm := range instance.BlockDeviceMappings {
		if bdm.Ebs != nil && bdm.Ebs.VolumeId != nil {
			// Note: VolumeSize is not available in EbsInstanceBlockDevice
			// Would need to call DescribeVolumes to get actual volume size
			// For now, we'll use a placeholder or skip this
			storageGB = 0 // TODO: Call DescribeVolumes to get actual volume sizes
		}
	}
	specs.StorageGB = storageGB
	
	// Build resource
	resource := interfaces.Resource{
		ID:            *instance.InstanceId,
		Name:          name,
		Type:          interfaces.ResourceTypeCompute,
		Vendor:        interfaces.VendorAWS,
		Region:        region,
		Tags:          tags,
		Specifications: specs,
		Configuration: map[string]interface{}{
			"instance_type":          string(instance.InstanceType),
			"state":                  string(instance.State.Name),
			"platform":               p.extractPlatform(instance.Platform),
			"vpc_id":                 p.extractStringPointer(instance.VpcId),
			"subnet_id":              p.extractStringPointer(instance.SubnetId),
			"availability_zone":      p.extractStringPointer(instance.Placement.AvailabilityZone),
			"launch_time":            instance.LaunchTime,
			"instance_lifecycle":     string(instance.InstanceLifecycle),
			"architecture":           string(instance.Architecture),
			"virtualization_type":    string(instance.VirtualizationType),
		},
		LastUpdated: time.Now(),
	}
	
	// Note: Account ID is not available directly from instance object
	// Would need to be passed from the calling context or derived from IAM/STS
	// For now, we'll leave it empty or use a placeholder
	// TODO: Get account ID from AWS STS GetCallerIdentity or pass from context
	
	return resource, nil
}

// transformEBSVolume converts AWS EBS volume to our standard resource format
func (p *AWSProvider) transformEBSVolume(volume ec2Types.Volume, region string) (interfaces.Resource, error) {
	if volume.VolumeId == nil {
		return interfaces.Resource{}, fmt.Errorf("volume ID is nil")
	}
	
	// Extract tags
	tags := make(map[string]string)
	for _, tag := range volume.Tags {
		if tag.Key != nil && tag.Value != nil {
			tags[*tag.Key] = *tag.Value
		}
	}
	
	// Determine volume name
	name := *volume.VolumeId
	if nameTag, exists := tags["Name"]; exists {
		name = nameTag
	}
	
	// Build resource specifications
	specs := interfaces.ResourceSpecs{
		StorageType: string(volume.VolumeType),
	}
	
	if volume.Size != nil {
		specs.StorageGB = float64(*volume.Size)
	}
	
	// Build resource
	resource := interfaces.Resource{
		ID:            *volume.VolumeId,
		Name:          name,
		Type:          interfaces.ResourceTypeStorage,
		Vendor:        interfaces.VendorAWS,
		Region:        region,
		Tags:          tags,
		Specifications: specs,
		Configuration: map[string]interface{}{
			"volume_type":       string(volume.VolumeType),
			"size":              volume.Size,
			"state":             string(volume.State),
			"encrypted":         volume.Encrypted,
			"availability_zone": p.extractStringPointer(volume.AvailabilityZone),
			"create_time":       volume.CreateTime,
		},
		LastUpdated: time.Now(),
	}
	
	// Add IOPS information if available
	if volume.Iops != nil {
		resource.Configuration["iops"] = *volume.Iops
	}
	
	// Add throughput information if available  
	if volume.Throughput != nil {
		resource.Configuration["throughput"] = *volume.Throughput
	}
	
	// Check if volume is attached
	if len(volume.Attachments) > 0 {
		resource.Configuration["attachments"] = volume.Attachments
		resource.Configuration["attached"] = true
	} else {
		resource.Configuration["attached"] = false
	}
	
	return resource, nil
}

// Helper functions

func (p *AWSProvider) extractPlatform(platform ec2Types.PlatformValues) string {
	if platform == "" {
		return "linux"
	}
	return string(platform)
}

func (p *AWSProvider) extractStringPointer(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// getEC2Recommendations implements the EC2-specific recommendation logic
func (p *AWSProvider) getEC2Recommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
	if p.computeOptimizerClient == nil {
		return nil, fmt.Errorf("compute optimizer client not initialized")
	}
	
	var allRecommendations []interfaces.Recommendation
	
	// Determine target regions
	regions := p.config.Regions
	if len(opts.Regions) > 0 {
		regions = opts.Regions
	}
	
	// Get recommendations for each finding type
	findings := []types.Finding{
		types.FindingOverProvisioned,
		types.FindingUnderProvisioned,
		types.FindingOptimized,
	}
	
	for _, finding := range findings {
		builder := NewEC2RecommendationsBuilder(p.computeOptimizerClient, finding)
		
		for _, region := range regions {
			recommendations, err := builder.GetRecommendations(ctx, region, opts.ARNs)
			if err != nil {
				return nil, fmt.Errorf("failed to get recommendations for region %s: %w", region, err)
			}
			allRecommendations = append(allRecommendations, recommendations...)
		}
	}
	
	return allRecommendations, nil
}

// getReservedInstanceRecommendations gets reserved instance recommendations
func (p *AWSProvider) getReservedInstanceRecommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
	if p.costExplorerClient == nil {
		return nil, fmt.Errorf("cost explorer client not initialized")
	}
	
	// TODO: Implement reserved instance recommendations using Cost Explorer
	// This would use GetReservationPurchaseRecommendation API
	return []interfaces.Recommendation{}, nil
}

// getSpotInstanceRecommendations gets spot instance recommendations  
func (p *AWSProvider) getSpotInstanceRecommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
	// TODO: Implement spot instance recommendations
	// This would analyze current on-demand instances and suggest spot alternatives
	return []interfaces.Recommendation{}, nil
}