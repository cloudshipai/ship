# Ship DevOps Tools: Use Cases and When to Use Each Tool

This guide explains what each tool in Ship's comprehensive DevOps toolkit does and when you should use them. Ship provides 75+ tools across security, infrastructure, Kubernetes, cloud, and development workflows.

## Quick Reference by Workflow

| Workflow | Primary Tools | Supporting Tools |
|----------|---------------|------------------|
| **Terraform Development** | `tflint`, `terraform-docs`, `checkov`, `tfsec` | `inframap`, `terrascan`, `openinfraquote` |
| **Container Security** | `trivy`, `grype`, `syft` | `dockle`, `cosign` |
| **Secret Management** | `trufflehog`, `gitleaks`, `git-secrets` | `sops`, `history-scrub` |
| **Kubernetes Operations** | `kubescape`, `kube-bench`, `velero` | `falco`, `kyverno`, `goldilocks` |
| **Cloud Security** | `prowler`, `scout-suite`, `steampipe` | `cloudquery`, `custodian` |
| **Web Application Testing** | `nuclei`, `zap`, `nikto` | `nmap` |
| **Compliance & Governance** | `checkov`, `conftest`, `openscap` | `gatekeeper`, `allstar` |

## Security Tools

### Vulnerability Scanning

#### **Trivy** - Universal Vulnerability Scanner
**When to use:** Your go-to scanner for containers, filesystems, git repos, and Kubernetes
- **Use cases:** 
  - Scan Docker images before deployment
  - Audit filesystem for vulnerabilities in CI/CD
  - Check git repositories for embedded vulnerabilities
  - Kubernetes cluster security scanning
- **Example:** `trivy_scan_image` for production container validation

#### **Grype** - Container-Focused Vulnerability Scanner  
**When to use:** When you need detailed SBOM-based container vulnerability analysis
- **Use cases:**
  - Deep container image analysis with SBOM correlation
  - CVE tracking and management
  - Supply chain security validation
- **Complements:** Trivy (Grype focuses more on SBOMs, Trivy is broader)

#### **Nuclei** - Fast Vulnerability Scanner with Templates
**When to use:** Web application and infrastructure vulnerability discovery
- **Use cases:**
  - Web app security testing with community templates
  - Infrastructure misconfiguration detection
  - Custom vulnerability template development
  - Automated security testing in CI/CD
- **Example:** Scan web apps for OWASP Top 10 vulnerabilities

### Secret Detection

#### **TruffleHog** - Advanced Secret Scanning with Verification
**When to use:** When you need comprehensive secret detection with verification
- **Use cases:**
  - Verify found secrets are actually valid (reduces false positives)
  - Deep git history scanning
  - Enterprise-grade secret detection
  - Integration with secret management platforms
- **Example:** Pre-commit hooks with secret verification

#### **Gitleaks** - Fast Git Secret Scanning
**When to use:** When you need speed and efficiency in CI/CD pipelines
- **Use cases:**
  - Fast pre-commit hooks
  - CI/CD pipeline integration (faster than TruffleHog)
  - Basic secret pattern detection
  - Git repository auditing
- **Example:** Quick CI checks before merge

#### **git-secrets** - AWS-Focused Secret Prevention  
**When to use:** AWS-centric environments with emphasis on prevention
- **Use cases:**
  - Prevent AWS credentials from being committed
  - Install git hooks for AWS secret prevention
  - AWS-specific patterns and validation
- **Complements:** TruffleHog/Gitleaks (AWS-specific focus)

### Infrastructure as Code Security

#### **Checkov** - Multi-Cloud IaC Security Scanner
**When to use:** Comprehensive security scanning across multiple IaC frameworks
- **Use cases:**
  - Terraform, CloudFormation, Kubernetes YAML security scanning
  - Multi-cloud security policy enforcement
  - Developer-friendly security feedback
  - Custom policy development
- **Example:** Scan Terraform before applying changes

#### **TFSec** - Terraform-Specific Security Scanner
**When to use:** Terraform-focused security with detailed Terraform knowledge
- **Use cases:**
  - Deep Terraform-specific security analysis
  - Terraform Cloud integration
  - Custom Terraform security rules
- **Complements:** Checkov (TFSec is Terraform-specific, Checkov is broader)

#### **Terrascan** - Policy-as-Code Security Scanner
**When to use:** When you need OPA/Rego-based custom policies
- **Use cases:**
  - Custom compliance policies using OPA
  - Enterprise policy enforcement
  - Multi-framework IaC scanning
- **Complements:** Checkov/TFSec (policy-focused approach)

## Terraform Development Tools

#### **TFLint** - Terraform Linter
**When to use:** Terraform syntax, best practices, and provider-specific issues
- **Use cases:**
  - Catch Terraform syntax errors before apply
  - Enforce Terraform best practices
  - Provider-specific validations (AWS, GCP, Azure)
  - Integration with Terraform Cloud/Enterprise
- **Different from security scanners:** Focuses on syntax/best practices, not security
- **Example:** Validate Terraform files in IDE or CI/CD

#### **terraform-docs** - Terraform Documentation Generator
**When to use:** Automatically generate and maintain Terraform module documentation
- **Use cases:**
  - Generate README.md for Terraform modules
  - Keep documentation in sync with code
  - Team collaboration on Terraform modules
  - Terraform registry publication
- **Example:** Auto-update module docs in PR workflows

#### **OpenInfraQuote** - Infrastructure Cost Estimation
**When to use:** Understand cost implications of infrastructure changes
- **Use cases:**
  - Cost estimation before Terraform apply
  - Budget planning for infrastructure changes
  - Cost optimization analysis
  - FinOps integration in DevOps workflows
- **Example:** Show cost diff in pull requests

#### **InfraMap** - Infrastructure Visualization
**When to use:** Visual representation of Terraform infrastructure
- **Use cases:**
  - Create architecture diagrams from Terraform
  - Infrastructure documentation
  - Team communication about infrastructure
  - Architecture reviews
- **Example:** Generate diagrams for design documents

## Kubernetes Tools

#### **Kubescape** - Kubernetes Security Platform
**When to use:** Comprehensive Kubernetes cluster security assessment
- **Use cases:**
  - CIS Kubernetes benchmark compliance
  - RBAC analysis and optimization
  - Kubernetes security posture management
  - Compliance reporting (SOC2, PCI DSS)
- **Example:** Regular cluster security health checks

#### **Kube-bench** - CIS Kubernetes Benchmark
**When to use:** Validate Kubernetes cluster against CIS benchmarks
- **Use cases:**
  - Compliance auditing (CIS benchmarks)
  - Kubernetes hardening validation
  - Security baseline establishment
- **Complements:** Kubescape (kube-bench is CIS-focused, Kubescape is broader)

#### **Velero** - Kubernetes Backup and Disaster Recovery
**When to use:** Production Kubernetes cluster backup and migration
- **Use cases:**
  - Kubernetes cluster backup and restore
  - Cluster migration between environments
  - Disaster recovery planning
  - Application-level backup strategies
- **Example:** Daily production cluster backups

#### **Goldilocks** - Kubernetes Resource Recommendations
**When to use:** Optimize Kubernetes resource requests and limits
- **Use cases:**
  - Right-size pod resource allocations
  - Cost optimization through better resource utilization
  - Performance optimization
  - Capacity planning
- **Example:** Monthly resource optimization reviews

#### **Falco** - Runtime Security Monitoring
**When to use:** Real-time security monitoring of Kubernetes workloads
- **Use cases:**
  - Detect anomalous behavior in running containers
  - Runtime security monitoring
  - Compliance monitoring (runtime)
  - Security incident detection
- **Example:** Alert on suspicious container behavior

## Cloud Security Tools

#### **Prowler** - Multi-Cloud Security Assessment
**When to use:** Comprehensive cloud security posture assessment
- **Use cases:**
  - AWS, GCP, Azure security assessment
  - Cloud compliance monitoring (CIS, SOC2, HIPAA)
  - Multi-account security management
  - Cloud security baseline establishment
- **Example:** Monthly cloud security health checks

#### **Scout-Suite** - Multi-Cloud Security Auditing
**When to use:** In-depth cloud configuration security analysis
- **Use cases:**
  - Cloud configuration audit and reporting
  - Security misconfiguration detection
  - Cloud governance enforcement
  - Compliance preparation
- **Complements:** Prowler (different approach to cloud security assessment)

#### **Steampipe** - Cloud Asset Querying with SQL
**When to use:** SQL-based cloud asset inventory and analysis
- **Use cases:**
  - Custom cloud asset queries
  - Cloud inventory management
  - Compliance reporting with SQL
  - Cloud resource optimization analysis
- **Example:** Find all unencrypted S3 buckets across accounts

## Web Application Security Tools

#### **OWASP ZAP** - Web Application Security Scanner
**When to use:** Comprehensive web application penetration testing
- **Use cases:**
  - OWASP Top 10 vulnerability testing
  - Web application security testing in CI/CD
  - Manual penetration testing
  - Security regression testing
- **Example:** Automated security testing after deployments

#### **Nikto** - Web Server Security Scanner
**When to use:** Web server and CGI security scanning
- **Use cases:**
  - Web server misconfiguration detection
  - Outdated software detection
  - CGI and script vulnerability scanning
- **Complements:** ZAP (Nikto focuses on server, ZAP on applications)

#### **Nmap** - Network Discovery and Security Auditing
**When to use:** Network reconnaissance and port scanning
- **Use cases:**
  - Network discovery and mapping
  - Port scanning and service detection
  - Network security assessment
  - Infrastructure reconnaissance
- **Example:** Validate firewall rules and network segmentation

## Supply Chain Security Tools

#### **Syft** - SBOM Generation
**When to use:** Generate Software Bill of Materials for containers and applications
- **Use cases:**
  - Container SBOM generation
  - Software inventory management
  - Supply chain visibility
  - License compliance tracking
- **Example:** Generate SBOMs for container images in CI/CD

#### **Cosign** - Container Signing and Verification
**When to use:** Cryptographically sign and verify container images
- **Use cases:**
  - Container image signing in CI/CD
  - Supply chain integrity validation
  - Image provenance verification
  - Trust policy enforcement
- **Example:** Sign images before pushing to registry

#### **OSSF Scorecard** - Open Source Project Security Assessment
**When to use:** Evaluate security posture of open source dependencies
- **Use cases:**
  - Assess security of OSS dependencies
  - Supply chain risk management
  - Dependency selection criteria
  - Security due diligence
- **Example:** Evaluate OSS libraries before adoption

## Development and CI/CD Tools

#### **Semgrep** - Static Analysis for Security
**When to use:** Code-level security vulnerability detection
- **Use cases:**
  - SAST (Static Application Security Testing)
  - Custom security rule enforcement
  - Code quality and security in CI/CD
  - Developer security feedback
- **Example:** Detect SQL injection patterns in code

#### **ActionLint** - GitHub Actions Workflow Linter
**When to use:** Validate and optimize GitHub Actions workflows
- **Use cases:**
  - GitHub Actions workflow validation
  - CI/CD pipeline security
  - Workflow optimization
  - Best practices enforcement
- **Example:** Validate GitHub Actions in pull requests

#### **Hadolint** - Dockerfile Linter
**When to use:** Dockerfile best practices and security validation
- **Use cases:**
  - Dockerfile security best practices
  - Container image optimization
  - Docker security in CI/CD
  - Container build validation
- **Example:** Validate Dockerfiles before image builds

## When to Use Multiple Tools Together

### **Secret Scanning Strategy**
- **TruffleHog**: Deep analysis with verification
- **Gitleaks**: Fast CI/CD integration
- **git-secrets**: AWS-specific prevention
- **Use together**: TruffleHog for thorough audits, Gitleaks for CI speed, git-secrets for AWS prevention

### **Terraform Security Pipeline**
- **TFLint**: Syntax and best practices (first)
- **Checkov/TFSec**: Security scanning (second)  
- **terraform-docs**: Documentation (always)
- **OpenInfraQuote**: Cost estimation (for significant changes)
- **Use together**: TFLint → Security scanners → Cost estimation → Documentation

### **Container Security Workflow**
- **Hadolint**: Dockerfile validation (build time)
- **Trivy/Grype**: Image vulnerability scanning (build time)
- **Syft**: SBOM generation (build time)
- **Cosign**: Image signing (deployment time)
- **Use together**: Build-time validation → Runtime verification

### **Kubernetes Security Stack**
- **Kubescape**: Cluster security assessment
- **Kube-bench**: CIS compliance validation
- **Falco**: Runtime monitoring
- **Velero**: Backup and disaster recovery
- **Use together**: Assessment → Compliance → Monitoring → Backup

## Choosing Between Similar Tools

### **Trivy vs Grype**
- **Trivy**: Broader scope (containers, filesystems, git), faster, simpler
- **Grype**: SBOM-focused, detailed analysis, better for supply chain tracking
- **Choose Trivy if**: You want one tool for most vulnerability scanning
- **Choose Grype if**: You need detailed SBOM analysis and supply chain focus

### **Checkov vs TFSec**
- **Checkov**: Multi-framework (Terraform, CloudFormation, K8s), broader scope
- **TFSec**: Terraform-specific, deeper Terraform knowledge
- **Choose Checkov if**: You use multiple IaC frameworks
- **Choose TFSec if**: You're Terraform-focused and want the deepest analysis

### **TruffleHog vs Gitleaks**
- **TruffleHog**: Verification capability, more accurate, enterprise features
- **Gitleaks**: Faster, simpler, better for CI/CD speed
- **Choose TruffleHog if**: Accuracy and verification are critical
- **Choose Gitleaks if**: Speed and CI/CD integration are priorities

### **ZAP vs Nuclei**
- **ZAP**: Full-featured web app testing suite, GUI available, manual testing
- **Nuclei**: Template-based, automation-focused, community templates
- **Choose ZAP if**: You need comprehensive web app penetration testing
- **Choose Nuclei if**: You want fast, automated template-based scanning

## Integration Patterns

### **CI/CD Pipeline Integration**
```
Code Commit → TFLint → Checkov/TFSec → Gitleaks → Trivy → Deploy
           ↓
    terraform-docs → OpenInfraQuote (for cost-sensitive changes)
```

### **Security Review Workflow**  
```
Weekly: Prowler + Scout-Suite (cloud security)
Daily: Trivy + Grype (new container images)
Per PR: TFLint + Checkov + Gitleaks
Monthly: TruffleHog + OSSF Scorecard (deep audit)
```

### **Kubernetes Operations**
```
New Cluster: kube-bench + Kubescape (security baseline)
Ongoing: Falco (runtime monitoring) + Velero (backup)
Monthly: Goldilocks (resource optimization)
```

This comprehensive toolkit allows teams to implement defense-in-depth across their entire DevOps lifecycle, from code development through production operations.