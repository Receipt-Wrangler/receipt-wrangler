package commands

import (
	"testing"
)

func TestBulkReceiptDeleteCommand_Validate(t *testing.T) {
	cases := map[string]struct {
		ids       []uint
		wantError bool
	}{
		"nil":       {ids: nil, wantError: true},
		"empty":     {ids: []uint{}, wantError: true},
		"duplicate": {ids: []uint{1, 2, 1}, wantError: true},
		"valid":     {ids: []uint{1, 2, 3}, wantError: false},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			command := BulkReceiptDeleteCommand{ReceiptIds: testCase.ids}
			errs := command.Validate()

			_, hasError := errs.Errors["receiptIds"]
			if hasError != testCase.wantError {
				t.Errorf("receiptIds error = %v, want %v (errors: %v)", hasError, testCase.wantError, errs.Errors)
			}
			if len(errs.Errors) > 1 {
				t.Errorf("unexpected extra errors: %v", errs.Errors)
			}
		})
	}
}
