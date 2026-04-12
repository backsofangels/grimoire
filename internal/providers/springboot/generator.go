package springboot

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

func renderTemplate(name string, data any) (string, error) {
	path := filepath.ToSlash(filepath.Join("templates", name))
	b, err := templateFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", name, err)
	}
	tmpl, err := template.New(name).Parse(string(b))
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %s: %w", name, err)
	}
	return buf.String(), nil
}

func writeFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

func initGit(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		// Non-fatal
		return fmt.Errorf("git init failed: %w", err)
	}
	return nil
}

// resolveGroupArtifact applies derivation rules for group/artifact/package.
// Returns group, artifact, package, or error if package cannot be determined.
func resolveGroupArtifact(cfg providers.ProviderConfig, appName string) (string, string, string, error) {
	// Gather explicit values
	group, _ := cfg["Group"].(string)
	if group == "" {
		if s, _ := cfg["group"].(string); s != "" {
			group = s
		}
	}
	artifact, _ := cfg["Artifact"].(string)
	if artifact == "" {
		if s, _ := cfg["artifact"].(string); s != "" {
			artifact = s
		}
	}
	packageName, _ := cfg["PackageName"].(string)
	if packageName == "" {
		if s, _ := cfg["package"].(string); s != "" {
			packageName = s
		}
	}

	// If package provided and group missing, derive group from package (all but last segment)
	if packageName != "" && group == "" {
		parts := strings.Split(packageName, ".")
		if len(parts) >= 2 {
			group = strings.Join(parts[:len(parts)-1], ".")
		}
	}

	// If artifact missing, derive from app name
	if artifact == "" && appName != "" {
		artifact = strings.ToLower(validator.SanitizeAppName(appName))
	}

	// If package missing but group+artifact present, derive package = group.artifact
	if packageName == "" && group != "" && artifact != "" {
		packageName = fmt.Sprintf("%s.%s", group, artifact)
	}

	// Final validation: require a package name
	if packageName == "" {
		return "", "", "", fmt.Errorf("package not provided; specify --package or provide --group and --artifact")
	}
	return group, artifact, packageName, nil
}

// checkJavaAvailable ensures Java is installed and on PATH and returns a
// user-friendly hint if not.
func checkJavaAvailable() error {
	if _, err := exec.LookPath("java"); err != nil {
		guide := `Java not found in PATH.
Grimoire requires Java (JDK or JRE) to build Java projects and run Gradle.

Quick install hints:
- Windows (Scoop): scoop install temurin
- Windows (Chocolatey): choco install temurin
- macOS (Homebrew): brew install temurin
- Linux (SDKMAN): curl -s "https://get.sdkman.io" | bash && source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install java
- Debian/Ubuntu (apt): sudo apt-get install default-jdk

After installing, ensure 'java' is on your PATH and re-run ` + "`grimoire new`" + `.
`
		return fmt.Errorf(guide)
	}
	return nil
}

// GenerateProject creates a Spring Boot or plain Java project based on cfg.
func GenerateProject(cfg providers.ProviderConfig) error {
	// Extract values from cfg
	appName, _ := cfg["AppName"].(string)
	if appName == "" {
		if s, _ := cfg["app-name"].(string); s != "" {
			appName = s
		}
	}
	if appName == "" {
		return fmt.Errorf("AppName is required in config")
	}

	// Resolve group, artifact and package according to rules:
	// - If --group not set, derive from --package (all segments except last)
	// - If --artifact not set, derive from app name
	// - If --package not set and group+artifact are set, derive package = group.artifact
	// - If none are set, return validation error asking for --package
	group, artifact, packageName, err := resolveGroupArtifact(cfg, appName)
	if err != nil {
		return err
	}

	templateKind, _ := cfg["Template"].(string)
	// Build system preference: prefer explicit BuildSystem, then provider flag 'build'
	buildSystem := "gradle"
	if b, ok := cfg["BuildSystem"].(string); ok && b != "" {
		buildSystem = strings.ToLower(b)
	} else if b2, ok := cfg["build"].(string); ok && b2 != "" {
		buildSystem = strings.ToLower(b2)
	}
	if templateKind == "" {
		if s, _ := cfg["template"].(string); s != "" {
			templateKind = s
		} else {
			templateKind = "springboot"
		}
	}

	outputDir, _ := cfg["OutputDir"].(string)
	if outputDir == "" {
		if s, _ := cfg["output-dir"].(string); s != "" {
			outputDir = s
		}
		if outputDir == "" {
			outputDir = filepath.Join(".", artifact)
		}
	}

	gitInit, _ := cfg["Git"].(bool)
	if !gitInit {
		if v, ok := cfg["git"].(bool); ok {
			gitInit = v
		}
	}

	// If Spring Boot template selected, ensure Java is available
	if strings.ToLower(templateKind) == "springboot" {
		if err := checkJavaAvailable(); err != nil {
			return err
		}
	}

	// Validation
	if err := validator.ValidateAppName(appName); err != nil {
		return err
	}
	if err := validator.ValidatePackageName(packageName); err != nil {
		return err
	}

	// Ensure output directory does not exist
	if _, err := os.Stat(outputDir); err == nil {
		return fmt.Errorf("output directory already exists: %s", outputDir)
	}

	baseName := validator.SanitizeAppName(appName)
	appClassName := baseName
	if strings.ToLower(templateKind) == "springboot" {
		appClassName = baseName + "Application"
	} else {
		appClassName = baseName + "Main"
	}

	data := map[string]any{
		"AppName":       appName,
		"AppNameLower":  strings.ToLower(appName),
		"Group":         group,
		"Artifact":      artifact,
		"PackageName":   packageName,
		"PackagePath":   validator.PackageToPath(packageName),
		"AppClassName":  appClassName,
		"Template":      templateKind,
		"TemplateLower": strings.ToLower(templateKind),
	}

	// write build files depending on build system
	if buildSystem == "maven" {
		if s, err := renderTemplate("pom.xml.tmpl", data); err == nil {
			if err := writeFile(filepath.Join(outputDir, "pom.xml"), s); err != nil {
				return err
			}
		}
	} else {
		if s, err := renderTemplate("settings_gradle.tmpl", data); err == nil {
			if err := writeFile(filepath.Join(outputDir, "settings.gradle"), s); err != nil {
				return err
			}
		}
		if s, err := renderTemplate("build_gradle.tmpl", data); err == nil {
			if err := writeFile(filepath.Join(outputDir, "build.gradle"), s); err != nil {
				return err
			}
		}
		if s, err := renderTemplate("gradle_properties.tmpl", data); err == nil {
			if err := writeFile(filepath.Join(outputDir, "gradle.properties"), s); err != nil {
				return err
			}
		}
	}
	if s, err := renderTemplate("gitignore.tmpl", data); err == nil {
		if err := writeFile(filepath.Join(outputDir, ".gitignore"), s); err != nil {
			return err
		}
	}
	if s, err := renderTemplate("README.tmpl", data); err == nil {
		if err := writeFile(filepath.Join(outputDir, "README.md"), s); err != nil {
			return err
		}
	}

	// Source
	pkgPath := validator.PackageToPath(packageName)
	var appSrc string
	var rerr error
	if strings.ToLower(templateKind) == "springboot" {
		appSrc, rerr = renderTemplate("application_springboot.java.tmpl", data)
	} else {
		appSrc, rerr = renderTemplate("application_plain.java.tmpl", data)
	}
	if rerr != nil {
		return rerr
	}
	if err := writeFile(filepath.Join(outputDir, "src", "main", "java", pkgPath, appClassName+".java"), appSrc); err != nil {
		return err
	}

	// resources: create application.properties for all templates
	if s, err := renderTemplate("application.properties.tmpl", data); err == nil {
		if err := writeFile(filepath.Join(outputDir, "src", "main", "resources", "application.properties"), s); err != nil {
			return err
		}
	}

	if gitInit {
		if err := initGit(outputDir); err != nil {
			// non-fatal
		}
	}

	return nil
}
