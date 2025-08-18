# Check SSL Cert MCP Tool

Check SSL Cert is a shell script that validates SSL/TLS certificates for expiration, validity, and other security properties.

## Description

Check SSL Cert performs comprehensive SSL/TLS certificate validation including:
- Certificate expiration checking
- Chain validation
- Protocol support verification
- Cipher suite analysis
- Certificate transparency log checking

## MCP Functions

### `check_ssl_cert_validate`
Validate SSL certificate using real check_ssl_cert script.

**Parameters:**
- `host` (required): Hostname to check
- `port`: Port number (default: 443)
- `warning_days`: Days before expiration to warn (default: 20)
- `critical_days`: Days before expiration for critical alert (default: 15)

**CLI Command:** `check_ssl_cert -H <host> [-p <port>] [-w <warning_days>] [-c <critical_days>]`

### `check_ssl_cert_validate_file`
Validate SSL certificate from file using real check_ssl_cert script.

**Parameters:**
- `cert_file` (required): Path to certificate file
- `warning_days`: Days before expiration to warn
- `critical_days`: Days before expiration for critical alert

**CLI Command:** `check_ssl_cert -f <cert_file> [-w <warning_days>] [-c <critical_days>]`

### `check_ssl_cert_check_chain`
Validate certificate chain using real check_ssl_cert script.

**Parameters:**
- `host` (required): Hostname to check
- `port`: Port number (default: 443)
- `root_cert`: Path to root certificate file

**CLI Command:** `check_ssl_cert -H <host> [-p <port>] [-r <root_cert>] --ignore-exp`

### `check_ssl_cert_check_protocols`
Check supported SSL/TLS protocols using real check_ssl_cert script.

**Parameters:**
- `host` (required): Hostname to check
- `port`: Port number (default: 443)
- `protocol`: Specific protocol to check (SSL2, SSL3, TLS1, TLS1_1, TLS1_2, TLS1_3)

**CLI Command:** `check_ssl_cert -H <host> [-p <port>] [-P <protocol>]`

### `check_ssl_cert_comprehensive`
Comprehensive SSL certificate check using real check_ssl_cert script.

**Parameters:**
- `host` (required): Hostname to check
- `port`: Port number (default: 443)
- `warning_days`: Days before expiration to warn
- `critical_days`: Days before expiration for critical alert
- `check_ocsp`: Enable OCSP checking
- `check_sct`: Check for certificate transparency

**CLI Command:** `check_ssl_cert -H <host> [-p <port>] [-w <warning_days>] [-c <critical_days>] [--ocsp] [--sct]`

### `check_ssl_cert_version`
Get check_ssl_cert version information.

**CLI Command:** `check_ssl_cert --version`

## Common Use Cases

1. **Certificate Expiration Monitoring**: Check certificates before they expire
2. **Chain Validation**: Verify complete certificate chains
3. **Protocol Compliance**: Ensure only secure protocols are enabled
4. **OCSP Validation**: Check certificate revocation status
5. **Certificate Transparency**: Verify CT log inclusion

## Integration with Ship CLI

All Check SSL Cert tools are integrated with the Ship CLI's MCP server, allowing AI assistants to:
- Monitor SSL/TLS certificate health
- Validate certificate chains
- Check protocol compliance
- Verify security configurations

The tools use containerized execution via Dagger for consistent, isolated certificate validation.