package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultProvider string `json:"default_provider"`
	DefaultLang     string `json:"default_lang"`
	DefaultMinSdk   int    `json:"default_min_sdk"`
	DefaultPackage  string `json:"default_package"`
	DefaultTemplate string `json:"default_template"`
	VSCode          bool   `json:"vscode"`
	Git             bool   `json:"git"`
}

func DefaultConfig() Config {
	return Config{
		DefaultProvider: "android",
		DefaultLang:     "kotlin",
		DefaultMinSdk:   26,
		DefaultPackage:  "com.example",
		DefaultTemplate: "basic",
		VSCode:          true,
		Git:             true,
	}
}

func Load() (Config, error) {
	cfg := DefaultConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}
	path := filepath.Join(home, ".grimoire", "config.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
