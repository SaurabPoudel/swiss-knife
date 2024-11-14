package spotifydownloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Download handles downloading tracks or playlists from Spotify
func Download(url, outputDir, format string, isPlaylist bool) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare spotdl command
	args := []string{
		"download",
		"--output", filepath.Join(outputDir, "{artist} - {title}.{ext}"),
		"--format", format,
		"--threads", "4", // For faster downloads
	}

	if isPlaylist {
		args = append(args, "--playlist")
	}

	// Add URL as the last argument
	args = append(args, url)

	// Check if spotdl is installed
	if err := checkSpotDL(); err != nil {
		return fmt.Errorf("spotdl not found: %w\nPlease install it using: pip install spotdl", err)
	}

	// Execute spotdl command
	cmd := exec.Command("spotdl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("spotdl not found. Please install it using: pip install spotdl")
		}
		return fmt.Errorf("failed to download: %w", err)
	}

	return nil
}

// checkSpotDL verifies if spotdl is installed
func checkSpotDL() error {
	cmd := exec.Command("spotdl", "--version")
	return cmd.Run()
}

// Helper function to clean filenames
func sanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename

	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}

	return result
} 