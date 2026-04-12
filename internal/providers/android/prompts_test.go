package android

import (
	"strconv"
	"testing"

	"github.com/backsofangels/grimoire/internal/config"
)

func TestInitialStateDefaults(t *testing.T) {
	cfg := config.DefaultConfig()
	appName, pkg, out, lang, minSdkStr, tmpl, git, vscode := initialState(cfg)
	if appName != "" {
		t.Fatalf("expected empty appName, got %q", appName)
	}
	if pkg != cfg.DefaultPackage+".myapp" {
		t.Fatalf("expected package %q, got %q", cfg.DefaultPackage+".myapp", pkg)
	}
	if out != "" {
		t.Fatalf("expected empty outputDir, got %q", out)
	}
	if lang != cfg.DefaultLang {
		t.Fatalf("expected lang %q, got %q", cfg.DefaultLang, lang)
	}
	if minSdkStr != strconv.Itoa(cfg.DefaultMinSdk) {
		t.Fatalf("expected minSdk %q, got %q", strconv.Itoa(cfg.DefaultMinSdk), minSdkStr)
	}
	if tmpl != cfg.DefaultTemplate {
		t.Fatalf("expected template %q, got %q", cfg.DefaultTemplate, tmpl)
	}
	if git != cfg.Git {
		t.Fatalf("expected git %v, got %v", cfg.Git, git)
	}
	if vscode != cfg.VSCode {
		t.Fatalf("expected vscode %v, got %v", cfg.VSCode, vscode)
	}
}

func TestBoolLabel(t *testing.T) {
	if boolLabel(true) != "yes" {
		t.Fatal("boolLabel(true) should return 'yes'")
	}
	if boolLabel(false) != "no" {
		t.Fatal("boolLabel(false) should return 'no'")
	}
}

func TestMustAtoi(t *testing.T) {
	if got := mustAtoi(""); got != 26 {
		t.Fatalf(`mustAtoi("") => %d, want 26`, got)
	}
	if got := mustAtoi("not-a-number"); got != 26 {
		t.Fatalf("mustAtoi('not-a-number') => %d, want 26", got)
	}
	if got := mustAtoi("35"); got != 35 {
		t.Fatalf("mustAtoi('35') => %d, want 35", got)
	}
}
