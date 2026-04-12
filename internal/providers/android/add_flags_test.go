package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

func TestAdd_NoNavDoesNotCreateNavGraph(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	// prepare minimal module structure and files
	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.nonav" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	build := "plugins {\n    id 'com.android.application'\n}"
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte(build), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "NoNavFragment",
		"PackageName": "com.example.nonav",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// nav_graph should NOT exist
	navFile := filepath.Join(module, "src", "main", "res", "navigation", "nav_graph.xml")
	if _, err := os.Stat(navFile); err == nil {
		t.Fatalf("unexpected nav_graph.xml created at %s", navFile)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat nav_graph error: %v", err)
	}

	// build.gradle should be unchanged
	b, err := os.ReadFile(filepath.Join(module, "build.gradle"))
	if err != nil {
		t.Fatalf("read build.gradle: %v", err)
	}
	if string(b) != build {
		t.Fatalf("expected build.gradle to be unchanged; got: %s", string(b))
	}
}

func TestAdd_NoDI_DoesNotCreateApplicationOrDeps(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.nodi" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	build := "plugins {\n    id 'com.android.application'\n}"
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte(build), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "NoDIActivity",
		"PackageName": "com.example.nodi",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.nodi")
	appFile := filepath.Join(module, "src", "main", "java", pkgPath, "MyApplication.kt")
	if _, err := os.Stat(appFile); err == nil {
		t.Fatalf("unexpected MyApplication.kt created: %s", appFile)
	}

	// module build.gradle should not include Hilt or Koin markers
	b, err := os.ReadFile(filepath.Join(module, "build.gradle"))
	if err != nil {
		t.Fatalf("read build.gradle: %v", err)
	}
	s := string(b)
	if strings.Contains(s, "hilt-android") || strings.Contains(s, "hilt-android-gradle-plugin") {
		t.Fatalf("unexpected Hilt markers in build.gradle: %s", s)
	}
	if strings.Contains(s, "io.insert-koin") || strings.Contains(s, "koin-android") {
		t.Fatalf("unexpected Koin markers in build.gradle: %s", s)
	}
}

func TestAdd_ViewModelFlagGeneratesViewModel(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.vmtest" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "VMFrag",
		"PackageName": "com.example.vmtest",
		"Module":      module,
		"Lang":        "kotlin",
		"ViewModel":   true,
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.vmtest")
	vmFile := filepath.Join(module, "src", "main", "java", pkgPath, "VMFragViewModel.kt")
	if _, err := os.Stat(vmFile); err != nil {
		t.Fatalf("expected ViewModel file at %s: %v", vmFile, err)
	}
}

func TestAdd_ViewModelAbsentDoesNotGenerate(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.novm" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "NoVMFrag",
		"PackageName": "com.example.novm",
		"Module":      module,
		"Lang":        "kotlin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.novm")
	vmFile := filepath.Join(module, "src", "main", "java", pkgPath, "NoVMFragViewModel.kt")
	if _, err := os.Stat(vmFile); err == nil {
		t.Fatalf("unexpected ViewModel file created: %s", vmFile)
	}
}
