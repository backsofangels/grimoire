package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

func TestAdd_JavaCreatesJavaSources(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.javapkg" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "JavaActivity",
		"PackageName": "com.example.javapkg",
		"Module":      module,
		"Lang":        "java",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.javapkg")
	src := filepath.Join(module, "src", "main", "java", pkgPath, "JavaActivity.java")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("expected java activity source at %s: %v", src, err)
	}
}

func TestAdd_KoinSetupCreatesKoinApplication(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.kointest" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	// minimal module build.gradle so ensureKoinSetup can append deps
	build := "plugins {\n    id 'com.android.application'\n}"
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte(build), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "KoinFragment",
		"PackageName": "com.example.kointest",
		"Module":      module,
		"Lang":        "kotlin",
		"DI":          "koin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// MyApplication.kt should be created with startKoin
	pkgPath := validator.PackageToPath("com.example.kointest")
	appFile := filepath.Join(module, "src", "main", "java", pkgPath, "MyApplication.kt")
	b, err := os.ReadFile(appFile)
	if err != nil {
		t.Fatalf("expected MyApplication.kt: %v", err)
	}
	if !strings.Contains(string(b), "startKoin") {
		t.Fatalf("MyApplication.kt missing startKoin() content: %s", string(b))
	}

	// module build.gradle should include Koin marker
	mb, err := os.ReadFile(filepath.Join(module, "build.gradle"))
	if err != nil {
		t.Fatalf("read module build.gradle: %v", err)
	}
	if !strings.Contains(string(mb), "koin-android") {
		t.Fatalf("module build.gradle missing Koin deps: %s", string(mb))
	}

	// fragment source should exist
	frag := filepath.Join(module, "src", "main", "java", pkgPath, "KoinFragment.kt")
	if _, err := os.Stat(frag); err != nil {
		t.Fatalf("expected fragment source at %s: %v", frag, err)
	}
}
