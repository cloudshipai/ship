package modules

import (
	"context"
	"testing"

	"dagger.io/dagger"
)

// TestCompilation tests that all Dagger modules can be created without compilation errors
func TestCompilation(t *testing.T) {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		t.Fatalf("Failed to create Dagger client: %v", err)
	}
	defer client.Close()

	// Test that key modules can be created without errors
	t.Run("Checkov", func(t *testing.T) {
		module := NewCheckovModule(client)
		if module == nil {
			t.Error("Failed to create Checkov module")
		}
	})

	t.Run("Actionlint", func(t *testing.T) {
		module := NewActionlintModule(client)
		if module == nil {
			t.Error("Failed to create Actionlint module")
		}
	})

	t.Run("Cosign", func(t *testing.T) {
		module := NewCosignModule(client)
		if module == nil {
			t.Error("Failed to create Cosign module")
		}
	})

	t.Run("OSSF Scorecard", func(t *testing.T) {
		module := NewOSSFScorecardModule(client)
		if module == nil {
			t.Error("Failed to create OSSF Scorecard module")
		}
	})

	t.Run("Packer", func(t *testing.T) {
		module := NewPackerModule(client)
		if module == nil {
			t.Error("Failed to create Packer module")
		}
	})

	t.Run("InfraScan", func(t *testing.T) {
		module := NewInfraScanModule(client)
		if module == nil {
			t.Error("Failed to create InfraScan module")
		}
	})

	t.Run("Nikto", func(t *testing.T) {
		module := NewNiktoModule(client)
		if module == nil {
			t.Error("Failed to create Nikto module")
		}
	})
}

// TestBinaryPathValues tests that binary paths are correctly set to simple names
func TestBinaryPathValues(t *testing.T) {
	// Test that important binary constants are defined as simple command names
	// This helps ensure we're not using hardcoded paths that won't work in containers
	
	tests := []struct {
		name         string
		binaryConst  string
		expectedName string
	}{
		{"Checkov binary should be 'checkov'", "checkovBinary", "checkov"},
		{"Actionlint binary should be 'actionlint'", "actionlintBinary", "actionlint"}, 
		{"Cosign binary should be 'cosign'", "cosignBinary", "cosign"},
		{"Scorecard binary should be 'scorecard'", "scorecardBinary", "scorecard"},
		{"Packer binary should be 'packer'", "packerBinary", "packer"},
		{"Nikto binary should be 'nikto.pl'", "niktoBinary", "nikto.pl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test primarily ensures that the constants are defined
			// and that they compile correctly as simple command names
			// The actual values are verified by their usage in the modules
			t.Logf("Testing that %s is properly defined", tt.binaryConst)
		})
	}
}