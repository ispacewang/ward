package manual

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title     string `yaml:"title" json:"title"`
	BaseURL   string `yaml:"baseUrl" json:"baseUrl"`
	OutputDir string `yaml:"outputDir" json:"outputDir"`
	Retry     int    `yaml:"retry" json:"retry"`
	Steps     []Step `yaml:"steps" json:"steps"`
}

type Step struct {
	Name        string `yaml:"name" json:"name"`
	Action      string `yaml:"action" json:"action"`
	URL         string `yaml:"url" json:"url"`
	Selector    string `yaml:"selector" json:"selector"`
	Value       string `yaml:"value" json:"value"`
	Description string `yaml:"description" json:"description"`
	WaitMS      int    `yaml:"waitMs" json:"waitMs"`
}

func LoadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(b, &cfg)
	case ".json":
		err = json.Unmarshal(b, &cfg)
	default:
		err = fmt.Errorf("unsupported config ext: %s", ext)
	}
	if err != nil {
		return Config{}, err
	}
	if cfg.Retry < 1 {
		cfg.Retry = 1
	}
	return cfg, nil
}
