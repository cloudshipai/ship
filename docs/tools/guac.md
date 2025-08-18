# GUAC

Graph for Understanding Artifact Composition - Supply chain metadata aggregation and analysis.

## Description

GUAC (Graph for Understanding Artifact Composition) is an open-source project that aggregates software security metadata into a high-fidelity graph database. It normalizes and maps relationships between software security entities, enabling advanced querying for supply chain transparency, vulnerability analysis, and compliance verification. GUAC occupies the "aggregation and synthesis" layer of software supply chain transparency.

## MCP Tools

### Data Collection & Ingestion
- **`guac_collect`** - Run GUAC collector to gather supply chain metadata
- **`guac_certify`** - Run GUAC certifier to add metadata attestations
- **`guac_ingest`** - Run GUAC ingestor service connected to NATS and GraphQL

### Analysis & Querying
- **`guac_query`** - Run GUAC canned queries for analysis
- **`guac_start_graphql`** - Start GUAC GraphQL server

### Service Operations
- **`guac_collector_subscriber`** - Run GUAC collector-subscriber service

## Real CLI Commands Used

### Core Commands
- `guacone collector <source>` - Collect metadata from various sources
- `guacone certify <type>` - Add certifications and attestations
- `guacone query <type>` - Run predefined queries for analysis
- `guacgql` - Start GraphQL server with backend configuration
- `guacingest` - Run ingestor service for data processing
- `guaccsub` - Run collector-subscriber service

### Command Parameters
- `--poll` - Enable continuous polling mode for certifiers
- `--path <path>` - Specify file path for file-based collectors
- `--gql-backend <type>` - Set GraphQL backend (keyvalue, postgresql, redis, tikv)
- `--gql-listen-port <port>` - Configure GraphQL server port
- `--nats-addr <address>` - NATS server address for messaging
- `--gql-addr <address>` - GraphQL server address for connections

## Supported Input Sources

### Metadata Collection Sources
- **files** - Local files and directories containing SBOMs, attestations
- **deps_dev** - Dependencies from deps.dev API
- **osv** - Open Source Vulnerabilities database
- **scorecard** - OpenSSF Scorecard security metrics

### Certification Types
- **osv** - OSV vulnerability certifications
- **scorecard** - Security scorecard certifications
- **vulns** - General vulnerability certifications

### Query Types
- **vulnerabilities** - Query vulnerability information
- **dependencies** - Analyze dependency relationships
- **packages** - Package metadata and relationships

## Architecture Components

### Backend Options
- **keyvalue** - In-memory backend (recommended for beginners)
- **postgresql** - Persistent PostgreSQL backend
- **redis** - Redis-based backend (experimental)
- **tikv** - TiKV distributed backend (experimental)

### Data Processing Pipeline
1. **Collection**: Gather metadata from various sources
2. **Ingestion**: Process and normalize collected data
3. **Certification**: Add attestations and verifications
4. **Storage**: Store in graph database backend
5. **Querying**: Analyze relationships and patterns

### Integration Architecture
- **NATS Messaging**: Asynchronous communication between components
- **GraphQL API**: Unified query interface for applications
- **Microservices**: Modular, scalable component architecture
- **Event-driven**: Real-time processing of supply chain events

## Use Cases

### Supply Chain Transparency
- **Dependency Mapping**: Visualize complete dependency graphs
- **Vulnerability Tracking**: Monitor vulnerabilities across supply chains
- **Risk Assessment**: Identify high-risk components and paths
- **Impact Analysis**: Understand blast radius of security issues

### Compliance and Governance
- **Audit Trails**: Maintain comprehensive artifact lineage
- **Policy Enforcement**: Validate compliance against organizational policies
- **Reporting**: Generate compliance reports for various frameworks
- **Attestation Verification**: Validate build and security attestations

### Security Operations
- **Threat Intelligence**: Correlate vulnerabilities with deployed systems
- **Incident Response**: Quickly identify affected components
- **Risk Prioritization**: Focus remediation efforts on critical paths
- **Security Metrics**: Track security posture improvements

### Development Workflow
- **Build Verification**: Validate build integrity and provenance
- **Dependency Analysis**: Understand transitive dependency risks
- **Security Gates**: Implement security checkpoints in CI/CD
- **Continuous Monitoring**: Real-time supply chain security monitoring

## Configuration Examples

### Basic Data Collection
```bash
# Collect from local SBOM files
guacone collector files --path ./sboms/

# Collect OpenSSF Scorecard data
guacone collector scorecard

# Collect OSV vulnerability data
guacone collector osv
```

### Certification and Attestation
```bash
# Add OSV vulnerability certifications
guacone certify osv

# Add scorecard certifications with polling
guacone certify scorecard --poll

# Add general vulnerability certifications
guacone certify vulns
```

### Querying and Analysis
```bash
# Query vulnerabilities for a package
guacone query vulnerabilities --subject "npm:lodash"

# Analyze dependencies
guacone query dependencies --subject "docker.io/library/ubuntu:latest"

# Query package metadata
guacone query packages --subject "maven:org.apache:log4j"
```

## Service Deployment

### GraphQL Server Setup
```bash
# Start with in-memory backend
guacgql --gql-backend keyvalue --gql-listen-port 8080

# Start with PostgreSQL backend
guacgql --gql-backend postgresql --gql-listen-port 8080
```

### Distributed Architecture
```bash
# Start ingestor service
guacingest --nats-addr nats://localhost:4222 --gql-addr http://localhost:8080

# Start collector-subscriber
guaccsub --nats-addr nats://localhost:4222 --source files

# Collect with specific configuration
guaccollect files --nats-addr nats://localhost:4222 --gql-addr http://localhost:8080
```

## Integration Patterns

### Docker Compose Deployment
```yaml
version: '3.8'
services:
  graphql:
    image: guac-graphql
    command: guacgql --gql-backend keyvalue
    ports:
      - "8080:8080"
  
  collector:
    image: guac-collector
    command: guacone collector files --path /data
    volumes:
      - ./sboms:/data
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: guac-graphql
spec:
  template:
    spec:
      containers:
      - name: graphql
        image: guac:latest
        command: ["guacgql"]
        args: ["--gql-backend", "postgresql"]
        ports:
        - containerPort: 8080
```

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Collect SBOM Data
  run: |
    guacone collector files --path ./build-artifacts/
    guacone certify osv
```

## Data Model and Standards

### Supported Formats
- **CycloneDX** - SBOM standard for dependency tracking
- **SPDX** - Software Package Data Exchange format
- **SLSA** - Supply-chain Levels for Software Artifacts
- **In-toto** - Build attestation framework
- **OpenSSF Scorecard** - Security metrics and ratings
- **OSV** - Open Source Vulnerabilities format
- **OpenVEX** - Vulnerability Exploitability eXchange

### Graph Relationships
- **Package Dependencies** - Direct and transitive relationships
- **Vulnerability Associations** - CVEs and security advisories
- **Build Provenance** - Source to artifact relationships
- **Attestations** - Security and quality certifications
- **Scorecard Metrics** - Security practice evaluations

## Best Practices

### Data Collection Strategy
- **Comprehensive Coverage**: Collect from multiple authoritative sources
- **Regular Updates**: Implement continuous collection and certification
- **Source Verification**: Validate data source authenticity
- **Incremental Processing**: Handle large datasets efficiently

### Query Optimization
- **Indexing Strategy**: Optimize backend for common query patterns
- **Caching**: Implement appropriate caching for frequent queries
- **Batch Processing**: Group related queries for efficiency
- **Result Pagination**: Handle large result sets appropriately

### Security Considerations
- **Access Control**: Implement authentication and authorization
- **Data Encryption**: Protect sensitive supply chain information
- **Audit Logging**: Track all data access and modifications
- **Network Security**: Secure inter-service communications

### Operational Excellence
- **Monitoring**: Track system health and performance metrics
- **Backup Strategy**: Implement regular data backup procedures
- **Scaling**: Plan for horizontal scaling of services
- **Documentation**: Maintain comprehensive operational documentation

GUAC provides a powerful foundation for understanding and securing software supply chains through comprehensive metadata aggregation and intelligent relationship mapping.