# Infracost

Cloud cost estimates for Terraform projects.

## Description

Infracost shows cloud cost estimates for Terraform projects, helping engineers understand costs before making changes. It supports AWS, Azure, and Google Cloud across 400+ pricing resources. The tool integrates with CI/CD pipelines to provide cost breakdowns, diffs, and automated commenting on pull requests, enabling teams to make cost-conscious infrastructure decisions.

## MCP Tools

### Cost Analysis
- **`infracost_breakdown`** - Generate cost breakdown for Terraform projects
- **`infracost_diff`** - Show diff of monthly costs between current and planned state

### Output & Formatting
- **`infracost_output`** - Combine and output Infracost JSON files in different formats

### CI/CD Integration
- **`infracost_upload`** - Upload Infracost JSON file to Infracost Cloud
- **`infracost_comment_github`** - Post Infracost comment to GitHub pull requests

### Configuration & Setup
- **`infracost_configure`** - Set global configuration settings
- **`infracost_auth_login`** - Authenticate with Infracost Cloud
- **`infracost_generate_config`** - Generate Infracost config file from template

## Real CLI Commands Used

### Core Cost Commands
- `infracost breakdown --path <dir>` - Show cost breakdown
- `infracost breakdown --path <dir> --format json` - JSON output
- `infracost breakdown --terraform-var-file <file>` - Use variable files
- `infracost breakdown --terraform-workspace <name>` - Specify workspace
- `infracost breakdown --usage-file <file>` - Include usage estimates
- `infracost breakdown --show-skipped` - Show unsupported resources

### Cost Comparison
- `infracost diff --path <dir> --compare-to <baseline.json>` - Compare costs
- `infracost diff --format json` - JSON diff output
- `infracost diff --format diff` - Human-readable diff

### Output Processing
- `infracost output --path <files> --format table` - Table format
- `infracost output --path <files> --format html` - HTML report
- `infracost output --path <files> --format github-comment` - GitHub comment
- `infracost output --path <files> --format gitlab-comment` - GitLab comment
- `infracost output --path <files> --format slack-message` - Slack message

### CI/CD Integration
- `infracost upload --path <file>` - Upload to Infracost Cloud
- `infracost comment github --repo <repo> --pull-request <pr>` - GitHub comments

### Configuration
- `infracost configure set api_key <key>` - Set API key
- `infracost configure set currency <currency>` - Set currency
- `infracost auth login` - Interactive login
- `infracost generate config --repo-path <path>` - Generate config

## Use Cases

### Development Workflow
- **Cost Awareness**: See cost impact before applying changes
- **Budget Planning**: Understand infrastructure costs upfront
- **Resource Optimization**: Identify expensive resources
- **Team Education**: Learn about cloud costs during development

### CI/CD Pipeline Integration
- **Automated Cost Checks**: Include cost analysis in pipelines
- **Pull Request Comments**: Show cost changes in PR reviews
- **Cost Gates**: Block deployments exceeding budget thresholds
- **Historical Tracking**: Track cost trends over time

### Infrastructure Management
- **Multi-environment Costing**: Compare costs across environments
- **Resource Rightsizing**: Identify oversized resources
- **Usage-based Estimation**: Include realistic usage patterns
- **Cost Forecasting**: Project future infrastructure costs

### Team Collaboration
- **Cost Transparency**: Share cost information across teams
- **Budget Accountability**: Assign cost responsibility
- **Cost Reviews**: Include costs in architectural decisions
- **Financial Planning**: Support budget allocation processes

## Configuration Examples

### Basic Cost Breakdown
```bash
# Simple cost breakdown
infracost breakdown --path .

# Breakdown with variables
infracost breakdown --path . --terraform-var-file prod.tfvars

# JSON output for processing
infracost breakdown --path . --format json --out-file costs.json

# Include unsupported resources
infracost breakdown --path . --show-skipped
```

### Cost Comparison Workflow
```bash
# Generate baseline
infracost breakdown --path . --format json --out-file baseline.json

# Make infrastructure changes...

# Compare against baseline
infracost diff --path . --compare-to baseline.json

# Save diff for reporting
infracost diff --path . --compare-to baseline.json --out-file diff.json
```

### Multiple Output Formats
```bash
# HTML report
infracost output --path costs.json --format html --out-file report.html

# GitHub comment format
infracost output --path costs.json --format github-comment

# Table format for terminal
infracost output --path costs.json --format table

# Combine multiple JSON files
infracost output --path "costs*.json" --format json
```

### Configuration Setup
```bash
# Set API key
infracost configure set api_key your-api-key-here

# Set default currency
infracost configure set currency EUR

# Interactive authentication
infracost auth login

# Generate config from template
infracost generate config --repo-path . --template-path infracost.yml.tmpl
```

## Advanced Features

### Usage File Integration
```bash
# Generate usage file template
infracost breakdown --path . --usage-file usage.yml

# Use custom usage estimates
infracost breakdown --path . --usage-file custom-usage.yml

# Sync usage from Infracost Cloud
infracost breakdown --path . --sync-usage-file
```

### Terraform Integration
```bash
# Use with Terraform plan
terraform plan -out tfplan.binary
terraform show -json tfplan.binary > plan.json
infracost breakdown --path plan.json

# Multiple terraform variables
infracost breakdown --path . --terraform-var "instance_type=t3.large" --terraform-var "region=us-west-2"

# Specific workspace
infracost breakdown --path . --terraform-workspace production
```

### CI/CD Pipeline Examples
```yaml
# GitHub Actions
- name: Run Infracost
  run: |
    infracost breakdown --path . --format json --out-file infracost.json
    infracost comment github --repo ${{ github.repository }} --pull-request ${{ github.event.number }} --path infracost.json

# GitLab CI
infracost:
  script:
    - infracost breakdown --path . --format json --out-file infracost.json
    - infracost output --path infracost.json --format gitlab-comment
```

## Integration Patterns

### GitHub Actions Workflow
```yaml
name: Infracost
on: [pull_request]
jobs:
  infracost:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: infracost/actions/setup@v2
      with:
        api-key: ${{ secrets.INFRACOST_API_KEY }}
    
    - name: Generate baseline
      run: |
        infracost breakdown --path . --format json --out-file /tmp/baseline.json
        
    - name: Checkout PR branch
      uses: actions/checkout@v2
      
    - name: Generate cost diff
      run: |
        infracost diff --path . --compare-to /tmp/baseline.json --format json --out-file /tmp/diff.json
        
    - name: Post PR comment
      run: |
        infracost comment github \
          --repo ${{ github.repository }} \
          --pull-request ${{ github.event.number }} \
          --path /tmp/diff.json \
          --behavior update
```

### GitLab CI Integration
```yaml
.infracost: &infracost
  image: infracost/infracost:latest
  before_script:
    - infracost configure set api_key $INFRACOST_API_KEY

infracost:breakdown:
  <<: *infracost
  stage: plan
  script:
    - infracost breakdown --path . --format json --out-file costs.json
  artifacts:
    paths: [costs.json]

infracost:comment:
  <<: *infracost
  stage: deploy
  script:
    - infracost comment gitlab --repo $CI_PROJECT_PATH --merge-request $CI_MERGE_REQUEST_IID --path costs.json
  only: [merge_requests]
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    environment {
        INFRACOST_API_KEY = credentials('infracost-api-key')
    }
    stages {
        stage('Infracost') {
            steps {
                sh '''
                    infracost breakdown --path . --format json --out-file costs.json
                    infracost output --path costs.json --format table
                '''
                archiveArtifacts artifacts: 'costs.json'
            }
        }
    }
}
```

### Multi-Project Configuration
```yaml
# infracost.yml
version: 0.1
projects:
  - path: environments/dev
    name: development
    terraform_var_files:
      - dev.tfvars
    usage_file: usage-dev.yml
    
  - path: environments/staging
    name: staging
    terraform_var_files:
      - staging.tfvars
    usage_file: usage-staging.yml
    
  - path: environments/prod
    name: production
    terraform_var_files:
      - prod.tfvars
    usage_file: usage-prod.yml
```

## Best Practices

### Cost Management
- **Regular Monitoring**: Run cost checks on every infrastructure change
- **Budget Alerts**: Set up notifications for cost threshold breaches
- **Resource Tagging**: Use consistent tagging for cost allocation
- **Usage Estimation**: Include realistic usage patterns in estimates

### Team Workflow
- **PR Integration**: Always include cost comments in pull requests
- **Baseline Tracking**: Maintain cost baselines for comparison
- **Cost Reviews**: Include cost considerations in architecture reviews
- **Documentation**: Document cost assumptions and decisions

### Accuracy Improvement
- **Usage Files**: Create accurate usage profiles for resources
- **Variable Files**: Use production-like variables for estimates
- **Reserved Instances**: Account for commitment discounts
- **Data Transfer**: Include realistic data transfer estimates

### Security Considerations
- **API Key Management**: Store API keys securely in CI/CD systems
- **Access Control**: Limit who can modify cost configurations
- **Audit Trail**: Track cost-related configuration changes
- **Sensitive Data**: Avoid exposing sensitive resource details

## Error Handling

### Common Issues
```bash
# Missing API key
export INFRACOST_API_KEY=your-key-here
infracost breakdown --path .

# Invalid Terraform
terraform validate
infracost breakdown --path .

# Unsupported resources
infracost breakdown --path . --show-skipped

# Missing provider credentials
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
infracost breakdown --path .
```

### Troubleshooting
- **Debug Mode**: Set `INFRACOST_LOG_LEVEL=debug` for detailed output
- **Plan Validation**: Ensure Terraform plan is valid before cost analysis
- **Provider Setup**: Verify cloud provider credentials are configured
- **Version Compatibility**: Check Terraform version compatibility

Infracost provides essential cost visibility for infrastructure teams, enabling informed decisions about cloud spending before changes are deployed.