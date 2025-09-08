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
	// PostHog configuration
	defaultAPIKey   = "phc_3h5yqMKKJsnxofsqFxCEoUFmn3vbm2UFXDDKuhdai9f"
	defaultHost     = "https://us-i.posthog.com"
	defaultTimeout  = 3 * time.Second
	defaultProjectID = "203190"
)

type Client struct {
	posthog     posthog.Client
	enabled     bool
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
	trackEvent("shp_mcp_command_executed", posthog.NewProperties().
		Set("tool_name", toolName).
		Set("command", "mcp").
		Set("version", getVersion()))
}

// TrackCLICommand tracks when a user runs a CLI command
func TrackCLICommand(commandName string, subcommand string, args []string) {
	props := posthog.NewProperties().
		Set("command", commandName).
		Set("version", getVersion())
	
	if subcommand != "" {
		props.Set("subcommand", subcommand)
	}
	
	if len(args) > 0 {
		props.Set("arg_count", len(args))
	}
	
	trackEvent("shp_cli_command_executed", props)
}

// TrackBuildXOperation tracks BuildX operations
func TrackBuildXOperation(operation string, platform string, success bool) {
	trackEvent("shp_buildx_operation", posthog.NewProperties().
		Set("operation", operation).
		Set("platform", platform).
		Set("success", success).
		Set("version", getVersion()))
}

// TrackToolExecution tracks when a security/infrastructure tool is executed
func TrackToolExecution(toolName string, executionTime time.Duration, success bool, errorType string) {
	props := posthog.NewProperties().
		Set("tool_name", toolName).
		Set("execution_time_ms", executionTime.Milliseconds()).
		Set("success", success).
		Set("version", getVersion())
	
	if !success && errorType != "" {
		props.Set("error_type", errorType)
	}
	
	trackEvent("shp_tool_execution", props)
}

// TrackDaggerOperation tracks Dagger-related operations
func TrackDaggerOperation(operation string, module string, success bool, executionTime time.Duration) {
	trackEvent("shp_dagger_operation", posthog.NewProperties().
		Set("operation", operation).
		Set("module", module).
		Set("success", success).
		Set("execution_time_ms", executionTime.Milliseconds()).
		Set("version", getVersion()))
}

// TrackError tracks significant errors
func TrackError(errorType string, component string, errorMessage string) {
	trackEvent("shp_error_occurred", posthog.NewProperties().
		Set("error_type", errorType).
		Set("component", component).
		Set("error_message_hash", hashString(errorMessage)). // Hash for privacy
		Set("version", getVersion()))
}

// TrackAppStart tracks when the application starts
func TrackAppStart(command string) {
	trackEvent("shp_app_started", posthog.NewProperties().
		Set("entry_command", command).
		Set("version", getVersion()))
}

// TrackConfigOperation tracks config-related operations
func TrackConfigOperation(operation string, success bool) {
	trackEvent("shp_config_operation", posthog.NewProperties().
		Set("operation", operation).
		Set("success", success).
		Set("version", getVersion()))
}

// trackEvent is a helper function for sending events asynchronously
func trackEvent(eventName string, properties posthog.Properties) {
	if globalClient == nil || !globalClient.enabled {
		return
	}

	// Track asynchronously to not block the command
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Capture the event
		err := globalClient.posthog.Enqueue(posthog.Capture{
			DistinctId: globalClient.anonymousID,
			Event:      eventName,
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

// hashString creates a hash of a string for privacy (used for error messages)
func hashString(input string) string {
	if input == "" {
		return ""
	}
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)[:16] // First 16 chars of hash
}

// IsEnabled returns whether telemetry is enabled
func IsEnabled() bool {
	return globalClient != nil && globalClient.enabled
}
