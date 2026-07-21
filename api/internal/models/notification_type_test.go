package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestNotificationType_Value(t *testing.T) {
	// NotificationType.Value does not validate; every value is returned as-is.
	cases := []NotificationType{NOTIFICATION_TYPE_NORMAL, NOTIFICATION_TYPE_URGENT, "", "anything"}
	for _, v := range cases {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestNotificationType_Scan(t *testing.T) {
	var notificationType NotificationType
	err := notificationType.Scan("NORMAL")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if notificationType != NOTIFICATION_TYPE_NORMAL {
		utils.PrintTestError(t, notificationType, NOTIFICATION_TYPE_NORMAL)
	}
}
