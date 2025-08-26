# Terraform configuration with various quality issues for TFLint testing
# This demonstrates common Terraform anti-patterns and best practices

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
  region = "us-east-1" # Hardcoded region - should be variable
}

# Missing variable descriptions
variable "instance_type" {
  type = string
}

variable "environment" {
  # Missing description and default
}

# Unused variable
variable "unused_var" {
  description = "This variable is never used"
  type        = string
  default     = "unused"
}

# Resource with deprecated argument
resource "aws_instance" "example" {
  ami           = "ami-12345678" # Hardcoded AMI - should use data source
  instance_type = var.instance_type

  # Deprecated argument (should use vpc_security_group_ids)
  security_groups = ["default"]

  # Missing required tags
  tags = {
    Name = "example-instance"
    # Missing Environment tag
  }
}

# Resource with inefficient naming
resource "aws_s3_bucket" "bucket" { # Vague naming
  bucket = "my-bucket-name" # Hardcoded bucket name
}

# Missing lifecycle rules
resource "aws_s3_bucket" "logs_bucket" {
  bucket = "application-logs-bucket"
  
  # Should have lifecycle configuration for log retention
}

# Resource without backup configuration
resource "aws_db_instance" "database" {
  identifier     = "myapp-db"
  engine         = "mysql"
  engine_version = "5.7" # Outdated version
  instance_class = "db.t2.micro" # Old generation instance
  
  allocated_storage = 20
  
  db_name  = "myapp"
  username = "root"
  password = "password123" # Should use random password
  
  # Missing backup configuration
  backup_retention_period = 0 # No backups
  
  skip_final_snapshot = true
}

# Missing data source for availability zones
resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-1a" # Hardcoded AZ
}

resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
  
  # Missing DNS settings
  enable_dns_hostnames = false
  enable_dns_support   = false
}

# Resource with potential naming conflicts
resource "aws_security_group" "sg" {
  name = "web-sg" # Should include environment/project prefix
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Output without description
output "instance_ip" {
  value = aws_instance.example.public_ip
}

# Output with sensitive data not marked
output "database_password" {
  value = aws_db_instance.database.password
  # Missing sensitive = true
}

# Local value that could be simplified
locals {
  common_tags = {
    Project     = "sandbox"
    Environment = var.environment != null ? var.environment : "development"
    ManagedBy   = "terraform"
  }
}