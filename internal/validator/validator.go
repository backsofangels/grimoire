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
