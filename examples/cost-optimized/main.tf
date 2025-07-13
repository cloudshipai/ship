terraform {
  required_version = ">= 1.0"
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

# Spot Fleet for cost optimization
resource "aws_spot_fleet_request" "web" {
  iam_fleet_role      = aws_iam_role.fleet.arn
  allocation_strategy = "diversified"
  target_capacity     = var.target_capacity
  spot_price          = var.spot_price

  launch_specification {
    image_id               = var.ami_id
    instance_type          = "t3.micro"
    subnet_id              = aws_subnet.public[0].id
    vpc_security_group_ids = [aws_security_group.web.id]
    
    user_data = base64encode(<<-EOF
      #!/bin/bash
      yum update -y
      yum install -y httpd
      systemctl start httpd
      systemctl enable httpd
      echo "<h1>Cost-Optimized Web Server</h1>" > /var/www/html/index.html
    EOF
    )
  }

  launch_specification {
    image_id               = var.ami_id
    instance_type          = "t3a.micro"  # AMD instances for cost savings
    subnet_id              = aws_subnet.public[1].id
    vpc_security_group_ids = [aws_security_group.web.id]
    
    user_data = base64encode(<<-EOF
      #!/bin/bash
      yum update -y
      yum install -y httpd
      systemctl start httpd
      systemctl enable httpd
      echo "<h1>Cost-Optimized Web Server</h1>" > /var/www/html/index.html
    EOF
    )
  }

  # COST ISSUE: Spot fleet not configured for termination
  terminate_instances_with_expiration = false

  tags = {
    Name        = "${var.project_name}-spot-fleet"
    Environment = var.environment
  }
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "${var.project_name}-vpc"
    Environment = var.environment
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "${var.project_name}-igw"
    Environment = var.environment
  }
}

# Public Subnets
resource "aws_subnet" "public" {
  count = 2

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet("10.0.0.0/16", 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name        = "${var.project_name}-public-${count.index + 1}"
    Environment = var.environment
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name        = "${var.project_name}-public-rt"
    Environment = var.environment
  }
}

resource "aws_route_table_association" "public" {
  count = length(aws_subnet.public)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Security Group
resource "aws_security_group" "web" {
  name        = "${var.project_name}-web-sg"
  description = "Security group for web servers"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.project_name}-web-sg"
    Environment = var.environment
  }
}

# IAM Role for Spot Fleet
resource "aws_iam_role" "fleet" {
  name = "${var.project_name}-fleet-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "spotfleet.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "fleet" {
  role       = aws_iam_role.fleet.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetTaggingRole"
}

# S3 Bucket with Intelligent Tiering for cost optimization
resource "aws_s3_bucket" "data" {
  bucket = "${var.project_name}-data-${random_string.bucket_suffix.result}"

  tags = {
    Name        = "${var.project_name}-data"
    Environment = var.environment
  }
}

resource "random_string" "bucket_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "aws_s3_bucket_intelligent_tiering_configuration" "data" {
  bucket = aws_s3_bucket.data.id
  name   = "entire-bucket"

  status = "Enabled"

  # COST OPTIMIZATION: Move to Deep Archive after 180 days
  optional_fields = ["BucketKeyStatus"]

  tiering {
    access_tier = "DEEP_ARCHIVE_ACCESS"
    days        = 180
  }

  tiering {
    access_tier = "ARCHIVE_ACCESS"
    days        = 90
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "data" {
  bucket = aws_s3_bucket.data.id

  rule {
    id     = "delete-old-versions"
    status = "Enabled"

    noncurrent_version_expiration {
      noncurrent_days = 30
    }

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
}

# COST ISSUE: RDS instance that could be optimized
resource "aws_db_instance" "main" {
  identifier                = "${var.project_name}-db"
  allocated_storage         = 100  # Could be over-provisioned
  storage_type             = "gp2"  # Could use gp3 for cost savings
  engine                   = "mysql"
  engine_version          = "8.0"
  instance_class          = "db.t3.medium"  # Could be t3.micro for dev
  db_name                 = "webapp"
  username                = "admin"
  password                = "changeme123!"
  backup_retention_period = 35  # Could be reduced for dev environments
  skip_final_snapshot     = true

  tags = {
    Name        = "${var.project_name}-db"
    Environment = var.environment
  }
}