package modules

import (
	"fmt"
	"os"
	"strings"
)

// getImageTag returns the image tag from environment variable override or falls back to default
func getImageTag(toolName, defaultImage string) string {
	envKey := fmt.Sprintf("SHIP_IMAGE_TAG_%s", strings.ToUpper(toolName))
	if override := os.Getenv(envKey); override != "" {
		return override
	}
	return defaultImage
}