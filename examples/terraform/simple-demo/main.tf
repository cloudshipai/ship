terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Example EC2 instances for cost estimation
resource "aws_instance" "web" {
  count         = 3
  ami           = "ami-0c55b159cbfafe1f0" # Amazon Linux 2
  instance_type = "t3.medium"
  
  root_block_device {
    volume_size = 30
    volume_type = "gp3"
  }
  
  tags = {
    Name        = "web-server-${count.index}"
    Environment = "production"
  }
}

# RDS Database for cost estimation
resource "aws_db_instance" "main" {
  identifier     = "myapp-database"
  engine         = "postgres"
  engine_version = "15.3"
  instance_class = "db.t3.medium"
  
  allocated_storage     = 100
  storage_type          = "gp3"
  storage_encrypted     = true
  
  db_name  = "myapp"
  username = "dbadmin"
  password = "temporary-password-change-me!"
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  skip_final_snapshot = true
  
  tags = {
    Name        = "myapp-database"
    Environment = "production"
  }
}

# S3 Bucket for storage costs
resource "aws_s3_bucket" "data" {
  bucket = "myapp-data-storage-demo"
  
  tags = {
    Name        = "myapp-data"
    Environment = "production"
  }
}

# Load Balancer
resource "aws_lb" "main" {
  name               = "myapp-alb"
  load_balancer_type = "application"
  subnets            = ["subnet-12345", "subnet-67890"] # Placeholder subnets
  
  tags = {
    Name        = "myapp-alb"
    Environment = "production"
  }
}