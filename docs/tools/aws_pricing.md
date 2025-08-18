# AWS Pricing - Official AWS CLI Pricing API

AWS Pricing tools for querying service pricing information using the official AWS Price List API through the AWS CLI.

## Overview

The AWS Price List API provides centralized access to AWS service pricing information. These MCP functions wrap the official AWS CLI pricing commands to enable programmatic access to pricing data for cost analysis, budgeting, and scenario planning.

## Available MCP Functions

### 1. `aws_pricing_describe_services`
**Description**: Get metadata for AWS services and their pricing attributes

**Parameters**:
- `service_code` (optional): AWS service code (e.g., AmazonEC2, AmazonS3) - leave empty to list all services
- `profile` (optional): AWS profile to use
- `max_items` (optional): Maximum number of items to return

**Example Usage**:
```bash
# List all available AWS services with pricing data
aws_pricing_describe_services()

# Get metadata for EC2 service
aws_pricing_describe_services(service_code="AmazonEC2")

# Get metadata for S3 service with item limit
aws_pricing_describe_services(service_code="AmazonS3", max_items="10")
```

### 2. `aws_pricing_get_attribute_values`
**Description**: Get available attribute values for AWS service pricing filters

**Parameters**:
- `service_code` (required): AWS service code (e.g., AmazonEC2, AmazonS3)
- `attribute_name` (required): Attribute name (e.g., instanceType, volumeType, location)
- `profile` (optional): AWS profile to use
- `max_items` (optional): Maximum number of items to return

**Example Usage**:
```bash
# Get all EC2 instance types
aws_pricing_get_attribute_values(
  service_code="AmazonEC2",
  attribute_name="instanceType"
)

# Get all EBS volume types with limit
aws_pricing_get_attribute_values(
  service_code="AmazonEC2",
  attribute_name="volumeType",
  max_items="20"
)

# Get all available regions
aws_pricing_get_attribute_values(
  service_code="AmazonEC2",
  attribute_name="location"
)
```

### 3. `aws_pricing_get_products`
**Description**: Get AWS pricing information for products that match filter criteria

**Parameters**:
- `service_code` (required): AWS service code (e.g., AmazonEC2, AmazonS3)
- `filters` (optional): JSON string of filter criteria
- `format_version` (optional): Format version for response (aws_v1)
- `profile` (optional): AWS profile to use
- `max_items` (optional): Maximum number of items to return

**Example Usage**:
```bash
# Get all EC2 products (limited results)
aws_pricing_get_products(
  service_code="AmazonEC2",
  max_items="5"
)

# Get EC2 pricing for specific instance type and region
aws_pricing_get_products(
  service_code="AmazonEC2",
  filters='[{"Type":"TERM_MATCH","Field":"instanceType","Value":"t3.micro"},{"Type":"TERM_MATCH","Field":"location","Value":"US East (N. Virginia)"}]',
  format_version="aws_v1",
  max_items="1"
)

# Get S3 storage pricing
aws_pricing_get_products(
  service_code="AmazonS3",
  filters='[{"Type":"TERM_MATCH","Field":"storageClass","Value":"General Purpose"}]'
)
```

### 4. `aws_pricing_get_version`
**Description**: Get AWS CLI version information

**Parameters**: None

**Example Usage**:
```bash
aws_pricing_get_version()
```

## Filter Examples

The `filters` parameter uses JSON format with the following structure:

### Filter Types
- **TERM_MATCH**: Exact match
- **EQUALS**: Exact equality
- **CONTAINS**: Contains the value
- **ANY_OF**: Matches any of the values
- **NONE_OF**: Matches none of the values

### Common Filter Examples

#### EC2 Instance Pricing
```json
[
  {
    "Type": "TERM_MATCH",
    "Field": "instanceType",
    "Value": "t3.medium"
  },
  {
    "Type": "TERM_MATCH",
    "Field": "location",
    "Value": "US East (N. Virginia)"
  },
  {
    "Type": "TERM_MATCH",
    "Field": "tenancy",
    "Value": "Shared"
  },
  {
    "Type": "TERM_MATCH",
    "Field": "operating-system",
    "Value": "Linux"
  }
]
```

#### EBS Volume Pricing
```json
[
  {
    "Type": "TERM_MATCH",
    "Field": "volumeType",
    "Value": "Provisioned IOPS"
  },
  {
    "Type": "TERM_MATCH",
    "Field": "location",
    "Value": "US West (Oregon)"
  }
]
```

#### S3 Storage Pricing
```json
[
  {
    "Type": "TERM_MATCH",
    "Field": "storageClass",
    "Value": "General Purpose"
  },
  {
    "Type": "TERM_MATCH",
    "Field": "location",
    "Value": "US East (N. Virginia)"
  }
]
```

## Real CLI Capabilities

All MCP functions use the official AWS CLI pricing commands:

### Describe Services
```bash
aws pricing describe-services --service-code AmazonEC2
```

### Get Attribute Values
```bash
aws pricing get-attribute-values --service-code AmazonEC2 --attribute-name instanceType --max-items 10
```

### Get Products
```bash
aws pricing get-products --service-code AmazonEC2 --filters file://filters.json --format-version aws_v1
```

### Version Information
```bash
aws --version
```

## Common AWS Service Codes

| Service | Service Code | Common Attributes |
|---------|--------------|-------------------|
| EC2 | AmazonEC2 | instanceType, location, operating-system, tenancy |
| S3 | AmazonS3 | storageClass, location |
| RDS | AmazonRDS | instanceType, location, databaseEngine |
| Lambda | AWSLambda | location |
| CloudFront | AmazonCloudFront | location |
| EBS | AmazonEC2 | volumeType, location |
| ELB | AWSELB | location, loadBalancerType |
| Route 53 | AmazonRoute53 | location |
| VPC | AmazonVPC | location |
| CloudWatch | AmazonCloudWatch | location |

## Prerequisites

### AWS CLI Installation
```bash
# Install AWS CLI v2
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Verify installation
aws --version
```

### AWS Configuration
```bash
# Configure AWS credentials (optional for pricing API)
aws configure

# Or use environment variables
export AWS_DEFAULT_REGION=us-east-1
```

**Note**: The AWS Pricing API does not require AWS credentials and can be accessed without an AWS account.

## Best Practices

### Efficient Queries
- **Use Filters**: Always filter results to get relevant data
- **Limit Results**: Use `max_items` to avoid large responses
- **Cache Results**: Price data changes infrequently, cache when possible
- **Pagination**: Handle paginated responses for large datasets

### Cost Analysis Workflow
1. **Discover Services**: Use `describe-services` to find available services
2. **Explore Attributes**: Use `get-attribute-values` to see available options
3. **Query Pricing**: Use `get-products` with specific filters
4. **Compare Options**: Query different configurations for comparison

### Filter Strategy
- **Start Broad**: Begin with service-level queries
- **Add Specificity**: Progressively add more specific filters
- **Regional Variations**: Always include location filters for accurate pricing
- **Multiple Scenarios**: Query different configurations for analysis

## Response Format

### Service Description Response
```json
{
  "Services": [
    {
      "ServiceCode": "AmazonEC2",
      "AttributeNames": [
        "instanceType",
        "location",
        "operating-system",
        "tenancy"
      ]
    }
  ]
}
```

### Attribute Values Response
```json
{
  "AttributeValues": [
    {
      "Value": "t3.micro"
    },
    {
      "Value": "t3.small"
    }
  ]
}
```

### Product Pricing Response
```json
{
  "PriceList": [
    {
      "product": {
        "productFamily": "Compute Instance",
        "attributes": {
          "instanceType": "t3.micro",
          "vcpu": "2",
          "memory": "1 GiB"
        }
      },
      "terms": {
        "OnDemand": {
          "priceDimensions": {
            "pricePerUnit": {
              "USD": "0.0104000000"
            }
          }
        }
      }
    }
  ]
}
```

## Use Cases

### Cost Estimation
- **Infrastructure Planning**: Estimate costs for new deployments
- **Migration Analysis**: Compare on-premises vs. cloud costs
- **Budgeting**: Create accurate budget forecasts
- **Optimization**: Find cost-effective instance types and storage

### Automation
- **Cost Monitoring**: Automated cost tracking and alerting
- **Right-sizing**: Programmatic instance type recommendations
- **Procurement**: Automated pricing for procurement processes
- **Reporting**: Generate cost reports and dashboards

### Comparison Analysis
- **Regional Pricing**: Compare costs across AWS regions
- **Service Options**: Compare different service configurations
- **Reserved vs. On-Demand**: Analyze different pricing models
- **Competitive Analysis**: Compare AWS pricing with other providers

## Troubleshooting

### Common Issues

1. **No Results Returned**
   - Verify service code spelling
   - Check filter values are exact matches
   - Ensure attribute names are correct

2. **Too Many Results**
   - Add more specific filters
   - Use `max_items` to limit results
   - Consider pagination for large datasets

3. **Invalid Filter Values**
   - Use `get-attribute-values` to see valid options
   - Check case sensitivity in filter values
   - Verify JSON syntax in filter strings

### Error Messages
- **InvalidParameterValue**: Check parameter values and format
- **Throttling**: Reduce request frequency
- **AccessDenied**: Usually not applicable to pricing API

## Integration with Ship CLI

These MCP functions integrate with Ship CLI's containerized execution:
- Pricing API calls are executed through the Ship CLI's Dagger engine
- No AWS credentials required for pricing queries
- Results can be processed and formatted by other Ship CLI tools

## References

- **AWS Pricing API User Guide**: https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/price-changes.html
- **AWS CLI Pricing Commands**: https://docs.aws.amazon.com/cli/latest/reference/pricing/
- **Price List API Examples**: https://docs.aws.amazon.com/cli/v1/userguide/cli_pricing_code_examples.html
- **AWS Service Codes**: https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/using-price-list-query-api.html