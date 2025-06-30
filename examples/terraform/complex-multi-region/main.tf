terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Primary region provider
provider "aws" {
  region = var.primary_region
  alias  = "primary"
}

# Secondary region provider
provider "aws" {
  region = var.secondary_region
  alias  = "secondary"
}

# Primary region infrastructure
module "primary_networking" {
  source = "./modules/networking"
  providers = {
    aws = aws.primary
  }

  region       = var.primary_region
  vpc_cidr     = var.primary_vpc_cidr
  project_name = var.project_name
  environment  = var.environment
}

module "primary_compute" {
  source = "./modules/compute"
  providers = {
    aws = aws.primary
  }

  vpc_id              = module.primary_networking.vpc_id
  public_subnet_ids   = module.primary_networking.public_subnet_ids
  private_subnet_ids  = module.primary_networking.private_subnet_ids
  project_name        = var.project_name
  environment         = var.environment
  instance_type       = var.instance_type
  min_size           = var.min_size
  max_size           = var.max_size
  database_endpoint   = module.primary_database.endpoint
}

module "primary_database" {
  source = "./modules/database"
  providers = {
    aws = aws.primary
  }

  vpc_id             = module.primary_networking.vpc_id
  subnet_ids         = module.primary_networking.database_subnet_ids
  project_name       = var.project_name
  environment        = var.environment
  instance_class     = var.db_instance_class
  allocated_storage  = var.db_allocated_storage
  multi_az          = true
  security_group_ids = [module.primary_compute.app_security_group_id]
}

# Secondary region infrastructure (disaster recovery)
module "secondary_networking" {
  source = "./modules/networking"
  providers = {
    aws = aws.secondary
  }

  region       = var.secondary_region
  vpc_cidr     = var.secondary_vpc_cidr
  project_name = var.project_name
  environment  = "${var.environment}-dr"
}

module "secondary_compute" {
  source = "./modules/compute"
  providers = {
    aws = aws.secondary
  }

  vpc_id              = module.secondary_networking.vpc_id
  public_subnet_ids   = module.secondary_networking.public_subnet_ids
  private_subnet_ids  = module.secondary_networking.private_subnet_ids
  project_name        = var.project_name
  environment         = "${var.environment}-dr"
  instance_type       = var.instance_type
  min_size           = var.dr_min_size
  max_size           = var.dr_max_size
  database_endpoint   = module.secondary_database.endpoint
}

module "secondary_database" {
  source = "./modules/database"
  providers = {
    aws = aws.secondary
  }

  vpc_id             = module.secondary_networking.vpc_id
  subnet_ids         = module.secondary_networking.database_subnet_ids
  project_name       = var.project_name
  environment        = "${var.environment}-dr"
  instance_class     = var.dr_db_instance_class
  allocated_storage  = var.db_allocated_storage
  multi_az          = false
  security_group_ids = [module.secondary_compute.app_security_group_id]
}

# Global monitoring
module "monitoring" {
  source = "./modules/monitoring"
  providers = {
    aws = aws.primary
  }

  project_name = var.project_name
  environment  = var.environment
  
  monitored_resources = {
    primary_alb_arn     = module.primary_compute.alb_arn
    primary_asg_name    = module.primary_compute.asg_name
    primary_db_id       = module.primary_database.db_instance_id
    secondary_alb_arn   = module.secondary_compute.alb_arn
    secondary_asg_name  = module.secondary_compute.asg_name
    secondary_db_id     = module.secondary_database.db_instance_id
  }

  alert_email = var.alert_email
}

# S3 bucket for cross-region replication
resource "aws_s3_bucket" "primary" {
  provider = aws.primary
  bucket   = "${var.project_name}-${var.environment}-primary-${data.aws_caller_identity.current.account_id}"

  tags = {
    Name        = "${var.project_name}-primary-bucket"
    Environment = var.environment
  }
}

resource "aws_s3_bucket" "secondary" {
  provider = aws.secondary
  bucket   = "${var.project_name}-${var.environment}-secondary-${data.aws_caller_identity.current.account_id}"

  tags = {
    Name        = "${var.project_name}-secondary-bucket"
    Environment = var.environment
  }
}

# S3 bucket versioning
resource "aws_s3_bucket_versioning" "primary" {
  provider = aws.primary
  bucket   = aws_s3_bucket.primary.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_versioning" "secondary" {
  provider = aws.secondary
  bucket   = aws_s3_bucket.secondary.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Get current AWS account ID
data "aws_caller_identity" "current" {}

# IAM role for replication
resource "aws_iam_role" "replication" {
  provider = aws.primary
  name     = "${var.project_name}-s3-replication-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
      }
    ]
  })
}

# IAM policy for replication
resource "aws_iam_role_policy" "replication" {
  provider = aws.primary
  name     = "${var.project_name}-s3-replication-policy"
  role     = aws_iam_role.replication.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetReplicationConfiguration",
          "s3:ListBucket"
        ]
        Resource = aws_s3_bucket.primary.arn
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObjectVersionForReplication",
          "s3:GetObjectVersionAcl",
          "s3:GetObjectVersionTagging"
        ]
        Resource = "${aws_s3_bucket.primary.arn}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ReplicateObject",
          "s3:ReplicateDelete",
          "s3:ReplicateTags"
        ]
        Resource = "${aws_s3_bucket.secondary.arn}/*"
      }
    ]
  })
}

# S3 bucket replication configuration
resource "aws_s3_bucket_replication_configuration" "replication" {
  provider = aws.primary
  role     = aws_iam_role.replication.arn
  bucket   = aws_s3_bucket.primary.id

  rule {
    id     = "replicate-all"
    status = "Enabled"

    destination {
      bucket        = aws_s3_bucket.secondary.arn
      storage_class = "STANDARD_IA"
    }
  }

  depends_on = [aws_s3_bucket_versioning.primary]
}