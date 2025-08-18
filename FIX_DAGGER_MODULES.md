# Dagger Module Container Execution Fix Guide

## Problem Statement
Many Dagger modules in Ship are failing because they incorrectly handle Docker container entrypoints. When `WithExec()` is called in Dagger, it overrides the container's entrypoint, so we need to use the full path to executables.

## Systematic Fix Process

### Step 1: Identify the Container's Executable Path
For each tool's Docker image, determine where the executable is located:

```bash
# Check the entrypoint
docker inspect <image:tag> | grep -A5 'Entrypoint'

# Or run a shell to explore
docker run --rm --entrypoint sh <image:tag> -c "which <tool>"
```

### Step 2: Update the Dagger Module
Add a constant for the binary path and update all `WithExec()` calls:

```go
// Add constant after the struct definition
const toolBinary = "/path/to/tool"  // e.g., "/usr/local/bin/tflint"

// Update all WithExec calls to use the full path
// BEFORE:
WithExec([]string{"tool", "command", "args"})

// AFTER:
WithExec([]string{toolBinary, "command", "args"})
```

### Step 3: Test the Fix
Create a test file to verify the module works:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/cloudshipai/ship/internal/dagger/modules"
    "dagger.io/dagger"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
    if err != nil {
        log.Fatalf("Failed to connect to Dagger: %v", err)
    }
    defer client.Close()
    
    // Test the module
    module := modules.NewToolModule(client)
    result, err := module.GetVersion(ctx)
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    fmt.Printf("✅ Success: %s\n", result)
}
```

## Tools to Fix (Alphabetically)

### ✅ COMPLETED FIXES (28/32)

1. **Custodian**
   - File: `internal/dagger/modules/custodian.go`
   - Image: `cloudcustodian/c7n:latest`
   - Binary Path: `/src/.venv/bin/custodian`
   - Status: FIXED ✅

2. **Actionlint** 
   - File: `internal/dagger/modules/actionlint.go`
   - Image: `rhysd/actionlint:latest`
   - Binary Path: `/usr/local/bin/actionlint`
   - Status: FIXED ✅

3. **AWS IAM Rotation**
   - File: `internal/dagger/modules/aws_iam_rotation.go`
   - Image: `amazon/aws-cli:latest`
   - Binary Path: `/usr/local/bin/aws`
   - Status: FIXED ✅

4. **AWS Pricing**
   - File: `internal/dagger/modules/aws_pricing.go`
   - Image: `amazon/aws-cli:latest`
   - Binary Path: `/usr/local/bin/aws`
   - Status: FIXED ✅

5. **CFN Nag**
   - File: `internal/dagger/modules/cfn_nag.go`
   - Image: `stelligent/cfn_nag:latest`
   - Binary Path: `/usr/local/bundle/bin/cfn_nag_scan`
   - Status: FIXED ✅

6. **Checkov**
   - File: `internal/dagger/modules/checkov.go`
   - Image: `bridgecrew/checkov:latest`
   - Binary Path: `/usr/local/bin/checkov`
   - Status: FIXED ✅

7. **CloudQuery**
   - File: `internal/dagger/modules/cloudquery.go`
   - Image: `ghcr.io/cloudquery/cloudquery:latest`
   - Binary Path: `/app/cloudquery`
   - Status: FIXED ✅

8. **Conftest**
   - File: `internal/dagger/modules/conftest.go`
   - Image: `openpolicyagent/conftest:latest`
   - Binary Path: `/conftest`
   - Status: FIXED ✅

9. **Cosign**
   - File: `internal/dagger/modules/cosign.go`
   - Image: `gcr.io/projectsigstore/cosign:latest`
   - Binary Path: `/ko-app/cosign`
   - Status: FIXED ✅

10. **Cosign Golden**
    - File: `internal/dagger/modules/cosign_golden.go`
    - Image: `gcr.io/projectsigstore/cosign:latest`
    - Binary Path: `/ko-app/cosign`
    - Status: FIXED ✅

11. **Allstar**
    - File: `internal/dagger/modules/allstar.go`
    - Status: SKIPPED (GitHub App, not a CLI tool) ⏭️

12. **Dockle**
    - File: `internal/dagger/modules/dockle.go`
    - Image: `goodwithtech/dockle:*`
    - Binary Path: `/usr/bin/dockle`
    - Status: FIXED ✅

13. **Gitleaks**
    - File: `internal/dagger/modules/gitleaks.go`
    - Image: `zricethezav/gitleaks:latest`
    - Binary Path: `/usr/bin/gitleaks`
    - Status: FIXED ✅

14. **Grype**
    - File: `internal/dagger/modules/grype.go`
    - Image: `anchore/grype:latest`
    - Binary Path: `/grype`
    - Status: FIXED ✅

15. **Hadolint**
    - File: `internal/dagger/modules/hadolint.go`
    - Image: `hadolint/hadolint:latest`
    - Binary Path: `/bin/hadolint`
    - Status: FIXED ✅

16. **Infracost**
    - File: `internal/dagger/modules/infracost.go`
    - Image: `infracost/infracost:latest`
    - Binary Path: `/usr/bin/infracost`
    - Status: FIXED ✅

17. **InfraMap**
    - File: `internal/dagger/modules/inframap.go`
    - Image: Custom (downloads binary)
    - Binary Path: `/usr/local/bin/inframap`
    - Status: FIXED ✅

18. **Falco**
    - File: `internal/dagger/modules/falco.go`
    - Image: `falcosecurity/falco:latest`
    - Binary Path: `/usr/bin/falco`
    - Status: FIXED ✅

19. **Fleet**
    - File: `internal/dagger/modules/fleet.go`
    - Image: `rancher/fleet:latest`
    - Binary Path: `/usr/local/bin/kubectl`
    - Status: FIXED ✅

20. **Git Secrets**
    - File: `internal/dagger/modules/git_secrets.go`
    - Image: `cloudshipai/git-secrets:latest`
    - Binary Path: `/usr/bin/git`
    - Status: FIXED ✅

21. **Goldilocks**
    - File: `internal/dagger/modules/goldilocks.go`
    - Image: `us-docker.pkg.dev/fairwinds-ops/oss/goldilocks:latest`
    - Binary Path: `/goldilocks`
    - Status: FIXED ✅

22. **Cert Manager**
    - File: `internal/dagger/modules/cert_manager.go`
    - Image: `quay.io/jetstack/cert-manager-ctl:latest`
    - Binary Path: `/usr/local/bin/cmctl`
    - Status: FIXED ✅

23. **Dependency Track**
    - File: `internal/dagger/modules/dependency_track.go`
    - Image: `node:alpine` (installs tools)
    - Binary Paths: Multiple (`/usr/local/bin/dtrack-cli`, `/usr/local/bin/syft`, etc.)
    - Status: FIXED ✅

24. **Gatekeeper**
    - File: `internal/dagger/modules/gatekeeper.go`
    - Images: Various (OPA, Helm, etc.)
    - Binary Paths: `/usr/local/bin/kubectl`, `/usr/local/bin/helm`, `/manager`
    - Status: FIXED ✅

25. **GitHub Admin**
    - File: `internal/dagger/modules/github_admin.go`
    - Image: `alpine:latest`
    - Binary Paths: Standard system tools (curl, jq)
    - Status: NO FIX NEEDED ✅

26. **GUAC**
    - File: `internal/dagger/modules/guac.go`
    - Image: `ghcr.io/guacsec/guac:latest`
    - Binary Path: `/usr/local/bin/guacone`
    - Status: FIXED ✅

27. **IAC Plan**
    - File: `internal/dagger/modules/iac_plan.go`
    - Images: `hashicorp/terraform:latest`, `opentofu/opentofu:latest`
    - Binary Paths: `/bin/terraform`, `/usr/local/bin/tofu`
    - Status: FIXED ✅

### 🔧 Remaining Tools to Fix

1. **actionlint**
   - File: `internal/dagger/modules/actionlint.go`
   - Image: `rhysd/actionlint:latest`
   - Binary Path: TBD

2. **allstar**
   - File: `internal/dagger/modules/allstar.go`
   - Image: `ossf/allstar:latest`
   - Binary Path: TBD

3. **aws_iam_rotation**
   - File: `internal/dagger/modules/aws_iam_rotation.go`
   - Image: TBD
   - Binary Path: TBD

4. **aws_pricing**
   - File: `internal/dagger/modules/aws_pricing.go`
   - Image: TBD
   - Binary Path: TBD

5. **cert_manager**
   - File: `internal/dagger/modules/cert_manager.go`
   - Image: `bitnami/kubectl:latest` or custom
   - Binary Path: TBD

6. **cfn_nag**
   - File: `internal/dagger/modules/cfn_nag.go`
   - Image: `stelligent/cfn_nag:latest`
   - Binary Path: TBD

7. **checkov**
   - File: `internal/dagger/modules/checkov.go`
   - Image: `bridgecrew/checkov:latest`
   - Binary Path: TBD

8. **cloudquery**
   - File: `internal/dagger/modules/cloudquery.go`
   - Image: `cloudquery/cloudquery:latest`
   - Binary Path: TBD

9. **cloudsplaining**
   - File: `internal/dagger/modules/cloudsplaining.go`
   - Image: TBD
   - Binary Path: TBD

10. **conftest**
    - File: `internal/dagger/modules/conftest.go`
    - Image: `instrumenta/conftest:latest`
    - Binary Path: TBD

11. **cosign**
    - File: `internal/dagger/modules/cosign.go`
    - Image: `gcr.io/projectsigstore/cosign:latest`
    - Binary Path: TBD

12. **cosign_golden**
    - File: `internal/dagger/modules/cosign_golden.go`
    - Image: `gcr.io/projectsigstore/cosign:latest`
    - Binary Path: TBD

13. **dependency_track**
    - File: `internal/dagger/modules/dependency_track.go`
    - Image: `owasp/dependency-track:latest`
    - Binary Path: TBD

14. **dockle**
    - File: `internal/dagger/modules/dockle.go`
    - Image: `goodwithtech/dockle:latest`
    - Binary Path: TBD

15. **falco**
    - File: `internal/dagger/modules/falco.go`
    - Image: `falcosecurity/falco:latest`
    - Binary Path: TBD

16. **fleet**
    - File: `internal/dagger/modules/fleet.go`
    - Image: `rancher/fleet:latest`
    - Binary Path: TBD

17. **gatekeeper**
    - File: `internal/dagger/modules/gatekeeper.go`
    - Image: `openpolicyagent/gatekeeper:latest`
    - Binary Path: TBD

18. **git_secrets**
    - File: `internal/dagger/modules/git_secrets.go`
    - Image: Custom or base image with git-secrets
    - Binary Path: TBD

19. **github_admin**
    - File: `internal/dagger/modules/github_admin.go`
    - Image: Custom with gh CLI
    - Binary Path: TBD

20. **github_packages**
    - File: `internal/dagger/modules/github_packages.go`
    - Image: Custom with gh CLI
    - Binary Path: TBD

21. **gitleaks**
    - File: `internal/dagger/modules/gitleaks.go`
    - Image: `zricethezav/gitleaks:latest`
    - Binary Path: TBD

22. **goldilocks**
    - File: `internal/dagger/modules/goldilocks.go`
    - Image: `fairwinds/goldilocks:latest`
    - Binary Path: TBD

23. **grype**
    - File: `internal/dagger/modules/grype.go`
    - Image: `anchore/grype:latest`
    - Binary Path: TBD

24. **guac**
    - File: `internal/dagger/modules/guac.go`
    - Image: `ghcr.io/guacsec/guac:latest`
    - Binary Path: TBD

25. **hadolint**
    - File: `internal/dagger/modules/hadolint.go`
    - Image: `hadolint/hadolint:latest`
    - Binary Path: TBD

26. **iac_plan**
    - File: `internal/dagger/modules/iac_plan.go`
    - Image: Custom Terraform image
    - Binary Path: TBD

27. **in_toto**
    - File: `internal/dagger/modules/in_toto.go`
    - Image: Custom with in-toto
    - Binary Path: TBD

28. **infracost**
    - File: `internal/dagger/modules/infracost.go`
    - Image: `infracost/infracost:latest`
    - Binary Path: TBD

29. **inframap**
    - File: `internal/dagger/modules/inframap.go`
    - Image: `cycloid/inframap:latest`
    - Binary Path: TBD

30. **infrascan**
    - File: `internal/dagger/modules/infrascan.go`
    - Image: Custom
    - Binary Path: TBD

## Fix Template

```go
// Example fix for any module
const <tool>Binary = "/usr/local/bin/<tool>"  // Update with actual path

// In each function, replace:
// WithExec([]string{"<tool>", "arg1", "arg2"})
// With:
// WithExec([]string{<tool>Binary, "arg1", "arg2"})
```

## Testing Checklist
- [ ] Identify Docker image and version
- [ ] Find executable path in container
- [ ] Add binary constant to module
- [ ] Update all WithExec calls
- [ ] Test GetVersion() if available
- [ ] Test at least one main function
- [ ] Verify no local installation needed

## Notes
- Some tools may not have issues if they don't use entrypoints
- Some tools may need additional environment variables or working directory settings
- Complex tools may need more extensive refactoring beyond just the binary path

## Modules Confirmed to Need Fixing (Based on Docker Entrypoints)

These modules use Docker images with entrypoints that require full binary paths:

1. **cfn_nag** - `stelligent/cfn_nag:latest` - Entrypoint: `cfn_nag`
   - Need to find path for `cfn_nag_scan` command
   
2. **checkov** - `bridgecrew/checkov:latest` - Entrypoint: `/entrypoint.sh`
   - Complex - uses shell script entrypoint
   
3. **cloudquery** - `ghcr.io/cloudquery/cloudquery:latest` - Entrypoint: `/app/cloudquery`
   - Binary at `/app/cloudquery`

4. **aws_pricing** - `amazon/aws-cli:latest` - Same as aws_iam_rotation
   - Binary at `/usr/local/bin/aws`

5. **cert_manager** - `quay.io/jetstack/cert-manager-ctl:latest`
   - Need to check entrypoint

Additional modules to investigate for similar issues.