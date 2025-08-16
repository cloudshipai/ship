# MCP Tool Conversion Status

## Overview
This document tracks the conversion of Ship CLI's 77 Dagger modules from single parameterized tools to multi-tool MCP patterns with discrete functions.

## Conversion Progress

### Completed Conversions (32 modules)
Successfully converted the following Dagger modules to modular MCP tool implementations:

| Module | MCP File | Discrete Tools | Category |
|--------|----------|----------------|----------|
| checkov | checkov.go | 7 | security |
| cloudsplaining | cloudsplaining.go | 4 | aws |
| cosign_golden | cosign_golden.go | 8 | supply-chain |
| dockle | dockle.go | 6 | security |
| gitleaks | gitleaks.go | 4 | security |
| grype | grype.go | 4 | security |
| guac | guac.go | 8 | supply-chain |
| iac_plan | iac_plan.go | 7 | terraform |
| in_toto | in_toto.go | 4 | supply-chain |
| infracost | infracost.go | 8 | terraform |
| inframap | inframap.go | 4 | terraform |
| kubescape | kubescape.go | 5 | security |
| ossf_scorecard | ossf_scorecard.go | 4 | security |
| parliament | parliament.go | 6 | aws |
| pmapper | pmapper.go | 6 | aws |
| policy_sentry | policy_sentry.go | 7 | aws |
| prowler | prowler.go | 6 | aws |
| rekor | rekor.go | 4 | supply-chain |
| sigstore_policy_controller | sigstore_policy_controller.go | 8 | supply-chain |
| sops | sops.go | 6 | security |
| steampipe | steampipe.go | 4 | security |
| syft | syft.go | 3 | security |
| terraform_docs | terraform_docs.go | 5 | terraform |
| trivy | trivy.go | 6 | security |
| trufflehog | trufflehog.go | 7 | security |
| velero | velero.go | 4 | kubernetes |
| tflint | tflint.go | 7 | terraform |
| terrascan | terrascan.go | 7 | security |
| semgrep | semgrep.go | 5 | security |
| cosign | cosign.go | 8 | supply-chain |

**Total Discrete MCP Tools: 172**

### Tool Distribution by Category
- **Security**: 14 modules, 95 discrete tools
- **AWS**: 5 modules, 29 discrete tools  
- **Supply Chain**: 6 modules, 40 discrete tools
- **Terraform**: 5 modules, 31 discrete tools
- **Kubernetes**: 1 module, 4 discrete tools

### Remaining Modules (45 modules)
The following 45 Dagger modules still need conversion to multi-tool MCP patterns:

| Module | Status | Notes |
|--------|--------|-------|
| actionlint | Pending | GitHub Actions linter |
| allstar | Pending | GitHub security policy |
| aws_iam_rotation | Pending | IAM key rotation |
| aws_pricing | Pending | AWS cost calculator |
| cert_manager | Pending | K8s certificate management |
| cfn_nag | Pending | CloudFormation linter |
| check_ssl_cert | Pending | SSL certificate checker |
| cloudquery | Pending | Cloud asset inventory |
| conftest | Pending | OPA policy testing |
| cosign | Pending | Container signing |
| custodian | Pending | Cloud governance |
| dependency_track | Pending | Component analysis |
| falco | Pending | Runtime security |
| fleet | Pending | GitOps management |
| gatekeeper | Pending | K8s policy controller |
| git_secrets | Pending | Git secret scanner |
| github_admin | Pending | GitHub administration |
| github_packages | Pending | GitHub packages |
| goldilocks | Pending | K8s resource recommendations |
| hadolint | Pending | Dockerfile linter |
| history_scrub | Pending | Git history cleaner |
| infrascan | Pending | Infrastructure scanner |
| k8s_network_policy | Pending | K8s network policies |
| kube_bench | Pending | CIS benchmark |
| kube_hunter | Pending | K8s penetration testing |
| kuttl | Pending | K8s test framework |
| kyverno | Pending | K8s policy management |
| kyverno_multitenant | Pending | Multi-tenant policies |
| license_detector | Pending | License scanner |
| litmus | Pending | Chaos engineering |
| modules_test | Pending | Test utilities |
| nikto | Pending | Web scanner |
| openinfraquote | Pending | Cost analysis |
| openscap | Pending | Security compliance |
| osv_scanner | Pending | OSV vulnerability scanner |
| packer | Pending | Image builder |
| powerpipe | Pending | Security dashboards |
| query_fixer | Pending | Query optimization |
| query_templates | Pending | Query templates |
| registry | Pending | Container registry |
| scout_suite | Pending | Multi-cloud security |
| semgrep | Pending | Static analysis |
| slsa_verifier | Pending | SLSA verification |
| step_ca | Pending | Certificate authority |
| terraformer | Pending | Infrastructure import |
| terrascan | Pending | IaC security scanner |
| tflint | Pending | Terraform linter |
| tfstate_reader | Pending | State file reader |
| tool_services | Pending | Tool service utilities |
| trivy_golden | Pending | Enhanced Trivy workflows |
| zap | Pending | OWASP ZAP scanner |

## Architecture Benefits

### Before Conversion
- Tools implemented as single parameterized functions
- Limited discoverability through MCP protocol
- Naming conflicts with periods (.) causing errors
- Monolithic mcp_cmd.go file (4000+ lines)

### After Conversion
- Each tool exposes multiple discrete functions
- Improved tool discoverability for AI assistants
- Clean naming with underscores (_) instead of periods
- Modular architecture with separate files
- Tool registry for organized management

## Next Steps

1. **Complete Remaining Conversions**: Convert all 49 remaining Dagger modules
2. **Test Discovery**: Verify all tools are discoverable via `ship mcp`
3. **Documentation**: Update user documentation with new tool patterns
4. **Integration Testing**: Test MCP client integration with all discrete tools

## Usage Examples

### Single Tool Server
```bash
# Start MCP server for specific tool
ship mcp gitleaks
```

### Tool Discovery
```bash
# Discover all available tools
ship mcp --list-tools
```

### Individual Function Usage
```bash
# Access discrete tool functions via MCP client
gitleaks_detect_secrets
gitleaks_scan_repository  
gitleaks_scan_directory
gitleaks_generate_baseline
```