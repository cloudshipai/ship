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

### Terraform Tools (7 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `tflint` | Terraform syntax and best practices | Before every Terraform apply |
| `terraform-docs` | Generate module documentation | Team collaboration, module publishing |
| `checkov` | Multi-framework IaC security | Comprehensive security scanning |
| `tfsec` | Terraform-specific security | Deep Terraform security analysis |
| `terrascan` | Policy-as-code security | Custom OPA policies |
| `inframap` | Infrastructure visualization | Architecture documentation |
| `openinfraquote` | Cost estimation | Before expensive infrastructure changes |

### Security Tools (41 tools)

#### Vulnerability Scanning
| Tool | Purpose | Use When |
|------|---------|----------|
| `trivy` | Universal vulnerability scanner | Containers, filesystems, git repos |
| `grype` | Container-focused with SBOM | Detailed container analysis |
| `nuclei` | Fast template-based scanning | Web apps, infrastructure misconfigs |
| `osv-scanner` | OSS vulnerability scanning | Official OSSF tool for dependencies |

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
| `nmap` | Network discovery/scanning | Network reconnaissance |

#### Cloud Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `prowler` | Multi-cloud security assessment | AWS/GCP/Azure security |
| `scout-suite` | Cloud configuration auditing | In-depth config analysis |
| `steampipe` | Cloud asset querying (SQL) | Custom cloud queries |
| `cloudquery` | Cloud asset inventory | Asset management |
| `custodian` | Cloud governance automation | Policy enforcement |

#### Static Analysis & Code Security
| Tool | Purpose | Use When |
|------|---------|----------|
| `semgrep` | Static analysis (SAST) | Code-level security |
| `actionlint` | GitHub Actions linting | Workflow validation |
| `cfn-nag` | CloudFormation security | AWS CloudFormation |
| `conftest` | OPA policy testing | Policy validation |
| `gatekeeper` | OPA policy enforcement | K8s policy enforcement |
| `kyverno` | K8s policy management | K8s native policies |

#### Compliance & Governance
| Tool | Purpose | Use When |
|------|---------|----------|
| `openscap` | Security compliance scanning | SCAP compliance |
| `ossf-scorecard` | OSS project security scoring | Dependency evaluation |
| `scancode` | License detection | License compliance |
| `license-detector` | Software license detection | Legal compliance |

### Cloud Infrastructure Tools (11 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `packer` | Machine image building | Golden image creation |
| `terraformer` | Infrastructure import | Import existing infra to Terraform |
| `tfstate-reader` | Terraform state analysis | State file analysis |
| `iac-plan` | IaC planning | Infrastructure planning |
| `infrascan` | Infrastructure security | Infrastructure scanning |
| `powerpipe` | Infrastructure benchmarking | Compliance benchmarking |
| `fleet` | GitOps for Kubernetes | K8s GitOps |
| `kuttl` | Kubernetes testing | K8s test framework |
| `litmus` | Chaos engineering | K8s resilience testing |
| `cert-manager` | Certificate management | K8s certificate automation |
| `k8s-network-policy` | Network policy management | K8s network security |

### AWS Tools (4 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `cloudsplaining` | AWS IAM security assessment | IAM policy analysis |
| `parliament` | AWS IAM policy linting | IAM policy validation |
| `pmapper` | AWS IAM privilege mapping | Privilege escalation analysis |
| `policy-sentry` | AWS IAM policy generation | Least-privilege policies |

### Supply Chain Tools (3 tools)
| Tool | Purpose | Use When |
|------|---------|----------|
| `cosign` | Container signing/verification | Supply chain integrity |
| `dependency-track` | SBOM analysis platform | Enterprise SBOM management |
| `syft` | SBOM generation | Software bill of materials |

### Development Tools (1 tool)
| Tool | Purpose | Use When |
|------|---------|----------|
| `opencode` | AI coding assistant | AI-powered development |

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

**Total: 82 tools** (63 built-in + 16 external MCPs + 3 collections)