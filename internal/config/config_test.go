package config

import "testing"

func TestNormalizeAndValidateDebounceMinimum(t *testing.T) {
	c := &Config{Debounce: 1}
	if err := c.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate returned error: %v", err)
	}

	if c.Debounce != 50 {
		t.Fatalf("expected debounce to be clamped to 50, got %d", c.Debounce)
	}
}

func TestNormalizeAndValidateExtensions(t *testing.T) {
	c := &Config{WatchExts: []string{"go", ".MOD", "  .sum "}, Debounce: 100}
	if err := c.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate returned error: %v", err)
	}

	want := []string{".go", ".mod", ".sum"}
	for i := range want {
		if c.WatchExts[i] != want[i] {
			t.Fatalf("watch ext mismatch at %d: want %s got %s", i, want[i], c.WatchExts[i])
		}
	}
}

func TestServicesForProfile(t *testing.T) {
	c := &Config{
		Project: ProjectConfig{DefaultProfile: "dev"},
		Services: ServicesMap{
			"api":    {Type: "go", Package: "./cmd/api"},
			"worker": {Type: "go", Package: "./cmd/worker"},
		},
		Profiles: ProfilesMap{
			"dev": {Services: []string{"api", "worker"}},
		},
		Debounce: 100,
	}

	if err := c.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate returned error: %v", err)
	}

	got, err := c.ServicesForProfile("")
	if err != nil {
		t.Fatalf("ServicesForProfile returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 services, got %d", len(got))
	}
}

func TestNormalizeAndValidateServiceHealthcheckValidation(t *testing.T) {
	c := &Config{
		Services: ServicesMap{
			"api": {
				Type:    "go",
				Package: "./cmd/api",
				Healthcheck: HealthcheckConfig{
					Type: "http",
				},
			},
		},
		Debounce: 100,
	}

	if err := c.NormalizeAndValidate(); err == nil {
		t.Fatal("expected validation error for missing healthcheck.url")
	}
}

func TestNormalizeAndValidateLogFormat(t *testing.T) {
	c := &Config{Project: ProjectConfig{LogFormat: "json"}, Debounce: 100}
	if err := c.NormalizeAndValidate(); err != nil {
		t.Fatalf("expected valid json log format, got error: %v", err)
	}

	c2 := &Config{Project: ProjectConfig{LogFormat: "xml"}, Debounce: 100}
	if err := c2.NormalizeAndValidate(); err == nil {
		t.Fatal("expected error for invalid log format")
	}
}
