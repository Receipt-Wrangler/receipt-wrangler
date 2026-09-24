package models

import (
	"database/sql/driver"
	"errors"
)

// ReceiptSummaryPosition is where a group's receipt summary renders relative to its
// receipts list: above it or below it. Presentation only -- it changes nothing about
// which receipts the block describes.
//
// It is a per-group setting rather than a per-user preference because the rest of the
// summary configuration is (see GroupReceiptSettings), and a block whose contents one
// admin controls but whose placement each viewer controls would be two settings with
// one name.
type ReceiptSummaryPosition string

const (
	RECEIPT_SUMMARY_POSITION_TOP    ReceiptSummaryPosition = "TOP"
	RECEIPT_SUMMARY_POSITION_BOTTOM ReceiptSummaryPosition = "BOTTOM"
)

func (receiptSummaryPosition *ReceiptSummaryPosition) Scan(value string) error {
	*receiptSummaryPosition = ReceiptSummaryPosition(value)
	return nil
}

func (receiptSummaryPosition ReceiptSummaryPosition) Value() (driver.Value, error) {
	if len(receiptSummaryPosition) == 0 {
		return "", nil
	}

	if receiptSummaryPosition != RECEIPT_SUMMARY_POSITION_TOP &&
		receiptSummaryPosition != RECEIPT_SUMMARY_POSITION_BOTTOM {
		return nil, errors.New("invalid receipt summary position")
	}
	return string(receiptSummaryPosition), nil
}

// OrDefault resolves the empty value to BOTTOM, which is where the summary rendered
// before this setting existed.
//
// Empty is reachable even though the column carries a GORM default: the summary
// service maps a missing settings row to a ZERO GroupReceiptSettings, and the
// projection loader hydrates settings that predate the column on an install whose
// AutoMigrate has not run yet. Both of those would otherwise emit "" on the wire,
// where a closed Dart EnumClass throws and fails the WHOLE payload -- the mechanism
// behind the two documented login outages. Normalizing here is strictly stronger than
// teaching every client to tolerate a value the server should never send.
func (receiptSummaryPosition ReceiptSummaryPosition) OrDefault() ReceiptSummaryPosition {
	if receiptSummaryPosition != RECEIPT_SUMMARY_POSITION_TOP {
		return RECEIPT_SUMMARY_POSITION_BOTTOM
	}
	return RECEIPT_SUMMARY_POSITION_TOP
}

func ReceiptSummaryPositions() []interface{} {
	return []interface{}{RECEIPT_SUMMARY_POSITION_TOP, RECEIPT_SUMMARY_POSITION_BOTTOM}
}
