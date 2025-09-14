package interfaces

import (
	"context"
	"time"
)

// ResourceProvider defines resource-specific operations
type ResourceProvider interface {
	// Identity
	GetID() string
	GetName() string
	GetType() ResourceType
	GetVendor() VendorType
	GetRegion() string
	GetTags() map[string]string
	
	// Cost information
	GetCurrentCost(ctx context.Context) (*CostData, error)
	GetCostHistory(ctx context.Context, period TimePeriod) ([]CostRecord, error)
	
	// Performance metrics
	GetUtilization(ctx context.Context) (*UtilizationMetrics, error)
	GetMetricsHistory(ctx context.Context, period TimePeriod) ([]UtilizationRecord, error)
	
	// Optimization capabilities
	GetRecommendations(ctx context.Context) ([]Recommendation, error)
	CanOptimize() bool
	
	// Resource-specific data
	GetSpecifications() ResourceSpecs
	GetConfiguration() map[string]interface{}
}

// ComputeResource extends ResourceProvider for compute instances
type ComputeResource interface {
	ResourceProvider
	
	// Compute-specific methods
	GetInstanceType() string
	GetCPUCount() int
	GetMemoryGB() float64
	GetStorageGB() float64
	
	// Operations
	CanResize() bool
	GetResizeOptions() []ResizeOption
	GetRightsizingRecommendation() (*RightsizingRecommendation, error)
	
	// State management
	GetState() InstanceState
	GetUptime() time.Duration
}

// StorageResource extends ResourceProvider for storage resources
type StorageResource interface {
	ResourceProvider
	
	// Storage-specific methods
	GetSizeGB() float64
	GetStorageType() string // gp2, gp3, io1, standard, etc.
	GetIOPS() int
	GetThroughput() float64
	
	// Optimization
	GetStorageClassOptions() []StorageClassOption
	GetStorageOptimizationRecommendation() (*StorageRecommendation, error)
	
	// Usage patterns
	GetAccessPattern() AccessPattern
	IsAttached() bool
	GetLastAccessTime() *time.Time
}

// DatabaseResource extends ResourceProvider for database instances
type DatabaseResource interface {
	ResourceProvider
	
	// Database-specific methods
	GetEngine() string
	GetEngineVersion() string
	GetInstanceClass() string
	GetAllocatedStorage() int
	GetStorageType() string
	GetIOPS() int
	
	// Operations
	GetUpgradeOptions() []DatabaseUpgradeOption
	GetPerformanceInsights() *DatabasePerformanceMetrics
	
	// Backup and maintenance
	GetBackupRetention() int
	GetMaintenanceWindow() string
}

// KubernetesResource extends ResourceProvider for K8s resources
type KubernetesResource interface {
	ResourceProvider
	
	// K8s-specific methods
	GetNamespace() string
	GetCluster() string
	GetLabels() map[string]string
	GetAnnotations() map[string]string
	
	// Cost allocation
	GetAllocatedCost() (*AllocationCost, error)
	GetIdleCost() (*IdleCost, error)
	
	// Resource requests/limits
	GetResourceRequests() ResourceRequirements
	GetResourceLimits() ResourceRequirements
	GetActualUsage() ResourceRequirements
	
	// Optimization
	GetRequestsOptimizationRecommendation() (*ResourceOptimizationRecommendation, error)
}

// UtilizationRecord represents utilization data for a specific time point
type UtilizationRecord struct {
	Timestamp time.Time          `json:"timestamp"`
	Metrics   UtilizationMetrics `json:"metrics"`
}

// ResizeOption represents a compute instance resize option
type ResizeOption struct {
	InstanceType    string  `json:"instance_type"`
	CPUCount       int     `json:"cpu_count"`
	MemoryGB       float64 `json:"memory_gb"`
	HourlyRate     float64 `json:"hourly_rate_usd"`
	ExpectedSavings float64 `json:"expected_savings_monthly_usd"`
	PerformanceRisk string  `json:"performance_risk"` // Low, Medium, High
	MigrationEffort string  `json:"migration_effort"` // Low, Medium, High
}

// RightsizingRecommendation provides rightsizing guidance
type RightsizingRecommendation struct {
	CurrentType      string                 `json:"current_type"`
	RecommendedType  string                 `json:"recommended_type"`
	Confidence       float64                `json:"confidence"` // 0.0 to 1.0
	MonthlySavings   float64                `json:"monthly_savings_usd"`
	Justification    string                 `json:"justification"`
	RiskFactors      []string               `json:"risk_factors"`
	UtilizationData  *UtilizationMetrics    `json:"utilization_data"`
}

// InstanceState represents compute instance state
type InstanceState string

const (
	InstanceStateRunning   InstanceState = "running"
	InstanceStateStopped   InstanceState = "stopped"
	InstanceStatePending   InstanceState = "pending"
	InstanceStateTerminated InstanceState = "terminated"
)

// StorageClassOption represents a storage class optimization option
type StorageClassOption struct {
	ClassName       string  `json:"class_name"`
	Description     string  `json:"description"`
	CostPerGB       float64 `json:"cost_per_gb_monthly"`
	AccessTime      string  `json:"access_time"`
	MinimumDuration string  `json:"minimum_duration"`
	ExpectedSavings float64 `json:"expected_savings_monthly_usd"`
}

// StorageRecommendation provides storage optimization guidance
type StorageRecommendation struct {
	CurrentClass       string             `json:"current_class"`
	RecommendedClass   string             `json:"recommended_class"`
	Confidence         float64            `json:"confidence"`
	MonthlySavings     float64            `json:"monthly_savings_usd"`
	Justification      string             `json:"justification"`
	AccessPattern      AccessPattern      `json:"access_pattern"`
}

// AccessPattern describes how storage is accessed
type AccessPattern struct {
	Frequency       string    `json:"frequency"` // frequent, infrequent, archive
	LastAccessed    time.Time `json:"last_accessed"`
	AccessCount30d  int       `json:"access_count_30d"`
	ReadWriteRatio  float64   `json:"read_write_ratio"`
}

// DatabaseUpgradeOption represents database upgrade possibilities
type DatabaseUpgradeOption struct {
	EngineVersion   string  `json:"engine_version"`
	InstanceClass   string  `json:"instance_class"`
	StorageType     string  `json:"storage_type"`
	ExpectedSavings float64 `json:"expected_savings_monthly_usd"`
	PerformanceGain string  `json:"performance_gain"`
	UpgradeEffort   string  `json:"upgrade_effort"`
}

// DatabasePerformanceMetrics contains database-specific metrics
type DatabasePerformanceMetrics struct {
	ConnectionsUsed    int     `json:"connections_used"`
	ConnectionsMax     int     `json:"connections_max"`
	QueryLatencyP95    float64 `json:"query_latency_p95_ms"`
	ThroughputQPS      float64 `json:"throughput_qps"`
	StorageUtilization float64 `json:"storage_utilization"`
	IOPSUtilization    float64 `json:"iops_utilization"`
}

// AllocationCost represents Kubernetes cost allocation
type AllocationCost struct {
	CPU         float64 `json:"cpu_cost_usd"`
	Memory      float64 `json:"memory_cost_usd"`
	Storage     float64 `json:"storage_cost_usd"`
	Network     float64 `json:"network_cost_usd"`
	Total       float64 `json:"total_cost_usd"`
	Period      string  `json:"period"`
	Efficiency  float64 `json:"efficiency"`  // How much of requested resources are used
}

// IdleCost represents unallocated Kubernetes costs
type IdleCost struct {
	CPU     float64 `json:"cpu_idle_cost_usd"`
	Memory  float64 `json:"memory_idle_cost_usd"`
	Storage float64 `json:"storage_idle_cost_usd"`
	Network float64 `json:"network_idle_cost_usd"`
	Total   float64 `json:"total_idle_cost_usd"`
	Period  string  `json:"period"`
}

// ResourceRequirements represents K8s resource requests/limits
type ResourceRequirements struct {
	CPU    string `json:"cpu"`    // "500m", "1", "2"
	Memory string `json:"memory"` // "512Mi", "1Gi", "2Gi"
	Storage string `json:"storage,omitempty"`
}

// ResourceOptimizationRecommendation provides K8s resource optimization guidance
type ResourceOptimizationRecommendation struct {
	CurrentRequests    ResourceRequirements `json:"current_requests"`
	RecommendedRequests ResourceRequirements `json:"recommended_requests"`
	MonthlySavings     float64              `json:"monthly_savings_usd"`
	Confidence         float64              `json:"confidence"`
	Justification      string               `json:"justification"`
}