# Ship CLI Example Terraform Scenarios

This directory contains comprehensive Terraform examples designed to test Ship CLI's MCP server tools with Claude Code. Each scenario represents different infrastructure patterns with varying complexity levels and intentional issues for testing security scanning tools.

## 🏗️ Available Scenarios

### 1. **aws-web-app** - Production Web Application
- **Architecture**: VPC, ALB, Auto Scaling, RDS, S3
- **Complexity**: High
- **Purpose**: Test comprehensive infrastructure analysis
- **Features**: Multi-tier architecture, proper security practices
- **Cost**: ~$30-50/month

### 2. **serverless-api** - Serverless REST API
- **Architecture**: Lambda, API Gateway, DynamoDB, S3
- **Complexity**: Medium  
- **Purpose**: Test serverless security and cost analysis
- **Features**: Pay-per-use model, event-driven architecture
- **Cost**: ~$5-15/month

### 3. **security-hardened** - Security-Focused Infrastructure
- **Architecture**: KMS, VPC Flow Logs, CloudWatch, IAM
- **Complexity**: High
- **Purpose**: Test security scanning with intentional vulnerabilities
- **Features**: Encryption, logging, monitoring, **intentional security issues**
- **Cost**: ~$25-40/month

### 4. **cost-optimized** - Budget-Conscious Setup
- **Architecture**: Spot Fleet, Intelligent Tiering, Lifecycle Policies
- **Complexity**: Medium
- **Purpose**: Test cost analysis and optimization recommendations
- **Features**: Spot instances, storage optimization, **cost inefficiencies**
- **Cost**: ~$10-20/month

### 5. **multi-cloud-setup** - Multi-Cloud Deployment
- **Architecture**: AWS + Azure + GCP resources
- **Complexity**: High
- **Purpose**: Test multi-cloud analysis capabilities
- **Features**: Cross-cloud networking, provider comparison
- **Cost**: Varies by provider

### 6. **kubernetes-cluster** - Container Orchestration
- **Architecture**: EKS, Node Groups, ALB Ingress
- **Complexity**: High
- **Purpose**: Test Kubernetes security and networking
- **Features**: Container security, network policies
- **Cost**: ~$70-100/month

### 7. **data-pipeline** - Data Processing Infrastructure  
- **Architecture**: Kinesis, Lambda, Redshift, S3
- **Complexity**: Medium
- **Purpose**: Test data pipeline security and compliance
- **Features**: Stream processing, data warehousing
- **Cost**: ~$50-80/month

### 8. **compliance-ready** - Regulatory Compliance
- **Architecture**: Config Rules, CloudTrail, GuardDuty
- **Complexity**: High
- **Purpose**: Test compliance scanning (SOC2, HIPAA, etc.)
- **Features**: Audit logging, compliance controls
- **Cost**: ~$40-60/month

## 🚀 Testing with Ship CLI MCP Server

### Prerequisites

1. **Configure Claude Code MCP**: Copy `.mcp.json` to your Claude Code config
2. **Set OpenAI API Key**: Replace placeholder in `.mcp.json`
3. **AWS Credentials**: Ensure AWS CLI is configured

### Available MCP Tools

Once configured, Claude Code will have access to these Ship CLI tools:

#### 🔧 Terraform Analysis Tools
- **`terraform_lint`** - Syntax and best practices checking
- **`terraform_checkov_scan`** - Security policy compliance
- **`terraform_security_scan`** - Alternative security scanning
- **`terraform_cost_analysis`** - Infrastructure cost analysis
- **`terraform_generate_docs`** - Documentation generation
- **`terraform_generate_diagram`** - Infrastructure visualization

#### 🤖 AI-Powered Investigation
- **`ai_investigate`** - Natural language infrastructure queries

#### ☁️ Cloud Operations
- **`cloudship_push`** - Upload artifacts for AI analysis

### Example Claude Code Conversations

#### Security Analysis
```
"Please analyze the security-hardened example for vulnerabilities using the terraform security scanning tools"
```

#### Cost Optimization  
```
"Review the cost-optimized example and identify opportunities for further cost savings"
```

#### Documentation Generation
```
"Generate comprehensive documentation for the aws-web-app example"
```

#### AI Investigation
```
"Use ai_investigate to show me all S3 buckets in my AWS account and check which ones have public access"
```

## 📁 Example Structure

Each example includes:
- **`main.tf`** - Primary infrastructure definition
- **`variables.tf`** - Input variables with defaults
- **`outputs.tf`** - Output values for integration
- **`README.md`** - Detailed usage instructions
- **Supporting files** - User data scripts, configs, etc.

## 🧪 Testing Scenarios

### Security Testing
1. **Run security scans** on `security-hardened/` - should find intentional issues
2. **Compare results** between `terraform_checkov_scan` and `terraform_security_scan`
3. **Verify fixes** - remediate issues and re-scan

### Cost Analysis Testing
1. **Analyze costs** for `cost-optimized/` - should identify inefficiencies  
2. **Compare scenarios** - web-app vs serverless-api cost profiles
3. **Optimization recommendations** - ask AI for cost reduction strategies

### Documentation Testing
1. **Generate docs** for complex scenarios like `aws-web-app`
2. **Verify completeness** - ensure all resources are documented
3. **Format testing** - try different output formats (markdown, table)

### Multi-Tool Workflows
1. **Full pipeline**: lint → security scan → cost analysis → docs
2. **Push results** to Cloudship using `--push` flag
3. **AI analysis** of existing infrastructure with `ai_investigate`

## ⚠️ Important Notes

### Security Warnings
- **Never deploy** these examples to production without review
- **Intentional vulnerabilities** exist in some examples for testing
- **Review all security settings** before any real deployment

### Cost Management
- **Monitor costs** if deploying to AWS
- **Use terraform destroy** when testing is complete
- **Set billing alerts** for safety

### MCP Configuration
- **.mcp.json is gitignored** - contains personal credentials
- **Replace placeholder values** with real API keys
- **Test connectivity** before extensive use

## 🎯 Best Practices

1. **Start simple** - test with `serverless-api` first
2. **Use plan files** - generate with `terraform plan -out=tf.plan`
3. **Test incrementally** - one tool at a time initially
4. **Compare outputs** - manual vs MCP tool results
5. **Document findings** - note any tool limitations or issues

## 🔧 Troubleshooting

### Common Issues
- **MCP server connection** - check .mcp.json configuration
- **AWS credentials** - ensure proper AWS CLI setup
- **Tool timeouts** - adjust timeout values in tool calls
- **Permission errors** - verify IAM permissions for analysis

### Getting Help
- **Check tool outputs** - most tools provide detailed error messages
- **Test individual tools** - isolate issues by testing one tool at a time
- **Review logs** - Ship CLI provides detailed logging with `--log-level debug`

---

These examples provide comprehensive testing scenarios for the Ship CLI MCP server integration with Claude Code. Each scenario tests different aspects of infrastructure analysis, security scanning, cost optimization, and documentation generation. 🚀