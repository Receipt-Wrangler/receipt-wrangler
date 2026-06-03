package services

import (
	"errors"
	"fmt"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"

	"gorm.io/gorm"
)

// ErrNoRequiredPermissions is returned when a permission check is invoked
// without specifying any required permission.
var ErrNoRequiredPermissions = errors.New("at least one required permission must be provided")

// matcher is the matching strategy used by a check: permissions.HasAll (AND,
// the default) or permissions.HasAny (OR).
type matcher func(granted []string, required ...string) bool

// PermissionService resolves a user's effective permissions from the database
// at check time (never trusting JWT contents for authorization) and matches
// them against the required permissions. App-level and group-level checks are
// kept as separate entry points; both share the permissions package matcher.
type PermissionService struct {
	BaseService
}

func NewPermissionService(tx *gorm.DB) PermissionService {
	service := PermissionService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

// HasAppPermissions reports whether the user's app role grants ALL of the
// required app-level permissions (logical AND — the default). A single-permission
// check is the common case: HasAppPermissions(userId, permissions.AppUsersRead).
func (service PermissionService) HasAppPermissions(userId uint, required ...string) (bool, error) {
	return service.checkApp(userId, permissions.HasAll, required...)
}

// HasAnyAppPermission reports whether the user's app role grants AT LEAST ONE of
// the required app-level permissions (logical OR).
func (service PermissionService) HasAnyAppPermission(userId uint, required ...string) (bool, error) {
	return service.checkApp(userId, permissions.HasAny, required...)
}

// HasGroupPermissions reports whether the user's role in the group grants ALL of
// the required group-level permissions (logical AND — the default).
func (service PermissionService) HasGroupPermissions(userId uint, groupId uint, required ...string) (bool, error) {
	return service.checkGroup(userId, groupId, permissions.HasAll, required...)
}

// HasAnyGroupPermission reports whether the user's role in the group grants AT
// LEAST ONE of the required group-level permissions (logical OR).
func (service PermissionService) HasAnyGroupPermission(userId uint, groupId uint, required ...string) (bool, error) {
	return service.checkGroup(userId, groupId, permissions.HasAny, required...)
}

func (service PermissionService) checkApp(userId uint, match matcher, required ...string) (bool, error) {
	if err := validateRequiredPermissions(permissions.ScopeApp, required); err != nil {
		return false, err
	}

	granted, err := service.resolveAppPermissions(userId)
	if err != nil {
		return false, err
	}

	return match(granted, required...), nil
}

func (service PermissionService) checkGroup(userId uint, groupId uint, match matcher, required ...string) (bool, error) {
	if err := validateRequiredPermissions(permissions.ScopeGroup, required); err != nil {
		return false, err
	}

	granted, err := service.resolveGroupPermissions(userId, groupId)
	if err != nil {
		return false, err
	}

	return match(granted, required...), nil
}

// validateRequiredPermissions guards against programmer error: an empty
// requirement, an unknown permission key (typo), or a key whose scope does not
// match the check being performed (e.g. an app permission passed to a group
// check).
func validateRequiredPermissions(scope permissions.Scope, required []string) error {
	if len(required) == 0 {
		return ErrNoRequiredPermissions
	}

	for _, key := range required {
		descriptor, ok := permissions.Get(key)
		if !ok {
			return fmt.Errorf("unknown permission %q", key)
		}
		if descriptor.Scope != scope {
			return fmt.Errorf("permission %q is %s-scoped, expected %s", key, descriptor.Scope, scope)
		}
	}

	return nil
}

// resolveAppPermissions loads the user's current app role permissions. A user
// with no app role assigned resolves to no permissions (deny).
func (service PermissionService) resolveAppPermissions(userId uint) ([]string, error) {
	roleRepository := repositories.NewRoleRepository(service.TX)

	appRoleId, err := roleRepository.GetUserAppRoleId(userId)
	if err != nil {
		return nil, err
	}
	if appRoleId == nil {
		return []string{}, nil
	}

	return loadRolePermissions(roleRepository, permissions.ScopeApp, *appRoleId)
}

// resolveGroupPermissions loads the user's current group role permissions for a
// group. A user who is not a member, or whose membership has no group role,
// resolves to no permissions (deny).
func (service PermissionService) resolveGroupPermissions(userId uint, groupId uint) ([]string, error) {
	roleRepository := repositories.NewRoleRepository(service.TX)

	groupRoleId, err := roleRepository.GetGroupMemberRoleId(userId, groupId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	if groupRoleId == nil {
		return []string{}, nil
	}

	return loadRolePermissions(roleRepository, permissions.ScopeGroup, *groupRoleId)
}

// loadRolePermissions returns a role's permission strings, consulting the cache
// first and populating it on a miss. Callers must not mutate the returned slice.
func loadRolePermissions(roleRepository repositories.RoleRepository, scope permissions.Scope, roleId uint) ([]string, error) {
	if cached, ok := getCachedRolePermissions(scope, roleId); ok {
		return cached, nil
	}

	var perms []string
	var err error
	if scope == permissions.ScopeApp {
		perms, err = roleRepository.GetAppRolePermissions(roleId)
	} else {
		perms, err = roleRepository.GetGroupRolePermissions(roleId)
	}
	if err != nil {
		return nil, err
	}

	setCachedRolePermissions(scope, roleId, perms)
	return perms, nil
}
