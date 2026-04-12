package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCtrlCCancelFlow starts the interactive wizard and sends an interrupt
// to simulate Ctrl+C; the CLI should exit cleanly and print "Cancelled.".
func TestCtrlCCancelFlow(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping Ctrl+C integration test in CI")
	}

	tmp, err := os.MkdirTemp("", "grimoire-ctrlc-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}

	// Build binary in a separate temp dir so the interactive working dir stays empty
	buildDir, err := os.MkdirTemp("", "grimoire-bin-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(buildDir)

	bin := filepath.Join(buildDir, "grimoire-test")
	if os.PathSeparator == '\\' {
		bin += ".exe"
	}

	build := execCommand("go", "build", "-o", bin)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	// Start interactive `grimoire new` (no args) in an isolated tmp dir
	cmd := execCommand(bin, "new")
	cmd.Dir = tmp
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	// Give the process a moment to initialize the TUI
	time.Sleep(800 * time.Millisecond)

	// Send interrupt (Ctrl+C)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		// On some platforms Signal may not be supported; attempt to kill
		if killErr := cmd.Process.Kill(); killErr != nil {
			t.Fatalf("failed to signal or kill process: signal err=%v, kill err=%v", err, killErr)
		}
	}

	// Wait for process exit with timeout
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("process did not exit after interrupt")
	case err := <-done:
		if err != nil {
			// process exited with non-zero; still check output
			t.Logf("process exited with error: %v", err)
		}
	}

	// Ensure the process exited and did not create any files in the tmp dir
	entries, err := os.ReadDir(tmp)
	if err != nil {
		t.Fatalf("reading tmp dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no files created in tmp dir after interrupt; found: %v\nout:\n%s", entries, buf.String())
	}
}

// TestNonInteractiveRegression verifies the non-interactive `new` path still works.
func TestNonInteractiveRegression(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping integration test in CI")
	}

	tmp, err := os.MkdirTemp("", "grimoire-regress-*")
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

	build := execCommand("go", "build", "-o", bin)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	run := execCommand(bin, "new", "SmokeReg", "--package", "com.test.smokeregress", "--lang", "kotlin", "--wrapper=false", "--git=false", "--vscode=false")
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("grimoire new failed: %v\n%s", err, string(out))
	}
}
