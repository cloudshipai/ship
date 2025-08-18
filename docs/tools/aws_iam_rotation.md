# AWS IAM Access Key Management

AWS IAM access key management tools for rotating, creating, and managing AWS IAM user access keys using the official AWS CLI.

## Overview

AWS IAM access key rotation is a critical security practice that involves regularly replacing access keys to reduce the risk of unauthorized access. This tool provides MCP functions that wrap the official AWS CLI commands for comprehensive access key lifecycle management.

## Available MCP Functions

### 1. `aws_iam_list_access_keys`
**Description**: List access keys for an IAM user

**Parameters**:
- `user_name` (required): IAM username to list access keys for
- `profile` (optional): AWS profile to use

**Example Usage**:
```bash
# List access keys for a specific user
aws_iam_list_access_keys(user_name="myuser")

# List access keys using a specific AWS profile
aws_iam_list_access_keys(user_name="myuser", profile="production")
```

### 2. `aws_iam_create_access_key`
**Description**: Create a new access key for an IAM user

**Parameters**:
- `user_name` (required): IAM username to create access key for
- `profile` (optional): AWS profile to use

**Example Usage**:
```bash
# Create new access key for user
aws_iam_create_access_key(user_name="myuser")

# Create access key using specific profile
aws_iam_create_access_key(user_name="myuser", profile="development")
```

### 3. `aws_iam_update_access_key`
**Description**: Update the status of an access key (Active/Inactive)

**Parameters**:
- `access_key_id` (required): Access key ID to update
- `status` (required): New status for the access key (Active or Inactive)
- `user_name` (required): IAM username that owns the access key
- `profile` (optional): AWS profile to use

**Example Usage**:
```bash
# Deactivate an access key
aws_iam_update_access_key(
  access_key_id="AKIAIOSFODNN7EXAMPLE",
  status="Inactive",
  user_name="myuser"
)

# Reactivate an access key
aws_iam_update_access_key(
  access_key_id="AKIAIOSFODNN7EXAMPLE",
  status="Active",
  user_name="myuser"
)
```

### 4. `aws_iam_delete_access_key`
**Description**: Delete an access key for an IAM user

**Parameters**:
- `access_key_id` (required): Access key ID to delete
- `user_name` (required): IAM username that owns the access key
- `profile` (optional): AWS profile to use

**Example Usage**:
```bash
# Delete an access key (irreversible!)
aws_iam_delete_access_key(
  access_key_id="AKIAIOSFODNN7EXAMPLE",
  user_name="myuser"
)
```

### 5. `aws_iam_get_access_key_last_used`
**Description**: Get information about when an access key was last used

**Parameters**:
- `access_key_id` (required): Access key ID to check
- `profile` (optional): AWS profile to use

**Example Usage**:
```bash
# Check when access key was last used
aws_iam_get_access_key_last_used(access_key_id="AKIAIOSFODNN7EXAMPLE")
```

### 6. `aws_iam_get_version`
**Description**: Get AWS CLI version information

**Parameters**: None

**Example Usage**:
```bash
aws_iam_get_version()
```

## Access Key Rotation Process

The recommended process for rotating access keys safely:

### Step 1: List Current Keys
```bash
aws_iam_list_access_keys(user_name="myuser")
```

### Step 2: Create New Access Key
```bash
aws_iam_create_access_key(user_name="myuser")
```
**Note**: Save the returned SecretAccessKey immediately - it cannot be retrieved later.

### Step 3: Update Applications
Update your applications, scripts, and configuration files to use the new access key.

### Step 4: Test New Key
Verify that applications work correctly with the new access key.

### Step 5: Deactivate Old Key
```bash
aws_iam_update_access_key(
  access_key_id="OLD_ACCESS_KEY_ID",
  status="Inactive",
  user_name="myuser"
)
```

### Step 6: Monitor and Test
Monitor applications to ensure they're working with the new key.

### Step 7: Delete Old Key
```bash
aws_iam_delete_access_key(
  access_key_id="OLD_ACCESS_KEY_ID",
  user_name="myuser"
)
```

## Real CLI Capabilities

All MCP functions are based on the official AWS CLI commands:

### List Access Keys
```bash
aws iam list-access-keys --user-name myuser
```

### Create Access Key
```bash
aws iam create-access-key --user-name myuser
```

### Update Access Key Status
```bash
aws iam update-access-key --access-key-id AKIAIOSFODNN7EXAMPLE --status Inactive --user-name myuser
```

### Delete Access Key
```bash
aws iam delete-access-key --access-key-id AKIAIOSFODNN7EXAMPLE --user-name myuser
```

### Get Access Key Last Used
```bash
aws iam get-access-key-last-used --access-key-id AKIAIOSFODNN7EXAMPLE
```

### Get AWS CLI Version
```bash
aws --version
```

## Prerequisites

### AWS CLI Installation
```bash
# Install AWS CLI v2 (recommended)
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Verify installation
aws --version
```

### AWS Configuration
```bash
# Configure AWS credentials
aws configure

# Or set environment variables
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1
```

### IAM Permissions Required
The IAM user or role performing these operations needs these permissions:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListAccessKeys",
                "iam:CreateAccessKey",
                "iam:UpdateAccessKey",
                "iam:DeleteAccessKey",
                "iam:GetAccessKeyLastUsed"
            ],
            "Resource": "*"
        }
    ]
}
```

## Best Practices

### Security
- **Rotate Regularly**: Rotate access keys every 30-90 days
- **Monitor Usage**: Use `get-access-key-last-used` to track key usage
- **Least Privilege**: Grant minimum necessary permissions
- **Never Commit Keys**: Don't store access keys in code repositories

### Operational
- **Test Thoroughly**: Always test applications after key rotation
- **Document Process**: Maintain runbooks for key rotation procedures
- **Automate When Possible**: Use automation tools for regular rotation
- **Monitor Failures**: Set up alerts for failed authentication attempts

### Access Key Limits
- **Maximum Keys**: Each IAM user can have a maximum of 2 access keys
- **This Enables Rotation**: You can create a new key before deleting the old one
- **Zero Downtime**: Proper rotation allows for zero-downtime key updates

## Third-Party Rotation Tools

Several open-source tools can automate AWS access key rotation:

### aws-rotate-iam-keys
Simple CLI tool for rotating IAM access keys:
```bash
# Install via homebrew
brew install aws-rotate-iam-keys

# Rotate default profile
aws-rotate-iam-keys

# Rotate specific profile
aws-rotate-iam-keys --profile myprofile
```

### AWS Samples Auto-Rotation
CloudFormation-based solution for automated rotation:
- Automatically rotates keys every 90 days
- Stores new keys in AWS Secrets Manager
- Sends email notifications via SES
- Repository: https://github.com/aws-samples/aws-iam-access-key-auto-rotation

## Troubleshooting

### Common Issues

1. **Permission Denied**
   - Verify IAM permissions for the user/role
   - Check if MFA is required for IAM operations

2. **Maximum Access Keys Reached**
   - Delete or deactivate an existing access key first
   - Each user can have maximum 2 access keys

3. **Access Key Not Found**
   - Verify the access key ID is correct
   - Ensure the key belongs to the specified user

4. **Profile Not Found**
   - Check AWS CLI configuration: `aws configure list-profiles`
   - Verify profile configuration: `aws configure list --profile myprofile`

### Error Codes
- **AccessDenied**: Insufficient IAM permissions
- **LimitExceeded**: Maximum number of access keys reached
- **NoSuchEntity**: Access key or user not found
- **InvalidUserType**: Operation not supported for user type

## Integration with Ship CLI

These MCP functions integrate with Ship CLI's containerized execution:
- Commands are executed through the Ship CLI's Dagger engine
- AWS CLI runs in a containerized environment
- Credentials can be passed via environment variables or mounted volumes

## References

- **AWS IAM User Guide**: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
- **AWS CLI IAM Commands**: https://docs.aws.amazon.com/cli/latest/reference/iam/
- **AWS Security Blog - Key Rotation**: https://aws.amazon.com/blogs/security/how-to-rotate-access-keys-for-iam-users/
- **AWS Best Practices**: https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html