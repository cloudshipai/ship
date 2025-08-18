# IaC Plan

Infrastructure as Code planning and management using Terraform CLI.

## Description

IaC Plan provides comprehensive Infrastructure as Code planning and management capabilities using the official Terraform CLI. It enables teams to plan, validate, format, and manage Terraform configurations and workspaces through a unified interface. The tool focuses on the planning phase of infrastructure management, providing essential operations for safe and reliable infrastructure deployments.

## MCP Tools

### Core Planning Operations
- **`iac_plan_terraform_plan`** - Generate Terraform execution plan
- **`iac_plan_terraform_show`** - Show Terraform plan in human-readable format
- **`iac_plan_terraform_init`** - Initialize Terraform working directory

### Configuration Management
- **`iac_plan_terraform_validate`** - Validate Terraform configuration syntax
- **`iac_plan_terraform_format`** - Format Terraform configuration files

### Workspace Management
- **`iac_plan_terraform_workspace`** - Manage Terraform workspaces

### Visualization
- **`iac_plan_terraform_graph`** - Generate Terraform dependency graph

## Real CLI Commands Used

### Core Terraform Commands
- `terraform plan` - Create execution plan
- `terraform plan -var-file=<file>` - Plan with variable file
- `terraform plan -out=<file>` - Save plan to file
- `terraform plan -destroy` - Create destroy plan
- `terraform plan -detailed-exitcode` - Enable detailed exit codes

### Configuration Commands
- `terraform validate` - Validate configuration syntax
- `terraform validate -json` - JSON validation output
- `terraform fmt` - Format configuration files
- `terraform fmt -check` - Check formatting without changes
- `terraform fmt -diff` - Show formatting differences

### Workspace Commands
- `terraform workspace list` - List workspaces
- `terraform workspace show` - Show current workspace
- `terraform workspace new <name>` - Create new workspace
- `terraform workspace select <name>` - Select workspace
- `terraform workspace delete <name>` - Delete workspace

### Advanced Commands
- `terraform show` - Show state or plan
- `terraform show <plan-file>` - Show specific plan file
- `terraform show -json` - JSON output format
- `terraform graph` - Generate dependency graph
- `terraform init` - Initialize working directory
- `terraform init -upgrade` - Upgrade providers and modules

## Use Cases

### Development Workflow
- **Plan Validation**: Verify changes before applying
- **Code Review**: Generate plans for review processes
- **Testing**: Validate configurations in development
- **Formatting**: Maintain consistent code style

### CI/CD Integration
- **Automated Planning**: Generate plans in pipelines
- **Validation Gates**: Block invalid configurations
- **Change Detection**: Identify infrastructure changes
- **Approval Workflows**: Human review of planned changes

### Infrastructure Management
- **Change Planning**: Understand impact before deployment
- **Workspace Management**: Organize environments
- **Dependency Analysis**: Visualize resource relationships
- **Configuration Maintenance**: Keep code clean and validated

### Team Collaboration
- **Shared Planning**: Standardized planning process
- **Code Standards**: Enforce formatting and validation
- **Environment Isolation**: Separate workspaces per environment
- **Documentation**: Generate visual representations

## Configuration Examples

### Basic Planning
```bash
# Initialize working directory
terraform init

# Create execution plan
terraform plan

# Save plan to file
terraform plan -out=tfplan

# Show saved plan
terraform show tfplan

# Show plan in JSON format
terraform show -json tfplan
```

### Variable Management
```bash
# Plan with variable file
terraform plan -var-file=prod.tfvars

# Plan with multiple variable files
terraform plan -var-file=common.tfvars -var-file=prod.tfvars

# Plan with inline variables
terraform plan -var="instance_type=t3.micro"
```

### Workspace Operations
```bash
# List available workspaces
terraform workspace list

# Create new workspace
terraform workspace new production

# Select workspace
terraform workspace select production

# Show current workspace
terraform workspace show

# Delete workspace
terraform workspace delete staging
```

### Validation and Formatting
```bash
# Validate configuration
terraform validate

# Format configuration files
terraform fmt

# Check if files need formatting
terraform fmt -check

# Show formatting differences
terraform fmt -diff
```

## Advanced Features

### Plan Analysis
```bash
# Generate detailed plan
terraform plan -detailed-exitcode

# Create destroy plan
terraform plan -destroy

# Plan specific resources
terraform plan -target=aws_instance.web

# Plan with parallelism control
terraform plan -parallelism=5
```

### Graph Generation
```bash
# Generate dependency graph
terraform graph > graph.dot

# Generate plan graph
terraform graph -type=plan > plan-graph.dot

# Generate apply graph
terraform graph -type=apply > apply-graph.dot

# Convert to visual format (requires Graphviz)
dot -Tpng graph.dot -o graph.png
```

### JSON Output Processing
```bash
# Get plan in JSON format
terraform show -json tfplan > plan.json

# Extract specific information
terraform show -json tfplan | jq '.resource_changes[].type' | sort | uniq

# Analyze planned changes
terraform show -json tfplan | jq '.resource_changes[] | select(.change.actions[] == "create")'
```

## Integration Patterns

### GitHub Actions
```yaml
name: Terraform Plan
on: [pull_request]
jobs:
  plan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: hashicorp/setup-terraform@v1
    
    - name: Terraform Init
      run: terraform init
      
    - name: Terraform Validate
      run: terraform validate
      
    - name: Terraform Format Check
      run: terraform fmt -check
      
    - name: Terraform Plan
      run: terraform plan -out=tfplan
      
    - name: Save Plan
      uses: actions/upload-artifact@v2
      with:
        name: terraform-plan
        path: tfplan
```

### GitLab CI
```yaml
terraform:plan:
  stage: plan
  script:
    - terraform init
    - terraform validate
    - terraform fmt -check
    - terraform plan -out=tfplan
  artifacts:
    paths:
      - tfplan
    expire_in: 1 week
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Terraform Plan') {
            steps {
                sh 'terraform init'
                sh 'terraform validate'
                sh 'terraform fmt -check'
                sh 'terraform plan -out=tfplan'
                archiveArtifacts artifacts: 'tfplan', fingerprint: true
            }
        }
    }
}
```

### Plan Comparison
```bash
#!/bin/bash
# compare-plans.sh

# Generate baseline plan
terraform plan -out=baseline.tfplan

# Make changes and generate new plan
terraform plan -out=current.tfplan

# Compare plans (requires custom tooling)
terraform show -json baseline.tfplan > baseline.json
terraform show -json current.tfplan > current.json

# Use jq or other tools to compare
diff <(jq -S . baseline.json) <(jq -S . current.json)
```

## Best Practices

### Planning Workflow
- **Always Plan First**: Never apply without reviewing plan
- **Save Plans**: Store plans for audit and review
- **Review Changes**: Understand all planned modifications
- **Validate Inputs**: Check variables and configuration

### Code Quality
- **Format Consistently**: Use terraform fmt regularly
- **Validate Syntax**: Run terraform validate in CI/CD
- **Use Variables**: Parameterize configurations properly
- **Document Changes**: Include plan outputs in PRs

### Workspace Management
- **Environment Isolation**: One workspace per environment
- **Naming Conventions**: Consistent workspace naming
- **Access Control**: Restrict workspace operations
- **State Management**: Understand workspace state isolation

### Security Considerations
- **Plan Review**: Always review plans before apply
- **Sensitive Data**: Be careful with plan outputs
- **Access Control**: Limit who can create/modify plans
- **Audit Trail**: Maintain records of planned changes

### Performance Optimization
- **Targeted Planning**: Use -target for specific resources
- **Parallelism**: Adjust based on provider limits
- **Resource Limits**: Consider large infrastructure impacts
- **Caching**: Leverage provider and module caching

## Error Handling

### Common Planning Issues
```bash
# Invalid configuration
terraform validate
# Fix syntax errors before planning

# Missing providers
terraform init
# Initialize before planning

# Variable validation
terraform plan -var-file=terraform.tfvars
# Ensure all required variables are provided

# State lock conflicts
terraform force-unlock <lock-id>
# Only use when certain no other operations are running
```

### Troubleshooting
- **Debug Mode**: Set TF_LOG=DEBUG for detailed output
- **Plan Refresh**: Use -refresh-only to update state
- **Resource Targeting**: Use -target to isolate issues
- **Graph Analysis**: Use terraform graph to understand dependencies

IaC Plan provides essential infrastructure planning capabilities through proven Terraform CLI commands, enabling safe and reliable infrastructure management workflows.