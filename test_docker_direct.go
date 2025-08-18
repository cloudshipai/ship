package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// Test version directly with Docker
	cmd := exec.Command("docker", "run", "--rm", "cloudcustodian/c7n:latest", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running version: %v\n", err)
	}
	fmt.Printf("Version output: %s\n", output)

	// Test validate directly with Docker
	cmd2 := exec.Command("docker", "run", "--rm", "-v", "/Users/jaredwolff/projects/ship/test-policy.yml:/policy.yml", "cloudcustodian/c7n:latest", "validate", "/policy.yml")
	output2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		fmt.Printf("Error running validate: %v\n", err2)
	}
	fmt.Printf("Validate output: %s\n", output2)
}