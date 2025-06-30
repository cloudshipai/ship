output "primary_alb_dns" {
  description = "DNS name of the primary ALB"
  value       = module.primary_compute.alb_dns_name
}

output "secondary_alb_dns" {
  description = "DNS name of the secondary ALB"
  value       = module.secondary_compute.alb_dns_name
}

output "primary_database_endpoint" {
  description = "Primary database endpoint"
  value       = module.primary_database.endpoint
}

output "secondary_database_endpoint" {
  description = "Secondary database endpoint"
  value       = module.secondary_database.endpoint
}

output "primary_s3_bucket" {
  description = "Primary S3 bucket name"
  value       = aws_s3_bucket.primary.id
}

output "secondary_s3_bucket" {
  description = "Secondary S3 bucket name"
  value       = aws_s3_bucket.secondary.id
}

output "monitoring_dashboard" {
  description = "CloudWatch dashboard URL"
  value       = module.monitoring.dashboard_url
}

output "sns_topic_arn" {
  description = "SNS topic ARN for alerts"
  value       = module.monitoring.sns_topic_arn
}