package generator

import (
	"bytes"
	"text/template"

	"github.com/example/docgen/internal/model"
)

const apiTemplate = `# {{.ProjectName}}

> 更新时间：{{.GeneratedAt.Format "2006-01-02 15:04:05"}}

## 接口目录
{{range .Endpoints}}- [{{.Title}}](#{{.Function}})
{{end}}
{{range .Endpoints}}
## <a id="{{.Function}}"></a>{{.Title}}

**请求方式：** {{.Method}}  
**请求路径：** {{.Path}}  
**Controller：** {{.Controller}}  
**方法名：** {{.Function}}  
**说明：** {{.Description}}

### 请求参数
| 参数名 | 类型 | 位置 | 必填 | 说明 |
|---|---|---|---|---|
{{- range .RequestParams }}| {{.Name}} | {{.Type}} | {{.In}} | {{if .Required}}是{{else}}否{{end}} | {{.Description}} |
{{- end}}

### 响应结构
{{- if .ResponseBody}}
- 类型：` + "`{{.ResponseBody.TypeName}}`" + `
{{- else}}
- 未识别
{{- end}}

### 源码
` + "`{{.SourceFile}}:{{.SourceLine}}`" + `

---
{{end}}`

func RenderMarkdown(doc *model.APIDocument) (string, error) {
	tmpl, err := template.New("api").Parse(apiTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, doc); err != nil {
		return "", err
	}
	return buf.String(), nil
}
