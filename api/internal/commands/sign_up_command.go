package commands

type SignUpCommand struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	IsDummyUser bool   `json:"isDummyUser"`
	AppRoleID   *uint  `json:"appRoleId"`
}
