#!/bin/bash

# Export AWS credentials from profile
eval $(aws configure export-credentials --profile default --format env 2>/dev/null || echo "")

if [ -z "$AWS_ACCESS_KEY_ID" ]; then
    echo "Failed to export AWS credentials. Trying different approach..."
    # Try to read directly from credentials file
    AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id --profile default 2>/dev/null)
    AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key --profile default 2>/dev/null)
    AWS_SESSION_TOKEN=$(aws configure get aws_session_token --profile default 2>/dev/null)
fi

# Test if we have credentials
if [ -n "$AWS_ACCESS_KEY_ID" ]; then
    echo "AWS credentials found!"
    echo "AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID:0:10}..."
    
    # Run ship with environment variables
    AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
    AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
    AWS_SESSION_TOKEN="$AWS_SESSION_TOKEN" \
    ./ship steampipe-test --provider aws
else
    echo "No AWS credentials found for default profile"
    echo "Available profiles:"
    aws configure list-profiles
fi