package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey     string          `mapstructure:"api_key"`
	OrgID      string          `mapstructure:"org_id"`
	DefaultEnv string          `mapstructure:"default_env"`
	BaseURL    string          `mapstructure:"base_url"`
	FleetID    string          `mapstructure:"fleet_id"`
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
	v          *viper.Viper
	configDir  string
	configPath string
)

func init() {
	v = viper.New()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %w", err))
	}

	configDir = filepath.Join(homeDir, ".ship")
	configPath = filepath.Join(configDir, "config.yaml")

	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(configDir)

	// Set defaults
	v.SetDefault("base_url", defaultBaseURL)
	v.SetDefault("telemetry.enabled", false)

	// Environment variable binding
	v.SetEnvPrefix("SHIP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Backwards compatibility for env vars
	v.BindEnv("api_key", "CLOUDSHIP_API_KEY")
	v.BindEnv("base_url", "SHIP_API_URL", "CLOUDSHIP_API_URL")
	v.BindEnv("fleet_id", "CLOUDSHIP_FLEET_ID")
}

func GetConfigDir() string {
	return configDir
}

func GetConfigPath() string {
	return configPath
}

func Load() (*Config, error) {
	return loadViper(v)
}

func loadViper(v *viper.Viper) (*Config, error) {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Backwards compatibility for `token` field
	if token := v.GetString("token"); token != "" && cfg.APIKey == "" {
		cfg.APIKey = token
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	return saveViper(v, cfg)
}

func saveViper(v *viper.Viper, cfg *Config) error {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	v.Set("api_key", cfg.APIKey)
	v.Set("org_id", cfg.OrgID)
	v.Set("default_env", cfg.DefaultEnv)
	v.Set("base_url", cfg.BaseURL)
	v.Set("fleet_id", cfg.FleetID)
	v.Set("telemetry.enabled", cfg.Telemetry.Enabled)
	v.Set("telemetry.session_id", cfg.Telemetry.SessionID)

	// Remove deprecated `token` field
	v.Set("token", nil)

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set config permissions: %w", err)
	}

	return nil
}

func Clear() error {
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}
	cacheDir := filepath.Join(configDir, "cache")
	if err := os.RemoveAll(cacheDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}
	return nil
}

func getBaseURL() string {
	// This function is no longer needed as Viper handles it.
	// Kept for reference during refactoring.
	return ""
}
