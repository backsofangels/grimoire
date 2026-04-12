package common

import (
	"strconv"

	"github.com/backsofangels/grimoire/internal/providers"
)

// GetString retrieves a string config value with fallback to multiple keys.
// Returns empty string if not found.
func GetString(cfg providers.ProviderConfig, keys ...string) string {
	for _, k := range keys {
		if v, ok := cfg[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// GetStringDefault retrieves a string config value with a default fallback.
func GetStringDefault(cfg providers.ProviderConfig, defaultVal string, keys ...string) string {
	if s := GetString(cfg, keys...); s != "" {
		return s
	}
	return defaultVal
}

// GetInt retrieves an int config value with type coercion from string.
// Returns 0 if not found or not a valid positive integer.
func GetInt(cfg providers.ProviderConfig, keys ...string) int {
	for _, k := range keys {
		if v, ok := cfg[k].(int); ok && v > 0 {
			return v
		}
		if v, ok := cfg[k].(string); ok {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				return n
			}
		}
	}
	return 0
}

// GetIntDefault retrieves an int config value with a default fallback.
func GetIntDefault(cfg providers.ProviderConfig, defaultVal int, keys ...string) int {
	if n := GetInt(cfg, keys...); n > 0 {
		return n
	}
	return defaultVal
}

// GetBool retrieves a bool config value.
func GetBool(cfg providers.ProviderConfig, keys ...string) bool {
	for _, k := range keys {
		if v, ok := cfg[k].(bool); ok {
			return v
		}
		if v, ok := cfg[k].(string); ok && (v == "true" || v == "1" || v == "yes") {
			return true
		}
	}
	return false
}

// GetBoolDefault retrieves a bool config value with a default fallback.
func GetBoolDefault(cfg providers.ProviderConfig, defaultVal bool, keys ...string) bool {
	for _, k := range keys {
		if v, ok := cfg[k].(bool); ok {
			return v
		}
		if v, ok := cfg[k].(string); ok && (v == "true" || v == "1" || v == "yes") {
			return true
		}
		if v, ok := cfg[k].(string); ok && (v == "false" || v == "0" || v == "no") {
			return false
		}
	}
	return defaultVal
}
