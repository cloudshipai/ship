# Station Data Schemas

This document describes the data schemas that Station will receive from the Ship FinOps tools. This file is not committed to the repository.

## Overview

Ship FinOps tools send data to Station (CloudshipAI proxy) via HTTP POST requests to `/api/v1/finops/*` endpoints. All data follows standardized schemas for consistent processing.

## Core Data Types

### 1. Resource Schema

Resources represent discovered cloud infrastructure:

```json
{
  "id": "i-1234567890abcdef0",
  "name": "web-server-prod-01",
  "type": "compute",
  "vendor": "aws",
  "region": "us-east-1",
  "account_id": "123456789012",
  "tags": {
    "Environment": "production",
    "Application": "web-frontend",
    "Team": "platform"
  },
  "cost": {
    "hourly_rate_usd": 0.063,
    "daily_rate_usd": 1.512,
    "monthly_rate_usd": 45.60,
    "currency": "USD",
    "last_updated": "2025-09-14T06:58:44Z"
  },
  "utilization": {
    "cpu": {
      "average": 15.5,
      "peak": 45.2,
      "minimum": 2.1,
      "unit": "percent",
      "percentiles": {
        "50": 12.3,
        "95": 38.7,
        "99": 44.1
      }
    },
    "memory": {
      "average": 32.1,
      "peak": 67.8,
      "minimum": 18.9,
      "unit": "percent"
    },
    "period": {
      "start": "2025-08-14T06:58:44Z",
      "end": "2025-09-14T06:58:44Z"
    },
    "last_updated": "2025-09-14T06:58:44Z"
  },
  "specifications": {
    "instance_type": "t3.medium",
    "cpu_count": 2,
    "memory_gb": 4,
    "storage_gb": 20,
    "storage_type": "gp3"
  },
  "configuration": {
    "state": "running",
    "launch_time": "2025-01-15T10:30:00Z",
    "availability_zone": "us-east-1a",
    "security_groups": ["sg-12345", "sg-67890"]
  },
  "last_updated": "2025-09-14T06:58:44Z"
}
```

### 2. Opportunity Schema

Opportunities represent cost optimization recommendations:

```json
{
  "id": "opp-aws-rightsize-i1234567890abcdef0-20250914",
  "resource_id": "i-1234567890abcdef0",
  "resource_arn": "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
  "opportunity_type": "rightsizing",
  "vendor": "aws",
  "title": "Rightsize EC2 instance from t3.medium to t3.small",
  "description": "This EC2 instance shows consistently low CPU and memory utilization over the past 30 days. CPU averages 15.5% with peaks of 45.2%. Memory utilization averages 32.1%. Rightsizing to t3.small would maintain performance while reducing costs.",
  "estimated_savings_monthly_usd": 22.80,
  "currency": "USD",
  "complexity": "S",
  "risk_level": "Low",
  "implementation_effort": "Low",
  "tags": {
    "Environment": "production",
    "Application": "web-frontend"
  },
  "region": "us-east-1",
  "service": "EC2",
  "account_id": "123456789012",
  "recommendation_details": {
    "current_configuration": {
      "instance_type": "t3.medium",
      "vcpu": 2,
      "memory_gb": 4,
      "monthly_cost_usd": 45.60
    },
    "recommended_configuration": {
      "instance_type": "t3.small", 
      "vcpu": 2,
      "memory_gb": 2,
      "monthly_cost_usd": 22.80
    },
    "implementation_steps": [
      "1. Create AMI snapshot of current instance",
      "2. Schedule maintenance window",
      "3. Stop instance and change instance type",
      "4. Start instance and verify functionality",
      "5. Monitor performance for 48 hours"
    ],
    "prerequisites": [
      "Verify application memory requirements",
      "Coordinate with application team",
      "Backup current instance state"
    ],
    "rollback_plan": [
      "Stop instance",
      "Change back to t3.medium",
      "Restart and verify"
    ],
    "testing_guidance": "Monitor CPU and memory metrics for 48 hours. Verify application response times remain within SLA."
  },
  "created_at": "2025-09-14T06:58:44Z",
  "valid_until": "2025-10-14T06:58:44Z",
  "tenant_id": "tenant-123",
  "status": "active"
}
```

### 3. Cost Record Schema

Cost records represent time-series cost data:

```json
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
    "Environment": "production",
    "CostCenter": "engineering"
  },
  "metadata": {
    "instance_type": "t3.medium",
    "usage_type": "BoxUsage:t3.medium",
    "operation": "RunInstances",
    "billing_mode": "on-demand"
  }
}
```

### 4. Insight Schema

Insights represent cost anomalies and patterns:

```json
{
  "id": "insight-cost-anomaly-20250914-001",
  "type": "cost_anomaly",
  "severity": "medium",
  "title": "Unusual EC2 spend increase in us-east-1",
  "description": "EC2 costs in us-east-1 increased by 45% compared to last month average. Investigation shows 3 new t3.large instances launched without proper tagging.",
  "category": "cost_management",
  "resources_affected": [
    "i-1234567890abcdef0",
    "i-0987654321fedcba0",
    "i-abcdef1234567890"
  ],
  "cost_impact_usd": 234.56,
  "time_range": {
    "start": "2025-09-01T00:00:00Z",
    "end": "2025-09-14T00:00:00Z"
  },
  "vendor": "aws",
  "region": "us-east-1",
  "service": "EC2",
  "tags": {
    "alert_type": "cost_anomaly",
    "auto_generated": "true"
  },
  "metadata": {
    "threshold_exceeded": 1.45,
    "baseline_cost": 520.30,
    "current_cost": 754.86,
    "detection_method": "statistical_anomaly"
  },
  "created_at": "2025-09-14T06:58:44Z",
  "tenant_id": "tenant-123",
  "status": "active"
}
```

## API Request Schemas

### 1. Resource Inventory Request

**Endpoint:** `POST /api/v1/finops/inventory`

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
    // Array of Resource objects (see Resource Schema above)
  ]
}
```

### 2. Opportunities Request

**Endpoint:** `POST /api/v1/finops/opportunities`

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
    // Array of Opportunity objects (see Opportunity Schema above)
  ]
}
```

### 3. Cost Data Request

**Endpoint:** `POST /api/v1/finops/costs`

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
    "granularity": "daily",
    "group_by": ["service", "region"]
  },
  "cost_records": [
    // Array of CostRecord objects (see Cost Record Schema above)
  ]
}
```

### 4. Insights Request

**Endpoint:** `POST /api/v1/finops/insights`

```json
{
  "tenant_id": "tenant-123",
  "data_type": "insights",
  "timestamp": "2025-09-14T06:58:44Z",
  "source": "ship-finops-analyze",
  "version": "v1.0.0", 
  "metadata": {
    "analysis_type": "cost_anomaly_detection",
    "detection_sensitivity": "medium"
  },
  "insights": [
    // Array of Insight objects (see Insight Schema above)
  ]
}
```

## API Response Schema

All Station endpoints return standardized responses:

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

## Data Validation Rules

### Required Fields
- All `*_id` fields must be non-empty strings
- `vendor` must be one of: `aws`, `gcp`, `azure`, `kubernetes`
- `estimated_savings_monthly_usd` must be >= 0
- `opportunity_type` must be one of: `rightsizing`, `idle-resource`, `commitment`, `storage-optimization`, `spot-instances`

### Field Formats
- All timestamps in RFC 3339 format (`2025-09-14T06:58:44Z`)
- Currency amounts as decimal numbers (not strings)
- Percentages as decimal (0.0 to 100.0)
- ARNs must follow AWS ARN format when provided

### Size Limits
- Maximum 1000 resources per inventory request
- Maximum 500 opportunities per request
- Maximum 5000 cost records per request
- Maximum 100 insights per request

## gRPC Alternative

For high-volume scenarios, Station also supports gRPC using the schema in `/internal/tools/finops/grpc/finops.proto`.

## Authentication

Since Ship and Station run on the same server:
- No authentication headers required
- All requests are local HTTP calls to `http://localhost:8080/api/v1/*`
- `Content-Type: application/json` header only

## Rate Limits

- 100 requests per minute per tenant
- 10MB maximum request size
- 30 second request timeout

This schema documentation ensures Station can properly ingest, validate, and process all FinOps data from Ship tools.