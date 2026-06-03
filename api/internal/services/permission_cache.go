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
//
// rolePermissionCacheGeneration is bumped on every eviction. A cache write
// carries the generation observed before its DB read (see loadRolePermissions)
// and is dropped if the generation has since advanced — this prevents an
// in-flight cache miss from resurrecting stale permissions after a concurrent
// eviction, which would otherwise leave revoked permissions authorized.
var (
	rolePermissionCacheMutex      sync.RWMutex
	rolePermissionCache           = map[rolePermissionCacheEntry][]string{}
	rolePermissionCacheGeneration uint64
)

func getCachedRolePermissions(scope permissions.Scope, roleId uint) ([]string, bool) {
	rolePermissionCacheMutex.RLock()
	defer rolePermissionCacheMutex.RUnlock()

	perms, ok := rolePermissionCache[rolePermissionCacheEntry{scope: scope, roleId: roleId}]
	return perms, ok
}

// rolePermissionCacheGen returns the current eviction generation. Capture it
// before reading permissions from the database and pass it to
// setCachedRolePermissions so a concurrent eviction invalidates the write.
func rolePermissionCacheGen() uint64 {
	rolePermissionCacheMutex.RLock()
	defer rolePermissionCacheMutex.RUnlock()

	return rolePermissionCacheGeneration
}

// setCachedRolePermissions stores perms only if no eviction has happened since
// observedGen was captured; otherwise the value is considered potentially stale
// and discarded (the next check re-reads from the database).
func setCachedRolePermissions(scope permissions.Scope, roleId uint, perms []string, observedGen uint64) {
	rolePermissionCacheMutex.Lock()
	defer rolePermissionCacheMutex.Unlock()

	if rolePermissionCacheGeneration != observedGen {
		return
	}

	rolePermissionCache[rolePermissionCacheEntry{scope: scope, roleId: roleId}] = perms
}

// clearRolePermissionCache evicts a single role's cached permissions and bumps
// the generation. Call this after a role's permissions change or the role is
// deleted.
func clearRolePermissionCache(scope permissions.Scope, roleId uint) {
	rolePermissionCacheMutex.Lock()
	defer rolePermissionCacheMutex.Unlock()

	rolePermissionCacheGeneration++
	delete(rolePermissionCache, rolePermissionCacheEntry{scope: scope, roleId: roleId})
}

// clearRolePermissionCacheAll empties the cache entirely and bumps the
// generation. Used by tests to keep cases independent.
func clearRolePermissionCacheAll() {
	rolePermissionCacheMutex.Lock()
	defer rolePermissionCacheMutex.Unlock()

	rolePermissionCacheGeneration++
	rolePermissionCache = map[rolePermissionCacheEntry][]string{}
}
