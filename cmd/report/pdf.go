package report

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

// outputPDF generates a PDF report from the BlockReport data
func outputPDF(report *BlockReport, outputFile string) error {
	log.Info().Msg("Generating PDF report from HTML")

	// Generate HTML from the existing template
	html := generateHTML(report)

	// Create chromedp context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Allocate a new browser context
	ctx, cancelChrome := chromedp.NewContext(ctx)
	defer cancelChrome()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Get the frame tree to set document content
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to get frame tree: %w", err)
			}

			// Set the HTML content
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Wait a bit for any dynamic content to settle, respecting context cancellation
			timer := time.NewTimer(500 * time.Millisecond)
			defer timer.Stop()
			select {
			case <-timer.C:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Print to PDF with appropriate settings
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithScale(0.8).
				WithPreferCSSPageSize(false).
				WithPaperWidth(8.5).
				WithPaperHeight(11).
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to print to PDF: %w", err)
			}
			return nil
		}),
	)

	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w\n\nPDF generation requires Google Chrome or Chromium to be installed on your system.\nPlease install Chrome/Chromium and try again. See documentation for installation instructions", err)
	}

	// Write PDF to file
	if err := os.WriteFile(outputFile, buf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	log.Info().Str("file", outputFile).Msg("PDF report written")
	return nil
}
