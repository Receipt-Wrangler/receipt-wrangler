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
"Roles & Permissions" section in `api/CLAUDE.md` for the permission model). It is routed behind the
global `UserRole.Admin` check in `RoleGuard` (`src/guards/role.guard.ts`).

- Talks to the backend via the **generated** clients in `src/open-api/` — `RoleService` for role
  CRUD and `PermissionService` (`GET /permission`) to load the permission catalog that populates the
  role editor's permission picker. `role-presets.ts` holds the role templates. Never hand-edit
  `src/open-api/` — regenerate it from `swagger.yml` instead.
- Built on existing shared patterns: `app-breadcrumb` and the segmented `app-filter-bar` (see "Use
  Established Patterns" above).
- **Default roles:** the role-list page shows two `app-select` controls above the filter bar —
  "Default application role" and "Default group role". Each is pre-selected from the role flagged
  `isDefault` for its scope and, on change, calls `RoleService.setDefaultRole(scope, roleId)` then
  reloads (setting one default clears the previous one). The default role per scope is what new
  accounts / group creators receive (see `api/CLAUDE.md` → "Default roles"); the current default
  cannot be deleted. Default rows also carry a "Default" badge next to the System badge.
- **Permission-based UI gating is not implemented yet.** There is no `*hasPermission` directive or
  `permissionService.has(...)` helper; UI access is still gated by the global admin `RoleGuard` and
  the legacy `UserRole`. When permission-driven gating is added, document the directive/guard here.

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

**Locators — use user-facing, auto-retrying selectors.**
- Prefer `page.getByRole('button', { name: 'Login' })`, `page.getByLabel('Password')`, `page.getByPlaceholder(...)`, `page.getByText(...)`.
- Use `page.getByTestId(...)` only when no accessible role/label exists. The codebase has no `data-testid` convention yet — add one on a component only when role/label truly can't identify the element.
- Avoid raw CSS/XPath (`page.locator('.btn-primary')`) — brittle to refactors.

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