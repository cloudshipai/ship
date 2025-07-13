variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "serverless-api"
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