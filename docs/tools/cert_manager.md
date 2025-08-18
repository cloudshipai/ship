# Cert-Manager - Kubernetes Certificate Management

Cert-manager is a Kubernetes add-on that automates the management and issuance of TLS certificates from various issuing sources. This tool provides MCP functions for managing cert-manager installations and certificates.

## Overview

Cert-manager runs within your Kubernetes cluster as a set of deployment resources. It automatically issues TLS certificates and ensures they are valid and up to date. These MCP functions wrap the standard kubectl and cmctl commands for cert-manager operations.

## Available MCP Functions

### 1. `cert_manager_install`
**Description**: Install cert-manager using kubectl apply

**Parameters**:
- `version` (optional): Cert-manager version to install (default: v1.18.2)
- `dry_run` (optional): Perform dry run without actually installing

**Example Usage**:
```bash
# Install latest supported version
cert_manager_install()

# Install specific version
cert_manager_install(version="v1.17.0")

# Dry run installation
cert_manager_install(version="v1.18.2", dry_run=true)
```

### 2. `cert_manager_check_installation`
**Description**: Check cert-manager installation status

**Parameters**:
- `namespace` (optional): Namespace to check (default: cert-manager)

**Example Usage**:
```bash
# Check default cert-manager namespace
cert_manager_check_installation()

# Check specific namespace
cert_manager_check_installation(namespace="cert-manager")
```

### 3. `cert_manager_create_certificate_request`
**Description**: Create a CertificateRequest using cmctl

**Parameters**:
- `name` (required): Name of the CertificateRequest
- `from_certificate_file` (optional): Path to certificate file to create request from
- `fetch_certificate` (optional): Fetch the certificate once issued
- `timeout` (optional): Timeout for the operation (e.g., 20m)

**Example Usage**:
```bash
# Create basic certificate request
cert_manager_create_certificate_request(name="my-cert-request")

# Create from certificate file with fetch
cert_manager_create_certificate_request(
  name="my-cert-request",
  from_certificate_file="/path/to/certificate.yaml",
  fetch_certificate=true,
  timeout="20m"
)
```

### 4. `cert_manager_list_certificates`
**Description**: List certificates using kubectl

**Parameters**:
- `namespace` (optional): Kubernetes namespace to list certificates from
- `all_namespaces` (optional): List certificates from all namespaces

**Example Usage**:
```bash
# List certificates in current namespace
cert_manager_list_certificates()

# List certificates in specific namespace
cert_manager_list_certificates(namespace="production")

# List certificates in all namespaces
cert_manager_list_certificates(all_namespaces=true)
```

### 5. `cert_manager_renew_certificate`
**Description**: Mark certificate for manual renewal using cmctl

**Parameters**:
- `cert_name` (required): Name of the certificate to renew (not used if all=true)
- `namespace` (optional): Kubernetes namespace
- `all` (optional): Renew all certificates in namespace

**Example Usage**:
```bash
# Renew specific certificate
cert_manager_renew_certificate(cert_name="my-certificate")

# Renew specific certificate in namespace
cert_manager_renew_certificate(
  cert_name="my-certificate",
  namespace="production"
)

# Renew all certificates in namespace
cert_manager_renew_certificate(all=true, namespace="production")
```

### 6. `cert_manager_status`
**Description**: Get status of certificate using cmctl

**Parameters**:
- `certificate_name` (required): Name of the certificate to check
- `namespace` (optional): Kubernetes namespace

**Example Usage**:
```bash
# Check certificate status
cert_manager_status(certificate_name="my-certificate")

# Check certificate status in specific namespace
cert_manager_status(
  certificate_name="my-certificate",
  namespace="production"
)
```

### 7. `cert_manager_get_version`
**Description**: Get cmctl version information

**Parameters**: None

**Example Usage**:
```bash
cert_manager_get_version()
```

## Real CLI Capabilities

All MCP functions are based on actual kubectl and cmctl commands:

### Installation
```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.18.2/cert-manager.yaml

# Check installation
kubectl get pods -n cert-manager
```

### Certificate Management with cmctl
```bash
# Create certificate request
cmctl create certificaterequest my-cr --from-certificate-file my-certificate.yaml --fetch-certificate --timeout 20m

# Renew certificates
cmctl renew my-certificate
cmctl renew --all -n my-namespace

# Check status
cmctl status certificate my-certificate -n my-namespace

# Get version
cmctl version
```

### Certificate Operations with kubectl
```bash
# List certificates
kubectl get certificates
kubectl get certificates -n my-namespace
kubectl get certificates --all-namespaces

# Describe certificate
kubectl describe certificate my-certificate -n my-namespace

# List all cert-manager resources
kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges --all-namespaces
```

## Prerequisites

### Kubernetes Cluster
- Running Kubernetes cluster (v1.22+)
- kubectl configured and connected
- Sufficient permissions to install and manage cert-manager

### cmctl Installation
```bash
# Install cmctl on Linux
curl -L -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/latest/download/cmctl-linux-amd64.tar.gz
tar xzf cmctl.tar.gz
sudo mv cmctl /usr/local/bin

# Install cmctl on macOS
brew install cmctl

# Verify installation
cmctl version
```

## Cert-Manager Components

### Core Components
- **cert-manager**: Main controller that manages certificates
- **cert-manager-cainjector**: Injects CA data into webhooks and APIServices
- **cert-manager-webhook**: Provides validation and mutation webhooks

### Custom Resources
- **Certificate**: Defines desired TLS certificate
- **Issuer/ClusterIssuer**: Defines certificate authorities
- **CertificateRequest**: Request for certificate issuance
- **Order/Challenge**: ACME protocol resources

## Common Workflows

### Basic Installation Workflow
1. **Install cert-manager**: Use `cert_manager_install`
2. **Verify installation**: Use `cert_manager_check_installation`
3. **Check version**: Use `cert_manager_get_version`

### Certificate Management Workflow
1. **Create Issuer**: Apply Issuer/ClusterIssuer YAML
2. **Request certificate**: Apply Certificate YAML or use cmctl
3. **Monitor status**: Use `cert_manager_status`
4. **List certificates**: Use `cert_manager_list_certificates`
5. **Renew when needed**: Use `cert_manager_renew_certificate`

### Troubleshooting Workflow
1. **Check installation**: Verify all pods are running
2. **Check certificate status**: Use cmctl status
3. **Review logs**: Check cert-manager pod logs
4. **Validate resources**: Ensure Issuers and Certificates are correct

## Issuer Types

### Self-Signed Issuer
```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
```

### Let's Encrypt Issuer
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: user@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

### CA Issuer
```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
spec:
  ca:
    secretName: ca-key-pair
```

## Certificate Examples

### Basic Certificate
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
spec:
  secretName: example-com-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
  - example.com
  - www.example.com
```

### Wildcard Certificate
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: wildcard-example-com
spec:
  secretName: wildcard-example-com-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
  - "*.example.com"
```

## Best Practices

### Installation
- **Use specific versions**: Pin cert-manager to tested versions
- **Monitor resources**: Set up monitoring for cert-manager components
- **Backup configuration**: Backup Issuer and Certificate resources
- **Test in staging**: Always test certificate workflows in non-production

### Security
- **Least privilege**: Use RBAC to limit cert-manager permissions
- **Secure secrets**: Protect private key secrets
- **Network policies**: Apply network policies to cert-manager pods
- **Regular updates**: Keep cert-manager updated for security patches

### Operations
- **Monitor expiry**: Set up alerts for certificate expiration
- **Automate renewal**: Let cert-manager handle automatic renewal
- **Document issuers**: Maintain documentation of issuer configurations
- **Test recovery**: Regularly test certificate recovery procedures

## Troubleshooting

### Common Issues

1. **Installation Fails**
   - Check Kubernetes version compatibility
   - Verify cluster has sufficient resources
   - Check network connectivity for webhook

2. **Certificate Not Issued**
   - Verify Issuer is ready: `kubectl describe issuer`
   - Check Certificate events: `kubectl describe certificate`
   - Review cert-manager logs

3. **ACME Challenges Fail**
   - Verify DNS/HTTP01 challenge configuration
   - Check ingress controller setup
   - Ensure domain is publicly accessible

4. **cmctl Commands Fail**
   - Verify cmctl installation and version
   - Check Kubernetes context and permissions
   - Ensure cert-manager is installed and running

### Useful Commands
```bash
# Check all cert-manager resources
kubectl get certificates,certificaterequests,issuers,clusterissuers --all-namespaces

# Check cert-manager logs
kubectl logs -n cert-manager -l app=cert-manager

# Describe certificate for troubleshooting
kubectl describe certificate <cert-name> -n <namespace>

# Force certificate renewal
cmctl renew <cert-name> -n <namespace>
```

## Integration with Ship CLI

These MCP functions integrate with Ship CLI's containerized execution:
- kubectl and cmctl commands are executed through Ship CLI's Dagger engine
- Kubernetes access requires proper kubeconfig configuration
- Certificate management operations respect Kubernetes RBAC

## References

- **Cert-Manager Documentation**: https://cert-manager.io/docs/
- **Installation Guide**: https://cert-manager.io/docs/installation/kubectl/
- **cmctl Reference**: https://cert-manager.io/docs/reference/cmctl/
- **GitHub Repository**: https://github.com/cert-manager/cert-manager
- **Kubectl Plugin**: https://cert-manager.io/v1.0-docs/usage/kubectl-plugin/