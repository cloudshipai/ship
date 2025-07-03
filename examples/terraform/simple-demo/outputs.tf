output "web_instance_count" {
  value       = length(aws_instance.web)
  description = "Number of web instances"
}

output "database_endpoint" {
  value       = aws_db_instance.main.address
  description = "RDS instance endpoint"
}

output "s3_bucket_name" {
  value       = aws_s3_bucket.data.id
  description = "S3 bucket name"
}

output "load_balancer_dns" {
  value       = aws_lb.main.dns_name
  description = "Load balancer DNS name"
}