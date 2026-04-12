package springboot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/backsofangels/grimoire/internal/providers"
)

// SpringBootProvider is a minimal skeleton implementation of the Provider
// interface. It is intentionally lightweight and acts as a registration
// point for future implementation of generation and prompts.
type SpringBootProvider struct{}

func (s *SpringBootProvider) Name() string { return "springboot" }

func (s *SpringBootProvider) Description() string {
	return "Spring Boot project provider (Java)"
}

func (s *SpringBootProvider) Flags() []providers.ProviderFlag {
	return []providers.ProviderFlag{
		{Name: "group", Short: "", Usage: "Group ID (e.g. com.example)", Default: "com.example"},
		{Name: "artifact", Short: "", Usage: "Artifact ID / module name", Default: "app"},
		{Name: "package", Short: "", Usage: "Base package name (overrides group+artifact)", Default: ""},
		{Name: "lang", Short: "", Usage: "Language (java)", Default: "java"},
		{Name: "template", Short: "", Usage: "Template (springboot|plain)", Default: "springboot"},
		{Name: "build", Short: "", Usage: "Build system (gradle|maven)", Default: "gradle"},
		{Name: "git", Short: "", Usage: "Initialize git", Default: true},
		{Name: "output-dir", Short: "o", Usage: "Output directory", Default: ""},
	}
}

// Prompt: interactive prompts will be added separately. Return an empty
// config for now so the cmd layer can continue to call Prompt() safely.
func (s *SpringBootProvider) Prompt() (providers.ProviderConfig, error) {
	return RunPrompt()
}

func (s *SpringBootProvider) Validate(cfg providers.ProviderConfig) error {
	// Provider-level validation will be implemented later.
	return nil
}

func (s *SpringBootProvider) Generate(cfg providers.ProviderConfig) error {
	return GenerateProject(cfg)
}

// Add creates a single resource inside an existing project (not applicable
// to the Spring Boot provider yet).
func (s *SpringBootProvider) Add(cfg providers.ProviderConfig) error {
	return fmt.Errorf("add: not supported for springboot provider")
}

func (s *SpringBootProvider) DoctorChecks() []providers.Check {
	checks := []providers.Check{
		{
			Name: "JDK available (javac)",
			Run: func() error {
				if _, err := exec.LookPath("javac"); err != nil {
					return fmt.Errorf("javac not found in PATH: install a JDK and ensure 'javac' is on your PATH")
				}
				return nil
			},
			Fix: nil,
		},
		{
			Name: "JAVA_HOME points to JDK",
			Run: func() error {
				v, ok := os.LookupEnv("JAVA_HOME")
				if !ok || v == "" {
					return fmt.Errorf("JAVA_HOME not set")
				}
				javac := filepath.Join(v, "bin", "javac")
				if _, err := os.Stat(javac); err != nil {
					return fmt.Errorf("JAVA_HOME does not appear to point to a JDK (missing %s)", javac)
				}
				return nil
			},
			Fix: func() error {
				// Attempt to detect a reasonable JAVA_HOME and set it non-destructively.
				// Prefer 'javac' location, fall back to 'java'.
				var binPath string
				if p, err := exec.LookPath("javac"); err == nil {
					binPath = p
				} else if p2, err2 := exec.LookPath("java"); err2 == nil {
					binPath = p2
				} else {
					return fmt.Errorf("cannot detect JDK: please install a JDK and ensure 'java' or 'javac' is on PATH")
				}

				cand := filepath.Dir(filepath.Dir(binPath)) // go up two levels from bin/java or bin/javac

				if runtime.GOOS == "windows" {
					// Use setx to set JAVA_HOME for the current user.
					setxPath, err := exec.LookPath("setx")
					if err != nil {
						return fmt.Errorf("setx not found: unable to set JAVA_HOME automatically on Windows; please set it manually to %s", cand)
					}
					cmd := exec.Command(setxPath, "JAVA_HOME", cand)
					if out, err := cmd.CombinedOutput(); err != nil {
						return fmt.Errorf("setx failed: %v - %s", err, string(out))
					}
					return nil
				}

				// Unix-like: append export to a likely shell profile (non-destructive)
				file, err := runShellExport("JAVA_HOME", cand)
				if err != nil {
					return fmt.Errorf("unable to append JAVA_HOME to profile: %w", err)
				}
				_ = file
				return nil
			},
		},
		{
			Name: "JDK version >= 11",
			Run: func() error {
				out, err := exec.Command("javac", "-version").CombinedOutput()
				if err != nil {
					return fmt.Errorf("unable to run 'javac -version': %w", err)
				}
				txt := strings.TrimSpace(string(out))
				// Expected: 'javac X.Y.Z' — parse X
				parts := strings.Fields(txt)
				if len(parts) < 2 {
					return fmt.Errorf("unexpected javac version output: %s", txt)
				}
				ver := parts[1]
				majorStr := strings.SplitN(ver, ".", 2)[0]
				major, err := strconv.Atoi(majorStr)
				if err != nil {
					return fmt.Errorf("cannot parse javac version: %s", ver)
				}
				if major < 11 {
					return fmt.Errorf("javac version %d detected; Grimoire requires JDK 11 or newer", major)
				}
				return nil
			},
			Fix: nil,
		},
	}
	return checks
}

// runShellExport appends a non-destructive `export` line to a likely shell
// profile (e.g. ~/.bashrc or ~/.zshrc). Returns the file that was written.
func runShellExport(name, value string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	shell := os.Getenv("SHELL")
	var candidates []string
	if strings.HasSuffix(shell, "zsh") {
		candidates = []string{filepath.Join(home, ".zshrc"), filepath.Join(home, ".zprofile"), filepath.Join(home, ".profile")}
	} else if strings.HasSuffix(shell, "bash") {
		candidates = []string{filepath.Join(home, ".bashrc"), filepath.Join(home, ".bash_profile"), filepath.Join(home, ".profile")}
	} else {
		candidates = []string{filepath.Join(home, ".profile"), filepath.Join(home, ".bashrc"), filepath.Join(home, ".zshrc")}
	}

	var target string
	for _, f := range candidates {
		if _, err := os.Stat(f); err == nil {
			target = f
			break
		}
	}
	if target == "" {
		target = candidates[0]
		if err := os.WriteFile(target, []byte{}, 0644); err != nil {
			return "", err
		}
	}

	b, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	content := string(b)
	marker := fmt.Sprintf("# grimoire: set %s", name)
	if strings.Contains(content, marker) || strings.Contains(content, "export "+name) || strings.Contains(content, name+"=") {
		// already set, no-op
		return target, nil
	}

	exportLine := fmt.Sprintf("\n%s\nexport %s=\"%s\"\n", marker, name, value)
	f, err := os.OpenFile(target, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(exportLine); err != nil {
		return "", err
	}
	return target, nil
}

func init() {
	providers.Register(&SpringBootProvider{})
}
