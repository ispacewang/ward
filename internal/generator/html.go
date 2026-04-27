package generator

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func MarkdownToHTML(md string) (string, error) {
	var out bytes.Buffer
	if err := goldmark.Convert([]byte(md), &out); err != nil {
		return "", err
	}
	prefix := `<!doctype html><html><head><meta charset="UTF-8"><style>
body{font-family:"Noto Sans CJK SC","Microsoft YaHei",Arial,sans-serif;padding:24px;line-height:1.6}
table{border-collapse:collapse;width:100%}td,th{border:1px solid #ddd;padding:8px}
code,pre{background:#f6f8fa;padding:2px 4px}
</style></head><body>`
	suffix := `</body></html>`
	return prefix + out.String() + suffix, nil
}
