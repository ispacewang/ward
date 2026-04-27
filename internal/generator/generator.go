package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/docgen/internal/model"
)

func GenerateAPIDocs(format, outDir string, doc *model.APIDocument) error {
	switch format {
	case "markdown":
		md, err := RenderMarkdown(doc)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(outDir, "api.md"), []byte(md), 0o644)
	case "html":
		md, err := RenderMarkdown(doc)
		if err != nil {
			return err
		}
		html, err := MarkdownToHTML(md)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(outDir, "api.html"), []byte(html), 0o644)
	case "pdf":
		md, err := RenderMarkdown(doc)
		if err != nil {
			return err
		}
		html, err := MarkdownToHTML(md)
		if err != nil {
			return err
		}
		tmp := filepath.Join(outDir, "api.html")
		if err := os.WriteFile(tmp, []byte(html), 0o644); err != nil {
			return err
		}
		return ExportPDF(tmp, filepath.Join(outDir, "api.pdf"))
	case "openapi":
		content, err := RenderOpenAPI(doc)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(outDir, "openapi.json"), []byte(content), 0o644)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteIntermediateJSON(path string, doc *model.APIDocument) error {
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
