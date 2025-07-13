package modules

import (
	"fmt"
	"regexp"
	"strings"
)

// FixSteampipeQuery attempts to fix common AI-generated query mistakes
func FixSteampipeQuery(query string, provider string) string {
	if provider != "aws" {
		return query
	}

	fixed := query

	// Fix common EC2 column mistakes
	// Replace "WHERE running" with proper syntax
	if strings.Contains(strings.ToLower(fixed), "where running") {
		fixed = regexp.MustCompile(`(?i)where\s+running`).ReplaceAllString(fixed, "WHERE instance_state = 'running'")
	}

	// Replace "WHERE state = 'running'" with correct column
	if strings.Contains(strings.ToLower(fixed), "where state") {
		fixed = regexp.MustCompile(`(?i)where\s+state\s*=`).ReplaceAllString(fixed, "WHERE instance_state =")
	}

	// Replace "WHERE state_name" with correct column
	if strings.Contains(strings.ToLower(fixed), "state_name") {
		fixed = strings.ReplaceAll(fixed, "state_name", "instance_state")
	}

	// Fix security group joins
	if strings.Contains(fixed, "sg.group_id") && !strings.Contains(fixed, "jsonb_array_elements") {
		// This is trying to join security groups incorrectly
		fixed = strings.ReplaceAll(fixed, "sg.group_id", "sg->>'GroupId'")
		fixed = strings.ReplaceAll(fixed, "sg.group_name", "sg->>'GroupName'")
	}

	// Fix common S3 mistakes
	fixed = strings.ReplaceAll(fixed, "bucket_name", "name")

	// Fix common RDS mistakes
	fixed = strings.ReplaceAll(fixed, "instance_identifier", "db_instance_identifier")

	return fixed
}

// ValidateQuery does basic validation
func ValidateQuery(query string) error {
	// Check for basic SQL injection attempts
	if strings.Contains(strings.ToLower(query), "drop table") ||
		strings.Contains(strings.ToLower(query), "delete from") {
		return fmt.Errorf("potentially dangerous query detected")
	}

	return nil
}
