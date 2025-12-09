# OpenVAS - Greenbone Vulnerability Management Integration

OpenVAS (Open Vulnerability Assessment Scanner) is a full-featured vulnerability scanner that is part of the Greenbone Vulnerability Management (GVM) framework. This integration provides tools to interact with a running GVM instance via the GMP (Greenbone Management Protocol).

## Overview

- **Type**: Network vulnerability scanner
- **Docker Image**: `registry.community.greenbone.net/community/gvm-tools:latest`
- **Category**: Security
- **Requires Variables**: Yes (GVM credentials)

## Prerequisites

OpenVAS/GVM requires a complete stack to be running, which includes:
- PostgreSQL database
- Redis server
- GVM daemon (gvmd)
- OpenVAS Scanner
- Greenbone Security Assistant (GSA) web interface

### Quick Setup with Docker Compose

Greenbone provides a docker-compose setup for running the complete GVM stack:

```bash
# Download the docker-compose file
curl -f -O https://greenbone.github.io/docs/latest/_static/docker-compose.yml

# Start the GVM stack
docker compose up -d

# Wait for services to initialize (can take several minutes on first run)
docker compose logs -f
```

The web interface will be available at `http://localhost:9392` (default credentials: admin/admin).

## MCP Tools

### openvas_get_system_status
Get the status of a running GVM instance.

**Parameters:**
- `socket_path`: Path to GVM socket (default: /run/gvmd/gvmd.sock)
- `username` (required): GVM username
- `password` (required): GVM password

### openvas_list_scans
List all vulnerability scans from the GVM instance.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password

### openvas_list_targets
List all scan targets from the GVM instance.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password

### openvas_create_target
Create a new scan target.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `target_name` (required): Name for the target
- `hosts` (required): Comma-separated list of hosts/IPs to scan
- `port_list_id`: ID of port list to use (optional)

**Example:**
```json
{
  "name": "openvas_create_target",
  "arguments": {
    "username": "admin",
    "password": "admin",
    "target_name": "Web Servers",
    "hosts": "192.168.1.10,192.168.1.11"
  }
}
```

### openvas_create_scan
Create a new vulnerability scan task.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `scan_name` (required): Name for the scan task
- `target_id` (required): ID of the target to scan
- `config_id`: Scan configuration ID (default: Full and fast)

### openvas_start_scan
Start an existing scan task.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `task_id` (required): ID of the scan task to start

### openvas_get_scan_status
Get the status of a scan task.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `task_id` (required): ID of the scan task

### openvas_get_scan_results
Get the results of a completed scan.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `report_id` (required): ID of the scan report
- `min_severity`: Minimum severity to include (0.0-10.0)

### openvas_list_scan_configs
List available scan configurations.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password

### openvas_stop_scan
Stop a running scan task.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `task_id` (required): ID of the scan task to stop

### openvas_delete_scan
Delete a scan task.

**Parameters:**
- `socket_path`: Path to GVM socket
- `username` (required): GVM username
- `password` (required): GVM password
- `task_id` (required): ID of the scan task to delete

## Workflow Example

1. **Create a target:**
```json
{
  "name": "openvas_create_target",
  "arguments": {
    "username": "admin",
    "password": "admin",
    "target_name": "Production Servers",
    "hosts": "192.168.1.0/24"
  }
}
```

2. **Create a scan task:** (use the target_id from step 1)
```json
{
  "name": "openvas_create_scan",
  "arguments": {
    "username": "admin",
    "password": "admin",
    "scan_name": "Weekly Production Scan",
    "target_id": "<target-id-from-step-1>"
  }
}
```

3. **Start the scan:**
```json
{
  "name": "openvas_start_scan",
  "arguments": {
    "username": "admin",
    "password": "admin",
    "task_id": "<task-id-from-step-2>"
  }
}
```

4. **Check scan status:**
```json
{
  "name": "openvas_get_scan_status",
  "arguments": {
    "username": "admin",
    "password": "admin",
    "task_id": "<task-id>"
  }
}
```

5. **Get results:** (once scan is complete)
```json
{
  "name": "openvas_get_scan_results",
  "arguments": {
    "username": "admin",
    "password": "admin",
    "report_id": "<report-id-from-status>",
    "min_severity": "7.0"
  }
}
```

## Common Scan Configuration IDs

- `daba56c8-73ec-11df-a475-002264764cea` - Full and fast (default)
- `698f691e-7489-11df-9d8c-002264764cea` - Full and fast ultimate
- `708f25c4-7489-11df-8094-002264764cea` - Full and very deep
- `74db13d6-7489-11df-91b9-002264764cea` - Full and very deep ultimate

## Environment Variables

Override the default Docker image:
```bash
export SHIP_IMAGE_TAG_OPENVAS=registry.community.greenbone.net/community/gvm-tools:stable
```

## References

- [Greenbone OpenVAS](https://github.com/greenbone/openvas-scanner)
- [GVM Documentation](https://greenbone.github.io/docs/)
- [GVM Docker Compose Setup](https://greenbone.github.io/docs/latest/22.4/container/)
- [GMP Protocol Documentation](https://docs.greenbone.net/API/GMP/gmp-22.04.html)
