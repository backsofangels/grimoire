package validator

import (
	"strings"
	"testing"
)

func TestValidateAppName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple name", "MyApp", false},
		{"valid lowercase with underscore", "my_app", false},
		{"valid lowercase with dash", "my-app", false},
		{"valid with numbers", "App123", false},
		{"invalid empty", "", true},
		{"invalid starts with digit", "1Bad", true},
		{"invalid with special chars", "name!", true},
		{"invalid too long", string(make([]byte, 60)), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAppName(%q) wantErr=%v gotErr=%v", tt.input, tt.wantErr, err)
			}
		})
	}
}

func TestValidatePackageName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid standard package", "com.example.app", false},
		{"invalid uppercase segment", "Com.Example.App", true},
		{"invalid missing segment", "com.example", true},
		{"invalid numeric start in segment", "com.1bad.app", true},
		{"invalid dash in segment", "com.-bad.app", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackageName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePackageName(%q) wantErr=%v gotErr=%v", tt.input, tt.wantErr, err)
			}
		})
	}
}

func TestSanitizeAppName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"dash to camel case", "my-app", "MyApp"},
		{"underscore to camel case", "my_app", "MyApp"},
		{"mixed separators", "my-_app", "MyApp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeAppName(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeAppName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPackageToPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"standard package", "com.example.myapp", "com/example/myapp"},
		{"short package", "com.app", "com/app"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PackageToPath(tt.input)
			if got != tt.expected {
				t.Errorf("PackageToPath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSdkVersionLabel(t *testing.T) {
	tests := []struct {
		name     string
		sdk      int
		expected string
	}{
		{"SDK 26", 26, "Android 8.0 Oreo"},
		{"SDK 33", 33, "Android 13"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SdkVersionLabel(tt.sdk)
			if got != tt.expected {
				t.Errorf("SdkVersionLabel(%d) = %q, want %q", tt.sdk, got, tt.expected)
			}
		})
	}
}

func TestValidateLanguage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty is valid (defaults to kotlin)", "", false},
		{"lowercase kotlin", "kotlin", false},
		{"uppercase kotlin", "KOTLIN", false},
		{"mixed case kotlin", "KoTlIn", false},
		{"lowercase java", "java", false},
		{"uppercase java", "JAVA", false},
		{"invalid swift", "swift", true},
		{"invalid c++", "c++", true},
		{"invalid go", "go", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLanguage(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLanguage(%q) wantErr=%v gotErr=%v", tt.input, tt.wantErr, err)
			}
		})
	}
}

func TestValidateUI(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty is valid", "", false},
		{"lowercase xml", "xml", false},
		{"uppercase xml", "XML", false},
		{"lowercase compose", "compose", false},
		{"uppercase compose", "COMPOSE", false},
		{"lowercase none", "none", false},
		{"uppercase none", "NONE", false},
		{"mixed case none", "NoNe", false},
		{"invalid flutter", "flutter", true},
		{"invalid swiftui", "swiftui", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUI(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUI(%q) wantErr=%v gotErr=%v", tt.input, tt.wantErr, err)
			}
		})
	}
}

func TestValidateDI(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty is valid", "", false},
		{"lowercase none", "none", false},
		{"uppercase none", "NONE", false},
		{"lowercase hilt", "hilt", false},
		{"uppercase hilt", "HILT", false},
		{"mixed case hilt", "HiLt", false},
		{"lowercase koin", "koin", false},
		{"uppercase koin", "KOIN", false},
		{"invalid dagger", "dagger", true},
		{"invalid guice", "guice", true},
		{"invalid spring", "spring", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDI(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDI(%q) wantErr=%v gotErr=%v", tt.input, tt.wantErr, err)
			}
		})
	}
}

func TestValidateClassName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty is invalid", "", true},
		{"simple class name", "MyClass", false},
		{"lowercase class name", "myClass", false},
		{"underscore prefix", "_MyClass", false},
		{"with numbers", "MyClass123", false},
		{"number prefix is invalid", "123MyClass", true},
		{"dash is invalid", "My-Class", true},
		{"space is invalid", "My Class", true},
		{"too long (129 chars)", strings.Repeat("a", 129), true},
		{"max length (128 chars)", strings.Repeat("a", 128), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClassName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClassName(%q) wantErr=%v gotErr=%v", tt.input, tt.wantErr, err)
			}
		})
	}
}
