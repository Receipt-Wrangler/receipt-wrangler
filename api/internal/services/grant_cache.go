package services

import "sync"

// grantEntry holds the category/tag ids a group role grants its members, as sets
// for O(1) membership tests. An EMPTY set means that resource is unrestricted
// (the role grants every category/tag) — so categories and tags are independent:
// a role may restrict categories while leaving tags unrestricted, or vice versa.
// Callers must treat the maps as read-only (they are shared cache state).
type grantEntry struct {
	categoryIds map[uint]struct{}
	tagIds      map[uint]struct{}
}

// groupRoleGrantCache memoizes a group role's category/tag grant id sets, keyed
// by group-role id (grants only exist on group roles). Like the permission
// cache, only a role's grant *lists* are cached — a user's role *assignment* is
// resolved fresh on every check, so re-assigning a user takes effect
// immediately. Invalidated whenever a group role is updated or deleted.
//
// groupRoleGrantCacheGeneration is bumped on every eviction; a cache write
// carries the generation observed before its DB read and is dropped if the
// generation has since advanced, preventing an in-flight miss from resurrecting
// stale grants after a concurrent eviction (which would otherwise leave a
// just-removed restriction — or a just-added one — applied incorrectly).
var (
	groupRoleGrantCacheMutex      sync.RWMutex
	groupRoleGrantCache           = map[uint]*grantEntry{}
	groupRoleGrantCacheGeneration uint64
)

func getCachedGroupRoleGrants(roleId uint) (*grantEntry, bool) {
	groupRoleGrantCacheMutex.RLock()
	defer groupRoleGrantCacheMutex.RUnlock()

	entry, ok := groupRoleGrantCache[roleId]
	return entry, ok
}

// groupRoleGrantCacheGen returns the current eviction generation. Capture it
// before reading grants from the database and pass it to setCachedGroupRoleGrants
// so a concurrent eviction invalidates the write.
func groupRoleGrantCacheGen() uint64 {
	groupRoleGrantCacheMutex.RLock()
	defer groupRoleGrantCacheMutex.RUnlock()

	return groupRoleGrantCacheGeneration
}

// setCachedGroupRoleGrants stores entry only if no eviction has happened since
// observedGen was captured; otherwise the value is considered potentially stale
// and discarded (the next check re-reads from the database).
func setCachedGroupRoleGrants(roleId uint, entry *grantEntry, observedGen uint64) {
	groupRoleGrantCacheMutex.Lock()
	defer groupRoleGrantCacheMutex.Unlock()

	if groupRoleGrantCacheGeneration != observedGen {
		return
	}

	groupRoleGrantCache[roleId] = entry
}

// clearGroupRoleGrantCache evicts a single group role's cached grants and bumps
// the generation. Call this after a group role's grants change or the role is
// deleted.
func clearGroupRoleGrantCache(roleId uint) {
	groupRoleGrantCacheMutex.Lock()
	defer groupRoleGrantCacheMutex.Unlock()

	groupRoleGrantCacheGeneration++
	delete(groupRoleGrantCache, roleId)
}

// clearGroupRoleGrantCacheAll empties the cache entirely and bumps the
// generation. Used by tests to keep cases independent.
func clearGroupRoleGrantCacheAll() {
	groupRoleGrantCacheMutex.Lock()
	defer groupRoleGrantCacheMutex.Unlock()

	groupRoleGrantCacheGeneration++
	groupRoleGrantCache = map[uint]*grantEntry{}
}

// ClearGroupRoleGrantCacheForTests empties the entire grant cache. Exported so
// tests in other packages can keep cases independent: a truncated test database
// reuses role ids across cases, which would otherwise return another case's
// cached grants.
func ClearGroupRoleGrantCacheForTests() {
	clearGroupRoleGrantCacheAll()
}
