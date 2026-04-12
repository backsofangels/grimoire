package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// findRepoRoot walks up until it finds go.mod
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

// execCommand abstracts exec.Command for easier testing/mocking.
func execCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func TestNewIntegrationCLI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping integration test in CI")
	}

	tmp, err := os.MkdirTemp("", "grimoire-integ-*")
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

	cmd := execCommand("go", "build", "-o", bin)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	run := execCommand(bin, "new", "SmokeCLI2", "--package", "com.test.smokecli2", "--lang", "kotlin", "--no-wrapper", "--git=false", "--vscode=false")
	run.Dir = tmp
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("grimoire new failed: %v\n%s", err, string(out))
	}

	if _, err := os.Stat(filepath.Join(tmp, "SmokeCLI2", "build.gradle")); err != nil {
		t.Fatalf("expected project build.gradle: %v", err)
	}
}
