# Velero

Kubernetes-native backup and disaster recovery solution for protecting cluster resources and persistent volumes.

## Description

Velero is a comprehensive open-source tool designed to safely backup and restore, perform disaster recovery, and migrate Kubernetes cluster resources and persistent volumes. Developed by VMware Tanzu, Velero provides enterprise-grade backup capabilities for Kubernetes environments, enabling organizations to protect their workloads against data loss, cluster failures, and facilitate cluster migrations. The tool supports multiple cloud providers and storage backends, offers flexible scheduling options, and includes advanced features like namespace remapping and backup hooks for stateful applications.

## MCP Tools

### Installation and Setup
- **`velero_install`** - Install Velero with storage provider configuration using real velero CLI
- **`velero_version`** - Get Velero version information using velero CLI

### Backup Management
- **`velero_backup_create`** - Create an on-demand backup using velero CLI
- **`velero_backup_get`** - Get list of backups using velero CLI
- **`velero_backup_describe`** - Describe a specific backup using velero CLI
- **`velero_backup_delete`** - Delete a backup using velero CLI
- **`velero_backup_logs`** - Get backup logs using velero CLI

### Restore Management
- **`velero_restore_create`** - Create a restore from backup using velero CLI
- **`velero_restore_get`** - Get list of restores using velero CLI
- **`velero_restore_describe`** - Describe a specific restore using velero CLI
- **`velero_restore_logs`** - Get restore logs using velero CLI

### Schedule Management
- **`velero_create_schedule`** - Create a backup schedule using velero CLI
- **`velero_schedule_get`** - Get list of backup schedules using velero CLI
- **`velero_schedule_delete`** - Delete a backup schedule using velero CLI

### Storage Management
- **`velero_backup_location_create`** - Create backup storage location using velero CLI
- **`velero_backup_location_get`** - Get backup storage locations using velero CLI

## Real CLI Commands Used

### Installation Commands
- `velero install --provider <provider> --bucket <bucket> --plugins <plugins>` - Install Velero on cluster
- `velero version --client-only` - Get client version
- `velero version` - Get client and server version

### Backup Commands
- `velero backup create <name>` - Create backup
- `velero backup get` - List backups
- `velero backup describe <name>` - Get backup details
- `velero backup delete <name>` - Delete backup
- `velero backup logs <name>` - Get backup logs

### Restore Commands
- `velero restore create <name> --from-backup <backup>` - Create restore
- `velero restore get` - List restores
- `velero restore describe <name>` - Get restore details
- `velero restore logs <name>` - Get restore logs

### Schedule Commands
- `velero schedule create <name> --schedule <cron>` - Create backup schedule
- `velero schedule get` - List schedules
- `velero schedule delete <name>` - Delete schedule

### Storage Location Commands
- `velero backup-location create <name> --provider <provider> --bucket <bucket>` - Create storage location
- `velero backup-location get` - List storage locations

## Use Cases

### Disaster Recovery
- **Cluster Backup**: Complete cluster state backup including resources and data
- **Point-in-Time Recovery**: Restore cluster to specific point in time
- **Cross-Region Recovery**: Restore cluster in different region for disaster scenarios
- **Application Recovery**: Selective application and namespace recovery

### Data Protection
- **Automated Backups**: Scheduled backups of critical workloads and data
- **Persistent Volume Protection**: Backup and restore of persistent storage
- **Configuration Backup**: Kubernetes resource configuration preservation
- **Incremental Backups**: Efficient storage using incremental backup strategies

### Migration and Mobility
- **Cluster Migration**: Move workloads between different Kubernetes clusters
- **Cloud Migration**: Migrate from on-premises to cloud or between cloud providers
- **Environment Cloning**: Clone production environments to staging/development
- **Multi-Cloud Strategy**: Implement multi-cloud backup and mobility strategies

### Compliance and Governance
- **Retention Policies**: Implement backup retention for compliance requirements
- **Audit Trail**: Track backup and restore operations for compliance
- **Data Sovereignty**: Control backup storage location for regulatory compliance
- **Business Continuity**: Ensure business continuity with reliable backup strategies

## Configuration Examples

### Basic Velero Installation
```bash
# Install Velero with AWS provider
velero install \
    --provider aws \
    --plugins velero/velero-plugin-for-aws:v1.8.0 \
    --bucket my-backup-bucket \
    --backup-location-config region=us-west-2 \
    --secret-file ./credentials-velero

# Install with Google Cloud provider
velero install \
    --provider gcp \
    --plugins velero/velero-plugin-for-gcp:v1.8.0 \
    --bucket my-gcp-backup-bucket \
    --secret-file ./gcp-credentials.json

# Install with Azure provider
velero install \
    --provider azure \
    --plugins velero/velero-plugin-for-microsoft-azure:v1.8.0 \
    --bucket my-azure-container \
    --backup-location-config \
        resourceGroup=velero-backups,storageAccount=velerostorage \
    --secret-file ./azure-credentials

# Verify installation
velero version
kubectl get pods -n velero
```

### Basic Backup Operations
```bash
# Create immediate backup of entire cluster
velero backup create full-cluster-backup

# Create backup of specific namespace
velero backup create app-backup \
    --include-namespaces production

# Create backup with labels
velero backup create labeled-backup \
    --labels environment=production,app=web

# Create backup excluding certain namespaces
velero backup create selective-backup \
    --exclude-namespaces kube-system,kube-public

# List all backups
velero backup get

# Get detailed backup information
velero backup describe full-cluster-backup

# Download backup logs
velero backup logs full-cluster-backup
```

### Backup Scheduling
```bash
# Create daily backup schedule
velero schedule create daily-backup \
    --schedule="0 2 * * *" \
    --ttl 720h

# Create weekly backup with timezone
velero schedule create weekly-backup \
    --schedule="CRON_TZ=America/New_York 0 1 * * 0" \
    --ttl 2160h

# Create schedule for specific namespace
velero schedule create app-schedule \
    --schedule="0 6 * * *" \
    --include-namespaces production \
    --ttl 168h

# List schedules
velero schedule get

# Delete schedule
velero schedule delete daily-backup --confirm
```

### Restore Operations
```bash
# Restore from backup
velero restore create restore-20241117 \
    --from-backup full-cluster-backup

# Restore with namespace mapping
velero restore create restore-to-staging \
    --from-backup production-backup \
    --namespace-mappings production:staging

# Restore specific namespaces only
velero restore create partial-restore \
    --from-backup full-cluster-backup \
    --include-namespaces app1,app2

# Restore excluding certain resources
velero restore create selective-restore \
    --from-backup full-cluster-backup \
    --exclude-resources secrets

# Check restore status
velero restore get
velero restore describe restore-20241117

# View restore logs
velero restore logs restore-20241117
```

## Advanced Usage

### Multi-Cloud Backup Strategy
```bash
#!/bin/bash
# multi-cloud-backup-setup.sh

echo "Setting up multi-cloud backup strategy with Velero..."

# Primary backup location (AWS)
velero backup-location create aws-primary \
    --provider aws \
    --bucket primary-backup-bucket \
    --config region=us-west-2

# Secondary backup location (GCP)
velero backup-location create gcp-secondary \
    --provider gcp \
    --bucket secondary-backup-bucket \
    --config projectId=my-gcp-project

# Create schedule with primary location
velero schedule create primary-daily \
    --schedule="0 2 * * *" \
    --storage-location aws-primary \
    --ttl 720h

# Create schedule with secondary location
velero schedule create secondary-weekly \
    --schedule="0 3 * * 0" \
    --storage-location gcp-secondary \
    --ttl 2160h

echo "Multi-cloud backup strategy configured successfully!"
```

### Application-Specific Backup with Hooks
```bash
# Create backup with pre/post hooks for database
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: mysql-backup
  annotations:
    pre.hook.backup.velero.io/command: '["/bin/bash", "-c", "mysqldump -u root -p\$MYSQL_ROOT_PASSWORD --all-databases > /backup/db-dump.sql"]'
    post.hook.backup.velero.io/command: '["/bin/bash", "-c", "rm -f /backup/db-dump.sql"]'
spec:
  containers:
  - name: mysql
    image: mysql:8.0
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: secretpassword
    volumeMounts:
    - name: backup-volume
      mountPath: /backup
  volumes:
  - name: backup-volume
    persistentVolumeClaim:
      claimName: mysql-backup-pvc
EOF

# Create backup that will execute hooks
velero backup create mysql-app-backup \
    --include-namespaces database \
    --wait
```

### Automated Backup Validation
```bash
#!/bin/bash
# backup-validation.sh

BACKUP_NAME="$1"
if [[ -z "$BACKUP_NAME" ]]; then
    echo "Usage: $0 <backup-name>"
    exit 1
fi

echo "Validating backup: $BACKUP_NAME"

# Check backup status
STATUS=$(velero backup describe $BACKUP_NAME --output json | jq -r '.status.phase')

if [[ "$STATUS" == "Completed" ]]; then
    echo "✅ Backup completed successfully"
    
    # Get backup details
    RESOURCES=$(velero backup describe $BACKUP_NAME --output json | jq -r '.status.itemsBackedUp')
    ERRORS=$(velero backup describe $BACKUP_NAME --output json | jq -r '.status.errors')
    
    echo "Resources backed up: $RESOURCES"
    echo "Errors: $ERRORS"
    
    # Test restore to temporary namespace
    RESTORE_NAME="validation-restore-$(date +%s)"
    TEMP_NAMESPACE="velero-validation-$(date +%s)"
    
    kubectl create namespace $TEMP_NAMESPACE
    
    velero restore create $RESTORE_NAME \
        --from-backup $BACKUP_NAME \
        --namespace-mappings default:$TEMP_NAMESPACE \
        --wait
    
    # Check restore status
    RESTORE_STATUS=$(velero restore describe $RESTORE_NAME --output json | jq -r '.status.phase')
    
    if [[ "$RESTORE_STATUS" == "Completed" ]]; then
        echo "✅ Backup validation successful - restore completed"
        
        # Cleanup
        kubectl delete namespace $TEMP_NAMESPACE
        velero restore delete $RESTORE_NAME --confirm
        
        exit 0
    else
        echo "❌ Backup validation failed - restore did not complete"
        echo "Restore status: $RESTORE_STATUS"
        exit 1
    fi
else
    echo "❌ Backup validation failed - backup status: $STATUS"
    exit 1
fi
```

### Comprehensive Disaster Recovery Plan
```bash
#!/bin/bash
# disaster-recovery-plan.sh

BACKUP_NAME="$1"
TARGET_CLUSTER="$2"

if [[ -z "$BACKUP_NAME" || -z "$TARGET_CLUSTER" ]]; then
    echo "Usage: $0 <backup-name> <target-cluster-context>"
    exit 1
fi

echo "Executing disaster recovery plan..."
echo "Source backup: $BACKUP_NAME"
echo "Target cluster: $TARGET_CLUSTER"

# Switch to target cluster
kubectl config use-context $TARGET_CLUSTER

# Verify Velero is installed on target cluster
if ! kubectl get namespace velero > /dev/null 2>&1; then
    echo "Installing Velero on target cluster..."
    velero install \
        --provider aws \
        --plugins velero/velero-plugin-for-aws:v1.8.0 \
        --bucket disaster-recovery-backups \
        --backup-location-config region=us-east-1 \
        --secret-file ./dr-credentials
    
    # Wait for Velero to be ready
    kubectl wait --for=condition=Available deployment/velero -n velero --timeout=300s
fi

# Verify backup exists
if ! velero backup get $BACKUP_NAME > /dev/null 2>&1; then
    echo "❌ Backup $BACKUP_NAME not found in target cluster"
    exit 1
fi

echo "Creating comprehensive restore..."

# Phase 1: Restore core infrastructure
velero restore create dr-infrastructure \
    --from-backup $BACKUP_NAME \
    --include-resources persistentvolumes,persistentvolumeclaims,storageclass \
    --wait

# Phase 2: Restore application configurations
velero restore create dr-configs \
    --from-backup $BACKUP_NAME \
    --include-resources configmaps,secrets \
    --exclude-namespaces kube-system,kube-public,velero \
    --wait

# Phase 3: Restore applications
velero restore create dr-applications \
    --from-backup $BACKUP_NAME \
    --exclude-resources persistentvolumes,persistentvolumeclaims,storageclass,configmaps,secrets \
    --exclude-namespaces kube-system,kube-public,velero \
    --wait

# Verify restoration
echo "Verifying disaster recovery..."

# Check all restores completed
for restore in dr-infrastructure dr-configs dr-applications; do
    status=$(velero restore describe $restore --output json | jq -r '.status.phase')
    if [[ "$status" != "Completed" ]]; then
        echo "❌ Restore $restore failed with status: $status"
        velero restore logs $restore
        exit 1
    fi
done

# Check pod status
TOTAL_PODS=$(kubectl get pods --all-namespaces --no-headers | grep -v "kube-system\|kube-public\|velero" | wc -l)
READY_PODS=$(kubectl get pods --all-namespaces --no-headers | grep -v "kube-system\|kube-public\|velero" | grep "Running\|Completed" | wc -l)

echo "Pod status: $READY_PODS/$TOTAL_PODS ready"

if [[ $READY_PODS -eq $TOTAL_PODS ]]; then
    echo "✅ Disaster recovery completed successfully!"
    echo "All applications restored and running on target cluster"
else
    echo "⚠️ Some pods are not ready yet. Monitor the cluster for full recovery."
    kubectl get pods --all-namespaces | grep -v "Running\|Completed"
fi

echo "Disaster recovery summary:"
echo "- Infrastructure restore: ✅"
echo "- Configuration restore: ✅"
echo "- Application restore: ✅"
echo "- Target cluster: $TARGET_CLUSTER"
```

### Backup Monitoring and Alerting
```bash
#!/bin/bash
# backup-monitoring.sh

SLACK_WEBHOOK="$1"
MONITORING_INTERVAL=3600  # 1 hour

if [[ -z "$SLACK_WEBHOOK" ]]; then
    echo "Usage: $0 <slack-webhook-url>"
    exit 1
fi

echo "Starting Velero backup monitoring..."

while true; do
    echo "Running backup health check at $(date)"
    
    # Check for failed backups in last 24 hours
    FAILED_BACKUPS=$(velero backup get --output json | jq -r '.items[] | select(.status.phase == "Failed" and (.metadata.creationTimestamp | strptime("%Y-%m-%dT%H:%M:%SZ") | mktime) > (now - 86400)) | .metadata.name')
    
    if [[ -n "$FAILED_BACKUPS" ]]; then
        MESSAGE="🚨 Velero backup failures detected:\n$FAILED_BACKUPS"
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"$MESSAGE\"}" \
            $SLACK_WEBHOOK
    fi
    
    # Check for old backups that haven't been cleaned up
    OLD_BACKUPS=$(velero backup get --output json | jq -r '.items[] | select(.status.phase == "Completed" and (.metadata.creationTimestamp | strptime("%Y-%m-%dT%H:%M:%SZ") | mktime) < (now - 2592000)) | .metadata.name')
    
    if [[ -n "$OLD_BACKUPS" ]]; then
        OLD_COUNT=$(echo "$OLD_BACKUPS" | wc -l)
        MESSAGE="🧹 Found $OLD_COUNT backups older than 30 days that may need cleanup"
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"$MESSAGE\"}" \
            $SLACK_WEBHOOK
    fi
    
    # Check schedule health
    SCHEDULES=$(velero schedule get --output json | jq -r '.items[].metadata.name')
    for schedule in $SCHEDULES; do
        LAST_BACKUP=$(velero schedule describe $schedule --output json | jq -r '.status.lastBackup // empty')
        if [[ -z "$LAST_BACKUP" ]]; then
            MESSAGE="⚠️ Schedule '$schedule' has no recent backups"
            curl -X POST -H 'Content-type: application/json' \
                --data "{\"text\":\"$MESSAGE\"}" \
                $SLACK_WEBHOOK
        fi
    done
    
    # Check Velero pod health
    VELERO_PODS=$(kubectl get pods -n velero --no-headers | grep -v Running | wc -l)
    if [[ $VELERO_PODS -gt 0 ]]; then
        MESSAGE="🚨 $VELERO_PODS Velero pods are not running properly"
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"$MESSAGE\"}" \
            $SLACK_WEBHOOK
    fi
    
    echo "Health check complete. Next check in $MONITORING_INTERVAL seconds."
    sleep $MONITORING_INTERVAL
done
```

## Integration Patterns

### CI/CD Pipeline Integration
```yaml
# .github/workflows/backup-on-deploy.yml
name: Backup on Deployment
on:
  workflow_run:
    workflows: ["Deploy to Production"]
    types:
      - completed

jobs:
  backup:
    runs-on: ubuntu-latest
    if: ${{ github.event.workflow_run.conclusion == 'success' }}
    steps:
    - name: Setup kubectl
      uses: azure/setup-kubectl@v1
      
    - name: Configure cluster access
      run: |
        echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
        
    - name: Install Velero CLI
      run: |
        wget https://github.com/vmware-tanzu/velero/releases/latest/download/velero-linux-amd64.tar.gz
        tar -xzf velero-linux-amd64.tar.gz
        sudo mv velero-*/velero /usr/local/bin/
        
    - name: Create Post-Deployment Backup
      run: |
        BACKUP_NAME="post-deploy-$(date +%Y%m%d-%H%M%S)"
        velero backup create $BACKUP_NAME \
          --include-namespaces production \
          --labels deployment=post-deploy,ci=github-actions \
          --wait
          
        # Verify backup completed
        STATUS=$(velero backup describe $BACKUP_NAME --output json | jq -r '.status.phase')
        if [[ "$STATUS" != "Completed" ]]; then
          echo "Backup failed with status: $STATUS"
          exit 1
        fi
        
        echo "✅ Post-deployment backup created: $BACKUP_NAME"
```

### Terraform Integration
```hcl
# terraform/velero-setup.tf
resource "kubernetes_namespace" "velero" {
  metadata {
    name = "velero"
  }
}

resource "helm_release" "velero" {
  name       = "velero"
  repository = "https://vmware-tanzu.github.io/helm-charts/"
  chart      = "velero"
  namespace  = kubernetes_namespace.velero.metadata[0].name

  values = [
    yamlencode({
      configuration = {
        provider = "aws"
        backupStorageLocation = {
          bucket = var.backup_bucket
          config = {
            region = var.aws_region
          }
        }
        volumeSnapshotLocation = {
          config = {
            region = var.aws_region
          }
        }
      }
      credentials = {
        useSecret = true
        secretContents = {
          cloud = base64encode(var.velero_credentials)
        }
      }
      deployRestic = true
      schedules = {
        daily = {
          schedule = "0 2 * * *"
          template = {
            ttl = "720h"
            includedNamespaces = ["production", "staging"]
          }
        }
      }
    })
  ]

  depends_on = [kubernetes_namespace.velero]
}

# Create backup storage bucket
resource "aws_s3_bucket" "velero_backups" {
  bucket = var.backup_bucket
}

resource "aws_s3_bucket_versioning" "velero_backups" {
  bucket = aws_s3_bucket.velero_backups.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "velero_backups" {
  bucket = aws_s3_bucket.velero_backups.id

  rule {
    id     = "backup_lifecycle"
    status = "Enabled"

    expiration {
      days = 90
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}
```

### ArgoCD Integration
```yaml
# argocd/velero-application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: velero
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://vmware-tanzu.github.io/helm-charts/
    chart: velero
    targetRevision: 5.1.4
    helm:
      values: |
        configuration:
          provider: aws
          backupStorageLocation:
            bucket: production-velero-backups
            config:
              region: us-west-2
        credentials:
          useSecret: true
        deployRestic: true
        schedules:
          daily-backup:
            schedule: "0 3 * * *"
            template:
              ttl: "720h"
              includedNamespaces:
              - production
              - staging
  destination:
    server: https://kubernetes.default.svc
    namespace: velero
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

## Best Practices

### Backup Strategy
- **Backup Frequency**: Daily backups for production, weekly for development
- **Retention Policy**: Keep 30 days of daily backups, 12 weeks of weekly backups
- **Multi-Location**: Store backups in multiple geographic locations
- **Encryption**: Use encryption at rest and in transit for sensitive data

### Scheduling Optimization
- **Off-Peak Hours**: Schedule backups during low-activity periods
- **Staggered Schedules**: Avoid overlapping backup windows
- **Resource Limits**: Set appropriate resource limits for backup pods
- **Monitoring**: Implement comprehensive backup monitoring and alerting

### Storage Management
- **Lifecycle Policies**: Implement automated cleanup of old backups
- **Cost Optimization**: Use appropriate storage classes for different retention periods
- **Compression**: Enable compression to reduce storage costs
- **Deduplication**: Use storage solutions that support deduplication

### Security Considerations
- **RBAC**: Implement least-privilege access for Velero operations
- **Secrets Management**: Securely manage cloud provider credentials
- **Network Policies**: Restrict network access to backup components
- **Audit Logging**: Enable comprehensive audit logging for backup operations

## Error Handling

### Common Issues
```bash
# Backup stuck in progress
velero backup describe stuck-backup
velero backup logs stuck-backup
# Solution: Check resource constraints and network connectivity

# Restore failing
velero restore describe failed-restore
velero restore logs failed-restore
# Solution: Check destination cluster compatibility and resource availability

# Storage location unavailable
velero backup-location get
kubectl describe backupstoragelocation -n velero
# Solution: Verify cloud credentials and network connectivity

# Plugin errors
kubectl logs deployment/velero -n velero
# Solution: Check plugin compatibility and configuration
```

### Troubleshooting
- **Permission Issues**: Verify cloud provider IAM permissions and Kubernetes RBAC
- **Network Connectivity**: Check network policies and firewall rules
- **Resource Constraints**: Ensure sufficient CPU, memory, and storage for operations
- **Version Compatibility**: Verify Velero version compatibility with Kubernetes cluster

Velero provides enterprise-grade backup and disaster recovery capabilities for Kubernetes environments, enabling organizations to protect their workloads with automated, reliable, and scalable backup solutions across multiple cloud providers and storage backends.