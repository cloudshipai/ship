# Ship CLI MCP Server Testing Guide

## 🎯 Quick Start

1. **Configure Claude Code MCP Server**:
   - Copy the `.mcp.json` from this directory to your Claude Code config
   - Replace `sk-proj-your-openai-key-here-replace-this` with your real OpenAI API key
   - Save and restart Claude Code

2. **Test Basic Connectivity**:
   ```
   "Please use the terraform_lint tool to check the aws-web-app example"
   ```

3. **Verify All Tools Work**:
   ```
   "Show me all available Ship CLI tools and their capabilities"
   ```

## 🛠️ MCP Tools Available

### Terraform Analysis Tools (6 tools)
- **terraform_lint** - TFLint syntax checking
- **terraform_checkov_scan** - Checkov security compliance  
- **terraform_security_scan** - Trivy security scanning
- **terraform_cost_analysis** - OpenInfraQuote cost analysis
- **terraform_generate_docs** - terraform-docs documentation
- **terraform_generate_diagram** - Infrastructure diagrams

### AI Investigation (1 tool)
- **ai_investigate** - Natural language infrastructure queries

### Cloud Operations (1 tool)  
- **cloudship_push** - Upload artifacts for analysis

## 🧪 Test Scenarios

### 1. Security Analysis Testing
```
"Please run both terraform_checkov_scan and terraform_security_scan on the security-hardened example and compare the results"
```

**Expected**: Should find intentional security issues like:
- SSH access from 0.0.0.0/0
- Missing S3 encryption
- Hardcoded secrets in user data
- Overly broad IAM permissions

### 2. Cost Analysis Testing
```
"Analyze the cost-optimized example using terraform_cost_analysis and identify potential cost savings"
```

**Expected**: Should identify:
- Over-provisioned RDS storage
- Suboptimal storage types
- Excessive backup retention
- Spot fleet configuration issues

### 3. Documentation Generation
```
"Generate comprehensive documentation for the aws-web-app example using terraform_generate_docs"
```

**Expected**: Should produce:
- Resource table with types and names
- Input variables with descriptions
- Output values and descriptions
- Provider requirements

### 4. Multi-Tool Workflow
```
"Please run a complete analysis of the serverless-api example: lint it, scan for security issues, and generate documentation"
```

**Expected**: Sequential execution of multiple tools with comprehensive output. (Push to Cloudship disabled during staging)

### 5. AI Infrastructure Investigation
```
"Use ai_investigate to show me all my AWS S3 buckets and check which ones have public access or are unencrypted"
```

**Expected**: Real-time analysis of your AWS account with security recommendations.

## 🔍 Testing Checklist

### Basic Functionality
- [ ] MCP server starts successfully
- [ ] Claude Code can connect to Ship CLI
- [ ] All 8 tools are available and accessible
- [ ] Tool parameters work correctly

### Terraform Tools
- [ ] `terraform_lint` - Check syntax and best practices
- [ ] `terraform_checkov_scan` - Security policy validation
- [ ] `terraform_security_scan` - Alternative security scanning  
- [ ] `terraform_cost_analysis` - Cost estimation
- [ ] `terraform_generate_docs` - Documentation creation
- [ ] `terraform_generate_diagram` - Visual diagrams

### AI Features
- [ ] `ai_investigate` - Natural language queries work
- [ ] AWS credential detection works
- [ ] Steampipe queries execute successfully
- [ ] Results are comprehensive and actionable

### Integration Features
- [ ] `cloudship_push` - Artifact upload works
- [ ] `--push` flags work with terraform tools
- [ ] Error handling is graceful
- [ ] Tool timeouts work appropriately

## 🐛 Troubleshooting

### MCP Connection Issues
```
Error: Cannot connect to Ship CLI MCP server
```
**Fix**: 
1. Check `.mcp.json` syntax is valid
2. Ensure `go` is in PATH
3. Verify working directory path in config
4. Restart Claude Code

### AWS Credentials Issues
```
Error: AWS credentials not found
```
**Fix**:
1. Run `aws sts get-caller-identity` to verify credentials
2. Check AWS_PROFILE in `.mcp.json` matches your profile
3. Ensure AWS_REGION is set correctly

### OpenAI API Issues
```
Error: OpenAI API key required
```
**Fix**:
1. Replace placeholder in `.mcp.json` with real API key
2. Verify API key has sufficient credits
3. Check API key permissions

### Tool Execution Issues
```
Error: Tool execution failed
```
**Fix**:
1. Check Docker is running (required for some tools)
2. Verify Terraform files are valid
3. Check file paths are correct
4. Review tool parameters

## 📊 Expected Results

### Security Scans
- **security-hardened**: ~5-10 findings including SSH, encryption, secrets
- **aws-web-app**: ~2-5 findings (minor security improvements)
- **serverless-api**: ~1-3 findings (serverless-specific issues)

### Cost Analysis
- **cost-optimized**: ~$10-20/month with optimization opportunities
- **aws-web-app**: ~$30-50/month baseline infrastructure
- **serverless-api**: ~$5-15/month pay-per-use model

### Documentation
- **aws-web-app**: ~25+ resources documented
- **serverless-api**: ~10+ resources documented  
- **All examples**: Complete variable and output documentation

## 🚀 Advanced Testing

### Custom Scenarios
1. **Create your own Terraform** and test tools against it
2. **Test with real infrastructure** using `ai_investigate`
3. **Compare tool outputs** with manual analysis
4. **Test error conditions** with invalid Terraform

### Performance Testing
1. **Large Terraform files** - test tool scalability
2. **Multiple simultaneous calls** - test concurrency
3. **Long-running investigations** - test timeout handling
4. **Network issues** - test retry mechanisms

### Integration Testing
1. **End-to-end workflows** - scan → fix → re-scan
2. **Multi-environment testing** - dev/staging/prod configs
3. **CI/CD integration** - automate tool execution
4. **Compliance workflows** - full audit trails

## 📈 Success Metrics

### Functionality
- ✅ All 8 tools execute successfully
- ✅ Results are accurate and actionable
- ✅ Error handling is graceful
- ✅ Performance is acceptable

### User Experience
- ✅ Natural language interaction works well
- ✅ Tool selection is appropriate for requests
- ✅ Results are well-formatted and readable
- ✅ Workflows are intuitive and efficient

### Integration
- ✅ MCP protocol implementation is stable
- ✅ Claude Code integration is seamless
- ✅ AWS credential handling is reliable
- ✅ Cloudship integration works correctly

---

## 🎉 Ready to Test!

Your Ship CLI MCP server is ready for testing with Claude Code. Start with simple tool tests and work your way up to complex multi-tool workflows. The combination of terraform-tools and AI investigation provides powerful infrastructure analysis capabilities directly in your AI assistant! 🚀