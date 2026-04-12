package common

import "strings"

// Normalize applies lowercase and trimspace.
func Normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// NormalizeLang normalizes a language identifier.
func NormalizeLang(s string) string {
	return Normalize(s)
}

// NormalizeUI normalizes a UI type identifier.
func NormalizeUI(s string) string {
	return Normalize(s)
}

// NormalizeDI normalizes a DI framework identifier.
func NormalizeDI(s string) string {
	return Normalize(s)
}

// NormalizeTemplate normalizes a template identifier.
func NormalizeTemplate(s string) string {
	return Normalize(s)
}
