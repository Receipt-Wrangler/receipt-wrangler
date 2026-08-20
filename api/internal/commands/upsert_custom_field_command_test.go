package commands

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpsertCustomFieldCommand_Validate_ValidInputs(t *testing.T) {
	tests := map[string]struct {
		command UpsertCustomFieldCommand
	}{
		"valid TEXT type": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.TEXT,
			},
		},
		"valid DATE type": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.DATE,
			},
		},
		"valid SELECT with options": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Value: "Option 1"},
				},
			},
		},
		"valid CURRENCY type": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.CURRENCY,
			},
		},
		"valid BOOLEAN type": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.BOOLEAN,
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

func TestUpsertCustomFieldCommand_Validate_InvalidInputs(t *testing.T) {
	tests := map[string]struct {
		command       UpsertCustomFieldCommand
		expectedError string
	}{
		"missing name": {
			command: UpsertCustomFieldCommand{
				Type: models.TEXT,
			},
			expectedError: "name",
		},
		"missing type": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
			},
			expectedError: "type",
		},
		"SELECT without options": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.SELECT,
			},
			expectedError: "options",
		},
		"SELECT with empty options": {
			command: UpsertCustomFieldCommand{
				Name:    "Test Field",
				Type:    models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{},
			},
			expectedError: "options",
		},
		// An option's value is its only label, and an update renames options in
		// place by id -- a blank value would keep the id and leave every receipt
		// that selected it showing nothing.
		"SELECT with a blank option value": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Value: "Option 1"},
					{Value: ""},
				},
			},
			expectedError: "options",
		},
		"SELECT with a whitespace-only option value": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Value: "   "},
				},
			},
			expectedError: "options",
		},
		"SELECT renaming an existing option to blank": {
			command: UpsertCustomFieldCommand{
				Name: "Test Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Id: 10, Value: ""},
				},
			},
			expectedError: "options",
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

func TestUpsertCustomFieldCommand_Validate_MultipleErrors(t *testing.T) {
	command := UpsertCustomFieldCommand{}

	vErr := command.Validate()

	if len(vErr.Errors) != 2 {
		utils.PrintTestError(t, len(vErr.Errors), 2)
	}

	if _, exists := vErr.Errors["name"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "name")
	}

	if _, exists := vErr.Errors["type"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "type")
	}
}

func buildExistingSelectCustomField() models.CustomField {
	existing := models.CustomField{
		Name: "Existing Field",
		Type: models.SELECT,
	}
	existing.ID = 1
	existing.Options = []models.CustomFieldOption{
		{BaseModel: models.BaseModel{ID: 10}, Value: "Option 1", CustomFieldId: 1},
		{BaseModel: models.BaseModel{ID: 11}, Value: "Option 2", CustomFieldId: 1},
	}

	return existing
}

func TestUpsertCustomFieldCommand_ValidateUpdate_ValidInputs(t *testing.T) {
	textField := models.CustomField{Name: "Existing Field", Type: models.TEXT}
	textField.ID = 2

	tests := map[string]struct {
		existing models.CustomField
		command  UpsertCustomFieldCommand
	}{
		"renaming a non-select field": {
			existing: textField,
			command: UpsertCustomFieldCommand{
				Name:        "Renamed Field",
				Type:        models.TEXT,
				Description: "Updated description",
			},
		},
		"clearing the description": {
			existing: textField,
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.TEXT,
			},
		},
		"renaming existing options by id": {
			existing: buildExistingSelectCustomField(),
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Id: 10, Value: "Renamed Option 1"},
					{Id: 11, Value: "Renamed Option 2"},
				},
			},
		},
		"appending a new option alongside existing ones": {
			existing: buildExistingSelectCustomField(),
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Id: 10, Value: "Option 1"},
					{Id: 11, Value: "Option 2"},
					{Value: "Option 3"},
				},
			},
		},
		"omitting an existing option is not a removal": {
			existing: buildExistingSelectCustomField(),
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Id: 10, Value: "Option 1"},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.ValidateUpdate(test.existing)

			if len(vErr.Errors) > 0 {
				utils.PrintTestError(t, vErr.Errors, 0)
			}
		})
	}
}

func TestUpsertCustomFieldCommand_ValidateUpdate_InvalidInputs(t *testing.T) {
	tests := map[string]struct {
		existing      models.CustomField
		command       UpsertCustomFieldCommand
		expectedError map[string]string
	}{
		"changing the type is rejected": {
			existing: buildExistingSelectCustomField(),
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.TEXT,
			},
			expectedError: map[string]string{
				"type": "Type cannot be changed",
			},
		},
		"an option id from another custom field is rejected": {
			existing: buildExistingSelectCustomField(),
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Id: 10, Value: "Option 1"},
					{Id: 999, Value: "Someone else's option"},
				},
			},
			expectedError: map[string]string{
				"options": "One or more options do not belong to this custom field",
			},
		},
		"an option id on a field with no options is rejected": {
			existing: models.CustomField{Name: "Existing Field", Type: models.SELECT},
			command: UpsertCustomFieldCommand{
				Name: "Existing Field",
				Type: models.SELECT,
				Options: []UpsertCustomFieldOptionCommand{
					{Id: 10, Value: "Option 1"},
				},
			},
			expectedError: map[string]string{
				"options": "One or more options do not belong to this custom field",
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.ValidateUpdate(test.existing)

			if len(vErr.Errors) != len(test.expectedError) {
				utils.PrintTestError(t, vErr.Errors, test.expectedError)
			}

			for key, expectedMessage := range test.expectedError {
				if vErr.Errors[key] != expectedMessage {
					utils.PrintTestError(t, vErr.Errors[key], expectedMessage)
				}
			}
		})
	}
}
