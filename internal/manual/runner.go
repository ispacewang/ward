package manual

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/example/docgen/internal/common"
)

type Runner struct {
	cfg Config
}

type StepResult struct {
	Index       int
	Name        string
	Description string
	Screenshot  string
	Success     bool
	Error       string
}

type Report struct {
	Markdown string
	Steps    []StepResult
}

func NewRunner(cfg Config) *Runner { return &Runner{cfg: cfg} }

func (r *Runner) Run() (*Report, error) {
	logger := common.NewLogger()
	if err := os.MkdirAll(filepath.Join(r.cfg.OutputDir, "screenshots"), 0o755); err != nil {
		return nil, err
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	res := make([]StepResult, 0, len(r.cfg.Steps))
	for i, step := range r.cfg.Steps {
		var runErr error
		shot := filepath.Join("screenshots", fmt.Sprintf("step-%d.png", i+1))
		for attempt := 1; attempt <= r.cfg.Retry; attempt++ {
			runErr = executeStep(ctx, r.cfg, step, filepath.Join(r.cfg.OutputDir, shot))
			if runErr == nil {
				break
			}
			logger.Warn().Err(runErr).Int("step", i+1).Int("attempt", attempt).Msg("step failed, retry")
		}
		item := StepResult{Index: i + 1, Name: step.Name, Description: step.Description, Screenshot: shot, Success: runErr == nil}
		if runErr != nil {
			item.Error = runErr.Error()
		}
		res = append(res, item)
	}

	md, err := renderManualMarkdown(r.cfg.Title, res)
	if err != nil {
		return nil, err
	}
	_ = logger.Info().Int("steps", len(res)).Msg("manual done")
	return &Report{Markdown: md, Steps: res}, nil
}

func executeStep(ctx context.Context, cfg Config, step Step, screenshotPath string) error {
	actions := make([]chromedp.Action, 0)
	assertExpected := ""
	assertActual := ""
	switch strings.ToLower(step.Action) {
	case "goto":
		actions = append(actions, chromedp.Navigate(joinURL(cfg.BaseURL, step.URL)))
	case "click":
		actions = append(actions, chromedp.Click(step.Selector, chromedp.ByQuery))
	case "input":
		actions = append(actions, chromedp.SetValue(step.Selector, step.Value, chromedp.ByQuery))
	case "wait":
		if step.Selector != "" {
			actions = append(actions, chromedp.WaitVisible(step.Selector, chromedp.ByQuery))
		} else {
			time.Sleep(time.Duration(step.WaitMS) * time.Millisecond)
		}
	case "screenshot":
		if step.Selector != "" {
			var buf []byte
			actions = append(actions, chromedp.Screenshot(step.Selector, &buf, chromedp.ByQuery), chromedp.ActionFunc(func(ctx context.Context) error {
				return os.WriteFile(screenshotPath, buf, 0o644)
			}))
		} else {
			var buf []byte
			actions = append(actions, chromedp.FullScreenshot(&buf, 90), chromedp.ActionFunc(func(ctx context.Context) error {
				return os.WriteFile(screenshotPath, buf, 0o644)
			}))
		}
	case "assert_text":
		assertExpected = step.Value
		actions = append(actions, chromedp.Text(step.Selector, &assertActual, chromedp.ByQuery))
	default:
		return fmt.Errorf("unsupported action: %s", step.Action)
	}
	if len(actions) > 0 {
		if err := chromedp.Run(ctx, actions...); err != nil {
			return err
		}
	}
	if assertExpected != "" && !strings.Contains(assertActual, assertExpected) {
		return fmt.Errorf("assert_text failed: %q not in %q", assertExpected, assertActual)
	}
	if strings.ToLower(step.Action) != "screenshot" {
		var buf []byte
		if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 90)); err != nil {
			return err
		}
		return os.WriteFile(screenshotPath, buf, 0o644)
	}
	return nil
}

func joinURL(base, path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return strings.TrimSuffix(base, "/") + "/" + strings.TrimPrefix(path, "/")
}

const manualTmpl = `# {{.Title}}

{{range .Steps}}
## {{.Index}}. {{.Name}}

{{.Description}}

{{if .Success}}![{{.Name}}](./{{.Screenshot}}){{else}}> 失败：{{.Error}}{{end}}

{{end}}`

func renderManualMarkdown(title string, steps []StepResult) (string, error) {
	if title == "" {
		title = "操作手册"
	}
	var buf bytes.Buffer
	err := template.Must(template.New("manual").Parse(manualTmpl)).Execute(&buf, map[string]interface{}{
		"Title": title,
		"Steps": steps,
	})
	return buf.String(), err
}
