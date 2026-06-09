package commands

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func validSystemSettingsCommand() UpsertSystemSettingsCommand {
	queueNames := models.GetQueueNames()
	configs := make([]UpsertTaskQueueConfigurationCommand, len(queueNames))

	return UpsertSystemSettingsCommand{
		CurrencySymbolPosition:       models.START,
		CurrencyThousandthsSeparator: models.COMMA,
		CurrencyDecimalSeparator:     models.DOT,
		TaskConcurrency:              1,
		EmailPollingInterval:         60,
		TaskQueueConfigurations:      configs,
	}
}

func TestUpsertSystemSettingsCommand_Validate_ValidInputs(t *testing.T) {
	primaryId := uint(1)
	fallbackId := uint(2)

	tests := map[string]struct {
		command UpsertSystemSettingsCommand
	}{
		"valid minimal": {
			command: validSystemSettingsCommand(),
		},
		"valid with processing settings": {
			command: func() UpsertSystemSettingsCommand {
				cmd := validSystemSettingsCommand()
				cmd.ReceiptProcessingSettingsId = &primaryId
				cmd.FallbackReceiptProcessingSettingsId = &fallbackId
				return cmd
			}(),
		},
		"valid with zero email polling interval": {
			command: func() UpsertSystemSettingsCommand {
				cmd := validSystemSettingsCommand()
				cmd.EmailPollingInterval = 0
				return cmd
			}(),
		},
		"valid with zero task concurrency": {
			command: func() UpsertSystemSettingsCommand {
				cmd := validSystemSettingsCommand()
				cmd.TaskConcurrency = 0
				return cmd
			}(),
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

func TestUpsertSystemSettingsCommand_Validate_InvalidInputs(t *testing.T) {
	zeroId := uint(0)
	sameId := uint(5)

	tests := map[string]struct {
		modify        func(cmd *UpsertSystemSettingsCommand)
		expectedError string
	}{
		"negative email polling interval": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.EmailPollingInterval = -1 },
			expectedError: "emailPollingInterval",
		},
		"invalid receipt processing settings id": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.ReceiptProcessingSettingsId = &zeroId },
			expectedError: "receiptProcessingSettingsId",
		},
		"invalid fallback receipt processing settings id": {
			modify: func(cmd *UpsertSystemSettingsCommand) {
				id := uint(1)
				cmd.ReceiptProcessingSettingsId = &id
				cmd.FallbackReceiptProcessingSettingsId = &zeroId
			},
			expectedError: "fallbackReceiptProcessingSettingsId",
		},
		"fallback without primary": {
			modify: func(cmd *UpsertSystemSettingsCommand) {
				id := uint(1)
				cmd.FallbackReceiptProcessingSettingsId = &id
			},
			expectedError: "fallbackReceiptProcessingSettingsId",
		},
		"fallback same as primary": {
			modify: func(cmd *UpsertSystemSettingsCommand) {
				cmd.ReceiptProcessingSettingsId = &sameId
				sameIdCopy := sameId
				cmd.FallbackReceiptProcessingSettingsId = &sameIdCopy
			},
			expectedError: "fallbackReceiptProcessingSettingsId",
		},
		"missing currency symbol position": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.CurrencySymbolPosition = "" },
			expectedError: "currencySymbolPosition",
		},
		"missing currency thousandths separator": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.CurrencyThousandthsSeparator = "" },
			expectedError: "currencyThousandthsSeparator",
		},
		"missing currency decimal separator": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.CurrencyDecimalSeparator = "" },
			expectedError: "currencyDecimalSeparator",
		},
		"negative task concurrency": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.TaskConcurrency = -1 },
			expectedError: "taskConcurrency",
		},
		"wrong queue config count": {
			modify:        func(cmd *UpsertSystemSettingsCommand) { cmd.TaskQueueConfigurations = []UpsertTaskQueueConfigurationCommand{} },
			expectedError: "taskQueueConfigurations",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			cmd := validSystemSettingsCommand()
			test.modify(&cmd)

			vErr := cmd.Validate()

			if len(vErr.Errors) == 0 {
				utils.PrintTestError(t, len(vErr.Errors), "greater than 0")
			}

			if _, exists := vErr.Errors[test.expectedError]; !exists {
				utils.PrintTestError(t, "error should exist for field", test.expectedError)
			}
		})
	}
}

func TestUpsertSystemSettingsCommand_Validate_MultipleErrors(t *testing.T) {
	command := UpsertSystemSettingsCommand{
		EmailPollingInterval: -1,
		TaskConcurrency:      -1,
	}

	vErr := command.Validate()

	if len(vErr.Errors) < 5 {
		utils.PrintTestError(t, len(vErr.Errors), "at least 5")
	}
}

func TestUpsertSystemSettingsCommand_Validate_PdfDpi(t *testing.T) {
	// 0 means "unset / use default" and must be allowed; in-range values pass;
	// out-of-range values are rejected with a pdfDpi error.
	validValues := []int{0, 72, 300, 1200}
	for _, v := range validValues {
		cmd := validSystemSettingsCommand()
		cmd.PdfDpi = v
		vErr := cmd.Validate()
		if _, ok := vErr.Errors["pdfDpi"]; ok {
			utils.PrintTestError(t, "pdfDpi error for valid value "+utils.UintToString(uint(v)), "no error")
		}
	}

	invalidValues := []int{71, 1201, 5000}
	for _, v := range invalidValues {
		cmd := validSystemSettingsCommand()
		cmd.PdfDpi = v
		vErr := cmd.Validate()
		if _, ok := vErr.Errors["pdfDpi"]; !ok {
			utils.PrintTestError(t, "no pdfDpi error for invalid value "+utils.UintToString(uint(v)), "pdfDpi error")
		}
	}
}
