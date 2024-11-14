/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/SaurabPoudel/swiss-knife/internal/md2pdf"
)

var (
	outputFile string
	pageSize   string
	margins    string
	css        string
)

var mdtopdfCmd = &cobra.Command{
	Use:   "mdtopdf [input.md]",
	Short: "Markdown to PDF that just works",
	Long: `A command-line tool to convert Markdown files to PDF format.
This tool provides a simple and reliable way to convert your Markdown documents
into professional-looking PDF files while preserving formatting and styling.

Example: mdtopdf input.md
         mdtopdf input.md -o output.pdf --page-size A4 --margins "1in"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]

		// If no output file specified, use input filename with .pdf extension
		if outputFile == "" {
			outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".pdf"
		}

		// Convert the file
		if err := md2pdf.MarkdownToPDF(inputFile, outputFile, pageSize, margins, css); err != nil {
			fmt.Printf("Error converting file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(mdtopdfCmd)

	mdtopdfCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output PDF file (default: input filename with .pdf extension)")
	mdtopdfCmd.Flags().StringVar(&pageSize, "page-size", "A4", "Page size (A4, Letter, Legal)")
	mdtopdfCmd.Flags().StringVar(&margins, "margins", "1in", "Page margins (e.g., 1in, 2cm)")
	mdtopdfCmd.Flags().StringVar(&css, "css", "", "Custom CSS file path for styling")
}
