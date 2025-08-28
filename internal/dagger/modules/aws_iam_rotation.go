package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// AWSIAMRotationModule manages AWS IAM credential rotation
type AWSIAMRotationModule struct {
	client *dagger.Client
	name   string
}

const awsIamBinary = "aws"

// NewAWSIAMRotationModule creates a new AWS IAM rotation module
func NewAWSIAMRotationModule(client *dagger.Client) *AWSIAMRotationModule {
	return &AWSIAMRotationModule{
		client: client,
		name:   "aws-iam-rotation",
	}
}

// RotateAccessKeys rotates AWS access keys for a user
func (m *AWSIAMRotationModule) RotateAccessKeys(ctx context.Context, username string, profile string) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{
		awsIamBinary, "iam", "create-access-key",
		"--user-name", username,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to rotate access keys: %w", err)
	}

	return output, nil
}

// ListAccessKeys lists access keys for a user
func (m *AWSIAMRotationModule) ListAccessKeys(ctx context.Context, username string, profile string) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{
		awsIamBinary, "iam", "list-access-keys",
		"--user-name", username,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list access keys: %w", err)
	}

	return output, nil
}

// DeleteAccessKey deletes an access key
func (m *AWSIAMRotationModule) DeleteAccessKey(ctx context.Context, username string, accessKeyId string, profile string) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{
		awsIamBinary, "iam", "delete-access-key",
		"--user-name", username,
		"--access-key-id", accessKeyId,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to delete access key: %w", err)
	}

	return output, nil
}

// UpdateAccessKey updates access key status
func (m *AWSIAMRotationModule) UpdateAccessKey(ctx context.Context, username string, accessKeyId string, status string, profile string) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{
		"aws", "iam", "update-access-key",
		"--user-name", username,
		"--access-key-id", accessKeyId,
		"--status", status,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to update access key: %w", err)
	}

	return output, nil
}

// GetAccessKeyLastUsed gets access key last used info
func (m *AWSIAMRotationModule) GetAccessKeyLastUsed(ctx context.Context, accessKeyId string, profile string) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{
		"aws", "iam", "get-access-key-last-used",
		"--access-key-id", accessKeyId,
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get access key last used: %w", err)
	}

	return output, nil
}

// GetVersion returns the AWS CLI version
func (m *AWSIAMRotationModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("amazon/aws-cli:latest").
		WithExec([]string{"aws", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get AWS CLI version: %w", err)
	}

	return output, nil
}
