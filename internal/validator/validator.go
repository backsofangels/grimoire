package validator

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func ValidateAppName(name string) error {
	if name == "" {
		return errors.New("app name cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("app name too long")
	}
	for i, r := range name {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return errors.New("app name contains invalid characters")
		}
		if i == 0 && unicode.IsDigit(r) {
			return errors.New("app name cannot start with a digit")
		}
	}
	return nil
}

var pkgSegment = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

func ValidatePackageName(pkg string) error {
	if pkg == "" {
		return errors.New("package cannot be empty")
	}
	parts := strings.Split(pkg, ".")
	if len(parts) < 3 {
		return errors.New("package must have at least 3 segments")
	}
	for _, seg := range parts {
		if !pkgSegment.MatchString(seg) {
			return errors.New("invalid package segment: " + seg)
		}
	}
	return nil
}

func SanitizeAppName(name string) string {
	// Convert separators to spaces, split, capitalize, join
	sep := func(r rune) bool {
		return r == '-' || r == '_' || unicode.IsSpace(r)
	}
	parts := strings.FieldsFunc(name, sep)
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}

func PackageToPath(pkg string) string {
	return strings.ReplaceAll(pkg, ".", "/")
}

func SdkVersionLabel(sdk int) string {
	labels := map[int]string{
		21: "Android 5.0 Lollipop",
		22: "Android 5.1 Lollipop",
		23: "Android 6.0 Marshmallow",
		24: "Android 7.0 Nougat",
		25: "Android 7.1 Nougat",
		26: "Android 8.0 Oreo",
		27: "Android 8.1 Oreo",
		28: "Android 9 Pie",
		29: "Android 10",
		30: "Android 11",
		31: "Android 12",
		32: "Android 12L",
		33: "Android 13",
		34: "Android 14",
		35: "Android 14.1",
	}
	if s, ok := labels[sdk]; ok {
		return s
	}
	return "Android API " + string(rune(sdk))
}

// ValidateLanguage checks if language is valid (kotlin|java).
func ValidateLanguage(lang string) error {
	s := strings.ToLower(strings.TrimSpace(lang))
	if s == "" {
		return nil // Empty is OK (uses default)
	}
	switch s {
	case "kotlin", "java":
		return nil
	default:
		return errors.New("invalid language: " + s + " (allowed: kotlin|java)")
	}
}

// ValidateUI checks if UI type is valid (xml|compose).
func ValidateUI(ui string) error {
	s := strings.ToLower(strings.TrimSpace(ui))
	if s == "" || s == "none" {
		return nil // Empty or 'none' is OK (disables UI)
	}
	switch s {
	case "xml", "compose":
		return nil
	default:
		return errors.New("invalid UI type: " + s + " (allowed: xml|compose)")
	}
}

// ValidateDI checks if DI framework is valid (none|hilt|koin).
func ValidateDI(di string) error {
	s := strings.ToLower(strings.TrimSpace(di))
	if s == "" {
		return nil // Empty is OK (uses default)
	}
	switch s {
	case "none", "hilt", "koin":
		return nil
	default:
		return errors.New("invalid DI framework: " + s + " (allowed: none|hilt|koin)")
	}
}

// ValidateClassName checks if a class name is valid Kotlin/Java.
func ValidateClassName(name string) error {
	if name == "" {
		return errors.New("class name cannot be empty")
	}
	if !isValidIdentifier(name) {
		return errors.New("invalid class name: " + name + " (must start with letter and contain only alphanumeric/underscore)")
	}
	if len(name) > 128 {
		return errors.New("class name too long (max 128 characters)")
	}
	return nil
}

// isValidIdentifier checks if a string is a valid Java/Kotlin identifier.
func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	// First character must be letter or underscore
	if !unicode.IsLetter(rune(s[0])) && s[0] != '_' {
		return false
	}
	// Remaining characters must be alphanumeric or underscore
	for _, r := range s[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
