package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

func TestJavaOverwriteBehavior_NoOverride(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	pkg := "com.example.jover"
	dstDir := filepath.Join(module, "src", "main", "java", validator.PackageToPath(pkg))
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	file := filepath.Join(dstDir, "ExistsActivity.java")
	if err := os.WriteFile(file, []byte("// original marker"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "ExistsActivity",
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "java",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err == nil {
		t.Fatalf("expected error when target exists and override not set")
	}
}

func TestJavaOverwriteBehavior_WithOverride(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	pkg := "com.example.jover2"
	dstDir := filepath.Join(module, "src", "main", "java", validator.PackageToPath(pkg))
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	file := filepath.Join(dstDir, "ExistsActivity.java")
	if err := os.WriteFile(file, []byte("// original marker"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "activity",
		"Name":        "ExistsActivity",
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "java",
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
	if strings.Contains(string(b), "// original marker") {
		t.Fatalf("expected file to be overwritten")
	}
}

func TestNavDuplicationPrevention(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.dupnav" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte("plugins {\n    id 'com.android.application'\n}"), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "DupNav",
		"PackageName": "com.example.dupnav",
		"Module":      module,
		"Lang":        "kotlin",
		"Nav":         true,
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}
	// allow overwrite on second add so the provider proceeds to nav wiring
	cfg["Override"] = true
	if err := p.Add(cfg); err != nil {
		t.Fatalf("second Add failed: %v", err)
	}

	navFile := filepath.Join(module, "src", "main", "res", "navigation", "nav_graph.xml")
	b, err := os.ReadFile(navFile)
	if err != nil {
		t.Fatalf("expected nav_graph.xml: %v", err)
	}
	s := string(b)
	// ensure only one fragment entry for the class
	nameEntry := "android:name=\"com.example.dupnav.DupNav\""
	if cnt := strings.Count(s, nameEntry); cnt != 1 {
		t.Fatalf("expected exactly one nav entry for class, found %d: %s", cnt, s)
	}
	// ensure id occurs only once
	id := toSnake("DupNav")
	if cnt := strings.Count(s, "@+id/"+id); cnt != 1 {
		t.Fatalf("expected exactly one nav id for fragment, found %d: %s", cnt, s)
	}
}

func TestKoinApplicationNotOverwritten(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	pkg := "com.example.koinkeep"
	appDir := filepath.Join(module, "src", "main", "java", validator.PackageToPath(pkg))
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	appFile := filepath.Join(appDir, "MyApplication.kt")
	orig := "// preserved application"
	if err := os.WriteFile(appFile, []byte(orig), 0o644); err != nil {
		t.Fatalf("write app file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(`<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.koinkeep" />`), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte("plugins {\n    id 'com.android.application'\n}"), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "KoinFragKeep",
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "kotlin",
		"DI":          "koin",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	b, err := os.ReadFile(appFile)
	if err != nil {
		t.Fatalf("read app file: %v", err)
	}
	if string(b) != orig {
		t.Fatalf("expected MyApplication.kt to be unchanged, got: %s", string(b))
	}
}

func TestJavaIncludeViewModelGeneration(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(`<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.javavmfrag" />`), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	pkg := "com.example.javavmfrag"
	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "JavaVMFrag",
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "java",
		"ViewModel":   true,
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath(pkg)
	vmFile := filepath.Join(module, "src", "main", "java", pkgPath, "JavaVMFragViewModel.java")
	if _, err := os.Stat(vmFile); err != nil {
		t.Fatalf("expected Java ViewModel file at %s: %v", vmFile, err)
	}
}
