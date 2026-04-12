package providers

import (
	"fmt"
)

// BaseProvider provides common functionality for all providers.
// Providers can embed this struct to inherit default implementations.
type BaseProvider struct {
	name        string
	description string
	flags       []ProviderFlag
}

// NewBaseProvider creates a new BaseProvider with the given name and description.
func NewBaseProvider(name, description string) *BaseProvider {
	return &BaseProvider{
		name:        name,
		description: description,
		flags:       []ProviderFlag{},
	}
}

// Name returns the provider name.
func (bp *BaseProvider) Name() string {
	return bp.name
}

// Description returns the provider description.
func (bp *BaseProvider) Description() string {
	return bp.description
}

// Flags returns the provider flags (can be overridden).
func (bp *BaseProvider) Flags() []ProviderFlag {
	return bp.flags
}

// SetFlags sets the provider flags.
func (bp *BaseProvider) SetFlags(flags []ProviderFlag) {
	bp.flags = flags
}

// AddFlag adds a flag to the provider.
func (bp *BaseProvider) AddFlag(flag ProviderFlag) {
	bp.flags = append(bp.flags, flag)
}

// Validate is a default validation method (can be overridden by providers).
// Returns nil by default.
func (bp *BaseProvider) Validate(cfg ProviderConfig) error {
	return nil
}

// Prompt is a default prompt method (can be overridden by providers).
// Returns empty config by default.
func (bp *BaseProvider) Prompt() (ProviderConfig, error) {
	return make(ProviderConfig), nil
}

// Generate is a default generate method (must be overridden by providers).
func (bp *BaseProvider) Generate(cfg ProviderConfig) error {
	return fmt.Errorf("Generate not implemented for provider %s", bp.name)
}

// Add is a default add method (can be overridden by providers).
func (bp *BaseProvider) Add(cfg ProviderConfig) error {
	return fmt.Errorf("Add not implemented for provider %s", bp.name)
}

// DoctorChecks is a default doctor checks method (can be overridden by providers).
func (bp *BaseProvider) DoctorChecks() []Check {
	return []Check{}
}

// GetString retrieves a string value from config with fallback keys.
func (bp *BaseProvider) GetString(cfg ProviderConfig, keys ...string) string {
	for _, k := range keys {
		if v, ok := cfg[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// GetStringOrDefault retrieves a string value from config with a default fallback.
func (bp *BaseProvider) GetStringOrDefault(cfg ProviderConfig, defaultVal string, keys ...string) string {
	for _, k := range keys {
		if v, ok := cfg[k].(string); ok && v != "" {
			return v
		}
	}
	return defaultVal
}

// GetInt retrieves an integer value from config with fallback keys.
// Handles both int and string types (string is converted via strconv.Atoi).
func (bp *BaseProvider) GetInt(cfg ProviderConfig, keys ...string) int {
	for _, k := range keys {
		if v, ok := cfg[k].(int); ok {
			return v
		}
		// Try string conversion
		if v, ok := cfg[k].(string); ok && v != "" {
			if n, err := bp.stringToInt(v); err == nil {
				return n
			}
		}
	}
	return 0
}

// GetIntOrDefault retrieves an integer value from config with a default fallback.
func (bp *BaseProvider) GetIntOrDefault(cfg ProviderConfig, defaultVal int, keys ...string) int {
	for _, k := range keys {
		if v, ok := cfg[k].(int); ok {
			return v
		}
		if v, ok := cfg[k].(string); ok && v != "" {
			if n, err := bp.stringToInt(v); err == nil {
				return n
			}
		}
	}
	return defaultVal
}

// GetBool retrieves a boolean value from config with fallback keys.
func (bp *BaseProvider) GetBool(cfg ProviderConfig, keys ...string) bool {
	for _, k := range keys {
		if v, ok := cfg[k].(bool); ok {
			return v
		}
	}
	return false
}

// GetBoolOrDefault retrieves a boolean value from config with a default fallback.
func (bp *BaseProvider) GetBoolOrDefault(cfg ProviderConfig, defaultVal bool, keys ...string) bool {
	for _, k := range keys {
		if v, ok := cfg[k].(bool); ok {
			return v
		}
	}
	return defaultVal
}

// stringToInt converts a string to int, attempting various formats.
func (bp *BaseProvider) stringToInt(s string) (int, error) {
	// Try standard integer parsing
	for _, ch := range s {
		if (ch < '0' || ch > '9') && ch != '-' {
			return 0, fmt.Errorf("invalid integer: %s", s)
		}
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
