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
	v.SetDefault("telemetry.enabled", false)
	v.SetEnvPrefix("SHIP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
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

		if cfg.DefaultEnv != "" {
			t.Errorf("Expected empty default env, got %s", cfg.DefaultEnv)
		}
	})

	t.Run("Save and load config", func(t *testing.T) {
		v := newTestViper(tempDir)
		testCfg := &Config{
			DefaultEnv: "prod",
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
		if loadedCfg.DefaultEnv != testCfg.DefaultEnv {
			t.Errorf("Expected default env %s, got %s", testCfg.DefaultEnv, loadedCfg.DefaultEnv)
		}

		if loadedCfg.Telemetry.Enabled != testCfg.Telemetry.Enabled {
			t.Errorf("Expected telemetry enabled %v, got %v", testCfg.Telemetry.Enabled, loadedCfg.Telemetry.Enabled)
		}
	})

	t.Run("Environment variable override", func(t *testing.T) {
		v := newTestViper(tempDir)
		// Set environment variable
		os.Setenv("SHIP_DEFAULT_ENV", "env-test")
		defer os.Unsetenv("SHIP_DEFAULT_ENV")

		cfg, err := loadViper(v)
		if err != nil {
			t.Fatalf("Failed to load config with env var: %v", err)
		}

		if cfg.DefaultEnv != "env-test" {
			t.Errorf("Expected default env from env var, got %s", cfg.DefaultEnv)
		}
	})

	t.Run("Clear config", func(t *testing.T) {
		v := newTestViper(tempDir)
		// Create a config file
		testCfg := &Config{DefaultEnv: "test"}
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
