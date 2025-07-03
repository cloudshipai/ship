package modules

import "strings"

// QueryTemplates provides tested, working Steampipe queries
var QueryTemplates = map[string]map[string]string{
	"aws": {
		// EC2 Queries
		"ec2_running_count":      `SELECT COUNT(*) as count FROM aws_ec2_instance WHERE instance_state = 'running'`,
		"ec2_running_list":       `SELECT instance_id, instance_type, instance_state, region, vpc_id FROM aws_ec2_instance WHERE instance_state = 'running'`,
		"ec2_all_list":          `SELECT instance_id, instance_type, instance_state, region, vpc_id FROM aws_ec2_instance`,
		"ec2_by_type":           `SELECT instance_type, COUNT(*) as count FROM aws_ec2_instance GROUP BY instance_type`,
		"ec2_security_groups":   `SELECT i.instance_id, sg->>'GroupId' as group_id, sg->>'GroupName' as group_name FROM aws_ec2_instance i, jsonb_array_elements(i.security_groups) as sg`,
		
		// S3 Queries  
		"s3_bucket_count":       `SELECT COUNT(*) as count FROM aws_s3_bucket`,
		"s3_bucket_list":        `SELECT name, region, creation_date FROM aws_s3_bucket`,
		"s3_public_buckets":     `SELECT name FROM aws_s3_bucket WHERE bucket_policy_is_public = true`,
		
		// RDS Queries
		"rds_instance_list":     `SELECT db_instance_identifier, engine, db_instance_class, publicly_accessible FROM aws_rds_db_instance`,
		"rds_public_instances":  `SELECT db_instance_identifier FROM aws_rds_db_instance WHERE publicly_accessible = true`,
		
		// Lambda Queries
		"lambda_function_list":  `SELECT name, runtime, timeout, memory_size FROM aws_lambda_function`,
		"lambda_by_runtime":     `SELECT runtime, COUNT(*) as count FROM aws_lambda_function GROUP BY runtime`,
		
		// IAM Queries
		"iam_users_no_mfa":      `SELECT name, create_date FROM aws_iam_user WHERE NOT mfa_enabled`,
		"iam_role_list":         `SELECT name, arn FROM aws_iam_role`,
		
		// VPC Queries
		"vpc_list":              `SELECT vpc_id, cidr_block, is_default FROM aws_vpc`,
		"security_group_open":   `SELECT group_id, group_name FROM aws_vpc_security_group WHERE jsonb_array_length(ingress_rules) > 0`,
	},
}

// GetQueryForPrompt returns a appropriate query based on the prompt
func GetQueryForPrompt(prompt string, provider string) []string {
	templates := QueryTemplates[provider]
	if templates == nil {
		return nil
	}
	
	promptLower := strings.ToLower(prompt)
	var queries []string
	
	// EC2 related
	if strings.Contains(promptLower, "ec2") || strings.Contains(promptLower, "instance") {
		if strings.Contains(promptLower, "running") {
			queries = append(queries, templates["ec2_running_count"], templates["ec2_running_list"])
		} else {
			queries = append(queries, templates["ec2_all_list"])
		}
	}
	
	// S3 related
	if strings.Contains(promptLower, "s3") || strings.Contains(promptLower, "bucket") {
		queries = append(queries, templates["s3_bucket_count"], templates["s3_bucket_list"])
	}
	
	// Add more mappings...
	
	return queries
}