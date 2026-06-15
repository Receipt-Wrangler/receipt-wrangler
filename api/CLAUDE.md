# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Receipt Wrangler API is a Go-based backend service for a receipt management and splitting application. It provides OCR-powered receipt scanning, AI-assisted data extraction, email integration, and multi-user support with group management capabilities.

## Development Commands

### Building and Running
- `go build` - Build the application
- `go run main.go` - Run the application directly
- `./set-up-dependencies.sh` - Install system dependencies (tesseract, ImageMagick, Chromium, Python deps)

### Testing
- `go test -v ./...` - Run all Go tests with verbose output
- `go test -coverprofile=coverage.out -covermode=atomic -v ./...` - Run tests with coverage
- `python3 -m unittest discover -s ./imap-client` - Run Python IMAP client tests

### API Client Generation
- `./generate-client.sh desktop <output-dir>` - Generate TypeScript Angular client
- `./generate-client.sh mobile <output-dir>` - Generate Dart Dio client

## Architecture Overview

### Core Structure
- **main.go** - Application entry point, initializes logging, config, database, and starts HTTP server
- **internal/** - Core application code organized by domain
- **imap-client/** - Python-based email processing client

### Key Directories
- **internal/handlers/** - HTTP request handlers for each API endpoint
- **internal/repositories/** - Database access layer using GORM
- **internal/services/** - Business logic layer
- **internal/models/** - Database models and domain objects
- **internal/commands/** - Command objects for API requests/responses
- **internal/routers/** - Route definitions and middleware setup
- **internal/wranglerasynq/** - Background job processing using Asynq
- **internal/ai/** - AI client implementations (OpenAI, Gemini, Ollama)

### Database
- Uses GORM ORM with support for SQLite, MySQL, and PostgreSQL
- Migrations are handled automatically on startup via `repositories.MakeMigrations()`
- Test databases are set up in `repositories/main_test.go`

### Background Processing
- Uses Hibiken's Asynq library for background job processing
- Email processing, OCR, and cleanup tasks run as background jobs
- Queue configurations defined in `internal/wranglerasynq/`

### AI Integration
- Supports multiple AI providers: OpenAI, Google Gemini, and Ollama
- AI clients implement a common interface defined in `internal/ai/base_client.go`
- Used for receipt data extraction and processing

### Configuration
- Configuration loaded from JSON files in `config/` directory
- Environment variables override config file settings
- Sample configuration in `config/config.sample.json`

## Testing Patterns

Each package typically has:
- `main_test.go` - Test setup and teardown
- `*_test.go` - Unit tests for specific functionality
- Test utilities in `internal/utils/testing.go` and `internal/repositories/testing.go`

Tests use dependency injection patterns and mock implementations for external services.

## Testing Guidelines for Claude

When working with tests in this codebase, follow these critical requirements:

### Test Execution Requirements
- **ALWAYS run tests after writing them** - When asked to write tests, you MUST run them to verify they pass
- **Report coverage** - Always report the coverage of files impacted by the tests using `go test -coverprofile=coverage.out -covermode=atomic`
- **Verify all tests pass** - Never consider test writing complete until all tests are verified to pass

### Test Database Cleanup
- **Failed tests may leave behind `app.db` files** in test directories (e.g., `services/app.db`, `handlers/app.db`)
- **These MUST be removed** before rerunning tests to avoid conflicts
- **CRITICAL**: Only remove `app.db` files from test directories, NEVER delete anything from the `sqlite/` directory
- Example cleanup locations: `internal/services/app.db`, `internal/handlers/app.db`, etc.

### Test Workflow
1. Write tests following existing patterns in the codebase
2. Run tests to verify they pass: `go test -v ./...`
3. Generate and report coverage: `go test -coverprofile=coverage.out -covermode=atomic -v ./...`
4. If tests fail, check for and remove any `app.db` files in test directories
5. Re-run tests until all pass
6. Report final coverage results for impacted files

## OCR and Image Processing

- Tesseract OCR integration via `otiai10/gosseract`
- ImageMagick integration for image processing and format conversion
- Supports HEIC format conversion to standard image formats
- Python dependencies for additional image processing capabilities

## Email HTML to PDF Rendering

- HTML email bodies are rendered to PDF via `chromedp` running headless Chromium
- The Chromium binary path is read from the `CHROMIUM_BINARY_PATH` env var
  (defaults to `/usr/bin/chromium`); installed by `set-up-dependencies.sh`
- The Chromium process sandbox is **off by default** because the supported
  docker images run as root, where chromium's sandbox refuses to start.
  Operators running the API as a non-root user can opt back in by setting
  `CHROMIUM_SANDBOX=true`
- External network resource loads (remote images, CSS, fonts) are
  **blocked by default** to remove an SSRF / tracking-pixel surface.
  Inline `data:` URIs and the file:// page itself remain allowed. To
  permit remote loads (useful when receipts depend on remote logos or
  product imagery), set `CHROMIUM_ALLOW_EXTERNAL_RESOURCES=true`
- Implementation: `internal/services/html_to_pdf.go` (HtmlToPdfService.Render)
- The rendered PDF is saved on the receipt as a `FileData` and routed
  through the existing `repositories.ConvertPdfToJpg` pipeline so vision and
  OCR models receive an image, exactly like a PDF email attachment
- Gating: only runs when `EmailBodyProcessingEnabled` is true on at least one
  consuming group (per `shouldRenderEmailBodyPdf` in `wranglerasynq/email.go`)
- For an email with both an attachment and an HTML body, each per-attachment
  receipt is augmented with a copy of the body PDF and both images are sent
  to the LLM together; when the body is sent as an image, its text is dropped
  from the prompt to avoid duplication

## Testing Requirements

**All new code must have accompanying unit tests.**

Before considering any work complete:

1. Write unit tests for all new functions and endpoints
2. Follow existing test patterns in the codebase (see `main_test.go` files for setup)
3. Mock external dependencies (database, services, etc.)
4. Run the full test suite: `go test -v ./...`
5. Ensure all tests pass before submitting changes

Tests should cover:

- Happy path scenarios
- Error handling and edge cases
- Input validation
- Authentication/authorization logic

## API Documentation

- OpenAPI 3.1 specification in `swagger.yml`
- API serves on port 8081 by default
- All endpoints require JWT authentication except login/signup

## Roles & Permissions

A configurable role/permission system. Administrators can define roles from granular permission
strings at two scopes — **application** and **group** — and assign them to users / group members.
**Handlers now enforce these permissions** (see "Enforcement status" below). The legacy
`models.UserRole` (`ADMIN`/`USER`) and `models.GroupRole` (`OWNER`/`EDITOR`/`VIEWER`) enums have
been **removed from the backend** — the Go types, the `User.UserRole`/`GroupMember.GroupRole` model
fields, the `Claims.userRole` JWT field, and the `DeriveLegacy*` shims are all gone. "Admin" is now
defined by the app permission `app.users.read` (the seeded **Legacy Admin** role grants it; **Legacy
User** omits it). The only legacy remnant is the **physical** `user_role`/`group_role` DB columns,
intentionally retained on existing installs so the one-time data migration can still back-fill from
them (see "Legacy role assignment" below); GORM never drops them and fresh installs never create
them.

### Permission registry

- `internal/permissions/registry.go` is the **hardcoded source of truth** for every permission.
  Each entry is a `Descriptor{ Key, Label, Description, Category, Scope }`; `Scope` is `APP` or
  `GROUP`. Helpers: `All()`, `Get(key)`, `Exists(key)`.
- **String format:** `scope.domain[.subdomain].action` — e.g. `app.users.create`,
  `group.receipts.read`. Permissions are **CRUD-granular** (`create`/`read`/`update`/`delete` per
  domain); distinct non-CRUD actions stand alone (e.g. `app.system-settings.restart-task-server`,
  `group.activities.rerun`, `group.receipts.duplicate`).
- Exposed to clients via `GET /permission` (returns the descriptors for the role-editor UI) and
  mirrored in the `swagger.yml` `Permission` enum. `permissions/registry_test.go` enforces that the
  registry and the swagger enum stay in sync.
- **Adding a permission:** add the constant + `Descriptor` in `registry.go`, add the key to the
  `Permission` enum in `swagger.yml`, then regenerate clients (see "API Client Generation").

### Matcher

- `internal/permissions/matcher.go` is a **pure** matcher over a granted `[]string`:
  - `HasAll(granted, required...)` — logical AND (the default; a single-permission check is just
    `HasAll(granted, "x")`).
  - `HasAny(granted, required...)` — logical OR.
- Wildcards in a *granted* string are honored: `*` matches anything, a trailing `group.*` matches
  any deeper key, and a mid-segment `*` matches exactly one segment. Both helpers deny when no
  required permission is supplied. The `:sub-scope` suffix (e.g. `read:any`) is matched literally —
  `:any` superset semantics are **not implemented yet**.

### Data model

- **App roles:** `AppRole` + `AppRolePermission` (permission strings in a child table).
- **Group roles:** `GroupRoleDefinition` + `GroupRolePermission`, plus `GroupRoleCategoryGrant` /
  `GroupRoleTagGrant` (composite-PK join rows) for per-role category/tag visibility. A group role
  with **no** grant rows is **unrestricted** (sees every category/tag); a non-empty grant set
  restricts members to exactly those ids — restriction is opt-in, so legacy/system roles (no grants)
  keep seeing everything and no data migration is needed. Categories/tags are **global** (no
  `GroupId`); a grant is a per-group-role slice of the global pool. Enforcement of these grants is
  rolled out in later slices; Phase 1 only persists/returns them through role CRUD.
- **Assignment:** nullable FKs `User.AppRoleID` and `GroupMember.GroupRoleID` (one app role per
  user; one group role per group membership). Nullable because per-create assignment is best-effort
  (the FK is left `nil` rather than failing creation when no role can be resolved, e.g. an unseeded
  test DB) and because the one-time migration back-fills pre-existing rows that start `nil`.
- `IsSystem` marks protected, non-editable/non-deletable roles; `IsDefault` (on **both** `AppRole`
  and `GroupRoleDefinition`) marks the single default role for that scope — the role assigned to new
  accounts (app) or to group creators (group). See "Default roles" below.

### Role CRUD

- Data access in `repositories/roles.go`; business logic in `services/roles.go`. Guards: system
  roles are immutable/undeletable, a role's scope can't be changed (type-mismatch error), an
  assigned role can't be deleted, and the **default role for a scope can't be deleted**
  (`ErrRoleIsDefault`) — pick another default first.
- Endpoints (`routers/role.go`, gated by `app.roles.*`): `GET /role`, `POST /role`,
  `PUT /role/{roleId}`, `PUT /role/{roleId}/default?scope=APP|GROUP` (make this role the scope's
  default; allowed on system roles; gated by `app.roles.update`), `DELETE /role/{roleId}?scope=APP|GROUP`.
- `commands.UpsertRoleCommand` validates that every permission exists in the registry and matches
  the role's scope. It also carries `categoryGrants` / `tagGrants` (category/tag ids): grants are
  rejected on APP scope and dedup-checked in `Validate()`, and their existence is verified against
  the DB in `RoleService` (`ErrInvalidGrant` → 400). `repositories/roles.go` syncs grant rows with
  the same delete-then-insert pattern as permissions (`replaceGroupRoleGrants`, with the nested
  Category/Tag association `Omit`-ted so only join rows are written), preloads them in
  `GetGroupRoleById` / `GetAllRoles`, and cascades them on role delete. `structs.RoleView` is the
  read model (includes `assignedCount`, `isDefault`, and `categoryGrants` / `tagGrants` — empty
  slices for app roles).

### Seeded system roles (legacy-equivalent)

- `repositories.SeedSystemRoles` (`repositories/seed_roles.go`) seeds five immutable
  (`IsSystem = true`) roles on startup — wired into `InitDB`, runs in all deploy envs (it is
  structural, unlike the bootstrap admin user). Their permission sets reproduce the capabilities of
  the historical `ADMIN`/`USER` app roles and `OWNER`/`EDITOR`/`VIEWER` group roles **exactly**, so
  upgrading installs see **zero behavior change**: **Legacy Admin** (every app permission), **Legacy
  User** (the app actions a plain `USER` could do), **Legacy Viewer** / **Legacy Editor** / **Legacy
  Owner** (the group VIEWER / EDITOR / OWNER tiers; Owner = every group permission). The sets live in
  `permissions/legacy.go` (`Legacy*Keys()` helpers) and were derived from the actual handler-level
  gating, not the desktop UI presets. **One deliberate exception:** `app.users.read` is omitted from
  Legacy User — it gates only the admin `GET /user/` listing, which no client calls (user dropdowns read
  from AppData via `app.account.read`), so granting it would only expose the admin "Manage Users" page to
  normal users. Removing it keeps normal users off that page with zero functional loss.
- `SeedSystemRoles` creates the roles with `IsDefault = false`; the **default** per scope is set
  separately by `EnsureDefaultRoles` (see "Default roles" below), the one-time data migration assigns
  the roles to existing users/members, and enforcement is wired in `HandleRequest`.
- Idempotent: keyed on role `Name` (a `uniqueIndex`), safe on every boot; a pre-existing
  same-named role is left untouched. The five role names are shared constants
  (`repositories/system_role_names.go`, `Legacy*RoleName`) used by both the seeder and the
  migration.
- **Known limitation:** because system roles are immutable and seeding skips existing names, a
  permission added to the registry later will **not** flow into an already-seeded Legacy Admin /
  Legacy Owner. Re-syncing system roles would need a dedicated reconciliation step (out of scope).

### Default roles

- Exactly one app role and one group role are the **default** (`IsDefault = true`): the role assigned
  to a new account on signup/admin-create, and to a group's creator on group creation. This is a
  required invariant — there is always exactly one default per scope.
- `repositories.EnsureDefaultRoles` (`repositories/seed_roles.go`) enforces it on every boot,
  immediately after `SeedSystemRoles` in `InitDB`: if a scope has **no** default, it flags the
  legacy-equivalent role (**Legacy User** for app, **Legacy Owner** for group), so upgrades and fresh
  installs behave exactly as before. It only acts when no default exists, so it never overrides a
  default an admin chose, and it self-heals dev DBs created before the `app_roles.is_default` column.
- Admins change the default via `PUT /role/{roleId}/default?scope=…`
  (`RoleService.SetDefaultRole` → `SetDefaultAppRole`/`SetDefaultGroupRole`), which clears the prior
  default and sets the new one in one transaction. System roles are eligible (the legacy defaults are
  system roles). The current default cannot be deleted (`ErrRoleIsDefault`).
- **Per-create assignment** uses the defaults: `UserRepository.CreateUser` sets `User.AppRoleID`
  (Legacy Admin for the first user — `isAdmin := usrCnt == 0`, resolved by `resolveAppRoleId(tx,
  isAdmin)` — so the bootstrap admin is never locked out; the default app role otherwise), and
  `GroupRepository.CreateGroup` sets the creator's `GroupMember.GroupRoleID` to the default group
  role. Both are **best-effort**: if the role can't be resolved (e.g. an
  unseeded test DB), the FK is left `nil` rather than failing creation. Members added to *existing*
  groups (and explicit role choices in the admin user/group-member forms) are assigned via the
  modern-role authoring flow — see "Modern role assignment in authoring flows" under "Enforcement
  status".

### Legacy role assignment (one-time data migration)

- A startup data migration back-fills the new role assignments from the legacy role values so
  existing installs upgrade with **zero behavior change**: each user's `user_role` maps onto the
  matching `User.AppRoleID` (`ADMIN` → Legacy Admin, `USER` → Legacy User) and each member's
  `group_role` onto `GroupMember.GroupRoleID` (`OWNER`/`EDITOR`/`VIEWER` → Legacy
  Owner/Editor/Viewer). Lives in `repositories/data_migrations.go` (`assignLegacyEquivalentRoles`).
- **Reads the retained physical columns, not Go fields.** The `UserRole`/`GroupRole` Go struct
  fields were removed, so the migration matches the `user_role`/`group_role` values as plain strings
  and **guards each back-fill loop with `tx.Migrator().HasColumn(...)`** — upgrading installs keep
  the physical columns (GORM never drops them) so the back-fill runs, while fresh installs never
  create them so the guard skips cleanly instead of erroring with "no such column". There is
  deliberately **no drop-column migration**, to preserve this upgrade path.
- **Tracking:** one-time data migrations are recorded in a `data_migrations` ledger
  (`models.DataMigration`, keyed by unique `name`) — distinct from GORM schema AutoMigrate. The
  runner `RunDataMigrations` skips any migration already in the ledger and otherwise runs it **and**
  writes the ledger row in a single `db.Transaction`, so a failure rolls back and retries next boot.
  Append new one-time migrations to the `dataMigrations` registry slice.
- Wired into `InitDB` **after** the bootstrap-admin step (roles must be seeded first, and this order
  also assigns a fresh install's first admin). `InitDB` is only called from `main.go`, never the
  test harness, so the migration does not auto-run in tests.
- Updates are guarded with `... IS NULL`, so a role an admin has already set through the new UI is
  never overwritten — defense-in-depth on top of the ledger. The migration is one-time over existing
  rows; per-create assignment for *new* users / group creators is handled by the default-role wiring
  (see "Default roles" above).
- Tests: `repositories/data_migrations_test.go` (assignment, idempotency, ledger short-circuit,
  no-clobber, and the `HasColumn`-guard skip path when the legacy columns are absent). The tests
  seed the legacy columns via raw `ALTER TABLE ... ADD COLUMN` DDL since AutoMigrate no longer
  creates them.

### Permission checks (`PermissionService`)

- `services/permission.go` exposes four **scope-separated** entry points:
  `HasAppPermissions` / `HasAnyAppPermission` and `HasGroupPermissions` / `HasAnyGroupPermission`
  (each App/Group pair = AND default + OR variant).
- Each call **resolves the user's current permissions from the database** (the user's app role for
  app checks; the group membership's group role for group checks) and matches them with the pure
  matcher. The **JWT is never trusted for authorization** — permissions are always re-read. A user
  with no assigned role, or a non-member of the group, resolves to no permissions (deny).
- Required keys are scope-guarded: passing an `app.*` key to a group check (or an unknown key)
  returns an error, catching call-site bugs.
- Backed by a small role-permission cache (`services/permission_cache.go`) keyed by `scope + roleId`
  and invalidated in `RoleService.UpdateRole` / `DeleteRole`. Only a role's permission *list* is
  cached; a user's role *assignment* is resolved fresh on every check, so re-assigning a user takes
  effect immediately.

### Enforcement status

Authorization is enforced centrally in `HandleRequest` (`handlers/generic_handler.go`) via the
`PermissionService`. Each handler declares its requirement on the `structs.Handler` it builds:

- `AppPermissions []string` — app-scoped permissions the caller must hold (logical AND).
- `GroupPermissions []string` — group-scoped permissions the caller must hold (AND) in **each**
  group resolved from `GroupId` / `GroupIds`, or from `ReceiptId` / `ReceiptIds` (the receipt's
  group is looked up automatically).
- `OrAppPermissions []string` — an app-scoped fallback; holding **any** of them bypasses the group
  check (e.g. an administrator viewing a group they aren't a member of). Replaces the old
  `OrUserRole`.

`HandleRequest` resolves the caller's effective permissions from the database (never the JWT) and
denies with `403` on any failure. The legacy `UserRole` / `GroupRole` / `OrUserRole` handler fields
and their checks have been **removed**. Every authenticated endpoint that previously had a legacy
role gate now has an equivalent permission gate; endpoints that touch only the caller's own data
(notifications, user preferences, own profile/claims/app-data, API keys, group lists) are gated by
dedicated self-service permissions (`app.notifications.*`, `app.user-preferences.*`,
`app.account.*`, `app.receipts.search`) included in the Legacy User set so existing users are
unaffected. Two endpoints are intentionally **not** permission-gated: the username-availability
lookup (used pre-auth during signup) and `ConvertToJpg` (a stateless image utility with no stored
resource to scope against). The role/permission management endpoints (`/role`, `/permission`) are
gated by `app.roles.*`.

**Effective permissions on AppData (desktop UI gating):** `GetAppData` (`services/auth.go`) includes the
caller's resolved permissions so the desktop can gate UI with them — `AppPermissions []string` and
`GroupPermissions map[uint][]string` (keyed by group id) — built via
`PermissionService.GetAppPermissionsForUser` / `GetGroupPermissionsForUser` (thin exported wrappers over
the cached `resolveAppPermissions` / `resolveGroupPermissions`). The JWT no longer carries any role
field at all (it holds only identity claims); the server always re-checks real permissions from the
DB on every request.

**`app.api-keys.read-any` (Security):** listing *all* users' API keys (`GetPagedApiKeys` with
`associatedApiKeys=ALL`) requires `app.api-keys.read-any`, checked in the handler body via the
`PermissionService` — the legacy `token.UserRole == ADMIN` check in
`commands/paged_api_key_request_command.go` was removed. Legacy Admin auto-includes the new permission
(its set is every app permission); Legacy User does not.

**Per-create role assignment (done):** a new account (signup or admin-create) is assigned the
default app role, and a group's creator is assigned the default group role, via the default-role
wiring (see "Default roles" above), so accounts created after the one-time migration are no longer
locked out.

**Modern role assignment in authoring flows (done):** admin user-create/update and group-member
create/update assign **modern roles directly**. `SignUpCommand` carries `AppRoleID` and
`UpsertGroupMemberCommand` carries `GroupRoleID`; `UserRepository.CreateUser`/`UpdateUser` and
`GroupRepository.CreateGroup`/`UpdateGroup` honor them when present. The legacy-enum bridge has been
**removed**: `DeriveLegacyUserRole`/`DeriveLegacyGroupRole` are deleted and no `UserRole`/`GroupRole`
is derived or written anymore. The admin create endpoint's role-required validation
(`middleware.ValidateUserData`) now accepts **only** `appRoleId`. Public `SignUp` strips a
caller-supplied `AppRoleID` so a sign-up can never self-assign a role.

### Tests

`permissions/matcher_test.go`, `permissions/registry_test.go`, `permissions/legacy_test.go`,
`services/permission_test.go`, `services/roles_test.go`, `repositories/roles_test.go`,
`repositories/seed_roles_test.go` (default-role seeding), the per-create assignment tests in
`repositories/users_test.go` / `repositories/groups_test.go`, and the handler authorization tests in
`handlers/generic_handler_test.go` (with shared helpers in `handlers/auth_test_helpers_test.go`).
