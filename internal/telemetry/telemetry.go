package telemetry

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudshipai/ship/internal/config"
	"github.com/posthog/posthog-go"
)

const (
	// PostHog configuration from your project
	defaultAPIKey  = "phc_3h5yqMKKJsnxofspAcEouPsmJkbm2UfX0DGuhaa19f"
	defaultHost    = "https://us-i.posthog.com"
	defaultTimeout = 3 * time.Second
)

type Client struct {
	posthog    posthog.Client
	enabled    bool
	anonymousID string
}

var globalClient *Client

// Init initializes the global telemetry client
func Init() error {
	// Check if telemetry is disabled via environment variable
	if strings.ToLower(os.Getenv("SHIP_TELEMETRY")) == "false" {
		globalClient = &Client{enabled: false}
		return nil
	}

	// Load config to check telemetry settings
	cfg, err := config.Load()
	if err != nil {
		// If config fails to load, continue with default settings
		cfg = &config.Config{
			Telemetry: config.TelemetryConfig{
				Enabled: true, // Default to enabled unless explicitly disabled
			},
		}
	}

	// Check if telemetry is disabled in config
	if !cfg.Telemetry.Enabled {
		globalClient = &Client{enabled: false}
		return nil
	}

	// Create PostHog client
	client, err := posthog.NewWithConfig(
		defaultAPIKey,
		posthog.Config{
			Endpoint: defaultHost,
		},
	)
	if err != nil {
		// If PostHog client creation fails, disable telemetry
		globalClient = &Client{enabled: false}
		return nil
	}

	// Generate anonymous ID based on machine characteristics
	anonymousID := generateAnonymousID()

	globalClient = &Client{
		posthog:     client,
		enabled:     true,
		anonymousID: anonymousID,
	}

	return nil
}

// Close closes the telemetry client
func Close() {
	if globalClient != nil && globalClient.enabled && globalClient.posthog != nil {
		globalClient.posthog.Close()
	}
}

// TrackMCPCommand tracks when a user runs an MCP command
func TrackMCPCommand(toolName string) {
	if globalClient == nil || !globalClient.enabled {
		return
	}

	// Track asynchronously to not block the command
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create properties
		properties := posthog.NewProperties().
			Set("tool_name", toolName).
			Set("command", "mcp").
			Set("version", getVersion())

		// Capture the event
		err := globalClient.posthog.Enqueue(posthog.Capture{
			DistinctId: globalClient.anonymousID,
			Event:      "shp_mcp_command_executed",
			Properties: properties,
		})

		if err != nil {
			// Silently ignore telemetry errors to not affect user experience
			return
		}

		// Wait for the event to be sent or timeout
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
			return
		}
	}()
}

// generateAnonymousID creates a stable anonymous ID based on machine characteristics
func generateAnonymousID() string {
	// Use hostname and other stable machine characteristics
	hostname, _ := os.Hostname()
	userHomeDir, _ := os.UserHomeDir()
	
	// Create a hash of stable machine characteristics
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("ship-cli-%s-%s", hostname, userHomeDir)))
	hash := hasher.Sum(nil)
	
	// Return first 16 chars of hex encoded hash
	return fmt.Sprintf("ship-%x", hash)[:21]
}

// getVersion returns the current version (can be set via build flags)
func getVersion() string {
	if version := os.Getenv("SHIP_VERSION"); version != "" {
		return version
	}
	return "dev"
}

// IsEnabled returns whether telemetry is enabled
func IsEnabled() bool {
	return globalClient != nil && globalClient.enabled
}