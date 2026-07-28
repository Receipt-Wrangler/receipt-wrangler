package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

// TestAuthorizeGroupMemberChanges exercises the two-layer guard behind UpdateGroup
// (GHSA-89mm-9qfv-cjg3): the group.members.* CRUD gate and the privilege ceiling
// that keeps a caller from granting or stripping a role beyond their own.
func TestAuthorizeGroupMemberChanges(t *testing.T) {
	defer repositories.TruncateTestDb()
	ClearRolePermissionCacheForTests()

	db := repositories.GetDB()
	roleRepository := repositories.NewRoleRepository(nil)

	group := models.Group{Name: "authz-group"}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}

	mustRole := func(name string, perms []string) uint {
		role, err := roleRepository.CreateGroupRole(name, "", perms, nil, nil, nil, false, false)
		if err != nil {
			t.Fatalf("seed role %q: %v", name, err)
		}
		return role.ID
	}

	ownerRole := mustRole("Owner", permissions.LegacyGroupOwnerKeys())
	// Manager holds the full member-management set but a limited overall scope, so
	// it can manage members yet must not be able to hand out the owner role.
	managerRole := mustRole("Manager", []string{
		permissions.GroupUpdate, permissions.GroupView, permissions.GroupReceiptsRead,
		permissions.GroupMembersCreate, permissions.GroupMembersUpdate, permissions.GroupMembersDelete,
	})
	viewerRole := mustRole("Viewer", []string{permissions.GroupView, permissions.GroupReceiptsRead})
	// Low holds group.update (enough to reach the endpoint) but no member perms.
	lowRole := mustRole("Low", []string{permissions.GroupUpdate, permissions.GroupView, permissions.GroupReceiptsRead})

	const owner, manager, low = uint(1), uint(2), uint(3)
	seed := func(userId uint, roleId uint) {
		if err := db.Create(&models.GroupMember{GroupID: group.ID, UserID: userId, GroupRoleID: &roleId}).Error; err != nil {
			t.Fatalf("seed member %d: %v", userId, err)
		}
	}
	seed(owner, ownerRole)
	seed(manager, managerRole)
	seed(low, lowRole)

	member := func(userId, roleId uint) commands.UpsertGroupMemberCommand {
		id := roleId
		return commands.UpsertGroupMemberCommand{UserID: userId, GroupID: group.ID, GroupRoleID: &id}
	}
	memberNilRole := func(userId uint) commands.UpsertGroupMemberCommand {
		return commands.UpsertGroupMemberCommand{UserID: userId, GroupID: group.ID}
	}

	current := []commands.UpsertGroupMemberCommand{member(owner, ownerRole), member(manager, managerRole), member(low, lowRole)}

	service := NewGroupService(nil)

	tests := []struct {
		name      string
		caller    uint
		submitted []commands.UpsertGroupMemberCommand
		forbidden bool
	}{
		{
			name:      "manager cannot promote self to owner (escalation) despite members.update",
			caller:    manager,
			submitted: []commands.UpsertGroupMemberCommand{member(owner, ownerRole), member(manager, ownerRole), member(low, lowRole)},
			forbidden: true,
		},
		{
			name:      "manager can add a member at or below their own privilege",
			caller:    manager,
			submitted: append(append([]commands.UpsertGroupMemberCommand{}, current...), member(4, viewerRole)),
			forbidden: false,
		},
		{
			name:      "manager cannot add a member with a role above their own",
			caller:    manager,
			submitted: append(append([]commands.UpsertGroupMemberCommand{}, current...), member(4, ownerRole)),
			forbidden: true,
		},
		{
			name:      "manager can add a member with no role",
			caller:    manager,
			submitted: append(append([]commands.UpsertGroupMemberCommand{}, current...), memberNilRole(4)),
			forbidden: false,
		},
		{
			name:      "manager can remove a member at or below their own privilege",
			caller:    manager,
			submitted: []commands.UpsertGroupMemberCommand{member(owner, ownerRole), member(manager, managerRole)},
			forbidden: false,
		},
		{
			name:      "manager cannot remove the owner",
			caller:    manager,
			submitted: []commands.UpsertGroupMemberCommand{member(manager, managerRole), member(low, lowRole)},
			forbidden: true,
		},
		{
			name:      "group.update-only member cannot add a member (missing members.create)",
			caller:    low,
			submitted: append(append([]commands.UpsertGroupMemberCommand{}, current...), member(4, viewerRole)),
			forbidden: true,
		},
		{
			name:      "group.update-only member may submit the roster unchanged",
			caller:    low,
			submitted: current,
			forbidden: false,
		},
		{
			name:      "owner may reassign any member's role",
			caller:    owner,
			submitted: []commands.UpsertGroupMemberCommand{member(owner, ownerRole), member(manager, managerRole), member(low, managerRole)},
			forbidden: false,
		},
		{
			// A duplicate userId must not let an escalating entry hide behind a
			// "safe" one that would win a last-wins dedup: the repository persists
			// the raw slice, so the whole roster is rejected.
			name:   "duplicate userId entry is rejected (cannot smuggle an owner role)",
			caller: manager,
			submitted: []commands.UpsertGroupMemberCommand{
				member(owner, ownerRole),
				member(manager, ownerRole),   // escalating entry
				member(manager, managerRole), // "safe" entry a last-wins map would keep
				member(low, lowRole),
			},
			forbidden: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.AuthorizeGroupMemberChanges(test.caller, group.ID, test.submitted)
			if test.forbidden {
				if !errors.Is(err, ErrGroupMemberChangeForbidden) {
					t.Fatalf("expected ErrGroupMemberChangeForbidden, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected the change to be authorized, got %v", err)
			}
		})
	}
}
