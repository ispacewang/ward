package model

import "time"

type ScanConfig struct {
	Path        string
	OutDir      string
	Formats     string
	ProjectName string
	IgnoreDirs  []string
	BaseURL     string
}

type APIDocument struct {
	ProjectName string     `json:"projectName"`
	GeneratedAt time.Time  `json:"generatedAt"`
	Endpoints   []Endpoint `json:"endpoints"`
	BaseURL     string     `json:"baseUrl,omitempty"`
}

type Endpoint struct {
	Title          string      `json:"title"`
	Description    string      `json:"description,omitempty"`
	Method         string      `json:"method"`
	Path           string      `json:"path"`
	Controller     string      `json:"controller"`
	Function       string      `json:"function"`
	SourceFile     string      `json:"sourceFile"`
	SourceLine     int         `json:"sourceLine"`
	RequestParams  []Param     `json:"requestParams,omitempty"`
	RequestBody    *TypeRef    `json:"requestBody,omitempty"`
	ResponseBody   *TypeRef    `json:"responseBody,omitempty"`
	Headers        []Param     `json:"headers,omitempty"`
	PathParams     []Param     `json:"pathParams,omitempty"`
	QueryParams    []Param     `json:"queryParams,omitempty"`
	StatusCodes    []StatusDef `json:"statusCodes,omitempty"`
	PossibleErrors []StatusDef `json:"possibleErrors,omitempty"`
}

type Param struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	Description  string `json:"description,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	ExampleValue string `json:"exampleValue,omitempty"`
	In           string `json:"in,omitempty"`
}

type TypeRef struct {
	TypeName   string      `json:"typeName"`
	RawType    string      `json:"rawType,omitempty"`
	Fields     []Field     `json:"fields,omitempty"`
	SampleJSON interface{} `json:"sampleJson,omitempty"`
}

type Field struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	Description  string `json:"description,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	ExampleValue string `json:"exampleValue,omitempty"`
}

type StatusDef struct {
	Code        int    `json:"code"`
	Description string `json:"description,omitempty"`
}
