package commands

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpdateGroupReceiptSettingsCommand_Validate_ValidInputs(t *testing.T) {
	userId := uint(2)

	tests := map[string]struct {
		command UpdateGroupReceiptSettingsCommand
	}{
		"paid by and status shown+required need no defaults": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:  true,
				QuickScanPaidByRequired: true,
				QuickScanStatusEnabled:  true,
				QuickScanStatusRequired: true,
			},
		},
		// Unlike paid-by/status, a comment can legitimately be empty, so requiring it needs no
		// configured default.
		"comment shown+required needs no default": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:   true,
				QuickScanPaidByRequired:  true,
				QuickScanStatusEnabled:   true,
				QuickScanStatusRequired:  true,
				QuickScanCommentEnabled:  true,
				QuickScanCommentRequired: true,
			},
		},
		"paid by optional with uploader default": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:     true,
				QuickScanPaidByRequired:    false,
				QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
				QuickScanStatusEnabled:     true,
				QuickScanStatusRequired:    true,
			},
		},
		"paid by optional with specific user default": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:     true,
				QuickScanPaidByRequired:    false,
				QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_USER,
				QuickScanDefaultPaidById:   &userId,
				QuickScanStatusEnabled:     true,
				QuickScanStatusRequired:    true,
			},
		},
		"paid by hidden with uploader default": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:     false,
				QuickScanPaidByRequired:    false,
				QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_UPLOADER,
				QuickScanStatusEnabled:     true,
				QuickScanStatusRequired:    true,
			},
		},
		"status optional with default status": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:  true,
				QuickScanPaidByRequired: true,
				QuickScanStatusEnabled:  true,
				QuickScanStatusRequired: false,
				QuickScanDefaultStatus:  models.OPEN,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate()

			if len(vErr.Errors) > 0 {
				utils.PrintTestError(t, vErr.Errors, "no errors")
			}
		})
	}
}

func TestUpdateGroupReceiptSettingsCommand_Validate_InvalidInputs(t *testing.T) {
	zero := uint(0)

	tests := map[string]struct {
		command     UpdateGroupReceiptSettingsCommand
		expectedKey string
	}{
		"paid by optional with no default type": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:  true,
				QuickScanPaidByRequired: false,
				QuickScanStatusEnabled:  true,
				QuickScanStatusRequired: true,
			},
			expectedKey: "quickScanDefaultPaidByType",
		},
		"paid by optional type user without id": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:     true,
				QuickScanPaidByRequired:    false,
				QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_USER,
				QuickScanStatusEnabled:     true,
				QuickScanStatusRequired:    true,
			},
			expectedKey: "quickScanDefaultPaidById",
		},
		"paid by optional type user with zero id": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:     true,
				QuickScanPaidByRequired:    false,
				QuickScanDefaultPaidByType: models.QUICK_SCAN_PAID_BY_USER,
				QuickScanDefaultPaidById:   &zero,
				QuickScanStatusEnabled:     true,
				QuickScanStatusRequired:    true,
			},
			expectedKey: "quickScanDefaultPaidById",
		},
		"status optional with no default": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:  true,
				QuickScanPaidByRequired: true,
				QuickScanStatusEnabled:  true,
				QuickScanStatusRequired: false,
			},
			expectedKey: "quickScanDefaultStatus",
		},
		"status optional with invalid default": {
			command: UpdateGroupReceiptSettingsCommand{
				QuickScanPaidByEnabled:  true,
				QuickScanPaidByRequired: true,
				QuickScanStatusEnabled:  true,
				QuickScanStatusRequired: false,
				QuickScanDefaultStatus:  models.ReceiptStatus("NOT_A_STATUS"),
			},
			expectedKey: "quickScanDefaultStatus",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate()

			if _, ok := vErr.Errors[test.expectedKey]; !ok {
				utils.PrintTestError(t, vErr.Errors, test.expectedKey)
			}
		})
	}
}
