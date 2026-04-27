package generator

import (
	"encoding/json"

	"github.com/example/docgen/internal/model"
)

func RenderOpenAPI(doc *model.APIDocument) (string, error) {
	paths := map[string]map[string]interface{}{}
	for _, ep := range doc.Endpoints {
		if paths[ep.Path] == nil {
			paths[ep.Path] = map[string]interface{}{}
		}
		params := make([]map[string]interface{}, 0, len(ep.RequestParams))
		for _, p := range ep.RequestParams {
			params = append(params, map[string]interface{}{
				"name":     p.Name,
				"in":       p.In,
				"required": p.Required,
				"schema": map[string]string{
					"type": p.Type,
				},
			})
		}
		paths[ep.Path][stringsLower(ep.Method)] = map[string]interface{}{
			"summary":     ep.Title,
			"description": ep.Description,
			"operationId": ep.Function,
			"parameters":  params,
			"responses": map[string]interface{}{
				"200": map[string]string{"description": "OK"},
			},
		}
	}
	obj := map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]string{
			"title":   doc.ProjectName,
			"version": "1.0.0",
		},
		"paths": paths,
	}
	if doc.BaseURL != "" {
		obj["servers"] = []map[string]string{{"url": doc.BaseURL}}
	}
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func stringsLower(in string) string {
	bs := []byte(in)
	for i, b := range bs {
		if b >= 'A' && b <= 'Z' {
			bs[i] = b + 32
		}
	}
	return string(bs)
}
