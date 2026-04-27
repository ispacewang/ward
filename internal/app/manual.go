package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/docgen/internal/generator"
	"github.com/example/docgen/internal/manual"
	"github.com/spf13/cobra"
)

func newManualCmd() *cobra.Command {
	var configPath string
	var outDir string
	var formats string

	cmd := &cobra.Command{
		Use:   "manual",
		Short: "Generate operation manuals from YAML/JSON steps",
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("--config is required")
			}
			cfg, err := manual.LoadConfig(configPath)
			if err != nil {
				return err
			}
			if outDir != "" {
				cfg.OutputDir = outDir
			}
			if cfg.OutputDir == "" {
				cfg.OutputDir = "./manuals"
			}
			if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
				return err
			}

			report, err := manual.NewRunner(cfg).Run()
			if err != nil {
				return err
			}
			for _, f := range parseFormats(formats) {
				switch f {
				case "markdown":
					if err := os.WriteFile(filepath.Join(cfg.OutputDir, "manual.md"), []byte(report.Markdown), 0o644); err != nil {
						return err
					}
				case "html":
					html, err := generator.MarkdownToHTML(report.Markdown)
					if err != nil {
						return err
					}
					if err := os.WriteFile(filepath.Join(cfg.OutputDir, "manual.html"), []byte(html), 0o644); err != nil {
						return err
					}
				case "pdf":
					html, err := generator.MarkdownToHTML(report.Markdown)
					if err != nil {
						return err
					}
					tmpHTML := filepath.Join(cfg.OutputDir, "manual.html")
					if err := os.WriteFile(tmpHTML, []byte(html), 0o644); err != nil {
						return err
					}
					if err := generator.ExportPDF(tmpHTML, filepath.Join(cfg.OutputDir, "manual.pdf")); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unsupported format: %s", f)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "Path to manual config (yaml/json)")
	cmd.Flags().StringVar(&outDir, "out", "", "Output directory")
	cmd.Flags().StringVar(&formats, "format", "markdown,html", "Output format(s): markdown,html,pdf")
	return cmd
}

func parseFormatsOrDefault(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{"markdown"}
	}
	return parseFormats(s)
}
