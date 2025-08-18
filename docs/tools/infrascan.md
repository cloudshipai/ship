# Infrascan

Generate system maps by connecting to your AWS account.

## Description

Infrascan is a tool that generates comprehensive system maps by scanning your AWS infrastructure across multiple regions. It connects to your AWS account using existing AWS credentials and profiles to discover and map all resources, creating visual representations of your cloud infrastructure topology. The tool helps teams understand their AWS infrastructure layout, dependencies, and relationships between resources.

## MCP Tools

### Infrastructure Mapping
- **`infrascan_scan`** - Scan AWS infrastructure across regions and generate system map
- **`infrascan_graph`** - Generate graph from scan results
- **`infrascan_render`** - Render infrastructure graph for visualization

## Real CLI Commands Used

### Core Commands
- `infrascan scan` - Scan AWS infrastructure
- `infrascan scan --region <region>` - Scan specific AWS region
- `infrascan scan -o <output-dir>` - Specify output directory
- `infrascan graph -i <input-dir>` - Generate graph from scan results
- `infrascan render -i <graph-file>` - Render graph visualization
- `infrascan render --browser -i <graph-file>` - Open in browser

### Multi-Region Scanning
```bash
# Scan multiple regions
AWS_PROFILE=readonly infrascan scan --region us-east-1 --region eu-west-1 -o scan-output

# Basic scan with output directory
infrascan scan -o my-scan-results

# Scan all regions with specific profile
AWS_PROFILE=production infrascan scan --region us-east-1 --region us-west-2 --region eu-central-1 -o prod-scan
```

## Use Cases

### Infrastructure Discovery
- **Resource Inventory**: Discover all AWS resources across regions
- **Architecture Documentation**: Generate visual maps of infrastructure
- **Compliance Auditing**: Document infrastructure for compliance requirements
- **Cost Analysis**: Understand resource distribution for cost optimization

### Cloud Migration Planning
- **Migration Assessment**: Map existing infrastructure before migration
- **Dependency Mapping**: Understand resource relationships and dependencies
- **Architecture Planning**: Design new architectures based on current state
- **Risk Assessment**: Identify critical infrastructure components

### Security and Governance
- **Security Review**: Visualize infrastructure for security assessment
- **Access Control**: Understand resource access patterns
- **Network Topology**: Map network connections and security groups
- **Compliance Documentation**: Generate reports for auditing

### Team Collaboration
- **Knowledge Sharing**: Share infrastructure maps with team members
- **Onboarding**: Help new team members understand infrastructure
- **Change Planning**: Visualize impact of infrastructure changes
- **Troubleshooting**: Use maps to understand system relationships

## Configuration Examples

### Basic Infrastructure Scanning
```bash
# Scan current AWS account default region
infrascan scan -o scan-results

# Scan specific regions
infrascan scan --region us-east-1 --region us-west-2 -o multi-region-scan

# Scan with specific AWS profile
AWS_PROFILE=production infrascan scan --region us-east-1 -o prod-scan

# Scan development environment
AWS_PROFILE=dev infrascan scan --region us-west-1 -o dev-scan
```

### Graph Generation and Rendering
```bash
# Generate graph from scan results
infrascan graph -i scan-results

# Render graph in browser
infrascan render --browser -i scan-results/graph.json

# Render graph to file (requires additional tools)
infrascan render -i scan-results/graph.json > infrastructure-map.html
```

### AWS Profile Configuration
```bash
# Configure AWS profiles for different environments
aws configure --profile readonly
aws configure --profile production
aws configure --profile development

# Use specific profile for scanning
AWS_PROFILE=readonly infrascan scan --region us-east-1 -o readonly-scan
AWS_PROFILE=production infrascan scan --region us-east-1 --region eu-west-1 -o prod-scan
```

## Advanced Usage

### Comprehensive Multi-Region Scanning
```bash
#!/bin/bash
# comprehensive-scan.sh

PROFILE="readonly"
OUTPUT_DIR="comprehensive-scan-$(date +%Y%m%d)"
REGIONS=("us-east-1" "us-west-1" "us-west-2" "eu-west-1" "eu-central-1" "ap-southeast-1")

echo "Starting comprehensive AWS infrastructure scan..."
echo "Profile: $PROFILE"
echo "Output Directory: $OUTPUT_DIR"
echo "Regions: ${REGIONS[*]}"

# Build scan command
SCAN_CMD="AWS_PROFILE=$PROFILE infrascan scan -o $OUTPUT_DIR"
for region in "${REGIONS[@]}"; do
    SCAN_CMD="$SCAN_CMD --region $region"
done

# Execute scan
echo "Executing: $SCAN_CMD"
eval $SCAN_CMD

# Generate graph
echo "Generating infrastructure graph..."
infrascan graph -i $OUTPUT_DIR

# Open in browser
echo "Opening graph in browser..."
infrascan render --browser -i $OUTPUT_DIR/graph.json

echo "Scan complete! Results saved to: $OUTPUT_DIR"
```

### Automated Daily Scans
```bash
#!/bin/bash
# daily-infrastructure-scan.sh

DATE=$(date +%Y%m%d)
OUTPUT_BASE="/infrastructure-scans"
ENVIRONMENTS=("production" "staging" "development")

for env in "${ENVIRONMENTS[@]}"; do
    echo "Scanning $env environment..."
    
    OUTPUT_DIR="$OUTPUT_BASE/$env-$DATE"
    
    # Scan infrastructure
    AWS_PROFILE=$env infrascan scan \
        --region us-east-1 \
        --region us-west-2 \
        -o $OUTPUT_DIR
    
    # Generate graph
    infrascan graph -i $OUTPUT_DIR
    
    # Create documentation
    echo "# Infrastructure Scan Report - $env" > $OUTPUT_DIR/report.md
    echo "Date: $(date)" >> $OUTPUT_DIR/report.md
    echo "Environment: $env" >> $OUTPUT_DIR/report.md
    echo "Regions: us-east-1, us-west-2" >> $OUTPUT_DIR/report.md
    
    echo "Completed scan for $env environment"
done

echo "All environment scans completed!"
```

### CI/CD Integration
```yaml
# GitHub Actions
name: Infrastructure Mapping
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly on Monday at 2 AM
  workflow_dispatch:

jobs:
  infrastructure-scan:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        environment: [production, staging]
        
    steps:
    - uses: actions/checkout@v2
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v1
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-east-1
    
    - name: Install Infrascan
      run: npm install -g @infrascan/cli
    
    - name: Scan Infrastructure
      run: |
        mkdir -p scans/${{ matrix.environment }}
        infrascan scan \
          --region us-east-1 \
          --region us-west-2 \
          -o scans/${{ matrix.environment }}
    
    - name: Generate Graph
      run: infrascan graph -i scans/${{ matrix.environment }}
    
    - name: Upload Scan Results
      uses: actions/upload-artifact@v2
      with:
        name: infrastructure-scan-${{ matrix.environment }}
        path: scans/${{ matrix.environment }}
        retention-days: 30
```

### Infrastructure Monitoring Script
```bash
#!/bin/bash
# infrastructure-monitor.sh

SLACK_WEBHOOK="$1"
if [[ -z "$SLACK_WEBHOOK" ]]; then
    echo "Usage: $0 <slack-webhook-url>"
    exit 1
fi

# Scan current infrastructure
DATE=$(date +%Y%m%d-%H%M)
SCAN_DIR="monitoring-scan-$DATE"

echo "Scanning infrastructure for monitoring..."
AWS_PROFILE=monitoring infrascan scan \
    --region us-east-1 \
    --region us-west-2 \
    -o $SCAN_DIR

# Generate graph
infrascan graph -i $SCAN_DIR

# Count resources (example analysis)
if [[ -f "$SCAN_DIR/graph.json" ]]; then
    RESOURCE_COUNT=$(jq '.nodes | length' $SCAN_DIR/graph.json)
    
    # Send notification to Slack
    curl -X POST -H 'Content-type: application/json' \
        --data "{
            \"text\": \"Infrastructure Scan Complete\",
            \"attachments\": [{
                \"color\": \"good\",
                \"fields\": [{
                    \"title\": \"Resources Discovered\",
                    \"value\": \"$RESOURCE_COUNT\",
                    \"short\": true
                }, {
                    \"title\": \"Scan Time\",
                    \"value\": \"$(date)\",
                    \"short\": true
                }]
            }]
        }" \
        $SLACK_WEBHOOK
    
    echo "Infrastructure scan completed. $RESOURCE_COUNT resources discovered."
else
    echo "Error: Failed to generate graph.json"
    exit 1
fi
```

## Integration Patterns

### Terraform Integration
```bash
# terraform-infrascan.sh
#!/bin/bash

# Apply Terraform changes
echo "Applying Terraform configuration..."
terraform apply -auto-approve

if [[ $? -eq 0 ]]; then
    echo "Terraform apply successful. Scanning updated infrastructure..."
    
    # Wait for resources to be ready
    sleep 30
    
    # Scan updated infrastructure
    infrascan scan --region us-east-1 -o post-terraform-scan
    infrascan graph -i post-terraform-scan
    
    echo "Updated infrastructure map generated!"
    echo "View at: post-terraform-scan/graph.json"
else
    echo "Terraform apply failed. Skipping infrastructure scan."
    exit 1
fi
```

### Documentation Generation
```bash
# generate-docs.sh
#!/bin/bash

ENVIRONMENTS=("production" "staging" "development")
DOCS_DIR="docs/infrastructure"

mkdir -p $DOCS_DIR

for env in "${ENVIRONMENTS[@]}"; do
    echo "Generating documentation for $env..."
    
    # Scan infrastructure
    infrascan scan --region us-east-1 -o temp-scan-$env
    infrascan graph -i temp-scan-$env
    
    # Copy graph to docs
    cp temp-scan-$env/graph.json $DOCS_DIR/$env-infrastructure.json
    
    # Generate markdown documentation
    cat > $DOCS_DIR/$env-README.md << EOF
# $env Environment Infrastructure

Generated: $(date)

## Overview
This document contains the infrastructure mapping for the $env environment.

## Visualization
Open the infrastructure graph: [$env-infrastructure.json](./$env-infrastructure.json)

## Usage
Use the Infrascan render command to visualize:
\`\`\`bash
infrascan render --browser -i $env-infrastructure.json
\`\`\`

## Last Updated
$(date)
EOF
    
    # Cleanup
    rm -rf temp-scan-$env
    
    echo "Documentation generated for $env environment"
done

echo "All documentation generated in $DOCS_DIR/"
```

## Best Practices

### Scanning Strategy
- **Regular Scanning**: Schedule regular scans to track infrastructure changes
- **Environment Separation**: Use separate AWS profiles for different environments
- **Region Coverage**: Include all regions where resources are deployed
- **Output Organization**: Use date-stamped directories for scan results

### Security Considerations
- **Read-Only Access**: Use AWS profiles with read-only permissions for scanning
- **Credential Management**: Store AWS credentials securely
- **Access Control**: Limit access to scan results containing sensitive information
- **Profile Isolation**: Use separate AWS profiles for different environments

### Automation Integration
- **CI/CD Pipelines**: Integrate scanning into deployment pipelines
- **Monitoring**: Set up automated scans for infrastructure monitoring
- **Alerting**: Configure notifications for scan failures or anomalies
- **Documentation**: Automatically update documentation with scan results

### Data Management
- **Retention Policies**: Implement retention policies for scan results
- **Version Control**: Track infrastructure changes over time
- **Backup**: Store important scan results in secure locations
- **Cleanup**: Regularly clean up old scan data

## Error Handling

### Common Issues
```bash
# AWS credentials not configured
aws configure list
# Solution: Configure AWS credentials or set AWS_PROFILE

# Insufficient permissions
# Solution: Ensure IAM user/role has read permissions for AWS resources

# Network connectivity issues
# Solution: Check network connectivity to AWS APIs

# Invalid region specified
# Solution: Use valid AWS region names (us-east-1, eu-west-1, etc.)
```

### Troubleshooting
- **AWS Profile**: Verify AWS profile is configured correctly
- **Permissions**: Ensure sufficient read permissions for resource discovery
- **Network**: Check network connectivity to AWS services
- **Node.js**: Verify Node.js is installed for npm package installation

Infrascan provides essential visibility into AWS infrastructure through automated discovery and visual mapping, enabling teams to understand, document, and manage their cloud environments effectively.