package commands

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
)

func decodeClaimCommand(t *testing.T, body string) UpsertClaimCommand {
	t.Helper()
	command := UpsertClaimCommand{}
	if err := json.Unmarshal([]byte(body), &command); err != nil {
		t.Fatalf("decode %s: %v", body, err)
	}
	return command
}

func TestUpsertClaimCommand_ValidCommandMapsToClaim(t *testing.T) {
	command := decodeClaimCommand(t, `{
		"name": "  PCS Rome  ",
		"referenceNumber": " TA-123 ",
		"type": "Travel voucher",
		"status": "SUBMITTED",
		"submittedDate": "2026-09-01",
		"dueDate": "2026-10-01T00:00:00-04:00",
		"claimedAmount": "412.50",
		"notes": "sent by email"
	}`)

	if vErr := command.Validate(); len(vErr.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", vErr.Errors)
	}

	claim := command.ToClaim()
	if claim.Name != "PCS Rome" || claim.ReferenceNumber != "TA-123" || claim.ClaimType != "Travel voucher" {
		utils.PrintTestError(t, claim, "trimmed text fields")
	}
	if claim.SubmittedDate == nil || !claim.SubmittedDate.Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)) {
		utils.PrintTestError(t, claim.SubmittedDate, "2026-09-01 midnight UTC")
	}
	// An RFC 3339 value is read as the day written, never converted to UTC first.
	if claim.DueDate == nil || !claim.DueDate.Equal(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)) {
		utils.PrintTestError(t, claim.DueDate, "2026-10-01 midnight UTC")
	}
	if claim.ClaimedAmount == nil || claim.ClaimedAmount.StringFixed(2) != "412.50" {
		utils.PrintTestError(t, claim.ClaimedAmount, "412.50")
	}
}

// The desktop's masked currency input sends "" for an emptied field; that, null,
// and an omitted key all mean "not set", as do blank dates.
func TestUpsertClaimCommand_BlankOptionalFieldsAreNotSet(t *testing.T) {
	bodies := []string{
		`{"name":"a","status":"DRAFT","claimedAmount":"","submittedDate":"","dueDate":null}`,
		`{"name":"a","status":"DRAFT","claimedAmount":null}`,
		`{"name":"a","status":"DRAFT","claimedAmount":"   "}`,
		`{"name":"a","status":"DRAFT"}`,
	}

	for _, body := range bodies {
		command := decodeClaimCommand(t, body)
		if vErr := command.Validate(); len(vErr.Errors) > 0 {
			t.Errorf("%s: unexpected errors %v", body, vErr.Errors)
			continue
		}
		claim := command.ToClaim()
		if claim.ClaimedAmount != nil || claim.SubmittedDate != nil || claim.DueDate != nil {
			t.Errorf("%s: optional fields should be unset, got %+v", body, claim)
		}
	}
}

func TestUpsertClaimCommand_NumericAmountIsAccepted(t *testing.T) {
	command := decodeClaimCommand(t, `{"name":"a","status":"DRAFT","claimedAmount":12.5}`)
	if command.ClaimedAmount.Value == nil || command.ClaimedAmount.Value.StringFixed(2) != "12.50" {
		utils.PrintTestError(t, command.ClaimedAmount.Value, "12.50")
	}
}

func TestUpsertClaimCommand_UnparseableAmountFailsToDecode(t *testing.T) {
	command := UpsertClaimCommand{}
	if err := json.Unmarshal([]byte(`{"name":"a","status":"DRAFT","claimedAmount":"abc"}`), &command); err == nil {
		t.Errorf("expected a decode error for a non-numeric amount")
	}
}

func TestUpsertClaimCommand_Rejections(t *testing.T) {
	cases := map[string]struct {
		body string
		key  string
	}{
		"blank name":         {`{"name":"   ","status":"DRAFT"}`, "name"},
		"long name":          {`{"name":"` + strings.Repeat("x", 256) + `","status":"DRAFT"}`, "name"},
		"unknown status":     {`{"name":"a","status":"LOST"}`, "status"},
		"missing status":     {`{"name":"a"}`, "status"},
		"bad submitted date": {`{"name":"a","status":"DRAFT","submittedDate":"09/01/2026"}`, "submittedDate"},
		"bad due date":       {`{"name":"a","status":"DRAFT","dueDate":"tomorrow"}`, "dueDate"},
		"negative amount":    {`{"name":"a","status":"DRAFT","claimedAmount":"-1"}`, "claimedAmount"},
		"three decimals":     {`{"name":"a","status":"DRAFT","claimedAmount":"1.005"}`, "claimedAmount"},
		"too large":          {`{"name":"a","status":"DRAFT","claimedAmount":"100000000"}`, "claimedAmount"},
		"long reference":     {`{"name":"a","status":"DRAFT","referenceNumber":"` + strings.Repeat("x", 101) + `"}`, "referenceNumber"},
		"long type":          {`{"name":"a","status":"DRAFT","type":"` + strings.Repeat("x", 101) + `"}`, "type"},
		"long notes":         {`{"name":"a","status":"DRAFT","notes":"` + strings.Repeat("x", 10001) + `"}`, "notes"},
	}

	for name, test := range cases {
		command := decodeClaimCommand(t, test.body)
		vErr := command.Validate()
		if _, ok := vErr.Errors[test.key]; !ok {
			t.Errorf("%s: expected an error keyed %q, got %v", name, test.key, vErr.Errors)
		}
	}
}

func TestUpsertClaimCommand_ZeroAmountAndEveryStatusAreAccepted(t *testing.T) {
	for _, status := range []string{
		models.ClaimStatusDraft, models.ClaimStatusSubmitted, models.ClaimStatusPartiallyPaid,
		models.ClaimStatusPaid, models.ClaimStatusDisputed, models.ClaimStatusClosed,
	} {
		command := decodeClaimCommand(t, `{"name":"a","status":"`+status+`","claimedAmount":"0"}`)
		if vErr := command.Validate(); len(vErr.Errors) > 0 {
			t.Errorf("%s: unexpected errors %v", status, vErr.Errors)
		}
	}
}

func TestUpsertClaimLedgerEntryCommand_Rules(t *testing.T) {
	cases := map[string]struct {
		body    string
		wantKey string // "" means valid
	}{
		"payment":                {`{"type":"PAYMENT","date":"2026-10-01","amount":"100.00"}`, ""},
		"recoupment":             {`{"type":"RECOUPMENT","date":"2026-10-01","amount":"5"}`, ""},
		"note":                   {`{"type":"NOTE","date":"2026-10-01","note":"called the office"}`, ""},
		"note with amount":       {`{"type":"NOTE","date":"2026-10-01","note":"disallowed","amount":"45.00"}`, ""},
		"payment without amount": {`{"type":"PAYMENT","date":"2026-10-01","amount":""}`, "amount"},
		"recoupment no amount":   {`{"type":"RECOUPMENT","date":"2026-10-01"}`, "amount"},
		"zero amount":            {`{"type":"PAYMENT","date":"2026-10-01","amount":"0"}`, "amount"},
		"negative amount":        {`{"type":"PAYMENT","date":"2026-10-01","amount":"-5"}`, "amount"},
		"note without text":      {`{"type":"NOTE","date":"2026-10-01","note":"  "}`, "note"},
		"unknown type":           {`{"type":"REFUND","date":"2026-10-01","amount":"1"}`, "type"},
		"missing date":           {`{"type":"PAYMENT","amount":"1"}`, "date"},
		"bad date":               {`{"type":"PAYMENT","date":"Oct 1","amount":"1"}`, "date"},
		"long note":              {`{"type":"NOTE","date":"2026-10-01","note":"` + strings.Repeat("x", 2001) + `"}`, "note"},
	}

	for name, test := range cases {
		command := UpsertClaimLedgerEntryCommand{}
		if err := json.Unmarshal([]byte(test.body), &command); err != nil {
			t.Fatalf("%s: decode: %v", name, err)
		}
		vErr := command.Validate()
		if test.wantKey == "" {
			if len(vErr.Errors) > 0 {
				t.Errorf("%s: unexpected errors %v", name, vErr.Errors)
			}
			continue
		}
		if _, ok := vErr.Errors[test.wantKey]; !ok {
			t.Errorf("%s: expected an error keyed %q, got %v", name, test.wantKey, vErr.Errors)
		}
	}
}

func TestUpsertClaimLedgerEntryCommand_ToEntry(t *testing.T) {
	command := UpsertClaimLedgerEntryCommand{}
	if err := json.Unmarshal([]byte(`{"type":"PAYMENT","date":"2026-10-01","amount":"100.10","note":"EFT"}`), &command); err != nil {
		t.Fatalf("decode: %v", err)
	}

	entry := command.ToClaimLedgerEntry()
	if entry.EntryType != models.ClaimLedgerEntryPayment || entry.Note != "EFT" {
		utils.PrintTestError(t, entry, "PAYMENT with its note")
	}
	if !entry.Date.Equal(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)) {
		utils.PrintTestError(t, entry.Date, "2026-10-01 midnight UTC")
	}
	if entry.Amount == nil || entry.Amount.StringFixed(2) != "100.10" {
		utils.PrintTestError(t, entry.Amount, "100.10")
	}
}

func TestClaimReceiptsCommand_Validate(t *testing.T) {
	tooMany := make([]uint, MaxClaimReceiptIds+1)
	for i := range tooMany {
		tooMany[i] = uint(i + 1)
	}

	cases := map[string]struct {
		ids   []uint
		valid bool
	}{
		"one":        {[]uint{1}, true},
		"empty":      {[]uint{}, false},
		"nil":        {nil, false},
		"duplicates": {[]uint{1, 1}, false},
		"too many":   {tooMany, false},
	}

	for name, test := range cases {
		command := ClaimReceiptsCommand{ReceiptIds: test.ids}
		vErr := command.Validate()
		if test.valid != (len(vErr.Errors) == 0) {
			t.Errorf("%s: valid=%v, errors %v", name, test.valid, vErr.Errors)
		}
	}
}
