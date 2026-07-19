package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestSystemTaskStatus_Value(t *testing.T) {
	valid := []SystemTaskStatus{SYSTEM_TASK_SUCCEEDED, SYSTEM_TASK_FAILED}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestSystemTaskStatus_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "empty", SystemTaskStatus(""))
	assertValuerInvalid(t, "bogus", SystemTaskStatus("bogus"))
}

func TestSystemTaskStatus_Scan(t *testing.T) {
	var status SystemTaskStatus
	err := status.Scan("SUCCEEDED")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if status != SYSTEM_TASK_SUCCEEDED {
		utils.PrintTestError(t, status, SYSTEM_TASK_SUCCEEDED)
	}
}

func TestSystemTaskType_Value(t *testing.T) {
	valid := []SystemTaskType{
		RECEIPT_UPLOADED,
		OCR_PROCESSING,
		CHAT_COMPLETION,
		MAGIC_FILL,
		QUICK_SCAN,
		EMAIL_UPLOAD,
		EMAIL_READ,
		SYSTEM_EMAIL_CONNECTIVITY_CHECK,
		RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK,
		PROMPT_GENERATED,
		RECEIPT_UPDATED,
		API_KEY_DELETED,
		HTML_TO_PDF,
	}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestSystemTaskType_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "empty", SystemTaskType(""))
	assertValuerInvalid(t, "bogus", SystemTaskType("bogus"))
}

func TestSystemTaskType_Scan(t *testing.T) {
	var taskType SystemTaskType
	err := taskType.Scan("RECEIPT_UPLOADED")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if taskType != RECEIPT_UPLOADED {
		utils.PrintTestError(t, taskType, RECEIPT_UPLOADED)
	}
}

func TestAssociatedEntityType_Value(t *testing.T) {
	valid := []AssociatedEntityType{
		RECEIPT,
		SYSTEM_EMAIL,
		PROMPT,
		RECEIPT_PROCESSING_SETTINGS,
		NOOP_ENTITY_TYPE,
		API_KEY,
	}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestAssociatedEntityType_Value_Invalid(t *testing.T) {
	assertValuerInvalid(t, "empty", AssociatedEntityType(""))
	assertValuerInvalid(t, "bogus", AssociatedEntityType("bogus"))
}

func TestAssociatedEntityType_Scan(t *testing.T) {
	var entityType AssociatedEntityType
	err := entityType.Scan("RECEIPT")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if entityType != RECEIPT {
		utils.PrintTestError(t, entityType, RECEIPT)
	}
}
