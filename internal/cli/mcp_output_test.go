package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPOutputFileLogging(t *testing.T) {
	if os.Getenv("SHIP_OUTPUT_TEST") == "" {
		t.Skip("Set SHIP_OUTPUT_TEST=1 to run output file tests")
	}

	// Check if ship binary exists
	shipPath := "/home/epuerta/.local/bin/ship"
	if _, err := os.Stat(shipPath); err != nil {
		t.Skip("Ship binary not found, run 'make local-install' first")
	}
	t.Logf("Using ship binary at: %s", shipPath)

	// Create test directory with vulnerable Terraform
	tempDir := t.TempDir()
	vulnerableTerraform := `
terraform {
  required_version = ">= 1.0"
}

provider "aws" {
  region = "us-west-2"
}

# INTENTIONALLY VULNERABLE S3 bucket for testing
resource "aws_s3_bucket" "test_bucket" {
  bucket = "test-vulnerable-bucket-checkov-test"
}

resource "aws_s3_bucket_public_access_block" "test_bucket" {
  bucket = aws_s3_bucket.test_bucket.id

  # INSECURE: Allows public access
  block_public_acls       = false
  block_public_policy     = false  
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# INSECURE: Security group with SSH open to world
resource "aws_security_group" "test_sg" {
  name_prefix = "test-sg"
  description = "Test security group"

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # INSECURE!
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# INSECURE: EC2 instance without encryption
resource "aws_instance" "test_instance" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t3.micro"
  
  root_block_device {
    volume_type = "gp3"
    volume_size = 20
    encrypted   = false  # INSECURE!
  }

  vpc_security_group_ids = [aws_security_group.test_sg.id]
  
  tags = {
    Name = "test-instance"
  }
}
`

	err := os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(vulnerableTerraform), 0644)
	require.NoError(t, err)

	// Create output files
	outputFile := filepath.Join(tempDir, "checkov-output.txt")
	executionLog := filepath.Join(tempDir, "execution.log")

	t.Run("Checkov with Output File", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Create MCP client with output file flags
		clientTransport := transport.NewStdio(shipPath, nil, 
			"mcp", "checkov", 
			"--output-file", outputFile, 
			"--execution-log", executionLog)
		mcpClient := client.NewClient(clientTransport)

		// Start the client
		err := mcpClient.Start(ctx)
		require.NoError(t, err)
		defer mcpClient.Close()

		// Initialize MCP session
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "ship-test-client",
			Version: "1.0.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = mcpClient.Initialize(ctx, initRequest)
		require.NoError(t, err)

		// List available tools
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
		require.NoError(t, err)
		
		// Find checkov_scan_directory tool
		var checkovTool *mcp.Tool
		for i, tool := range toolsResult.Tools {
			if tool.Name == "checkov_scan_directory" {
				checkovTool = &toolsResult.Tools[i]
				break
			}
		}
		require.NotNil(t, checkovTool, "checkov_scan_directory tool should be available")
		t.Logf("Found checkov tool: %s", checkovTool.Name)

		t.Log("Starting Checkov scan with output file logging...")
		startTime := time.Now()

		// Call the checkov_scan_directory tool
		callRequest := mcp.CallToolRequest{}
		callRequest.Params.Name = "checkov_scan_directory"
		callRequest.Params.Arguments = map[string]interface{}{
			"directory": tempDir,
			"framework": "terraform",
			"output":    "cli",
		}

		callResult, err := mcpClient.CallTool(ctx, callRequest)
		elapsed := time.Since(startTime)
		
		require.NoError(t, err)
		assert.NotEmpty(t, callResult.Content)

		t.Logf("Checkov scan completed in: %v", elapsed)

		// Verify MCP response contains results
		if len(callResult.Content) > 0 {
			if textContent, ok := callResult.Content[0].(mcp.TextContent); ok {
				assert.NotEmpty(t, textContent.Text)
				assert.Contains(t, textContent.Text, "checkov", "Response should contain checkov output")
				t.Logf("MCP response size: %d bytes", len(textContent.Text))
				
				// Log a preview of the response
				preview := textContent.Text
				if len(preview) > 500 {
					preview = preview[:500] + "..."
				}
				t.Logf("MCP response preview: %s", preview)
			}
		}

		// Give a moment for file writes to complete
		time.Sleep(2 * time.Second)

		// Verify output file was created and contains results
		if _, err := os.Stat(outputFile); err != nil {
			t.Errorf("Output file was not created: %s", outputFile)
		} else {
			outputContent, err := os.ReadFile(outputFile)
			require.NoError(t, err)
			
			outputStr := string(outputContent)
			assert.NotEmpty(t, outputStr, "Output file should not be empty")
			assert.Contains(t, outputStr, "Ship CLI Output", "Output file should contain Ship header")
			assert.Contains(t, outputStr, "checkov", "Output file should contain checkov results")
			assert.Contains(t, outputStr, "Duration:", "Output file should contain timing info")
			
			t.Logf("Output file created successfully: %s (%d bytes)", outputFile, len(outputContent))
			
			// Check for security findings
			if strings.Contains(outputStr, "failed") || strings.Contains(outputStr, "FAILED") {
				t.Log("✅ Checkov found security issues as expected")
			}
			
			// Log a preview
			lines := strings.Split(outputStr, "\n")
			previewLines := 10
			if len(lines) > previewLines {
				preview := strings.Join(lines[:previewLines], "\n")
				t.Logf("Output file preview (first %d lines):\n%s", previewLines, preview)
			}
		}

		// Verify execution log was created
		if _, err := os.Stat(executionLog); err != nil {
			t.Errorf("Execution log was not created: %s", executionLog)
		} else {
			logContent, err := os.ReadFile(executionLog)
			require.NoError(t, err)
			
			logStr := string(logContent)
			assert.NotEmpty(t, logStr, "Execution log should not be empty")
			assert.Contains(t, logStr, "checkov", "Execution log should contain command")
			assert.Contains(t, logStr, "Duration:", "Execution log should contain timing")
			assert.Contains(t, logStr, "Success:", "Execution log should contain success status")
			
			t.Logf("Execution log created successfully: %s (%d bytes)", executionLog, len(logContent))
			t.Logf("Execution log content: %s", strings.TrimSpace(logStr))
		}
	})

	t.Run("Multiple Tool Executions", func(t *testing.T) {
		t.Skip("Skipping multi-execution test for now - focus on single tool test first")
	})
}

// TestMCPOutputFilePerformance tests performance impact of file logging
func TestMCPOutputFilePerformance(t *testing.T) {
	if os.Getenv("SHIP_PERFORMANCE_TEST") == "" {
		t.Skip("Set SHIP_PERFORMANCE_TEST=1 to run performance tests")
	}

	shipPath := "/home/epuerta/.local/bin/ship"
	if _, err := os.Stat(shipPath); err != nil {
		t.Skip("Ship binary not found")
	}

	// Create test directory
	tempDir := t.TempDir()
	tfContent := `
resource "aws_instance" "test" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"
  
  root_block_device {
    encrypted = false
  }
}
`
	err := os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(tfContent), 0644)
	require.NoError(t, err)

	runTest := func(name string, useLogging bool) time.Duration {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var args []string
		if useLogging {
			outputFile := filepath.Join(tempDir, "perf-output.txt")
			executionLog := filepath.Join(tempDir, "perf-execution.log") 
			args = []string{"mcp", "checkov", "--output-file", outputFile, "--execution-log", executionLog}
		} else {
			args = []string{"mcp", "checkov"}
		}

		clientTransport := transport.NewStdio(shipPath, nil, args...)
		mcpClient := client.NewClient(clientTransport)

		start := time.Now()
		
		err := mcpClient.Start(ctx)
		require.NoError(t, err)
		defer mcpClient.Close()

		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{Name: "perf-test", Version: "1.0.0"}

		_, err = mcpClient.Initialize(ctx, initRequest)
		require.NoError(t, err)

		callRequest := mcp.CallToolRequest{}
		callRequest.Params.Name = "checkov_scan_directory"
		callRequest.Params.Arguments = map[string]interface{}{
			"directory": tempDir,
			"framework": "terraform",
			"output":    "cli",
		}

		_, err = mcpClient.CallTool(ctx, callRequest)
		require.NoError(t, err)

		elapsed := time.Since(start)
		t.Logf("%s took: %v", name, elapsed)
		return elapsed
	}

	// Run without logging
	normalTime := runTest("Without file logging", false)
	
	// Run with logging
	loggingTime := runTest("With file logging", true)

	// Calculate overhead
	overhead := loggingTime - normalTime
	overheadPercent := float64(overhead) / float64(normalTime) * 100

	t.Logf("File logging overhead: %v (%.1f%%)", overhead, overheadPercent)
	
	// Overhead should be minimal (less than 50%)
	if overheadPercent > 50 {
		t.Logf("WARNING: File logging adds significant overhead: %.1f%%", overheadPercent)
	} else {
		t.Logf("✅ File logging overhead is acceptable: %.1f%%", overheadPercent)
	}
}