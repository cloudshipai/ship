# FinOps Tools Package

A comprehensive FinOps (Financial Operations) toolkit for multi-cloud cost optimization, resource discovery, and recommendations. Built with pluggable vendor providers and CloudshipAI Station integration.

## Station API Endpoints (for Station Team)

**Station should implement these endpoints on port 8585:**

### Required Endpoints for FinOps Integration

#### 1. Resource Inventory Endpoint
```
POST http://localhost:8585/api/v1/finops/inventory
Content-Type: application/json
```

**Request Body:**
```json
{
  "tenant_id": "tenant-123",
  "data_type": "resource_inventory", 
  "timestamp": "2025-09-14T06:58:44Z",
  "source": "ship-finops-discover",
  "version": "v1.0.0",
  "metadata": {
    "provider": "aws",
    "region": "us-east-1",
    "discovery_options": {
      "resource_types": ["compute", "storage"],
      "tags": {"Environment": "production"}
    }
  },
  "resources": [
    {
      "id": "i-1234567890abcdef0",
      "name": "web-server-prod-01",
      "type": "compute",
      "vendor": "aws",
      "region": "us-east-1",
      "account_id": "123456789012",
      "tags": {
        "Environment": "production",
        "Application": "web-frontend"
      },
      "cost": {
        "hourly_rate_usd": 0.063,
        "monthly_rate_usd": 45.60,
        "currency": "USD"
      },
      "utilization": {
        "cpu": {"average": 15.5, "peak": 45.2},
        "memory": {"average": 32.1, "peak": 67.8}
      },
      "specifications": {
        "instance_type": "t3.medium",
        "cpu_count": 2,
        "memory_gb": 4
      },
      "configuration": {
        "state": "running",
        "availability_zone": "us-east-1a"
      }
    }
  ]
}
```

#### 2. Opportunities/Recommendations Endpoint
```
POST http://localhost:8585/api/v1/finops/opportunities
Content-Type: application/json
```

**Request Body:**
```json
{
  "tenant_id": "tenant-123",
  "data_type": "opportunities",
  "timestamp": "2025-09-14T06:58:44Z", 
  "source": "ship-finops-recommend",
  "version": "v1.0.0",
  "metadata": {
    "provider": "aws",
    "finding_types": ["rightsizing", "reserved_instances"],
    "min_savings_filter": 25.0
  },
  "opportunities": [
    {
      "id": "opp-aws-rightsize-i1234567890abcdef0-20250914",
      "resource_id": "i-1234567890abcdef0",
      "resource_arn": "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
      "opportunity_type": "rightsizing",
      "vendor": "aws",
      "title": "Rightsize EC2 instance from t3.medium to t3.small",
      "description": "This EC2 instance shows consistently low CPU and memory utilization",
      "estimated_savings_monthly_usd": 22.80,
      "currency": "USD",
      "complexity": "S",
      "risk_level": "Low",
      "implementation_effort": "Low",
      "recommendation_details": {
        "current_configuration": {
          "instance_type": "t3.medium",
          "monthly_cost_usd": 45.60
        },
        "recommended_configuration": {
          "instance_type": "t3.small", 
          "monthly_cost_usd": 22.80
        },
        "implementation_steps": [
          "1. Create AMI snapshot of current instance",
          "2. Schedule maintenance window",
          "3. Stop instance and change instance type"
        ]
      },
      "created_at": "2025-09-14T06:58:44Z"
    }
  ]
}
```

#### 3. Cost Data Endpoint
```
POST http://localhost:8585/api/v1/finops/costs
Content-Type: application/json
```

**Request Body:**
```json
{
  "tenant_id": "tenant-123",
  "data_type": "cost_data",
  "timestamp": "2025-09-14T06:58:44Z",
  "source": "ship-finops-analyze", 
  "version": "v1.0.0",
  "metadata": {
    "provider": "aws",
    "time_window": "30d",
    "granularity": "daily"
  },
  "cost_records": [
    {
      "resource_id": "i-1234567890abcdef0",
      "period": {
        "start": "2025-09-13T00:00:00Z",
        "end": "2025-09-14T00:00:00Z"
      },
      "amount": 1.512,
      "currency": "USD",
      "service": "EC2-Instance",
      "region": "us-east-1",
      "tags": {
        "Environment": "production"
      }
    }
  ]
}
```

#### 4. Insights Endpoint
```
POST http://localhost:8585/api/v1/finops/insights
Content-Type: application/json
```

**Request Body:**
```json
{
  "tenant_id": "tenant-123",
  "data_type": "insights",
  "timestamp": "2025-09-14T06:58:44Z",
  "source": "ship-finops-analyze",
  "version": "v1.0.0", 
  "metadata": {
    "analysis_type": "cost_anomaly_detection"
  },
  "insights": [
    {
      "id": "insight-cost-anomaly-20250914-001",
      "type": "cost_anomaly",
      "severity": "medium",
      "title": "Unusual EC2 spend increase in us-east-1",
      "description": "EC2 costs increased by 45% compared to last month",
      "cost_impact_usd": 234.56,
      "time_range": {
        "start": "2025-09-01T00:00:00Z",
        "end": "2025-09-14T00:00:00Z"
      },
      "vendor": "aws",
      "region": "us-east-1"
    }
  ]
}
```

### Expected Response Format

All Station endpoints should return:

```json
{
  "success": true,
  "request_id": "req-1234567890abcdef",
  "records_processed": 15,
  "errors": [],
  "warnings": [
    {
      "code": "MISSING_TAG",
      "message": "Resource i-abc123 missing Environment tag",
      "field": "tags.Environment"
    }
  ],
  "metadata": {
    "processing_time_ms": 145,
    "storage_location": "s3://cloudshipai-data/tenant-123/finops/",
    "next_sync_recommended": "2025-09-15T06:58:44Z"
  },
  "timestamp": "2025-09-14T06:58:44Z"
}
```

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Missing required field: provider",
    "details": {
      "field": "metadata.provider",
      "expected": "string",
      "received": "null"
    }
  },
  "request_id": "req-1234567890abcdef",
  "timestamp": "2025-09-14T06:58:44Z"
}
```

## Architecture Overview

The FinOps package follows a modular, vendor-agnostic architecture that enables easy extension to new cloud providers:

```
internal/tools/finops/
├── interfaces/           # Core abstractions and contracts
│   ├── vendor.go        # VendorProvider interface and data types
│   ├── lighthouse.go    # CloudshipAI Station integration
│   └── types.go         # Common data structures
├── providers/           # Vendor-specific implementations
│   ├── aws/            # Amazon Web Services provider
│   ├── gcp/            # Google Cloud Platform (future)
│   ├── azure/          # Microsoft Azure (future)
│   └── kubernetes/     # Kubernetes cost analysis (future)
├── mcp/                # MCP (Model Context Protocol) tools
├── cloudshipai/        # Station API client implementation
└── grpc/               # gRPC schema definitions for Station
```

## Core Interfaces

### VendorProvider Interface

All cloud providers implement the `VendorProvider` interface, ensuring consistent behavior across different clouds:

```go
type VendorProvider interface {
    // Basic provider information
    Name() string
    Type() VendorType
    GetCapabilities() VendorCapabilities
    
    // Configuration and authentication
    SetCredentials(creds interface{}) error
    ValidateConfig(config map[string]interface{}) error
    
    // Core FinOps operations
    DiscoverResources(ctx context.Context, opts DiscoveryOptions) ([]Resource, error)
    GetRecommendations(ctx context.Context, opts RecommendationOptions) ([]Recommendation, error)
    GetCostData(ctx context.Context, opts CostOptions) ([]CostRecord, error)
}
```

## Adding New Cloud Providers

The interfaces are designed to be generic enough for multi-cloud support. Here's how each service maps:

### 1. Google Cloud Platform (GCP) Implementation

#### Service Mappings
| GCP Service | ResourceType | Discovery API | Recommendations API |
|-------------|--------------|---------------|---------------------|
| Compute Engine | `compute` | `compute.instances.list` | `recommender.projects.locations.recommenders.recommendations.list` |
| Persistent Disks | `storage` | `compute.disks.list` | `recommender` (disk size/type) |
| Cloud SQL | `database` | `sqladmin.instances.list` | `recommender` (machine type) |
| Cloud Storage | `storage` | `storage.buckets.list` | `recommender` (lifecycle rules) |
| Cloud Functions | `serverless` | `cloudfunctions.projects.locations.functions.list` | Performance insights |

#### GCP Provider Structure
```go
type GCPProvider struct {
    computeClient     *compute.Service
    recommenderClient *recommender.Service
    billingClient     *billing.CloudBillingClient
    monitoringClient  *monitoring.Service
}

func (p *GCPProvider) DiscoverResources(ctx context.Context, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
    // Use compute.instances.list to get VM instances
    // Map to generic Resource with GCP-specific tags as labels
    // Extract utilization from Cloud Monitoring API
}

func (p *GCPProvider) GetRecommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
    // Use Recommender API to get rightsizing insights
    // Convert recommenderType: "google.compute.instance.MachineTypeRecommender" 
    // to RecommendationType: "RightSizeInstance"
}
```

### 2. Microsoft Azure Implementation

#### Service Mappings
| Azure Service | ResourceType | Discovery API | Recommendations API |
|---------------|--------------|---------------|---------------------|
| Virtual Machines | `compute` | `compute.VirtualMachinesClient.List` | `advisor.RecommendationsClient.List` |
| Managed Disks | `storage` | `compute.DisksClient.List` | Azure Advisor (cost recommendations) |
| SQL Database | `database` | `sql.DatabasesClient.ListByServer` | Azure Advisor (performance) |
| Storage Accounts | `storage` | `storage.AccountsClient.List` | Azure Advisor (cost/performance) |
| Functions | `serverless` | `web.AppsClient.List` | Application Insights |

#### Azure Provider Structure
```go
type AzureProvider struct {
    computeClient     compute.VirtualMachinesClient
    advisorClient     advisor.RecommendationsClient
    consumptionClient consumption.UsageDetailsClient
    monitorClient     monitor.MetricsClient
}

func (p *AzureProvider) GetRecommendations(ctx context.Context, opts interfaces.RecommendationOptions) ([]interfaces.Recommendation, error) {
    // Use Azure Advisor API to get cost recommendations
    // Map Category: "Cost" + Impact: "High" to our confidence levels
    // Convert advisor recommendation to generic format
}
```

### 3. Kubernetes Cost Analysis Implementation

#### Service Mappings
| K8s Resource | ResourceType | Discovery API | Cost Source |
|--------------|--------------|---------------|-------------|
| Pods | `compute` | `k8s.CoreV1().Pods().List` | OpenCost/KubeCost metrics |
| PersistentVolumes | `storage` | `k8s.CoreV1().PersistentVolumes().List` | Cloud provider billing |
| Services | `network` | `k8s.CoreV1().Services().List` | Load balancer costs |
| Nodes | `compute` | `k8s.CoreV1().Nodes().List` | Cloud instance costs |

```go
type KubernetesProvider struct {
    clientset         kubernetes.Interface
    metricsClient     metrics.Interface
    prometheusClient  promapi.Client  // For OpenCost integration
}

func (p *KubernetesProvider) DiscoverResources(ctx context.Context, opts interfaces.DiscoveryOptions) ([]interfaces.Resource, error) {
    // List pods with resource requests/limits and actual usage
    // Calculate cost per pod based on node costs and resource allocation
    // Map K8s labels to generic tags
}
```

## Generic Data Compatibility

The interfaces are designed to work across all cloud providers:

### Resource Discovery
- **AWS**: `DescribeInstances` → `Resource{Type: compute, Vendor: aws}`
- **GCP**: `instances.list` → `Resource{Type: compute, Vendor: gcp}`
- **Azure**: `VirtualMachines.List` → `Resource{Type: compute, Vendor: azure}`
- **K8s**: `Pods.List` → `Resource{Type: compute, Vendor: kubernetes}`

### Cost Optimization
- **AWS**: `GetEC2InstanceRecommendations` → `Recommendation{RecommendationType: "RightSizeInstance"}`
- **GCP**: `Recommender.MachineTypeRecommender` → `Recommendation{RecommendationType: "RightSizeInstance"}`
- **Azure**: `Advisor.CostRecommendation` → `Recommendation{RecommendationType: "RightSizeInstance"}`
- **K8s**: `Resource Quotas Analysis` → `Recommendation{RecommendationType: "RightSizeContainer"}`

### Cost Data
- **AWS**: `GetCostAndUsage` → `CostRecord{Service: "EC2-Instance"}`
- **GCP**: `CloudBilling` → `CostRecord{Service: "Compute Engine"}`
- **Azure**: `ConsumptionClient` → `CostRecord{Service: "Virtual Machines"}`
- **K8s**: `OpenCost metrics` → `CostRecord{Service: "Pod CPU/Memory"}`

## Station Schema Compliance

All providers normalize their data to Station-compatible schemas:

```go
// Generic Resource → Station Resource Schema
func convertToStationResource(resource interfaces.Resource) StationResource {
    return StationResource{
        ID:       resource.ID,
        Name:     resource.Name,
        Type:     string(resource.Type),           // "compute", "storage", etc.
        Vendor:   string(resource.Vendor),        // "aws", "gcp", "azure", "kubernetes"
        Region:   resource.Region,
        Tags:     resource.Tags,                  // All providers support tags/labels
        Cost: StationCost{
            HourlyRateUSD:  resource.Cost.HourlyRate,
            MonthlyRateUSD: resource.Cost.MonthlyRate,
            Currency:       "USD",               // Normalized across providers
        },
        // ... other fields mapped generically
    }
}
```

## Provider Registration

New providers are automatically integrated into all MCP tools:

```go
// In mcp/discover_tool.go
func NewFinOpsDiscoverTool(config interfaces.LighthouseConfig) (*FinOpsDiscoverTool, error) {
    providers := make(map[interfaces.VendorType]interfaces.VendorProvider)
    
    // AWS Provider
    if awsProvider, err := aws.NewProvider(); err == nil {
        providers[interfaces.VendorAWS] = awsProvider
    }
    
    // GCP Provider (future)
    if gcpProvider, err := gcp.NewProvider(); err == nil {
        providers[interfaces.VendorGCP] = gcpProvider
    }
    
    // Azure Provider (future)
    if azureProvider, err := azure.NewProvider(); err == nil {
        providers[interfaces.VendorAzure] = azureProvider
    }
    
    // Kubernetes Provider (future)
    if k8sProvider, err := kubernetes.NewProvider(); err == nil {
        providers[interfaces.VendorKubernetes] = k8sProvider
    }
    
    return &FinOpsDiscoverTool{providers: providers}, nil
}
```

## Usage Examples

### Multi-Provider Discovery
```bash
# AWS resources
ship finops discover --provider aws --region us-east-1

# Future: GCP resources
ship finops discover --provider gcp --project my-gcp-project --region us-central1

# Future: Azure resources  
ship finops discover --provider azure --subscription-id xxx --resource-group production

# Future: Kubernetes resources
ship finops discover --provider kubernetes --namespace production
```

### Cross-Provider Recommendations
```bash
# Get recommendations from all available providers
ship mcp finops --var AWS_PROFILE=prod --var GCP_PROJECT=my-project --var KUBECONFIG=~/.kube/config
```

The architecture ensures that **all providers send the same normalized data format to Station**, making it easy for CloudshipAI to provide consistent cost optimization insights across any cloud environment!