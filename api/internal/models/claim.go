package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Claim statuses. A claim's status is set by hand and any status may follow any
// other; the server only checks that a written value is one of these. It is a
// plain string on the wire (never an OpenAPI enum), so a client tolerates a value
// added later, and a stored value this build does not know is still served.
const (
	ClaimStatusDraft         = "DRAFT"
	ClaimStatusSubmitted     = "SUBMITTED"
	ClaimStatusPartiallyPaid = "PARTIALLY_PAID"
	ClaimStatusPaid          = "PAID"
	ClaimStatusDisputed      = "DISPUTED"
	ClaimStatusClosed        = "CLOSED"
)

// IsValidClaimStatus reports whether status may be written to a claim.
func IsValidClaimStatus(status string) bool {
	switch status {
	case ClaimStatusDraft, ClaimStatusSubmitted, ClaimStatusPartiallyPaid,
		ClaimStatusPaid, ClaimStatusDisputed, ClaimStatusClosed:
		return true
	}
	return false
}

// ClaimStatusAwaitsSettlement reports whether a claim in this status is waiting
// on the other party: it is what daysWaiting counts and what the list summary's
// outstanding total covers. An unknown status awaits nothing.
func ClaimStatusAwaitsSettlement(status string) bool {
	switch status {
	case ClaimStatusSubmitted, ClaimStatusPartiallyPaid, ClaimStatusDisputed:
		return true
	}
	return false
}

// Claim is a reimbursement claim: a bundle of one group's receipts submitted
// together, plus a ledger of what happened after submission (ClaimLedgerEntry).
// Receipts are linked through ClaimReceipt rather than a column on the receipt,
// so the receipt model, its payloads and its change history are untouched.
//
// SubmittedDate and DueDate are calendar days, stored as midnight UTC of the day
// the client wrote, exactly like Receipt.Date, and always read back in UTC.
//
// Deleting the group cascades to its claims (and on to their entries and links),
// because GroupService.DeleteGroup removes the group row outright. A RESTRICT
// here would make group deletion fail with an opaque error instead.
type Claim struct {
	BaseModel
	GroupId         uint             `gorm:"not null;index" json:"groupId"`
	Group           *Group           `gorm:"foreignKey:GroupId;constraint:OnDelete:CASCADE" json:"-"`
	Name            string           `gorm:"size:255;not null" json:"name"`
	ReferenceNumber string           `gorm:"size:100" json:"referenceNumber"`
	ClaimType       string           `gorm:"column:claim_type;size:100" json:"type"`
	Status          string           `gorm:"size:32;not null;default:'DRAFT'" json:"status"`
	SubmittedDate   *time.Time       `json:"submittedDate"`
	DueDate         *time.Time       `json:"dueDate"`
	ClaimedAmount   *decimal.Decimal `gorm:"type:decimal(10,2)" json:"claimedAmount"`
	Notes           string           `gorm:"type:text" json:"notes"`
}

// ClaimReceipt links a receipt to the claim it was submitted in. ReceiptId is
// the primary key, which is what makes "a receipt is in at most one claim" a
// database guarantee rather than a convention.
//
// Both foreign keys cascade. The receipt one only fires on paths that delete a
// receipt without asking (group delete, user delete): the receipt delete
// endpoints refuse a receipt that is in a claim first, in the service layer (see
// ClaimService.ReceiptDeletionRefusals). A receipt moved to another group is
// unlinked by ReceiptRepository.UpdateReceipt, since a claim only holds receipts
// of its own group.
type ClaimReceipt struct {
	ReceiptId uint      `gorm:"primaryKey;autoIncrement:false" json:"receiptId"`
	Receipt   *Receipt  `gorm:"foreignKey:ReceiptId;constraint:OnDelete:CASCADE" json:"-"`
	ClaimId   uint      `gorm:"not null;index" json:"claimId"`
	Claim     *Claim    `gorm:"foreignKey:ClaimId;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

// Claim ledger entry types. PAYMENT adds to what was paid, RECOUPMENT (money the
// other party took back) subtracts from it, and NOTE affects no figure. Plain
// strings on the wire, like the claim status; an entry of a type this build does
// not know counts toward nothing.
const (
	ClaimLedgerEntryPayment    = "PAYMENT"
	ClaimLedgerEntryRecoupment = "RECOUPMENT"
	ClaimLedgerEntryNote       = "NOTE"
)

// IsValidClaimLedgerEntryType reports whether entryType may be written.
func IsValidClaimLedgerEntryType(entryType string) bool {
	switch entryType {
	case ClaimLedgerEntryPayment, ClaimLedgerEntryRecoupment, ClaimLedgerEntryNote:
		return true
	}
	return false
}

// ClaimLedgerEntry is one dated event on a claim: a payment received, a
// recoupment, or a note. Each is its own row so the audit trail is individually
// addressable. Amount is always positive; the type gives the direction. Date is a
// calendar day stored as midnight UTC.
type ClaimLedgerEntry struct {
	BaseModel
	ClaimId   uint             `gorm:"not null;index" json:"claimId"`
	Claim     *Claim           `gorm:"foreignKey:ClaimId;constraint:OnDelete:CASCADE" json:"-"`
	EntryType string           `gorm:"column:entry_type;size:32;not null" json:"type"`
	Date      time.Time        `gorm:"not null" json:"date"`
	Amount    *decimal.Decimal `gorm:"type:decimal(10,2)" json:"amount"`
	Note      string           `gorm:"size:2000" json:"note"`
}
