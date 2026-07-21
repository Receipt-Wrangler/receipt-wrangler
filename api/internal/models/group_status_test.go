package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestGroupStatus_Value(t *testing.T) {
	// GroupStatus.Value does not validate; every value is returned as-is.
	cases := []GroupStatus{GROUP_ACTIVE, GROUP_ARCHIVED, "", "anything"}
	for _, v := range cases {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestGroupStatus_Scan(t *testing.T) {
	var status GroupStatus
	err := status.Scan("ACTIVE")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if status != GROUP_ACTIVE {
		utils.PrintTestError(t, status, GROUP_ACTIVE)
	}
}
