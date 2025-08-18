# MCP→Dagger Direct Conversion Status

This file tracks the conversion of all MCP tools from `executeShipCommand()` calls to direct Dagger module calls.

## Overview
**Goal:** Convert all MCP tools to bypass CLI layer and call Dagger modules directly.

**Pattern:**
```
BEFORE: MCP → executeShipCommand() → CLI (broken) → Dagger
AFTER:  MCP → dagger.Connect() → Dagger Module → Container
```

## Conversion Progress: 62/72 (86.1%)

### 📋 **ALPHABETICAL CONVERSION LIST**

| # | Tool Name | Status | MCP File | Dagger Module | Functions | Notes |
|---|-----------|--------|----------|---------------|-----------|-------|
| 1 | **actionlint** | ✅ DONE | actionlint.go | actionlint.go | 3→7 | GitHub Actions linter |
| 2 | **allstar** | ✅ DONE | allstar.go | allstar.go | 1→1 | GitHub App (info only) |
| 3 | **aws-iam-rotation** | ✅ DONE | aws_iam_rotation.go | aws_iam_rotation.go | 6 | IAM credential rotation |
| 4 | **aws-pricing** | ✅ DONE | aws_pricing.go | aws_pricing.go | 4→7 | AWS pricing calculator |
| 5 | **cert-manager** | ✅ DONE | cert_manager.go | cert_manager.go | 8→9 | Certificate management |
| 6 | **cfn-nag** | ✅ DONE | cfn_nag.go | cfn_nag.go | 6→9 | CloudFormation linter |
| 7 | **check-ssl-cert** | ✅ DONE | check_ssl_cert.go | check_ssl_cert.go | 6→9 | SSL cert validation |
| 8 | **checkov** | ✅ DONE | checkov.go | checkov.go | 10→16 | IaC security scanner |
| 9 | **cloudquery** | ✅ DONE | cloudquery.go | cloudquery.go | 10→14 | Cloud asset inventory |
| 10 | **cloudsplaining** | ✅ DONE | cloudsplaining.go | cloudsplaining.go | 6→9 | AWS IAM scanner |
| 11 | **conftest** | ✅ DONE | conftest.go | conftest.go | 5→8 | OPA policy testing |
| 12 | **container-registry** | ✅ DONE | container_registry.go | container_registry.go | 6→6 | Registry operations |
| 13 | **cosign** | ✅ DONE | cosign.go | cosign.go | 12→18 | Container signing |
| 14 | **cosign-advanced** | ❌ N/A | - | - | - | File doesn't exist |
| 15 | **cosign-golden** | ✅ DONE | cosign_golden.go | cosign_golden.go | 10→15 | Golden image cosign |
| 16 | **custodian** | ✅ DONE | custodian.go | custodian.go | 5 | Cloud governance |
| 17 | **dependency-track** | ✅ DONE | dependency_track.go | dependency_track.go | 3→8 | SBOM analysis |
| 18 | **dockle** | ✅ DONE | dockle.go | dockle.go | 4→10 | Container linter |
| 19 | **falco** | ✅ DONE | falco.go | falco.go | 7→15 | Runtime security |
| 20 | **fleet** | ✅ DONE | fleet.go | fleet.go | 6→13 | GitOps for K8s |
| 21 | **gatekeeper** | ✅ DONE | gatekeeper.go | gatekeeper.go | 7→14 | OPA Gatekeeper |
| 22 | **git-secrets** | ✅ DONE | git_secrets.go | git_secrets.go | 7→9 | Git secret scanner |
| 23 | **github-admin** | ✅ DONE | github_admin.go | github_admin.go | 6→15 | GitHub admin tools |
| 24 | **github-packages** | ✅ DONE | github_packages.go | github_packages.go | 6→14 | GitHub packages |
| 25 | **gitleaks** | ✅ DONE | gitleaks.go | gitleaks.go | 4→6 | Secret detection |
| 26 | **goldilocks** | ✅ DONE | goldilocks.go | goldilocks.go | 5→7 | K8s resource sizing |
| 27 | **grype** | ✅ DONE | grype.go | grype.go | 8→11 | Vulnerability scanner |
| 28 | **guac** | ✅ DONE | guac.go | guac.go | 6→9 | Supply chain analysis |
| 29 | **hadolint** | ✅ DONE | hadolint.go | hadolint.go | 5→6 | Dockerfile linter |
| 30 | **history-scrub** | ✅ DONE | history_scrub.go | history_scrub.go | 6→6 | Git history cleaning |
| 31 | **iac-plan** | ✅ DONE | iac_plan.go | iac_plan.go | 7→8 | IaC planning |
| 32 | **in-toto** | ✅ DONE | in_toto.go | in_toto.go | 4→12 | Supply chain attestation |
| 33 | **infracost** | ✅ DONE | infracost.go | infracost.go | 8→9 | Cost estimation |
| 34 | **inframap** | ✅ DONE | inframap.go | inframap.go | 2→4 | Infrastructure viz |
| 35 | **infrascan** | ✅ DONE | infrascan.go | infrascan.go | 3→7 | AWS infrastructure mapping + Trivy security |
| 36 | **k8s-network-policy** | ✅ DONE | k8s_network_policy.go | k8s_network_policy.go | 6→9 | K8s network policies + enhanced |
| 37 | **kube-bench** | ✅ DONE | kube_bench.go | kube_bench.go | 6→8 | K8s CIS benchmark |
| 38 | **kube-hunter** | ✅ DONE | kube_hunter.go | kube_hunter.go | 6→7 | K8s penetration test + enhanced |
| 39 | **kubescape** | ✅ DONE | kubescape.go | kubescape.go | 3→8 | K8s security scanner + fixed Dagger module |
| 40 | **kuttl** | ✅ DONE | kuttl.go | kuttl.go | 4→5 | K8s testing + enhanced |
| 41 | **kyverno** | ✅ DONE | kyverno.go | kyverno.go | 9→11 | K8s policy mgmt + enhanced |
| 42 | **kyverno-multitenant** | ✅ DONE | kyverno_multitenant.go | kyverno_multitenant.go | 6→7 | Multi-tenant policies + enhanced |
| 43 | **license-detector** | ✅ DONE | license_detector.go | license_detector.go | 9→10 | License detection + enhanced |
| 44 | **litmus** | ✅ DONE | litmus.go | litmus.go | 12→14 | Chaos engineering + enhanced |
| 45 | **nikto** | ✅ DONE | nikto.go | nikto.go | 11→11 | Web scanner + enhanced |
| 46 | **openinfraquote** | ✅ DONE | openinfraquote.go | openinfraquote.go | 5→9 | Cost analysis + enhanced |
| 47 | **openscap** | ✅ DONE | openscap.go | openscap.go | 10→11 | Compliance scanning + NIST certified |
| 48 | **ossf-scorecard** | ✅ DONE | ossf_scorecard.go | ossf_scorecard.go | 4→4 | OSSF scorecard + enhanced |
| 49 | **osv-scanner** | ✅ DONE | osv_scanner.go | osv_scanner.go | 10→11 | OSV scanning |
| 50 | **packer** | ✅ DONE | packer.go | packer.go | 10→10 | Image building |
| 51 | **parliament** | ✅ DONE | parliament.go | parliament.go | 9→10 | IAM policy linting |
| 52 | **pmapper** | ✅ DONE | pmapper.go | pmapper.go | 6→6 | IAM privilege mapping |
| 53 | **policy-sentry** | ✅ DONE | policy_sentry.go | policy_sentry.go | 8→8 | IAM policy generation |
| 54 | **powerpipe** | ✅ DONE | powerpipe.go | powerpipe.go | 8→9 | Infrastructure bench |
| 55 | **prowler** | ✅ DONE | prowler.go | prowler.go | 9→11 | Multi-cloud security |
| 56 | **rekor** | ✅ DONE | rekor.go | rekor.go | 7→7 | Transparency log |
| 57 | **scout-suite** | ✅ DONE | scout_suite.go | scout_suite.go | 5→6 | Cloud auditing |
| 58 | **semgrep** | ✅ DONE | semgrep.go | semgrep.go | 11→14 | Static analysis |
| 59 | **sigstore-policy-controller** | ✅ DONE | sigstore_policy_controller.go | sigstore_policy_controller.go | 9→16 | Sigstore policy |
| 60 | **slsa-verifier** | ✅ DONE | slsa_verifier.go | slsa_verifier.go | 4→4 | SLSA verification |
| 61 | **sops** | ✅ DONE | sops.go | sops.go | 6→6 | Secrets management |
| 62 | **steampipe** | ✅ DONE | steampipe.go | steampipe.go | 10→11 | Cloud queries |
| 63 | **step-ca** | ✅ DONE | step_ca.go | step_ca.go | 6→6 | CA operations |
| 64 | **syft** | ⬜ TODO | syft.go | syft.go | 12 | SBOM generation |
| 65 | **terraform-docs** | ⬜ TODO | terraform_docs.go | terraform_docs.go | 6 | TF documentation |
| 66 | **terraformer** | ⬜ TODO | terraformer.go | terraformer.go | 6 | Infrastructure import |
| 67 | **terrascan** | ⬜ TODO | terrascan.go | terrascan.go | 15 | IaC scanner |
| 68 | **tflint** | ⬜ TODO | tflint.go | tflint.go | 8 | Terraform linter |
| 69 | **tfstate-reader** | ⬜ TODO | tfstate_reader.go | tfstate_reader.go | 7 | State analysis |
| 70 | **trivy** | ⬜ TODO | trivy.go | trivy.go | 10+ | Vulnerability scanner |
| 71 | **trivy-golden** | ⬜ TODO | trivy_golden.go | trivy_golden.go | 6 | Golden image trivy |
| 72 | **trufflehog** | ⬜ TODO | trufflehog.go | trufflehog.go | 15 | Secret scanner |
| 73 | **velero** | ⬜ TODO | velero.go | velero.go | 15 | K8s backup/restore |
| 74 | **zap** | ⬜ TODO | zap.go | zap.go | 6 | Web app scanner |

## Conversion Checklist (Per Tool)

For each tool conversion:

### 🔍 **1. Analysis Phase**
- [ ] Check if Dagger module exists in `internal/dagger/modules/`
- [ ] Verify Dagger module functions match MCP tool functions
- [ ] Check official CLI documentation for correct syntax
- [ ] Identify any missing parameters or features

### 🔧 **2. Implementation Phase** 
- [ ] Update imports (add `dagger.io/dagger`, remove unused)
- [ ] Add wrapper function to maintain registry compatibility:
  ```go
  func AddXxxTools(s *server.MCPServer, executeShipCommand ExecuteShipCommandFunc) {
      // Ignore executeShipCommand - we use direct Dagger calls
      addXxxToolsDirect(s)
  }
  ```
- [ ] Replace all `executeShipCommand(args)` calls with:
  ```go
  client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
  if err != nil { return mcp.NewToolResultError(...) }
  defer client.Close()
  
  module := modules.NewXxxModule(client)
  result, err := module.SomeFunction(ctx, params...)
  ```
- [ ] Handle unsupported parameters gracefully (add warnings)

### ✅ **3. Verification Phase**
- [ ] Run `make build` - compiles successfully  
- [ ] Check for import or syntax errors
- [ ] Verify all MCP tool functions converted
- [ ] Update Dagger module if CLI syntax needs fixing

### 📋 **4. Documentation Phase**
- [ ] Update status in this file (⬜ TODO → ✅ DONE)
- [ ] Note any limitations or issues
- [ ] Commit with clear message

## Next Tool: **github-packages** (#24)

### ✅ Recently Completed: **github-admin** (#23) 
- **Status**: ✅ DONE  
- **GitHub**: N/A (GitHub CLI wrapper tool)
- **MCP Functions**: 6 (github_admin_list_org_repos, github_admin_create_org_repo, github_admin_get_repo_info, github_admin_list_org_issues, github_admin_list_org_prs, github_admin_get_version)
- **Dagger Functions**: 15 (enhanced with ListOrgReposSimple, CreateOrgRepoSimple, GetRepoInfoSimple, ListOrgIssuesSimple, ListOrgPRsSimple, GetVersionSimple + existing 9 functions)
- **Build**: ✅ Successful
- **Impact**: Complete GitHub organization administration using official GitHub CLI with repo management, issue tracking, and PR monitoring

---
**Last Updated**: 2025-08-18  
**Status**: 22/72 tools converted (30.6% complete) - 🎉 30% MILESTONE REACHED!