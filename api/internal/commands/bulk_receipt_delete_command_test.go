package commands

import (
	"testing"
)

// sequentialIds returns the ids 1..count, so a list of any length has no duplicates.
func sequentialIds(count int) []uint {
	ids := make([]uint, count)
	for i := range ids {
		ids[i] = uint(i + 1)
	}
	return ids
}

func TestBulkReceiptDeleteCommand_Validate(t *testing.T) {
	cases := map[string]struct {
		ids       []uint
		wantError bool
	}{
		"nil":            {ids: nil, wantError: true},
		"empty":          {ids: []uint{}, wantError: true},
		"duplicate":      {ids: []uint{1, 2, 1}, wantError: true},
		"valid":          {ids: []uint{1, 2, 3}, wantError: false},
		"at the limit":   {ids: sequentialIds(MaxBulkReceiptDeleteIds), wantError: false},
		"over the limit": {ids: sequentialIds(MaxBulkReceiptDeleteIds + 1), wantError: true},
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
