package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run environment checks (JDK, Gradle, Android SDK)",
	Run: func(cmd *cobra.Command, args []string) {
		fix, _ := cmd.Flags().GetBool("fix")

		logging.Info("→ Starting environment checks")

		// Java
		javaPath, javaErr := exec.LookPath("java")
		if javaErr != nil {
			logging.Error("✗ Missing requirement: Java JDK", "hint", "Install a JDK and ensure 'java' is on your PATH")
		} else {
			out, _ := exec.Command(javaPath, "-version").CombinedOutput()
			ver := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
			logging.Info("✓ Java found", "path", javaPath, "version", ver)
		}
		// JAVA_HOME
		if v, ok := os.LookupEnv("JAVA_HOME"); ok && v != "" {
			logging.Info("✓ Environment variable detected: JAVA_HOME", "value", v)
		} else {
			logging.Warn("! Warning: JAVA_HOME not set", "hint", "Consider setting JAVA_HOME to your JDK root")
		}

		// Gradle
		gradlePath, gradleErr := exec.LookPath("gradle")
		if gradleErr != nil {
			logging.Error("✗ Missing requirement: Gradle CLI", "hint", "Install Gradle and ensure 'gradle' is on your PATH")
		} else {
			out, _ := exec.Command(gradlePath, "--version").CombinedOutput()
			ver := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
			logging.Info("✓ Gradle found", "path", gradlePath, "version", ver)
		}

		// Android SDK / adb
		adbPath, adbErr := exec.LookPath("adb")
		if adbErr != nil {
			// sdkmanager as fallback
			if _, err := exec.LookPath("sdkmanager"); err == nil {
				logging.Info("✓ Android SDK tools found (sdkmanager)")
			} else {
				logging.Warn("! Warning: Android SDK not found", "hint", "Install Android SDK and add platform-tools to PATH (adb)")
			}
		} else {
			logging.Info("✓ adb found", "path", adbPath)
		}

		// ANDROID_HOME / ANDROID_SDK_ROOT
		if v, ok := os.LookupEnv("ANDROID_HOME"); ok && v != "" {
			logging.Info("✓ Environment variable detected: ANDROID_HOME", "value", v)
		} else if v, ok := os.LookupEnv("ANDROID_SDK_ROOT"); ok && v != "" {
			logging.Info("✓ Environment variable detected: ANDROID_SDK_ROOT", "value", v)
		} else {
			logging.Warn("! Warning: Android SDK env vars not set", "hint", "Set ANDROID_SDK_ROOT or ANDROID_HOME to your SDK root")
		}

		// Provider-specific checks
		providerName, _ := cmd.Flags().GetString("provider")
		if providerName == "" {
			providerName = "android"
		}
		if p, err := providers.Get(providerName); err == nil {
			if checks := p.DoctorChecks(); checks != nil {
				for _, c := range checks {
					if c.Run != nil {
						if err := c.Run(); err != nil {
							logging.Error("✗ Check failed", "name", c.Name, "error", err)
							if c.Fix != nil && fix {
								if err := c.Fix(); err != nil {
									logging.Error("✗ Automatic fix failed", "name", c.Name, "error", err)
								} else {
									logging.Info("✓ Automatic fix applied", "name", c.Name)
								}
							}
						} else {
							logging.Info("✓ Check passed", "name", c.Name)
						}
					}
				}
			}
		}

		// Attempt non-destructive fixes if requested (Windows only for now)
		if fix {
			if runtime.GOOS == "windows" {
				// Set JAVA_HOME from java.exe location if missing
				if _, ok := os.LookupEnv("JAVA_HOME"); !ok {
					if javaErr == nil {
						cand := filepath.Dir(filepath.Dir(javaPath))
						if err := runSetx("JAVA_HOME", cand); err == nil {
							logging.Info("✓ Set: JAVA_HOME", "value", cand)
						} else {
							logging.Error("✗ Error: unable to set JAVA_HOME", "error", err)
						}
					}
				}
				// Set ANDROID_SDK_ROOT / ANDROID_HOME from adb location if missing
				if _, ok := os.LookupEnv("ANDROID_SDK_ROOT"); !ok {
					if adbErr == nil {
						sdkroot := filepath.Dir(filepath.Dir(adbPath))
						if err := runSetx("ANDROID_SDK_ROOT", sdkroot); err == nil {
							logging.Info("✓ Set: ANDROID_SDK_ROOT", "value", sdkroot)
						} else {
							logging.Error("✗ Error: unable to set ANDROID_SDK_ROOT", "error", err)
						}
						// also set ANDROID_HOME for compatibility
						if err := runSetx("ANDROID_HOME", sdkroot); err == nil {
							logging.Info("✓ Set: ANDROID_HOME", "value", sdkroot)
						}
					}
				}
			} else {
				// Attempt non-destructive fixes for Unix-like systems by appending
				// `export` lines to a likely shell profile (e.g. ~/.bashrc or ~/.zshrc).
				// These are best-effort and will not overwrite existing files.

				// JAVA_HOME
				if _, ok := os.LookupEnv("JAVA_HOME"); !ok {
					if javaErr == nil {
						// try to detect java home via `java -XshowSettings:properties -version`
						if cand, derr := detectJavaHome(javaPath); derr == nil && cand != "" {
							if file, err := runShellExport("JAVA_HOME", cand); err == nil {
								logging.Info("✓ Created: JAVA_HOME added to profile file", "file", file, "value", cand)
							} else {
								logging.Error("✗ Error: unable to append JAVA_HOME to profile", "error", err)
							}
						} else if cand2 := filepath.Dir(filepath.Dir(javaPath)); cand2 != "" {
							if file, err := runShellExport("JAVA_HOME", cand2); err == nil {
								logging.Info("✓ Created: JAVA_HOME added to profile file", "file", file, "value", cand2)
							} else {
								logging.Error("✗ Error: unable to append JAVA_HOME to profile", "error", err)
							}
						}
					}
				}

				// ANDROID_SDK_ROOT / ANDROID_HOME
				if _, ok := os.LookupEnv("ANDROID_SDK_ROOT"); !ok {
					if adbErr == nil {
						cand := filepath.Dir(filepath.Dir(adbPath))
						if info, err := os.Stat(cand); err == nil && info.IsDir() {
							if file, err := runShellExport("ANDROID_SDK_ROOT", cand); err == nil {
								logging.Info("✓ Created: ANDROID_SDK_ROOT added to profile file", "file", file, "value", cand)
							} else {
								logging.Error("✗ Error: unable to append ANDROID_SDK_ROOT to profile", "error", err)
							}
							if file, err := runShellExport("ANDROID_HOME", cand); err == nil {
								logging.Info("✓ Created: ANDROID_HOME added to profile file", "file", file, "value", cand)
							}
						}
					}
				}

				logging.Info("→ Fixes complete. Open a new shell or source the profile to apply changes.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().Bool("fix", false, "Attempt non-destructive fixes (set env vars)")
}

func runSetx(name, value string) error {
	setxPath, err := exec.LookPath("setx")
	if err != nil {
		return fmt.Errorf("setx not found: %w", err)
	}
	cmd := exec.Command(setxPath, name, value)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("setx failed: %v - %s", err, string(out))
	}
	return nil
}

// detectJavaHome attempts to determine the Java home directory by running
// `java -XshowSettings:properties -version` and parsing the `java.home` value.
// It is used on Unix-like systems to find a reliable JAVA_HOME candidate.
func detectJavaHome(javaPath string) (string, error) {
	if javaPath == "" {
		return "", fmt.Errorf("java path empty")
	}
	cmd := exec.Command(javaPath, "-XshowSettings:properties", "-version")
	out, _ := cmd.CombinedOutput() // ignore error; some JDKs print to stderr
	if len(out) == 0 {
		return "", fmt.Errorf("no output from java command")
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "java.home =") || strings.HasPrefix(line, "java.home=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("java.home not found in java output")
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

	// Pick the first existing profile, or create the first candidate if none exist.
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

	exportLine := fmt.Sprintf("\n# grimoire: set %s\nexport %s=\"%s\"\n", name, name, value)
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
