package cli

import (
	"fmt"
	"strings"

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
	Short: "Authenticate with Cloudship",
	Long:  `Authenticate with Cloudship using an API token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if logout {
			return handleLogout()
		}

		if authToken == "" {
			return fmt.Errorf("--token flag is required")
		}

		return handleAuth(authToken)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVar(&authToken, "token", "", "API token (required)")
	authCmd.Flags().BoolVar(&logout, "logout", false, "Log out and clear stored credentials")
}

func handleAuth(token string) error {
	// Load existing config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// For now, skip validation since we don't have a real API endpoint
	// In production, uncomment this:
	/*
		// Validate token
		fmt.Println("Validating token...")
		authClient := auth.NewClient(cfg.BaseURL)
		authResp, err := authClient.ValidateToken(token)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Update config with validated info
		cfg.Token = token
		cfg.OrgID = authResp.OrgID
	*/

	// For development, just save the token
	cfg.Token = token
	cfg.OrgID = "dev-org" // Mock org ID

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Success message
	green := color.New(color.FgGreen)
	green.Printf("✓ Successfully authenticated!\n")
	fmt.Printf("Configuration saved to: %s\n", config.GetConfigPath())

	// Show masked token for confirmation
	maskedToken := maskToken(token)
	fmt.Printf("Token: %s\n", maskedToken)

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
