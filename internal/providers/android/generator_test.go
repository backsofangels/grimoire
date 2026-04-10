package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
)

func TestGenerate_Basic_Kotlin(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "SmokeKotlin")
	cfg := providers.ProviderConfig{
		"AppName":     "SmokeKotlin",
		"PackageName": "com.test.smokekotlin",
		"Lang":        "kotlin",
		"Template":    "basic",
		"OutputDir":   out,
		"Git":         false,
		"Vscode":      true,
		"NoWrapper":   true,
		"MinSdk":      26,
		"TargetSdk":   35,
	}
	if err := GenerateProject(cfg); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}

	manifestPath := filepath.Join(out, "app", "src", "main", "AndroidManifest.xml")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("manifest missing: %v", err)
	}
	// Newer AGP versions expect `namespace` in module build files instead of
	// `package` in the manifest. Verify namespace is present in app's build.gradle.
	appBuild := filepath.Join(out, "app", "build.gradle")
	b, _ := os.ReadFile(appBuild)
	if !strings.Contains(string(b), "namespace \"com.test.smokekotlin\"") {
		t.Fatalf("namespace not in app build.gradle")
	}

	// Kotlin main
	mainKt := filepath.Join(out, "app", "src", "main", "java", "com", "test", "smokekotlin", "MainActivity.kt")
	if _, err := os.Stat(mainKt); err != nil {
		t.Fatalf("kotlin main missing: %v", err)
	}

	// layout
	layout := filepath.Join(out, "app", "src", "main", "res", "layout", "activity_main.xml")
	if _, err := os.Stat(layout); err != nil {
		t.Fatalf("layout missing: %v", err)
	}
}

func TestGenerate_Basic_Java(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "SmokeJava")
	cfg := providers.ProviderConfig{
		"AppName":     "SmokeJava",
		"PackageName": "com.test.smokejava",
		"Lang":        "java",
		"Template":    "basic",
		"OutputDir":   out,
		"Git":         false,
		"Vscode":      true,
		"NoWrapper":   true,
		"MinSdk":      26,
		"TargetSdk":   35,
	}
	if err := GenerateProject(cfg); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}
	mainJava := filepath.Join(out, "app", "src", "main", "java", "com", "test", "smokejava", "MainActivity.java")
	if _, err := os.Stat(mainJava); err != nil {
		t.Fatalf("java main missing: %v", err)
	}
	mainKt := filepath.Join(out, "app", "src", "main", "java", "com", "test", "smokejava", "MainActivity.kt")
	if _, err := os.Stat(mainKt); err == nil {
		t.Fatalf("unexpected kotlin file present")
	}
}

func TestGenerate_Empty(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "SmokeEmpty")
	cfg := providers.ProviderConfig{
		"AppName":     "SmokeEmpty",
		"PackageName": "com.test.smokeempty",
		"Lang":        "kotlin",
		"Template":    "empty",
		"OutputDir":   out,
		"Git":         false,
		"Vscode":      true,
		"NoWrapper":   true,
	}
	if err := GenerateProject(cfg); err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	layout := filepath.Join(out, "app", "src", "main", "res", "layout")
	if _, err := os.Stat(layout); err == nil {
		// layout should not exist for empty template
		t.Fatalf("layout directory present for empty template")
	}
}

func TestGenerate_NoVSCode(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "SmokeNoCode")
	cfg := providers.ProviderConfig{
		"AppName":     "SmokeNoCode",
		"PackageName": "com.test.smokenocode",
		"Lang":        "kotlin",
		"Template":    "basic",
		"OutputDir":   out,
		"Git":         false,
		"Vscode":      false,
		"NoWrapper":   true,
	}
	if err := GenerateProject(cfg); err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	vs := filepath.Join(out, ".vscode")
	if _, err := os.Stat(vs); err == nil {
		t.Fatalf(".vscode should not exist when Vscode=false")
	}
}

func TestGenerate_OutputAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	// create dir to simulate existing output
	if err := os.MkdirAll(filepath.Join(tmp, "existing"), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	cfg := providers.ProviderConfig{
		"AppName":     "Existing",
		"PackageName": "com.test.existing",
		"Lang":        "kotlin",
		"Template":    "basic",
		"OutputDir":   filepath.Join(tmp, "existing"),
		"Git":         false,
		"Vscode":      true,
		"NoWrapper":   true,
	}
	if err := GenerateProject(cfg); err == nil {
		t.Fatalf("expected error when output dir exists")
	}
}
