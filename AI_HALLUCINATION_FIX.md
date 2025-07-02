# AI Hallucination Fix - Ship CLI

## Problem
The AI was hallucinating non-existent Steampipe table names like `aws_resource_inventory` when generating infrastructure investigation queries. This caused all AI-powered queries to fail.

## Root Cause
1. **No context about real tables**: The AI prompts didn't include any information about actual Steampipe tables
2. **Hardcoded fallback**: `llm_dagger.go` had a hardcoded fake table name in the fallback logic
3. **Generic prompts**: The AI was told to "use appropriate Steampipe tables" without knowing what tables exist

## Solution Implemented

### 1. Created Table Reference (`steampipe_tables.go`)
- Added comprehensive list of real Steampipe tables for AWS, Azure, and GCP
- Included example queries for common use cases
- Functions: `GetCommonSteampipeTables()` and `GetSteampipeTableExamples()`

### 2. Enhanced AI Prompts
- Modified `CreateInvestigationPlan` in `llm_dagger.go`
- Now includes actual table names in the prompt
- Added example queries to guide the AI
- Explicit instruction: "NEVER use made-up table names"

### 3. Fixed Fallback Logic
- Replaced hardcoded `aws_resource_inventory` with real tables
- Now uses `aws_account`, `azure_subscription`, or `gcp_project` as fallbacks
- Added JSON parsing with markdown code block handling

### 4. AWS Credentials Fix
- Fixed credential mounting in Steampipe containers
- Now parses `~/.aws/credentials` and passes as environment variables
- Removed problematic profile-based authentication

## Results

### Before:
```sql
-- AI generated invalid queries:
SELECT * FROM aws_resource_inventory  -- Doesn't exist!
```

### After:
```sql
-- AI now generates valid queries:
SELECT name, versioning_enabled FROM aws_s3_bucket
SELECT instance_id, instance_type FROM aws_ec2_instance
SELECT name, create_date FROM aws_iam_user
```

## Testing
```bash
# Test AI investigation
./ship ai-investigate --prompt "List all S3 buckets" --provider aws

# Test with execution
./ship ai-investigate --prompt "Check IAM users" --provider aws --execute

# Direct query (always worked, now credentials work too)
./ship query "SELECT * FROM aws_s3_bucket" --provider aws
```

## Future Improvements
1. Add more Steampipe tables as they're released
2. Implement query validation before execution
3. Add caching for table schemas
4. Improve error messages when queries fail