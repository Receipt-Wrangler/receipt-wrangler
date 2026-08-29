package commands

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpsertPromptCommand_Validate_ValidInputs(t *testing.T) {
	tests := map[string]struct {
		command UpsertPromptCommand
	}{
		"valid with no template vars": {
			command: UpsertPromptCommand{
				Name:   "Test Prompt",
				Prompt: "Extract receipt data from the image",
			},
		},
		"valid with all template vars": {
			command: UpsertPromptCommand{
				Name:   "Test Prompt",
				Prompt: "Use @categories @tags @ocrText @customFields and @currentYear",
			},
		},
		"valid with single template var": {
			command: UpsertPromptCommand{
				Name:   "Test Prompt",
				Prompt: "Categories are @categories",
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

func TestUpsertPromptCommand_Validate_InvalidInputs(t *testing.T) {
	tests := map[string]struct {
		command       UpsertPromptCommand
		expectedError string
	}{
		"missing name": {
			command: UpsertPromptCommand{
				Prompt: "Extract receipt data",
			},
			expectedError: "name",
		},
		"missing prompt": {
			command: UpsertPromptCommand{
				Name: "Test Prompt",
			},
			expectedError: "prompt",
		},
		"invalid template variable": {
			command: UpsertPromptCommand{
				Name:   "Test Prompt",
				Prompt: "Use @invalidVar to extract data",
			},
			expectedError: "prompt",
		},
		"mix of valid and invalid template vars": {
			command: UpsertPromptCommand{
				Name:   "Test Prompt",
				Prompt: "Use @categories and @badVar",
			},
			expectedError: "prompt",
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

func TestUpsertPromptCommand_Validate_MultipleErrors(t *testing.T) {
	command := UpsertPromptCommand{}

	vErr := command.Validate()

	if len(vErr.Errors) != 2 {
		utils.PrintTestError(t, len(vErr.Errors), 2)
	}

	if _, exists := vErr.Errors["name"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "name")
	}

	if _, exists := vErr.Errors["prompt"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "prompt")
	}
}
