package repositories

import (
	"receipt-wrangler/api/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ClaimRepository is the data access for reimbursement claims, their receipt
// links and their ledger entries. Visibility and authorization live in
// services.ClaimService; nothing here filters by caller.
type ClaimRepository struct {
	BaseRepository
}

func NewClaimRepository(tx *gorm.DB) ClaimRepository {
	return ClaimRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
}

// makeClaimMigrations creates the claims tables. It runs from MakeMigrations,
// after the main list, so the receipts and groups tables its foreign keys
// reference already exist. Kept apart from that list so the fork's models do not
// sit in a block upstream edits.
func makeClaimMigrations() error {
	return db.AutoMigrate(
		&models.Claim{},
		&models.ClaimReceipt{},
		&models.ClaimLedgerEntry{},
	)
}

// LinkedReceipt is a receipt as a claim sees it: just what visibility and the
// figures need.
type LinkedReceipt struct {
	ClaimId      uint
	ReceiptId    uint
	GroupId      uint
	PaidByUserId uint
	Amount       decimal.Decimal
}

// ReceiptClaimLink names the claim a receipt is in.
type ReceiptClaimLink struct {
	ReceiptId      uint
	ClaimId        uint
	ClaimName      string
	ClaimGroupId   uint
	ClaimCreatedBy *uint
}

func (repository ClaimRepository) GetClaimById(claimId uint) (models.Claim, error) {
	var claim models.Claim
	err := repository.GetDB().Where("id = ?", claimId).First(&claim).Error
	return claim, err
}

// GetClaimsByGroupIds returns the claims of the given groups, newest first.
func (repository ClaimRepository) GetClaimsByGroupIds(groupIds []uint) ([]models.Claim, error) {
	claims := []models.Claim{}
	if len(groupIds) == 0 {
		return claims, nil
	}

	err := repository.GetDB().
		Where("group_id IN ?", groupIds).
		Order("id DESC").
		Find(&claims).Error
	return claims, err
}

func (repository ClaimRepository) CreateClaim(claim models.Claim) (models.Claim, error) {
	err := repository.GetDB().Omit(clause.Associations).Create(&claim).Error
	return claim, err
}

// UpdateClaim replaces a claim's own columns. The map form is used so a field
// cleared to nil or "" is actually written: the struct form skips zero values.
func (repository ClaimRepository) UpdateClaim(claimId uint, claim models.Claim) error {
	return repository.GetDB().Model(&models.Claim{}).Where("id = ?", claimId).Updates(map[string]interface{}{
		"name":             claim.Name,
		"reference_number": claim.ReferenceNumber,
		"claim_type":       claim.ClaimType,
		"status":           claim.Status,
		"submitted_date":   nullableValue(claim.SubmittedDate),
		"due_date":         nullableValue(claim.DueDate),
		"claimed_amount":   nullableValue(claim.ClaimedAmount),
		"notes":            claim.Notes,
	}).Error
}

// DeleteClaim deletes a claim with its ledger entries and receipt links. The
// children are deleted explicitly rather than left to the cascading foreign keys,
// so a connection without foreign key enforcement still leaves nothing behind.
func (repository ClaimRepository) DeleteClaim(claimId uint) error {
	return repository.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("claim_id = ?", claimId).Delete(&models.ClaimLedgerEntry{}).Error; err != nil {
			return err
		}
		if err := tx.Where("claim_id = ?", claimId).Delete(&models.ClaimReceipt{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", claimId).Delete(&models.Claim{}).Error
	})
}

// GetLinkedReceipts returns the receipts linked to the given claims, by claim,
// ordered by receipt id. A link whose receipt now lives in another group than the
// claim is left out: a claim only ever holds receipts of its own group, and the
// group move that would break that unlinks the receipt (DetachReceiptFromClaims),
// so this is a second line of defence rather than the rule itself.
func (repository ClaimRepository) GetLinkedReceipts(claimIds []uint) (map[uint][]LinkedReceipt, error) {
	result := make(map[uint][]LinkedReceipt)
	if len(claimIds) == 0 {
		return result, nil
	}

	var rows []LinkedReceipt
	err := repository.GetDB().
		Table("claim_receipts").
		Select("claim_receipts.claim_id AS claim_id, receipts.id AS receipt_id, receipts.group_id AS group_id, receipts.paid_by_user_id AS paid_by_user_id, receipts.amount AS amount").
		Joins("JOIN receipts ON receipts.id = claim_receipts.receipt_id").
		Joins("JOIN claims ON claims.id = claim_receipts.claim_id").
		Where("claim_receipts.claim_id IN ? AND receipts.group_id = claims.group_id", claimIds).
		Order("receipts.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.ClaimId] = append(result[row.ClaimId], row)
	}
	return result, nil
}

// GetLedgerEntries returns the given claims' ledger entries, by claim, oldest
// first (by date, then id).
func (repository ClaimRepository) GetLedgerEntries(claimIds []uint) (map[uint][]models.ClaimLedgerEntry, error) {
	result := make(map[uint][]models.ClaimLedgerEntry)
	if len(claimIds) == 0 {
		return result, nil
	}

	var entries []models.ClaimLedgerEntry
	err := repository.GetDB().
		Where("claim_id IN ?", claimIds).
		Order("date ASC").
		Order("id ASC").
		Find(&entries).Error
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		result[entry.ClaimId] = append(result[entry.ClaimId], entry)
	}
	return result, nil
}

// GetClaimLinksForReceipts returns, for each of the given receipts that is linked
// to a claim, which claim that is. A link whose receipt has left the claim's group
// is included; callers that care compare ClaimGroupId with the receipt's group.
func (repository ClaimRepository) GetClaimLinksForReceipts(receiptIds []uint) ([]ReceiptClaimLink, error) {
	links := []ReceiptClaimLink{}
	if len(receiptIds) == 0 {
		return links, nil
	}

	err := repository.GetDB().
		Table("claim_receipts").
		Select("claim_receipts.receipt_id AS receipt_id, claims.id AS claim_id, claims.name AS claim_name, claims.group_id AS claim_group_id, claims.created_by AS claim_created_by").
		Joins("JOIN claims ON claims.id = claim_receipts.claim_id").
		Where("claim_receipts.receipt_id IN ?", receiptIds).
		Scan(&links).Error
	return links, err
}

// GetClaimLinksForGroupReceipts is GetClaimLinksForReceipts restricted to links
// whose receipt still lives in the claim's group: the links a claim actually
// shows, and so the only ones that may block deleting a receipt.
func (repository ClaimRepository) GetClaimLinksForGroupReceipts(receiptIds []uint) ([]ReceiptClaimLink, error) {
	links := []ReceiptClaimLink{}
	if len(receiptIds) == 0 {
		return links, nil
	}

	err := repository.GetDB().
		Table("claim_receipts").
		Select("claim_receipts.receipt_id AS receipt_id, claims.id AS claim_id, claims.name AS claim_name, claims.group_id AS claim_group_id, claims.created_by AS claim_created_by").
		Joins("JOIN claims ON claims.id = claim_receipts.claim_id").
		Joins("JOIN receipts ON receipts.id = claim_receipts.receipt_id").
		Where("claim_receipts.receipt_id IN ? AND receipts.group_id = claims.group_id", receiptIds).
		Scan(&links).Error
	return links, err
}

// LinkReceipts links each receipt to the claim, replacing any link it already
// has. Callers decide which receipts may be linked; the replace only ever applies
// to a stale link (one whose receipt has left that claim's group), because the
// service refuses a receipt that is live in another claim.
func (repository ClaimRepository) LinkReceipts(claimId uint, receiptIds []uint) error {
	if len(receiptIds) == 0 {
		return nil
	}

	return repository.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("receipt_id IN ?", receiptIds).Delete(&models.ClaimReceipt{}).Error; err != nil {
			return err
		}

		links := make([]models.ClaimReceipt, 0, len(receiptIds))
		for _, receiptId := range receiptIds {
			links = append(links, models.ClaimReceipt{ReceiptId: receiptId, ClaimId: claimId})
		}
		return tx.Omit(clause.Associations).Create(&links).Error
	})
}

// UnlinkReceipts removes the given receipts from the claim. A receipt linked to a
// different claim is untouched.
func (repository ClaimRepository) UnlinkReceipts(claimId uint, receiptIds []uint) error {
	if len(receiptIds) == 0 {
		return nil
	}

	return repository.GetDB().
		Where("claim_id = ? AND receipt_id IN ?", claimId, receiptIds).
		Delete(&models.ClaimReceipt{}).Error
}

// DetachReceiptFromClaims unlinks a receipt from whatever claim holds it. Called
// by ReceiptRepository.UpdateReceipt, inside its transaction, when a receipt moves
// to another group: a claim only holds receipts of its own group.
func DetachReceiptFromClaims(tx *gorm.DB, receiptId uint) error {
	return tx.Where("receipt_id = ?", receiptId).Delete(&models.ClaimReceipt{}).Error
}

func (repository ClaimRepository) GetLedgerEntryById(entryId uint) (models.ClaimLedgerEntry, error) {
	var entry models.ClaimLedgerEntry
	err := repository.GetDB().Where("id = ?", entryId).First(&entry).Error
	return entry, err
}

func (repository ClaimRepository) CreateLedgerEntry(entry models.ClaimLedgerEntry) (models.ClaimLedgerEntry, error) {
	err := repository.GetDB().Omit(clause.Associations).Create(&entry).Error
	return entry, err
}

// UpdateLedgerEntry replaces an entry's own columns, with the map form so a
// cleared amount or note is written.
func (repository ClaimRepository) UpdateLedgerEntry(entryId uint, entry models.ClaimLedgerEntry) error {
	return repository.GetDB().Model(&models.ClaimLedgerEntry{}).Where("id = ?", entryId).Updates(map[string]interface{}{
		"entry_type": entry.EntryType,
		"date":       entry.Date,
		"amount":     nullableValue(entry.Amount),
		"note":       entry.Note,
	}).Error
}

func (repository ClaimRepository) DeleteLedgerEntry(entryId uint) error {
	return repository.GetDB().Where("id = ?", entryId).Delete(&models.ClaimLedgerEntry{}).Error
}

// nullableValue turns a nil pointer into an untyped nil, so an Updates map writes
// NULL, and a set one into its value.
func nullableValue[T any](value *T) interface{} {
	if value == nil {
		return nil
	}
	return *value
}
