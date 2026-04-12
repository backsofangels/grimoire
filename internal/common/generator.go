package common

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// RenderTemplate renders a Go template string with the given data.
func RenderTemplate(tmplName, tmplContent string, data any) (string, error) {
	t, err := template.New(tmplName).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", tmplName, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %s: %w", tmplName, err)
	}
	return buf.String(), nil
}

// WriteFile writes content to a file, creating directories as needed.
func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// InitGit initializes a git repository in the given directory.
func InitGit(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git init: %w", err)
	}
	return nil
}

// RenderAndWriteTemplate renders a template from an embedded filesystem and writes to file.
func RenderAndWriteTemplate(embedFS embed.FS, outputDir, templateName, outputPath string, data any) error {
	content, err := fs.ReadFile(embedFS, filepath.Join("templates", templateName))
	if err != nil {
		return fmt.Errorf("read template %s: %w", templateName, err)
	}

	rendered, err := RenderTemplate(templateName, string(content), data)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(outputDir, outputPath)
	if err := WriteFile(fullPath, rendered); err != nil {
		return err
	}

	return nil
}
