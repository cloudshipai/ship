# Ship CLI MCP Tools Reference

This document provides a comprehensive reference of all available tools accessible through the Ship CLI's MCP server.

## Usage

Start the MCP server with all tools:
```bash
ship mcp
```

Start the MCP server with specific tool category:
```bash
ship mcp security    # Security tools only
ship mcp terraform   # Terraform tools only
ship mcp kubernetes  # Kubernetes tools only
```

Start the MCP server with specific tool:
```bash
ship mcp gitleaks     # Only Gitleaks tools
ship mcp trivy        # Only Trivy tools
```

## Tool Categories and Available Discrete Functions

### Security Tools (34 tools)

#### **Gitleaks** - Secret detection in code and git history
- `gitleaks_scan_directory` - Scan directory for secrets
- `gitleaks_scan_file` - Scan specific file for secrets  
- `gitleaks_scan_git_repo` - Scan git repository for secrets
- `gitleaks_scan_with_config` - Scan using custom configuration

#### **Trivy** - Comprehensive vulnerability scanner
- `trivy_scan_image` - Scan container image for vulnerabilities
- `trivy_scan_filesystem` - Scan filesystem for vulnerabilities
- `trivy_scan_repository` - Scan git repository for vulnerabilities
- `trivy_scan_config` - Scan configuration files for security issues
- `trivy_scan_sbom` - Scan SBOM file for vulnerabilities
- `trivy_get_version` - Get Trivy version information

#### **TruffleHog** - Advanced secret scanning with verification
- `trufflehog_scan_directory` - Scan directory for secrets
- `trufflehog_scan_git_repo` - Scan git repository for secrets
- `trufflehog_scan_github` - Scan GitHub repository for secrets
- `trufflehog_scan_github_org` - Scan GitHub organization for secrets
- `trufflehog_scan_docker` - Scan Docker image for secrets
- `trufflehog_scan_s3` - Scan S3 bucket for secrets
- `trufflehog_scan_verified` - Scan with secret verification

#### **Checkov** - Infrastructure as code static analysis
- `checkov_scan_directory` - Scan directory for security issues
- `checkov_scan_file` - Scan specific file for security issues
- `checkov_scan_with_policy` - Scan with custom policy
- `checkov_scan_multi_framework` - Scan with multiple frameworks
- `checkov_scan_with_severity` - Scan with severity threshold
- `checkov_scan_with_skips` - Scan with skipped checks
- `checkov_get_version` - Get Checkov version information

#### **Terrascan** - IaC security scanner
- `terrascan_scan_directory` - Scan directory for IaC security issues
- `terrascan_scan_terraform` - Scan Terraform files specifically
- `terrascan_scan_kubernetes` - Scan Kubernetes manifests
- `terrascan_scan_cloudformation` - Scan CloudFormation templates
- `terrascan_scan_dockerfiles` - Scan Dockerfiles for security issues
- `terrascan_scan_with_severity` - Scan with specific severity threshold

#### **Semgrep** - Static analysis for security
- `semgrep_scan_directory` - Scan directory with Semgrep rules
- `semgrep_scan_file` - Scan specific file with Semgrep
- `semgrep_scan_with_config` - Scan using custom configuration
- `semgrep_scan_with_rules` - Scan with specific rule sets
- `semgrep_get_version` - Get Semgrep version information

#### **Kubescape** - Kubernetes security scanner
- `kubescape_scan_cluster` - Scan Kubernetes cluster
- `kubescape_scan_manifests` - Scan Kubernetes manifests
- `kubescape_scan_helm` - Scan Helm chart
- `kubescape_scan_repository` - Scan repository for Kubernetes security issues
- `kubescape_generate_report` - Generate security report

#### **Nikto** - Web server security scanner
- `nikto_scan_host` - Scan web host for vulnerabilities
- `nikto_scan_ssl` - Scan host with SSL/TLS analysis
- `nikto_scan_tuned` - Scan with specific tuning options
- `nikto_get_version` - Get Nikto version

#### **Actionlint** - GitHub Actions workflow linter
- `actionlint_scan_directory` - Scan directory for workflow issues
- `actionlint_scan_file` - Scan specific workflow file
- `actionlint_get_version` - Get Actionlint version

#### **Hadolint** - Dockerfile linter
- `hadolint_scan_dockerfile` - Scan specific Dockerfile
- `hadolint_scan_directory` - Scan all Dockerfiles in directory
- `hadolint_get_version` - Get Hadolint version

#### **Conftest** - Policy testing with OPA
- `conftest_test_with_policy` - Test files against OPA policies
- `conftest_test_file` - Test specific file against policies
- `conftest_get_version` - Get Conftest version

#### **Git-secrets** - Git repository secret scanner
- `git_secrets_scan_repository` - Scan git repository for secrets
- `git_secrets_scan_aws` - Scan with AWS secret patterns
- `git_secrets_get_version` - Get git-secrets version

#### **ZAP** - OWASP ZAP web application scanner
- `zap_baseline_scan` - Run baseline security scan
- `zap_full_scan` - Run comprehensive security scan
- `zap_api_scan` - Scan API with OpenAPI/Swagger spec
- `zap_scan_with_context` - Scan with ZAP context file
- `zap_get_version` - Get ZAP version

#### **Falco** - Runtime security monitoring
- `falco_run_default` - Run Falco with default rules
- `falco_run_custom` - Run Falco with custom rules
- `falco_validate_rules` - Validate Falco rules
- `falco_get_version` - Get Falco version

#### **Prowler** - Multi-cloud security assessment
- `prowler_scan_aws` - Scan AWS account for security issues
- `prowler_scan_azure` - Scan Azure subscription for security issues
- `prowler_scan_gcp` - Scan GCP project for security issues
- `prowler_scan_kubernetes` - Scan Kubernetes cluster for security issues
- `prowler_scan_compliance` - Scan with specific compliance framework
- `prowler_scan_services` - Scan specific cloud services

#### **Dockle** - Container image linter
- `dockle_scan_image` - Scan container image for security and best practices
- `dockle_scan_tarball` - Scan container image tarball
- `dockle_scan_dockerfile` - Scan Dockerfile for best practices
- `dockle_generate_config` - Generate Dockle configuration template
- `dockle_list_checks` - List all available security checks
- `dockle_scan_with_policy` - Scan with custom policy

#### **SOPS** - Secrets management
- `sops_encrypt_file` - Encrypt file using SOPS
- `sops_decrypt_file` - Decrypt file using SOPS
- `sops_rotate_keys` - Rotate encryption keys for SOPS file
- `sops_edit_file` - Edit encrypted file using SOPS
- `sops_generate_config` - Generate SOPS configuration file
- `sops_validate_file` - Validate SOPS encrypted file integrity

### Terraform Tools (5 tools)

#### **TFLint** - Terraform linter
- `tflint_lint_directory` - Lint Terraform files in directory
- `tflint_lint_file` - Lint specific Terraform file
- `tflint_lint_with_config` - Lint using custom configuration
- `tflint_lint_with_rules` - Lint with specific rule sets
- `tflint_init_plugins` - Initialize TFLint plugins
- `tflint_get_version` - Get TFLint version information

#### **Terraform Docs** - Terraform documentation generator
- `terraform_docs_generate_markdown` - Generate documentation in Markdown format
- `terraform_docs_generate_json` - Generate documentation in JSON format
- `terraform_docs_generate_config` - Generate documentation using configuration file
- `terraform_docs_generate_table` - Generate documentation in table format
- `terraform_docs_get_version` - Get terraform-docs version information

#### **Infracost** - Infrastructure cost estimation
- `infracost_breakdown_directory` - Generate cost breakdown for Terraform directory
- `infracost_breakdown_plan` - Generate cost breakdown from Terraform plan file
- `infracost_diff` - Show cost difference between current and planned state
- `infracost_breakdown_config` - Generate cost breakdown using config file
- `infracost_generate_html` - Generate HTML cost report
- `infracost_generate_table` - Generate table format cost report
- `infracost_get_version` - Get Infracost version
- `infracost_get_pricing` - Get cloud pricing information

### Cloud & Infrastructure Tools (4 tools)

#### **CloudQuery** - Cloud asset inventory
- `cloudquery_sync_with_config` - Sync cloud resources using config
- `cloudquery_validate_config` - Validate CloudQuery configuration file
- `cloudquery_list_providers` - List available CloudQuery providers
- `cloudquery_get_version` - Get CloudQuery version

#### **Custodian** - Cloud governance engine
- `custodian_run_policy` - Run Cloud Custodian policy
- `custodian_validate_policy` - Validate Cloud Custodian policy syntax
- `custodian_dry_run` - Dry run Cloud Custodian policy

### Kubernetes Tools (7 tools)

#### **Kyverno** - Kubernetes policy management
- `kyverno_apply_policies` - Apply Kyverno policies to cluster
- `kyverno_validate_policies` - Validate Kyverno policy syntax
- `kyverno_test_policies` - Test Kyverno policies against resources
- `kyverno_get_version` - Get Kyverno version

### Supply Chain Security Tools (8 tools)

#### **GUAC** - Graph for Understanding Artifact Composition
- `guac_ingest_sbom` - Ingest SBOM into GUAC knowledge graph
- `guac_analyze_artifact` - Analyze artifact using GUAC
- `guac_query_dependencies` - Query package dependencies
- `guac_query_vulnerabilities` - Query vulnerabilities for package
- `guac_generate_graph` - Generate dependency graph
- `guac_analyze_impact` - Analyze vulnerability impact
- `guac_collect_files` - Collect and analyze project files
- `guac_validate_attestation` - Validate attestation

#### **Cosign Golden** - Advanced container signing workflows
- `cosign_golden_sign_keyless` - Sign container image using keyless signing
- `cosign_golden_verify_keyless` - Verify container image signature using keyless verification
- `cosign_golden_sign_pipeline` - Sign container image using golden pipeline workflow
- `cosign_golden_generate_attestation` - Generate attestation for container image
- `cosign_golden_verify_attestation` - Verify attestation for container image
- `cosign_golden_copy_signatures` - Copy signatures from source to destination image
- `cosign_golden_tree_view` - Display tree view of image signatures and attestations
- `cosign_golden_get_version` - Get Cosign Golden version information

#### **Sigstore Policy Controller** - Sigstore policy enforcement
- `sigstore_validate_policy` - Validate Sigstore policy syntax
- `sigstore_test_policy` - Test Sigstore policy against an image
- `sigstore_verify_signature` - Verify container image signature
- `sigstore_generate_template` - Generate Sigstore policy template
- `sigstore_validate_manifest` - Validate Kubernetes manifest against policy
- `sigstore_check_compliance` - Check manifests compliance with policy
- `sigstore_list_policies` - List available Sigstore policies
- `sigstore_audit_images` - Audit container images in namespace

### AWS Tools (6 tools)

#### **Parliament** - AWS IAM policy linter
- `parliament_lint_file` - Lint AWS IAM policy file
- `parliament_lint_directory` - Lint AWS IAM policy files in directory
- `parliament_lint_string` - Lint AWS IAM policy JSON string
- `parliament_lint_community` - Lint with community auditors
- `parliament_lint_private` - Lint with private auditors
- `parliament_lint_severity` - Lint with severity filter

#### **PMapper** - AWS IAM privilege escalation analysis
- `pmapper_create_graph` - Create IAM privilege graph
- `pmapper_query_access` - Query IAM access permissions
- `pmapper_find_privilege_escalation` - Find privilege escalation paths
- `pmapper_visualize_graph` - Visualize IAM privilege graph
- `pmapper_list_principals` - List IAM principals
- `pmapper_check_admin_access` - Check if principal has admin access

#### **Policy Sentry** - AWS IAM policy generator
- `policy_sentry_create_template` - Create IAM policy template
- `policy_sentry_write_policy` - Write IAM policy from input file
- `policy_sentry_write_from_template` - Write IAM policy from template YAML
- `policy_sentry_write_with_actions` - Write IAM policy with specific actions
- `policy_sentry_write_with_crud` - Write IAM policy with CRUD access levels
- `policy_sentry_query_actions` - Query AWS service action table
- `policy_sentry_query_conditions` - Query AWS service condition table

#### **Prowler** - Multi-cloud security assessment
- (Listed above in Security Tools section)

## Total Tool Count

- **Total Tools**: 64 tool categories
- **Total Discrete Functions**: 300+ individual MCP tools
- **Categories**: 6 (Security, Terraform, Cloud, Kubernetes, Supply Chain, AWS)

## Architecture

Each tool is implemented as:
1. **Dagger Module**: Containerized execution in `/internal/dagger/modules/`
2. **MCP Tool Registration**: Individual functions in `/internal/cli/mcp/`
3. **Registry Entry**: Categorized in the tool registry
4. **CLI Integration**: Accessible via `ship mcp [tool-name]`

## Adding New Tools

To add a new tool:
1. Create Dagger module in `/internal/dagger/modules/newtool.go`
2. Create MCP tools file in `/internal/cli/mcp/newtool.go`
3. Add entry to registry in `/internal/cli/mcp/registry.go`
4. Tool automatically becomes available via MCP server

All tools maintain the same containerized execution model while exposing their discrete functionality as individual MCP tools for AI assistants and automation.