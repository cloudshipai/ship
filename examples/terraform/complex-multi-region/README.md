# Complex Example: Multi-Region Production Architecture

This Terraform configuration demonstrates a complex, production-ready multi-region architecture on AWS.

## Architecture Overview

### Primary Region (us-east-1)
- **Networking**: VPC with public, private, and database subnets across 2 AZs
- **Compute**: Auto Scaling Group (3-10 instances) with Application Load Balancer
- **Database**: RDS MySQL Multi-AZ instance
- **Storage**: S3 bucket with cross-region replication

### Secondary Region (us-west-2) - Disaster Recovery
- **Networking**: Identical VPC setup for DR
- **Compute**: Smaller Auto Scaling Group (1-5 instances)
- **Database**: RDS MySQL single-AZ read replica
- **Storage**: S3 bucket receiving replicated data

### Global Components
- **Monitoring**: CloudWatch dashboards and alarms
- **Alerting**: SNS topics for operational alerts
- **Cost Management**: Cost anomaly detection
- **Security**: VPC Flow Logs, security groups, IAM roles

## Module Structure

```
.
├── main.tf                 # Root module orchestrating everything
├── variables.tf            # Input variables
├── outputs.tf             # Output values
└── modules/
    ├── networking/        # VPC, subnets, NAT gateways
    ├── compute/           # ALB, ASG, launch templates
    ├── database/          # RDS instances and security
    └── monitoring/        # CloudWatch, SNS, alarms
```

## Testing with Ship CLI

```bash
# Navigate to this directory
cd examples/terraform/complex-multi-region

# Run comprehensive analysis
ship terraform-tools lint
ship terraform-tools security-scan
ship terraform-tools checkov-scan
ship terraform-tools generate-docs
ship terraform-tools cost-estimate
ship terraform-tools cost-analysis

# Generate infrastructure diagram
ship terraform-tools generate-diagram . --hcl -o infrastructure.png

# Run all with automatic push to CloudShip
ship terraform-tools security-scan --push --push-tags "complex,multi-region"
ship terraform-tools cost-estimate --push --push-metadata "regions=2,environment=prod"
```

## Estimated Monthly Costs

### Primary Region (us-east-1)
- **VPC**: 2 NAT Gateways × $45 = $90
- **ALB**: $16 + usage
- **EC2**: 3 × t3.small × $15 = $45
- **RDS**: db.t3.micro Multi-AZ = $29
- **S3**: Storage + replication costs
- **CloudWatch**: Logs, metrics, dashboards = ~$10
- **Subtotal**: ~$190/month

### Secondary Region (us-west-2)
- **VPC**: 2 NAT Gateways × $45 = $90
- **ALB**: $16 + usage
- **EC2**: 1 × t3.small = $15
- **RDS**: db.t3.micro = $15
- **S3**: Storage costs
- **Subtotal**: ~$136/month

**Total Base Cost**: ~$326/month (excluding data transfer and storage)

## Security Features

- ✅ Network isolation with public/private/database subnets
- ✅ Security groups with least privilege access
- ✅ Encrypted RDS instances
- ✅ S3 bucket versioning and replication
- ✅ VPC Flow Logs for network monitoring
- ✅ IAM roles for EC2 instances
- ✅ SSM Parameter Store for secrets
- ✅ CloudWatch agent for detailed monitoring

## High Availability Features

- ✅ Multi-AZ deployment in primary region
- ✅ Auto Scaling based on CPU metrics
- ✅ Cross-region disaster recovery setup
- ✅ S3 cross-region replication
- ✅ Health checks and automatic recovery
- ✅ Multiple NAT gateways for redundancy

## Customization

Key variables to adjust:
- `primary_region` / `secondary_region`: Change regions
- `instance_type`: Scale compute resources
- `db_instance_class`: Scale database resources
- `min_size` / `max_size`: Adjust scaling limits
- `alert_email`: Set monitoring recipient

## Deployment Notes

1. This example creates real AWS resources that will incur costs
2. Ensure you have appropriate AWS credentials configured
3. The SNS email subscription requires confirmation
4. S3 bucket names must be globally unique
5. Some resources may take 10-15 minutes to create (RDS)

## Clean Up

To avoid ongoing charges:
```bash
terraform destroy
```

Note: Ensure S3 buckets are empty before destroying.