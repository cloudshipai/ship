package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func newTestViper(tempDir string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileType)
	v.AddConfigPath(tempDir)
	v.SetDefault("base_url", defaultBaseURL)
	v.SetDefault("telemetry.enabled", false)
	v.SetEnvPrefix("SHIP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.BindEnv("api_key", "CLOUDSHIP_API_KEY")
	v.BindEnv("base_url", "SHIP_API_URL", "CLOUDSHIP_API_URL")
	v.BindEnv("fleet_id", "CLOUDSHIP_FLEET_ID")
	return v
}

func TestConfig(t *testing.T) {
	// Create temp directory for test config
	tempDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath

	// Override config paths for testing
	configDir = tempDir
	configPath = filepath.Join(tempDir, "config.yaml")

	// Restore original paths after test
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	t.Run("Load empty config", func(t *testing.T) {
		v := newTestViper(tempDir)
		cfg, err := loadViper(v)
		if err != nil {
			t.Fatalf("Failed to load empty config: %v", err)
		}

		if cfg.BaseURL != defaultBaseURL {
			t.Errorf("Expected base URL %s, got %s", defaultBaseURL, cfg.BaseURL)
		}

		if cfg.APIKey != "" {
			t.Errorf("Expected empty token, got %s", cfg.APIKey)
		}
	})

	t.Run("Save and load config", func(t *testing.T) {
		v := newTestViper(tempDir)
		testCfg := &Config{
			APIKey:     "test-token",
			OrgID:      "org-123",
			DefaultEnv: "prod",
			BaseURL:    "https://test.api.com",
			Telemetry: TelemetryConfig{
				Enabled:   true,
				SessionID: "session-123",
			},
		}

		// Save config
		if err := saveViper(v, testCfg); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Check file permissions
		info, err := os.Stat(configPath)
		if err != nil {
			t.Fatalf("Failed to stat config file: %v", err)
		}

		if info.Mode().Perm() != 0600 {
			t.Errorf("Expected permissions 0600, got %v", info.Mode().Perm())
		}

		// Debug: Check what was actually saved
		savedContent, _ := os.ReadFile(configPath)
		t.Logf("Saved config content: %s", savedContent)

		// Load config
		loadedCfg, err := loadViper(v)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify loaded config
		if loadedCfg.APIKey != testCfg.APIKey {
			t.Errorf("Expected token %s, got %s", testCfg.APIKey, loadedCfg.APIKey)
		}

		if loadedCfg.OrgID != testCfg.OrgID {
			t.Errorf("Expected org ID %s, got %s", testCfg.OrgID, loadedCfg.OrgID)
		}

		if loadedCfg.DefaultEnv != testCfg.DefaultEnv {
			t.Errorf("Expected default env %s, got %s", testCfg.DefaultEnv, loadedCfg.DefaultEnv)
		}

		if loadedCfg.BaseURL != testCfg.BaseURL {
			t.Errorf("Expected base URL %s, got %s", testCfg.BaseURL, loadedCfg.BaseURL)
		}

		if loadedCfg.Telemetry.Enabled != testCfg.Telemetry.Enabled {
			t.Errorf("Expected telemetry enabled %v, got %v", testCfg.Telemetry.Enabled, loadedCfg.Telemetry.Enabled)
		}
	})

	t.Run("Environment variable override", func(t *testing.T) {
		v := newTestViper(tempDir)
		// Set environment variable
		os.Setenv("SHIP_API_KEY", "env-token")
		defer os.Unsetenv("SHIP_API_KEY")

		cfg, err := loadViper(v)
		if err != nil {
			t.Fatalf("Failed to load config with env var: %v", err)
		}

		if cfg.APIKey != "env-token" {
			t.Errorf("Expected token from env var, got %s", cfg.APIKey)
		}
	})

	t.Run("Clear config", func(t *testing.T) {
		v := newTestViper(tempDir)
		// Create a config file
		testCfg := &Config{APIKey: "test"}
		if err := saveViper(v, testCfg); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		// Clear config
		if err := Clear(); err != nil {
			t.Fatalf("Failed to clear config: %v", err)
		}

		// Verify file is removed
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			t.Error("Config file should be removed")
		}
	})
}
