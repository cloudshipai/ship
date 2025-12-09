# Nuclei - Fast Vulnerability Scanner

Nuclei is a fast, customizable vulnerability scanner powered by a global security community. It uses simple YAML-based templates (called nuclei-templates) to detect vulnerabilities, misconfigurations, and security issues.

## Overview

- **Type**: Web application vulnerability scanner
- **Docker Image**: `projectdiscovery/nuclei:latest`
- **Category**: Security
- **Requires Variables**: No

## Features

- Fast and highly configurable
- Community-driven templates library with 10,000+ templates
- Supports multiple protocols: HTTP, DNS, TCP, SSL, WebSocket, etc.
- YAML-based DSL for writing custom templates
- Low false positive rate
- JSON/JSONL/Markdown/SARIF output formats

## MCP Tools

### nuclei_scan_url
Scan a single URL for vulnerabilities.

**Parameters:**
- `url` (required): Target URL to scan
- `severity`: Filter by severity (info, low, medium, high, critical)

**Example:**
```json
{
  "name": "nuclei_scan_url",
  "arguments": {
    "url": "https://example.com",
    "severity": "high,critical"
  }
}
```

### nuclei_scan_with_tags
Scan using specific vulnerability tags.

**Parameters:**
- `url` (required): Target URL to scan
- `tags` (required): Comma-separated tags (e.g., cve, sqli, xss, tech)

**Example:**
```json
{
  "name": "nuclei_scan_with_tags",
  "arguments": {
    "url": "https://example.com",
    "tags": "cve,sqli,xss"
  }
}
```

### nuclei_scan_with_template
Scan using specific Nuclei template(s).

**Parameters:**
- `url` (required): Target URL to scan
- `template`: Template path or name

### nuclei_update_templates
Update Nuclei vulnerability templates to the latest version.

### nuclei_validate_template
Validate a custom Nuclei template file.

**Parameters:**
- `template_path` (required): Path to template file to validate

### nuclei_generate_report
Generate a vulnerability scan report in various formats.

**Parameters:**
- `url` (required): Target URL to scan
- `report_type`: Report format (json, markdown, sarif)

## Usage Examples

### Basic URL Scan
```bash
ship mcp nuclei
# Then call: nuclei_scan_url with url="https://target.com"
```

### Scan for CVEs Only
```bash
# Use nuclei_scan_with_tags with tags="cve"
```

### Technology Detection
```bash
# Use nuclei_scan_with_tags with tags="tech"
```

## Template Tags Reference

Common template tags:
- `cve` - CVE vulnerabilities
- `sqli` - SQL injection
- `xss` - Cross-site scripting
- `lfi` - Local file inclusion
- `rfi` - Remote file inclusion
- `rce` - Remote code execution
- `ssrf` - Server-side request forgery
- `tech` - Technology detection
- `misconfig` - Misconfigurations
- `exposure` - Information exposure
- `default-login` - Default credentials

## Output Format

Nuclei returns results in JSONL format with details including:
- Template ID and name
- Severity level
- Matched URL
- Request/response details
- CVE information (if applicable)

## Environment Variables

Override the default Docker image:
```bash
export SHIP_IMAGE_TAG_NUCLEI=projectdiscovery/nuclei:v3.6.0
```

## References

- [Nuclei GitHub](https://github.com/projectdiscovery/nuclei)
- [Nuclei Templates](https://github.com/projectdiscovery/nuclei-templates)
- [Nuclei Documentation](https://docs.projectdiscovery.io/tools/nuclei)
