package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudshipai/ship/internal/cloudship"
	"github.com/cloudshipai/ship/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	pushFleetID  string
	pushScanType string
	pushTags     []string
	pushMetadata map[string]string
)

var pushCmd = &cobra.Command{
	Use:   "push [file]",
	Short: "Upload an artifact to CloudShip",
	Long: `Upload an artifact file to CloudShip for analysis.

Examples:
  # Upload a Terraform plan
  ship push terraform.tfplan --fleet-id your-fleet-id --type terraform_plan

  # Upload a security scan result
  ship push security-scan.json --fleet-id your-fleet-id --type security_scan --tags "production,critical"

  # Upload with custom metadata
  ship push inventory.json --fleet-id your-fleet-id --type aws_inventory --metadata "region=us-east-1,account=123456789"

  # Use environment variables
  export CLOUDSHIP_FLEET_ID=your-fleet-id
  ship push cost-analysis.json --type cost_analysis

You can also pipe data:
  terraform plan -out=tfplan && ship push tfplan --type terraform_plan
  ship terraform-tools security-scan | ship push - --type security_scan`,
	Args: cobra.ExactArgs(1),
	RunE: runPush,
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&pushFleetID, "fleet-id", "", "Fleet ID (can also use CLOUDSHIP_FLEET_ID env var)")
	pushCmd.Flags().StringVar(&pushScanType, "type", "", "Artifact type (e.g., terraform_plan, security_scan, cost_analysis)")
	pushCmd.Flags().StringSliceVar(&pushTags, "tags", []string{}, "Tags for the artifact (comma-separated)")
	pushCmd.Flags().StringToStringVar(&pushMetadata, "metadata", map[string]string{}, "Additional metadata as key=value pairs")

	pushCmd.MarkFlagRequired("type")
}

func runPush(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check for API key
	if cfg.APIKey == "" {
		return fmt.Errorf("not authenticated - run 'ship auth --api-key YOUR_KEY' first")
	}

	// Get fleet ID from flag or config
	fleetID := pushFleetID
	if fleetID == "" {
		fleetID = cfg.FleetID
	}
	if fleetID == "" {
		return fmt.Errorf("fleet ID required - use --fleet-id flag or set CLOUDSHIP_FLEET_ID")
	}

	// Get file path
	filePath := args[0]

	// Handle stdin
	var data []byte
	if filePath == "-" {
		data, err = readStdin()
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
		filePath = fmt.Sprintf("stdin-%s.json", time.Now().Format("20060102-150405"))
	} else {
		// Check if file exists
		if _, err := os.Stat(filePath); err != nil {
			return fmt.Errorf("file not found: %s", filePath)
		}

		// Read file
		data, err = os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	}

	// Check file size
	if len(data) > cloudship.MaxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", len(data), cloudship.MaxFileSize)
	}

	// Create client
	client := cloudship.NewClient(cfg.APIKey)

	// Prepare metadata
	metadata := map[string]interface{}{
		"scan_type":      pushScanType,
		"scan_timestamp": time.Now().UTC().Format(time.RFC3339),
		"source":         fmt.Sprintf("ship-cli/v%s", getVersion()),
		"tags":           pushTags,
	}

	// Add custom metadata
	for k, v := range pushMetadata {
		metadata[k] = v
	}

	// Add ship ID if available
	if shipID := os.Getenv("SHIP_ID"); shipID != "" {
		metadata["ship_id"] = shipID
	}

	// Add execution ID if available
	if execID := os.Getenv("EXECUTION_ID"); execID != "" {
		metadata["execution_id"] = execID
	}

	// Determine file type
	fileType := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		fileType = "application/json"
	case ".yaml", ".yml":
		fileType = "application/yaml"
	case ".tfplan":
		fileType = "application/x-terraform-plan"
	case ".log":
		fileType = "text/plain"
	}

	// Create upload request
	req := &cloudship.UploadArtifactRequest{
		FleetID:  fleetID,
		FileName: filepath.Base(filePath),
		FileType: fileType,
		Content:  base64Encode(data),
		Metadata: metadata,
	}

	// Upload artifact
	fmt.Printf("Uploading %s to CloudShip...\n", filepath.Base(filePath))

	resp, err := client.UploadArtifact(req)
	if err != nil {
		return fmt.Errorf("failed to upload artifact: %w", err)
	}

	// Success
	green := color.New(color.FgGreen)
	green.Printf("✓ Successfully uploaded artifact!\n")

	fmt.Printf("\nArtifact ID: %s\n", resp.ArtifactID)
	fmt.Printf("Version: %d\n", resp.Version)
	fmt.Printf("Download URL: %s\n", resp.DownloadURL)
	fmt.Printf("Created At: %s\n", resp.CreatedAt.Format(time.RFC3339))

	return nil
}

func readStdin() ([]byte, error) {
	var data []byte
	buf := make([]byte, 1024)

	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
	}

	return data, nil
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func getVersion() string {
	// This will be replaced with actual version
	return "0.3.5"
}
