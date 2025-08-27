# Ship Tools Quick Reference

This table provides a quick overview of all Ship tools organized by workflow and use case.

## Quick Reference by Workflow

| Workflow | Primary Tools | Supporting Tools | When to Use |
|----------|---------------|------------------|-------------|
| **Terraform Development** | `tflint`, `terraform-docs`, `checkov`, `tfsec` | `inframap`, `terrascan`, `openinfraquote` | Terraform syntax → Security → Cost → Documentation |
| **Container Security** | `trivy`, `grype`, `syft` | `dockle`, `cosign`, `hadolint` | Dockerfile → Build → Scan → Sign |
| **Secret Management** | `trufflehog`, `gitleaks`, `git-secrets` | `sops`, `history-scrub` | Fast CI checks vs Deep audits |
| **Kubernetes Operations** | `kubescape`, `kube-bench`, `velero` | `falco`, `kyverno`, `goldilocks` | Assessment → Compliance → Backup → Monitor |
| **Cloud Security** | `prowler`, `scout-suite`, `steampipe` | `cloudquery`, `custodian` | Multi-cloud security posture |
| **Web Application Testing** | `nuclei`, `zap`, `nikto` | `nmap` | Template-based vs Manual testing |
| **Development & CI/CD** | `semgrep`, `actionlint`, `gitleaks` | `hadolint`, `conftest` | SAST → Workflow validation → Fast scans |
| **Supply Chain Security** | `syft`, `cosign`, `dependency-track` | `ossf-scorecard` | SBOM → Signing → Analysis |
| **Compliance & Governance** | `checkov`, `conftest`, `openscap` | `gatekeeper`, `kyverno` | Policy validation → Enforcement |

## All Tools by Category

### Terraform Tools (11 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `lint` | General Terraform linting | Basic syntax validation |
| `tflint` | Advanced Terraform linting | Comprehensive Terraform validation |
| `cost` | Infrastructure cost estimation | Budget planning |
| `docs` | Generate documentation | Module documentation |
| `diagram` | Infrastructure visualization | Architecture diagrams |
| `terraform-docs` | Advanced module documentation | Detailed module publishing |
| `openinfraquote` | Detailed cost analysis | Financial impact assessment |
| `aws-pricing-builtin` | AWS-specific pricing | AWS cost optimization |

### Security Tools (31 tools)

#### Vulnerability Scanning
| Tool | Purpose | Use When |
|------|---------|----------|
| `trivy` | Universal vulnerability scanner | Containers, filesystems, git repos |
| `grype` | Container-focused with SBOM | Detailed container analysis |
| `osv-scanner` | OSS vulnerability scanning | Official OSSF tool for dependencies |
| `dependency-track` | SBOM analysis platform | Enterprise vulnerability management |
| `checkov` | Multi-framework IaC security | Infrastructure as Code security scanning |

#### Secret Detection
| Tool | Purpose | Use When |
|------|---------|----------|
| `trufflehog` | Advanced secret scanning + verification | Deep audits, reduces false positives |
| `gitleaks` | Fast git secret scanning | CI/CD pipelines, pre-commit hooks |
| `git-secrets` | AWS-specific secret prevention | AWS environments, prevention focus |
| `sops` | Secrets management | Encrypted secrets in git |

#### Container Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `syft` | SBOM generation | Software inventory, supply chain |
| `dockle` | Container image linting | Dockerfile best practices |
| `cosign` | Container signing/verification | Supply chain integrity |
| `hadolint` | Dockerfile linter | Build-time Dockerfile validation |

#### Kubernetes Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `kubescape` | K8s security platform | Comprehensive cluster assessment |
| `kube-bench` | CIS K8s benchmarks | Compliance validation |
| `kube-hunter` | K8s penetration testing | Security testing |
| `falco` | Runtime security monitoring | Production monitoring |
| `velero` | K8s backup/disaster recovery | Production cluster protection |
| `goldilocks` | Resource recommendations | Cost/performance optimization |

#### Web Application Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `zap` | Web app penetration testing | Comprehensive web security |
| `nikto` | Web server security | Server misconfiguration |

#### Cloud Security & Compliance
| Tool | Purpose | Use When |
|------|---------|----------|
| `prowler` | Multi-cloud security assessment | AWS/GCP/Azure security |
| `scout-suite` | Multi-cloud security auditing | Cloud security posture |
| `powerpipe` | Infrastructure benchmarking | Compliance and security benchmarks |
| `openscap` | Security compliance scanning | SCAP compliance validation |
| `ossf-scorecard` | Open source project security | Supply chain security assessment |

#### Static Analysis & Code Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `semgrep` | Static analysis (SAST) | Code-level security |
| `actionlint` | GitHub Actions linting | Workflow validation |
| `cfn-nag` | CloudFormation security | AWS CloudFormation |
| `conftest` | OPA policy testing | Policy validation |
| `gatekeeper` | OPA policy enforcement | K8s policy enforcement |
| `license-detector` | Software license detection | License compliance |
| `iac-plan` | Infrastructure as Code planning | IaC security planning |
| `tfsec` | Terraform security scanning | Terraform-specific security |

#### Infrastructure & Platform Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `terrascan` | Infrastructure as Code security | Terraform/CloudFormation scanning |



### Cloud Infrastructure Tools (17 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `cloudquery` | Cloud asset inventory | Multi-cloud asset discovery |
| `custodian` | Cloud governance automation | Policy enforcement |
| `terraformer` | Infrastructure import | Import existing infra to Terraform |
| `inframap` | Infrastructure visualization | Architecture documentation |
| `infrascan` | Infrastructure security scanning | Security assessment |
| `aws-iam-rotation` | AWS IAM key rotation | Security automation |
| `tfstate-reader` | Terraform state analysis | State file inspection |
| `packer` | Machine image building | Golden image creation |
| `fleet` | GitOps for Kubernetes | K8s GitOps delivery |
| `kuttl` | Kubernetes testing | K8s test automation |
| `litmus` | Chaos engineering | Reliability testing |
| `cert-manager` | Certificate management | TLS automation |
| `k8s-network-policy` | Network policy management | K8s network security |
| `kyverno` | Kubernetes policy engine | Policy management |
| `kyverno-multitenant` | Multi-tenant K8s policies | Tenant isolation |
| `github-admin` | GitHub administration | Repository management |
| `github-packages` | GitHub packages management | Package operations |

### AWS Tools (4 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `cloudsplaining` | AWS IAM security assessment | IAM policy analysis |
| `parliament` | AWS IAM policy linting | IAM policy validation |
| `pmapper` | AWS IAM privilege mapping | Privilege escalation analysis |
| `policy-sentry` | AWS IAM policy generation | Least-privilege policies |


## External MCP Servers (16 servers)

### Development & CI/CD
| Server | Purpose | Variables Needed |
|--------|---------|------------------|
| `bitbucket` | Repository management | USERNAME, APP_PASSWORD, WORKSPACE |
| `trello` | Project management | API_KEY, TOKEN |
| `github` | GitHub operations | PERSONAL_ACCESS_TOKEN |

### Browser Automation
| Server | Purpose | Variables Needed |
|--------|---------|------------------|
| `playwright` | Browser automation | BROWSER, HEADLESS |

### Database Integration
| Server | Purpose | Variables Needed |
|--------|---------|------------------|
| `postgresql` | Database operations | CONNECTION_STRING |
| `supabase` | Backend-as-a-service | ACCESS_TOKEN, PROJECT_REF |

### Infrastructure & Monitoring
| Server | Purpose | Variables Needed |
|--------|---------|------------------|
| `steampipe` | Cloud asset queries | DATABASE_CONNECTIONS |
| `grafana` | Monitoring integration | URL, API_KEY |

### Productivity & Storage
| Server | Purpose | Variables Needed |
|--------|---------|------------------|
| `filesystem` | File operations | FILESYSTEM_ROOT |
| `memory` | Knowledge storage | STORAGE_PATH, MAX_SIZE |
| `brave-search` | Web search | API_KEY |
| `slack` | Team communication | Multiple token options |
| `desktop-commander` | Desktop operations | ROOT, SAFE_MODE |

### AWS Operations (6 servers)
| Server | Purpose | Variables Needed |
|--------|---------|------------------|
| `aws-core` | General AWS operations | PROFILE, REGION |
| `aws-iam` | IAM management | PROFILE, REGION |
| `aws-pricing` | Cost analysis | PROFILE, REGION |
| `aws-eks` | EKS operations | PROFILE, REGION |
| `aws-ec2` | EC2 operations | PROFILE, REGION |
| `aws-s3` | S3 operations | PROFILE, REGION |

## Usage Patterns

### CI/CD Pipeline Integration
```
Code → tflint → checkov → gitleaks → trivy → Deploy
    ↓
terraform-docs → openinfraquote (cost-sensitive)
```

### Security Review Workflow
```
Weekly: prowler + scout-suite (cloud)
Daily: trivy + grype (containers)
Per PR: tflint + checkov + gitleaks
Monthly: trufflehog + ossf-scorecard
```

### Multi-Tool Security Stack
```
Build: hadolint → trivy → syft → cosign
Runtime: falco → kubescape
Compliance: kube-bench → openscap
```

## Tool Selection Guide

### **Choose Trivy if**: You want one tool for most vulnerability scanning
### **Choose Grype if**: You need detailed SBOM analysis
### **Choose Checkov if**: You use multiple IaC frameworks  
### **Choose TFSec if**: You're Terraform-focused
### **Choose TruffleHog if**: Accuracy and verification are critical
### **Choose Gitleaks if**: Speed and CI/CD integration are priorities
### **Choose ZAP if**: You need comprehensive web app testing
### **Choose Nuclei if**: You want fast, automated template-based scanning

---

**Total: 79 tools** (63 built-in + 16 external MCPs)

### Summary by Category
- **Terraform Tools**: 11 tools
- **Security Tools**: 31 tools  
- **Cloud Infrastructure Tools**: 17 tools
- **AWS Tools**: 4 tools
- **Total Built-in Tools**: 63 tools
- **External MCP Servers**: 16 servers