package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

// A stored report template's Configuration is json.Marshal of the command, then read
// back into the desktop builder, which reads the filter's lowercase `value` / `tags`
// keys (matching swagger). If the struct marshals a capital `Value` / `Tags`, the
// operation hydrates but the value silently drops — the filter appears empty in the
// builder. This pins the json tags so that regression can't return.
func TestReceiptPagedRequestFilter_MarshalsLowercaseKeys(t *testing.T) {
	filter := ReceiptPagedRequestFilter{
		Categories: PagedRequestField{Operation: CONTAINS, Value: []int{1, 2}},
		Tags:       PagedRequestField{Operation: CONTAINS, Value: []int{7}},
	}

	bytes, err := json.Marshal(filter)
	if err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}
	got := string(bytes)

	for _, want := range []string{`"value":[1,2]`, `"tags":`, `"categories":`} {
		if !strings.Contains(got, want) {
			utils.PrintTestError(t, got, "to contain "+want)
		}
	}
	for _, unwanted := range []string{`"Value"`, `"Tags"`} {
		if strings.Contains(got, unwanted) {
			utils.PrintTestError(t, got, "to NOT contain "+unwanted)
		}
	}
}

// Loading a filter and re-marshaling it (the template store path) must keep the
// category value, and a legacy blob that carried a capital `Value` must still
// deserialize (Go matches json keys case-insensitively) so old templates still parse.
func TestReceiptPagedRequestFilter_ValueRoundTrips(t *testing.T) {
	tests := map[string]string{
		"lowercase value":      `{"filter":{"categories":{"operation":"CONTAINS","value":[3,4]}}}`,
		"legacy capital Value": `{"filter":{"categories":{"operation":"CONTAINS","Value":[3,4]}}}`,
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			command := ReceiptPagedRequestCommand{}
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			w := httptest.NewRecorder()
			if err := command.LoadDataFromRequest(w, r); err != nil {
				utils.PrintTestError(t, err, nil)
				return
			}

			bytes, err := json.Marshal(command.Filter)
			if err != nil {
				utils.PrintTestError(t, err, nil)
				return
			}
			if !strings.Contains(string(bytes), `"categories":{"operation":"CONTAINS","value":[3,4]}`) {
				utils.PrintTestError(t, string(bytes), `categories with lowercase value [3,4]`)
			}
		})
	}
}

func TestReceiptPagedRequestCommand_LoadDataFromRequest_NormalizesGroup(t *testing.T) {
	tests := map[string]string{
		"empty string group value": `{"filter":{"group":{"operation":"CONTAINS","value":""}}}`,
		"null group value":         `{"filter":{"group":{"operation":"CONTAINS","value":null}}}`,
		"missing group field":      `{"filter":{}}`,
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			command := ReceiptPagedRequestCommand{}
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			w := httptest.NewRecorder()

			if err := command.LoadDataFromRequest(w, r); err != nil {
				utils.PrintTestError(t, err, nil)
				return
			}

			value, ok := command.Filter.Group.Value.([]interface{})
			if !ok {
				utils.PrintTestError(t, command.Filter.Group.Value, "[]interface{}")
				return
			}
			if len(value) != 0 {
				utils.PrintTestError(t, len(value), 0)
			}
		})
	}
}

func TestPagedRequestCommand_Validate_ValidInputs(t *testing.T) {
	tests := map[string]struct {
		command PagedRequestCommand
	}{
		"valid ascending": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      10,
				SortDirection: ASCENDING,
			},
		},
		"valid descending": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      10,
				SortDirection: DESCENDING,
			},
		},
		"valid default sort direction": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      10,
				SortDirection: DEFAULT,
			},
		},
		"valid pageSize -1 (no limit)": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      -1,
				SortDirection: ASCENDING,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate()

			if len(vErr.Errors) > 0 {
				utils.PrintTestError(t, len(vErr.Errors), 0)
			}
		})
	}
}

func TestPagedRequestCommand_Validate_InvalidInputs(t *testing.T) {
	tests := map[string]struct {
		command       PagedRequestCommand
		expectedError string
	}{
		"page less than 1": {
			command: PagedRequestCommand{
				Page:          0,
				PageSize:      10,
				SortDirection: ASCENDING,
			},
			expectedError: "page",
		},
		"negative page": {
			command: PagedRequestCommand{
				Page:          -1,
				PageSize:      10,
				SortDirection: ASCENDING,
			},
			expectedError: "page",
		},
		"pageSize zero": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      0,
				SortDirection: ASCENDING,
			},
			expectedError: "pageSize",
		},
		"negative pageSize (not -1)": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      -2,
				SortDirection: ASCENDING,
			},
			expectedError: "pageSize",
		},
		"invalid sort direction": {
			command: PagedRequestCommand{
				Page:          1,
				PageSize:      10,
				SortDirection: "invalid",
			},
			expectedError: "sortDirection",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate()

			if len(vErr.Errors) == 0 {
				utils.PrintTestError(t, len(vErr.Errors), "greater than 0")
			}

			if _, exists := vErr.Errors[test.expectedError]; !exists {
				utils.PrintTestError(t, "error should exist for field", test.expectedError)
			}
		})
	}
}

func TestPagedRequestCommand_Validate_MultipleErrors(t *testing.T) {
	command := PagedRequestCommand{
		Page:          0,
		PageSize:      0,
		SortDirection: "invalid",
	}

	vErr := command.Validate()

	if len(vErr.Errors) != 3 {
		utils.PrintTestError(t, len(vErr.Errors), 3)
	}

	if _, exists := vErr.Errors["page"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "page")
	}

	if _, exists := vErr.Errors["pageSize"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "pageSize")
	}

	if _, exists := vErr.Errors["sortDirection"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "sortDirection")
	}
}
