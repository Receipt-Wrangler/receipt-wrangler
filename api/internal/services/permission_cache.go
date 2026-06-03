package services

import (
	"sync"

	"receipt-wrangler/api/internal/permissions"
)

// rolePermissionCacheEntry keys the cache by scope + role id. App and group role
// ids overlap (separate tables), so the scope is part of the key.
type rolePermissionCacheEntry struct {
	scope  permissions.Scope
	roleId uint
}

// rolePermissionCache memoizes the permission strings granted by a role. Only
// the (rarely changing) permission list of a role is cached here — a user's
// role *assignment* is resolved fresh on every check, so re-assigning a user to
// a different role always takes effect immediately. The cache is invalidated
// whenever a role's permissions are updated or the role is deleted.
var (
	rolePermissionCacheMutex sync.RWMutex
	rolePermissionCache      = map[rolePermissionCacheEntry][]string{}
)

func getCachedRolePermissions(scope permissions.Scope, roleId uint) ([]string, bool) {
	rolePermissionCacheMutex.RLock()
	defer rolePermissionCacheMutex.RUnlock()

	perms, ok := rolePermissionCache[rolePermissionCacheEntry{scope: scope, roleId: roleId}]
	return perms, ok
}

func setCachedRolePermissions(scope permissions.Scope, roleId uint, perms []string) {
	rolePermissionCacheMutex.Lock()
	defer rolePermissionCacheMutex.Unlock()

	rolePermissionCache[rolePermissionCacheEntry{scope: scope, roleId: roleId}] = perms
}

// clearRolePermissionCache evicts a single role's cached permissions. Call this
// after a role's permissions change or the role is deleted.
func clearRolePermissionCache(scope permissions.Scope, roleId uint) {
	rolePermissionCacheMutex.Lock()
	defer rolePermissionCacheMutex.Unlock()

	delete(rolePermissionCache, rolePermissionCacheEntry{scope: scope, roleId: roleId})
}

// clearRolePermissionCacheAll empties the cache entirely. Used by tests to keep
// cases independent.
func clearRolePermissionCacheAll() {
	rolePermissionCacheMutex.Lock()
	defer rolePermissionCacheMutex.Unlock()

	rolePermissionCache = map[rolePermissionCacheEntry][]string{}
}
