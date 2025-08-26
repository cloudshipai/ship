# Intentionally vulnerable Terraform configuration for security scanning demos
# This is for testing Trivy and Checkov security scanning capabilities

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
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

# Vulnerable S3 bucket - public read access
resource "aws_s3_bucket" "vulnerable_bucket" {
  bucket = "my-vulnerable-test-bucket-${random_string.suffix.result}"

  tags = {
    Environment = "test"
    Purpose     = "vulnerability-demo"
  }
}

resource "aws_s3_bucket_public_access_block" "vulnerable_bucket_pab" {
  bucket = aws_s3_bucket.vulnerable_bucket.id

  # These should be true for security, but left false for demo
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# Vulnerable EC2 instance - no encryption, open security group
resource "aws_instance" "vulnerable_instance" {
  ami           = "ami-0c02fb55956c7d316" # Amazon Linux 2
  instance_type = "t2.micro"
  
  # No encryption enabled - vulnerability
  ebs_block_device {
    device_name = "/dev/sda1"
    volume_size = 8
    encrypted   = false
  }

  vpc_security_group_ids = [aws_security_group.vulnerable_sg.id]

  tags = {
    Name = "vulnerable-instance"
  }
}

# Overly permissive security group
resource "aws_security_group" "vulnerable_sg" {
  name_prefix = "vulnerable-sg-"
  description = "Vulnerable security group for testing"

  # Open to the world - major vulnerability
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# RDS instance without encryption
resource "aws_db_instance" "vulnerable_db" {
  identifier     = "vulnerable-test-db"
  engine         = "mysql"
  engine_version = "8.0"
  instance_class = "db.t3.micro"
  
  allocated_storage = 20
  
  # No encryption - vulnerability
  storage_encrypted = false
  
  # Publicly accessible - vulnerability
  publicly_accessible = true
  
  db_name  = "testdb"
  username = "admin"
  password = "password123" # Hardcoded password - vulnerability
  
  skip_final_snapshot = true
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

output "bucket_name" {
  value = aws_s3_bucket.vulnerable_bucket.bucket
}

output "instance_id" {
  value = aws_instance.vulnerable_instance.id
}