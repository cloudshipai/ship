variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "cost-opt"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "ami_id" {
  description = "AMI ID for EC2 instances"
  type        = string
  default     = "ami-0c02fb55956c7d316"
}

variable "target_capacity" {
  description = "Target capacity for spot fleet"
  type        = number
  default     = 2
}

variable "spot_price" {
  description = "Maximum spot price"
  type        = string
  default     = "0.05"
}