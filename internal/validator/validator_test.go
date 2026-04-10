package validator

import "testing"

func TestValidateAppName(t *testing.T) {
	valid := []string{"MyApp", "my_app", "my-app", "App123"}
	for _, s := range valid {
		if err := ValidateAppName(s); err != nil {
			t.Fatalf("expected valid app name %q: %v", s, err)
		}
	}

	invalid := []string{"", "1Bad", "name!", string(make([]byte, 60))}
	for _, s := range invalid {
		if err := ValidateAppName(s); err == nil {
			t.Fatalf("expected invalid app name %q", s)
		}
	}
}

func TestValidatePackageName(t *testing.T) {
	if err := ValidatePackageName("com.example.app"); err != nil {
		t.Fatalf("valid package failed: %v", err)
	}
	invalid := []string{"Com.Example.App", "com.example", "com.1bad.app", "com.-bad.app"}
	for _, s := range invalid {
		if err := ValidatePackageName(s); err == nil {
			t.Fatalf("expected invalid package %q", s)
		}
	}
}

func TestSanitizeAppName(t *testing.T) {
	if got := SanitizeAppName("my-app"); got != "MyApp" {
		t.Fatalf("sanitize failed: %s", got)
	}
	if got := SanitizeAppName("my_app"); got != "MyApp" {
		t.Fatalf("sanitize failed: %s", got)
	}
}

func TestPackageToPath(t *testing.T) {
	if got := PackageToPath("com.example.myapp"); got != "com/example/myapp" {
		t.Fatalf("package to path failed: %s", got)
	}
}

func TestSdkVersionLabel(t *testing.T) {
	if got := SdkVersionLabel(26); got != "Android 8.0 Oreo" {
		t.Fatalf("sdk label wrong: %s", got)
	}
}
