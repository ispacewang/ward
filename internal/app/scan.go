package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/docgen/internal/generator"
	"github.com/example/docgen/internal/model"
	"github.com/example/docgen/internal/scanner/java"
	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	var cfg model.ScanConfig
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan Java Spring controllers and generate API docs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Path == "" {
				return fmt.Errorf("--path is required")
			}
			if cfg.OutDir == "" {
				cfg.OutDir = "./docs"
			}
			formats := parseFormats(cfg.Formats)
			if len(formats) == 0 {
				formats = []string{"markdown"}
			}
			scanner := java.NewSpringScanner(cfg)
			doc, err := scanner.Scan()
			if err != nil {
				return err
			}
			if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
				return err
			}
			if err := generator.WriteIntermediateJSON(filepath.Join(cfg.OutDir, "api-index.json"), doc); err != nil {
				return err
			}

			for _, f := range formats {
				if err := generator.GenerateAPIDocs(f, cfg.OutDir, doc); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&cfg.Path, "path", "", "Path to Java project")
	cmd.Flags().StringVar(&cfg.OutDir, "out", "./docs", "Output directory")
	cmd.Flags().StringVar(&cfg.Formats, "format", "markdown", "Output format(s): markdown,html,pdf,openapi")
	cmd.Flags().StringVar(&cfg.ProjectName, "project", "API 文档", "Project name")
	cmd.Flags().StringSliceVar(&cfg.IgnoreDirs, "ignore", []string{"target", "build", ".git", "node_modules"}, "Ignored directory names")
	cmd.Flags().StringVar(&cfg.BaseURL, "base-url", "", "Base URL for OpenAPI server")
	return cmd
}

func parseFormats(s string) []string {
	parts := strings.Split(s, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}
