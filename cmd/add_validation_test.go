package cmd

import (
	"testing"

	"github.com/backsofangels/grimoire/internal/validator"
)

func TestValidateUIFromValidator(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"", false},
		{"xml", false},
		{"XML", false},
		{"compose", false},
		{"none", false},  // "none" is now valid (disables UI)
		{"invalid", true},
	}
	for _, c := range cases {
		if err := validator.ValidateUI(c.input); (err != nil) != c.wantErr {
			t.Errorf("ValidateUI(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
		}
	}
}

func TestValidateLangFromValidator(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"", false},
		{"kotlin", false},
		{"KOTLIN", false},
		{"java", false},
		{"py", true},
	}
	for _, c := range cases {
		if err := validator.ValidateLanguage(c.input); (err != nil) != c.wantErr {
			t.Errorf("ValidateLanguage(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
		}
	}
}

func TestValidateDIFromValidator(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"", false},
		{"none", false},
		{"hilt", false},
		{"koin", false},
		{"bad", true},
	}
	for _, c := range cases {
		if err := validator.ValidateDI(c.input); (err != nil) != c.wantErr {
			t.Errorf("ValidateDI(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
		}
	}
}
