package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
)

// MemorySchemaLearner implements SchemaLearner using in-memory storage
type MemorySchemaLearner struct {
	client    *dagger.Client
	module    *modules.SteampipeModule
	schemas   map[string]*TableSchema // key: provider.table_name
	mutex     sync.RWMutex
	memory    *AgentMemory
}

// NewMemorySchemaLearner creates a new in-memory schema learner
func NewMemorySchemaLearner(client *dagger.Client, memory *AgentMemory) *MemorySchemaLearner {
	return &MemorySchemaLearner{
		client:  client,
		module:  modules.NewSteampipeModule(client),
		schemas: make(map[string]*TableSchema),
		memory:  memory,
	}
}

// LearnSchema discovers and caches table schemas for a provider
func (l *MemorySchemaLearner) LearnSchema(ctx context.Context, provider string, credentials map[string]string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	slog.Info("Learning schemas for provider", "provider", provider)

	// Get list of common tables for the provider
	tables := l.getCommonTables(provider)
	
	for _, tableName := range tables {
		schema, err := l.discoverTableSchema(ctx, provider, tableName, credentials)
		if err != nil {
			slog.Debug("Failed to discover schema for table", "table", tableName, "error", err)
			continue
		}
		
		key := fmt.Sprintf("%s.%s", provider, tableName)
		l.schemas[key] = schema
		
		// Also store in agent memory
		if l.memory != nil {
			if l.memory.Schemas == nil {
				l.memory.Schemas = make(map[string]TableSchema)
			}
			l.memory.Schemas[key] = *schema
		}
		
		slog.Debug("Learned schema for table", "table", tableName, "columns", len(schema.Columns))
	}

	slog.Info("Schema learning completed", "provider", provider, "tables_learned", len(tables))
	return nil
}

// GetSchema retrieves a cached table schema
func (l *MemorySchemaLearner) GetSchema(provider, tableName string) (*TableSchema, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	key := fmt.Sprintf("%s.%s", provider, tableName)
	schema, exists := l.schemas[key]
	return schema, exists
}

// RefreshSchema updates the schema for a specific table
func (l *MemorySchemaLearner) RefreshSchema(ctx context.Context, provider, tableName string, credentials map[string]string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	schema, err := l.discoverTableSchema(ctx, provider, tableName, credentials)
	if err != nil {
		return fmt.Errorf("failed to refresh schema for %s.%s: %w", provider, tableName, err)
	}

	key := fmt.Sprintf("%s.%s", provider, tableName)
	l.schemas[key] = schema

	// Update agent memory
	if l.memory != nil && l.memory.Schemas != nil {
		l.memory.Schemas[key] = *schema
	}

	slog.Info("Refreshed schema for table", "table", tableName, "provider", provider)
	return nil
}

// discoverTableSchema queries Steampipe to discover the actual table schema
func (l *MemorySchemaLearner) discoverTableSchema(ctx context.Context, provider, tableName string, credentials map[string]string) (*TableSchema, error) {
	// Query the information_schema to get column information
	schemaQuery := fmt.Sprintf(`
		SELECT 
			column_name,
			data_type,
			is_nullable
		FROM information_schema.columns 
		WHERE table_name = '%s' 
		ORDER BY ordinal_position
	`, tableName)

	result, err := l.module.RunQuery(ctx, provider, schemaQuery, credentials, "json")
	if err != nil {
		return nil, fmt.Errorf("failed to query table schema: %w", err)
	}

	// Parse the results to build schema
	columns := l.parseSchemaResult(result)
	
	// Enhance with known column descriptions and examples
	enhancedColumns := l.enhanceColumnInfo(provider, tableName, columns)

	schema := &TableSchema{
		TableName:   tableName,
		Columns:     enhancedColumns,
		Provider:    provider,
		Description: l.getTableDescription(provider, tableName),
		LastUpdated: "now", // Would use actual timestamp
	}

	return schema, nil
}

// parseSchemaResult parses the JSON result from information_schema query
func (l *MemorySchemaLearner) parseSchemaResult(result string) []ColumnInfo {
	// For now, return basic columns. In production, this would parse JSON
	// and extract actual column information from Steampipe
	return []ColumnInfo{
		{Name: "name", Type: "text", Description: "Resource name"},
		{Name: "id", Type: "text", Description: "Resource identifier"},
		{Name: "region", Type: "text", Description: "AWS region"},
	}
}

// enhanceColumnInfo adds descriptions and examples to column information
func (l *MemorySchemaLearner) enhanceColumnInfo(provider, tableName string, columns []ColumnInfo) []ColumnInfo {
	enhanced := make([]ColumnInfo, len(columns))
	copy(enhanced, columns)

	// Add known information based on common patterns
	for i, col := range enhanced {
		enhanced[i] = l.enhanceColumn(provider, tableName, col)
	}

	return enhanced
}

// enhanceColumn adds metadata to a single column
func (l *MemorySchemaLearner) enhanceColumn(provider, tableName string, col ColumnInfo) ColumnInfo {
	enhanced := col

	// Add common descriptions and examples
	switch col.Name {
	case "instance_id":
		enhanced.Description = "EC2 instance identifier"
		enhanced.Examples = []string{"i-1234567890abcdef0"}
	case "instance_state":
		enhanced.Description = "Current state of the EC2 instance"
		enhanced.Examples = []string{"running", "stopped", "pending", "terminated"}
	case "instance_type":
		enhanced.Description = "EC2 instance type/size"
		enhanced.Examples = []string{"t3.micro", "m5.large", "c5.xlarge"}
	case "vpc_id":
		enhanced.Description = "VPC identifier where resource is located"
		enhanced.Examples = []string{"vpc-12345678"}
	case "region", "aws_region":
		enhanced.Description = "AWS region where resource is located"
		enhanced.Examples = []string{"us-east-1", "us-west-2", "eu-west-1"}
	case "name":
		enhanced.Description = "Resource name or identifier"
	case "tags":
		enhanced.Description = "Resource tags as JSON object"
		enhanced.Examples = []string{`{"Environment": "prod", "Team": "platform"}`}
	}

	return enhanced
}

// getTableDescription returns a description for a known table
func (l *MemorySchemaLearner) getTableDescription(provider, tableName string) string {
	descriptions := map[string]string{
		"aws.aws_ec2_instance":        "EC2 virtual machine instances",
		"aws.aws_s3_bucket":          "S3 storage buckets",
		"aws.aws_rds_db_instance":    "RDS database instances",
		"aws.aws_vpc_security_group": "VPC security groups and rules",
		"aws.aws_iam_user":           "IAM users and their configurations",
		"aws.aws_iam_role":           "IAM roles and their policies",
		"aws.aws_vpc":                "Virtual Private Clouds (VPCs)",
		"aws.aws_lambda_function":    "Lambda serverless functions",
	}

	key := fmt.Sprintf("%s.%s", provider, tableName)
	if desc, exists := descriptions[key]; exists {
		return desc
	}

	return fmt.Sprintf("%s table for %s provider", tableName, provider)
}

// getCommonTables returns a list of commonly used tables for a provider
func (l *MemorySchemaLearner) getCommonTables(provider string) []string {
	switch provider {
	case "aws":
		return []string{
			"aws_account",
			"aws_ec2_instance",
			"aws_s3_bucket",
			"aws_rds_db_instance",
			"aws_vpc_security_group",
			"aws_iam_user",
			"aws_iam_role",
			"aws_vpc",
			"aws_lambda_function",
			"aws_iam_policy",
			"aws_iam_access_key",
		}
	case "azure":
		return []string{
			"azure_subscription",
			"azure_compute_virtual_machine",
			"azure_storage_account",
			"azure_sql_server",
		}
	case "gcp":
		return []string{
			"gcp_project",
			"gcp_compute_instance",
			"gcp_storage_bucket",
			"gcp_sql_database_instance",
		}
	default:
		return []string{}
	}
}

// GetProviderTables returns all learned tables for a provider
func (l *MemorySchemaLearner) GetProviderTables(provider string) []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var tables []string
	prefix := provider + "."
	
	for key := range l.schemas {
		if strings.HasPrefix(key, prefix) {
			tableName := strings.TrimPrefix(key, prefix)
			tables = append(tables, tableName)
		}
	}

	return tables
}

// GetColumnNames returns just the column names for a table
func (l *MemorySchemaLearner) GetColumnNames(provider, tableName string) []string {
	schema, exists := l.GetSchema(provider, tableName)
	if !exists {
		return []string{}
	}

	names := make([]string, len(schema.Columns))
	for i, col := range schema.Columns {
		names[i] = col.Name
	}

	return names
}

// GenerateSchemaPrompt creates a prompt section with schema information
func (l *MemorySchemaLearner) GenerateSchemaPrompt(provider string, relevantTables []string) string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var prompt strings.Builder
	prompt.WriteString(fmt.Sprintf("\nAvailable %s tables and their schemas:\n", provider))

	for _, tableName := range relevantTables {
		key := fmt.Sprintf("%s.%s", provider, tableName)
		schema, exists := l.schemas[key]
		if !exists {
			continue
		}

		prompt.WriteString(fmt.Sprintf("\nTable: %s\n", tableName))
		prompt.WriteString(fmt.Sprintf("Description: %s\n", schema.Description))
		prompt.WriteString("Columns:\n")

		for _, col := range schema.Columns {
			prompt.WriteString(fmt.Sprintf("  - %s (%s): %s\n", col.Name, col.Type, col.Description))
			if len(col.Examples) > 0 {
				prompt.WriteString(fmt.Sprintf("    Examples: %s\n", strings.Join(col.Examples, ", ")))
			}
		}
	}

	return prompt.String()
}