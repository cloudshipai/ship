output "endpoint" {
  description = "Database endpoint"
  value       = aws_db_instance.main.endpoint
}

output "db_instance_id" {
  description = "Database instance ID"
  value       = aws_db_instance.main.id
}

output "db_name" {
  description = "Database name"
  value       = aws_db_instance.main.db_name
}

output "password_ssm_parameter" {
  description = "SSM parameter name for database password"
  value       = aws_ssm_parameter.db_password.name
}