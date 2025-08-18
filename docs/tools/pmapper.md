# PMapper MCP Tool

PMapper is an AWS IAM privilege escalation analysis tool that identifies potential privilege escalation paths in AWS accounts.

## Description

PMapper creates a graph of IAM permissions and relationships to identify:
- Privilege escalation paths through role assumptions
- Cross-account access relationships
- IAM permissions and their implications
- Potentially dangerous permission combinations

## MCP Functions

### `pmapper_graph_create`
Create IAM privilege graph using real pmapper CLI.

**Parameters:**
- `profile`: AWS profile to use
- `account`: AWS account number

**CLI Command:** `pmapper [--profile <profile>] [--account <account>] graph create`

### `pmapper_query`
Query IAM permissions using real pmapper CLI.

**Parameters:**
- `query_string` (required): Query string (e.g., 'who can do iam:CreateUser')
- `profile`: AWS profile to use
- `account`: AWS account number

**CLI Command:** `pmapper [--profile <profile>] [--account <account>] query <query_string>`

### `pmapper_query_privesc`
Find privilege escalation paths using real pmapper CLI preset query.

**Parameters:**
- `target` (required): Target principal or wildcard (*) for all principals
- `profile`: AWS profile to use
- `account`: AWS account number

**CLI Command:** `pmapper [--profile <profile>] [--account <account>] query "preset privesc <target>"`

### `pmapper_visualize`
Visualize IAM privilege graph using real pmapper CLI.

**Parameters:**
- `filetype`: Output file type (svg, png, etc.)
- `profile`: AWS profile to use
- `account`: AWS account number

**CLI Command:** `pmapper [--profile <profile>] [--account <account>] visualize [--filetype <filetype>]`

### `pmapper_query_who_can`
Query who can perform specific action using real pmapper CLI.

**Parameters:**
- `action` (required): AWS action to check (e.g., iam:CreateUser)
- `profile`: AWS profile to use
- `account`: AWS account number

**CLI Command:** `pmapper [--profile <profile>] [--account <account>] query "who can do <action>"`

### `pmapper_argquery`
Advanced query with conditions using real pmapper CLI.

**Parameters:**
- `action` (required): AWS action to check (e.g., ec2:RunInstances)
- `condition`: Condition to check (e.g., ec2:InstanceType=c6gd.16xlarge)
- `profile`: AWS profile to use
- `account`: AWS account number
- `skip_admin`: Skip reporting current admin users (-s flag)

**CLI Command:** `pmapper [--profile <profile>] [--account <account>] argquery [-s] --action <action> [--condition <condition>]`

## Common Use Cases

1. **Graph Creation**: Use `pmapper_graph_create` to build the IAM permission graph for analysis
2. **Permission Queries**: Use `pmapper_query` with custom query strings to explore permissions
3. **Privilege Escalation**: Use `pmapper_query_privesc` to find escalation paths
4. **Visualization**: Use `pmapper_visualize` to create visual representations of permissions
5. **Action Analysis**: Use `pmapper_query_who_can` to see who has specific permissions
6. **Advanced Conditions**: Use `pmapper_argquery` for complex permission analysis with conditions

## Query Examples

- Find privilege escalation paths: `preset privesc *`
- Check who can create users: `who can do iam:CreateUser`
- Find admin users: `who can do iam:*`
- Check EC2 permissions: `who can do ec2:RunInstances`

## Integration with Ship CLI

All PMapper tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Perform automated IAM privilege escalation analysis
- Create visual maps of AWS account permissions
- Identify potential security risks in IAM configurations
- Generate reports on permission relationships

The tools use containerized execution via Dagger for consistent, isolated analysis environments.