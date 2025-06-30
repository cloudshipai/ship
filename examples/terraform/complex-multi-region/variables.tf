variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "ship-complex"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "prod"
}

variable "primary_region" {
  description = "Primary AWS region"
  type        = string
  default     = "us-east-1"
}

variable "secondary_region" {
  description = "Secondary AWS region for DR"
  type        = string
  default     = "us-west-2"
}

variable "primary_vpc_cidr" {
  description = "CIDR block for primary VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "secondary_vpc_cidr" {
  description = "CIDR block for secondary VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "instance_type" {
  description = "EC2 instance type for application servers"
  type        = string
  default     = "t3.small"
}

variable "min_size" {
  description = "Minimum number of instances in primary region"
  type        = number
  default     = 3
}

variable "max_size" {
  description = "Maximum number of instances in primary region"
  type        = number
  default     = 10
}

variable "dr_min_size" {
  description = "Minimum number of instances in DR region"
  type        = number
  default     = 1
}

variable "dr_max_size" {
  description = "Maximum number of instances in DR region"
  type        = number
  default     = 5
}

variable "db_instance_class" {
  description = "RDS instance class for primary database"
  type        = string
  default     = "db.t3.micro"
}

variable "dr_db_instance_class" {
  description = "RDS instance class for DR database"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "Allocated storage for RDS in GB"
  type        = number
  default     = 20
}

variable "alert_email" {
  description = "Email address for monitoring alerts"
  type        = string
  default     = "alerts@example.com"
}