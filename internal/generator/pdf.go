package generator

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ExportPDF(htmlFile, outPDF string) error {
	abs, err := filepath.Abs(htmlFile)
	if err != nil {
		return err
	}
	url := "file://" + abs
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithDisplayHeaderFooter(true).
				WithHeaderTemplate(`<span style="font-size:8px;margin-left:8px">docgen</span>`).
				WithFooterTemplate(`<span style="font-size:8px;margin-left:8px">第 <span class="pageNumber"></span> / <span class="totalPages"></span> 页</span>`).
				Do(ctx)
			if err != nil {
				return err
			}
			pdf = buf
			return nil
		}),
	); err != nil {
		return err
	}
	return os.WriteFile(outPDF, pdf, 0o644)
}
