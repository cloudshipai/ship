package cloudship

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultAPIURL is the default CloudShip API endpoint
	DefaultAPIURL = "https://staging.cloudshipai.com/api/v1"
	// MaxFileSize is the maximum file size for artifacts (100MB)
	MaxFileSize = 100 * 1024 * 1024
)

// Client represents a CloudShip API client
type Client struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// NewClient creates a new CloudShip API client
func NewClient(apiKey string) *Client {
	apiURL := os.Getenv("CLOUDSHIP_API_URL")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}

	return &Client{
		apiKey: apiKey,
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadArtifactRequest represents the request to upload an artifact
type UploadArtifactRequest struct {
	FleetID  string                 `json:"fleet_id"`
	FileName string                 `json:"file_name"`
	FileType string                 `json:"file_type"`
	Content  string                 `json:"content"` // base64 encoded
	Metadata map[string]interface{} `json:"metadata"`
}

// UploadArtifactResponse represents the response from uploading an artifact
type UploadArtifactResponse struct {
	ArtifactID  string    `json:"artifact_id"`
	DownloadURL string    `json:"download_url"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
}

// UploadArtifact uploads an artifact to CloudShip
func (c *Client) UploadArtifact(req *UploadArtifactRequest) (*UploadArtifactResponse, error) {
	url := fmt.Sprintf("%s/artifacts/upload", c.apiURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var uploadResp UploadArtifactResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &uploadResp, nil
}

// UploadFile uploads a file as an artifact
func (c *Client) UploadFile(fleetID, filePath, scanType string, tags []string) (*UploadArtifactResponse, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Check file size
	if len(data) > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum of 100MB")
	}

	// Determine file type
	fileType := "application/octet-stream"
	if ext := getFileExtension(filePath); ext == ".json" {
		fileType = "application/json"
	} else if ext == ".yaml" || ext == ".yml" {
		fileType = "application/yaml"
	}

	// Create metadata
	metadata := map[string]interface{}{
		"scan_type":      scanType,
		"scan_timestamp": time.Now().UTC().Format(time.RFC3339),
		"source":         fmt.Sprintf("ship-cli/v%s", getVersion()),
		"tags":           tags,
	}

	// Add ship ID if available
	if shipID := os.Getenv("SHIP_ID"); shipID != "" {
		metadata["ship_id"] = shipID
	}

	// Add execution ID if available
	if execID := os.Getenv("EXECUTION_ID"); execID != "" {
		metadata["execution_id"] = execID
	}

	// Create request
	req := &UploadArtifactRequest{
		FleetID:  fleetID,
		FileName: getFileName(filePath),
		FileType: fileType,
		Content:  base64.StdEncoding.EncodeToString(data),
		Metadata: metadata,
	}

	return c.UploadArtifact(req)
}

// ListArtifactsRequest represents the request to list artifacts
type ListArtifactsRequest struct {
	FleetID string
	Limit   int
	Offset  int
	Type    string
}

// ListArtifactsResponse represents the response from listing artifacts
type ListArtifactsResponse struct {
	Artifacts []Artifact `json:"artifacts"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// Artifact represents an artifact in the list
type Artifact struct {
	ID        string                 `json:"id"`
	FileName  string                 `json:"file_name"`
	FileType  string                 `json:"file_type"`
	FileSize  int64                  `json:"file_size"`
	Version   int                    `json:"version"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ListArtifacts lists artifacts for a fleet
func (c *Client) ListArtifacts(req *ListArtifactsRequest) (*ListArtifactsResponse, error) {
	url := fmt.Sprintf("%s/artifacts?fleet_id=%s", c.apiURL, req.FleetID)

	if req.Limit > 0 {
		url += fmt.Sprintf("&limit=%d", req.Limit)
	}
	if req.Offset > 0 {
		url += fmt.Sprintf("&offset=%d", req.Offset)
	}
	if req.Type != "" {
		url += fmt.Sprintf("&type=%s", req.Type)
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var listResp ListArtifactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &listResp, nil
}

// DownloadArtifact downloads an artifact by ID
func (c *Client) DownloadArtifact(artifactID string) ([]byte, error) {
	url := fmt.Sprintf("%s/artifacts/%s/download", c.apiURL, artifactID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Helper functions

func getFileExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

func getFileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}

func getVersion() string {
	// This will be replaced with actual version
	return "0.3.5"
}