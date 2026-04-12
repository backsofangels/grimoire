package testing

import (
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
)

// ProviderTester provides a convenient interface for testing a provider.
type ProviderTester struct {
	helper *TestHelper
	cfg    *ConfigBuilder
}

// NewProviderTester creates a new ProviderTester with a test helper and config builder.
func NewProviderTester(t *testing.T) *ProviderTester {
	return &ProviderTester{
		helper: NewTestHelper(t),
		cfg:    NewConfig(),
	}
}

// TempDir returns the temporary directory.
func (pt *ProviderTester) TempDir() string {
	return pt.helper.TempDir()
}

// Config returns the ConfigBuilder for fluent configuration.
func (pt *ProviderTester) Config() *ConfigBuilder {
	return pt.cfg
}

// Helpers returns the TestHelper for assertions and path utilities.
func (pt *ProviderTester) Helpers() *TestHelper {
	return pt.helper
}

// GetConfig builds and returns the current configuration.
func (pt *ProviderTester) GetConfig() providers.ProviderConfig {
	return pt.cfg.Build()
}
