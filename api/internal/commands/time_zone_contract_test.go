package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gopkg.in/yaml.v3"
)

// timeZone rides on AppData, which every mobile build deserializes at login.
// The generated Dart client renders an enum as a closed EnumClass that throws
// on a value it does not know, failing the WHOLE payload — the shape of both
// documented login outages. IANA names change with every tzdata release, so the
// field must stay a plain, optional string on every schema that carries it.
func TestTimeZoneIsAnOpenOptionalStringOnTheContract(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "swagger.yml"))
	if err != nil {
		t.Fatalf("read swagger.yml: %v", err)
	}

	type property struct {
		Type string   `yaml:"type"`
		Enum []string `yaml:"enum"`
		Ref  string   `yaml:"$ref"`
	}
	type schemaPart struct {
		Required   []string            `yaml:"required"`
		Properties map[string]property `yaml:"properties"`
	}
	type schema struct {
		schemaPart `yaml:",inline"`
		AllOf      []schemaPart `yaml:"allOf"`
	}
	var doc struct {
		Components struct {
			Schemas map[string]schema `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse swagger.yml: %v", err)
	}

	for _, name := range []string{"SystemSettings", "UpsertSystemSettingsCommand", "AppData"} {
		t.Run(name, func(t *testing.T) {
			found, ok := doc.Components.Schemas[name]
			if !ok {
				t.Fatalf("swagger.yml is missing the %s schema", name)
			}

			parts := append([]schemaPart{found.schemaPart}, found.AllOf...)
			declared := false
			for _, part := range parts {
				for _, required := range part.Required {
					if required == "timeZone" {
						t.Error("timeZone must stay optional: older servers and clients omit it")
					}
				}
				timeZone, ok := part.Properties["timeZone"]
				if !ok {
					continue
				}
				declared = true
				if timeZone.Type != "string" || len(timeZone.Enum) > 0 || timeZone.Ref != "" {
					t.Errorf("timeZone must be a plain string, got type=%q enum=%v $ref=%q",
						timeZone.Type, timeZone.Enum, timeZone.Ref)
				}
			}
			if !declared {
				t.Errorf("%s does not declare timeZone", name)
			}
		})
	}
}
