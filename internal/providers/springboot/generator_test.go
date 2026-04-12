package springboot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

func TestGeneratePlainProject(t *testing.T) {
	tmp := t.TempDir()
	output := filepath.Join(tmp, "myapp-out")
	cfg := providers.ProviderConfig{
		"AppName":     "MyApp",
		"Group":       "com.example",
		"Artifact":    "myapp",
		"PackageName": "com.example.myapp",
		"OutputDir":   output,
		"Template":    "plain",
		"Git":         false,
	}

	if err := GenerateProject(cfg); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}

	files := []string{
		filepath.Join(output, "build.gradle"),
		filepath.Join(output, "settings.gradle"),
		filepath.Join(output, "gradle.properties"),
		filepath.Join(output, "README.md"),
		filepath.Join(output, "src", "main", "java", "com", "example", "myapp", "MyAppMain.java"),
		filepath.Join(output, "src", "main", "resources", "application.properties"),
	}

	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			t.Fatalf("expected file %s to exist: %v", f, err)
		}
		if info.IsDir() {
			t.Fatalf("expected %s to be a file, got directory", f)
		}
	}

	// verify java file content
	javaFile := files[4]
	b, err := os.ReadFile(javaFile)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, "package com.example.myapp;") {
		t.Fatalf("generated java file missing package declaration")
	}
	// class name is based on sanitized app name + "Main"
	expectedClass := validator.SanitizeAppName("MyApp") + "Main"
	if !strings.Contains(s, "public class "+expectedClass) {
		t.Fatalf("generated java file missing expected class declaration; want %s", expectedClass)
	}
}
