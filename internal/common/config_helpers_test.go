package common

import (
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
)

func TestGetString(t *testing.T) {
	cfg := providers.ProviderConfig{
		"Key1": "value1",
		"key2": "value2",
	}

	tests := []struct {
		name    string
		keys    []string
		want    string
		wantNil bool
	}{
		{
			name: "single key found",
			keys: []string{"Key1"},
			want: "value1",
		},
		{
			name: "fallback to second key",
			keys: []string{"Missing", "key2"},
			want: "value2",
		},
		{
			name: "key not found",
			keys: []string{"Missing"},
			want: "",
		},
		{
			name: "empty string ignored",
			keys: []string{"EmptyKey", "key2"},
			want: "value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetString(cfg, tt.keys...)
			if got != tt.want {
				t.Errorf("GetString got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetStringDefault(t *testing.T) {
	cfg := providers.ProviderConfig{
		"Key1": "value1",
	}

	tests := []struct {
		name       string
		defaultVal string
		keys       []string
		want       string
	}{
		{
			name:       "key found, ignores default",
			defaultVal: "default",
			keys:       []string{"Key1"},
			want:       "value1",
		},
		{
			name:       "key not found, returns default",
			defaultVal: "default",
			keys:       []string{"Missing"},
			want:       "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStringDefault(cfg, tt.defaultVal, tt.keys...)
			if got != tt.want {
				t.Errorf("GetStringDefault got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	cfg := providers.ProviderConfig{
		"IntKey":    42,
		"StringInt": "123",
		"ZeroKey":   0,
	}

	tests := []struct {
		name string
		keys []string
		want int
	}{
		{
			name: "int value found",
			keys: []string{"IntKey"},
			want: 42,
		},
		{
			name: "string int converted",
			keys: []string{"StringInt"},
			want: 123,
		},
		{
			name: "zero value ignored",
			keys: []string{"ZeroKey"},
			want: 0,
		},
		{
			name: "key not found",
			keys: []string{"Missing"},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInt(cfg, tt.keys...)
			if got != tt.want {
				t.Errorf("GetInt got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	cfg := providers.ProviderConfig{
		"BoolTrue": true,
		"BoolFalse": false,
		"StringTrue": "true",
		"StringOne": "1",
		"StringYes": "yes",
		"StringNo": "no",
	}

	tests := []struct {
		name string
		keys []string
		want bool
	}{
		{
			name: "bool true found",
			keys: []string{"BoolTrue"},
			want: true,
		},
		{
			name: "bool false found",
			keys: []string{"BoolFalse"},
			want: false,
		},
		{
			name: "string true",
			keys: []string{"StringTrue"},
			want: true,
		},
		{
			name: "string 1",
			keys: []string{"StringOne"},
			want: true,
		},
		{
			name: "string yes",
			keys: []string{"StringYes"},
			want: true,
		},
		{
			name: "key not found",
			keys: []string{"Missing"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBool(cfg, tt.keys...)
			if got != tt.want {
				t.Errorf("GetBool got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolDefault(t *testing.T) {
	cfg := providers.ProviderConfig{
		"BoolTrue": true,
	}

	tests := []struct {
		name       string
		defaultVal bool
		keys       []string
		want       bool
	}{
		{
			name:       "bool found, ignores default",
			defaultVal: false,
			keys:       []string{"BoolTrue"},
			want:       true,
		},
		{
			name:       "key not found, returns default true",
			defaultVal: true,
			keys:       []string{"Missing"},
			want:       true,
		},
		{
			name:       "key not found, returns default false",
			defaultVal: false,
			keys:       []string{"Missing"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBoolDefault(cfg, tt.defaultVal, tt.keys...)
			if got != tt.want {
				t.Errorf("GetBoolDefault got %v, want %v", got, tt.want)
			}
		})
	}
}
