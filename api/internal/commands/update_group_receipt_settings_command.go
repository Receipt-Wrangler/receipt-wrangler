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

	// Default custom fields. BOTH are pointers, and `nil` means LEAVE UNCHANGED — the desktop hides
	// this whole section from an admin without app.custom-fields.read, so its getRawValue() omits
	// both keys. A non-pointer bool would unmarshal as `false` and the repository's unconditional
	// field-by-field assignment would wipe the stored toggle (byte for byte the documented
	// hideComments bug). An explicit empty `defaultCustomFieldIds: []` clears the set.
	//
	// Validate() deliberately does not check these: the feature has no notion of a REQUIRED custom
	// field. The handler does check that every submitted id exists (400) and that the caller holds
	// app.custom-fields.read (403).
	DefaultCustomFieldIds            *[]uint `json:"defaultCustomFieldIds"`
	ApplyDefaultCustomFieldsOnIngest *bool   `json:"applyDefaultCustomFieldsOnIngest"`

	// Receipt summary. All three are pointers with the same `nil` == LEAVE UNCHANGED meaning, for
	// three separate reasons worth keeping straight:
	//
	//   - ReceiptSummaryCustomFieldIds MUST be one: the desktop hides that control from an admin
	//     without app.custom-fields.read, so its getRawValue() omits the key entirely.
	//   - ReceiptSummaryStatuses must be one because a non-pointer slice cannot tell "the client
	//     omitted this" from the legitimate, explicit "clear every status" — both unmarshal to nil.
	//   - ReceiptSummaryEnabled must be one because a non-pointer bool unmarshals as `false` for any
	//     caller that does not send the key, and the repository's unconditional assignment would
	//     then silently switch a configured summary off.
	//
	// An explicit empty array clears the corresponding set.
	//
	// Validate() checks only the statuses (no DB needed). The handler checks that every submitted
	// custom field id exists AND is a CURRENCY field (400), and that the caller holds
	// app.custom-fields.read (403) — but only for the custom field key; see the handler.
	//
	// ReceiptSummaryPosition is a pointer for the ReceiptSummaryEnabled reason, in its enum form: a
	// non-pointer would unmarshal to "" for a caller that omits the key, and the repository's
	// assignment would then blank a configured position. It needs no permission of its own — see
	// the handler's gate, which covers the custom field key alone.
	ReceiptSummaryEnabled        *bool                          `json:"receiptSummaryEnabled"`
	ReceiptSummaryPosition       *models.ReceiptSummaryPosition `json:"receiptSummaryPosition"`
	ReceiptSummaryCustomFieldIds *[]uint                        `json:"receiptSummaryCustomFieldIds"`
	ReceiptSummaryStatuses       *[]models.ReceiptStatus        `json:"receiptSummaryStatuses"`
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

	// Every summary breakdown status must be a real one. Note "" is rejected: ReceiptStatuses()
	// does not list it, and an empty status is never a breakdown row — unlike the quick-scan
	// default above, where ReceiptStatus.Value() tolerates "" for an unset scalar. Checking
	// membership here rather than leaving it to the DB layer is the point: ReceiptStatus.Value()
	// would surface a bogus value as a generic 500 (see api/CLAUDE.md -> "Receipt statuses").
	if command.ReceiptSummaryStatuses != nil {
		for _, status := range *command.ReceiptSummaryStatuses {
			if !isValidReceiptStatus(status) {
				vErr.Errors["receiptSummaryStatuses"] = "One or more selected statuses is not a valid receipt status"
				break
			}
		}
	}

	// The position goes through its own Valuer rather than a membership walk, the
	// pie_chart_data_command.go style: two values make a list overkill, and routing it here turns
	// what the DB layer would surface as a generic 500 into a field-level 400.
	//
	// The EMPTY check is separate and load-bearing. Value() tolerates "" (every enum here does, for
	// unset scalars), and "" is a real member of the swagger enum -- it has to be, so the generated
	// Dart enum has a member to carry `fallback: true`. But on the WRITE side an explicit "" is not
	// "leave unchanged": the pointer being non-nil makes the repository assign it, silently
	// resetting a configured position to the default with a 200. Omitting the key is the only way
	// to leave it alone, so an explicit empty is a 400.
	if command.ReceiptSummaryPosition != nil {
		_, err := command.ReceiptSummaryPosition.Value()
		if err != nil || len(*command.ReceiptSummaryPosition) == 0 {
			vErr.Errors["receiptSummaryPosition"] = "Invalid receipt summary position. Must be TOP or BOTTOM"
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
