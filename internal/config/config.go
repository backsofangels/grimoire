package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the Grimoire configuration file.
type Config struct {
	DefaultProvider string `json:"default_provider"`
	DefaultLang     string `json:"default_lang"`
	DefaultMinSdk   int    `json:"default_min_sdk"`
	DefaultPackage  string `json:"default_package"`
	DefaultTemplate string `json:"default_template"`
	VSCode          bool   `json:"vscode"`
	Git             bool   `json:"git"`
}

// DefaultConfig returns a Config with default values.
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

// ConfigPath returns the path to the config file in the user's home directory.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, ".grimoire", "config.json"), nil
}

// ConfigDir returns the directory where the config file is stored.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, ".grimoire"), nil
}

// Load loads the config file from the user's home directory.
// If the file doesn't exist, returns default config without error.
// If the file exists but is invalid, returns an error.
func Load() (Config, error) {
	cfg := DefaultConfig()
	path, err := ConfigPath()
	if err != nil {
		// Can't determine path, return defaults
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Config file doesn't exist, return defaults
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("read config file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config file %s: %w", path, err)
	}

	return cfg, nil
}

// Save persists the config to the config file in the user's home directory.
// Creates the .grimoire directory if it doesn't exist.
func (c Config) Save() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config directory %s: %w", dir, err)
	}

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config file %s: %w", path, err)
	}

	return nil
}

// Validate checks if the config values are valid.
func (c Config) Validate() error {
	if c.DefaultProvider == "" {
		return fmt.Errorf("default_provider cannot be empty")
	}
	if c.DefaultMinSdk < 0 {
		return fmt.Errorf("default_min_sdk must be non-negative, got %d", c.DefaultMinSdk)
	}
	return nil
}

// Merge merges another config into this one, with the other config taking precedence
// for non-empty values. This is useful for overriding defaults.
func (c Config) Merge(other Config) Config {
	if other.DefaultProvider != "" {
		c.DefaultProvider = other.DefaultProvider
	}
	if other.DefaultLang != "" {
		c.DefaultLang = other.DefaultLang
	}
	if other.DefaultMinSdk > 0 {
		c.DefaultMinSdk = other.DefaultMinSdk
	}
	if other.DefaultPackage != "" {
		c.DefaultPackage = other.DefaultPackage
	}
	if other.DefaultTemplate != "" {
		c.DefaultTemplate = other.DefaultTemplate
	}
	// Note: Boolean fields are tricky to merge since false is the zero value.
	// We only merge if they're true (assuming true is an explicit override).
	if other.VSCode {
		c.VSCode = other.VSCode
	}
	if other.Git {
		c.Git = other.Git
	}
	return c
}
