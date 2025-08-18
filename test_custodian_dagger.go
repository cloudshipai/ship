package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudshipai/ship/internal/dagger/modules"
	"dagger.io/dagger"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to Dagger
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		log.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Connected to Dagger")

	// Test the fixed Custodian module
	custodianModule := modules.NewCustodianModule(client)
	
	// Test GetVersion
	version, err := custodianModule.GetVersion(ctx)
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("✅ Custodian version: %s\n", version)

	// Test ValidatePolicy with our test policy
	validationResult, err := custodianModule.ValidatePolicy(ctx, "/Users/jaredwolff/projects/ship/test-policy.yml")
	if err != nil {
		log.Printf("Failed to validate policy: %v", err)
	} else {
		fmt.Printf("✅ Policy validation result:\n%s\n", validationResult)
	}
}