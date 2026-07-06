# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Receipt Wrangler Mobile is a Flutter mobile application that provides a native interface for Receipt Wrangler, a receipt management and splitting system. The app enables users to manage receipts on the go with camera/gallery uploads, receipt scanning, group management, and receipt splitting capabilities.

## Development Commands

### Core Flutter Commands
- `flutter run` - Run the app on connected device/emulator
- `flutter build apk` - Build Android APK
- `flutter build ios` - Build iOS app
- `flutter test` - Run unit tests
- `flutter analyze` - Analyze Dart code for issues
- `dart format .` - Format Dart code
- `flutter clean` - Clean build artifacts
- `flutter pub get` - Install dependencies
- `flutter pub upgrade` - Upgrade dependencies

### API Client
The project uses a generated OpenAPI client located in the `api/` directory. The client is imported as a local package dependency in pubspec.yaml.

## Architecture Overview

### State Management
The app uses Provider pattern with ChangeNotifier models:
- **AuthModel**: Authentication state, JWT tokens, API client configuration
- **GroupModel**: Group management and selection
- **ReceiptModel**: Receipt data, form state, and image handling
- **UserModel**: User profile and preferences
- **CategoryModel**, **TagModel**: Metadata management
- **SearchModel**: Search functionality with RxDart streams
- **PermissionsModel**: The caller's effective permissions, used for UI gating (see "Permission-based UI gating")

### Navigation
Uses `go_router` with nested shell routes:
- **Group Selection Shell**: `/groups` with group selection UI
- **Group Context Shell**: `/groups/:groupId/*` with group-specific navigation
- **Search Shell**: `/search` with search interface
- Individual routes for receipt forms, viewing, and editing

The receipt form routes (`/receipts/:receiptId/{view,edit}` and the comments variants) use
`pageBuilder` with `NoTransitionPage(key: state.pageKey, ...)` instead of `builder`. During a
default transition the outgoing and incoming `ReceiptFormScreen`s coexist for the animation, and
both attach the same global `receiptFormKey` — a GlobalKey collision that corrupts the form. The
`NoTransitionPage` makes the swap atomic; `key: state.pageKey` forces a fresh element per
navigation so the new screen re-fetches rather than reusing the old element. Because the swap
removes the old page in the same frame, popup-menu items that navigate must capture the router
and defer the `go` until the menu finishes dismissing (see
`receipt_app_bar_action_builder.dart`'s `_goAfterMenu`).

### Core Directory Structure
- `lib/models/` - Provider-based state management models
- `lib/auth/` - Authentication screens and logic  
- `lib/groups/` - Group management, dashboards, receipts
- `lib/receipts/` - Receipt forms, viewing, image handling
- `lib/search/` - Search functionality
- `lib/shared/` - Reusable widgets and utilities
- `lib/client/` - OpenAPI client wrapper
- `lib/utils/` - Utility functions for auth, currency, dates, etc.

### Key Features
- **Receipt Management**: Create, edit, view receipts with items and images
- **Image Handling**: Camera/gallery upload with scanning capabilities
- **Group Management**: Multi-user groups with role-based access
- **Search**: Full-text search across receipts
- **Offline Support**: Secure token storage with refresh token flow

### Form Handling
Uses `flutter_form_builder` for complex forms with validation. Receipt forms support:
- Dynamic item lists with custom fields
- Image carousel with infinite scroll
- Category and tag selection
- Currency formatting and validation

### API Integration
- Generated OpenAPI client from backend specification
- JWT-based authentication with automatic token refresh
- Centralized client configuration in `OpenApiClient` singleton
- Secure token storage using `flutter_secure_storage`

### Permission-based UI gating

The mobile app gates UI on the caller's **effective permissions**, mirroring the desktop client
(`desktop/CLAUDE.md` → "Permission-based UI gating") and the backend's enforcement (`api/CLAUDE.md` →
"Roles & Permissions"). The JWT no longer carries any role field — effective permissions are
delivered on `AppData` and are a UI hint; the server re-checks real permissions on every request, so
a stale action at worst returns 403.

- **Delivery & hydration:** permissions arrive on `AppData` (`appPermissions`, and `groupPermissions`
  keyed by group id) and are stored in **`PermissionsModel`** (`lib/models/permissions_model.dart`)
  via `setPermissions`, called **only** from `storeAppData` (`lib/utils/auth.dart`) on login and
  app-init. Token-only refreshes never touch it. The `FutureBuilder` in `main.dart` blocks first
  paint until app-init completes, so permissions are present before any gated widget builds.
- **Per-group category/tag catalogs:** `AppData` also delivers `groupCategories` / `groupTags` (keyed
  by group id), each filtered to the caller's group-role grants (the full pool when unrestricted).
  `storeAppData` indexes them into `CategoryModel` / `TagModel` (`setGroupCategories` /
  `setGroupTags`); the receipt-form pickers (`category_select_field.dart` / `tag_select_field.dart`,
  passed the receipt's `groupId`) source options via `categoriesForGroup(groupId)` /
  `tagsForGroup(groupId)`. **Do not** source the pickers from the flat `categories` / `tags` arrays —
  those are admin-only (populated only for holders of `app.categories.read` / `app.tags.read`), so a
  normal user would see an empty picker. Mirrors desktop's per-group AppData catalog selectors.
- **Matcher:** `lib/utils/permission_matcher.dart` (`permissionMatches` / `hasAll` / `hasAny`) is a
  faithful port of the backend matcher (`api/internal/permissions/matcher.go`) and its desktop twin
  (`desktop/src/utils/permission.utils.ts`), wildcard semantics included, so UI gating === backend.
  The generated `Permission` enum is converted to its wire string at hydration (effective
  permissions are always concrete registry keys, so the enum round-trips safely).
- **Checks** (`PermissionsModel`): `hasAppPermission(p)`, `hasAnyAppPermission([..])`,
  `hasGroupPermission(groupId, p, {orApp})` — the group one applies the `orApp` app-scoped override
  first (the backend `OrAppPermissions` admin-not-a-member pattern) — and
  `hasGroupPermissionInAnyGroup(p)` for screens with no single current group.
- **Gating pattern:** read `Provider.of<PermissionsModel>(context, listen: false)` and conditionally
  render, referencing the generated `Permission` constants. Most gating is **widget-level**, but a
  couple of permission-scoped routes also have **go_router redirects** (see "Route redirects" below) —
  a `redirect:` callback reads the same `PermissionsModel` from the provider tree (the router is built
  under the app's `MultiProvider`, exactly like the existing `/receipts/add` redirect), so the check
  is identical to the widget-level one.
- **Current widget gates:**
  - `canEditReceipt(permissionsModel, groupId)` (`lib/shared/functions/permissions.dart`) →
    `group.receipts.update`, used by the receipt-edit popup (`receipt_edit_popup_menu.dart`) and the
    receipt swipe-to-edit (`receipt_list_item.dart`).
  - `canCommentCreate` / `canCommentDelete` (same file) → `group.comments.create` /
    `group.comments.delete`, scoped to the receipt's group. The comment-input bottom sheet
    (`receipt_comment_screen.dart`) is hidden without create; the swipe-to-delete
    (`receipt_comments.dart`) is disabled without delete (only in **edit** state — add-state comment
    removal is a local model edit with no API call, so it stays enabled).
  - Activity rerun (`group_activity_list_item.dart`) → `group.activities.rerun`. Covered by a widget
    test (`test/widgets/group_activity_list_item_test.dart`) only — there is **no e2e** for it because
    a failed activity's `canBeRestarted` backend state isn't deterministically seedable from the API.
  - The add menu (`show_add_menu.dart`) shows "Add Manual Receipt" on `group.receipts.create` and
    "Quick Scan" on `group.receipts.quick-scan` — per-group when inside a group, or "held in any
    group" on the group-select / all-groups view, where there is no single current group.
  - The **Search** bottom-nav destination (`group_bottom_nav.dart`, `group_select_bottom_nav.dart`) is
    shown only on `app.receipts.search`. It is the **trailing** destination in both navs, so gating it
    out doesn't shift the other indices and the `switch`/`setIndexSelected` logic is unchanged.
- **Route redirects** (`lib/guards/permission-guard.dart`, mirroring the desktop route guards; wired
  in `main.dart`'s `_buildAppRouter`):
  - `groupDashboardReadRedirect` on `/groups/:groupId/dashboards` → redirects to that group's
    `/receipts` when the caller lacks `group.dashboards.read` for the route's group (mirrors desktop's
    `groupDashboardReadGuard`). Group selection always lands on the dashboards route, so this covers
    both group entry and the Dashboards tab tap; the Dashboards tab itself stays visible (redirect-only,
    matching desktop).
  - `receiptsSearchRedirect` on `/search` → bounces deep links to the originating group's `/receipts`
    (or `/groups`) when the caller lacks `app.receipts.search`. Defense-in-depth behind the hidden
    search buttons; it also runs before the search shell's `state.extra as Map` cast.
- **403 handling (`lib/interceptors/auth_interceptor.dart`):** the backend returns **403 for both** an
  expired session and a permission denial (it never sends 401), so the interceptor distinguishes them
  by **token validity** (mirroring desktop's `http-interceptor.ts`): a 403 with a still-valid token is
  a permission denial and is surfaced untouched — **no token refresh, no logout** (refreshing can't
  grant a missing permission, and force-refreshing would burn the one-time refresh token / risk a
  logout). Only an expired/invalid token triggers a refresh + retry. Token freshness is otherwise kept
  current proactively (the 15-min timer in `main.dart` and the auth guard on navigation). **Paid-by
  visibility** rides on this: a group role limited to "their own receipts" gets a server-filtered
  receipts list, and any stray 403 on a hidden receipt is surfaced without disturbing the session.

## Development Notes

### Flutter SDK Setup (Claude Code Environment)

When working in the Claude Code environment, Flutter may not be pre-installed. To install the latest Flutter SDK on Debian/Ubuntu:

```bash
# Prereqs. curl/git/pkg-config/xz-utils are usually already present; the rest
# are required for Linux desktop builds (needed for integration_test runs).
apt-get update && apt-get install -y --no-install-recommends \
  unzip zip clang lld llvm cmake ninja-build libgtk-3-dev

# Download and extract Flutter SDK. Check the current stable version at
# https://storage.googleapis.com/flutter_infra_release/releases/releases_linux.json
# (the `current_release.stable` field names the hash; find its `version`).
cd /tmp && rm -rf flutter && \
curl -fL https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.41.7-stable.tar.xz -o flutter.tar.xz && \
tar xf flutter.tar.xz && rm flutter.tar.xz

# Fix git "dubious ownership" warning, add to PATH persistently, disable analytics.
git config --global --add safe.directory /tmp/flutter
grep -q '/tmp/flutter/bin' /root/.bashrc || echo 'export PATH="/tmp/flutter/bin:$PATH"' >> /root/.bashrc
export PATH="/tmp/flutter/bin:$PATH"
flutter config --no-analytics

# Verify and enable Linux desktop target (needed for ./run-e2e.sh).
flutter --version
flutter config --enable-linux-desktop
flutter devices  # should list "Linux (desktop)"
```

After installing Flutter, standard commands work from `mobile/`:
```bash
cd /app/mobile
flutter pub get      # Install dependencies
flutter analyze      # Check for errors (recommended before building)
flutter test         # Run unit/widget tests
./run-e2e.sh         # Run integration tests on Linux desktop (see E2E Testing below)
flutter build apk    # Build Android APK (requires Android SDK — not installed by default)
```

**Note:** The base environment does not include the Android SDK or Chrome, so `flutter build apk` and web targets will not work without additional setup. Linux desktop + `flutter analyze` + `flutter test` + integration_test are fully supported.

### Regenerating API Client Models

After regenerating the API client with `generate-client.sh`, you need to run build_runner to generate the `.g.dart` files:

```bash
cd /home/user/receipt-wrangler/mobile/api
flutter pub run build_runner build --delete-conflicting-outputs
```

**Known dart-dio default-value regressions (re-patch after every regen).** The `dart-dio`
generator emits invalid `_defaults` initializers for some fields, which `build_runner` then can't
compile (and it deletes the matching `.g.dart` first, so the package stops building). After
regenerating, restore these hand-fixes (precedent: commits `fad192a0`, `a2ec7479`):

- `model/user_preferences.dart` — `..quickScanDefaultStatus = 'OPEN'` →
  `..quickScanDefaultStatus = ReceiptStatus.OPEN`.
- `model/system_settings.dart` — `..currencyDisplay = '$'` → `..currencyDisplay = r'$'` (a bare `$`
  in a non-raw string is invalid Dart).

`model/claims.dart` **no longer** needs a `userRole` patch — the role rework dropped `userRole` from
the swagger, so `Claims` carries only identity claims and the field is gone from the generated model.

Run `flutter analyze` after a regen; these surface as compile errors. (Hand-editing generated files
is otherwise forbidden — these are the documented exception.)

### Quick Scan field configuration

The Quick Scan per-image form (`lib/receipts/widgets/quick_scan_form.dart`) respects the selected
group's quick-scan config on `GroupReceiptSettings` — the
`quickScan{PaidBy,Status,Categories,Tags}{Enabled,Required}` fields (mirrored from the API/desktop;
see `api/CLAUDE.md` → "Quick Scan Field Configuration"). Each field renders per its `*Enabled` flag
and gets a `required()` validator only when shown **and** `*Required`; Categories/Tags reuse the
shared `CategorySelectField` / `TagSelectField`, sourced from the per-group catalogs. Because the
backend **backfills** a default for a hidden/optional paid-by or status, `_submitQuickScan`
(`lib/shared/functions/quick_scan.dart`) sends the "unset" sentinel (`0` / `ReceiptStatus.empty`) for
those and per-file comma-joined `categoryIds` / `tagIds`, building **one aligned array entry per
image** (never skipping, so `files` and the parallel arrays stay 1:1). It requires a field only when
that group's config marks it shown+required, mirroring the backend's `resolveQuickScanFields`. Null
settings (no group selected yet) fall back to the backend defaults: paid-by/status shown,
categories/tags hidden.

### Android Gradle: cunning_document_scanner needs the Kotlin plugin + Kotlin 2.2.x

`cunning_document_scanner` **2.3.0** modernized its `android/build.gradle` to the
`kotlin { compilerOptions { ... } }` DSL but (a) **does not apply the Kotlin Gradle plugin
itself** and (b) requires **Kotlin 2.2.x** for that extension-level `compilerOptions` block.
This app historically applied `kotlin-android` only to `:app` and pinned Kotlin **2.1.0**, so a
clean Android build failed at the Gradle **configuration** step — first
`Could not find method kotlin()` then `Could not find method compilerOptions() ... on
KotlinAndroidProjectExtension` — for **every** build, blocking the whole Android e2e suite.

The fix (all in `mobile/android/`, committed app config — not generated):

- **`build.gradle` (root):** a `subprojects { sp -> sp.plugins.withId('com.android.library')
  { sp.apply plugin: 'org.jetbrains.kotlin.android' } }` block applies `kotlin-android` to any
  Android-library plugin subproject the moment it applies `com.android.library`, so the
  plugin's `kotlin {}` DSL resolves (re-applying an already-applied plugin is a no-op). The
  legacy `buildscript` Kotlin classpath was **removed** from this file — the version is now
  managed in `settings.gradle` (keeping it in both places triggers "plugin already on the
  classpath with an unknown version").
- **`settings.gradle`:** the `pluginManagement` `plugins {}` block declares
  `id "org.jetbrains.kotlin.android" version "2.2.21" apply false` (matching the modern Flutter
  template and the plugin's own example app). This is the lever that actually sets the Kotlin
  version the `kotlin-android` plugin resolves to — bumping the old `ext.kotlin_version` in the
  root `build.gradle` had **no effect**.

If a future dependency bump re-breaks the Android build with a Kotlin/AGP DSL error, the Kotlin
version to bump is in `settings.gradle`, not the root `build.gradle`. The simpler alternative
(if you don't need 2.3.0) is to pin `cunning_document_scanner: ^1.4.0`, whose `build.gradle`
applies `kotlin-android` itself and uses the legacy `kotlinOptions {}` DSL (works on Kotlin
2.1.0); its `getPictures(noOfPages:)` API — the only call the app uses (`lib/utils/scan.dart`) —
is identical.

### Testing

Run tests with `flutter test`. Run a single file with `flutter test test/path/to/file_test.dart`.

**All new code must have accompanying tests.** When adding a new widget, utility, model, or service, add a corresponding test in `test/` that exercises:
- The happy path
- Sign / boundary cases (negative, zero, empty) where applicable
- Wiring contracts (validators, keyboard types, transformers) that downstream code depends on

Existing reference tests:
- `test/services/token_refresh_service_test.dart` — service unit tests with mocktail
- `test/widgets/amount_field_test.dart` — widget tests with FormBuilder + Provider
- `test/utils/currency_test.dart` — pure utility tests
- `test/helpers/widget_test_helpers.dart` — shared widget-test setup helpers
- `test/helpers/auth_test_helpers.dart` — shared mocks and JWT builders

#### Directory layout
Mirror the `lib/` tree: `test/widgets/` for widget tests of `lib/shared/widgets/...`, `test/utils/` for `lib/utils/...`, `test/services/` for `lib/service[s]/...`, `test/interceptors/` for interceptors. Shared helpers go in `test/helpers/`.

#### Flutter widget-test best practices

These patterns are followed by the existing tests; new tests should keep to them:

- **Use `testWidgets` (not `test`) for widget tests.** It supplies the `WidgetTester` and binds the framework.
- **Locate by `Key`, not by widget type.** Pass a `ValueKey` to the widget under test and use `find.byKey(...)`. When you need a specific descendant (e.g. the inner `FormBuilderTextField` of an `AmountField`), use `find.descendant(of: find.byKey(...), matching: find.byType(...))`. `find.byType(...)` alone breaks as soon as another instance lands in the tree.
- **Prefer `pump()` over `pumpAndSettle()`.** `pumpAndSettle` waits for *all* frames to drain and will time out against any continuous animation or formatting-on-change controller (e.g. `currency_textfield`). Reach for `pumpAndSettle` only when a specific test introduces an animation that has to flush.
- **Inject ChangeNotifier dependencies with `ChangeNotifierProvider`.** Use the `create:` constructor when the test owns the instance (auto-disposes); use `.value(value: existing)` only when the test reuses a model created elsewhere.
- **Prefer real model instances over mocks** when the model has no I/O and reasonable defaults (e.g. `SystemSettingsModel`). Mocking via mocktail is for models with I/O or where you need to verify interactions.
- **Only call `registerFallbackValue` when stubs use `any()` matchers.** Concrete `when(() => mock.x()).thenReturn(...)` does not need fallback registration.
- **Don't `tester.enterText` against `currency_textfield` (or any input with a controller that intercepts/reformats keystrokes).** It's fragile across package versions. Test the read path via `initialAmount` round-tripped through `valueTransformer`, and test the write path by inspecting the widget's `keyboardType`.
- **Register the custom currency in `setUpAll`** before any test that calls `exchangeCustomToUSD` / `exchangeUSDToCustom`. The shared helper `registerCustomCurrencyForTests()` in `test/helpers/widget_test_helpers.dart` is idempotent — call it once per test file.
- **Skip golden tests** unless the component is visually critical and the team is set up to maintain reference images.

#### Workflow

1. Write the test alongside the change.
2. `flutter analyze` — must be clean on the new files (the codebase has pre-existing warnings; only check the files you touched).
3. `flutter test` — must be all green.
4. If a test surfaces a real production bug (it happens — e.g. `Money.parse` of a leading `-` against the USD pattern), fix the bug as part of the same change rather than skipping the test.

### E2E Testing

End-to-end tests live in `integration_test/` (sibling of `test/`) and use Flutter's first-party **`integration_test`** package. They drive the real app against a running Go API, mirroring the desktop Playwright suite under `desktop/e2e/`.

**Stack choice:** `integration_test` SDK package. Not Patrol (we don't need native permission dialogs yet). Not the deprecated `flutter_driver`.

**Supported targets:**
- **Local Android emulator** via `./run-e2e-android.sh` (macOS, auto-boots an AVD).
- **Local iOS Simulator** via `./run-e2e-ios.sh` (macOS, auto-boots a sim).
- **Local Linux desktop** via `./run-e2e.sh` (containers/CI Linux). Originally the primary target; kept for the dev container's headless flow.
- **CI Android + iOS** via `.github/workflows/mobile-e2e.yml`, currently **advisory** (`continue-on-error: true`). Triggers: `pull_request` against `main`, `push` to `main` (post-merge), `push` to `tech/mobile-e2e` (iteration on the e2e setup itself), and `workflow_dispatch`. The formerly skipped specs (`receipt_comments_test`, `receipt_cost_split_test`) are un-skipped and green — the product bugs they tracked are fixed — so nothing blocks flipping `continue-on-error` once CI demonstrates stability.

Screenshot/video capture on failure is still deferred — see the "Out of scope" note at the bottom of this section.

#### Prerequisites

**Per target:**

| Target          | Prereqs                                                                                                              |
| --------------- | -------------------------------------------------------------------------------------------------------------------- |
| Android (mac)   | Android SDK (Studio or `cmdline-tools`); at least one AVD; `coreutils` on PATH (`brew install coreutils`, for `gtimeout`). |
| iOS (mac)       | Xcode + iOS Simulator (`xcrun simctl` on PATH); `coreutils` on PATH. The iOS Simulator runtime matching your installed Xcode is required — see "iOS 26.x runtime gotcha" below. |
| Linux desktop   | `apt-get install -y --no-install-recommends libsecret-1-dev xvfb` (`libsecret-1-dev` to *build* `flutter_secure_storage_linux`, `xvfb` to *run* headlessly — `run-e2e.sh` auto-wraps in `xvfb-run` when `$DISPLAY` is unset). Plus `flutter config --enable-linux-desktop`. |

**Shared, one-time:**

1. `cd mobile && flutter pub get`
2. Seed the two e2e users by running `./api/dev/seed-e2e-users.sh` (idempotent — safe to re-run). It logs in as the default `admin/admin` user that `MakeMigrations()` auto-creates on a fresh DB, then calls the admin-protected `POST /user/` to create:
   - `e2e-admin` with role `ADMIN`
   - `e2e-user` with role `USER`

   The script uses creds matching `api/dev/switch-to-sqlite.sh` (and the mariadb/postgresql variants), so a later `source` of those scripts gives the runners credentials that line up with what's seeded. Override `ADMIN_USERNAME` / `ADMIN_PASSWORD` if the default admin's password has been changed; override `API_BASE_URL` to seed a non-local backend.

   Why API-seed instead of the UI: `enableLocalSignUp` is `false` locally so the UI signup 404s, and even if enabled, the auto-`admin/admin` already holds the "first user = ADMIN" slot — a subsequent UI signup for `e2e-admin` would land as USER. `POST /user/` requires admin auth but accepts an explicit `userRole`, so it bypasses both gotchas.

**Every run:** start the Go API separately (`cd api && go run main.go`). None of the runners start the API — same pattern as Playwright.

#### Running locally

Three runners, one per target. Each accepts optional spec paths; with no args it iterates every `*_test.dart` under `integration_test/`.

```bash
# Android emulator (mac dev primary target):
cd mobile && ./run-e2e-android.sh
cd mobile && ./run-e2e-android.sh integration_test/smoke_login_test.dart

# iOS Simulator:
cd mobile && ./run-e2e-ios.sh
cd mobile && ./run-e2e-ios.sh integration_test/smoke_login_test.dart

# Linux desktop (containers, CI Linux):
cd mobile && ./run-e2e.sh
```

All three runners source `api/dev/switch-to-sqlite.sh` for the four `E2E_*` credentials, write a temp dart-define JSON, and invoke flutter against the suite. **What differs:**

| Runner             | Target          | Invocation                       | Default `E2E_BASE_URL`           |
| ------------------ | --------------- | -------------------------------- | -------------------------------- |
| `run-e2e-android.sh` | Android emulator | `flutter drive` (one per spec) | `http://10.0.2.2:8081/api` (host loopback alias) |
| `run-e2e-ios.sh`   | iOS Simulator   | `flutter drive` (one per spec)   | `http://localhost:8081/api`      |
| `run-e2e.sh`       | Linux desktop   | `flutter test` (one per spec)    | `http://localhost:8081/api`      |

**Mobile-runner-specific env overrides:**
- `E2E_ANDROID_AVD` — AVD name (default `Pixel_3a_API_34_extension_level_7_arm64-v8a`). Auto-booted with `-no-snapshot-save -no-boot-anim` if no emulator is attached; the script waits for `sys.boot_completed=1` (5 min ceiling).
- `E2E_IOS_DEVICE` — simulator device name (default `iPhone 15`). Auto-booted via `xcrun simctl boot` + `bootstatus`. The grep anchor (`^ *<name> (`) excludes "iPhone 15 Pro" / "Pro Max" siblings, so set the exact name.
- `E2E_MOBILE_BASE_URL` — overrides the table defaults. Useful for pointing at a remote backend, e.g. `E2E_MOBILE_BASE_URL=https://demo.receiptwrangler.io/api ./run-e2e-android.sh`.

**Why `flutter drive` on mobile, `flutter test` on Linux:** on the Android/iOS path, `flutter test integration_test/...` repeatedly hit "Integration tests and unit tests cannot be run in a single invocation" even on a single file. Every working android-emulator-runner CI example uses `flutter drive` with an explicit driver + target (`test_driver/integration_test.dart` calls `integrationDriver()`), so we follow that pattern on mobile. Linux desktop works fine with `flutter test`.

**Why one-`flutter drive`-per-spec on mobile:** the top-level `GoRouter` in `lib/main.dart` is a final global — its location persists across `testWidgets` calls within the same flutter process. Spec N+1 inherits spec N's last URL and 403s on bootstrap. Looping one `flutter drive` per spec gives each its own process. Between specs, the mobile scripts force-stop and uninstall by bundle id `io.receiptwrangler` (flutter drive's own cleanup uninstalls by the namespace `com.example.receipt_wrangler_mobile`, which doesn't match the real package), `pkill -f dartvm` to free leaked dart isolate hosts that own the vmservice port, and wrap each spec in `gtimeout 600` so a hung run doesn't eat the whole job.

**Auto-discovery:** the mobile scripts walk a small list of candidate Flutter install dirs (`~/Documents/flutter/bin`, `~/flutter/bin`, `/opt/flutter/bin`, `/usr/local/flutter/bin`) if `flutter` is not on `$PATH`, and resolve the Android SDK from `$ANDROID_HOME` → `$ANDROID_SDK_ROOT` → `~/Library/Android/sdk`. They prepend `platform-tools` and `emulator` to PATH for the run.

**Devices are left running on exit.** The scripts don't shut down the emulator/simulator — reruns reuse the booted device for speed. To force a cold boot next time: `adb emu kill` (Android) or `xcrun simctl shutdown <udid>` (iOS).

#### How env vars reach the tests

`String.fromEnvironment` is a `const` constructor — the **key has to be a literal**, so you cannot build it dynamically per role. `integration_test/helpers/env.dart` declares all five `E2E_*` reads as `static const` fields and exposes `E2eEnv.assertAdmin()` / `assertUser()` to fail fast when vars are unset.

**Never use `Platform.environment`** — it returns an empty map on Android/iOS. `--dart-define` is the only portable mechanism.

**Base URL gotcha:** the desktop suite's `E2E_BASE_URL=http://localhost:4200` points at the Angular dev server, whose proxy forwards `/api` to the Go backend. The mobile app has no proxy — it hits the API directly. `run-e2e.sh` therefore reads `E2E_MOBILE_BASE_URL` (defaults to `http://localhost:8081/api`) and maps it into the `E2E_BASE_URL` dart-define the test sees. Override for remote targets: `E2E_MOBILE_BASE_URL=https://demo.receiptwrangler.io/api ./run-e2e.sh`.

#### Writing tests

- **Bootstrap:** call `await tester.pumpWidget(buildApp())` (imported as `import 'package:receipt_wrangler_mobile/main.dart' show buildApp;`). `buildApp()` returns a fresh `MultiProvider` + `ReceiptWrangler` widget tree, with a per-`State` `late final GoRouter` so router location does not leak across `testWidgets`. Do NOT call `app.main()` from a test — `main()` triggers `runApp()` and `FlutterNativeSplash.preserve`, which conflicts with the test binding.
- **`IntegrationTestWidgetsFlutterBinding.ensureInitialized()`** at the top of `main()` in every spec file. Required — `testWidgets` without it runs as a unit test and fails to reach native channels.
- **Gate `installLinuxDesktopMocks()` on `Platform.isLinux`** (from `integration_test/helpers/platform_mocks.dart`), right after the binding. It stubs three mobile-only plugins whose method channels are unimplemented on Linux desktop and would otherwise throw `MissingPluginException` during app bootstrap:
  - `permission_handler` (channel `flutter.baseflow.com/permissions/methods`) — camera permission request in `lib/utils/permissions.dart`.
  - `gal` (channel `gal`) — image-gallery access in the same helper.
  - `flutter_secure_storage` (channel `plugins.it_nomads.com/flutter_secure_storage`) — backed by an in-memory map; real libsecret would need an unlocked gnome-keyring + dbus session, which is fragile in containers/CI.

  On Android/iOS these plugins have real native implementations and must be hit directly, so the `if (Platform.isLinux) { installLinuxDesktopMocks(); }` gate in `smoke_login_test.dart` is the template — copy it into every new spec.
- **Never use `pumpAndSettle` on the bootstrap frame.** `main.dart` renders a `CircularProgressIndicator` inside a `FutureBuilder` during auth init; the indicator's animation means `pumpAndSettle` never returns. Use `pumpUntilFound` (from `integration_test/helpers/pump.dart`) instead — it polls until a target finder hits, with a timeout.
- **Locators:**
  - `FormBuilderTextField` has no Key; match by its `name` field:
    ```dart
    find.byWidgetPredicate((w) => w is FormBuilderTextField && w.name == 'username')
    ```
  - `CupertinoButton.filled` with a `Text` child is `find.widgetWithText(CupertinoButton, 'Log In')`.
- **Assert navigation by widget presence**, not URL. After login, `pumpUntilFound(find.byType(GroupSelect))` is stronger than reading the go_router state — the widget is present iff the `/groups` shell has mounted.
- **Each test cold-boots.** There is no Flutter equivalent of Playwright's `storageState`. When the suite grows past a handful of specs, either accept the per-test login cost or introduce a non-UI setup step. Don't hand-write state sharing between tests.

#### Caveats / things that will bite

- **Three tap-flake patterns** (each cost a debugging session on the Android emulator; recognize them by a
  "derived an Offset ... that would not hit test on the specified widget" warning followed by a downstream timeout).
  iOS makes pattern 2 **deterministic** where Android only flaked — Cupertino page/popup transitions are slower
  (~400ms), so any tap site without the hitTestable-wait + drain fails every run on the simulator:
  1. **`tester.ensureVisible` does not pump.** It jumps the scroll position but the widget's global offset only
     updates after a relayout, so an immediate `tap()` computes the stale (off-screen) center. Always
     `await tester.pump(...)` (or `pumpAndSettle`) between `ensureVisible` and `tap`.
  2. **Modal sheets and popup menus mount content on the animation's first frame.** `pumpUntilFound(find.text(...))`
     returns while the sheet is still sliding in / the popup is still scaling, and the tap lands where the item
     *was*. Wait on `finder.hitTestable()` and then drain a few frames (`for (...) pump(100ms)`) before tapping —
     see `addManualReceiptViaUI` and `receipt_cost_split_test._navigateToEdit` for the canonical shape.
  3. **Snackbars absorb taps on bottom-of-screen buttons.** The root ScaffoldMessenger renders snackbars ABOVE
     modal bottom sheets, so e.g. "Receipt added successfully" covers a sheet's bottom submit button for ~4s.
     `pumpUntilFound(tester, button.hitTestable())` resumes exactly when the snackbar departs.
- **Destination markers must be unique to the destination.** `find.text('Name')` matches on BOTH `/view` and
  `/edit` receipt forms, so it cannot prove an Edit navigation happened — use `find.byType(BottomSubmitButton)`
  (only mounted on edit/add paths) instead.
- **Quick Scan image input on Linux:** the sheet's gallery-upload icon (`getGalleryImages` in `lib/utils/scan.dart`) throws `"Unsupported platform"` on desktop via a `Platform.operatingSystem` switch **before** it reaches `file_selector`, so the file-selector mock can't help and `quick_scan_test.dart` (the gallery happy-path) is `skip: Platform.isLinux`. To reach the Quick Scan form headlessly, feed an image through the **document-scanner** icon (`Icons.add_a_photo`) with `installDocumentScannerMock()`, which works on **all** targets (Linux/iOS/Android).
- **Document-scanner mock must grant camera permission on every platform:** `CunningDocumentScanner.getPictures` requests `Permission.camera` **itself** (Dart-side, via `permission_handler`) before invoking its native channel, so `installDocumentScannerMock()` also calls `installCameraGalleryPermissionMocks()` (extracted from `installLinuxDesktopMocks`) on **all** platforms — not just Linux. Without it, iOS/Android hit the real permission_handler: the fire-and-forget `requestPermissions()` in `main.dart` leaves an app-init camera request **pending** (its native dialog is never dismissed in a headless test), and `getPictures`' own camera request then collides with it → `PlatformException(ERROR_ALREADY_REQUESTING_PERMISSIONS)`. Granting up front makes both requests resolve instantly with no dialog. Only this scan path needs the mobile grant; every other spec hits permission_handler natively (one fire-and-forget request, never a second) and stays green. `flutter_secure_storage` stays **real** on iOS/Android (only Linux mocks it).
- **Linux build linker/ar:** the desktop build resolves its toolchain from the installed clang's dir (e.g. `/usr/lib/llvm-19/bin`). With only `clang` installed you get `Failed to find any of [ld.lld, ld]` then `[llvm-ar, ar]` — install `lld` **and** the matching `llvm` package (see the Flutter SDK Setup apt line above) so `ld.lld` / `llvm-ar` land in that dir.
- **Headless display:** Flutter Linux desktop apps render through GTK and exit immediately without a display. `run-e2e.sh` auto-wraps in `xvfb-run` when `$DISPLAY` is unset. If you see "The log reader stopped unexpectedly, or never started," your display setup isn't working — check `xvfb-run --help` or set `DISPLAY` to a real X server.
- **`libsecret-1-dev` at build time:** the `flutter_secure_storage_linux` plugin's CMakeLists.txt does a `pkg_check_modules(libsecret-1>=0.18.4)` — if the dev headers aren't installed, the build fails with "The following required packages were not found: libsecret-1". Installed as a prereq above.
- **`libsecret` at runtime is avoided via mocks.** We don't bring up gnome-keyring + dbus for tests. `installLinuxDesktopMocks()` intercepts the platform channel with an in-memory map. If you ever want to exercise the real storage path (e.g. to reproduce a token-persistence bug), start a dbus session + gnome-keyring-daemon before the test — but don't do that by default; it adds a lot of fragile state.
- **Go API rate-limiter:** login is rate-limited. Rerunning the same test in tight succession can 429 — give it a few seconds between runs. The desktop suite notes the same issue in `desktop/e2e/helpers/auth.ts`.
- **DB accumulation:** tests write real rows (sessions, refresh tokens). Fine for a smoke test; when specs start creating receipts/groups/etc., build per-test uniqueness (UUIDs) into the data, mirroring the Playwright conventions.
- **Never commit credentials or the generated JSON.** `.e2e-env.json` is gitignored as belt-and-suspenders — the script already uses `mktemp`.
- **iOS 26.x runtime gotcha:** Xcode 26.x ships with only its own SDK (e.g. Xcode 26.5 → iOS 26.5 SDK only). Until that exact-version *simulator runtime* is installed, `xcodebuild -showdestinations` returns **zero** eligible destinations for this project — no sim on any iOS version is buildable, even though older runtimes (17.x / 18.x) show up under `xcrun simctl list runtimes`. Symptom: `run-e2e-ios.sh` reports `Unable to find a destination matching the provided destination specifier: { id:<udid> }` for whatever sim it picks. Fix: `xcodebuild -downloadPlatform iOS` (about 8GB, ~10–15 min). Once installed, the older simulators become buildable too.

#### Reference files

- `integration_test/smoke_login_test.dart` — canonical smoke test.
- `integration_test/permission_add_menu_test.dart` / `permission_receipt_edit_test.dart` — permission-gating coverage (add-menu gate, edit-popup gate, swipe-to-edit gate) using per-spec provisioned users/groups.
- `integration_test/permission_search_test.dart` — search bottom-nav destination gated on `app.receipts.search` (deny via a custom app role minus that permission; allow via a Legacy User).
- `integration_test/permission_dashboard_redirect_test.dart` — group dashboards route gated on `group.dashboards.read` (deny → redirected to the receipts list via a custom group role minus that permission; allow via a Legacy Viewer). Landing is told apart by `GroupReceiptsList` vs `GroupDashboardWrapper`.
- `integration_test/permission_comments_test.dart` — comment **deny** paths on the edit-state comment screen: `group.comments.create` hidden → no input; `group.comments.delete` hidden → swipe-to-delete disabled. Members are provisioned from the **Legacy Editor** baseline (holds `group.receipts.update`, needed to reach edit state) minus the permission under test, via `provisionGroupMemberWithoutPermission(..., baselineRole: 'Legacy Editor')`.
- `integration_test/permission_paid_by_visibility_test.dart` — group-role paid-by visibility: a member restricted to "their own receipts" (via `provisionPaidByOwnMember` → `createRole(..., includeOwnPaidReceipts: true)`) sees only their own receipt in the group list; the admin-paid receipt is filtered out server-side. Mirrors desktop `paid-by-visibility.spec.ts` (list axis).
- `integration_test/permission_receipt_category_visibility_test.dart` — non-admin sees the per-group **category and tag** catalogs in the receipt-form pickers (sourced from `groupCategories` / `groupTags`, not the admin-only flat lists).
- `integration_test/helpers/env.dart` — dart-define consumption + guards.
- `integration_test/helpers/pump.dart` — `pumpUntilFound` polling helper.
- `integration_test/helpers/platform_mocks.dart` — Linux-desktop platform-channel stubs for `permission_handler`, `gal`, `flutter_secure_storage`.
- `integration_test/helpers/login.dart` / `api.dart` — UI + API login as admin, the shared e2e-user, or arbitrary credentials (`loginAs` / `apiLoginAs`).
- `integration_test/helpers/permission_fixtures.dart` — admin-API provisioning for permission specs: fresh user + group with a chosen system group role ("Legacy Viewer"/"Legacy Editor"), optional seeded receipt, `addTearDown` cleanup. Also mints **custom roles** for negative specs: `createRole`/`deleteRole`, `rolePermissionsByName`, and the convenience `provisionUserWithoutAppPermission` / `provisionGroupMemberWithoutPermission` (build a role = a Legacy role **minus one permission**). The backend won't delete an assigned role, so the role-delete teardown is registered **before** the user/group ones — LIFO makes it run last, after the assignments are gone.
- `integration_test/helpers/feature_flags.dart` — `enableAiPoweredReceiptsForTest()`: toggles the Quick Scan flag by pointing systemSettings at a junk receiptProcessingSettings record, restoring on teardown.
- `integration_test/quick_scan_config_response_test.dart` — the Quick Scan per-image **form** shows/hides/requires fields per the selected group's `GroupReceiptSettings.quickScan*`. Green on **all** targets (Linux/iOS/Android): it feeds an image via the mocked document-scanner channel (not the desktop-blocked gallery path) and injects the group config by mutating the live `GroupModel` via Provider (the same technique `quick_scan_disabled_test` uses for the AI flag — deterministic, no reliance on the local API persisting the new fields).
- `integration_test/helpers/document_scanner_mock.dart` — `installDocumentScannerMock()`: stubs the `cunning_document_scanner` channel's `getPictures` to return a fixed on-disk PNG **and** grants camera/gallery permission via `installCameraGalleryPermissionMocks()` on every platform (see the "Document-scanner mock" caveat above). This is the **only** way to add an image to the Quick Scan sheet on Linux desktop.
- `integration_test/helpers/platform_mocks.dart` — `installLinuxDesktopMocks()` (Linux-only: permission_handler + gal + flutter_secure_storage) and the extracted `installCameraGalleryPermissionMocks()` (camera + gallery grant, shared with the document-scanner mock on all platforms).
- `integration_test/helpers/receipt_test_helpers.dart` — `addManualReceiptViaUI` (group-selectable), URL/id extraction, receipt cleanup.
- `integration_test/helpers/nav.dart` — group-entry and group-receipt navigation helpers.
- `test_driver/integration_test.dart` — `integrationDriver()` entrypoint that `flutter drive` uses.
- `run-e2e-android.sh` — Android runner; auto-boots an AVD, runs per-spec `flutter drive` loop with cleanup between specs.
- `run-e2e-ios.sh` — iOS runner; resolves a simulator UDID by name and boots it, runs per-spec `flutter drive` loop.
- `run-e2e.sh` — Linux desktop runner; wraps in `xvfb-run` when headless, invokes `flutter test`.
- `.github/workflows/mobile-e2e.yml` — CI counterpart; same command shapes the local mobile scripts mirror.
- `desktop/e2e/helpers/auth.ts` — Playwright counterpart; follow its conventions when adding new flows.

#### Out of scope (future work)

- Flipping the workflow from advisory (`continue-on-error: true`) to required — unblocked now that every spec runs un-skipped; do it after CI demonstrates stability.
- Screenshot / video artifact capture on failure.
- `storageState`-style auth warmup across a multi-spec suite.
- Additional specs (receipt CRUD, group management, logout).

### Build Configuration
- Android configuration in `android/` directory
- iOS configuration in `ios/` directory  
- Web configuration in `web/` directory
- Custom fonts (Raleway) configured in pubspec.yaml
- Native splash screen and launcher icons configured