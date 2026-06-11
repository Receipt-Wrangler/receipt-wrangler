# Backend: Remove the Legacy Role System — Plan & Progress

> **Status: WORK IN PROGRESS — does not compile yet.** This branch
> (`feature/role-rework-backend-legacy-removal`) is a mid-refactor checkpoint pushed so the
> environment can be reconfigured. Phase 1 is complete; Phase 2 is partially applied (see the
> checklist). The tree will not build until Phase 2 is finished, because some legacy-enum readers
> have been removed while others (and the enum types themselves) still exist.

## Context

Receipt Wrangler historically gated access with two hardcoded enums — `models.UserRole`
(`ADMIN`/`USER`) and `models.GroupRole` (`OWNER`/`EDITOR`/`VIEWER`). It has since migrated to a
configurable, permission-based role system (permission registry + `AppRole`/`GroupRoleDefinition`
models + `PermissionService`), and handler-level authorization is already fully enforced through the
modern system (`HandleRequest` → `enforcePermissions`). The legacy enums now survive only as a
transitional bridge: a few stray `== ADMIN` checks, the JWT `Claims.userRole` UI hint, the
serialized `userRole`/`groupRole` API fields, and `DeriveLegacy*` shims that keep those fields
populated. `api/CLAUDE.md` explicitly says to "delete these derivations once the remaining
legacy-enum readers are migrated."

This task does exactly that for the **backend**: rip out the legacy enums and every reader, deriver,
and contract field, leaving only the modern permission system. Desktop is a planned follow-up phase.

### Scope decisions (confirmed with user)

- **Full rip-out including the API contract.** Remove `userRole`/`groupRole`/`Claims.userRole` from
  `swagger.yml` and regenerate clients. The desktop still consumes these fields (user-list role
  column, group-form owner gating, auth-state role selectors), so it will **not compile** until the
  follow-up frontend phase — that is accepted and out of scope here. Mobile is already clean.
- **"Admin" is defined by the app permission `app.users.read`** (`permissions.AppUsersRead`). This
  reproduces the old ADMIN/USER split exactly: the seeded "Legacy Admin" role grants it, "Legacy
  User" deliberately omits it, and it is already the gate on the admin user-listing endpoint.
- **KEEP the "modern legacy roles"** — the seeded `Legacy Admin/User/Viewer/Editor/Owner` system
  roles (`repositories/seed_roles.go`, `system_role_names.go`), the `Legacy*Keys()` permission sets
  (`permissions/legacy.go`), and the one-time data migration (`repositories/data_migrations.go`).
  These are the modern bridge for existing installs and stay.

---

## Progress checklist

### ✅ Phase 1 — Migrate live enforcement to permissions (DONE)

- [x] `repositories/roles.go` — added `appRoleIdsWithPermission(perm string) ([]uint, error)`
  (unexported, wildcard-aware via `permissions.HasAll`) and
  `CountUsersWithAppPermission(perm string) (int64, error)`.
- [x] `services/users.go` `DeleteUser` — last-admin guard now checks
  `HasAppPermissions(user.ID, AppUsersRead)` + `CountUsersWithAppPermission(AppUsersRead) <= 1`.
  Added `permissions` import.
- [x] `repositories/users.go` `IsFirstAdminToLogin` — reimplemented via
  `appRoleIdsWithPermission(AppUsersRead)` + `app_role_id IN (...) AND last_login_date IS NOT NULL`.
  Added `permissions` import.
- [x] `services/auth.go` `LoginUser` — admin gate now
  `NewPermissionService(nil).HasAppPermissions(dbUser.ID, AppUsersRead)`. Added `permissions` import.
- [x] Paged-groups ALL gate — removed `token.UserRole != models.ADMIN` from
  `commands/paged_group_request_command.go` `Validate`; moved into `handlers/groups.go`
  `GetPagedGroups` body as `HasAppPermissions(token.UserId, AppGroupsRead)` (403 otherwise),
  mirroring the `app.api-keys.read-any` precedent. Removed now-unused `models` import from the
  command.

### 🚧 Phase 2 — Remove enums, fields, derivers, contract (Go) (PARTIAL)

Done:
- [x] `commands/sign_up_command.go` — removed `UserRole` field + `models` import.
- [x] `commands/upsert_group_memeber_command.go` — removed `GroupRole` field + `models` import.
- [x] `commands/constants.go` — removed `UserRole: models.ADMIN` + `models` import.
- [x] `repositories/users.go` `CreateUser` — removed `UserRole` from struct literal; `AppRoleID==nil`
  branch now computes `isAdmin := usrCnt == 0` and calls `resolveAppRoleId(tx, isAdmin)`; `AppRoleID!=nil`
  branch just sets `user.AppRoleID` (dropped `DeriveLegacyUserRole`).
- [x] `repositories/users.go` `resolveAppRoleId` — signature changed to `(tx, isAdmin bool)`.
- [x] `repositories/users.go` `UpdateUser` — dropped `user_role` from updates map and the
  `DeriveLegacyUserRole` derivation.

**Remaining (NOT yet applied — this is why the tree doesn't compile):**
- [ ] `repositories/groups.go` `buildGroupMemberFromCommand` — drop `GroupRole: command.GroupRole`
  and the `DeriveLegacyGroupRole` derivation; keep `GroupRoleID = command.GroupRoleID`. The
  `roleRepository` param and the `error` return likely become removable — check callers (grep
  `buildGroupMemberFromCommand`) and the `memberRoleRepository` local in `CreateGroup`.
- [ ] `repositories/groups.go` `CreateGroup` (~line 147) — **drop only `GroupRole: models.OWNER`**
  and KEEP `GroupRoleID: defaultGroupRoleId` (`GetDefaultGroupRoleId()`). Do NOT hardcode Legacy
  Owner — the creator already gets the configurable default group role (Legacy Owner by default via
  `EnsureDefaultRoles`); forcing Legacy Owner would regress admin-chosen defaults.
- [ ] `repositories/roles.go` — delete `DeriveLegacyUserRole` and `DeriveLegacyGroupRole` (now
  uncalled once groups.go is updated).
- [ ] `handlers/sign_up.go` — strip only `userData.AppRoleID = nil` (remove the `UserRole = ""`
  line); update the comment.
- [ ] `middleware/user.go` `ValidateUserData` — accept only `AppRoleID` (remove legacy `UserRole`
  acceptance); keep the `models` import (uses `models.User`).
- [ ] `middleware/group.go` — delete the dead `ValidateGroupRole` middleware; drop now-unused
  `models`/`services`/`chi` imports.
- [ ] `services/groups.go` — delete the `ValidateGroupRole` method (its only caller was the deleted
  middleware); drop the now-unused `errors` import if nothing else uses it.
- [ ] `services/auth.go` `GenerateJWT` — drop `UserRole:` from both Claims constructors
  (access ~line 115 + refresh ~line 140).
- [ ] `services/api_key.go` — drop `UserRole:` from the API-key Claims constructor (~line 136).
- [ ] `structs/claims.go` — remove the `UserRole` field and the enum-validity block in `Validate()`;
  keep the `models` import (`ApiKeyScope models.ApiKeyScope`).
- [ ] `structs/api_user.go` — remove `UserView.UserRole`; drop the `models` import if now unused.
- [ ] `models/user.go` — remove `User.UserRole`. `models/group_member.go` — remove
  `GroupMember.GroupRole`.
- [ ] Delete `models/user_role.go`, `models/group_role.go`, `models/utils.go`
  (`BuildGroupMap`/`BuildUserRoleMap`/`HasRole` — all legacy-only).

After this, `go build ./...` should pass; only `_test.go` files should fail to compile.

### ⬜ Phase 3 — Guard the data migration against dropped columns (NOT STARTED)

- [ ] `repositories/data_migrations.go` `assignLegacyEquivalentRoles` — it reads the legacy columns
  via raw SQL (`Where("user_role = ?")` / `Where("group_role = ?")`), which survives removing the Go
  struct fields. (a) Replace the `models.UserRole`/`models.GroupRole`-typed mapping values with plain
  string literals (`"ADMIN"`, `"USER"`, `"OWNER"`, `"EDITOR"`, `"VIEWER"`). (b) Wrap each backfill
  loop in `if tx.Migrator().HasColumn(&models.User{}, "user_role") { ... }` /
  `HasColumn(&models.GroupMember{}, "group_role")`. Existing installs keep the physical columns (GORM
  never drops them) so the backfill runs; fresh installs never create them so the guard skips cleanly
  instead of erroring with "no such column". **Do NOT add a drop-column migration** — keep the
  columns for the upgrade path.

### ⬜ Phase 4 — Swagger + client regeneration (NOT STARTED)

- [ ] `api/swagger.yml` — remove the `UserRole` and `GroupRole` enum schemas and every reference:
  `User.userRole`, `UserView.userRole`, `SignUpCommand.userRole`, `GroupMember.groupRole`,
  `UpsertGroupMemberCommand.groupRole`, `Claims.userRole`. Preserve `appRoleId`/`groupRoleId`,
  `AppRole`/`GroupRoleDefinition`, and the `Permission` enum.
- [ ] Regenerate clients from `api/`: `./generate-client.sh desktop ../desktop/src/open-api` and
  `./generate-client.sh mobile ../mobile/api`. Desktop will break compile — expected, deferred to the
  frontend phase.

### ⬜ Phase 5 — Tests (NOT STARTED — must end green: `go test ./...`)

Fix the legacy-enum fallout (~30 `_test.go` files). Key items:
- `repositories/data_migrations_test.go` (highest effort) — the test DB's AutoMigrate will no longer
  create `user_role`/`group_role`, so seed them via raw DDL
  (`ALTER TABLE users ADD COLUMN user_role text`, same for `group_members.group_role`) + raw
  `UPDATE`s, and assert the backfill still maps to the right role ids. Add a case asserting the
  `HasColumn` guard skips cleanly when the columns are absent (fresh-install path).
- `middleware/group_test.go` — delete the `ValidateGroupRole` tests (middleware is gone).
- `structs/claims_test.go` — remove the `UserRole` enum-validation assertions/fixtures.
- Everywhere else — convert `SignUpCommand{UserRole: ADMIN/USER}` fixtures to `AppRoleID` (or the
  first-user fallback + existing `grantAllAppPerms`/`grantGroupPerms` helpers) and
  `GroupMember{GroupRole: OWNER/EDITOR/VIEWER}` fixtures to `GroupRoleID`. Many grep hits
  (`GroupRoleDefinition`, `GroupRoleID`, `LegacyUserRoleName`, `LegacyOwnerRoleName`) are false
  positives needing no change — verify before editing.

Remove stray `app.db` files from test dirs between runs (per `api/CLAUDE.md`).

### ⬜ Phase 6 — Docs (NOT STARTED)

- [ ] Update `api/CLAUDE.md` "Roles & Permissions": legacy enums removed from the backend (no longer
  used by the JWT, model, or `DeriveLegacy*`); the JWT no longer carries `userRole`; admin is defined
  by `app.users.read`; `DeriveLegacy*` deleted; the data migration is column-existence-guarded.
- [ ] Update monorepo `CLAUDE.md` Authorization note ("Rollout is additive and in progress").
- [ ] Delete this file once the work lands (it is a transient checkpoint artifact).

---

## Key design notes (for whoever resumes)

- **Admin = `app.users.read`.** "Legacy User" deliberately omits it; "Legacy Admin" has it — so it
  reproduces the old ADMIN/USER split. Used for the last-admin guard and first-admin-login check.
- **Bootstrap admin & group creator assign concrete roles** (Legacy Admin id / default group role
  id), which is correct even though *checks* are permission-based — you must assign a real role.
- **Wildcard-aware counting.** `CountUsersWithAppPermission` enumerates app roles and matches each
  role's grants with `permissions.HasAll`, so a custom role granting `*`/`app.*` is counted — a raw
  `WHERE permission = 'app.users.read'` would miss it.
- **No import cycle.** The paged-groups permission check lives in the handler, not the `commands`
  validator (`services` imports `commands`, so the reverse would cycle).
- **DB columns intentionally orphaned.** `user_role`/`group_role` physical columns are left in place
  (no drop migration) so the guarded one-time data migration still works on upgrading installs.

## Verification (once Phases 2–3 land)

1. `cd api && go build ./...` (green after Phase 2/3).
2. `cd api && go test ./...` (green after Phase 5; clean `app.db` between runs).
3. Fresh-install boot against an empty SQLite DB — no "no such column: user_role" panic; bootstrap
   admin lands on Legacy Admin (admin pages reachable).
4. Upgrade-install backfill — seed legacy columns + null role FKs, boot, verify back-fill runs.
5. Last-admin guard with a `app.*` wildcard custom role — second-to-last delete allowed, last
   blocked with `ErrLastAdmin`.
6. First-admin-login UX — first admin login true, second false.
7. Paged groups — non-admin `associatedGroup=ALL` → 403; admin → success; `MINE` works for everyone.
