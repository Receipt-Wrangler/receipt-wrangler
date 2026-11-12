package commands

import (
	"encoding/json"
	"net/http"
	"strings"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

const MaxClaimLedgerEntryNoteLength = 2000

// UpsertClaimLedgerEntryCommand is a ledger entry's fields, for create and for a
// full replace on update. An omitted, null or "" optional field is stored as not
// set. The claim comes from the path, never the body.
type UpsertClaimLedgerEntryCommand struct {
	Type   string         `json:"type"`
	Date   string         `json:"date"`
	Amount OptionalAmount `json:"amount"`
	Note   *string        `json:"note"`
}

func (command *UpsertClaimLedgerEntryCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, command)
}

// Validate requires an amount on a PAYMENT or RECOUPMENT, since those are the
// entries the figures are built from, and a note on a NOTE, which has nothing
// else to say. Amounts are always positive: the type gives the direction.
func (command *UpsertClaimLedgerEntryCommand) Validate() structs.ValidatorError {
	errors := structs.ValidatorError{Errors: make(map[string]string)}

	if !models.IsValidClaimLedgerEntryType(command.Type) {
		errors.Errors["type"] = "Type must be one of PAYMENT, RECOUPMENT, NOTE"
	}

	if _, ok := parseClaimCalendarDay(command.Date); !ok {
		errors.Errors["date"] = "Date is required and must be a calendar day (yyyy-MM-dd)"
	}

	if command.Amount.Value != nil {
		if message := validateClaimAmount(*command.Amount.Value, false); message != "" {
			errors.Errors["amount"] = message
		}
	} else if command.Type == models.ClaimLedgerEntryPayment || command.Type == models.ClaimLedgerEntryRecoupment {
		errors.Errors["amount"] = "Amount is required for a payment or a recoupment"
	}

	note := ""
	if command.Note != nil {
		note = *command.Note
	}
	if tooLong(note, MaxClaimLedgerEntryNoteLength) {
		errors.Errors["note"] = "Note is too long"
	} else if command.Type == models.ClaimLedgerEntryNote && strings.TrimSpace(note) == "" {
		errors.Errors["note"] = "A note is required"
	}

	return errors
}

// ToClaimLedgerEntry maps a validated command onto an entry's columns. The claim,
// id and creator are the caller's to set.
func (command *UpsertClaimLedgerEntryCommand) ToClaimLedgerEntry() models.ClaimLedgerEntry {
	date, _ := parseClaimCalendarDay(command.Date)
	note := ""
	if command.Note != nil {
		note = *command.Note
	}

	return models.ClaimLedgerEntry{
		EntryType: command.Type,
		Date:      date,
		Amount:    command.Amount.Value,
		Note:      note,
	}
}
