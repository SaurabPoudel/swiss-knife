package md2pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/page"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const defaultTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            padding: 2em;
            max-width: 100%;
            margin: 0 auto;
        }
        pre {
            background-color: #f6f8fa;
            padding: 16px;
            border-radius: 6px;
            overflow-x: auto;
        }
        code {
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
        }
        img {
            max-width: 100%;
            height: auto;
        }
        {{if .CustomCSS}}
        {{.CustomCSS}}
        {{end}}
    </style>
</head>
<body>
    {{.Content}}
</body>
</html>
`

func MarkdownToPDF(inputFile, outputFile, pageSize, margins, cssFile string) error {
	// Read markdown file
	mdContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Initialize markdown parser
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := md.Convert(mdContent, &buf); err != nil {
		return fmt.Errorf("failed to convert markdown: %w", err)
	}

	// Read custom CSS if provided
	var customCSS string
	if cssFile != "" {
		cssContent, err := ioutil.ReadFile(cssFile)
		if err != nil {
			return fmt.Errorf("failed to read CSS file: %w", err)
		}
		customCSS = string(cssContent)
	}

	// Apply template
	tmpl, err := template.New("pdf").Parse(defaultTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, struct {
		Content   template.HTML
		CustomCSS string
	}{
		Content:   template.HTML(buf.String()),
		CustomCSS: customCSS,
	}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create temporary HTML file
	tmpFile, err := ioutil.TempFile("", "markdown-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(htmlBuf.Bytes()); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Setup Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Generate PDF
	var pdfBuf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("file://"+tmpFile.Name()),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(float64(1.0)).
				WithMarginBottom(float64(1.0)).
				WithMarginLeft(float64(1.0)).
				WithMarginRight(float64(1.0)).
				WithPaperWidth(8.27).
				WithPaperHeight(11.7).
				Do(ctx)
			return err
		}),
	); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Write PDF file
	if err := ioutil.WriteFile(outputFile, pdfBuf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
} 