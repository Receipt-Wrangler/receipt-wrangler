package commands

func GetDefaultAdminSignUpCommand() SignUpCommand {
	return SignUpCommand{
		Username:    "admin",
		DisplayName: "Admin",
		Password:    "admin",
		IsDummyUser: false,
	}
}
