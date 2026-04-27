# docgen (MVP)

一个用 Go 编写的“代码扫描 + 接口文档生成 + 操作手册自动生成 + PDF 导出”工具。

## 功能（当前 MVP）

- 扫描 Java Spring Controller（`@RestController/@Controller` + Mapping 注解）
- 提取基础接口信息：HTTP 方法、路径、方法名、参数、响应类型、源码位置
- 输出：`Markdown` / `HTML` / `OpenAPI JSON` / `PDF`
- 保存扫描中间产物：`api-index.json`
- 读取 `manual.yaml/json`，驱动浏览器执行基础步骤并逐步截图，生成操作手册（`Markdown/HTML/PDF`）
- 支持操作失败重试和结构化日志

## 快速开始

```bash
go mod tidy
go run ./cmd/docgen scan --path ./examples/spring-demo --out ./docs --format markdown,html,openapi
go run ./cmd/docgen manual --config ./examples/manual/manual.yaml --out ./manuals --format markdown,html
```

## CLI

```bash
docgen scan --path ./java-project --out ./docs --format markdown,pdf,openapi

docgen manual --config ./manual.yaml --out ./manuals --format markdown,pdf

docgen serve --dir ./docs --port 8899
```

## 目录说明

```text
cmd/docgen            # CLI 入口
internal/app          # 命令实现（scan/manual/serve）
internal/scanner/java # Java Spring 扫描器
internal/generator    # Markdown/HTML/OpenAPI/PDF 输出
internal/manual       # YAML 驱动的浏览器自动化与手册生成
internal/model        # 统一数据模型
examples              # 示例 Java 项目与手册配置
```

## 技术选型（MVP）

- Java 解析：先用注解+签名轻量解析器（预留 AST 替换点）
- 文档模板：Go `text/template`
- Markdown -> HTML：`goldmark`
- PDF：`chromedp` PrintToPDF（支持页眉页脚，适配中文字体）
- 浏览器自动化：`chromedp`
- CLI：`cobra`
- 配置：`yaml.v3` + `encoding/json`
- 日志：`zerolog`

> 下一阶段建议将 `internal/scanner/java` 切换到 `tree-sitter-java` AST 解析，补全 DTO 字段、继承、枚举和泛型展开。
