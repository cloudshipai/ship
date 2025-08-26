# Simple Terraform configuration for cost analysis demos
# This will be used with OpenInfraQuote (no API key required)

terraform {
  required_version = ">= 0.14"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  description = "AWS region for cost analysis"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "demo"
}

# Multiple EC2 instances of different sizes for cost comparison
resource "aws_instance" "small_instance" {
  count         = 2
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t3.micro"

  tags = {
    Name        = "small-instance-${count.index + 1}"
    Environment = var.environment
    CostCenter  = "demo"
  }
}

resource "aws_instance" "medium_instance" {
  count         = 1
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t3.medium"

  ebs_block_device {
    device_name = "/dev/sda1"
    volume_size = 50
    volume_type = "gp3"
  }

  tags = {
    Name        = "medium-instance-${count.index + 1}"
    Environment = var.environment
    CostCenter  = "demo"
  }
}

resource "aws_instance" "large_instance" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t3.large"

  ebs_block_device {
    device_name = "/dev/sda1"
    volume_size = 100
    volume_type = "gp3"
    iops        = 3000
  }

  tags = {
    Name        = "large-instance"
    Environment = var.environment
    CostCenter  = "demo"
  }
}

# RDS instance for database costs
resource "aws_db_instance" "demo_db" {
  identifier     = "demo-cost-analysis-db"
  engine         = "mysql"
  engine_version = "8.0"
  instance_class = "db.t3.micro"
  
  allocated_storage     = 20
  max_allocated_storage = 100
  
  db_name  = "demodb"
  username = "admin"
  password = "temporarypassword"
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  skip_final_snapshot = true
  
  tags = {
    Name        = "demo-database"
    Environment = var.environment
    CostCenter  = "demo"
  }
}

# Load Balancer for additional cost analysis
resource "aws_lb" "demo_alb" {
  name               = "demo-alb"
  internal           = false
  load_balancer_type = "application"
  subnets            = data.aws_subnets.default.ids

  tags = {
    Name        = "demo-load-balancer"
    Environment = var.environment
    CostCenter  = "demo"
  }
}

# Data sources
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Outputs for cost tracking
output "total_instances" {
  value = length(aws_instance.small_instance) + length(aws_instance.medium_instance) + 1
}

output "instance_types" {
  value = {
    small  = aws_instance.small_instance[*].instance_type
    medium = aws_instance.medium_instance[*].instance_type
    large  = [aws_instance.large_instance.instance_type]
  }
}

output "database_instance_class" {
  value = aws_db_instance.demo_db.instance_class
}

output "load_balancer_type" {
  value = aws_lb.demo_alb.load_balancer_type
}