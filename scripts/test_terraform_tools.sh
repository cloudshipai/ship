#!/bin/bash

echo "Testing Terraform Tools in Ship CLI"
echo "==================================="

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build the CLI first
echo -e "${BLUE}Building Ship CLI...${NC}"
go build -o ship ./cmd/ship

if [ $? -ne 0 ]; then
    echo "Failed to build Ship CLI"
    exit 1
fi

echo -e "${GREEN}✓ Ship CLI built successfully${NC}"
echo ""

# Test directory
TEST_DIR="test/terraform-mocks"

# Test 1: Generate documentation
echo -e "${BLUE}Test 1: Generating Terraform documentation${NC}"
echo "Command: ./ship terraform-tools generate-docs $TEST_DIR"
echo "---"
./ship terraform-tools generate-docs $TEST_DIR
echo ""

# Test 2: Security scan
echo -e "${BLUE}Test 2: Running security scan${NC}"
echo "Command: ./ship terraform-tools security-scan $TEST_DIR"
echo "---"
./ship terraform-tools security-scan $TEST_DIR
echo ""

# Test 3: Cost analysis with plan file
echo -e "${BLUE}Test 3: Running cost analysis on plan file${NC}"
echo "Command: ./ship terraform-tools cost-analysis $TEST_DIR/tfplan.json"
echo "---"
./ship terraform-tools cost-analysis $TEST_DIR/tfplan.json
echo ""

# Test 4: Cost analysis on directory
echo -e "${BLUE}Test 4: Running cost analysis on directory${NC}"
echo "Command: ./ship terraform-tools cost-analysis $TEST_DIR"
echo "---"
./ship terraform-tools cost-analysis $TEST_DIR
echo ""

echo -e "${GREEN}All tests completed!${NC}"
echo ""
echo "Summary:"
echo "- terraform-docs: Generates documentation from Terraform code"
echo "- InfraScan: Detects security issues in Terraform configurations"
echo "- OpenInfraQuote: Estimates costs for Terraform resources"