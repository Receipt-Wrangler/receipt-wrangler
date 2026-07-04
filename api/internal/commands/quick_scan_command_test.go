package commands

import (
	"mime/multipart"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
)

type mockFile struct {
	*strings.Reader
}

func (m mockFile) Close() error {
	return nil
}

func newMockFile() multipart.File {
	return mockFile{Reader: strings.NewReader("mock file content")}
}

func TestQuickScanCommand_Validate_ValidInputs(t *testing.T) {
	tests := map[string]struct {
		command QuickScanCommand
	}{
		"valid with single file": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile()},
				PaidByUserIds: []uint{1},
				GroupIds:      []uint{1},
				Statuses:      []models.ReceiptStatus{"OPEN"},
			},
		},
		"valid with multiple files": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1, 2},
				GroupIds:      []uint{1, 1},
				Statuses:      []models.ReceiptStatus{"OPEN", "OPEN"},
			},
		},
		"valid with category and tag selections": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1, 2},
				GroupIds:      []uint{1, 1},
				Statuses:      []models.ReceiptStatus{"OPEN", "OPEN"},
				CategoryIds:   [][]uint{{1, 2}, {}},
				TagIds:        [][]uint{{}, {5}},
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

func TestQuickScanCommand_Validate_InvalidInputs(t *testing.T) {
	tests := map[string]struct {
		command        QuickScanCommand
		expectedErrors []string
	}{
		"no files": {
			command:        QuickScanCommand{},
			expectedErrors: []string{"files", "paidByUserId", "groupId", "status"},
		},
		"mismatched paid by user ids": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1},
				GroupIds:      []uint{1, 1},
				Statuses:      []models.ReceiptStatus{"OPEN", "OPEN"},
			},
			expectedErrors: []string{"paidByUserId"},
		},
		"mismatched group ids": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1, 2},
				GroupIds:      []uint{1},
				Statuses:      []models.ReceiptStatus{"OPEN", "OPEN"},
			},
			expectedErrors: []string{"groupIds"},
		},
		"mismatched statuses": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1, 2},
				GroupIds:      []uint{1, 1},
				Statuses:      []models.ReceiptStatus{"OPEN"},
			},
			expectedErrors: []string{"statuses"},
		},
		"all arrays empty with no files": {
			command: QuickScanCommand{
				Files:         []multipart.File{},
				PaidByUserIds: []uint{},
				GroupIds:      []uint{},
				Statuses:      []models.ReceiptStatus{},
			},
			expectedErrors: []string{"files", "paidByUserId", "groupId", "status"},
		},
		"mismatched category ids": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1, 2},
				GroupIds:      []uint{1, 1},
				Statuses:      []models.ReceiptStatus{"OPEN", "OPEN"},
				CategoryIds:   [][]uint{{1}},
			},
			expectedErrors: []string{"categoryIds"},
		},
		"mismatched tag ids": {
			command: QuickScanCommand{
				Files:         []multipart.File{newMockFile(), newMockFile()},
				PaidByUserIds: []uint{1, 2},
				GroupIds:      []uint{1, 1},
				Statuses:      []models.ReceiptStatus{"OPEN", "OPEN"},
				TagIds:        [][]uint{{1}},
			},
			expectedErrors: []string{"tagIds"},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate()

			if len(vErr.Errors) == 0 {
				utils.PrintTestError(t, len(vErr.Errors), "greater than 0")
			}

			for _, expectedError := range test.expectedErrors {
				if _, exists := vErr.Errors[expectedError]; !exists {
					utils.PrintTestError(t, "error should exist for field", expectedError)
				}
			}
		})
	}
}

func TestQuickScanCommand_CategoryAndTagIdsForFile(t *testing.T) {
	command := QuickScanCommand{
		CategoryIds: [][]uint{{1, 2}},
		TagIds:      [][]uint{{9}},
	}

	if got := command.CategoryIdsForFile(0); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		utils.PrintTestError(t, got, []uint{1, 2})
	}

	// Out-of-range index (client omitted selections for this file) yields an empty slice.
	if got := command.CategoryIdsForFile(5); len(got) != 0 {
		utils.PrintTestError(t, got, []uint{})
	}

	if got := command.TagIdsForFile(0); len(got) != 1 || got[0] != 9 {
		utils.PrintTestError(t, got, []uint{9})
	}

	if got := command.TagIdsForFile(5); len(got) != 0 {
		utils.PrintTestError(t, got, []uint{})
	}
}

func TestParseCommaSeparatedUints(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []uint
	}{
		"empty string yields empty slice": {input: "", expected: []uint{}},
		"single id":                       {input: "7", expected: []uint{7}},
		"multiple ids":                    {input: "3,7,11", expected: []uint{3, 7, 11}},
		"blank segments skipped":          {input: "3,,7,", expected: []uint{3, 7}},
		"whitespace trimmed":              {input: " 3 , 7 ", expected: []uint{3, 7}},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			got, err := parseCommaSeparatedUints(test.input)
			if err != nil {
				utils.PrintTestError(t, err, nil)
			}

			if len(got) != len(test.expected) {
				utils.PrintTestError(t, got, test.expected)
				return
			}

			for i := range test.expected {
				if got[i] != test.expected[i] {
					utils.PrintTestError(t, got, test.expected)
				}
			}
		})
	}
}
