package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name       string
		tmplName   string
		tmplString string
		data       any
		want       string
		wantErr    bool
	}{
		{
			name:       "simple template",
			tmplName:   "test",
			tmplString: "Hello {{.Name}}",
			data: map[string]string{
				"Name": "World",
			},
			want:    "Hello World",
			wantErr: false,
		},
		{
			name:       "invalid template",
			tmplName:   "test",
			tmplString: "Hello {{.Name",
			data:       map[string]string{},
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTemplate(tt.tmplName, tt.tmplString, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RenderTemplate got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		content   string
		wantExist bool
		wantErr   bool
	}{
		{
			name:      "write file success",
			path:      filepath.Join(t.TempDir(), "test.txt"),
			content:   "Hello World",
			wantExist: true,
			wantErr:   false,
		},
		{
			name:      "create nested directories",
			path:      filepath.Join(t.TempDir(), "dir1", "dir2", "test.txt"),
			content:   "Nested",
			wantExist: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteFile(tt.path, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantExist {
				data, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("WriteFile failed to create file: %v", err)
					return
				}
				if string(data) != tt.content {
					t.Errorf("WriteFile content mismatch: got %q, want %q", string(data), tt.content)
				}
			}
		})
	}
}

func TestInitGit(t *testing.T) {
	tmpDir := t.TempDir()

	err := InitGit(tmpDir)
	if err != nil {
		t.Errorf("InitGit failed: %v", err)
		return
	}

	// Check if .git directory was created
	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("InitGit did not create .git directory")
	}
}
