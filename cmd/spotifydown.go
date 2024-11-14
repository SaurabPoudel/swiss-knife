/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/SaurabPoudel/swiss-knife/internal/spotifydownloader"
)

var (
	spotifyOutputDir string
	spotifyFormat    string
	playlistMode     bool
)

var spotifydownCmd = &cobra.Command{
	Use:   "spotifydown [URL]",
	Short: "Download Spotify tracks and playlists",
	Long: `Download music from Spotify URLs (tracks or playlists).
This command allows you to download both individual tracks and entire playlists
from Spotify by providing the URL.

Example: 
  # Download a single track
  spotifydown "https://open.spotify.com/track/xxxxx"
  
  # Download entire playlist
  spotifydown "https://open.spotify.com/playlist/xxxxx" --playlist
  
Note: Always wrap URLs in quotes to handle special characters correctly.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if err := spotifydownloader.Download(url, spotifyOutputDir, spotifyFormat, playlistMode); err != nil {
			fmt.Printf("Error downloading from Spotify: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Download completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(spotifydownCmd)

	homeDir, _ := os.UserHomeDir()
	defaultOutputDir := filepath.Join(homeDir, "Music", "Spotify")

	spotifydownCmd.Flags().StringVarP(&spotifyOutputDir, "output", "o", defaultOutputDir, "Output directory for downloaded music")
	spotifydownCmd.Flags().StringVarP(&spotifyFormat, "format", "f", "mp3", "Audio format (mp3, m4a)")
	spotifydownCmd.Flags().BoolVarP(&playlistMode, "playlist", "p", false, "Download entire playlist")
}
