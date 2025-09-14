package interfaces

import (
	"context"
	"time"
)

// VendorType represents supported cloud vendors
type VendorType string

const (
	VendorAWS        VendorType = "aws"
	VendorGCP        VendorType = "gcp"  
	VendorAzure      VendorType = "azure"
	VendorKubernetes VendorType = "kubernetes"
)

// VendorProvider is the core interface that all cloud vendor implementations must satisfy
type VendorProvider interface {
	// Identity
	Name() string
	Type() VendorType
	
	// Resource discovery with vendor-specific filtering
	DiscoverResources(ctx context.Context, opts DiscoveryOptions) ([]Resource, error)
	
	// Cost data with vendor-specific APIs
	GetCostData(ctx context.Context, opts CostOptions) ([]CostRecord, error)
	
	// Recommendations using vendor's native recommendation engines
	GetRecommendations(ctx context.Context, opts RecommendationOptions) ([]Recommendation, error)
	
	// Vendor-specific capabilities
	GetCapabilities() VendorCapabilities
	
	// Authentication - follows OpenOps credential patterns
	SetCredentials(creds interface{}) error
	
	// Vendor-specific configuration validation
	ValidateConfig(config map[string]interface{}) error
}

// VendorCapabilities defines what operations a vendor supports
type VendorCapabilities struct {
	SupportsRecommendations bool     `json:"supports_recommendations"`
	SupportsCostForecasting bool     `json:"supports_cost_forecasting"`
	SupportsRightsizing     bool     `json:"supports_rightsizing"`
	ResourceTypes          []string  `json:"resource_types"`
	Regions               []string  `json:"regions"`
}

// DiscoveryOptions configures resource discovery operations
type DiscoveryOptions struct {
	Region        string            `json:"region,omitempty"`
	ResourceTypes []string          `json:"resource_types,omitempty"`
	ARNs          []string          `json:"arns,omitempty"`        // Optional ARN filtering like OpenOps
	Tags          map[string]string `json:"tags,omitempty"`
	AccountID     string            `json:"account_id,omitempty"`  // Multi-account support
	AccountIDs    []string          `json:"account_ids,omitempty"` // Multiple account IDs
}

// CostOptions configures cost data retrieval
type CostOptions struct {
	TimeWindow    string            `json:"time_window"`     // "7d", "30d", "90d"  
	Granularity   string            `json:"granularity"`     // "daily", "monthly"
	GroupBy       []string          `json:"group_by"`        // ["service", "region", "account"]
	Filters       map[string]string `json:"filters,omitempty"`
	Currency      string            `json:"currency"`        // "USD", "EUR", etc.
}

// RecommendationOptions configures recommendation generation
type RecommendationOptions struct {
	FindingTypes []string `json:"finding_types"` // ["rightsizing", "reserved_instances", "spot_instances"]
	Regions      []string `json:"regions,omitempty"`     // Optional region filter
	ARNs         []string `json:"arns,omitempty"`        // Optional ARN filter (OpenOps pattern)
	AccountIDs   []string `json:"account_ids,omitempty"` // Multi-account support
}

// Resource represents a cloud resource with standardized properties
type Resource struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Type              ResourceType           `json:"type"`
	Vendor            VendorType             `json:"vendor"`
	Region            string                 `json:"region"`
	AccountID         string                 `json:"account_id,omitempty"`
	Tags              map[string]string      `json:"tags"`
	Cost              *CostData              `json:"cost,omitempty"`
	Utilization       *UtilizationMetrics    `json:"utilization,omitempty"`
	Specifications    ResourceSpecs          `json:"specifications"`
	Configuration     map[string]interface{} `json:"configuration"`
	LastUpdated       time.Time              `json:"last_updated"`
}

// ResourceType categorizes different types of cloud resources
type ResourceType string

const (
	ResourceTypeCompute      ResourceType = "compute"
	ResourceTypeStorage      ResourceType = "storage"
	ResourceTypeDatabase     ResourceType = "database"
	ResourceTypeServerless   ResourceType = "serverless"
	ResourceTypeNetwork      ResourceType = "network"
	ResourceTypeCache        ResourceType = "cache"
	ResourceTypeKubernetes   ResourceType = "kubernetes"
)

// CostRecord represents cost data for a specific time period
type CostRecord struct {
	ResourceID    string            `json:"resource_id"`
	Period        TimePeriod        `json:"period"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Service       string            `json:"service,omitempty"`
	Region        string            `json:"region,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Recommendation represents an optimization opportunity
type Recommendation struct {
	ID                     string            `json:"id"`
	ResourceID             string            `json:"resource_id"`
	ResourceARN            string            `json:"resource_arn,omitempty"`
	ResourceType           ResourceType      `json:"resource_type"`
	Provider               VendorType        `json:"provider"`
	RecommendationType     string            `json:"recommendation_type"`
	Title                  string            `json:"title"`
	Description            string            `json:"description"`
	EstimatedSavings       EstimatedSavings  `json:"estimated_savings"`
	PerformanceRisk        string            `json:"performance_risk"`    // Low, Medium, High
	MigrationEffort        string            `json:"migration_effort"`    // Low, Medium, High
	Confidence             float64           `json:"confidence"`          // 0.0 to 1.0
	Options                []RecommendationOption `json:"options"`
	Justification          string            `json:"justification"`
	RiskFactors            []string          `json:"risk_factors"`
	CreatedAt              time.Time         `json:"created_at"`
}

// EstimatedSavings represents cost savings potential
type EstimatedSavings struct {
	Currency string  `json:"currency"`
	Hourly   float64 `json:"hourly,omitempty"`
	Daily    float64 `json:"daily,omitempty"`
	Monthly  float64 `json:"monthly"`
	Yearly   float64 `json:"yearly,omitempty"`
}

// RecommendationOption represents a specific optimization action
type RecommendationOption struct {
	OptionType      string  `json:"option_type"`
	Description     string  `json:"description"`
	ExpectedSavings float64 `json:"expected_savings"`
	Parameters      map[string]interface{} `json:"parameters"`
}

// CostData represents current cost information for a resource
type CostData struct {
	HourlyRate   float64   `json:"hourly_rate_usd"`
	DailyRate    float64   `json:"daily_rate_usd"`
	MonthlyRate  float64   `json:"monthly_rate_usd"`
	Currency     string    `json:"currency"`
	LastUpdated  time.Time `json:"last_updated"`
}

// UtilizationMetrics represents resource utilization over time
type UtilizationMetrics struct {
	CPU         *MetricValue `json:"cpu,omitempty"`
	Memory      *MetricValue `json:"memory,omitempty"`
	Storage     *MetricValue `json:"storage,omitempty"`
	Network     *MetricValue `json:"network,omitempty"`
	IOPS        *MetricValue `json:"iops,omitempty"`
	Period      TimePeriod   `json:"period"`
	LastUpdated time.Time    `json:"last_updated"`
}

// MetricValue represents statistical data for a metric
type MetricValue struct {
	Average    float64            `json:"average"`
	Peak       float64            `json:"peak"`
	Minimum    float64            `json:"minimum"`
	Unit       string             `json:"unit"`
	Percentile map[int]float64    `json:"percentiles,omitempty"` // P50, P95, P99
}

// ResourceSpecs contains resource-specific specifications
type ResourceSpecs struct {
	InstanceType string  `json:"instance_type,omitempty"`
	CPUCount     int     `json:"cpu_count,omitempty"`
	MemoryGB     float64 `json:"memory_gb,omitempty"`
	StorageGB    float64 `json:"storage_gb,omitempty"`
	StorageType  string  `json:"storage_type,omitempty"`
}

// TimePeriod represents a time range
type TimePeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}