package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudship/ship/internal/dagger/modules"
)

// QueryTemplate represents a template for generating Steampipe queries
type QueryTemplate struct {
	ResourceType string
	Fields       []string
	Table        string
	Conditions   []string
}

// Common AWS query templates
var awsQueryTemplates = map[string]QueryTemplate{
	"s3_buckets": {
		ResourceType: "S3 Buckets",
		Fields:       []string{"name", "arn", "region", "creation_date", "versioning_enabled"},
		Table:        "aws_s3_bucket",
	},
	"ec2_instances": {
		ResourceType: "EC2 Instances",
		Fields:       []string{"instance_id", "instance_type", "instance_state", "public_ip_address", "private_ip_address", "tags"},
		Table:        "aws_ec2_instance",
	},
	"security_groups": {
		ResourceType: "Security Groups",
		Fields:       []string{"group_id", "group_name", "description", "vpc_id"},
		Table:        "aws_vpc_security_group",
	},
	"iam_users": {
		ResourceType: "IAM Users",
		Fields:       []string{"name", "user_id", "arn", "create_date", "mfa_enabled"},
		Table:        "aws_iam_user",
	},
	"rds_instances": {
		ResourceType: "RDS Instances",
		Fields:       []string{"db_instance_identifier", "engine", "engine_version", "db_instance_class", "allocated_storage", "publicly_accessible"},
		Table:        "aws_rds_db_instance",
	},
	"lambda_functions": {
		ResourceType: "Lambda Functions",
		Fields:       []string{"name", "arn", "runtime", "timeout", "memory_size", "last_modified"},
		Table:        "aws_lambda_function",
	},
	"vpc": {
		ResourceType: "VPCs",
		Fields:       []string{"vpc_id", "cidr_block", "state", "is_default", "tags"},
		Table:        "aws_vpc",
	},
	"ebs_volumes": {
		ResourceType: "EBS Volumes",
		Fields:       []string{"volume_id", "size", "volume_type", "encrypted", "state", "availability_zone"},
		Table:        "aws_ebs_volume",
	},
	"elb": {
		ResourceType: "Load Balancers",
		Fields:       []string{"name", "dns_name", "scheme", "type", "state"},
		Table:        "aws_ec2_load_balancer",
	},
	"cloudtrail": {
		ResourceType: "CloudTrail Trails",
		Fields:       []string{"name", "arn", "is_logging", "is_multi_region_trail", "log_file_validation_enabled"},
		Table:        "aws_cloudtrail_trail",
	},
}

// DynamicQueryGenerator generates Steampipe queries based on natural language prompts
type DynamicQueryGenerator struct {
	Provider string
}

// GenerateQuery creates a Steampipe query from a natural language prompt
func (g *DynamicQueryGenerator) GenerateQuery(prompt string) (string, string) {
	lowerPrompt := strings.ToLower(prompt)
	
	// Check for security-related queries
	if strings.Contains(lowerPrompt, "security") || strings.Contains(lowerPrompt, "vulnerable") || 
	   strings.Contains(lowerPrompt, "public") || strings.Contains(lowerPrompt, "exposed") {
		return g.generateSecurityQuery(lowerPrompt)
	}
	
	// Check for cost-related queries
	if strings.Contains(lowerPrompt, "cost") || strings.Contains(lowerPrompt, "expensive") || 
	   strings.Contains(lowerPrompt, "billing") {
		return g.generateCostQuery(lowerPrompt)
	}
	
	// Check for compliance queries
	if strings.Contains(lowerPrompt, "compliance") || strings.Contains(lowerPrompt, "encrypt") || 
	   strings.Contains(lowerPrompt, "audit") {
		return g.generateComplianceQuery(lowerPrompt)
	}
	
	// Resource-specific queries
	for key, template := range awsQueryTemplates {
		keywords := g.getResourceKeywords(key)
		for _, keyword := range keywords {
			if strings.Contains(lowerPrompt, keyword) {
				return g.generateResourceQuery(template, lowerPrompt)
			}
		}
	}
	
	// Default to a general resource listing
	return g.generateGeneralQuery(lowerPrompt)
}

// getResourceKeywords returns keywords for each resource type
func (g *DynamicQueryGenerator) getResourceKeywords(resourceType string) []string {
	keywords := map[string][]string{
		"s3_buckets":       {"s3", "bucket", "storage", "object storage"},
		"ec2_instances":    {"ec2", "instance", "server", "compute", "vm", "virtual machine"},
		"security_groups":  {"security group", "firewall", "sg"},
		"iam_users":        {"iam", "user", "identity", "access"},
		"rds_instances":    {"rds", "database", "db", "mysql", "postgres", "aurora"},
		"lambda_functions": {"lambda", "function", "serverless"},
		"vpc":              {"vpc", "network", "subnet"},
		"ebs_volumes":      {"ebs", "volume", "disk", "storage"},
		"elb":              {"elb", "load balancer", "alb", "nlb"},
		"cloudtrail":       {"cloudtrail", "audit", "logging", "trail"},
	}
	return keywords[resourceType]
}

// generateSecurityQuery creates security-focused queries
func (g *DynamicQueryGenerator) generateSecurityQuery(prompt string) (string, string) {
	queries := []string{}
	
	// Public resources
	if strings.Contains(prompt, "public") || strings.Contains(prompt, "exposed") {
		queries = append(queries,
			"SELECT name, arn, region, acl FROM aws_s3_bucket WHERE acl NOT LIKE '%private%'",
			"SELECT db_instance_identifier, publicly_accessible, endpoint_address FROM aws_rds_db_instance WHERE publicly_accessible = true",
			"SELECT instance_id, public_ip_address, state FROM aws_ec2_instance WHERE public_ip_address IS NOT NULL",
		)
	}
	
	// Security groups
	if strings.Contains(prompt, "security") || strings.Contains(prompt, "firewall") || strings.Contains(prompt, "0.0.0.0") {
		queries = append(queries,
			"SELECT group_name, group_id, from_port, to_port, cidr_ip FROM aws_vpc_security_group_rule WHERE cidr_ip = '0.0.0.0/0' AND type = 'ingress'",
		)
	}
	
	// General security scan
	if len(queries) == 0 {
		queries = append(queries,
			"SELECT name, arn, region FROM aws_s3_bucket WHERE bucket_policy_is_public = true OR acl NOT LIKE '%private%'",
			"SELECT group_name, group_id, from_port, to_port, cidr_ip FROM aws_vpc_security_group_rule WHERE cidr_ip = '0.0.0.0/0'",
			"SELECT volume_id, encrypted, size FROM aws_ebs_volume WHERE encrypted = false",
		)
	}
	
	if len(queries) > 0 {
		return queries[0], "Security Analysis"
	}
	return "SELECT 'No specific security query generated' as message", "Security Analysis"
}

// generateCostQuery creates cost-related queries
func (g *DynamicQueryGenerator) generateCostQuery(prompt string) (string, string) {
	queries := []string{
		"SELECT instance_id, instance_type, state, placement_availability_zone FROM aws_ec2_instance WHERE state = 'running' ORDER BY instance_type DESC",
		"SELECT volume_id, size, volume_type, state FROM aws_ebs_volume WHERE state != 'in-use' ORDER BY size DESC",
		"SELECT db_instance_identifier, db_instance_class, engine, allocated_storage FROM aws_rds_db_instance ORDER BY allocated_storage DESC",
	}
	
	if strings.Contains(prompt, "unused") || strings.Contains(prompt, "idle") {
		return "SELECT volume_id, size, volume_type, state, create_time FROM aws_ebs_volume WHERE state = 'available'", "Unused Resources"
	}
	
	return queries[0], "Cost Analysis"
}

// generateComplianceQuery creates compliance-focused queries
func (g *DynamicQueryGenerator) generateComplianceQuery(prompt string) (string, string) {
	if strings.Contains(prompt, "encrypt") {
		return "SELECT volume_id, encrypted, size, volume_type FROM aws_ebs_volume WHERE encrypted = false", "Encryption Compliance"
	}
	
	if strings.Contains(prompt, "mfa") {
		return "SELECT name, user_id, mfa_enabled, password_last_used FROM aws_iam_user WHERE mfa_enabled = false", "MFA Compliance"
	}
	
	if strings.Contains(prompt, "logging") || strings.Contains(prompt, "audit") {
		return "SELECT name, is_logging, is_multi_region_trail, log_file_validation_enabled FROM aws_cloudtrail_trail", "Audit Compliance"
	}
	
	// General compliance check
	return "SELECT name, user_id, mfa_enabled FROM aws_iam_user WHERE mfa_enabled = false", "Compliance Check"
}

// generateResourceQuery creates resource-specific queries
func (g *DynamicQueryGenerator) generateResourceQuery(template QueryTemplate, prompt string) (string, string) {
	fields := strings.Join(template.Fields, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", fields, template.Table)
	
	// Add conditions based on prompt
	conditions := []string{}
	
	// Filter by region if mentioned
	if strings.Contains(prompt, "us-east-1") || strings.Contains(prompt, "us-west-2") || 
	   strings.Contains(prompt, "eu-west-1") {
		for _, region := range []string{"us-east-1", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1"} {
			if strings.Contains(prompt, region) {
				conditions = append(conditions, fmt.Sprintf("region = '%s'", region))
				break
			}
		}
	}
	
	// Add state filters
	if strings.Contains(prompt, "running") {
		if strings.Contains(template.Table, "ec2") {
			conditions = append(conditions, "instance_state = 'running'")
		} else {
			conditions = append(conditions, "state = 'running'")
		}
	} else if strings.Contains(prompt, "stopped") {
		if strings.Contains(template.Table, "ec2") {
			conditions = append(conditions, "instance_state = 'stopped'")
		} else {
			conditions = append(conditions, "state = 'stopped'")
		}
	}
	
	// Add date filters
	if strings.Contains(prompt, "recent") || strings.Contains(prompt, "new") {
		if strings.Contains(template.Table, "s3") {
			conditions = append(conditions, "creation_date > now() - interval '30 days'")
		} else if strings.Contains(fields, "create_time") {
			conditions = append(conditions, "create_time > now() - interval '30 days'")
		}
	}
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Add sorting
	if strings.Contains(prompt, "recent") || strings.Contains(prompt, "latest") {
		if strings.Contains(template.Table, "s3") {
			query += " ORDER BY creation_date DESC"
		} else if strings.Contains(fields, "create_time") {
			query += " ORDER BY create_time DESC"
		}
	} else if strings.Contains(prompt, "large") || strings.Contains(prompt, "big") {
		if strings.Contains(fields, "size") {
			query += " ORDER BY size DESC"
		}
	}
	
	// Add limit if requested
	if strings.Contains(prompt, "top") || strings.Contains(prompt, "first") {
		numbers := []string{"5", "10", "20", "50", "100"}
		for _, num := range numbers {
			if strings.Contains(prompt, num) {
				query += " LIMIT " + num
				break
			}
		}
		if !strings.Contains(query, "LIMIT") {
			query += " LIMIT 10"
		}
	}
	
	return query, template.ResourceType + " Query"
}

// generateGeneralQuery creates a general resource listing query
func (g *DynamicQueryGenerator) generateGeneralQuery(prompt string) (string, string) {
	// Default to listing main resources
	queries := []string{
		"SELECT 'EC2' as service, COUNT(*) as count FROM aws_ec2_instance UNION ALL SELECT 'S3' as service, COUNT(*) as count FROM aws_s3_bucket UNION ALL SELECT 'RDS' as service, COUNT(*) as count FROM aws_rds_db_instance",
	}
	
	return queries[0], "Resource Overview"
}

// GenerateInvestigationPlan creates a comprehensive investigation plan based on the prompt
func GenerateInvestigationPlan(ctx context.Context, prompt string, provider string) []modules.InvestigationStep {
	generator := &DynamicQueryGenerator{Provider: provider}
	steps := []modules.InvestigationStep{}
	
	// Analyze prompt to determine investigation type
	lowerPrompt := strings.ToLower(prompt)
	
	// Generate primary query based on prompt
	query1, desc1 := generator.GenerateQuery(prompt)
	steps = append(steps, modules.InvestigationStep{
		StepNumber:       1,
		Description:      desc1,
		Provider:         provider,
		Query:            query1,
		ExpectedInsights: "Primary investigation results based on your request",
	})
	
	// Add complementary queries based on context
	if strings.Contains(lowerPrompt, "s3") {
		// Add S3 security check
		steps = append(steps, modules.InvestigationStep{
			StepNumber:       len(steps) + 1,
			Description:      "S3 Security Check",
			Provider:         provider,
			Query:            "SELECT name, bucket_policy_is_public, acl, versioning_enabled, logging FROM aws_s3_bucket",
			ExpectedInsights: "S3 bucket security configuration",
		})
	}
	
	if strings.Contains(lowerPrompt, "ec2") || strings.Contains(lowerPrompt, "instance") {
		// Add instance security check
		steps = append(steps, modules.InvestigationStep{
			StepNumber:       len(steps) + 1,
			Description:      "Instance Security Analysis",
			Provider:         provider,
			Query:            "SELECT instance_id, instance_state, public_ip_address, security_groups FROM aws_ec2_instance WHERE instance_state = 'running'",
			ExpectedInsights: "Security configuration for running instances",
		})
	}
	
	if strings.Contains(lowerPrompt, "cost") || strings.Contains(lowerPrompt, "expensive") {
		// Add cost optimization checks
		steps = append(steps, modules.InvestigationStep{
			StepNumber:       len(steps) + 1,
			Description:      "Unused Resources Check",
			Provider:         provider,
			Query:            "SELECT 'EBS Volumes' as resource_type, COUNT(*) as unused_count, SUM(size) as total_gb FROM aws_ebs_volume WHERE state = 'available'",
			ExpectedInsights: "Identify resources that are provisioned but not in use",
		})
	}
	
	// Always add a general security overview unless already focused on security
	if !strings.Contains(lowerPrompt, "security") && len(steps) < 3 {
		steps = append(steps, modules.InvestigationStep{
			StepNumber:       len(steps) + 1,
			Description:      "Security Overview",
			Provider:         provider,
			Query:            "SELECT 'Public S3 Buckets' as issue, COUNT(*) as count FROM aws_s3_bucket WHERE bucket_policy_is_public = true UNION ALL SELECT 'Open Security Groups' as issue, COUNT(DISTINCT group_id) as count FROM aws_vpc_security_group_rule WHERE cidr_ip = '0.0.0.0/0' AND type = 'ingress'",
			ExpectedInsights: "High-level security posture assessment",
		})
	}
	
	return steps
}

// ParseQueryResults analyzes the results and generates insights
func ParseQueryResults(results map[string]interface{}, prompt string) string {
	insights := []string{}
	
	// Analyze each step's results
	for stepKey, stepResults := range results {
		if queryResults, ok := stepResults.([]map[string]interface{}); ok && len(queryResults) > 0 {
			stepInsight := analyzeStepResults(stepKey, queryResults, prompt)
			if stepInsight != "" {
				insights = append(insights, stepInsight)
			}
		}
	}
	
	if len(insights) == 0 {
		return "No significant findings based on the investigation."
	}
	
	return strings.Join(insights, "\n\n")
}

// analyzeStepResults generates insights for a specific query result
func analyzeStepResults(stepKey string, results []map[string]interface{}, prompt string) string {
	if len(results) == 0 {
		return ""
	}
	
	// Count results
	count := len(results)
	
	// Analyze based on content
	insights := []string{}
	
	// Check for security issues
	for _, result := range results {
		if cidr, ok := result["cidr_ip"].(string); ok && cidr == "0.0.0.0/0" {
			if groupName, ok := result["group_name"].(string); ok {
				insights = append(insights, fmt.Sprintf("⚠️  Security group '%s' allows traffic from anywhere (0.0.0.0/0)", groupName))
			}
		}
		
		if public, ok := result["publicly_accessible"].(bool); ok && public {
			if identifier, ok := result["db_instance_identifier"].(string); ok {
				insights = append(insights, fmt.Sprintf("⚠️  RDS instance '%s' is publicly accessible", identifier))
			}
		}
		
		if encrypted, ok := result["encrypted"].(bool); ok && !encrypted {
			if volumeID, ok := result["volume_id"].(string); ok {
				insights = append(insights, fmt.Sprintf("⚠️  EBS volume '%s' is not encrypted", volumeID))
			}
		}
	}
	
	// General count insights
	if count > 0 && len(insights) == 0 {
		insights = append(insights, fmt.Sprintf("Found %d results for %s", count, stepKey))
	}
	
	return strings.Join(insights, "\n")
}