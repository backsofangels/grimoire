package springboot

import (
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
)

func TestProviderRegistrationAndFlags(t *testing.T) {
	p, err := providers.Get("springboot")
	if err != nil {
		t.Fatalf("springboot provider not registered: %v", err)
	}
	if p.Name() != "springboot" {
		t.Fatalf("expected provider name 'springboot', got %s", p.Name())
	}

	flags := p.Flags()
	found := map[string]bool{}
	for _, f := range flags {
		found[f.Name] = true
	}
	for _, want := range []string{"group", "artifact", "template"} {
		if !found[want] {
			t.Fatalf("expected flag %s to be present", want)
		}
	}

	checks := p.DoctorChecks()
	if len(checks) < 3 {
		t.Fatalf("expected at least 3 doctor checks, got %d", len(checks))
	}
	names := map[string]bool{}
	for _, c := range checks {
		names[c.Name] = true
	}
	if !names["JDK available (javac)"] {
		t.Fatalf("missing check 'JDK available (javac)'")
	}
	if !names["JAVA_HOME points to JDK"] {
		t.Fatalf("missing check 'JAVA_HOME points to JDK'")
	}
	if !names["JDK version >= 11"] {
		t.Fatalf("missing check 'JDK version >= 11'")
	}
}
