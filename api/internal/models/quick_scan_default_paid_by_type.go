package models

import (
	"database/sql/driver"
	"errors"
)

type QuickScanDefaultPaidByType string

const (
	QUICK_SCAN_PAID_BY_UPLOADER QuickScanDefaultPaidByType = "UPLOADER"
	QUICK_SCAN_PAID_BY_USER     QuickScanDefaultPaidByType = "USER"
)

func (self *QuickScanDefaultPaidByType) Scan(value string) error {
	*self = QuickScanDefaultPaidByType(value)
	return nil
}

func (self QuickScanDefaultPaidByType) Value() (driver.Value, error) {
	if self != QUICK_SCAN_PAID_BY_UPLOADER && self != QUICK_SCAN_PAID_BY_USER && self != "" {
		return nil, errors.New("invalid quickScanDefaultPaidByType")
	}
	return string(self), nil
}

func QuickScanDefaultPaidByTypes() []interface{} {
	return []interface{}{QUICK_SCAN_PAID_BY_UPLOADER, QUICK_SCAN_PAID_BY_USER}
}
