package cmd

import "testing"

func TestValidateUI(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"", false},
		{"xml", false},
		{"XML", false},
		{"compose", false},
		{"none", false},
		{"invalid", true},
	}
	for _, c := range cases {
		if err := validateUI(c.input); (err != nil) != c.wantErr {
			t.Fatalf("validateUI(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
		}
	}
}

func TestValidateLang(t *testing.T) {
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
		if err := validateLang(c.input); (err != nil) != c.wantErr {
			t.Fatalf("validateLang(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
		}
	}
}

func TestValidateDI(t *testing.T) {
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
		if err := validateDI(c.input); (err != nil) != c.wantErr {
			t.Fatalf("validateDI(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
		}
	}
}
