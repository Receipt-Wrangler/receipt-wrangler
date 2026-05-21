package permissions

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRegistryInvariants(t *testing.T) {
	descriptors := All()

	if len(descriptors) == 0 {
		t.Fatal("registry is empty")
	}

	seen := map[string]bool{}
	for _, d := range descriptors {
		if d.Key == "" {
			t.Errorf("descriptor has empty key: %+v", d)
		}
		if d.Label == "" {
			t.Errorf("descriptor %q has empty label", d.Key)
		}
		if d.Description == "" {
			t.Errorf("descriptor %q has empty description", d.Key)
		}
		if d.Category == "" {
			t.Errorf("descriptor %q has empty category", d.Key)
		}
		if d.Scope != ScopeApp && d.Scope != ScopeGroup {
			t.Errorf("descriptor %q has invalid scope %q", d.Key, d.Scope)
		}
		if seen[d.Key] {
			t.Errorf("duplicate key %q", d.Key)
		}
		seen[d.Key] = true
	}
}

func TestGetAndExistsRoundTrip(t *testing.T) {
	for _, d := range All() {
		got, ok := Get(d.Key)
		if !ok {
			t.Errorf("Get(%q) returned ok=false", d.Key)
			continue
		}
		if got != d {
			t.Errorf("Get(%q) returned %+v, want %+v", d.Key, got, d)
		}
		if !Exists(d.Key) {
			t.Errorf("Exists(%q) returned false", d.Key)
		}
	}

	if _, ok := Get("not.a.real.permission"); ok {
		t.Error("Get returned ok=true for unknown key")
	}
	if Exists("not.a.real.permission") {
		t.Error("Exists returned true for unknown key")
	}
}

func TestSwaggerEnumMatchesRegistry(t *testing.T) {
	swaggerPath := findSwagger(t)

	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read swagger.yml: %v", err)
	}

	var doc struct {
		Components struct {
			Schemas map[string]struct {
				Type string   `yaml:"type"`
				Enum []string `yaml:"enum"`
			} `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse swagger.yml: %v", err)
	}

	permSchema, ok := doc.Components.Schemas["Permission"]
	if !ok {
		t.Fatal("swagger.yml is missing the Permission schema")
	}
	if permSchema.Type != "string" {
		t.Errorf("Permission schema type = %q, want %q", permSchema.Type, "string")
	}

	swaggerKeys := append([]string{}, permSchema.Enum...)
	sort.Strings(swaggerKeys)

	registryKeys := make([]string, 0, len(registry))
	for _, d := range registry {
		registryKeys = append(registryKeys, d.Key)
	}
	sort.Strings(registryKeys)

	if len(swaggerKeys) != len(registryKeys) {
		t.Fatalf("swagger Permission enum has %d values, registry has %d", len(swaggerKeys), len(registryKeys))
	}
	for i, key := range registryKeys {
		if swaggerKeys[i] != key {
			t.Errorf("swagger/registry mismatch at index %d: swagger=%q registry=%q", i, swaggerKeys[i], key)
		}
	}
}

func findSwagger(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// /app/api/internal/permissions/registry_test.go -> /app/api/swagger.yml
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "swagger.yml")
}
