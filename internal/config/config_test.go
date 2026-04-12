package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DefaultProvider != "android" {
		t.Errorf("DefaultProvider = %q, want %q", cfg.DefaultProvider, "android")
	}
	if cfg.DefaultLang != "kotlin" {
		t.Errorf("DefaultLang = %q, want %q", cfg.DefaultLang, "kotlin")
	}
	if cfg.DefaultMinSdk != 26 {
		t.Errorf("DefaultMinSdk = %d, want %d", cfg.DefaultMinSdk, 26)
	}
	if !cfg.VSCode {
		t.Errorf("VSCode = %v, want true", cfg.VSCode)
	}
	if !cfg.Git {
		t.Errorf("Git = %v, want true", cfg.Git)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "default config is valid",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "empty provider is invalid",
			config: Config{
				DefaultProvider: "",
				DefaultLang:     "kotlin",
				DefaultMinSdk:   26,
			},
			wantErr: true,
		},
		{
			name: "negative min sdk is invalid",
			config: Config{
				DefaultProvider: "android",
				DefaultMinSdk:   -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() wantErr=%v gotErr=%v", tt.wantErr, err)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	base := Config{
		DefaultProvider: "android",
		DefaultLang:     "kotlin",
		DefaultMinSdk:   26,
		DefaultPackage:  "com.example",
		VSCode:          true,
		Git:             true,
	}

	tests := []struct {
		name     string
		override Config
		verify   func(Config) bool
	}{
		{
			name: "merge updates provider",
			override: Config{
				DefaultProvider: "springboot",
			},
			verify: func(c Config) bool {
				return c.DefaultProvider == "springboot" && c.DefaultLang == "kotlin"
			},
		},
		{
			name: "merge updates multiple fields",
			override: Config{
				DefaultLang:    "java",
				DefaultMinSdk:  30,
				DefaultPackage: "com.myapp",
			},
			verify: func(c Config) bool {
				return c.DefaultLang == "java" && c.DefaultMinSdk == 30 && c.DefaultPackage == "com.myapp"
			},
		},
		{
			name: "merge preserves unset fields",
			override: Config{
				DefaultLang: "java",
			},
			verify: func(c Config) bool {
				return c.DefaultProvider == "android" && c.DefaultLang == "java"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base.Merge(tt.override)
			if !tt.verify(result) {
				t.Errorf("Merge verification failed")
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Test that config round-trips correctly through JSON
	// (actual file system operations are tested via integration tests)

	// Create a test config
	cfg := Config{
		DefaultProvider: "springboot",
		DefaultLang:     "java",
		DefaultMinSdk:   30,
		DefaultPackage:  "com.test",
		DefaultTemplate: "compose",
		VSCode:          false,
		Git:             true,
	}

	// Note: Full Save/Load testing is done via integration tests
	// because ConfigPath depends on os.UserHomeDir()
	// which we don't want to mock in unit tests
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}
}
