package modules

// GetCommonSteampipeTables returns common Steampipe table names for each provider
func GetCommonSteampipeTables(provider string) []string {
	tables := map[string][]string{
		"aws": {
			"aws_account",
			"aws_ec2_instance",
			"aws_s3_bucket",
			"aws_iam_user",
			"aws_iam_role",
			"aws_iam_policy",
			"aws_lambda_function",
			"aws_rds_db_instance",
			"aws_vpc",
			"aws_vpc_security_group",
			"aws_vpc_subnet",
			"aws_cloudwatch_log_group",
			"aws_cloudtrail_trail",
			"aws_config_configuration_recorder",
			"aws_ebs_volume",
			"aws_ebs_snapshot",
			"aws_ecs_cluster",
			"aws_ecs_service",
			"aws_eks_cluster",
			"aws_elasticache_cluster",
			"aws_kinesis_stream",
			"aws_sns_topic",
			"aws_sqs_queue",
			"aws_dynamodb_table",
			"aws_api_gateway_rest_api",
			"aws_cloudfront_distribution",
			"aws_route53_zone",
			"aws_acm_certificate",
			"aws_kms_key",
			"aws_secretsmanager_secret",
		},
		"azure": {
			"azure_subscription",
			"azure_compute_virtual_machine",
			"azure_storage_account",
			"azure_storage_container",
			"azure_app_service_web_app",
			"azure_sql_database",
			"azure_key_vault",
			"azure_network_virtual_network",
			"azure_network_security_group",
			"azure_cosmosdb_account",
			"azure_kubernetes_cluster",
		},
		"gcp": {
			"gcp_project",
			"gcp_compute_instance",
			"gcp_storage_bucket",
			"gcp_sql_database_instance",
			"gcp_kubernetes_cluster",
			"gcp_compute_network",
			"gcp_compute_firewall",
			"gcp_iam_service_account",
			"gcp_bigquery_dataset",
			"gcp_pubsub_topic",
		},
	}

	if providerTables, ok := tables[provider]; ok {
		return providerTables
	}

	// Default return empty if provider not found
	return []string{}
}

// GetSteampipeTableExamples returns example queries for common use cases
func GetSteampipeTableExamples(provider string) map[string]string {
	examples := map[string]map[string]string{
		"aws": {
			"Count S3 buckets":      "SELECT COUNT(*) as bucket_count FROM aws_s3_bucket",
			"List S3 buckets":       "SELECT name, region FROM aws_s3_bucket LIMIT 10",
			"Check S3 encryption":   "SELECT name, server_side_encryption_configuration FROM aws_s3_bucket LIMIT 5",
			"Count EC2 instances":   "SELECT COUNT(*) as instance_count FROM aws_ec2_instance",
			"Count running EC2s":    "SELECT COUNT(*) as running_count FROM aws_ec2_instance WHERE instance_state = 'running'",
			"List EC2 instances":    "SELECT instance_id, instance_type, instance_state, region, vpc_id FROM aws_ec2_instance",
			"List IAM users":        "SELECT name, create_date, mfa_enabled, password_last_used FROM aws_iam_user",
			"IAM users without MFA": "SELECT name, create_date FROM aws_iam_user WHERE NOT mfa_enabled",
			"List IAM roles":        "SELECT name, arn, create_date FROM aws_iam_role",
			"List IAM policies":     "SELECT name, arn, attachment_count FROM aws_iam_policy WHERE is_attachable = true",
			"Public RDS instances":  "SELECT db_instance_identifier, engine FROM aws_rds_db_instance WHERE publicly_accessible = true",
			"S3 bucket encryption":  "SELECT name, server_side_encryption_configuration FROM aws_s3_bucket",
			"Check account info":    "SELECT account_id, arn FROM aws_account",
			"EC2 security groups":   "SELECT i.instance_id, sg->>'GroupId' as group_id FROM aws_ec2_instance i, jsonb_array_elements(i.security_groups) as sg",
		},
		"azure": {
			"List VMs":                 "SELECT name, location, vm_size FROM azure_compute_virtual_machine",
			"Check storage encryption": "SELECT name, encryption FROM azure_storage_account",
			"List SQL databases":       "SELECT name, edition, service_level_objective FROM azure_sql_database",
		},
		"gcp": {
			"List compute instances":   "SELECT name, machine_type, status FROM gcp_compute_instance",
			"Check bucket permissions": "SELECT name, location, storage_class FROM gcp_storage_bucket",
			"List Kubernetes clusters": "SELECT name, location, status FROM gcp_kubernetes_cluster",
		},
	}

	if providerExamples, ok := examples[provider]; ok {
		return providerExamples
	}

	return map[string]string{}
}
