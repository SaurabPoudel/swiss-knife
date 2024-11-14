/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/SaurabPoudel/swiss-knife/internal/ytdownloader"
)

var (
	outputDir string
	format    string
)

// ytdownCmd represents the ytdown command
var ytdownCmd = &cobra.Command{
	Use:   "ytdown [URL]",
	Short: "Download YouTube videos",
	Long: `Download YouTube videos and shorts from provided URLs.
This command allows you to download both regular YouTube videos
and YouTube Shorts content by providing the video URL.

Example: ytdown "https://youtube.com/watch?v=xxxxx"
Note: Always wrap URLs in quotes to handle special characters correctly.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if err := ytdownloader.DownloadVideo(url, outputDir, format); err != nil {
			fmt.Printf("Error downloading video: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Download completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(ytdownCmd)

	homeDir, _ := os.UserHomeDir()
	defaultOutputDir := filepath.Join(homeDir, "Downloads")

	ytdownCmd.Flags().StringVarP(&outputDir, "output", "o", defaultOutputDir, "Output directory for downloaded videos")
	ytdownCmd.Flags().StringVarP(&format, "format", "f", "mp4", "Video format (mp4, webm)")
}
