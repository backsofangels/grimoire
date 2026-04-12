package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

func TestJavaHiltNavDuplicationPrevention(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.jhiltdup" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte("plugins {\n    id 'com.android.application'\n}"), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	pkg := "com.example.jhiltdup"
	name := "HiltDupJava"
	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        name,
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "java",
		"DI":          "hilt",
		"Nav":         true,
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}
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
	fq := pkg + "." + name
	if strings.Count(s, fq) != 1 {
		t.Fatalf("expected exactly one nav entry for %s, got: %s", fq, s)
	}
	id := toSnake(name)
	if strings.Count(s, "@+id/"+id) != 1 {
		t.Fatalf("expected exactly one nav id for %s, got: %s", id, s)
	}

	// Java fragment should include @AndroidEntryPoint import/annotation
	frag := filepath.Join(module, "src", "main", "java", validator.PackageToPath(pkg), name+".java")
	fb, err := os.ReadFile(frag)
	if err != nil {
		t.Fatalf("expected fragment source: %v", err)
	}
	if !strings.Contains(string(fb), "AndroidEntryPoint") {
		t.Fatalf("fragment source missing AndroidEntryPoint annotation/import: %s", string(fb))
	}
}

func TestJavaHiltOverwriteBehavior(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	pkg := "com.example.jhiltover"
	dstDir := filepath.Join(module, "src", "main", "java", validator.PackageToPath(pkg))
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	file := filepath.Join(dstDir, "ExistsFrag.java")
	if err := os.WriteFile(file, []byte("// original marker"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(`<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.jhiltover" />`), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte("plugins {\n    id 'com.android.application'\n}"), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "ExistsFrag",
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "java",
		"DI":          "hilt",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err == nil {
		t.Fatalf("expected error when target exists and override not set")
	}

	cfg["Override"] = true
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
	if !strings.Contains(string(b), "AndroidEntryPoint") {
		t.Fatalf("expected AndroidEntryPoint to be added to Java source: %s", string(b))
	}
}

func TestJavaHiltViewModelInjection(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	pkg := "com.example.jhiltvm"
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(`<?xml version="1.0" encoding="utf-8"?>\n<manifest package="`+pkg+`" />`), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte("plugins {\n    id 'com.android.application'\n}"), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "viewmodel",
		"Name":        "JHiltVM",
		"PackageName": pkg,
		"Module":      module,
		"Lang":        "java",
		"DI":          "hilt",
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath(pkg)
	vmFile := filepath.Join(module, "src", "main", "java", pkgPath, "JHiltVM.java")
	b, err := os.ReadFile(vmFile)
	if err != nil {
		t.Fatalf("expected java viewmodel source: %v", err)
	}
	sb := string(b)
	if !strings.Contains(sb, "HiltViewModel") || !strings.Contains(sb, "@Inject") {
		t.Fatalf("viewmodel missing Hilt annotations or Inject: %s", sb)
	}
}
