package structs

import (
	"time"
)

type UserView struct {
	CreatedAt          time.Time `json:"createdAt"`
	DefaultAvatarColor string    `json:"defaultAvatarColor"`
	DisplayName        string    `json:"displayName"`
	ID                 uint      `json:"id"`
	IsDummyUser        bool      `json:"isDummyUser"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Username           string    `json:"username"`
	AppRoleID          *uint     `json:"appRoleId"`
}
