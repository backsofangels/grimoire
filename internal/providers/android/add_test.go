package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

// simple toSnake helper to match provider's behavior
func toSnakeTest(s string) string {
	var parts []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			parts = append(parts, '_')
		}
		parts = append(parts, r)
	}
	return strings.ToLower(string(parts))
}

func TestAddActivityCreatesFiles(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "MyActivity",
		"PackageName": "com.example.test",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.test")
	src := filepath.Join(module, "src", "main", "java", pkgPath, "MyActivity.kt")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("expected activity source at %s: %v", src, err)
	}

	layoutName := "activity_" + toSnakeTest("MyActivity")
	layout := filepath.Join(module, "src", "main", "res", "layout", layoutName+".xml")
	if _, err := os.Stat(layout); err != nil {
		t.Fatalf("expected layout at %s: %v", layout, err)
	}
}

func TestAddFragmentCreatesFiles(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "DetailsFragment",
		"PackageName": "com.example.fragment",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add fragment failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.fragment")
	src := filepath.Join(module, "src", "main", "java", pkgPath, "DetailsFragment.kt")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("expected fragment source at %s: %v", src, err)
	}

	layoutName := "fragment_" + toSnakeTest("DetailsFragment")
	layout := filepath.Join(module, "src", "main", "res", "layout", layoutName+".xml")
	if _, err := os.Stat(layout); err != nil {
		t.Fatalf("expected fragment layout at %s: %v", layout, err)
	}
}

func TestAddViewModelCreatesFile(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	cfg := providers.ProviderConfig{
		"Kind":        "viewmodel",
		"Name":        "MainViewModel",
		"PackageName": "com.example.vm",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add viewmodel failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.vm")
	src := filepath.Join(module, "src", "main", "java", pkgPath, "MainViewModel.kt")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("expected viewmodel source at %s: %v", src, err)
	}
}

func TestAddDetectsPackageFromManifest(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")
	manifestDir := filepath.Join(module, "src", "main")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		t.Fatalf("mkdir manifest dir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>
<manifest package="com.detect.pkg" />`
	if err := os.WriteFile(filepath.Join(manifestDir, "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":   "activity",
		"Name":   "DetectedActivity",
		"Module": module,
		"Lang":   "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed when detecting package: %v", err)
	}

	pkgPath := validator.PackageToPath("com.detect.pkg")
	src := filepath.Join(module, "src", "main", "java", pkgPath, "DetectedActivity.kt")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("expected detected activity source at %s: %v", src, err)
	}
}

func TestAddRejectsInvalidClassName(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "1BadName",
		"PackageName": "com.example.invalid",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err == nil {
		t.Fatalf("expected error for invalid class name")
	}
}

func TestAddExistingFileNoOverride(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")
	pkgPath := validator.PackageToPath("com.example.exists")
	dstDir := filepath.Join(module, "src", "main", "java", pkgPath)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	file := filepath.Join(dstDir, "ExistsActivity.kt")
	if err := os.WriteFile(file, []byte("// existing"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "ExistsActivity",
		"PackageName": "com.example.exists",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err == nil {
		t.Fatalf("expected error when target file exists and override not set")
	}
}

func TestAddExistingFileWithOverride(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")
	pkgPath := validator.PackageToPath("com.example.exists2")
	dstDir := filepath.Join(module, "src", "main", "java", pkgPath)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	file := filepath.Join(dstDir, "ExistsActivity.kt")
	if err := os.WriteFile(file, []byte("// existing"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "ExistsActivity",
		"PackageName": "com.example.exists2",
		"Module":      module,
		"Lang":        "kotlin",
		"Override":    true,
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add with override failed: %v", err)
	}

	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read file after override: %v", err)
	}
	if strings.Contains(string(b), "// existing") {
		t.Fatalf("expected file to be overwritten")
	}
}
