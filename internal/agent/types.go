package agent

import (
	"context"
)

// InvestigationRequest represents a user's natural language investigation request
type InvestigationRequest struct {
	Prompt    string            `json:"prompt"`
	Provider  string            `json:"provider"`  // aws, azure, gcp
	Region    string            `json:"region,omitempty"`
	Credentials map[string]string `json:"-"` // Not serialized for security
}

// InvestigationResult represents the final investigation results
type InvestigationResult struct {
	Success     bool                     `json:"success"`
	Steps       []InvestigationStep     `json:"steps"`
	Summary     string                  `json:"summary"`
	Insights    []Insight               `json:"insights"`
	QueryCount  int                     `json:"query_count"`
	Duration    string                  `json:"duration"`
	Confidence  float64                 `json:"confidence"`
}

// InvestigationStep represents a single step in the investigation
type InvestigationStep struct {
	StepNumber       int                    `json:"step_number"`
	Description      string                 `json:"description"`
	Query            string                 `json:"query"`
	Results          []map[string]interface{} `json:"results"`
	Success          bool                   `json:"success"`
	Error            string                 `json:"error,omitempty"`
	ExecutionTime    string                 `json:"execution_time"`
	Insights         []string               `json:"insights"`
}

// Insight represents a key finding from the investigation
type Insight struct {
	Type        string  `json:"type"`        // security, cost, performance, compliance
	Severity    string  `json:"severity"`    // critical, high, medium, low, info
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Recommendation string `json:"recommendation"`
	Confidence  float64 `json:"confidence"`
}

// TableSchema represents Steampipe table schema information
type TableSchema struct {
	TableName   string        `json:"table_name"`
	Columns     []ColumnInfo  `json:"columns"`
	Provider    string        `json:"provider"`
	Description string        `json:"description"`
	LastUpdated string        `json:"last_updated"`
}

// ColumnInfo represents information about a table column
type ColumnInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Examples    []string `json:"examples,omitempty"`
}

// QueryPattern represents a learned query pattern for reuse
type QueryPattern struct {
	ID            string            `json:"id"`
	Intent        string            `json:"intent"`        // Natural language intent
	Template      string            `json:"template"`      // SQL template with placeholders
	Parameters    []string          `json:"parameters"`    // Parameter names
	Provider      string            `json:"provider"`
	SuccessRate   float64           `json:"success_rate"`
	UsageCount    int               `json:"usage_count"`
	Examples      []string          `json:"examples"`      // Example natural language inputs
	Tags          []string          `json:"tags"`
	CreatedAt     string            `json:"created_at"`
	LastUsed      string            `json:"last_used"`
}

// AgentMemory represents the agent's persistent knowledge
type AgentMemory struct {
	Schemas      map[string]TableSchema  `json:"schemas"`       // provider.table_name -> schema
	Patterns     map[string]QueryPattern `json:"patterns"`      // pattern_id -> pattern
	Successes    []QuerySuccess          `json:"successes"`     // Recent successful queries
	Failures     []QueryFailure          `json:"failures"`      // Recent failures for learning
	LastUpdate   string                  `json:"last_update"`
}

// QuerySuccess represents a successful query execution
type QuerySuccess struct {
	OriginalIntent string    `json:"original_intent"`
	GeneratedQuery string    `json:"generated_query"`
	ResultCount    int       `json:"result_count"`
	ExecutionTime  string    `json:"execution_time"`
	Provider       string    `json:"provider"`
	Timestamp      string    `json:"timestamp"`
	PatternUsed    string    `json:"pattern_used,omitempty"`
}

// QueryFailure represents a failed query with learning information
type QueryFailure struct {
	OriginalIntent string `json:"original_intent"`
	GeneratedQuery string `json:"generated_query"`
	ErrorMessage   string `json:"error_message"`
	ErrorType      string `json:"error_type"`      // syntax, schema, auth, timeout
	Provider       string `json:"provider"`
	Timestamp      string `json:"timestamp"`
	LessonLearned  string `json:"lesson_learned"`  // What the agent learned
}

// SteampipeExecutor defines the interface for executing Steampipe queries
type SteampipeExecutor interface {
	ExecuteQuery(ctx context.Context, provider, query string, credentials map[string]string) ([]map[string]interface{}, error)
	GetTableSchema(ctx context.Context, provider, tableName string, credentials map[string]string) (*TableSchema, error)
	GetAvailableTables(ctx context.Context, provider string, credentials map[string]string) ([]string, error)
	ValidateQuery(query string) error
}

// SchemaLearner defines the interface for learning table schemas
type SchemaLearner interface {
	LearnSchema(ctx context.Context, provider string, credentials map[string]string) error
	GetSchema(provider, tableName string) (*TableSchema, bool)
	RefreshSchema(ctx context.Context, provider, tableName string, credentials map[string]string) error
}

// PatternMatcher defines the interface for matching intents to query patterns
type PatternMatcher interface {
	FindBestPattern(intent, provider string) (*QueryPattern, float64)
	LearnPattern(intent, query, provider string, success bool) error
	UpdatePatternSuccess(patternID string, success bool) error
}

// InvestigationAgent represents the main Eino agent for infrastructure investigation
type InvestigationAgent interface {
	Investigate(ctx context.Context, request InvestigationRequest) (*InvestigationResult, error)
	GetMemory() *AgentMemory
	SaveMemory() error
	LoadMemory() error
}

// EinoTool represents a tool that can be used by the Eino agent
type EinoTool interface {
	Name() string
	Description() string
	InputSchema() interface{}
	Call(ctx context.Context, input string) (string, error)
}