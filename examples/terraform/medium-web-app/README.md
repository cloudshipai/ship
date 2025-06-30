# Medium Example: Web Application with Auto Scaling

This Terraform configuration creates a scalable web application infrastructure on AWS.

## Architecture

- **VPC**: Custom VPC with public and private subnets across 2 AZs
- **Load Balancer**: Application Load Balancer in public subnets
- **Auto Scaling**: EC2 instances in private subnets with auto scaling
- **NAT Gateways**: For outbound internet access from private subnets
- **Security**: Proper security groups limiting access

## Resources Created

- 1 VPC with 2 public and 2 private subnets
- 1 Internet Gateway
- 2 NAT Gateways (1 per AZ for high availability)
- 1 Application Load Balancer
- 1 Auto Scaling Group (2-6 instances)
- Security Groups for ALB and EC2 instances

## Testing with Ship CLI

```bash
# Navigate to this directory
cd examples/terraform/medium-web-app

# Run linting
ship terraform-tools lint

# Check for security issues
ship terraform-tools security-scan
ship terraform-tools checkov-scan

# Generate documentation
ship terraform-tools generate-docs

# Estimate costs
ship terraform-tools cost-estimate
ship terraform-tools cost-analysis

# Generate infrastructure diagram
ship terraform-tools generate-diagram . --hcl -o infrastructure.png
```

## Estimated Monthly Costs (us-east-1)

- **ALB**: ~$16/month + $0.008/LCU-hour
- **NAT Gateways**: 2 × ~$45/month = ~$90/month
- **EC2 Instances** (t3.micro): 2 × ~$7.50/month = ~$15/month
- **EBS Storage**: 2 × 8GB × $0.10/GB = ~$1.60/month
- **Data Transfer**: Variable based on usage

**Total Base Cost**: ~$125/month (excluding data transfer)

## Security Best Practices

- ✅ Private subnets for compute resources
- ✅ Security groups with least privilege
- ✅ NAT Gateways for secure outbound traffic
- ✅ Auto Scaling for availability
- ✅ Health checks configured

## Customization

You can modify the following variables:
- `instance_type`: Change EC2 instance size
- `min_size`, `max_size`: Adjust scaling limits
- `vpc_cidr`: Change network addressing
- `ami_id`: Use different AMI (update user_data accordingly)