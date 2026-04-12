package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/backsofangels/grimoire/internal/validator"
)

// TestNewIntegrationCLISpringBoot verifies non-interactive Spring Boot project generation.
func TestNewIntegrationCLISpringBoot(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping integration test in CI")
	}

	tmp, err := os.MkdirTemp("", "grimoire-integ-spring-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}

	bin := filepath.Join(tmp, "grimoire-test")
	if os.PathSeparator == '\\' {
		bin += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", bin)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	// Run non-interactive generation for springboot provider
	run := exec.Command(bin, "new", "MySpringApp", "--provider", "springboot", "--group", "com.test", "--artifact", "myspring", "--package", "com.test.myspring", "--template", "springboot", "--no-wrapper", "--git=false")
	run.Dir = tmp
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("grimoire new failed: %v\n%s", err, string(out))
	}

	outDir := filepath.Join(tmp, "myspring")
	if _, err := os.Stat(filepath.Join(outDir, "build.gradle")); err != nil {
		t.Fatalf("expected project build.gradle: %v", err)
	}

	// verify generated main class exists
	base := validator.SanitizeAppName("MySpringApp")
	appClass := base + "Application"
	javaPath := filepath.Join(outDir, "src", "main", "java", "com", "test", "myspring", appClass+".java")
	if _, err := os.Stat(javaPath); err != nil {
		t.Fatalf("expected generated java main class: %v", err)
	}
}
