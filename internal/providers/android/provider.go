package android

import (
	"github.com/backsofangels/grimoire/internal/providers"
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
		{Name: "no-wrapper", Short: "", Usage: "Skip Gradle wrapper generation", Default: false},
		{Name: "vscode", Short: "", Usage: "Generate .vscode", Default: true},
		{Name: "output-dir", Short: "o", Usage: "Output directory", Default: ""},
	}
}

// Prompt is implemented in prompts.go (interactive wizard).
// The implementation is separated to avoid cluttering the main provider file.

func (a *AndroidProvider) Validate(cfg providers.ProviderConfig) error {
	// Basic validation left to validator package; provider-level checks can go here.
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
