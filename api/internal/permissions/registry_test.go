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

// TestAppDataEffectivePermissionsAreUntypedStrings pins the AppData transport
// decision so a future "consistency cleanup" cannot re-arm it.
//
// The Permission schema is a CLOSED enum in generated clients. The dart-dio
// EnumClass's _$valueOf ends in `default: throw ArgumentError(name)`, so an
// unknown wire value fails the WHOLE AppData deserialization -- and since
// permissions hydrate on login, that hard-fails login on every already-released
// mobile build the moment a permission is added here. It shipped twice:
// 2026-07-24 (group.members.create) and 2026-08-06 (group.members.grants.update,
// PR #661).
//
// This rejects BOTH ways of closing the set -- a $ref to the Permission schema
// and an inline `enum:` on the items -- because the generator turns either one
// into the same throwing EnumClass.
//
// So a user's EFFECTIVE permissions ride as plain strings: that payload is
// server-resolved data the client only pattern-matches, and a granted entry may
// even be a wildcard (matcher.go), which is not an enum member at all. The
// CATALOG of which permissions exist keeps the enum -- Role.permissions,
// UpsertRoleCommand.permissions, PermissionDescriptor.key, and the `permission`
// query param, all still covered by TestSwaggerEnumMatchesRegistry above.
func TestAppDataEffectivePermissionsAreUntypedStrings(t *testing.T) {
	swaggerPath := findSwagger(t)

	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read swagger.yml: %v", err)
	}

	// Enum is checked alongside Ref: an INLINE closed set --
	// `items: {type: string, enum: [...]}` -- generates exactly the same
	// throwing EnumClass as a $ref to the Permission schema, so checking only
	// for $ref would let the outage back in through the other door.
	type itemSchema struct {
		Type string   `yaml:"type"`
		Ref  string   `yaml:"$ref"`
		Enum []string `yaml:"enum"`
	}

	var doc struct {
		Components struct {
			Schemas map[string]struct {
				Properties map[string]struct {
					Items                itemSchema `yaml:"items"`
					AdditionalProperties struct {
						Items itemSchema `yaml:"items"`
					} `yaml:"additionalProperties"`
				} `yaml:"properties"`
			} `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse swagger.yml: %v", err)
	}

	appData, ok := doc.Components.Schemas["AppData"]
	if !ok {
		t.Fatal("swagger.yml is missing the AppData schema")
	}

	// Both fields must be open string sets: no $ref to a schema (which would be
	// an enum) and no inline enum either.
	assertOpenStringItems := func(field string, items itemSchema) {
		t.Helper()
		if items.Ref != "" {
			t.Errorf("AppData.%s items use $ref %q; effective permissions must be plain strings", field, items.Ref)
		}
		if len(items.Enum) > 0 {
			t.Errorf("AppData.%s items declare an inline enum %v; effective permissions must be an open string set", field, items.Enum)
		}
		if items.Type != "string" {
			t.Errorf("AppData.%s items type = %q, want %q", field, items.Type, "string")
		}
	}

	appPermissions, ok := appData.Properties["appPermissions"]
	if !ok {
		t.Fatal("AppData is missing appPermissions")
	}
	assertOpenStringItems("appPermissions", appPermissions.Items)

	groupPermissions, ok := appData.Properties["groupPermissions"]
	if !ok {
		t.Fatal("AppData is missing groupPermissions")
	}
	assertOpenStringItems("groupPermissions", groupPermissions.AdditionalProperties.Items)
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
