package aws

import (
	"testing"

	"github.com/cloudshipai/ship/internal/tools/finops/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestAWSProvider_Name(t *testing.T) {
	config := AWSConfig{
		Region: "us-east-1",
	}
	
	provider, err := NewAWSProvider(config)
	assert.NoError(t, err)
	assert.Equal(t, "AWS", provider.Name())
}

func TestAWSProvider_Type(t *testing.T) {
	config := AWSConfig{
		Region: "us-east-1",
	}
	
	provider, err := NewAWSProvider(config)
	assert.NoError(t, err)
	assert.Equal(t, interfaces.VendorAWS, provider.Type())
}

func TestAWSProvider_GetCapabilities(t *testing.T) {
	config := AWSConfig{
		Region: "us-east-1",
	}
	
	provider, err := NewAWSProvider(config)
	assert.NoError(t, err)
	
	capabilities := provider.GetCapabilities()
	assert.True(t, capabilities.SupportsRecommendations)
	assert.True(t, capabilities.SupportsCostForecasting)
	assert.True(t, capabilities.SupportsRightsizing)
	assert.Contains(t, capabilities.ResourceTypes, "ec2-instance")
	assert.Contains(t, capabilities.Regions, "us-east-1")
}

func TestAWSProvider_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  AWSConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: AWSConfig{
				Region: "us-east-1",
			},
			wantErr: false,
		},
		{
			name: "invalid region",
			config: AWSConfig{
				Region: "invalid-region",
			},
			wantErr: true,
		},
		{
			name: "missing region",
			config: AWSConfig{
				Region: "",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &AWSProvider{config: tt.config}
			err := provider.ValidateConfig(map[string]interface{}{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetCredentials(t *testing.T) {
	config := AWSConfig{
		Region: "us-east-1",
	}
	
	provider, err := NewAWSProvider(config)
	assert.NoError(t, err)
	
	// Test valid credentials
	creds := AWSCredentials{
		AccessKeyID:     "AKIAEXAMPLE",
		SecretAccessKey: "secret",
		Profile:         "default",
	}
	
	err = provider.SetCredentials(creds)
	assert.NoError(t, err)
	
	// Test invalid credentials type
	err = provider.SetCredentials("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials type")
}