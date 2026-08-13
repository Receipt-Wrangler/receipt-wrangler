package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type UpdateGroupReceiptSettingsCommand struct {
	HideImages            bool `json:"hideImages"`
	HideReceiptCategories bool `json:"hideReceiptCategories"`
	HideReceiptTags       bool `json:"hideReceiptTags"`
	HideItemCategories    bool `json:"hideItemCategories"`
	HideItemTags          bool `json:"hideItemTags"`
	HideComments          bool `json:"hideComments"`
	HideShareCategories   bool `json:"hideShareCategories"`
	HideShareTags         bool `json:"hideShareTags"`

	QuickScanPaidByEnabled     bool                              `json:"quickScanPaidByEnabled"`
	QuickScanPaidByRequired    bool                              `json:"quickScanPaidByRequired"`
	QuickScanDefaultPaidByType models.QuickScanDefaultPaidByType `json:"quickScanDefaultPaidByType"`
	QuickScanDefaultPaidById   *uint                             `json:"quickScanDefaultPaidById"`

	QuickScanStatusEnabled  bool                 `json:"quickScanStatusEnabled"`
	QuickScanStatusRequired bool                 `json:"quickScanStatusRequired"`
	QuickScanDefaultStatus  models.ReceiptStatus `json:"quickScanDefaultStatus"`

	QuickScanCategoriesEnabled  bool `json:"quickScanCategoriesEnabled"`
	QuickScanCategoriesRequired bool `json:"quickScanCategoriesRequired"`

	QuickScanTagsEnabled  bool `json:"quickScanTagsEnabled"`
	QuickScanTagsRequired bool `json:"quickScanTagsRequired"`

	QuickScanCommentEnabled  bool `json:"quickScanCommentEnabled"`
	QuickScanCommentRequired bool `json:"quickScanCommentRequired"`
}

func (command *UpdateGroupReceiptSettingsCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
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

func (command UpdateGroupReceiptSettingsCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{
		Errors: make(map[string]string),
	}

	// When paid-by is not both shown and required, the user may skip it, so a default must be
	// configured to backfill the value (a receipt always has a paid-by).
	if !(command.QuickScanPaidByEnabled && command.QuickScanPaidByRequired) {
		switch command.QuickScanDefaultPaidByType {
		case models.QUICK_SCAN_PAID_BY_UPLOADER:
			// Uploader resolves at scan time; no id needed.
		case models.QUICK_SCAN_PAID_BY_USER:
			if command.QuickScanDefaultPaidById == nil || *command.QuickScanDefaultPaidById == 0 {
				vErr.Errors["quickScanDefaultPaidById"] = "A default paid by user is required when paid by is optional"
			}
		default:
			vErr.Errors["quickScanDefaultPaidByType"] = "A default paid by is required when paid by is optional"
		}
	}

	// Same rule for status.
	if !(command.QuickScanStatusEnabled && command.QuickScanStatusRequired) {
		if !isValidReceiptStatus(command.QuickScanDefaultStatus) {
			vErr.Errors["quickScanDefaultStatus"] = "A default status is required when status is optional"
		}
	}

	return vErr
}

func (command *UpdateGroupReceiptSettingsCommand) LoadDataFromRequestAndValidate(w http.ResponseWriter, r *http.Request) (structs.ValidatorError, error) {
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		return structs.ValidatorError{}, err
	}

	vErr := command.Validate()
	if len(vErr.Errors) > 0 {
		return vErr, nil
	}

	return structs.ValidatorError{}, nil
}

func isValidReceiptStatus(status models.ReceiptStatus) bool {
	for _, valid := range models.ReceiptStatuses() {
		if valid == status {
			return true
		}
	}
	return false
}
