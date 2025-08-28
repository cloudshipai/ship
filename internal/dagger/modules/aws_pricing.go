package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// AWSPricingModule provides AWS pricing information using AWS CLI
type AWSPricingModule struct {
	client *dagger.Client
	name   string
}

const awsPricingBinary = "aws"

// NewAWSPricingModule creates a new AWS pricing module
func NewAWSPricingModule(client *dagger.Client) *AWSPricingModule {
	return &AWSPricingModule{
		client: client,
		name:   "aws-pricing",
	}
}

// GetServicePricing gets pricing for a specific AWS service
func (m *AWSPricingModule) GetServicePricing(ctx context.Context, service, region string) (string, error) {
	if service == "" {
		service = "EC2-Instance"
	}
	if region == "" {
		region = "us-east-1"
	}

	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{awsPricingBinary, "pricing", "describe-services", 
			"--service-code", service,
			"--region", "us-east-1", // Pricing API is only available in us-east-1
			"--format-version", "aws_v1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get service pricing: %w", err)
	}

	return output, nil
}

// GetEC2Pricing gets EC2 instance pricing for a specific region and instance type
func (m *AWSPricingModule) GetEC2Pricing(ctx context.Context, instanceType, region string) (string, error) {
	if instanceType == "" {
		instanceType = "t3.micro"
	}
	if region == "" {
		region = "us-east-1"
	}

	// Use AWS CLI to get EC2 pricing
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{awsPricingBinary, "pricing", "get-products",
			"--service-code", "AmazonEC2",
			"--region", "us-east-1", // Pricing API is only available in us-east-1
			"--filters",
			fmt.Sprintf("Type=TERM_MATCH,Field=instanceType,Value=%s", instanceType),
			fmt.Sprintf("Type=TERM_MATCH,Field=location,Value=%s", getLocationFromRegion(region)),
			"Type=TERM_MATCH,Field=tenancy,Value=Shared",
			"Type=TERM_MATCH,Field=operating-system,Value=Linux",
			"Type=TERM_MATCH,Field=preInstalledSw,Value=NA",
			"--format-version", "aws_v1",
			"--max-results", "1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get EC2 pricing: %w", err)
	}

	return output, nil
}

// GetRDSPricing gets RDS pricing for a specific instance class and region
func (m *AWSPricingModule) GetRDSPricing(ctx context.Context, instanceClass, engine, region string) (string, error) {
	if instanceClass == "" {
		instanceClass = "db.t3.micro"
	}
	if engine == "" {
		engine = "MySQL"
	}
	if region == "" {
		region = "us-east-1"
	}

	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{awsPricingBinary, "pricing", "get-products",
			"--service-code", "AmazonRDS",
			"--region", "us-east-1", // Pricing API is only available in us-east-1
			"--filters",
			fmt.Sprintf("Type=TERM_MATCH,Field=instanceType,Value=%s", instanceClass),
			fmt.Sprintf("Type=TERM_MATCH,Field=location,Value=%s", getLocationFromRegion(region)),
			fmt.Sprintf("Type=TERM_MATCH,Field=databaseEngine,Value=%s", engine),
			"Type=TERM_MATCH,Field=deploymentOption,Value=Single-AZ",
			"--format-version", "aws_v1",
			"--max-results", "1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get RDS pricing: %w", err)
	}

	return output, nil
}

// ListServices lists available AWS services for pricing
func (m *AWSPricingModule) ListServices(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{awsPricingBinary, "pricing", "describe-services",
			"--region", "us-east-1", // Pricing API is only available in us-east-1
			"--format-version", "aws_v1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list services: %w", err)
	}

	return output, nil
}

// CalculateMonthlyCost calculates estimated monthly cost for common resources
func (m *AWSPricingModule) CalculateMonthlyCost(ctx context.Context, resourceType, size, region string) (string, error) {
	if region == "" {
		region = "us-east-1"
	}

	// Create a simple cost calculator script
	script := fmt.Sprintf(`#!/bin/bash
echo "=== AWS Monthly Cost Calculator ==="
echo "Resource Type: %s"
echo "Size/Type: %s" 
echo "Region: %s"
echo ""

case "%s" in
	"ec2")
		echo "EC2 Instance Pricing Estimate:"
		echo "Instance Type: %s"
		echo "Region: %s"
		echo ""
		# Get basic EC2 pricing info
		aws pricing get-products \
			--service-code AmazonEC2 \
			--region us-east-1 \
			--filters \
			"Type=TERM_MATCH,Field=instanceType,Value=%s" \
			"Type=TERM_MATCH,Field=location,Value=%s" \
			"Type=TERM_MATCH,Field=tenancy,Value=Shared" \
			"Type=TERM_MATCH,Field=operating-system,Value=Linux" \
			--format-version aws_v1 \
			--max-results 1 | jq -r '.PriceList[0]' | jq -r '.terms.OnDemand | to_entries[0].value.priceDimensions | to_entries[0].value.pricePerUnit.USD' 2>/dev/null || echo "Unable to fetch exact pricing"
		echo ""
		echo "Note: This is the hourly rate. Monthly estimate = hourly rate × 730 hours"
		;;
	"rds")
		echo "RDS Database Pricing Estimate:"
		echo "Instance Class: %s"
		echo "Region: %s"
		echo ""
		echo "Estimated monthly cost varies by instance type and storage."
		echo "Common estimates:"
		echo "- db.t3.micro: ~$15-20/month"
		echo "- db.t3.small: ~$30-40/month"
		echo "- db.t3.medium: ~$60-80/month"
		;;
	"s3")
		echo "S3 Storage Pricing Estimate:"
		echo "Storage Class: Standard"
		echo "Region: %s"
		echo ""
		echo "S3 Standard Storage: ~$0.023 per GB/month"
		echo "Example costs:"
		echo "- 100 GB: ~$2.30/month"
		echo "- 1 TB: ~$23/month"
		echo "- 10 TB: ~$230/month"
		;;
	*)
		echo "Resource type '%s' not supported yet."
		echo "Supported types: ec2, rds, s3"
		;;
esac
`, resourceType, size, region, resourceType, size, region, size, getLocationFromRegion(region), size, region, region, resourceType)

	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{"apk", "add", "--no-cache", "jq"}).
		WithNewFile("/tmp/calculate.sh", script, dagger.ContainerWithNewFileOpts{
			Permissions: 0755,
		}).
		WithExec([]string{"/tmp/calculate.sh"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to calculate costs: %w", err)
	}

	return output, nil
}

// DescribeServices gets metadata for AWS services and their pricing attributes
func (m *AWSPricingModule) DescribeServices(ctx context.Context, serviceCode string, maxItems string) (string, error) {
	args := []string{"aws", "pricing", "describe-services", "--region", "us-east-1"}
	
	if serviceCode != "" {
		args = append(args, "--service-code", serviceCode)
	}
	if maxItems != "" {
		args = append(args, "--max-items", maxItems)
	}

	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe services: %w", err)
	}

	return output, nil
}

// GetProducts gets AWS pricing information for products that match filter criteria
func (m *AWSPricingModule) GetProducts(ctx context.Context, serviceCode, filters, formatVersion, maxItems string) (string, error) {
	args := []string{"aws", "pricing", "get-products", "--service-code", serviceCode, "--region", "us-east-1"}
	
	if filters != "" {
		args = append(args, "--filters", filters)
	}
	if formatVersion != "" {
		args = append(args, "--format-version", formatVersion)
	}
	if maxItems != "" {
		args = append(args, "--max-items", maxItems)
	}

	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get products: %w", err)
	}

	return output, nil
}

// GetVersion returns the AWS CLI version
func (m *AWSPricingModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{awsPricingBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get AWS CLI version: %w", err)
	}

	return output, nil
}

// getLocationFromRegion converts AWS region to pricing API location format
func getLocationFromRegion(region string) string {
	regionMap := map[string]string{
		"us-east-1":      "US East (N. Virginia)",
		"us-east-2":      "US East (Ohio)",
		"us-west-1":      "US West (N. California)",
		"us-west-2":      "US West (Oregon)",
		"eu-west-1":      "Europe (Ireland)",
		"eu-west-2":      "Europe (London)",
		"eu-west-3":      "Europe (Paris)",
		"eu-central-1":   "Europe (Frankfurt)",
		"eu-north-1":     "Europe (Stockholm)",
		"ap-southeast-1": "Asia Pacific (Singapore)",
		"ap-southeast-2": "Asia Pacific (Sydney)",
		"ap-northeast-1": "Asia Pacific (Tokyo)",
		"ap-northeast-2": "Asia Pacific (Seoul)",
		"ap-south-1":     "Asia Pacific (Mumbai)",
		"ca-central-1":   "Canada (Central)",
		"sa-east-1":      "South America (Sao Paulo)",
	}

	if location, exists := regionMap[region]; exists {
		return location
	}
	return "US East (N. Virginia)" // Default fallback
}

// GetAttributeValues gets available attribute values for AWS service pricing filters
func (m *AWSPricingModule) GetAttributeValues(ctx context.Context, serviceCode, attributeName, maxItems string) (string, error) {
	args := []string{"aws", "pricing", "get-attribute-values", "--service-code", serviceCode, "--attribute-name", attributeName, "--region", "us-east-1"}
	
	if maxItems != "" {
		args = append(args, "--max-items", maxItems)
	}

	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get attribute values: %w", err)
	}

	return output, nil
}
