package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type UpsertPromptCommand struct {
	Name        string `gorm:"not null; uniqueIndex" json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

func (command *UpsertPromptCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	return nil
}

var allowedTemplateVariables = map[structs.PromptTemplateVariable]bool{
	structs.CATEGORIES:    true,
	structs.TAGS:          true,
	structs.OCR_TEXT:      true,
	structs.CURRENT_YEAR:  true,
	structs.CUSTOM_FIELDS: true,
}

func (command *UpsertPromptCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{}
	errorMap := make(map[string]string)

	if len(command.Name) == 0 {
		errorMap["name"] = "Name cannot be empty"
	}

	if len(command.Prompt) == 0 {
		errorMap["prompt"] = "Prompt cannot be empty"
	} else {
		regex := utils.GetTriggerRegex()
		templateVariables := regex.FindAllString(command.Prompt, -1)
		for i := 0; i < len(templateVariables); i++ {
			variable := templateVariables[i]
			if !allowedTemplateVariables[structs.PromptTemplateVariable(variable)] {
				errorMap["prompt"] = "Invalid template variables found"
				break
			}
		}
	}

	vErr.Errors = errorMap
	return vErr
}
