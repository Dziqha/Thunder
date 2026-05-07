package cli

import "testing"

func TestResolveVersionUsesInjectedVersion(t *testing.T) {
	orig := version
	version = "v9.9.9"
	t.Cleanup(func() { version = orig })

	got := resolveVersion()
	if got != "v9.9.9" {
		t.Fatalf("expected injected version, got %q", got)
	}
}
