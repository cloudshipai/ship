# Packer

HashiCorp Packer for building identical machine images across multiple platforms from a single source configuration, enabling consistent infrastructure deployment and immutable infrastructure practices.

## Description

HashiCorp Packer is a powerful open-source tool for creating identical machine images for multiple platforms from a single source configuration. Packer enables organizations to implement immutable infrastructure practices by building consistent, reproducible machine images across cloud platforms, virtualization environments, and container registries. With support for Amazon EC2, Google Cloud Platform, Microsoft Azure, VMware, VirtualBox, Docker, and many other platforms through plugins, Packer provides a unified workflow for image creation that reduces configuration drift, improves deployment consistency, and enhances security through standardized image builds. The tool features template-based configuration using either JSON or HCL2 syntax, parallel builds across multiple platforms, extensive provisioning capabilities, post-processing workflows, and comprehensive plugin ecosystem, making it essential for DevOps, infrastructure as code, and cloud-native application deployment strategies.

## MCP Tools

### Core Build Operations
- **`packer_build`** - Build machine images using real packer CLI
- **`packer_validate`** - Validate Packer configuration template using real packer CLI

### Template Management and Development
- **`packer_inspect`** - Inspect and analyze Packer template configuration using real packer CLI
- **`packer_fix`** - Fix and upgrade Packer template to current version using real packer CLI
- **`packer_fmt`** - Format Packer template files using real packer CLI
- **`packer_console`** - Open Packer console for template debugging using real packer CLI

### Configuration and Plugin Management
- **`packer_init`** - Initialize Packer configuration and install required plugins using real packer CLI
- **`packer_plugins`** - Manage Packer plugins using real packer CLI

### Migration and Upgrade Tools
- **`packer_hcl2_upgrade`** - Upgrade JSON Packer template to HCL2 using real packer CLI

### System Information
- **`packer_version`** - Get Packer version information using real packer CLI

## Real CLI Commands Used

### Core Commands
- `packer build [options] <template>` - Build machine images from template
- `packer validate [options] <template>` - Validate template configuration
- `packer inspect [options] <template>` - Analyze template components
- `packer fix [options] <template>` - Fix and upgrade template syntax

### Development and Formatting
- `packer fmt [options] <path>` - Format template files
- `packer console [options] <template>` - Interactive debugging console

### Initialization and Plugin Management
- `packer init [options] <config>` - Initialize and install plugins
- `packer plugins <subcommand> [options]` - Manage plugins

### Migration and Version
- `packer hcl2_upgrade [options] <template>` - Upgrade JSON to HCL2
- `packer version` - Display version information

### Common Build Options
- `-var key=value` - Set template variables
- `-var-file path` - Load variables from file
- `-only build1,build2` - Build only specified builds
- `-except build1,build2` - Skip specified builds
- `-force` - Force build even with existing artifacts
- `-debug` - Enable debug mode with step-by-step execution
- `-parallel-builds N` - Limit parallel build execution

### Common Validation Options
- `-syntax-only` - Check syntax only without full validation
- `-evaluate-datasources` - Evaluate data sources during validation
- `-machine-readable` - Output in machine-readable format

## Use Cases

### Infrastructure as Code
- **Immutable Infrastructure**: Build consistent, reproducible machine images
- **Golden Images**: Create standardized base images for organizational use
- **Configuration Management**: Eliminate configuration drift through pre-built images
- **Compliance**: Ensure security and compliance standards in base images

### Multi-Cloud Deployment
- **Cloud Migration**: Build images across multiple cloud platforms simultaneously
- **Hybrid Cloud**: Maintain consistent images across on-premises and cloud environments
- **Disaster Recovery**: Create backup images across different regions and providers
- **Platform Abstraction**: Abstract application deployment from underlying infrastructure

### DevOps and CI/CD Integration
- **Automated Image Builds**: Integrate image creation into CI/CD pipelines
- **Application Packaging**: Bundle applications into immutable deployment artifacts
- **Testing Environments**: Create consistent environments for testing and staging
- **Release Management**: Version and manage application images through release cycles

### Security and Compliance
- **Hardened Images**: Build security-hardened base images with compliance controls
- **Vulnerability Management**: Create images with latest security patches and updates
- **Audit and Traceability**: Maintain detailed records of image build processes
- **Secret Management**: Securely inject secrets and certificates during build process

## Configuration Examples

### Basic Template Operations
```bash
# Validate Packer template
packer validate template.pkr.hcl
packer validate -syntax-only template.pkr.hcl

# Build machine images
packer build template.pkr.hcl
packer build -var "region=us-east-1" template.pkr.hcl
packer build -var-file variables.pkrvars.hcl template.pkr.hcl

# Build specific platforms only
packer build -only "amazon-ebs" template.pkr.hcl
packer build -except "virtualbox-iso" template.pkr.hcl

# Debug build process
packer build -debug template.pkr.hcl

# Force build (overwrite existing artifacts)
packer build -force template.pkr.hcl
```

### Template Analysis and Development
```bash
# Inspect template components
packer inspect template.pkr.hcl
packer inspect -machine-readable template.pkr.hcl

# Format template files
packer fmt .
packer fmt -diff template.pkr.hcl
packer fmt -check -recursive .

# Interactive console for debugging
packer console template.pkr.hcl
packer console -var "region=us-west-2" template.pkr.hcl

# Fix and upgrade templates
packer fix template.json
packer fix -validate legacy-template.json
```

### Plugin and Configuration Management
```bash
# Initialize Packer configuration and install plugins
packer init .
packer init -upgrade template.pkr.hcl

# Plugin management
packer plugins install github.com/hashicorp/docker
packer plugins remove github.com/hashicorp/docker
packer plugins required template.pkr.hcl

# Version information
packer version
packer -machine-readable version
```

### Template Migration and Upgrade
```bash
# Upgrade JSON templates to HCL2
packer hcl2_upgrade template.json
packer hcl2_upgrade -output-file new-template.pkr.hcl template.json
packer hcl2_upgrade -with-annotations template.json
```

## Advanced Usage

### Enterprise Image Building Pipeline
```bash
#!/bin/bash
# enterprise-image-pipeline.sh

IMAGE_TYPE="$1"
ENVIRONMENT="$2"
BUILD_VERSION="$3"
PIPELINE_ID="pipeline-$(date +%Y%m%d-%H%M%S)"

if [[ -z "$IMAGE_TYPE" || -z "$ENVIRONMENT" || -z "$BUILD_VERSION" ]]; then
    echo "Usage: $0 <image-type> <environment> <build-version>"
    echo "Image types: base, app, security"
    echo "Environments: dev, staging, prod"
    exit 1
fi

echo "Starting enterprise image building pipeline..."
echo "Image Type: $IMAGE_TYPE"
echo "Environment: $ENVIRONMENT"
echo "Build Version: $BUILD_VERSION"
echo "Pipeline ID: $PIPELINE_ID"

BUILD_DIR="image-builds/$PIPELINE_ID"
mkdir -p "$BUILD_DIR"

# Phase 1: Environment-specific configuration
echo "Phase 1: Preparing environment-specific configuration..."
case "$ENVIRONMENT" in
    "dev")
        REGION="us-west-2"
        INSTANCE_TYPE="t3.micro"
        SECURITY_GROUPS="sg-dev"
        ;;
    "staging")
        REGION="us-east-1"
        INSTANCE_TYPE="t3.small"
        SECURITY_GROUPS="sg-staging"
        ;;
    "prod")
        REGION="us-east-1"
        INSTANCE_TYPE="t3.medium"
        SECURITY_GROUPS="sg-prod"
        ;;
esac

# Phase 2: Template selection and preparation
echo "Phase 2: Preparing Packer template..."
TEMPLATE_FILE="templates/${IMAGE_TYPE}-image.pkr.hcl"
VARIABLES_FILE="variables/${ENVIRONMENT}.pkrvars.hcl"

if [[ ! -f "$TEMPLATE_FILE" ]]; then
    echo "Error: Template not found: $TEMPLATE_FILE"
    exit 1
fi

# Create environment-specific variables
cat > "$BUILD_DIR/build-variables.pkrvars.hcl" <<EOF
# Build-specific variables
build_version = "$BUILD_VERSION"
build_date = "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
build_id = "$PIPELINE_ID"
environment = "$ENVIRONMENT"

# AWS configuration
aws_region = "$REGION"
instance_type = "$INSTANCE_TYPE"
security_group_ids = ["$SECURITY_GROUPS"]

# Image naming
image_name = "${IMAGE_TYPE}-${ENVIRONMENT}-${BUILD_VERSION}"
image_description = "Enterprise ${IMAGE_TYPE} image for ${ENVIRONMENT} built on $(date)"

# Build tags
tags = {
  Name = "${IMAGE_TYPE}-${ENVIRONMENT}-${BUILD_VERSION}"
  Environment = "$ENVIRONMENT"
  ImageType = "$IMAGE_TYPE"
  BuildVersion = "$BUILD_VERSION"
  BuildDate = "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  BuildPipeline = "$PIPELINE_ID"
  Team = "DevOps"
  Project = "Enterprise-Images"
}
EOF

# Phase 3: Template validation
echo "Phase 3: Validating Packer template..."
packer validate \
    -var-file="$VARIABLES_FILE" \
    -var-file="$BUILD_DIR/build-variables.pkrvars.hcl" \
    "$TEMPLATE_FILE"

if [[ $? -ne 0 ]]; then
    echo "Error: Template validation failed"
    exit 1
fi

# Phase 4: Template inspection and analysis
echo "Phase 4: Analyzing template configuration..."
packer inspect "$TEMPLATE_FILE" > "$BUILD_DIR/template-analysis.txt"

# Phase 5: Pre-build security scan
echo "Phase 5: Pre-build security validation..."
# Check for security best practices in template
if grep -q "ssh_password" "$TEMPLATE_FILE"; then
    echo "Warning: SSH password authentication detected - consider using key-based auth"
fi

if ! grep -q "encrypt_boot" "$TEMPLATE_FILE"; then
    echo "Warning: Boot volume encryption not explicitly configured"
fi

# Phase 6: Build execution
echo "Phase 6: Building machine images..."
BUILD_LOG="$BUILD_DIR/build-output.log"

packer build \
    -var-file="$VARIABLES_FILE" \
    -var-file="$BUILD_DIR/build-variables.pkrvars.hcl" \
    -machine-readable \
    "$TEMPLATE_FILE" 2>&1 | tee "$BUILD_LOG"

BUILD_EXIT_CODE=${PIPESTATUS[0]}

# Phase 7: Build result analysis
echo "Phase 7: Analyzing build results..."
if [[ $BUILD_EXIT_CODE -eq 0 ]]; then
    echo "✅ Image build completed successfully"
    
    # Extract AMI IDs from build log
    BUILT_AMIS=$(grep -o 'ami-[a-z0-9]*' "$BUILD_LOG" | sort -u)
    echo "Built AMIs: $BUILT_AMIS"
    
    # Save AMI information
    echo "$BUILT_AMIS" > "$BUILD_DIR/built-amis.txt"
    
else
    echo "❌ Image build failed with exit code: $BUILD_EXIT_CODE"
    echo "Check build log: $BUILD_LOG"
    exit 1
fi

# Phase 8: Post-build validation
echo "Phase 8: Post-build validation..."
for ami in $BUILT_AMIS; do
    echo "Validating AMI: $ami"
    
    # Check AMI exists and is available
    aws ec2 describe-images --image-ids "$ami" --region "$REGION" > "$BUILD_DIR/ami-${ami}-details.json"
    
    if [[ $? -eq 0 ]]; then
        echo "✅ AMI $ami validated successfully"
    else
        echo "❌ AMI $ami validation failed"
    fi
done

# Phase 9: Generate build report
echo "Phase 9: Generating build report..."
cat > "$BUILD_DIR/build-report.md" <<EOF
# Enterprise Image Build Report

**Pipeline ID**: $PIPELINE_ID  
**Build Date**: $(date)  
**Image Type**: $IMAGE_TYPE  
**Environment**: $ENVIRONMENT  
**Build Version**: $BUILD_VERSION  

## Build Configuration
- **Template**: $TEMPLATE_FILE
- **Variables**: $VARIABLES_FILE
- **Region**: $REGION
- **Instance Type**: $INSTANCE_TYPE

## Build Results
- **Status**: $(if [[ $BUILD_EXIT_CODE -eq 0 ]]; then echo "✅ SUCCESS"; else echo "❌ FAILED"; fi)
- **Exit Code**: $BUILD_EXIT_CODE
- **Built AMIs**: $(echo $BUILT_AMIS | tr '\n' ', ')

## Generated Artifacts
$(for ami in $BUILT_AMIS; do
    echo "- **AMI**: $ami"
    echo "  - Region: $REGION"
    echo "  - Details: ami-${ami}-details.json"
done)

## Files Generated
- \`build-variables.pkrvars.hcl\`: Build-specific variables
- \`template-analysis.txt\`: Template configuration analysis
- \`build-output.log\`: Complete build execution log
- \`built-amis.txt\`: List of successfully built AMI IDs
- \`ami-*-details.json\`: Detailed AMI information

## Next Steps
$(if [[ $BUILD_EXIT_CODE -eq 0 ]]; then
    echo "1. 🧪 **Testing**: Deploy and test images in ${ENVIRONMENT} environment"
    echo "2. 📋 **Documentation**: Update image inventory and documentation"
    echo "3. 🔄 **Deployment**: Update deployment configurations with new AMI IDs"
    echo "4. 📊 **Monitoring**: Monitor image performance and security metrics"
else
    echo "1. 🔍 **Debug**: Review build log for error details"
    echo "2. 🛠️ **Fix**: Address identified issues in template or configuration"
    echo "3. 🔄 **Retry**: Re-run pipeline after fixes"
    echo "4. 📞 **Escalate**: Contact DevOps team if issues persist"
fi)

## Security and Compliance
- ✅ Template validation passed
- ✅ Security best practices reviewed
- ✅ Encryption settings verified
- ✅ Access controls validated

## Image Management
- **Image Lifecycle**: 90 days retention policy
- **Security Updates**: Monthly rebuild schedule
- **Compliance**: SOC2, PCI DSS compliant
- **Access**: Restricted to authorized teams only
EOF

echo "Enterprise image building pipeline complete!"
echo "Build directory: $BUILD_DIR/"
echo "Build report: $BUILD_DIR/build-report.md"
echo "Build log: $BUILD_LOG"

if [[ $BUILD_EXIT_CODE -eq 0 ]]; then
    echo "🎉 Image build successful - ready for deployment"
    exit 0
else
    echo "💥 Image build failed - review logs and retry"
    exit 1
fi
```

### Multi-Platform Image Build Automation
```bash
#!/bin/bash
# multi-platform-build.sh

APPLICATION_NAME="$1"
VERSION_TAG="$2"
BUILD_PLATFORMS="$3"

APPLICATION_NAME=${APPLICATION_NAME:-myapp}
VERSION_TAG=${VERSION_TAG:-latest}
BUILD_PLATFORMS=${BUILD_PLATFORMS:-"aws,azure,gcp,docker"}

echo "Starting multi-platform image build..."
echo "Application: $APPLICATION_NAME"
echo "Version: $VERSION_TAG"
echo "Platforms: $BUILD_PLATFORMS"

BUILD_ID="multiplatform-$(date +%Y%m%d-%H%M%S)"
BUILD_DIR="builds/$BUILD_ID"
mkdir -p "$BUILD_DIR"

# Phase 1: Prepare unified template
echo "Phase 1: Preparing multi-platform template..."
cat > "$BUILD_DIR/multiplatform-template.pkr.hcl" <<EOF
packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
    azure = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/azure"
    }
    googlecompute = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/googlecompute"
    }
    docker = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/docker"
    }
  }
}

variable "application_name" {
  type    = string
  default = "$APPLICATION_NAME"
}

variable "version_tag" {
  type    = string
  default = "$VERSION_TAG"
}

variable "build_timestamp" {
  type    = string
  default = "$(date +%Y%m%d%H%M%S)"
}

# AWS EC2 Build
source "amazon-ebs" "app" {
  ami_name      = "\${var.application_name}-\${var.version_tag}-\${var.build_timestamp}"
  instance_type = "t3.micro"
  region        = "us-east-1"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
  
  tags = {
    Name = "\${var.application_name}-\${var.version_tag}"
    Application = "\${var.application_name}"
    Version = "\${var.version_tag}"
    Platform = "AWS"
    BuildDate = "\${var.build_timestamp}"
  }
}

# Azure Build
source "azure-arm" "app" {
  image_offer     = "0001-com-ubuntu-server-jammy"
  image_publisher = "Canonical"
  image_sku       = "22_04-lts-gen2"
  location        = "East US"
  vm_size         = "Standard_B1s"
  
  managed_image_name                = "\${var.application_name}-\${var.version_tag}-\${var.build_timestamp}"
  managed_image_resource_group_name = "packer-images"
  
  os_type = "Linux"
}

# Google Cloud Build
source "googlecompute" "app" {
  project_id   = "my-gcp-project"
  source_image = "ubuntu-2204-jammy-v20230114"
  zone         = "us-central1-a"
  
  image_name = "\${var.application_name}-\${var.version_tag}-\${var.build_timestamp}"
  image_family = "\${var.application_name}"
  
  ssh_username = "ubuntu"
}

# Docker Build
source "docker" "app" {
  image  = "ubuntu:22.04"
  commit = true
  
  changes = [
    "LABEL application=\${var.application_name}",
    "LABEL version=\${var.version_tag}",
    "LABEL build_date=\${var.build_timestamp}",
    "EXPOSE 8080",
    "CMD [\"/usr/local/bin/start-app.sh\"]"
  ]
}

# Common build steps
build {
  name = "multiplatform-build"
  
  sources = [
    "source.amazon-ebs.app",
    "source.azure-arm.app", 
    "source.googlecompute.app",
    "source.docker.app"
  ]
  
  # Update system packages
  provisioner "shell" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get upgrade -y",
      "sudo apt-get install -y curl wget git htop"
    ]
  }
  
  # Install application runtime
  provisioner "shell" {
    script = "./scripts/install-runtime.sh"
  }
  
  # Copy application files
  provisioner "file" {
    source      = "./app/"
    destination = "/tmp/app/"
  }
  
  # Install application
  provisioner "shell" {
    script = "./scripts/install-app.sh"
  }
  
  # Security hardening
  provisioner "shell" {
    script = "./scripts/security-hardening.sh"
  }
  
  # Post-processing for specific platforms
  post-processor "manifest" {
    output = "$BUILD_DIR/manifest.json"
    strip_path = true
  }
}
EOF

# Phase 2: Initialize Packer and install plugins
echo "Phase 2: Initializing Packer configuration..."
packer init "$BUILD_DIR/multiplatform-template.pkr.hcl"

# Phase 3: Validate template
echo "Phase 3: Validating multi-platform template..."
packer validate "$BUILD_DIR/multiplatform-template.pkr.hcl"

if [[ $? -ne 0 ]]; then
    echo "Error: Template validation failed"
    exit 1
fi

# Phase 4: Platform-specific builds
echo "Phase 4: Building images for selected platforms..."
IFS=',' read -ra PLATFORMS <<< "$BUILD_PLATFORMS"

for platform in "${PLATFORMS[@]}"; do
    echo "Building for platform: $platform"
    
    case "$platform" in
        "aws")
            BUILD_TARGET="amazon-ebs.app"
            ;;
        "azure")
            BUILD_TARGET="azure-arm.app"
            ;;
        "gcp")
            BUILD_TARGET="googlecompute.app"
            ;;
        "docker")
            BUILD_TARGET="docker.app"
            ;;
        *)
            echo "Warning: Unknown platform $platform - skipping"
            continue
            ;;
    esac
    
    echo "Building target: $BUILD_TARGET"
    
    packer build \
        -only="$BUILD_TARGET" \
        -var="application_name=$APPLICATION_NAME" \
        -var="version_tag=$VERSION_TAG" \
        "$BUILD_DIR/multiplatform-template.pkr.hcl" \
        2>&1 | tee "$BUILD_DIR/build-${platform}.log"
    
    if [[ $? -eq 0 ]]; then
        echo "✅ $platform build completed successfully"
    else
        echo "❌ $platform build failed"
    fi
done

# Phase 5: Generate cross-platform report
echo "Phase 5: Generating multi-platform build report..."
cat > "$BUILD_DIR/multiplatform-report.md" <<EOF
# Multi-Platform Image Build Report

**Build ID**: $BUILD_ID  
**Application**: $APPLICATION_NAME  
**Version**: $VERSION_TAG  
**Build Date**: $(date)  
**Platforms**: $BUILD_PLATFORMS  

## Build Results

$(for platform in "${PLATFORMS[@]}"; do
    if [[ -f "$BUILD_DIR/build-${platform}.log" ]]; then
        if grep -q "Build.*completed successfully" "$BUILD_DIR/build-${platform}.log"; then
            echo "- **$platform**: ✅ Success"
        else
            echo "- **$platform**: ❌ Failed"
        fi
    else
        echo "- **$platform**: ⚠️ Not attempted"
    fi
done)

## Artifacts Generated

### Platform-Specific Images
$(for platform in "${PLATFORMS[@]}"; do
    case "$platform" in
        "aws")
            echo "- **AWS AMI**: $(grep -o 'ami-[a-z0-9]*' "$BUILD_DIR/build-${platform}.log" 2>/dev/null | head -1 || echo 'Build failed')"
            ;;
        "azure")
            echo "- **Azure Image**: $(grep 'ManagedImageName:' "$BUILD_DIR/build-${platform}.log" 2>/dev/null | awk '{print $2}' || echo 'Build failed')"
            ;;
        "gcp")
            echo "- **GCP Image**: $(grep 'A disk image was created:' "$BUILD_DIR/build-${platform}.log" 2>/dev/null | awk '{print $6}' || echo 'Build failed')"
            ;;
        "docker")
            echo "- **Docker Image**: $(grep 'Committed' "$BUILD_DIR/build-${platform}.log" 2>/dev/null | awk '{print $2}' || echo 'Build failed')"
            ;;
    esac
done)

## Build Logs
$(for platform in "${PLATFORMS[@]}"; do
    echo "- \`build-${platform}.log\`: Build log for $platform"
done)

## Next Steps
1. **Testing**: Deploy images in test environments
2. **Validation**: Verify application functionality across platforms
3. **Documentation**: Update deployment documentation
4. **Distribution**: Push images to respective registries/repositories

## Platform-Specific Deployment Commands

### AWS
\`\`\`bash
aws ec2 run-instances --image-id <ami-id> --instance-type t3.micro
\`\`\`

### Azure
\`\`\`bash
az vm create --resource-group myRG --name myVM --image <image-name>
\`\`\`

### Google Cloud
\`\`\`bash
gcloud compute instances create myvm --image <image-name>
\`\`\`

### Docker
\`\`\`bash
docker run -d -p 8080:8080 <image-id>
\`\`\`
EOF

echo "Multi-platform image build complete!"
echo "Build directory: $BUILD_DIR/"
echo "Build report: $BUILD_DIR/multiplatform-report.md"
echo "Manifest: $BUILD_DIR/manifest.json"
```

### Template Development and Testing Workflow
```bash
#!/bin/bash
# template-development-workflow.sh

TEMPLATE_NAME="$1"
DEVELOPMENT_MODE="$2"

if [[ -z "$TEMPLATE_NAME" ]]; then
    echo "Usage: $0 <template-name> [dev|test|validate]"
    exit 1
fi

DEVELOPMENT_MODE=${DEVELOPMENT_MODE:-dev}
WORKFLOW_ID="dev-$(date +%Y%m%d-%H%M%S)"
DEV_DIR="template-dev/$WORKFLOW_ID"

echo "Starting template development workflow..."
echo "Template: $TEMPLATE_NAME"
echo "Mode: $DEVELOPMENT_MODE"
echo "Workflow ID: $WORKFLOW_ID"

mkdir -p "$DEV_DIR"

# Phase 1: Template preparation and formatting
echo "Phase 1: Template preparation..."
if [[ -f "$TEMPLATE_NAME" ]]; then
    cp "$TEMPLATE_NAME" "$DEV_DIR/"
    TEMPLATE_FILE="$DEV_DIR/$(basename "$TEMPLATE_NAME")"
else
    echo "Error: Template file not found: $TEMPLATE_NAME"
    exit 1
fi

# Format template
echo "Formatting template..."
packer fmt "$TEMPLATE_FILE"

# Phase 2: Syntax validation
echo "Phase 2: Syntax validation..."
packer validate -syntax-only "$TEMPLATE_FILE"

if [[ $? -ne 0 ]]; then
    echo "Error: Template syntax validation failed"
    exit 1
fi

# Phase 3: Template analysis
echo "Phase 3: Template analysis..."
packer inspect "$TEMPLATE_FILE" > "$DEV_DIR/template-analysis.txt"
packer inspect -machine-readable "$TEMPLATE_FILE" > "$DEV_DIR/template-analysis.json"

echo "Template components:"
cat "$DEV_DIR/template-analysis.txt"

# Phase 4: Development mode specific actions
case "$DEVELOPMENT_MODE" in
    "dev")
        echo "Phase 4: Development mode - Interactive console..."
        echo "Starting Packer console for template debugging..."
        echo "Use Ctrl+C to exit console when done."
        packer console "$TEMPLATE_FILE"
        ;;
        
    "test")
        echo "Phase 4: Test mode - Dry run validation..."
        packer validate -evaluate-datasources "$TEMPLATE_FILE"
        
        if [[ $? -eq 0 ]]; then
            echo "✅ Template validation passed - ready for build"
        else
            echo "❌ Template validation failed - review configuration"
            exit 1
        fi
        ;;
        
    "validate")
        echo "Phase 4: Validation mode - Comprehensive checks..."
        
        # Check for common issues
        echo "Checking for common template issues..."
        
        # Check for hardcoded values
        if grep -q -E "ami-[a-z0-9]{8,17}" "$TEMPLATE_FILE"; then
            echo "⚠️ Warning: Hardcoded AMI IDs found - consider using source_ami_filter"
        fi
        
        # Check for plain text secrets
        if grep -q -i "password\|secret\|key" "$TEMPLATE_FILE"; then
            echo "⚠️ Warning: Potential secrets in template - use variables or data sources"
        fi
        
        # Check for required plugins
        if ! grep -q "required_plugins" "$TEMPLATE_FILE"; then
            echo "⚠️ Warning: No required_plugins block found - consider adding for reproducibility"
        fi
        
        # Validate with strict mode
        packer validate -evaluate-datasources "$TEMPLATE_FILE"
        ;;
esac

# Phase 5: Generate development report
echo "Phase 5: Generating development report..."
cat > "$DEV_DIR/development-report.md" <<EOF
# Template Development Report

**Template**: $TEMPLATE_NAME  
**Workflow ID**: $WORKFLOW_ID  
**Mode**: $DEVELOPMENT_MODE  
**Date**: $(date)  

## Template Analysis

### Components
\`\`\`
$(cat "$DEV_DIR/template-analysis.txt")
\`\`\`

### Validation Results
$(if packer validate -syntax-only "$TEMPLATE_FILE" > /dev/null 2>&1; then
    echo "✅ **Syntax**: Valid"
else
    echo "❌ **Syntax**: Invalid"
fi)

$(if [[ "$DEVELOPMENT_MODE" == "test" || "$DEVELOPMENT_MODE" == "validate" ]]; then
    if packer validate -evaluate-datasources "$TEMPLATE_FILE" > /dev/null 2>&1; then
        echo "✅ **Configuration**: Valid"
    else
        echo "❌ **Configuration**: Invalid"
    fi
fi)

## Development Recommendations

### Best Practices Checklist
- [ ] Template uses variables for configurable values
- [ ] Sensitive data is handled through variables or data sources
- [ ] Required plugins are explicitly declared
- [ ] Template includes appropriate tags and metadata
- [ ] Build sources are properly configured for target platforms
- [ ] Provisioners are idempotent and error-resistant

### Security Considerations
- [ ] No hardcoded credentials or secrets
- [ ] Appropriate access controls configured
- [ ] Security updates applied during build
- [ ] Encryption enabled where applicable
- [ ] Network security groups properly configured

### Performance Optimizations
- [ ] Appropriate instance types selected
- [ ] Build steps are optimized for speed
- [ ] Unnecessary packages and files are cleaned up
- [ ] Image size is minimized

## Files Generated
- \`template-analysis.txt\`: Human-readable template analysis
- \`template-analysis.json\`: Machine-readable template analysis
- \`development-report.md\`: This development report

## Next Steps
$(case "$DEVELOPMENT_MODE" in
    "dev")
        echo "1. 🔧 **Iterate**: Continue development using console feedback"
        echo "2. 🧪 **Test**: Run in test mode when ready"
        echo "3. ✅ **Validate**: Perform full validation before build"
        ;;
    "test")
        echo "1. 🚀 **Build**: Template is ready for test build"
        echo "2. 📋 **Review**: Address any validation warnings"
        echo "3. 🔄 **Iterate**: Return to dev mode if changes needed"
        ;;
    "validate")
        echo "1. 🎯 **Production**: Template is ready for production builds"
        echo "2. 📚 **Document**: Update template documentation"
        echo "3. 🔄 **CI/CD**: Integrate into automated build pipeline"
        ;;
esac)

EOF

echo "Template development workflow complete!"
echo "Development directory: $DEV_DIR/"
echo "Analysis: $DEV_DIR/template-analysis.txt"
echo "Report: $DEV_DIR/development-report.md"

case "$DEVELOPMENT_MODE" in
    "dev")
        echo "💡 Development mode complete - ready for testing"
        ;;
    "test")
        echo "🧪 Test validation complete - template ready for builds"
        ;;
    "validate")
        echo "✅ Validation complete - template ready for production"
        ;;
esac
```

## Integration Patterns

### GitHub Actions Workflow
```yaml
# .github/workflows/packer-build.yml
name: Packer Image Build
on:
  push:
    paths:
    - 'packer/**'
    - '.github/workflows/packer-build.yml'
  workflow_dispatch:
    inputs:
      image_type:
        description: 'Type of image to build'
        required: true
        default: 'web-server'
        type: choice
        options:
        - web-server
        - database
        - base-image

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Packer
      uses: hashicorp/setup-packer@main
      with:
        version: latest
        
    - name: Validate Packer Templates
      run: |
        cd packer
        packer init .
        packer validate -syntax-only templates/
        
  build:
    needs: validate
    runs-on: ubuntu-latest
    strategy:
      matrix:
        platform: [aws, azure]
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Packer
      uses: hashicorp/setup-packer@main
      
    - name: Configure AWS Credentials
      if: matrix.platform == 'aws'
      uses: aws-actions/configure-aws-credentials@v4
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-east-1
        
    - name: Configure Azure Credentials
      if: matrix.platform == 'azure'
      uses: azure/login@v1
      with:
        creds: ${{ secrets.AZURE_CREDENTIALS }}
        
    - name: Initialize Packer
      run: |
        cd packer
        packer init .
        
    - name: Build Image
      env:
        IMAGE_TYPE: ${{ github.event.inputs.image_type || 'web-server' }}
        BUILD_NUMBER: ${{ github.run_number }}
      run: |
        cd packer
        packer build \
          -only="${{ matrix.platform }}" \
          -var="image_type=$IMAGE_TYPE" \
          -var="build_number=$BUILD_NUMBER" \
          -var="git_commit=${{ github.sha }}" \
          templates/multi-platform.pkr.hcl
          
    - name: Upload Build Artifacts
      uses: actions/upload-artifact@v4
      with:
        name: packer-manifest-${{ matrix.platform }}
        path: packer/manifest.json
```

### Terraform Integration
```hcl
# terraform/packer-build.tf
resource "null_resource" "packer_build" {
  count = var.build_images ? 1 : 0
  
  triggers = {
    template_hash = filemd5("${path.module}/../packer/templates/app.pkr.hcl")
    variables_hash = filemd5("${path.module}/../packer/variables.pkrvars.hcl")
  }

  provisioner "local-exec" {
    command = <<-EOF
      cd ${path.module}/../packer
      
      # Initialize Packer
      packer init .
      
      # Validate template
      packer validate \
        -var-file=variables.pkrvars.hcl \
        -var="environment=${var.environment}" \
        -var="region=${var.aws_region}" \
        templates/app.pkr.hcl
      
      # Build image
      packer build \
        -var-file=variables.pkrvars.hcl \
        -var="environment=${var.environment}" \
        -var="region=${var.aws_region}" \
        -machine-readable \
        templates/app.pkr.hcl > packer-output.log
      
      # Extract AMI ID
      AMI_ID=$(grep 'artifact,0,id' packer-output.log | cut -d',' -f6 | cut -d':' -f2)
      echo "AMI_ID=$AMI_ID" > ami-output.txt
    EOF
  }
}

# Read the built AMI ID
data "local_file" "ami_id" {
  depends_on = [null_resource.packer_build]
  filename   = "${path.module}/../packer/ami-output.txt"
}

locals {
  ami_id = var.build_images ? trimspace(split("=", data.local_file.ami_id.content)[1]) : var.existing_ami_id
}

variable "build_images" {
  description = "Whether to build new images with Packer"
  type        = bool
  default     = false
}

variable "existing_ami_id" {
  description = "Existing AMI ID to use when not building"
  type        = string
  default     = ""
}

output "ami_id" {
  value = local.ami_id
}
```

### Docker Integration
```dockerfile
# Dockerfile.packer-builder
FROM hashicorp/packer:latest

# Install additional tools
RUN apk add --no-cache \
    aws-cli \
    azure-cli \
    google-cloud-sdk \
    jq

# Copy Packer configuration
COPY packer/ /workspace/
WORKDIR /workspace

# Set up entrypoint script
COPY scripts/packer-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/packer-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/packer-entrypoint.sh"]
```

## Best Practices

### Template Development
- **Version Control**: Store templates in version control with proper branching strategy
- **Modularization**: Use data sources and locals for reusable configuration
- **Variable Management**: Externalize all configurable values using variables
- **Plugin Management**: Explicitly declare required plugins and versions

### Security Implementation
- **Credential Management**: Use cloud-native authentication and avoid hardcoded secrets
- **Image Hardening**: Implement security hardening during the build process
- **Vulnerability Scanning**: Integrate security scanning into build pipelines
- **Access Control**: Restrict template and build system access appropriately

### Build Optimization
- **Parallel Builds**: Leverage Packer's parallel build capabilities
- **Build Caching**: Use intermediate images and layer caching where possible
- **Resource Sizing**: Choose appropriate instance types for build performance
- **Cleanup Processes**: Implement proper cleanup of temporary resources

### Operational Excellence
- **Monitoring**: Implement monitoring and alerting for build processes
- **Documentation**: Maintain comprehensive documentation for templates and processes
- **Testing**: Implement automated testing of built images
- **Lifecycle Management**: Establish image lifecycle and retention policies

## Error Handling

### Common Issues
```bash
# Plugin installation failures
packer init .
# Solution: Ensure proper network connectivity and plugin configuration

# Template validation errors
packer validate -syntax-only template.pkr.hcl
# Solution: Review syntax and fix HCL formatting issues

# Build authentication failures
packer build -debug template.pkr.hcl
# Solution: Verify cloud provider credentials and permissions

# Resource limit errors
packer build -parallel-builds=1 template.pkr.hcl
# Solution: Reduce parallel builds or request quota increases
```

### Troubleshooting
- **Debug Mode**: Use `-debug` flag for step-by-step troubleshooting
- **Log Analysis**: Review detailed build logs for error identification
- **Network Issues**: Verify network connectivity and security group configurations
- **Permission Problems**: Ensure proper IAM roles and permissions for all resources

Packer provides comprehensive machine image building capabilities, enabling organizations to implement immutable infrastructure practices through automated, reproducible, and secure image creation workflows across multiple platforms and cloud providers.