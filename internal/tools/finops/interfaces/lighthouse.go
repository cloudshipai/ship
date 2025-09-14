package interfaces

import (
	"context"
	"time"
)

// LighthouseReporter defines the interface for reporting FinOps data to Station/CloudshipAI
type LighthouseReporter interface {
	// Report resource inventory discovered by finops tools
	ReportResources(ctx context.Context, resources []Resource) error
	
	// Report cost optimization opportunities
	ReportOpportunities(ctx context.Context, opportunities []Opportunity) error
	
	// Report cost data and trends
	ReportCostData(ctx context.Context, costData []CostRecord) error
	
	// Report insights and anomalies
	ReportInsights(ctx context.Context, insights []Insight) error
	
	// Health check for lighthouse connectivity
	HealthCheck(ctx context.Context) error
}

// Opportunity represents a cost optimization opportunity for CloudshipAI
type Opportunity struct {
	ID                    string            `json:"id"`
	ResourceID            string            `json:"resource_id" validate:"required"`
	ResourceARN           string            `json:"resource_arn,omitempty"`
	Type                  OpportunityType   `json:"opportunity_type" validate:"required,oneof=rightsizing idle-resource commitment storage-optimization"`
	Vendor                VendorType        `json:"vendor" validate:"required,oneof=aws gcp azure kubernetes"`
	Title                 string            `json:"title"`
	Description           string            `json:"description"`
	EstimatedSavings      float64           `json:"estimated_savings_monthly_usd" validate:"required,min=0"`
	Currency              string            `json:"currency"`
	Complexity            string            `json:"complexity,omitempty" validate:"omitempty,oneof=XS S M L XL"`
	RiskLevel             string            `json:"risk_level,omitempty" validate:"omitempty,oneof=Low Medium High"`
	ImplementationEffort  string            `json:"implementation_effort,omitempty"`
	Tags                  map[string]string `json:"tags,omitempty"`
	Region                string            `json:"region,omitempty"`
	Service               string            `json:"service,omitempty"`
	AccountID             string            `json:"account_id,omitempty"`
	
	// Recommendation details
	RecommendationDetails *RecommendationDetails `json:"recommendation_details,omitempty"`
	
	// Temporal data  
	CreatedAt time.Time `json:"created_at"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
	
	// Internal fields managed by CloudshipAI
	Status   string `json:"status,omitempty"`
}

// OpportunityType categorizes different optimization opportunities
type OpportunityType string

const (
	OpportunityRightsizing         OpportunityType = "rightsizing"
	OpportunityIdleResource        OpportunityType = "idle-resource"
	OpportunityCommitment          OpportunityType = "commitment"
	OpportunityStorageOptimization OpportunityType = "storage-optimization"
	OpportunitySpotInstances       OpportunityType = "spot-instances"
)

// RecommendationDetails provides specific implementation guidance
type RecommendationDetails struct {
	CurrentConfiguration  map[string]interface{} `json:"current_configuration"`
	RecommendedConfiguration map[string]interface{} `json:"recommended_configuration"`
	ImplementationSteps   []string               `json:"implementation_steps"`
	Prerequisites         []string               `json:"prerequisites,omitempty"`
	RollbackPlan          []string               `json:"rollback_plan,omitempty"`
	TestingGuidance       string                 `json:"testing_guidance,omitempty"`
}

// Insight represents a cost insight or anomaly
type Insight struct {
	ID              string                 `json:"id"`
	Type            InsightType            `json:"type"`
	Severity        InsightSeverity        `json:"severity"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Category        string                 `json:"category"`
	ResourcesAffected []string             `json:"resources_affected"`
	CostImpact      float64                `json:"cost_impact_usd"`
	TimeRange       TimePeriod             `json:"time_range"`
	Vendor          VendorType             `json:"vendor"`
	Region          string                 `json:"region,omitempty"`
	Service         string                 `json:"service,omitempty"`
	Tags            map[string]string      `json:"tags,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	
	// Internal fields
	Status   string `json:"status,omitempty"`
}

// InsightType categorizes different types of insights
type InsightType string

const (
	InsightTypeCostAnomaly      InsightType = "cost_anomaly"
	InsightTypeUsageAnomaly     InsightType = "usage_anomaly"
	InsightTypeTrend            InsightType = "trend"
	InsightTypeCompliance       InsightType = "compliance"
	InsightTypeGovernance       InsightType = "governance"
	InsightTypeOptimization     InsightType = "optimization"
)

// InsightSeverity indicates the importance level of an insight
type InsightSeverity string

const (
	InsightSeverityLow      InsightSeverity = "low"
	InsightSeverityMedium   InsightSeverity = "medium"
	InsightSeverityHigh     InsightSeverity = "high"
	InsightSeverityCritical InsightSeverity = "critical"
)

// LighthouseConfig configures the lighthouse reporter  
type LighthouseConfig struct {
	// Connection settings (Station is on localhost)
	Timeout       time.Duration `json:"timeout"`
	RetryAttempts int           `json:"retry_attempts"`
	BatchSize     int           `json:"batch_size"`
	
	// Data validation
	ValidateSchema bool `json:"validate_schema"`
	
	// Debugging
	EnableTracing bool `json:"enable_tracing"`
}

// ReportRequest is the base structure for all lighthouse reports
type ReportRequest struct {
	DataType  string                 `json:"data_type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Version   string                 `json:"version"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ResourceInventoryRequest reports discovered resources
type ResourceInventoryRequest struct {
	ReportRequest
	Resources []Resource `json:"resources"`
}

// OpportunityRequest reports optimization opportunities
type OpportunityRequest struct {
	ReportRequest
	Opportunities []Opportunity `json:"opportunities"`
}

// CostDataRequest reports cost information
type CostDataRequest struct {
	ReportRequest
	CostRecords []CostRecord `json:"cost_records"`
}

// InsightRequest reports insights and anomalies
type InsightRequest struct {
	ReportRequest
	Insights []Insight `json:"insights"`
}

// ReportResponse is returned by lighthouse operations
type ReportResponse struct {
	Success    bool                   `json:"success"`
	RequestID  string                 `json:"request_id"`
	RecordsProcessed int              `json:"records_processed"`
	Errors     []ReportError          `json:"errors,omitempty"`
	Warnings   []ReportWarning        `json:"warnings,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

// ReportError represents an error in lighthouse reporting
type ReportError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ReportWarning represents a warning in lighthouse reporting
type ReportWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}