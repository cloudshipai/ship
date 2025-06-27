#!/bin/bash

# Script to run AI investigation with proper AWS credentials

PROFILE=${AWS_PROFILE:-default}
echo "Using AWS profile: $PROFILE"

# Export AWS credentials from profile
eval $(aws configure export-credentials --profile $PROFILE --format env 2>/dev/null || echo "")

if [ -z "$AWS_ACCESS_KEY_ID" ]; then
    echo "Trying to read credentials directly..."
    AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id --profile $PROFILE 2>/dev/null)
    AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key --profile $PROFILE 2>/dev/null)
    AWS_SESSION_TOKEN=$(aws configure get aws_session_token --profile $PROFILE 2>/dev/null)
    AWS_REGION=$(aws configure get region --profile $PROFILE 2>/dev/null || echo "us-east-1")
fi

if [ -n "$AWS_ACCESS_KEY_ID" ]; then
    echo "✓ AWS credentials loaded successfully"
    echo "  Access Key: ${AWS_ACCESS_KEY_ID:0:10}..."
    echo "  Region: $AWS_REGION"
    
    # Export for child processes
    export AWS_ACCESS_KEY_ID
    export AWS_SECRET_ACCESS_KEY
    export AWS_SESSION_TOKEN
    export AWS_REGION
    
    # Run the command
    exec ./ship ai-investigate "$@"
else
    echo "❌ Failed to load AWS credentials for profile: $PROFILE"
    echo ""
    echo "Available profiles:"
    aws configure list-profiles
    echo ""
    echo "Usage: AWS_PROFILE=your-profile $0 --prompt \"your query\" --execute"
    exit 1
fi