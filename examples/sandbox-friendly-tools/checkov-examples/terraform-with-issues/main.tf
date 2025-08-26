# Terraform configuration with security issues for Checkov testing

# S3 bucket with security misconfigurations
resource "aws_s3_bucket" "example" {
  bucket = "my-insecure-bucket-12345"
  
  # Missing: bucket versioning, encryption, public access block
}

# S3 bucket with public read access (security issue)
resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "public-read"  # CKV_AWS_20: Public read access
}

# Security group with overly permissive rules
resource "aws_security_group" "example" {
  name        = "insecure-sg"
  description = "Security group with issues"

  ingress {
    description = "SSH from anywhere"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # CKV_AWS_24: SSH open to world
  }

  ingress {
    description = "HTTP from anywhere"
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
resource "aws_db_instance" "example" {
  allocated_storage    = 20
  storage_type        = "gp2"
  engine              = "mysql"
  engine_version      = "5.7"
  instance_class      = "db.t2.micro"
  db_name             = "mydb"
  username            = "foo"
  password            = "mustbeeightcharacters"  # CKV_AWS_17: Plain text password
  skip_final_snapshot = true
  
  # Missing: storage_encrypted = true
  # Missing: backup_retention_period
}

# EC2 instance with security issues
resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1d0"
  instance_type = "t2.micro"
  
  # Missing: encrypted EBS volumes
  # Missing: detailed monitoring
  
  vpc_security_group_ids = [aws_security_group.example.id]
  
  # Public IP assignment without justification
  associate_public_ip_address = true
}

# IAM policy with overly broad permissions
resource "aws_iam_policy" "example" {
  name        = "overly-broad-policy"
  description = "A policy with excessive permissions"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "*"  # CKV_AWS_60: Wildcard permissions
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}