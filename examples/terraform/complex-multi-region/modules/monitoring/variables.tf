variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "monitored_resources" {
  description = "Resources to monitor"
  type = object({
    primary_alb_arn    = string
    primary_asg_name   = string
    primary_db_id      = string
    secondary_alb_arn  = string
    secondary_asg_name = string
    secondary_db_id    = string
  })
}

variable "alert_email" {
  description = "Email address for alerts"
  type        = string
}