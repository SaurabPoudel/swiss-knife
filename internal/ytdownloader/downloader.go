package ytdownloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Rename the function to match what we're calling from ytdown.go
func DownloadVideo(url, outputDir, format string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare yt-dlp command
	outputTemplate := filepath.Join(outputDir, "%(title)s.%(ext)s")
	args := []string{
		"--format", "bestvideo[ext=" + format + "]+bestaudio[ext=m4a]/best[ext=" + format + "]/best",
		"--merge-output-format", format,
		"-o", outputTemplate,
		url,
	}

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
} 