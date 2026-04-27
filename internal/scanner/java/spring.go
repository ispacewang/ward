package java

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/example/docgen/internal/model"
)

var (
	classDeclRE     = regexp.MustCompile(`class\s+([A-Za-z0-9_]+)`)                                                          //nolint:lll
	methodDeclRE    = regexp.MustCompile(`(public|protected|private)\s+[A-Za-z0-9_<>, ?\[\]]+\s+([A-Za-z0-9_]+)\s*\((.*)\)`) //nolint:lll
	requestMapArgRE = regexp.MustCompile(`"([^"]+)"`)
)

type SpringScanner struct {
	cfg model.ScanConfig
}

func NewSpringScanner(cfg model.ScanConfig) *SpringScanner {
	return &SpringScanner{cfg: cfg}
}

func (s *SpringScanner) Scan() (*model.APIDocument, error) {
	doc := &model.APIDocument{ProjectName: s.cfg.ProjectName, GeneratedAt: time.Now(), BaseURL: s.cfg.BaseURL}
	err := filepath.WalkDir(s.cfg.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && shouldIgnore(d.Name(), s.cfg.IgnoreDirs) {
			return filepath.SkipDir
		}
		if d.IsDir() || filepath.Ext(path) != ".java" {
			return nil
		}
		eps, err := s.scanFile(path)
		if err != nil {
			return nil
		}
		doc.Endpoints = append(doc.Endpoints, eps...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func shouldIgnore(name string, ignore []string) bool {
	for _, item := range ignore {
		if name == item {
			return true
		}
	}
	return false
}

func (s *SpringScanner) scanFile(path string) ([]model.Endpoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	controller := ""
	isController := false
	basePath := ""
	pending := make([]string, 0)
	var endpoints []model.Endpoint

	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "@") {
			pending = append(pending, line)
			continue
		}
		if m := classDeclRE.FindStringSubmatch(line); len(m) > 1 {
			controller = m[1]
			isController = containsAny(pending, "@RestController", "@Controller")
			if isController {
				basePath = extractMappingPath(pending, "@RequestMapping")
			}
			pending = nil
			continue
		}
		if !isController {
			pending = nil
			continue
		}
		if m := methodDeclRE.FindStringSubmatch(line); len(m) > 3 {
			methodName := m[2]
			paramsRaw := m[3]
			httpMethod, relPath := extractMethodAndPath(pending)
			if httpMethod == "" {
				pending = nil
				continue
			}
			ep := model.Endpoint{
				Title:       humanizeName(methodName),
				Description: extractDocDescription(pending, methodName),
				Method:      httpMethod,
				Path:        normalizePath(basePath, relPath),
				Controller:  controller,
				Function:    methodName,
				SourceFile:  path,
				SourceLine:  i + 1,
			}
			ep.RequestParams, ep.QueryParams, ep.PathParams, ep.Headers, ep.RequestBody = parseMethodParams(paramsRaw)
			ep.ResponseBody = inferResponse(line)
			ep.StatusCodes = []model.StatusDef{{Code: 200, Description: "OK"}}
			endpoints = append(endpoints, ep)
		}
		pending = nil
	}
	return endpoints, nil
}

func containsAny(annotations []string, targets ...string) bool {
	for _, a := range annotations {
		for _, t := range targets {
			if strings.Contains(a, t) {
				return true
			}
		}
	}
	return false
}

func extractMappingPath(annotations []string, mapping string) string {
	for _, a := range annotations {
		if !strings.Contains(a, mapping) {
			continue
		}
		if m := requestMapArgRE.FindStringSubmatch(a); len(m) > 1 {
			return m[1]
		}
	}
	return ""
}

func extractMethodAndPath(annotations []string) (string, string) {
	mapping := map[string]string{
		"@GetMapping":    "GET",
		"@PostMapping":   "POST",
		"@PutMapping":    "PUT",
		"@DeleteMapping": "DELETE",
		"@PatchMapping":  "PATCH",
	}
	for _, a := range annotations {
		for key, method := range mapping {
			if strings.Contains(a, key) {
				return method, firstArg(a)
			}
		}
		if strings.Contains(a, "@RequestMapping") {
			m := "GET"
			if strings.Contains(a, "RequestMethod.POST") {
				m = "POST"
			} else if strings.Contains(a, "RequestMethod.PUT") {
				m = "PUT"
			} else if strings.Contains(a, "RequestMethod.DELETE") {
				m = "DELETE"
			} else if strings.Contains(a, "RequestMethod.PATCH") {
				m = "PATCH"
			}
			return m, firstArg(a)
		}
	}
	return "", ""
}

func firstArg(annotation string) string {
	if m := requestMapArgRE.FindStringSubmatch(annotation); len(m) > 1 {
		return m[1]
	}
	return ""
}

func normalizePath(basePath, relPath string) string {
	full := strings.TrimSuffix(basePath, "/") + "/" + strings.TrimPrefix(relPath, "/")
	full = strings.ReplaceAll(full, "//", "/")
	if full == "" {
		return "/"
	}
	if !strings.HasPrefix(full, "/") {
		full = "/" + full
	}
	return full
}

func parseMethodParams(raw string) ([]model.Param, []model.Param, []model.Param, []model.Param, *model.TypeRef) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil, nil, nil, nil
	}
	parts := splitParams(raw)
	all := make([]model.Param, 0, len(parts))
	query := make([]model.Param, 0)
	path := make([]model.Param, 0)
	headers := make([]model.Param, 0)
	var body *model.TypeRef
	for _, p := range parts {
		name, pType := parseParamNameType(p)
		param := model.Param{Name: name, Type: pType, Required: true}
		switch {
		case strings.Contains(p, "@RequestParam"):
			param.In = "query"
			query = append(query, param)
		case strings.Contains(p, "@PathVariable"):
			param.In = "path"
			path = append(path, param)
		case strings.Contains(p, "@RequestHeader"):
			param.In = "header"
			headers = append(headers, param)
		case strings.Contains(p, "@RequestBody"):
			body = &model.TypeRef{TypeName: pType, RawType: pType}
			param.In = "body"
		default:
			param.In = "query"
			query = append(query, param)
		}
		all = append(all, param)
	}
	return all, query, path, headers, body
}

func parseParamNameType(input string) (string, string) {
	clean := strings.TrimSpace(input)
	chunks := strings.Fields(clean)
	if len(chunks) < 2 {
		return clean, "Object"
	}
	name := chunks[len(chunks)-1]
	ptype := chunks[len(chunks)-2]
	if strings.HasPrefix(name, "final") && len(chunks) > 2 {
		name = chunks[len(chunks)-2]
		ptype = chunks[len(chunks)-3]
	}
	return strings.Trim(name, ","), ptype
}

func splitParams(raw string) []string {
	res := make([]string, 0)
	depth := 0
	last := 0
	for i, r := range raw {
		switch r {
		case '<', '(', '[':
			depth++
		case '>', ')', ']':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				res = append(res, strings.TrimSpace(raw[last:i]))
				last = i + 1
			}
		}
	}
	res = append(res, strings.TrimSpace(raw[last:]))
	return res
}

func inferResponse(line string) *model.TypeRef {
	chunks := strings.Fields(line)
	if len(chunks) < 3 {
		return nil
	}
	for i, c := range chunks {
		if c == "public" || c == "protected" || c == "private" {
			if i+1 < len(chunks) {
				t := strings.TrimSpace(chunks[i+1])
				return &model.TypeRef{TypeName: t, RawType: t}
			}
		}
	}
	return nil
}

func extractDocDescription(annotations []string, fallback string) string {
	for _, a := range annotations {
		if strings.Contains(a, "@Operation") || strings.Contains(a, "@ApiOperation") {
			if m := requestMapArgRE.FindStringSubmatch(a); len(m) > 1 {
				return m[1]
			}
		}
	}
	return humanizeName(fallback)
}

func humanizeName(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	var out []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, ' ')
		}
		out = append(out, r)
	}
	return strings.Title(strings.TrimSpace(string(out))) //nolint:staticcheck
}
