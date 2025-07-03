# Ship CLI Demo Gallery

> Interactive demonstrations of Ship CLI features - Infrastructure as Code analysis and automation

## 🚀 Quick Start
![Quick Start Demo](./ship-quick-start.gif)
*Get started with Ship CLI in under a minute*

## Core Features

### 🔍 Steampipe Query Engine
![Steampipe Query Demo](./ship-query.gif)
*Query your cloud infrastructure using SQL with live data*

```bash
ship query 'SELECT * FROM aws_ec2_instance WHERE instance_state = "running"'
```

### 🤖 AI-Powered Investigation
![AI Investigation Demo](./ship-ai-investigate.gif)
*Use natural language to investigate infrastructure issues*

```bash
ship ai-investigate --prompt "Find security vulnerabilities in my AWS setup" --execute
```

### 🤖 Autonomous AI Agent
![AI Agent Demo](./ship-ai-agent.gif)
*Run complex multi-tool investigations autonomously*

```bash
ship ai-agent --task "Perform complete security audit and cost optimization"
```

## Terraform Tools Suite

### 📊 Complete Terraform Analysis
![All Terraform Tools Demo](./ship-terraform-all-tools.gif)
*Run the full suite of Terraform analysis tools*

Individual tools available:
- `ship terraform-tools lint` - Lint your Terraform code
- `ship terraform-tools checkov-scan` - Security scanning
- `ship terraform-tools cost-estimate` - Cost estimation with Infracost
- `ship terraform-tools generate-diagram` - Visualize infrastructure
- `ship terraform-tools generate-docs` - Auto-generate documentation

## CloudShip Integration

### 🔐 Authentication
![Authentication Demo](./ship-auth.gif)
*Connect to CloudShip for centralized analysis*

```bash
ship auth --api-key YOUR_API_KEY
```

### 📤 Push Artifacts
![Push Demo](./ship-push.gif)
*Upload Terraform plans and analysis results*

```bash
ship push tfplan.binary --tags env=prod,team=platform
```

### 🔌 MCP Server
![MCP Server Demo](./ship-mcp.gif)
*Integrate Ship tools with LLMs via Model Context Protocol*

```bash
ship mcp  # Start MCP server for Claude Desktop or other LLM tools
```

## Feature Highlights

### 🎯 Key Capabilities
- **No Local Installation Required** - All tools run in containers via Dagger
- **Multi-Cloud Support** - AWS, Azure, and GCP
- **AI-Powered** - Natural language infrastructure queries
- **Security First** - Multiple security scanners integrated
- **Cost Optimization** - Real-time cost analysis and recommendations
- **Developer Friendly** - Simple CLI with sensible defaults

### 🛠️ Available Commands

```
ship
├── auth                  # Manage CloudShip authentication
├── push                  # Upload artifacts to CloudShip
├── query                 # Run Steampipe SQL queries
├── ai-investigate        # AI-powered infrastructure investigation
├── ai-agent             # Autonomous AI agent with multiple tools
├── mcp                  # Start MCP server for LLM integration
└── terraform-tools      # Terraform analysis suite
    ├── cost-analysis    # Analyze costs with OpenInfraQuote
    ├── security-scan    # Security scanning with Trivy
    ├── generate-docs    # Auto-generate documentation
    ├── lint            # Lint with TFLint
    ├── checkov-scan    # Security scan with Checkov
    ├── cost-estimate   # Estimate costs with Infracost
    └── generate-diagram # Visualize with InfraMap
```

## Getting Started

1. **Install Ship CLI**
   ```bash
   curl -sSL https://ship.cloudshipai.com/install | bash
   ```

2. **Set up authentication (optional)**
   ```bash
   ship auth --api-key YOUR_API_KEY
   ```

3. **Run your first query**
   ```bash
   ship query 'SELECT COUNT(*) FROM aws_s3_bucket'
   ```

4. **Try AI investigation**
   ```bash
   ship ai-investigate --prompt "What's my AWS monthly spend?" --execute
   ```

## Learn More

- 📚 [Full Documentation](https://docs.cloudshipai.com)
- 🐛 [Report Issues](https://github.com/cloudshipai/ship/issues)
- 💬 [Join our Discord](https://discord.gg/cloudship)
- 🌟 [Star on GitHub](https://github.com/cloudshipai/ship)

---

*Ship CLI - Making infrastructure analysis accessible to everyone* 🚢