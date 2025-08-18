# Falco

Runtime security monitoring for containers and Kubernetes.

## Description

Falco is the de facto Kubernetes threat detection engine. It detects unexpected application behavior and alerts on threats at runtime by observing Linux kernel calls and other data sources. Falco can detect and alert on any behavior that involves making Linux system calls.

## MCP Tools

### Runtime Monitoring
- **`falco_start_monitoring`** - Start Falco runtime security monitoring
- **`falco_dry_run`** - Run Falco in dry-run mode without processing events

### Rule Management
- **`falco_validate_rules`** - Validate Falco rules syntax
- **`falco_list_rules`** - List all loaded Falco rules
- **`falco_describe_rule`** - Show description of a specific Falco rule

### System Information
- **`falco_list_fields`** - List supported fields for Falco rules
- **`falco_get_version`** - Get Falco version information

## Real CLI Commands Used

- `falco` - Start Falco monitoring with default configuration
- `falco -c /path/to/config.yaml` - Start with custom configuration
- `falco -r /path/to/rules.yaml` - Start with custom rules
- `falco --dry-run` - Run in dry-run mode without processing events
- `falco -V /path/to/rules.yaml` - Validate rules file syntax
- `falco -L` - List all loaded rules
- `falco -l <rule_name>` - Show description of specific rule
- `falco --list [source]` - List supported fields for rules
- `falco --version` - Show version information

## Security Detection Capabilities

### Container Security
- Unexpected process execution in containers
- Privilege escalation attempts
- Container escape attempts
- Unauthorized file access in containers

### Kubernetes Security
- Pod creation with suspicious settings
- Unauthorized service account usage
- Suspicious network connections
- Kubernetes API server abuse

### System Security
- Suspicious system call patterns
- File integrity monitoring
- Network anomaly detection
- Process behavior analysis

### Custom Rules
- YAML-based rule definitions
- Flexible condition matching
- Custom output formatting
- Integration with external systems

## Use Cases

### Runtime Threat Detection
- Real-time security monitoring
- Behavioral anomaly detection
- Incident response automation
- Compliance monitoring

### Cloud-Native Security
- Kubernetes workload protection
- Container runtime security
- Multi-cloud security monitoring
- Serverless security

### DevSecOps Integration
- CI/CD security gates
- Security testing automation
- Compliance verification
- Security metrics collection

## Rule Examples

### Detect Shell Access in Container
```yaml
- rule: Shell in Container
  desc: Notice shell activity within a container
  condition: >
    spawned_process and container and
    shell_procs and proc.tty != 0
  output: Shell spawned in container (user=%user.name container_id=%container.id shell=%proc.name)
  priority: WARNING
```

### Detect Suspicious Network Activity
```yaml
- rule: Unexpected outbound connection destination
  desc: Detect connections to well-known mining services
  condition: >
    outbound and
    fd.sip in (mining_pool_addresses)
  output: Outbound connection to mining pool (command=%proc.cmdline connection=%fd.name)
  priority: CRITICAL
```

## Event Sources

- **Linux syscalls** - System call monitoring via kernel modules or eBPF
- **Kubernetes audit logs** - API server activity monitoring
- **Cloud provider logs** - AWS CloudTrail, GCP Audit Logs, etc.
- **Custom plugins** - Extensible data source integration

## Output Formats

- Standard output with formatted alerts
- JSON output for integration with SIEM systems
- gRPC output for real-time streaming
- File output for logging and archival

## Integration

Works as a Kubernetes DaemonSet, container, or standalone binary. Integrates with alerting systems like Slack, PagerDuty, and security orchestration platforms.