package commands

import (
	"fmt"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpsertGroupCommand_Validate_ValidInputs(t *testing.T) {
	tests := map[string]struct {
		command  UpsertGroupCommand
		isCreate bool
	}{
		"valid create without group members": {
			command: UpsertGroupCommand{
				Name:   "Test Group",
				Status: models.GROUP_ACTIVE,
			},
			isCreate: true,
		},
		"valid create with group members": {
			command: UpsertGroupCommand{
				Name:         "Test Group",
				Status:       models.GROUP_ACTIVE,
				GroupMembers: []UpsertGroupMemberCommand{{UserID: 1, GroupID: 1}},
			},
			isCreate: true,
		},
		"valid update with group members": {
			command: UpsertGroupCommand{
				Name:         "Test Group",
				Status:       models.GROUP_ACTIVE,
				GroupMembers: []UpsertGroupMemberCommand{{UserID: 1, GroupID: 1}},
			},
			isCreate: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate(test.isCreate)

			if len(vErr.Errors) > 0 {
				utils.PrintTestError(t, len(vErr.Errors), 0)
			}
		})
	}
}

func TestUpsertGroupCommand_Validate_InvalidInputs(t *testing.T) {
	tests := map[string]struct {
		command       UpsertGroupCommand
		isCreate      bool
		expectedError string
	}{
		"missing name": {
			command: UpsertGroupCommand{
				Status:       models.GROUP_ACTIVE,
				GroupMembers: []UpsertGroupMemberCommand{{UserID: 1, GroupID: 1}},
			},
			isCreate:      false,
			expectedError: "name",
		},
		"missing status": {
			command: UpsertGroupCommand{
				Name:         "Test Group",
				GroupMembers: []UpsertGroupMemberCommand{{UserID: 1, GroupID: 1}},
			},
			isCreate:      false,
			expectedError: "status",
		},
		"missing group members on update": {
			command: UpsertGroupCommand{
				Name:   "Test Group",
				Status: models.GROUP_ACTIVE,
			},
			isCreate:      false,
			expectedError: "groupMembers",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate(test.isCreate)

			if len(vErr.Errors) == 0 {
				utils.PrintTestError(t, len(vErr.Errors), "greater than 0")
			}

			if _, exists := vErr.Errors[test.expectedError]; !exists {
				utils.PrintTestError(t, "error should exist for field", test.expectedError)
			}
		})
	}
}

func TestUpsertGroupCommand_Validate_GroupMembersNotRequiredOnCreate(t *testing.T) {
	command := UpsertGroupCommand{
		Name:   "Test Group",
		Status: models.GROUP_ACTIVE,
	}

	vErr := command.Validate(true)

	if len(vErr.Errors) != 0 {
		utils.PrintTestError(t, len(vErr.Errors), 0)
	}
}

func TestUpsertGroupCommand_Validate_MultipleErrors(t *testing.T) {
	command := UpsertGroupCommand{}

	vErr := command.Validate(false)

	if len(vErr.Errors) != 3 {
		utils.PrintTestError(t, len(vErr.Errors), 3)
	}

	if _, exists := vErr.Errors["name"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "name")
	}

	if _, exists := vErr.Errors["status"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "status")
	}

	if _, exists := vErr.Errors["groupMembers"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "groupMembers")
	}
}

// A group name is used to build the group's on-disk storage directory, so a name
// containing path separators or traversal elements must be rejected (CWE-22).
func TestUpsertGroupCommand_Validate_RejectsUnsafeNames(t *testing.T) {
	unsafeNames := []string{
		"../../../../tmp/x",
		"..",
		".",
		"foo/bar",
		"foo\\bar",
		"a/../b",
		"/etc/passwd",
	}

	for _, name := range unsafeNames {
		for _, isCreate := range []bool{true, false} {
			t.Run(fmt.Sprintf("%q_isCreate_%v", name, isCreate), func(t *testing.T) {
				command := UpsertGroupCommand{
					Name:         name,
					Status:       models.GROUP_ACTIVE,
					GroupMembers: []UpsertGroupMemberCommand{{UserID: 1, GroupID: 1}},
				}

				vErr := command.Validate(isCreate)

				if _, exists := vErr.Errors["name"]; !exists {
					utils.PrintTestError(t, "no name error for unsafe name "+name, "name error")
				}
			})
		}
	}
}

func TestUpsertGroupCommand_Validate_AllowsNormalNames(t *testing.T) {
	names := []string{"My Receipts", "Reporting Load Test", "Mom & Dad's", "group-123"}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			command := UpsertGroupCommand{
				Name:   name,
				Status: models.GROUP_ACTIVE,
			}

			vErr := command.Validate(true)

			if _, exists := vErr.Errors["name"]; exists {
				utils.PrintTestError(t, "unexpected name error for "+name, "no name error")
			}
		})
	}
}
