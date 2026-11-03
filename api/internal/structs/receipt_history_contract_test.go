package structs

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type swaggerProperty struct {
	Type  string   `yaml:"type"`
	Enum  []string `yaml:"enum"`
	Ref   string   `yaml:"$ref"`
	Items *struct {
		Ref string `yaml:"$ref"`
	} `yaml:"items"`
}

type swaggerSchema struct {
	Required   []string                   `yaml:"required"`
	Properties map[string]swaggerProperty `yaml:"properties"`
}

func readSwaggerSchemas(t *testing.T) map[string]swaggerSchema {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(thisFile), "..", "..", "swagger.yml"))
	if err != nil {
		t.Fatalf("read swagger.yml: %v", err)
	}

	var doc struct {
		Components struct {
			Schemas map[string]swaggerSchema `yaml:"schemas"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse swagger.yml: %v", err)
	}
	return doc.Components.Schemas
}

// jsonKeys lists the JSON keys a struct serializes, and the ones it always
// serializes (no omitempty).
func jsonKeys(value interface{}) (keys []string, alwaysPresent []string) {
	structType := reflect.TypeOf(value)
	for i := 0; i < structType.NumField(); i++ {
		tag := structType.Field(i).Tag.Get("json")
		parts := strings.Split(tag, ",")
		if parts[0] == "" || parts[0] == "-" {
			continue
		}
		keys = append(keys, parts[0])
		if !slices.Contains(parts[1:], "omitempty") {
			alwaysPresent = append(alwaysPresent, parts[0])
		}
	}
	sort.Strings(keys)
	sort.Strings(alwaysPresent)
	return keys, alwaysPresent
}

// The receipt history is read by generated clients, including the Dart one,
// whose closed enums fail a whole payload on a value they do not know, the
// cause of two production login outages. Every code on it is therefore a plain
// string, so the server can add an action, a field or a value type without
// breaking a released build. The Go structs must also serialize exactly the
// keys swagger declares, and always send the ones it requires.
func TestReceiptHistorySchemasMatchTheContract(t *testing.T) {
	schemas := readSwaggerSchemas(t)

	cases := map[string]interface{}{
		"ReceiptHistory":          ReceiptHistory{},
		"ReceiptHistoryEntry":     ReceiptHistoryEntry{},
		"ReceiptHistoryActor":     ReceiptHistoryActor{},
		"ReceiptHistoryChange":    ReceiptHistoryChange{},
		"ReceiptHistoryValue":     ReceiptHistoryValue{},
		"ReceiptHistorySnapshots": ReceiptHistorySnapshots{},
	}

	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			schema, ok := schemas[name]
			if !ok {
				t.Fatalf("swagger.yml is missing the %s schema", name)
			}

			swaggerKeys := make([]string, 0, len(schema.Properties))
			for key, property := range schema.Properties {
				swaggerKeys = append(swaggerKeys, key)
				if len(property.Enum) > 0 {
					t.Errorf("%s.%s must be a plain string, not an enum %v", name, key, property.Enum)
				}
				if property.Ref != "" && !strings.HasPrefix(property.Ref, "#/components/schemas/ReceiptHistory") {
					t.Errorf("%s.%s refers to %s, outside the ReceiptHistory schemas", name, key, property.Ref)
				}
			}
			sort.Strings(swaggerKeys)

			goKeys, alwaysPresent := jsonKeys(value)
			if !slices.Equal(goKeys, swaggerKeys) {
				t.Errorf("Go keys %v, swagger properties %v", goKeys, swaggerKeys)
			}
			for _, required := range schema.Required {
				if !slices.Contains(alwaysPresent, required) {
					t.Errorf("%s.%s is required on the contract but may be omitted by Go", name, required)
				}
			}
		})
	}
}

// NewReceiptHistory must serialize its entries as [], never null: the
// generated Dart deserializer has no null guard.
func TestNewReceiptHistorySerializesEmptyEntries(t *testing.T) {
	history := NewReceiptHistory()
	if history.Entries == nil || len(history.Entries) != 0 {
		t.Errorf("entries = %#v, want an empty, non-nil slice", history.Entries)
	}
}

// The email origin rides on Receipt, which every client deserializes: it must
// stay a pair of optional, plain strings.
func TestReceiptEmailOriginIsOptionalPlainStrings(t *testing.T) {
	schema, ok := readSwaggerSchemas(t)["Receipt"]
	if !ok {
		t.Fatal("swagger.yml is missing the Receipt schema")
	}

	for _, key := range []string{"emailInbox", "emailSender"} {
		property, ok := schema.Properties[key]
		if !ok {
			t.Errorf("Receipt.%s is not declared", key)
			continue
		}
		if property.Type != "string" || len(property.Enum) > 0 || property.Ref != "" {
			t.Errorf("Receipt.%s must be a plain string, got type=%q enum=%v $ref=%q", key, property.Type, property.Enum, property.Ref)
		}
		if slices.Contains(schema.Required, key) {
			t.Errorf("Receipt.%s must stay optional", key)
		}
	}
}
