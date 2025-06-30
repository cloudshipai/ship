package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Token      string          `mapstructure:"token"`      // Deprecated, use APIKey
	APIKey     string          `mapstructure:"api_key"`
	OrgID      string          `mapstructure:"org_id"`
	DefaultEnv string          `mapstructure:"default_env"`
	BaseURL    string          `mapstructure:"base_url"`
	FleetID    string          `mapstructure:"fleet_id"`    // Default fleet ID for push operations
	Telemetry  TelemetryConfig `mapstructure:"telemetry"`
}

type TelemetryConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	SessionID string `mapstructure:"session_id"`
}

const (
	defaultBaseURL = "https://api.cloudship.ai/v1"
	configFileName = "config"
	configFileType = "yaml"
)

var (
	configDir  string
	configPath string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %w", err))
	}

	configDir = filepath.Join(homeDir, ".ship")
	configPath = filepath.Join(configDir, "config.yaml")
}

func GetConfigDir() string {
	return configDir
}

func GetConfigPath() string {
	return configPath
}

func Load() (*Config, error) {
	// Check environment variables first
	cfg := &Config{
		BaseURL: getBaseURL(),
	}
	
	// Check for API key in environment
	if apiKey := os.Getenv("CLOUDSHIP_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	} else if token := os.Getenv("SHIP_TOKEN"); token != "" {
		// Fallback to old env var for backwards compatibility
		cfg.APIKey = token
	}
	
	// Check for fleet ID in environment
	if fleetID := os.Getenv("CLOUDSHIP_FLEET_ID"); fleetID != "" {
		cfg.FleetID = fleetID
	}
	
	// If we have env vars, return early
	if cfg.APIKey != "" {
		return cfg, nil
	}

	// Create a new viper instance for loading
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(configDir)

	// Set defaults
	v.SetDefault("base_url", defaultBaseURL)
	v.SetDefault("telemetry.enabled", false)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; return empty config
			return &Config{
				BaseURL: defaultBaseURL,
			}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var fileCfg Config
	if err := v.Unmarshal(&fileCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Merge file config with env config (env takes precedence)
	if cfg.APIKey == "" {
		cfg.APIKey = fileCfg.APIKey
	}
	if cfg.FleetID == "" {
		cfg.FleetID = fileCfg.FleetID
	}
	if cfg.OrgID == "" {
		cfg.OrgID = fileCfg.OrgID
	}
	if cfg.DefaultEnv == "" {
		cfg.DefaultEnv = fileCfg.DefaultEnv
	}
	cfg.Telemetry = fileCfg.Telemetry

	return cfg, nil
}

func Save(cfg *Config) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create a new viper instance for saving
	v := viper.New()
	v.SetConfigType(configFileType)

	// Set all config values
	v.Set("token", cfg.Token)
	v.Set("api_key", cfg.APIKey)
	v.Set("org_id", cfg.OrgID)
	v.Set("default_env", cfg.DefaultEnv)
	v.Set("base_url", cfg.BaseURL)
	v.Set("fleet_id", cfg.FleetID)
	v.Set("telemetry.enabled", cfg.Telemetry.Enabled)
	v.Set("telemetry.session_id", cfg.Telemetry.SessionID)

	// Write config file
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Set secure permissions (owner read/write only)
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set config permissions: %w", err)
	}

	return nil
}

func Clear() error {
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	// Also remove cache directory if it exists
	cacheDir := filepath.Join(configDir, "cache")
	if err := os.RemoveAll(cacheDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}

	return nil
}

func getBaseURL() string {
	if url := os.Getenv("CLOUDSHIP_API_URL"); url != "" {
		return url
	} else if url := os.Getenv("SHIP_API_URL"); url != "" {
		// Fallback to old env var for backwards compatibility
		return url
	}
	return defaultBaseURL
}
