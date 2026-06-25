# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Development
- `npm start` - Start development server with proxy configuration (serves on localhost:4200, proxies /api to localhost:8081)
- `npm run build` - Build production application
- `npm run watch` - Build in watch mode for development
- `npm test` - Run unit tests with coverage
- `npm test:ci` - Run tests in CI mode with ChromeHeadless
- `npm run e2e` - Run Playwright end-to-end tests (see **E2E Testing** below)
- `npm run e2e:ui` - Run Playwright tests in interactive UI mode
- `npm run e2e:install` - Install Playwright browser binaries (one-time setup)

### Build Configuration
- Production builds go to `dist/receipt-wrangler/`
- Development server uses proxy configuration in `proxy.conf.json` to route API calls to backend
- Angular CLI configuration in `angular.json`

## Code Architecture

### Application Structure
Receipt Wrangler Desktop is an Angular 19 application with modular architecture using:

- **State Management**: NGXS store with persistent storage for application state
- **API Layer**: Auto-generated OpenAPI client in `src/open-api/` (do not manually edit these files)
- **Component Architecture**: Feature modules with lazy-loaded routing
- **UI Framework**: Angular Material + Bootstrap 5 + custom shared components

### Key Architectural Patterns

#### Module Organization
- Feature modules (receipts, dashboard, groups, etc.) with their own routing
- Shared UI components in `src/shared-ui/` for reusable elements
- Lazy-loaded modules for performance optimization
- Centralized store management with NGXS states

#### State Management (NGXS)
- All application state managed through NGXS store
- State persistence configured for key data (auth, user preferences, table states)
- Individual state files for each feature (receipt-table.state.ts, group.state.ts, etc.)
- Actions and state updates follow NGXS patterns

#### Component Structure
- Feature components organized by domain (receipts/, dashboard/, groups/)
- Shared UI components provide consistent design patterns
- Form components use reactive forms with custom validation
- Table components use base table service pattern for pagination and filtering

### Key Directories

#### Core Application
- `src/app/` - Main application module and routing
- `src/store/` - NGXS state management (18+ state files)
- `src/services/` - Application services and business logic
- `src/guards/` - Route guards for authentication and authorization

#### Features
- `src/receipts/` - Receipt management (forms, tables, processing)
- `src/dashboard/` - Customizable dashboard widgets and views
- `src/groups/` - Group management and member administration
- `src/categories/` and `src/tags/` - Receipt organization features
- `src/auth/` - Authentication and user management
- `src/roles/` - Role & permission management (admin-only Manage Roles UI)

#### Shared Infrastructure
- `src/shared-ui/` - 30+ reusable UI components (buttons, forms, tables, dialogs)
- `src/pipes/` - Custom Angular pipes for data transformation
- `src/utils/` - Utility functions and helpers
- `src/open-api/` - Generated API client (auto-generated, do not edit)

### Testing Strategy
- Unit tests use Jasmine/Karma framework
- Code coverage reporting with minimum thresholds
- Tests exclude auto-generated API code (`src/open-api/`)
- CI tests run in headless Chrome

### Development Environment
- Angular CLI 21 with TypeScript 5.9
- Bootstrap 5 + Angular Material for UI components
- NGXS for state management with Redux DevTools integration
- Strict TypeScript configuration with comprehensive compiler options

### API Integration
- Backend API proxied through development server
- OpenAPI client generated from backend specification
- API base path configurable through environment
- HTTP interceptors handle authentication and error responses

### Code Conventions
- SCSS for styling with component-scoped styles
- TypeScript strict mode enabled
- Angular style guide followed for component organization
- Lazy loading for feature modules to optimize bundle size

### Use Established Patterns (do not invent one-offs)
New UI MUST reuse the application's established patterns and shared components rather than
inventing a divergent, one-off implementation of something the app already standardizes — **unless
the user explicitly confirms the divergence**. Examples of standards to follow:
- **Form actions:** the floating save bar fixed to the bottom of the page — `<app-form>` (which wraps
  `app-form-button-bar` + `app-submit-button`), or, for bespoke layouts, a plain
  `<form (ngSubmit)="...">` ending in a standalone `<app-form-button-bar [mode]="...">` containing an
  `<app-submit-button>` (see `src/receipts/receipt-form/` for the bespoke-layout precedent). Do NOT
  place Save/Cancel buttons in the page header.
- **Form fields:** `app-input`, `app-textarea`, `app-select`, `app-checkbox`, grouped with
  `app-form-section`; bind via the `formGet` pipe.
- **Tables:** `app-table`; **dialogs:** `app-dialog` + `app-dialog-footer`.
- **Simple filters:** the segmented `app-filter-bar` (`src/shared-ui/filter-bar/`) — pass `FilterTab[]`
  (`{ value, label, icon?, count? }`) and two-way bind the selected `value`.
- **Breadcrumbs:** `app-breadcrumb` with `BreadcrumbItem[]`.
If a design appears to require a new pattern, confirm with the user before diverging.

### Roles & Permissions (Manage Roles)

The admin-only **Manage Roles** feature (`src/roles/` — `role-list`, `role-form`, `role-presets`,
`roles.module`) provides CRUD for the backend's app- and group-scoped roles (see the backend
"Roles & Permissions" section in `api/CLAUDE.md` for the permission model). The Manage Roles routes are
gated by `appPermissionGuard` requiring `app.roles.read` (see **Permission-based UI gating** below).

- Talks to the backend via the **generated** clients in `src/open-api/` — `RoleService` for role
  CRUD and `PermissionService` (`GET /permission`) to load the permission catalog that populates the
  role editor's permission picker. `role-presets.ts` holds the role templates. Never hand-edit
  `src/open-api/` — regenerate it from `swagger.yml` instead.
- Built on existing shared patterns: `app-breadcrumb` and the segmented `app-filter-bar` (see "Use
  Established Patterns" above).
- **Category/tag grants (group roles only):** the role editor shows a "Category & tag access" section
  (gated on `showGrants()` = group scope) with `app-category-autocomplete` / `app-tag-autocomplete`
  (both fed the full pool via `CategoryService.getAllCategories()` / `TagService.getAllTags()` — the
  editor is admin-only). Selecting grants restricts members to those categories/tags; **empty = all**.
  The selections drive `FormArray`s loaded from `role.categoryGrants` / `role.tagGrants` (resolved to
  pool objects by an effect once the pool arrives) and serialized back as id arrays on
  `UpsertRoleCommand` for group scope only. The grant pickers pass `[creatable]="false"` (pick from
  existing, never create). See `api/CLAUDE.md` → "Data model".
- **Paid-by visibility (group roles only):** the same group-scoped grants section also shows a
  "Paid-by visibility" picker — a single `app-autocomlete` multi-select over `paidByOptions()` (a
  pinned **"Their own receipts"** sentinel option, id `OWN_PAID_RECEIPTS_OPTION_ID = -1`, followed by
  every user from `UserState.users`). On submit the selections split into `includeOwnPaidReceipts`
  (the sentinel is present) and `paidByUserGrants` (the remaining user ids, sentinel excluded); on
  edit they rehydrate from `role.paidByUserGrants` / `role.includeOwnPaidReceipts` via an effect that
  filters the shared `paidByOptions()` (stable references so the autocomplete excludes selected
  options). Empty = members see every payer's receipts; it restricts which receipts a member can see,
  not what they can edit. See `api/CLAUDE.md` → "Paid-by visibility enforcement".
- **Default roles:** the role-list page shows two `app-select` controls above the filter bar —
  "Default application role" and "Default group role". Each is pre-selected from the role flagged
  `isDefault` for its scope and, on change, calls `RoleService.setDefaultRole(scope, roleId)` then
  reloads (setting one default clears the previous one). The default role per scope is what new
  accounts / group creators receive (see `api/CLAUDE.md` → "Default roles"); the current default
  cannot be deleted. Default rows also carry a "Default" badge next to the System badge.
- **Modern role assignment in authoring forms:** the user add/edit form (`src/user/user-form/`) and
  the group-member add/edit form (`src/group/group-member-form/`) assign **modern roles** — an
  `app-select` of `RoleService.getRoles()` filtered to the `APP` / `GROUP` scope, bound to
  `appRoleId` / `groupRoleId` (not the legacy enums). Add forms pre-select the configured default
  role. Each selector has a Preview icon button (`data-testid="role-preview"`) that opens the shared
  **`RolePreviewDialogComponent`** (`src/roles/role-preview/`, standalone, opened via
  `openRolePreviewDialog(dialog, role)`) — a read-only dialog rendering the role's scope, description
  and permissions (grouped by resource using the `role-presets.ts` helpers). The `getRoles()` calls
  use `catchError` so a non-admin who lacks `app.roles.read` (see **Permission-based UI gating** below)
  gets an empty selector rather than an error.
- **Member-table role display:** the admin user-list (`src/user/user-list/`) and `group-form`
  (`src/group/group-form/`, loaded by every group view including non-admins) resolve each user's
  `appRoleId` / each member's `groupRoleId` to the role **name** via the shared `RoleNamePipe`
  (`src/pipes/role-name.pipe.ts` — `{{ id | roleName : roles() : scope }}`). The pipe matches on
  **id _and_ `PermissionScope`** because app- and group-role ids are independent sequences and can
  collide (e.g. group role `id=1` vs app role `id=1`); callers pass `PermissionScope.App` (user-list)
  / `PermissionScope.Group` (group-form) so the wrong-scope role is never matched. Both load roles with
  `RoleService.getRoles()` wrapped in `catchError`, so a non-admin lacking `app.roles.read` sees a
  blank name rather than a 403-driven logout. There is no longer any legacy `groupRole` enum sync,
  and `group-form` no longer enforces a "keep an owner" rule — the backend dropped the owner concept,
  so group management is governed entirely by `group.*` permissions.
- **Permission-based UI gating.** The UI gates on the user's effective permissions, mirroring the
  backend's enforcement. Permissions are delivered on **AppData** (`appPermissions: string[]` and
  `groupPermissions: { [groupId]: string[] }`) and stored in `AuthState` via the dedicated
  `SetPermissions` action — dispatched **only** from `setAppData` (`utils/app-data.utill.ts`), never from
  `TokenRefreshService` (whose claims-only `SetAuthState` must not wipe permissions). They refresh on
  login + app-init; the server re-checks real permissions on every request, so the stored set is a UI
  hint (a stale button at worst 403s, handled by the interceptor).
  - **403 handling (`src/interceptors/http-interceptor.ts`).** The backend returns **403 for every
    access denial** (auth *and* permission — it never uses 401). With a still-valid token a 403 is a
    permission denial, so the interceptor surfaces it via a **Forbidden toast** (only for
    user-initiated mutations — non-`GET` — and never in `queueMode`) and **re-throws without logging
    the user out**. It does **not** refresh/retry on 403: token freshness is handled proactively
    elsewhere (15-min timer in `app.component.ts`, app-init, and `auth.guard`), and
    `TokenRefreshService` keeps its own logout-on-refresh-failure path for a truly dead session.
    Background `GET` 403s propagate silently for callers to handle (e.g. the `getRoles` +
    `catchError` reads above).
  - **Category/tag catalogs:** AppData also carries `groupCategories` / `groupTags` (keyed by group
    id, filtered to the user's grants), stored via `SetGroupCatalog` and read with the
    `AuthState.groupCategories(groupId)` / `groupTags(groupId)` selectors. The **receipt form** and
    **receipts-table filters** source their category/tag options from these per-group catalogs (not
    the global `GET /category` / `GET /tag`, which are now admin-only — the receipt routes no longer
    use the category/tag resolvers). The receipt form's pickers gate `[creatable]` on
    `app.categories.create` / `app.tags.create` so restricted users can only pick from the granted set.
  - **Matcher:** `src/utils/permission.utils.ts` — `matches`/`hasAll`/`hasAny`, a faithful port of the Go
    matcher (`api/internal/permissions/matcher.go`) including wildcard semantics, so UI gating === backend.
  - **Selectors** (`AuthState`): `hasAppPermission(perm)`, `hasAnyAppPermission(perms)`,
    `hasGroupPermission(groupId, perm, orApp = [])` — the group one applies the `orApp` app-scoped
    override first, mirroring the backend `OrAppPermissions` (admin-not-a-member) pattern.
  - **Directives** (`DirectivesModule`, signal/`effect`-driven so they re-render when AppData lands after
    first paint): `*hasAppPermission="Permission.X"` and
    `*hasGroupPermission="{ groupId, permission, orApp? }"`. Components expose the generated `Permission`
    const to reference it in templates.
  - **Route guards** (`src/guards/`): `appPermissionGuard` (`data: { appPermissions: [...] }`, ANY-of)
    and `groupPermissionGuard` (`data: { groupPermission, orAppPermissions?, useRouteGroupId? }`).
    `receiptGuardGuard` is unchanged (server-checked per-receipt access); `system-settings-landing.guard`
    redirects `/system-settings` to the first tab the user can read, and `settings-landing.guard`
    does the same for `/settings` (the avatar-menu "User Settings" link → first readable of
    User Profile `app.account.read` / User Preferences `app.user-preferences.read` / API Keys
    `app.api-keys.read`). The `/settings` shell and each tab route are `appPermissionGuard`-gated on
    those reads, the in-page tabs render conditionally on the same, and the avatar-menu button is
    gated by a `hasAnyAppPermission` signal.
  - **Retired** with this migration: `RoleGuard`, `GroupRoleGuard`, the `*appRole` `RoleDirective`, the
    `groupRole` `GroupRolePipe`, and `GroupUtil.hasGroupAccess`. The group-member legacy-enum bridge
    (`legacyGroupRoleFromRole`) and `AuthState.userRole`/`hasRole` are now **removed** as well, since
    the backend legacy `UserRole`/`GroupRole` enums (and the `userRole`/`groupRole` API fields) are gone.
  - **Behavior note:** create actions for categories/tags/custom-fields now gate on the granular
    `.create` permission, so a normal user (Legacy User holds `.create`) sees the **Add** button;
    **Edit/Delete** stay admin-only (`.update`/`.delete`). **Group creation** follows the same shape:
    the Create-Group FAB on the groups list (`group-table`), the sidebar speed-dial "Add Group"
    button, and the `/groups/create` route guard all gate on `app.groups.create`. Note the
    read/create asymmetry — Legacy User holds `app.groups.create` but **not** `app.groups.read`, so
    they create via the sidebar FAB (the groups-list page itself is `app.groups.read`-gated and off
    limits to them), exactly like categories/tags.
  - **Dashboard CRUD** (`group-dashboards.component.html`): the Add / Edit / Delete dashboard buttons
    gate on `group.dashboards.create` / `.update` / `.delete` via `*hasGroupPermission` (the group id
    comes from a `selectedGroupIdNum` computed). Previously ungated — the buttons rendered for every
    member and 403'd on the backend; now they only render for holders, matching the receipts-table.
  - **Notification delete** (`notification/notification.component.html`): the per-notification delete
    control gates on `app.notifications.delete` via `*hasAppPermission`.

## Signals & Zoneless Change Detection

This application uses Angular's signal-based reactivity model with zoneless change detection (`provideZonelessChangeDetection()`). All new code MUST follow these patterns.

### Signal Primitives — Decision Guide

| Need | Use | NOT |
|------|-----|-----|
| Mutable state | `signal()` | Plain class properties |
| Read-only derived value | `computed()` | `effect()` that copies signals |
| Writable derived state (resets on dependency change, can be overridden) | `linkedSignal()` | `effect()` that sets a signal |
| Sync signal state to imperative/external APIs (DOM, localStorage, canvas, analytics) | `effect()` | — |
| DOM measurement/manipulation after render | `afterRenderEffect()` | `effect()` + `setTimeout` |
| Async data fetching | `resource()` | Manual subscribe + signal set |
| Observable → Signal bridge | `toSignal()` | `subscribe()` + signal set |
| Signal → Observable bridge | `toObservable()` | — |

### signal() — Writable State
- Use for mutable, source-of-truth state in components or services.
- Prefer `signal()` over plain class properties — signals automatically notify Angular's change detection.
- Provide a custom equality function when needed to avoid unnecessary updates.

```typescript
count = signal(0);
items = signal<Item[]>([]);
```

### computed() — Derived State
- Use whenever a value is derived from other signals. Always prefer over `effect()` for derivations.
- Computed signals are lazy (not evaluated until read) and cached (not recalculated until dependencies change).
- Safe to perform expensive operations (e.g., filtering arrays) inside computed.

```typescript
fullName = computed(() => `${this.firstName()} ${this.lastName()}`);
filteredItems = computed(() => this.items().filter(i => i.active));
```

### linkedSignal() — Writable Derived State
- Use when a value normally follows a computation but can be manually overridden.
- Resets to the computed value when dependencies change, but allows `set()`/`update()`.
- Perfect for selections that reset when options change.

```typescript
// Resets to first option when options change, but user can select manually
selectedOption = linkedSignal(() => this.options()[0]);
```

### effect() — Side Effects (Last Resort)
- **NEVER** use `effect()` to derive state or copy signal values between signals. Use `computed()` or `linkedSignal()` instead.
- **ONLY** use for syncing to non-reactive/imperative APIs: logging, localStorage, canvas rendering, third-party UI libraries.
- Effects run during change detection. They do not need `allowSignalWrites` (removed in Angular 19).
- Use `afterRenderEffect()` instead when you need to read DOM properties (offsetWidth, etc.) after rendering.

```typescript
// GOOD: Syncing to localStorage
effect(() => {
  localStorage.setItem('theme', this.theme());
});

// BAD: Deriving state — use computed() instead
effect(() => {
  this.fullName.set(`${this.firstName()} ${this.lastName()}`); // ❌ NEVER DO THIS
});
```

### Signal Inputs — input() and input.required()
- Use `input()` for optional inputs with defaults. Use `input.required()` for required inputs.
- Signal inputs are read-only (`InputSignal`). Template binding syntax `[prop]="value"` is unchanged.
- Use `computed()` to derive values from inputs. Use `effect()` only for imperative side effects triggered by input changes.
- Use `model()` for two-way binding (component modifies a value based on user interaction, e.g., custom form controls).

```typescript
// Required input — no undefined in type
mode = input.required<FormMode>();

// Optional input with default
disabled = input(false);

// Optional input without default
tooltip = input<string>();

// Two-way binding
value = model<string>('');

// Deriving from inputs — use computed, NOT effect
displayText = computed(() => this.mode() === FormMode.Edit ? 'Save' : 'Create');
```

**Replacing ngOnChanges:** Convert input-watching logic from `ngOnChanges` to `computed()` (for derived values) or `effect()` (for imperative side effects like loading data).

```typescript
// Before (ngOnChanges)
ngOnChanges(changes: SimpleChanges) {
  if (changes['groupId']) this.loadData();
}

// After (effect for imperative side effect)
constructor() {
  effect(() => {
    const id = this.groupId();
    if (id) this.loadData(id);
  });
}
```

### Signal Outputs — output()
- Use `output()` instead of `@Output() + EventEmitter`. Template syntax `(event)="handler($event)"` is unchanged.
- Use `outputFromObservable()` when the source is an Observable.

```typescript
clicked = output<MouseEvent>();
// Emit: this.clicked.emit(event);
```

### Signal Queries — viewChild() / viewChildren()
- Use `viewChild()` / `viewChildren()` instead of `@ViewChild` / `@ViewChildren`.
- Access via signal call: `this.paginator()` instead of `this.paginator`.
- Use `viewChild.required()` when the element is guaranteed to exist (not behind `@if`).

```typescript
paginator = viewChild.required(MatPaginator);
optionalEl = viewChild<ElementRef>('myEl');
items = viewChildren(ItemComponent);
```

### RxJS Interop
- **`toSignal(observable)`**: Converts Observable to Signal. Creates a subscription — call once and reuse the signal, never call repeatedly. Automatically unsubscribes on destroy.
  - Provide `initialValue` for Observables that don't emit synchronously.
  - Use `requireSync: true` for BehaviorSubject or other synchronous sources.
- **`toObservable(signal)`**: Converts Signal to Observable. Only emits the latest stabilized value.
- **`takeUntilDestroyed()`**: Replaces `@UntilDestroy()` / `untilDestroyed(this)`. Use in constructor or pass `DestroyRef`.
- **`outputFromObservable()`**: Declares an output from an Observable source.

```typescript
// NGXS selector → signal (preferred pattern)
groups = this.store.selectSignal(GroupState.groups);

// HTTP Observable → signal
data = toSignal(this.http.get<Data>('/api/data'), { initialValue: [] });

// Cleanup subscriptions
constructor() {
  this.someObservable$.pipe(
    takeUntilDestroyed(),
  ).subscribe(val => this.doSomething(val));
}
```

### NGXS State Access
- Use `store.selectSignal()` instead of `@Select` decorator for template-bound state. Returns a `Signal<T>`.
- `store.selectSnapshot()` remains valid for synchronous one-time reads in methods.
- Remove `| async` pipe from templates — use signal reads `()` instead.

```typescript
// Before
@Select(AuthState.isLoggedIn) isLoggedIn!: Observable<boolean>;
// Template: *ngIf="isLoggedIn | async"

// After
isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);
// Template: @if (isLoggedIn()) { ... }
```

### Zoneless Change Detection Rules
Angular no longer uses zone.js. Change detection is triggered ONLY by:
1. **Signal writes** — `signal.set()`, `signal.update()`, `computed()` recalculation
2. **`ChangeDetectorRef.markForCheck()`** — for non-signal reactive patterns (AsyncPipe calls this automatically)
3. **Template event bindings** — `(click)="handler()"` automatically triggers CD
4. **`ComponentRef.setInput()`** — programmatic input setting

**Key implications:**
- Plain property mutations (`this.foo = 'bar'`) in async callbacks (subscribe, setTimeout, Promise.then) will NOT trigger change detection. Always use signals for state that affects templates.
- `ChangeDetectorRef.detectChanges()` still works but is rarely needed — prefer signals.
- `setTimeout` still works for delays but won't auto-trigger CD. The callback must write to a signal if the template needs updating.
- All `@HostListener` handlers automatically trigger CD (same as template events).

### Testing with Zoneless
- Add `provideZonelessChangeDetection()` to `TestBed.configureTestingModule` providers.
- Prefer `await fixture.whenStable()` over `fixture.detectChanges()` for most realistic test behavior.
- Use `TestBed.flushEffects()` when testing effect-based logic.

## E2E Testing

End-to-end tests live in `e2e/` and use **Playwright**. They drive the real Angular UI against a real Go API. Config is `playwright.config.ts`.

### Running locally

1. **One-time:** install browsers — `npm run e2e:install`.
2. **One-time:** sign up the two e2e accounts against your local DB. The **first** signup is auto-promoted to admin, so order matters. With the API running, go to `http://localhost:4200/auth/sign-up` and create:
   - Admin first: username `e2e-admin`, password `e2e-admin-password`
   - Then user: username `e2e-user`, password `e2e-user-password`
3. **Every run:** source the dev env script so the `E2E_*` vars are exported:
   ```bash
   cd ../api/dev && source switch-to-sqlite.sh && cd -
   ```
   (`switch-to-mariadb.sh` / `switch-to-postgresql.sh` work the same — all three export the same `E2E_*` defaults.)
4. Start the Go API separately (`cd ../api && go run main.go`). Playwright auto-starts the Angular dev server via its `webServer` config, but it cannot launch the API.
5. Run the tests: `npm run e2e` (or `npm run e2e:ui` for watch-style debugging).

### CI

In CI the same spec files run against the demo URL. GitHub secrets populate the `E2E_*` vars — point `E2E_BASE_URL` at `https://demo.receiptwrangler.io` and supply the secret credentials. When `E2E_BASE_URL` is remote, the config skips the `webServer` block and does not start a local dev server.

### Best practices (follow these when adding new e2e tests)

**Locators — prefer `data-testid`; auto-retrying selectors only.**
- **Use `page.getByTestId(...)` as the standard selector.** Icon-only controls (the shared
  `app-add-button` / `app-edit-button` / `app-delete-button` / `app-cancel-button`, filter/menu icon
  buttons, etc.) and any element without a stable accessible name **must** carry a `data-testid`. Name
  it `<resource>-<action>` — e.g. `group-delete`, `comment-delete`, `receipt-duplicate`,
  `add-group-member`, `dialog-submit-button`. The `data-testid` passes through the shared button
  components to the host element, so `getByTestId('comment-delete')` resolves it directly.
- `page.getByRole('button', { name: '...' })` / `page.getByLabel(...)` / `page.getByPlaceholder(...)`
  remain fine for elements that already have a real accessible name (text buttons, labelled inputs).
- **Never** use structural CSS chains (`page.locator('app-receipt-comments app-delete-button')`) or raw
  CSS/XPath (`page.locator('.btn-primary')`) — they're brittle to component-structure refactors. Add a
  `data-testid` to the control instead.

**Assertions — rely on web-first expects, never `waitForTimeout`.**
- Use `await expect(locator).toBeVisible()`, `toHaveText()`, `toHaveURL()`, `toHaveCount()` — they auto-retry until `expect.timeout`.
- Never `await page.waitForTimeout(ms)` — it's a fixed sleep and flakes.
- Prefer `await page.waitForURL(/.../)` or `await page.waitForResponse(...)` for navigation/network waits.

**Isolation — each test gets a fresh `BrowserContext`.**
- No cookies/localStorage/session leak between siblings.
- Do NOT hand-write state-sharing between tests. If two tests need a logged-in session, use Playwright's `storageState` pattern (see below), not module-level globals.

**Auth — reuse login state, don't re-login in every test.**
- Current suite is tiny (login IS the test), so each test logs in via the UI. Fine for now.
- When the suite grows, switch to the **setup project** pattern: a `*.setup.ts` file logs in once and saves `storageState` to `e2e/.auth/<role>.json`; other tests declare `test.use({ storageState: 'e2e/.auth/user.json' })`. Keep `.auth/` git-ignored — it contains session cookies.
- One storageState file per role (admin, user). Never share one login across roles.

**`webServer` — for processes Playwright can launch.**
- The config uses `webServer` to start `npm start` when `E2E_BASE_URL` is localhost, and skips it when the URL is remote. `reuseExistingServer: !process.env.CI` lets local devs keep `ng serve` running between runs.
- Playwright cannot launch the Go API — that's always a separate process.

**Env vars and secrets.**
- Read via `process.env.E2E_*` — never hardcode credentials.
- Local defaults come from `api/dev/switch-to-*.sh`. CI values come from GitHub secrets.
- Never commit `.env` files or `e2e/.auth/` artifacts.

**Parallelism and flake budget.**
- `fullyParallel: true` is on. Tests must not mutate shared server state in ways that collide (same DB row, same uploaded file, same group membership). When you need mutation, create unique data per test (timestamp/UUID in names) and clean up after.
- `retries: 2` in CI, `0` locally — a test that only passes with retries is a bug, not a feature. Fix the root cause.
- `trace: 'on-first-retry'` captures a trace file on the first retry; view with `npx playwright show-trace <file>`. Do not set `trace: 'on'` — too heavy.

**Writing selectors for this app.**
- Forms use a custom `<app-input>` wrapper over `<mat-form-field>`. `page.getByLabel('Username')` resolves through the `<mat-label>` association.
- Submit buttons use `<app-button>` rendering `<button>` with visible text — `page.getByRole('button', { name: '...' })` works directly.
- Error feedback is often a Material snackbar (not inline `<mat-error>`). When asserting errors, locate the snackbar container or its text, not the form.

### Permission-gating specs (provisioned roles/users/groups)

Negative-permission coverage needs an account/role that *lacks* the permission under test, which no
seeded account provides — so an admin context **provisions a custom role (and user/group) through the
real UI** in `beforeAll` and **tears it down through the admin API** in `afterAll`. Shared flows live
in `e2e/helpers/provisioning.ts`: `createRole` (role form — type, preset, category toggles, individual
toggle-offs), `createUserWithRole`, `createGroupWithMember`, `uniqueName`, and the API-teardown
helpers `withAdminApi` + `apiDeleteUserByName` / `apiDeleteGroupById` / `apiDeleteRoleByName`.

- An admin `BrowserContext` (`storageState: 'e2e/.auth/admin.json'`) provisions in `beforeAll`. Tests
  then run **either** as the default e2e-user — for a *group-scoped* member added to a fixture group
  (Angular re-fetches AppData on every navigation, so a membership added after the saved session still
  drives the route guards without re-login) — **or** as a freshly-provisioned *custom user* whose
  session is captured to a git-ignored `e2e/.auth/<name>.json` (wait for the held permission in
  `localStorage.auth` before saving) and `rmSync`-ed in `afterAll`.
- **Teardown is API-based, not UI.** The role-list delete button is disabled while a role is assigned,
  and the UI's *bulk* user-delete dummy-converts a group-owning user (so the role stays assigned and
  never deletable) — UI teardown leaks roles. `withAdminApi` logs in via `request.newContext` (through
  the dev-server `/api` proxy) and `DELETE /api/user/{id}` **hard-deletes** (freeing the app-role) /
  `DELETE /api/group/{id}` frees the group-role assignment, so the role then deletes. **Order:** delete
  the user/group first, then the role; best-effort `try/catch` so a cleanup error doesn't mask the result.
- Reference specs: `system-settings-tab-gating.spec.ts` (custom app role + user + storageState — note:
  its UI teardown leaks roles, the reason the new specs use API teardown), `group-viewer-visibility.spec.ts`
  (group member with a group role), `search-bar-visibility.spec.ts` (no `app.receipts.search` → header
  search bar never renders), `dashboard-read-redirect.spec.ts` (no `group.dashboards.read` →
  `/dashboard/group/:id` redirects to `/receipts/group/:id`; an owner contrast still sees the dashboard),
  `paid-by-visibility.spec.ts` (a group role limited to "their own receipts" → a hidden-payer receipt
  `GET` 403s and is absent from the list, the member's own 200s; uses `withApiAs`/`apiCreateReceipt` and
  the `createRole` `paidByOwn` option in `helpers/provisioning.ts`),
  `dashboard-crud-gating.spec.ts` (a Viewer holding `group.dashboards.read` but not create/update/delete
  sees no Add/Edit/Delete dashboard buttons; owner contrast does),
  `comment-gating.spec.ts` (Receipt-Editor-preset members minus `group.comments.create` / `.delete` →
  no composer / no delete control; uses `apiCreateComment`),
  `receipt-feature-gating.spec.ts` (Quick Scan / Poll Email / Magic Fill controls hidden for a Viewer —
  positive contrast is a `test.fixme` because all three also sit behind the `aiPoweredReceipts` feature
  flag, which is `false` in the dev/CI API),
  `receipt-action-gating.spec.ts` (a Legacy Viewer sees no duplicate/delete row action, the
  `/receipts/:id/edit` route redirects, and `POST /api/receipt` **403s** via `withApiAs('user')`).
  `legacy-user-visibility.spec.ts` likewise carries **API-403** assertions (`DELETE /api/category|tag/:id`)
  so server enforcement is proven, not just the hidden control. Note: the receipts-table **edit** action
  is not template-gated (only duplicate/delete are); the edit *route* is guarded, so the edit denial is
  asserted at the route level, not button absence. (`receipt-action-gating.spec.ts` is a standalone
  spec rather than an extension of `group-viewer-visibility.spec.ts`, whose serial block has a known
  pre-existing failure — a Legacy User can't load `/groups` — that would skip any test appended to it.)

## Testing Requirements

**All new code must have accompanying unit tests.**

Before considering any work complete:

1. Write unit tests for all new components, services, and pipes
2. Use Angular TestBed for component testing
3. Mock services and HTTP calls appropriately
4. Run the full test suite: `npm test`
5. Ensure all tests pass before submitting changes

Tests should cover:

- Component rendering and user interactions
- Component method inputs and outputs
- Service method behavior
- Form validation logic
- Error handling scenarios