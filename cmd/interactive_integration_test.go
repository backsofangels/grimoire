package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/creack/pty"
)

// TestInteractiveNew_SelectSpringBoot runs the interactive TUI, selects
// Java -> Spring Boot and answers prompts programmatically via a PTY.
func TestInteractiveNew_SelectSpringBoot(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping interactive integration test in CI")
	}

	tmp, err := os.MkdirTemp("", "grimoire-ptyin-*")
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

	build := exec.Command("go", "build", "-o", bin)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	// Start interactive binary under a PTY
	cmd := exec.Command(bin, "new")
	cmd.Dir = tmp
	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Skipf("pty start failed (skipping interactive test): %v", err)
	}
	defer func() { _ = ptmx.Close() }()

	// Helper to write and flush
	write := func(s string) {
		_, _ = io.WriteString(ptmx, s)
		time.Sleep(250 * time.Millisecond)
	}

	// Step 1: project type - press Down then Enter to select 'java'
	write("\x1b[B")
	write("\r")

	// SpringBoot prompts: App name, Group, Artifact, Package (blank), Template (Enter), Output dir (Enter), Git (n), Confirm (Enter)
	write("MySpringApp\r")
	write("com.test\r")
	write("myspring\r")
	// package: leave blank
	write("\r")
	// template: accept default
	write("\r")
	// output dir: leave blank
	write("\r")
	// git: send 'n' then Enter to disable
	write("n\r")
	// Confirm / create
	write("\r")

	// Wait for process to finish (with timeout)
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case <-time.After(20 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("process did not exit in time")
	case err := <-done:
		if err != nil {
			t.Fatalf("process exited with error: %v", err)
		}
	}

	// Verify output
	outDir := filepath.Join(tmp, "myspring")
	if _, err := os.Stat(filepath.Join(outDir, "build.gradle")); err != nil {
		t.Fatalf("expected build.gradle in %s: %v", outDir, err)
	}
	base := validator.SanitizeAppName("MySpringApp")
	appClass := base + "Application"
	javaPath := filepath.Join(outDir, "src", "main", "java", "com", "test", "myspring", appClass+".java")
	if _, err := os.Stat(javaPath); err != nil {
		t.Fatalf("expected generated java main class: %v", err)
	}
}
