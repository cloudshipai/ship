# Policy Sentry MCP Tool

Policy Sentry is an AWS IAM policy generator that helps create least-privilege IAM policies using a guided approach.

## Description

Policy Sentry provides:
- Template-based IAM policy generation
- Query capabilities for AWS service actions, conditions, and ARN formats
- Least-privilege policy creation from YAML templates
- IAM database initialization and maintenance

## MCP Functions

### `policy_sentry_create_template`
Create IAM policy template using real policy_sentry CLI.

**Parameters:**
- `template_type` (required): Type of template (crud, actions)
- `output_file`: Output file path for template

**CLI Command:** `policy_sentry create-template --template-type <type> [--output-file <file>]`

### `policy_sentry_write_policy`
Write IAM policy from input YAML file using real policy_sentry CLI.

**Parameters:**
- `input_file` (required): Path of the YAML file for generating policies
- `minimize`: Minimize policy statements
- `format`: Output format (yaml, json, terraform)

**CLI Command:** `policy_sentry write-policy --input-file <file> [--minimize] [--fmt <format>]`

### `policy_sentry_initialize`
Initialize Policy Sentry IAM database using real policy_sentry CLI.

**Parameters:**
- `fetch`: Fetch latest AWS documentation from AWS docs

**CLI Command:** `policy_sentry initialize [--fetch]`

### `policy_sentry_query_action_table`
Query AWS service action table using real policy_sentry CLI.

**Parameters:**
- `service` (required): AWS service name (e.g., s3, ec2, iam)
- `name`: IAM Action name
- `access_level`: Access level filter (read, write, list, tagging, permissions-management)
- `condition`: Condition key filter
- `resource_type`: Resource type filter
- `format`: Output format (yaml, json)

**CLI Command:** `policy_sentry query action-table --service <service> [options]`

### `policy_sentry_query_condition_table`
Query AWS service condition table using real policy_sentry CLI.

**Parameters:**
- `service` (required): AWS service name (e.g., s3, ec2, iam)
- `name`: Condition key name
- `format`: Output format (yaml, json)

**CLI Command:** `policy_sentry query condition-table --service <service> [options]`

### `policy_sentry_query_arn_table`
Query AWS service ARN table using real policy_sentry CLI.

**Parameters:**
- `service` (required): AWS service name (e.g., s3, ec2, iam)
- `name`: Resource ARN type name
- `list_arn_types`: List ARN types
- `format`: Output format (yaml, json)

**CLI Command:** `policy_sentry query arn-table --service <service> [options]`

### `policy_sentry_query_service_table`
Query AWS service table using real policy_sentry CLI.

**Parameters:**
- `format`: Output format (yaml, json, csv)

**CLI Command:** `policy_sentry query service-table [--fmt <format>]`

## Common Use Cases

1. **Template Creation**: Use `policy_sentry_create_template` to generate YAML templates for policy creation
2. **Policy Generation**: Use `policy_sentry_write_policy` to create least-privilege policies from templates
3. **Database Management**: Use `policy_sentry_initialize` to set up and update the IAM database
4. **Service Exploration**: Use the query functions to explore AWS service capabilities
5. **Action Discovery**: Use `policy_sentry_query_action_table` to find available actions for services
6. **Condition Analysis**: Use `policy_sentry_query_condition_table` to understand condition keys
7. **ARN Format Discovery**: Use `policy_sentry_query_arn_table` to understand resource ARN formats

## Template Types

- **crud**: Creates template for CRUD-based access levels (read, write, list, tagging, permissions-management)
- **actions**: Creates template for specific action-based policies

## Integration with Ship CLI

All Policy Sentry tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Generate least-privilege IAM policies automatically
- Query AWS service capabilities and restrictions
- Create policy templates for common use cases
- Explore IAM action, condition, and ARN requirements

The tools use containerized execution via Dagger for consistent, isolated policy generation environments.