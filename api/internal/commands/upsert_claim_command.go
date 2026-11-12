package commands

import (
	"encoding/json"
	"net/http"
	"strings"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

const (
	MaxClaimNameLength            = 255
	MaxClaimReferenceNumberLength = 100
	MaxClaimTypeLength            = 100
	MaxClaimNotesLength           = 10000
)

// UpsertClaimCommand is a claim's own fields, for create and for a FULL replace
// on update: an optional field that is omitted, null or "" is stored as not set.
// That is safe here, unlike the pointer-means-unchanged settings commands,
// because nothing a client might hide lives on this command — receipts and ledger
// entries have their own endpoints. A field added later must be a pointer that
// leaves the stored value alone when nil, or an older client's save would clear it.
//
// It carries no group: create takes it from the path, and a claim never moves.
type UpsertClaimCommand struct {
	Name            string         `json:"name"`
	ReferenceNumber *string        `json:"referenceNumber"`
	Type            *string        `json:"type"`
	Status          string         `json:"status"`
	SubmittedDate   *string        `json:"submittedDate"`
	DueDate         *string        `json:"dueDate"`
	ClaimedAmount   OptionalAmount `json:"claimedAmount"`
	Notes           *string        `json:"notes"`
}

func (command *UpsertClaimCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, command)
}

func (command *UpsertClaimCommand) Validate() structs.ValidatorError {
	errors := structs.ValidatorError{Errors: make(map[string]string)}

	name := strings.TrimSpace(command.Name)
	if name == "" {
		errors.Errors["name"] = "Name is required"
	} else if tooLong(name, MaxClaimNameLength) {
		errors.Errors["name"] = "Name is too long"
	}

	if tooLong(trimmedOrEmpty(command.ReferenceNumber), MaxClaimReferenceNumberLength) {
		errors.Errors["referenceNumber"] = "Reference number is too long"
	}
	if tooLong(trimmedOrEmpty(command.Type), MaxClaimTypeLength) {
		errors.Errors["type"] = "Type is too long"
	}
	if command.Notes != nil && tooLong(*command.Notes, MaxClaimNotesLength) {
		errors.Errors["notes"] = "Notes are too long"
	}

	if !models.IsValidClaimStatus(command.Status) {
		errors.Errors["status"] = "Status must be one of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED"
	}

	if _, ok := optionalClaimDay(command.SubmittedDate); !ok {
		errors.Errors["submittedDate"] = "Submitted date must be a calendar day (yyyy-MM-dd)"
	}
	if _, ok := optionalClaimDay(command.DueDate); !ok {
		errors.Errors["dueDate"] = "Due date must be a calendar day (yyyy-MM-dd)"
	}

	if command.ClaimedAmount.Value != nil {
		if message := validateClaimAmount(*command.ClaimedAmount.Value, true); message != "" {
			errors.Errors["claimedAmount"] = message
		}
	}

	return errors
}

// ToClaim maps a validated command onto a claim's own columns. The group, id and
// creator are the caller's to set.
func (command *UpsertClaimCommand) ToClaim() models.Claim {
	submittedDate, _ := optionalClaimDay(command.SubmittedDate)
	dueDate, _ := optionalClaimDay(command.DueDate)
	notes := ""
	if command.Notes != nil {
		notes = *command.Notes
	}

	return models.Claim{
		Name:            strings.TrimSpace(command.Name),
		ReferenceNumber: trimmedOrEmpty(command.ReferenceNumber),
		ClaimType:       trimmedOrEmpty(command.Type),
		Status:          command.Status,
		SubmittedDate:   submittedDate,
		DueDate:         dueDate,
		ClaimedAmount:   command.ClaimedAmount.Value,
		Notes:           notes,
	}
}
