package models

import (
	"encoding/json"
	"receipt-wrangler/api/internal/utils"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestReceipt_ToString(t *testing.T) {
	receipt := Receipt{
		BaseModel: BaseModel{ID: 1},
		Name:      "Test Receipt",
		Amount:    decimal.NewFromFloat(12.34),
		Date:      time.Now(),
		Status:    OPEN,
		GroupId:   2,
	}

	result, err := receipt.ToString()
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	// The result must be valid JSON that round-trips the key fields.
	var decoded map[string]interface{}
	if unmarshalErr := json.Unmarshal([]byte(result), &decoded); unmarshalErr != nil {
		utils.PrintTestError(t, unmarshalErr, nil)
	}

	if decoded["name"] != "Test Receipt" {
		utils.PrintTestError(t, decoded["name"], "Test Receipt")
	}
}
