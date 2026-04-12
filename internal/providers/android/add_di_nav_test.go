package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

// Test that adding a fragment with Hilt + ViewModel + Nav wiring creates
// the expected Application class, annotations, ViewModel modifications,
// and a navigation graph entry.
func TestAddWithHiltAndNavCreatesHiltAndNavFiles(t *testing.T) {
	tmp := t.TempDir()
	module := filepath.Join(tmp, "app")

	// prepare minimal module structure and files
	if err := os.MkdirAll(filepath.Join(module, "src", "main"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `<?xml version="1.0" encoding="utf-8"?>\n<manifest package="com.example.hilt" />`
	if err := os.WriteFile(filepath.Join(module, "src", "main", "AndroidManifest.xml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	// minimal build.gradle (module) so wiring code can append plugins/deps
	build := `plugins {\n    id 'com.android.application'\n}`
	if err := os.WriteFile(filepath.Join(module, "build.gradle"), []byte(build), 0o644); err != nil {
		t.Fatalf("write build.gradle: %v", err)
	}

	cfg := providers.ProviderConfig{
		"Kind":        "fragment",
		"Name":        "HiltFragment",
		"PackageName": "com.example.hilt",
		"Module":      module,
		"Lang":        "kotlin",
		"ViewModel":   true,
		"DI":          "hilt",
		"Nav":         true,
	}

	p := &AndroidProvider{}
	if err := p.Add(cfg); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	pkgPath := validator.PackageToPath("com.example.hilt")

	// MyApplication.kt with Hilt annotation
	appFile := filepath.Join(module, "src", "main", "java", pkgPath, "MyApplication.kt")
	b, err := os.ReadFile(appFile)
	if err != nil {
		t.Fatalf("expected MyApplication.kt: %v", err)
	}
	if !strings.Contains(string(b), "HiltAndroidApp") {
		t.Fatalf("MyApplication.kt missing @HiltAndroidApp: %s", string(b))
	}

	// Fragment should have @AndroidEntryPoint
	fragFile := filepath.Join(module, "src", "main", "java", pkgPath, "HiltFragment.kt")
	b, err = os.ReadFile(fragFile)
	if err != nil {
		t.Fatalf("expected fragment source: %v", err)
	}
	if !strings.Contains(string(b), "@AndroidEntryPoint") {
		t.Fatalf("fragment source missing @AndroidEntryPoint: %s", string(b))
	}

	// ViewModel: Hilt annotations and Inject
	vmFile := filepath.Join(module, "src", "main", "java", pkgPath, "HiltFragmentViewModel.kt")
	b, err = os.ReadFile(vmFile)
	if err != nil {
		t.Fatalf("expected viewmodel source: %v", err)
	}
	sb := string(b)
	if !strings.Contains(sb, "HiltViewModel") || !strings.Contains(sb, "@Inject") {
		t.Fatalf("viewmodel missing Hilt annotations or Inject: %s", sb)
	}

	// Navigation graph created and contains fragment entry
	navFile := filepath.Join(module, "src", "main", "res", "navigation", "nav_graph.xml")
	b, err = os.ReadFile(navFile)
	if err != nil {
		t.Fatalf("expected nav_graph.xml: %v", err)
	}
	ns := string(b)
	if !strings.Contains(ns, "com.example.hilt.HiltFragment") {
		t.Fatalf("nav_graph.xml missing fragment name: %s", ns)
	}
	if !strings.Contains(ns, "@+id/hilt_fragment") {
		t.Fatalf("nav_graph.xml missing fragment id: %s", ns)
	}

	// module build.gradle should contain Hilt and Navigation dep markers
	modBuild := filepath.Join(module, "build.gradle")
	b, err = os.ReadFile(modBuild)
	if err != nil {
		t.Fatalf("read module build.gradle: %v", err)
	}
	ms := string(b)
	if !strings.Contains(ms, "hilt-android") && !strings.Contains(ms, "hilt-android-gradle-plugin") {
		t.Fatalf("module build.gradle missing Hilt setup: %s", ms)
	}
	if !strings.Contains(ms, "navigation-fragment-ktx") {
		t.Fatalf("module build.gradle missing navigation deps: %s", ms)
	}
}
