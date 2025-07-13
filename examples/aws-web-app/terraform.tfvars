project_name = "demo-web-app"
environment  = "dev"
aws_region   = "us-east-1"
vpc_cidr     = "10.0.0.0/16"

# Database configuration
db_username = "admin"
db_password = "temporarypassword123!"

# Instance configuration
instance_type     = "t3.micro"
min_size         = 1
max_size         = 3
desired_capacity = 2

# Domain configuration (optional)
domain_name = ""
ssl_cert_arn = ""