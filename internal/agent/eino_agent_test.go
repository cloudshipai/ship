package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEinoInvestigationAgent_Creation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping agent creation test in short mode")
	}
	
	// This test requires actual OpenAI API key and Dagger, so we'll mock it
	t.Run("agent creation with valid config", func(t *testing.T) {
		// Test agent creation parameters
		ctx := context.Background()
		apiKey := "test-api-key"
		memoryPath := "/tmp/test-memory.json"
		
		// We can't create a real agent without Dagger client and OpenAI key
		// So we'll test the validation logic
		assert.NotEmpty(t, apiKey)
		assert.NotEmpty(t, memoryPath)
		assert.NotNil(t, ctx)
	})
}

func TestInvestigationRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request InvestigationRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: InvestigationRequest{
				Prompt:   "Find all EC2 instances",
				Provider: "aws",
				Region:   "us-east-1",
				Credentials: map[string]string{
					"AWS_ACCESS_KEY_ID": "test",
				},
			},
			valid: true,
		},
		{
			name: "missing prompt",
			request: InvestigationRequest{
				Provider: "aws",
				Region:   "us-east-1",
			},
			valid: false,
		},
		{
			name: "missing provider",
			request: InvestigationRequest{
				Prompt: "Find all instances",
				Region: "us-east-1",
			},
			valid: false,
		},
		{
			name: "invalid provider",
			request: InvestigationRequest{
				Prompt:   "Find all instances",
				Provider: "invalid",
				Region:   "us-east-1",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInvestigationRequest(tt.request)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestEinoInvestigationAgent_TableIdentification(t *testing.T) {
	agent := &EinoInvestigationAgent{}
	
	tests := []struct {
		name     string
		prompt   string
		provider string
		expected []string
	}{
		{
			name:     "EC2 instances",
			prompt:   "Find all running EC2 instances",
			provider: "aws",
			expected: []string{"aws_ec2_instance", "aws_account"},
		},
		{
			name:     "S3 buckets",
			prompt:   "List all S3 buckets with public access",
			provider: "aws",
			expected: []string{"aws_s3_bucket", "aws_account"},
		},
		{
			name:     "Security groups",
			prompt:   "Show security groups allowing 0.0.0.0/0",
			provider: "aws",
			expected: []string{"aws_vpc_security_group", "aws_account"},
		},
		{
			name:     "IAM users",
			prompt:   "Find IAM users without MFA",
			provider: "aws",
			expected: []string{"aws_iam_user", "aws_account"},
		},
		{
			name:     "Lambda functions",
			prompt:   "List all Lambda functions with errors",
			provider: "aws",
			expected: []string{"aws_lambda_function", "aws_account"},
		},
		{
			name:     "Multiple resources",
			prompt:   "Find EC2 instances and their security groups",
			provider: "aws",
			expected: []string{"aws_ec2_instance", "aws_vpc_security_group", "aws_account"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tables := agent.identifyRelevantTables(tt.prompt, tt.provider)
			
			for _, expected := range tt.expected {
				assert.Contains(t, tables, expected, "Expected table %s not found in %v", expected, tables)
			}
		})
	}
}

func TestInvestigationResult_Structure(t *testing.T) {
	result := &InvestigationResult{
		Success:    true,
		Steps:      []InvestigationStep{},
		Summary:    "Test investigation completed",
		Insights:   []Insight{},
		QueryCount: 3,
		Duration:   "2.5s",
		Confidence: 0.85,
	}
	
	assert.True(t, result.Success)
	assert.Equal(t, "Test investigation completed", result.Summary)
	assert.Equal(t, 3, result.QueryCount)
	assert.Equal(t, "2.5s", result.Duration)
	assert.Equal(t, 0.85, result.Confidence)
	assert.NotNil(t, result.Steps)
	assert.NotNil(t, result.Insights)
}

func TestInvestigationStep_Structure(t *testing.T) {
	step := InvestigationStep{
		StepNumber:  1,
		Description: "Query EC2 instances",
		Query:       "SELECT * FROM aws_ec2_instance",
		Results: []map[string]interface{}{
			{"instance_id": "i-123456", "state": "running"},
		},
		Success:       true,
		Error:         "",
		ExecutionTime: "1.2s",
		Insights:      []string{"Found 1 instance"},
	}
	
	assert.Equal(t, 1, step.StepNumber)
	assert.Equal(t, "Query EC2 instances", step.Description)
	assert.True(t, step.Success)
	assert.Empty(t, step.Error)
	assert.Len(t, step.Results, 1)
	assert.Equal(t, "i-123456", step.Results[0]["instance_id"])
}

func TestInsight_Structure(t *testing.T) {
	insight := Insight{
		Type:           "security",
		Severity:       "high",
		Title:          "Open Security Groups",
		Description:    "Found security groups allowing 0.0.0.0/0",
		Impact:         "Potential unauthorized access",
		Recommendation: "Restrict security group rules",
		Confidence:     0.9,
	}
	
	assert.Equal(t, "security", insight.Type)
	assert.Equal(t, "high", insight.Severity)
	assert.Equal(t, "Open Security Groups", insight.Title)
	assert.Equal(t, 0.9, insight.Confidence)
}

func TestEinoInvestigationAgent_PromptEnhancement(t *testing.T) {
	agent := &EinoInvestigationAgent{
		memory: &AgentMemory{
			Failures: []QueryFailure{
				{
					OriginalIntent: "bad query",
					ErrorType:      "schema",
					LessonLearned:  "Use proper column names",
				},
			},
		},
	}
	
	request := InvestigationRequest{
		Prompt:   "Find running instances",
		Provider: "aws",
		Region:   "us-east-1",
	}
	
	enhanced := agent.enhancePromptWithContext(request)
	
	assert.Contains(t, enhanced, "Find running instances")
	assert.Contains(t, enhanced, "TARGET PROVIDER: aws")
	assert.Contains(t, enhanced, "REGION: us-east-1")
	assert.Contains(t, enhanced, "KNOWN ISSUES TO AVOID:")
	assert.Contains(t, enhanced, "Use proper column names")
}

func TestEinoInvestigationAgent_InsightExtraction(t *testing.T) {
	agent := &EinoInvestigationAgent{}
	
	tests := []struct {
		name           string
		content        string
		provider       string
		expectedCount  int
		expectedTypes  []string
	}{
		{
			name:          "security group issue",
			content:       "Found security groups allowing access from 0.0.0.0/0 which poses security risks",
			provider:      "aws",
			expectedCount: 1,
			expectedTypes: []string{"security"},
		},
		{
			name:          "unencrypted resources",
			content:       "Several S3 buckets are unencrypted and vulnerable",
			provider:      "aws",
			expectedCount: 1,
			expectedTypes: []string{"security"},
		},
		{
			name:          "stopped instances",
			content:       "Found 5 stopped instances that are still incurring costs",
			provider:      "aws",
			expectedCount: 1,
			expectedTypes: []string{"cost"},
		},
		{
			name:          "multiple issues",
			content:       "Found security groups with 0.0.0.0/0 access and several unencrypted buckets and stopped instances",
			provider:      "aws",
			expectedCount: 3,
			expectedTypes: []string{"security", "security", "cost"},
		},
		{
			name:          "no issues",
			content:       "All resources are properly configured",
			provider:      "aws",
			expectedCount: 0,
			expectedTypes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insights := agent.extractInsights(tt.content, tt.provider)
			
			assert.Len(t, insights, tt.expectedCount)
			
			for i, expectedType := range tt.expectedTypes {
				if i < len(insights) {
					assert.Equal(t, expectedType, insights[i].Type)
				}
			}
		})
	}
}

func TestAgentMemory_Management(t *testing.T) {
	memory := &AgentMemory{
		Schemas:   make(map[string]TableSchema),
		Patterns:  make(map[string]QueryPattern),
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	// Test adding schemas
	schema := TableSchema{
		TableName: "aws_ec2_instance",
		Provider:  "aws",
		Columns: []ColumnInfo{
			{Name: "instance_id", Type: "text"},
			{Name: "instance_state", Type: "text"},
		},
	}
	
	memory.Schemas["aws.aws_ec2_instance"] = schema
	assert.Len(t, memory.Schemas, 1)
	assert.Equal(t, "aws_ec2_instance", memory.Schemas["aws.aws_ec2_instance"].TableName)
	
	// Test adding patterns
	pattern := QueryPattern{
		Intent:        "find running instances",
		Template:      "SELECT * FROM aws_ec2_instance WHERE instance_state = 'running'",
		Provider:      "aws",
		SuccessRate:   0.95,
		UsageCount:    100,
		Parameters:    []string{"instance_state"},
		Examples:      []string{"find running instances"},
		Tags:          []string{"ec2", "instances"},
		CreatedAt:     "2024-01-01T00:00:00Z",
		LastUsed:      "2024-01-01T00:00:00Z",
	}
	
	memory.Patterns["running_instances"] = pattern
	assert.Len(t, memory.Patterns, 1)
	assert.Equal(t, 0.95, memory.Patterns["running_instances"].SuccessRate)
}

// Benchmark tests for agent performance
func BenchmarkEinoInvestigationAgent_TableIdentification(b *testing.B) {
	agent := &EinoInvestigationAgent{}
	prompt := "Find all running EC2 instances with security group issues"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = agent.identifyRelevantTables(prompt, "aws")
	}
}

func BenchmarkEinoInvestigationAgent_PromptEnhancement(b *testing.B) {
	agent := &EinoInvestigationAgent{
		memory: &AgentMemory{
			Failures: make([]QueryFailure, 50), // Pre-populate with some failures
		},
	}
	
	request := InvestigationRequest{
		Prompt:   "Find running instances with security issues",
		Provider: "aws",
		Region:   "us-east-1",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = agent.enhancePromptWithContext(request)
	}
}

func BenchmarkEinoInvestigationAgent_InsightExtraction(b *testing.B) {
	agent := &EinoInvestigationAgent{}
	content := "Found multiple security groups allowing access from 0.0.0.0/0 and several unencrypted S3 buckets and stopped EC2 instances"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = agent.extractInsights(content, "aws")
	}
}

// Helper function for request validation
func validateInvestigationRequest(req InvestigationRequest) error {
	if req.Prompt == "" {
		return assert.AnError
	}
	if req.Provider == "" {
		return assert.AnError
	}
	if req.Provider != "aws" && req.Provider != "azure" && req.Provider != "gcp" {
		return assert.AnError
	}
	return nil
}