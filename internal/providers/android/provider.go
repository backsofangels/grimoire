package android

import (
	"fmt"
	"github.com/backsofangels/grimoire/internal/providers"
	"strings"
)

type AndroidProvider struct{}

func (a *AndroidProvider) Name() string { return "android" }

func (a *AndroidProvider) Description() string {
	return "Android project provider (Kotlin/Java)"
}

func (a *AndroidProvider) Flags() []providers.ProviderFlag {
	return []providers.ProviderFlag{
		{Name: "package", Short: "p", Usage: "Package name", Default: ""},
		{Name: "lang", Short: "l", Usage: "Language (kotlin|java)", Default: "kotlin"},
		{Name: "min-sdk", Short: "", Usage: "Minimum SDK", Default: 26},
		{Name: "target-sdk", Short: "", Usage: "Target SDK", Default: 35},
		{Name: "template", Short: "t", Usage: "Template (basic|empty)", Default: "basic"},
		{Name: "git", Short: "", Usage: "Initialize git", Default: true},
		{Name: "wrapper", Short: "", Usage: "Generate Gradle wrapper (boolean, default: true)", Default: true},
		{Name: "vscode", Short: "", Usage: "Generate .vscode", Default: true},
		{Name: "output-dir", Short: "o", Usage: "Output directory", Default: ""},
	}
}

// Prompt is implemented in prompts.go (interactive wizard).
// The implementation is separated to avoid cluttering the main provider file.

func (a *AndroidProvider) Validate(cfg providers.ProviderConfig) error {
	// Validate language/template compatibility and normalize some keys.
	var lang string
	if v, ok := cfg["Lang"].(string); ok && v != "" {
		lang = v
	} else if v2, ok := cfg["lang"].(string); ok && v2 != "" {
		lang = v2
	}
	var tmpl string
	if t, ok := cfg["Template"].(string); ok && t != "" {
		tmpl = t
	} else if t2, ok := cfg["template"].(string); ok && t2 != "" {
		tmpl = t2
	}
	lang = strings.ToLower(strings.TrimSpace(lang))
	tmpl = strings.ToLower(strings.TrimSpace(tmpl))

	if lang == "java" && tmpl == "compose" {
		return fmt.Errorf("compose template is not available for Java — use --lang kotlin")
	}
	return nil
}

func (a *AndroidProvider) Generate(cfg providers.ProviderConfig) error {
	return GenerateProject(cfg)
}

func (a *AndroidProvider) DoctorChecks() []providers.Check {
	return nil
}

func init() {
	providers.Register(&AndroidProvider{})
}
