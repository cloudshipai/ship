package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSteampipeTool_Info(t *testing.T) {
	tool := &SteampipeTool{
		toolInfo: &schema.ToolInfo{
			Name: "steampipe_query",
			Desc: "Test description",
		},
	}

	ctx := context.Background()
	info, err := tool.Info(ctx)
	
	require.NoError(t, err)
	assert.Equal(t, "steampipe_query", info.Name)
	assert.Equal(t, "Test description", info.Desc)
}

func TestSteampipeTool_ParameterSchema(t *testing.T) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	ctx := context.Background()
	info, err := tool.Info(ctx)
	
	require.NoError(t, err)
	assert.NotNil(t, info.ParamsOneOf)
	
	// Test that the tool has the expected parameters
	assert.Equal(t, "steampipe_query", info.Name)
	assert.Contains(t, info.Desc, "Execute SQL queries against cloud infrastructure")
}

func TestSteampipeRequest_JSONParsing(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected SteampipeRequest
		wantErr  bool
	}{
		{
			name:    "valid request",
			jsonStr: `{"provider": "aws", "query": "SELECT * FROM aws_ec2_instance", "credentials": {"key": "value"}}`,
			expected: SteampipeRequest{
				Provider:    "aws",
				Query:       "SELECT * FROM aws_ec2_instance",
				Credentials: map[string]string{"key": "value"},
			},
			wantErr: false,
		},
		{
			name:    "minimal request",
			jsonStr: `{"provider": "aws", "query": "SELECT 1"}`,
			expected: SteampipeRequest{
				Provider: "aws",
				Query:    "SELECT 1",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			jsonStr: `{"provider": "aws", "query":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req SteampipeRequest
			err := json.Unmarshal([]byte(tt.jsonStr), &req)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Provider, req.Provider)
			assert.Equal(t, tt.expected.Query, req.Query)
			
			if tt.expected.Credentials != nil {
				assert.Equal(t, tt.expected.Credentials, req.Credentials)
			}
		})
	}
}

func TestSteampipeTool_QueryImprovement(t *testing.T) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	tests := []struct {
		name     string
		query    string
		provider string
		expected string
	}{
		{
			name:     "fix instance state column",
			query:    "SELECT * FROM aws_ec2_instance WHERE state = 'running'",
			provider: "aws",
			expected: "SELECT * FROM aws_ec2_instance WHERE instance_state = 'running'",
		},
		{
			name:     "fix security group field",
			query:    "SELECT sg.group_id FROM aws_vpc_security_group sg",
			provider: "aws",
			expected: "SELECT sg->>'GroupId' FROM aws_vpc_security_group sg",
		},
		{
			name:     "handle multiple statements",
			query:    "SELECT 1; DROP TABLE users; SELECT 2",
			provider: "aws",
			expected: "SELECT 1",
		},
		{
			name:     "no changes needed",
			query:    "SELECT instance_id FROM aws_ec2_instance",
			provider: "aws",
			expected: "SELECT instance_id FROM aws_ec2_instance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			improved, err := tool.improveQuery(tt.query, tt.provider)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, improved)
		})
	}
}

func TestSteampipeTool_Memory(t *testing.T) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	// Test recording success
	tool.recordSuccess("original query", "executed query", "aws", 5)
	assert.Len(t, memory.Successes, 1)
	assert.Equal(t, "original query", memory.Successes[0].OriginalIntent)
	assert.Equal(t, "executed query", memory.Successes[0].GeneratedQuery)
	assert.Equal(t, 5, memory.Successes[0].ResultCount)
	assert.Equal(t, "aws", memory.Successes[0].Provider)
	
	// Test recording failure
	tool.recordFailure("bad query", "aws", "column does not exist")
	assert.Len(t, memory.Failures, 1)
	assert.Equal(t, "bad query", memory.Failures[0].OriginalIntent)
	assert.Equal(t, "schema", memory.Failures[0].ErrorType)
	assert.Contains(t, memory.Failures[0].LessonLearned, "schema understanding")
}

func TestSteampipeTool_ErrorClassification(t *testing.T) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	tests := []struct {
		name        string
		errorMsg    string
		expectedType string
	}{
		{
			name:        "schema error",
			errorMsg:    "column \"invalid_column\" does not exist",
			expectedType: "schema",
		},
		{
			name:        "syntax error",
			errorMsg:    "syntax error at or near \"INVALID\"",
			expectedType: "syntax",
		},
		{
			name:        "authentication error",
			errorMsg:    "authentication failed for AWS",
			expectedType: "auth",
		},
		{
			name:        "timeout error",
			errorMsg:    "query timeout exceeded",
			expectedType: "timeout",
		},
		{
			name:        "unknown error",
			errorMsg:    "some other error occurred",
			expectedType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool.recordFailure("test query", "aws", tt.errorMsg)
			
			require.Len(t, memory.Failures, 1)
			assert.Equal(t, tt.expectedType, memory.Failures[0].ErrorType)
			
			// Reset for next test
			memory.Failures = make([]QueryFailure, 0)
		})
	}
}

func TestSteampipeTool_LessonGeneration(t *testing.T) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	tests := []struct {
		name           string
		errorMsg       string
		expectedLesson string
	}{
		{
			name:           "state column error",
			errorMsg:       "column \"state\" does not exist",
			expectedLesson: "Use 'instance_state' instead of 'state' for EC2 instance queries",
		},
		{
			name:           "running column error",
			errorMsg:       "column \"running\" does not exist",
			expectedLesson: "Use 'instance_state = \"running\"' instead of 'running' column",
		},
		{
			name:           "group_id error",
			errorMsg:       "column \"group_id\" does not exist in security groups",
			expectedLesson: "Use JSONB operators for security group fields: sg->>'GroupId'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lesson := tool.generateLessonFromError(tt.errorMsg)
			assert.Equal(t, tt.expectedLesson, lesson)
		})
	}
}

func TestSteampipeTool_Insights(t *testing.T) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	tests := []struct {
		name     string
		results  []map[string]interface{}
		provider string
		query    string
		expected []string
	}{
		{
			name:     "no results",
			results:  []map[string]interface{}{},
			provider: "aws",
			query:    "SELECT * FROM aws_ec2_instance",
			expected: []string{"No results found - consider broadening the query scope"},
		},
		{
			name: "single result",
			results: []map[string]interface{}{
				{"instance_id": "i-123456"},
			},
			provider: "aws",
			query:    "SELECT * FROM aws_ec2_instance",
			expected: []string{"Found 1 result", "Consider checking instance security groups and tags"},
		},
		{
			name: "multiple results with S3",
			results: []map[string]interface{}{
				{"bucket_name": "bucket1"},
				{"bucket_name": "bucket2"},
				{"bucket_name": "bucket3"},
			},
			provider: "aws",
			query:    "SELECT * FROM aws_s3_bucket",
			expected: []string{"Found 3 results", "Consider checking bucket encryption and public access settings"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insights := tool.generateInsights(tt.results, tt.provider, tt.query)
			
			for _, expectedInsight := range tt.expected {
				assert.Contains(t, insights, expectedInsight)
			}
		})
	}
}

// Benchmark tests for performance measurement
func BenchmarkSteampipeTool_QueryImprovement(b *testing.B) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	query := "SELECT * FROM aws_ec2_instance WHERE state = 'running' AND sg.group_id = 'sg-123'"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tool.improveQuery(query, "aws")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSteampipeTool_MemoryOperations(b *testing.B) {
	memory := &AgentMemory{
		Successes: make([]QuerySuccess, 0),
		Failures:  make([]QueryFailure, 0),
	}
	
	tool := NewSteampipeTool(nil, memory)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate alternating success and failure recording
		if i%2 == 0 {
			tool.recordSuccess("query", "query", "aws", 10)
		} else {
			tool.recordFailure("query", "aws", "error")
		}
	}
}