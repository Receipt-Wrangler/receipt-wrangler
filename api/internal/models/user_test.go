package models

import (
	"encoding/json"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUser_MarshalDefaultAvatarColor(t *testing.T) {
	user := User{DefaultAvatarColor: "#27b1ff"}

	bytes, err := json.Marshal(user)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	var decoded map[string]interface{}
	if unmarshalErr := json.Unmarshal(bytes, &decoded); unmarshalErr != nil {
		utils.PrintTestError(t, unmarshalErr, nil)
	}

	// The field must serialize under the clean "defaultAvatarColor" key.
	if decoded["defaultAvatarColor"] != "#27b1ff" {
		utils.PrintTestError(t, decoded["defaultAvatarColor"], "#27b1ff")
	}

	// The previously malformed key must not exist.
	if _, ok := decoded["defaultAvatarColor gorm:"]; ok {
		utils.PrintTestError(t, "malformed json key present", "absent")
	}
}

func TestUser_UnmarshalDefaultAvatarColor(t *testing.T) {
	var user User
	err := json.Unmarshal([]byte(`{"defaultAvatarColor": "#abcdef"}`), &user)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if user.DefaultAvatarColor != "#abcdef" {
		utils.PrintTestError(t, user.DefaultAvatarColor, "#abcdef")
	}
}
