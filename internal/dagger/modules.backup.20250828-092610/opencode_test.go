package modules

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"dagger.io/dagger"
)

func TestOpenCodeModule_GetVersion(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	version, err := module.GetVersion(ctx)
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if version == "" {
		t.Error("Expected non-empty version string")
	}

	// Version should contain some expected patterns
	if !strings.Contains(strings.ToLower(version), "opencode") && 
	   !strings.Contains(version, ".") {
		t.Errorf("Version string doesn't look like a valid version: %s", version)
	}
}

func TestOpenCodeModule_Chat(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "hello.js")
	err = os.WriteFile(testFile, []byte(`console.log("Hello World!");`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.Chat(ctx, tempDir, "What does this code do?")
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty chat response")
	}
}

func TestOpenCodeModule_Generate(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	tempDir := t.TempDir()
	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.Generate(ctx, tempDir, "Create a simple hello world function in JavaScript", "hello.js")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty generation response")
	}
}

func TestOpenCodeModule_AnalyzeFile(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.py")
	testCode := `def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(10))`
	
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.AnalyzeFile(ctx, testFile, "What is the time complexity of this code?")
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty analysis response")
	}
}

func TestOpenCodeModule_Review(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a temporary directory with some code
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "review.js")
	testCode := `function add(a, b) {
    return a + b;
}

var result = add(5, 3);
console.log(result);`
	
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.Review(ctx, tempDir, "review.js")
	if err != nil {
		t.Fatalf("Review failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty review response")
	}
}

func TestOpenCodeModule_WithAuth(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	tempDir := t.TempDir()
	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.WithAuth(ctx, tempDir, "OPENAI", "test-key-123")
	if err != nil {
		t.Fatalf("WithAuth failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty auth response")
	}
}

func TestOpenCodeModule_Test(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a simple JavaScript project
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "math.js")
	testCode := `function multiply(a, b) {
    return a * b;
}

module.exports = { multiply };`
	
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.Test(ctx, tempDir, "unit", false)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty test response")
	}
}

func TestOpenCodeModule_Document(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create a simple code file to document
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "utils.py")
	testCode := `def calculate_area(length, width):
    """Calculate the area of a rectangle."""
    return length * width

def calculate_perimeter(length, width):
    """Calculate the perimeter of a rectangle."""
    return 2 * (length + width)`
	
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	output, err := module.Document(ctx, tempDir, "markdown", "docs")
	if err != nil {
		t.Fatalf("Document failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty documentation response")
	}
}

func TestOpenCodeModule_BatchProcess(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	// Create multiple test files
	tempDir := t.TempDir()
	
	files := map[string]string{
		"file1.js": `function hello() { console.log("Hello"); }`,
		"file2.js": `function world() { console.log("World"); }`,
		"file3.js": `function greet() { console.log("Greetings"); }`,
	}

	for filename, content := range files {
		err = os.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	module := NewOpenCodeModule(client)
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	output, err := module.BatchProcess(ctx, tempDir, "*.js", "analyze")
	if err != nil {
		t.Fatalf("BatchProcess failed: %v", err)
	}

	if output == "" {
		t.Error("Expected non-empty batch process response")
	}
}

func TestNewOpenCodeModule(t *testing.T) {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Skipf("Dagger not available: %v", err)
	}
	defer client.Close()

	module := NewOpenCodeModule(client)
	
	if module == nil {
		t.Fatal("NewOpenCodeModule returned nil")
	}
	
	if module.client != client {
		t.Error("Client not properly assigned")
	}
	
	if module.name != "opencode" {
		t.Errorf("Expected name 'opencode', got '%s'", module.name)
	}
}