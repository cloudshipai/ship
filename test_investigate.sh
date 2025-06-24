#!/bin/bash

echo "Testing Ship CLI investigate command with mock data"
echo "==================================================="

# Test with a simple query that doesn't require credentials
./ship investigate --provider aws --env dev

echo ""
echo "Note: This is using hardcoded queries for testing."
echo "In production, queries would come from the Cloudship API."