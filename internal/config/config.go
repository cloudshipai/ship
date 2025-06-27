package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Token      string          `mapstructure:"token"`
	OrgID      string          `mapstructure:"org_id"`
	DefaultEnv string          `mapstructure:"default_env"`
	BaseURL    string          `mapstructure:"base_url"`
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
	// Check environment variable first
	if token := os.Getenv("SHIP_TOKEN"); token != "" {
		return &Config{
			Token:   token,
			BaseURL: getBaseURL(),
		}, nil
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

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
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
	v.Set("org_id", cfg.OrgID)
	v.Set("default_env", cfg.DefaultEnv)
	v.Set("base_url", cfg.BaseURL)
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
	if url := os.Getenv("SHIP_API_URL"); url != "" {
		return url
	}
	return defaultBaseURL
}
