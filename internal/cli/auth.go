package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudshipai/ship/internal/cloudship"
	"github.com/cloudshipai/ship/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	authToken string
	logout    bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with CloudShip",
	Long:  `Authenticate with CloudShip using an API key from your CloudShip settings page.
	
Get your API key from: https://app.cloudshipai.com/settings/api-keys

Example:
  ship auth --api-key your-api-key-here
  
You can also set the API key via environment variable:
  export CLOUDSHIP_API_KEY=your-api-key-here`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if logout {
			return handleLogout()
		}

		// Check for API key from flag or environment
		if authToken == "" {
			authToken = os.Getenv("CLOUDSHIP_API_KEY")
		}

		if authToken == "" {
			return fmt.Errorf("--api-key flag is required or set CLOUDSHIP_API_KEY environment variable")
		}

		return handleAuth(authToken)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVar(&authToken, "api-key", "", "CloudShip API key")
	authCmd.Flags().BoolVar(&logout, "logout", false, "Log out and clear stored credentials")
	
	// Keep old flag for backwards compatibility
	authCmd.Flags().StringVar(&authToken, "token", "", "API token (deprecated, use --api-key)")
	authCmd.Flags().MarkDeprecated("token", "use --api-key instead")
}

func handleAuth(token string) error {
	// Load existing config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate API key by attempting a simple API call
	fmt.Println("Validating API key...")
	client := cloudship.NewClient(token)
	
	// Try to list artifacts with a dummy fleet ID to validate the API key
	// The API will return 401 if the key is invalid
	testReq := &cloudship.ListArtifactsRequest{
		FleetID: "test",
		Limit:   1,
	}
	
	_, err = client.ListArtifacts(testReq)
	if err != nil {
		// Check if it's an auth error
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "unauthorized") {
			return fmt.Errorf("invalid API key - please check your key and try again")
		}
		// Other errors might be OK (e.g., fleet not found)
		// For now, we'll accept the key if it's not explicitly unauthorized
	}

	// Update config with API key
	cfg.APIKey = token
	cfg.BaseURL = os.Getenv("CLOUDSHIP_API_URL")
	if cfg.BaseURL == "" {
		cfg.BaseURL = cloudship.DefaultAPIURL
	}

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Success message
	green := color.New(color.FgGreen)
	green.Printf("✓ Successfully authenticated with CloudShip!\n")
	fmt.Printf("Configuration saved to: %s\n", config.GetConfigPath())

	// Show masked token for confirmation
	maskedToken := maskToken(token)
	fmt.Printf("API Key: %s\n", maskedToken)
	fmt.Printf("API URL: %s\n", cfg.BaseURL)

	return nil
}

func handleLogout() error {
	// Load config to get token for logout API call
	_, err := config.Load()
	if err != nil {
		// If config doesn't exist, that's fine
		if err := config.Clear(); err != nil {
			return fmt.Errorf("failed to clear config: %w", err)
		}
		yellow := color.New(color.FgYellow)
		yellow.Println("✓ Already logged out")
		return nil
	}

	// For now, skip API logout since we don't have a real endpoint
	// In production, uncomment this:
	/*
		if cfg.Token != "" {
			authClient := auth.NewClient(cfg.BaseURL)
			if err := authClient.Logout(cfg.Token); err != nil {
				// Log error but continue with local cleanup
				fmt.Printf("Warning: failed to logout from server: %v\n", err)
			}
		}
	*/

	// Clear local config
	if err := config.Clear(); err != nil {
		return fmt.Errorf("failed to clear config: %w", err)
	}

	green := color.New(color.FgGreen)
	green.Println("✓ Successfully logged out")
	fmt.Println("Cleared configuration and cache")

	return nil
}

func maskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}

	// Show first 7 chars and last 3 chars
	prefix := token[:7]
	suffix := token[len(token)-3:]
	masked := strings.Repeat("*", len(token)-10)

	return fmt.Sprintf("%s%s%s", prefix, masked, suffix)
}
