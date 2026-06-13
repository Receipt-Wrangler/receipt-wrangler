# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Receipt Wrangler is a full-stack receipt management and splitting application with OCR-powered scanning, AI-assisted data extraction, and multi-user group management. This is a **monorepo** containing three main components:

- **api/** - Go backend service (port 8081)
- **desktop/** - Angular 19 web interface (port 4200 dev, port 80 production)
- **mobile/** - Flutter cross-platform mobile app
- **docker/** - Monolith Docker build configuration

Each component has its own CLAUDE.md with detailed component-specific guidance. This file covers monorepo-level architecture and workflows.

## Monorepo Architecture

### Component Communication
- **API Contract**: OpenAPI 3.1 specification in `api/swagger.yml` defines the API contract
- **Client Generation**: API clients are auto-generated from swagger.yml using `api/generate-client.sh`
  - Desktop: TypeScript Angular client → `desktop/src/open-api/`
  - Mobile: Dart Dio client → `mobile/api/`
  - MCP: TypeScript client for MCP integration
- **Development Flow**: Changes to API → update swagger.yml → regenerate clients → update frontend

### Technology Stack
- **Backend**: Go 1.24 with Chi router, GORM ORM, Asynq background jobs
- **Frontend**: Angular 19 with NGXS state management, Material + Bootstrap UI
- **Mobile**: Flutter with Provider state management, go_router navigation
- **Infrastructure**: Docker, nginx, PostgreSQL/MySQL/SQLite

## Docker Deployment

### Production Build (Monolith)
The `docker/Dockerfile` builds a single container with both API and web interface:
- Stage 1: Build Angular desktop app
- Stage 2: Build Go API and install dependencies (Tesseract, ImageMagick, Python)
- Final: nginx serves frontend, proxies `/api` to Go backend on port 80

### Development Build
The `docker/dev/Dockerfile` includes:
- All production components plus development tools
- SSH access for debugging (port 22, password: "development")
- Documentation site build from receipt-wrangler-doc repo
- Java runtime for OpenAPI generator
- Flutter SDK at `/opt/flutter` (on `PATH` via `ENV` and `/root/.bashrc`) with Linux desktop enabled and the `mobile/` pub cache warmed, plus `xvfb` + `libsecret-1-dev` so `mobile/run-e2e.sh` works out of the box

### Build Commands
```bash
# Production monolith
docker build -f docker/Dockerfile -t receipt-wrangler .

# Development container
docker build -f docker/dev/Dockerfile -t receipt-wrangler-dev .
```

## API Client Regeneration

When the API swagger.yml changes, regenerate clients:

```bash
# From api/ directory
./generate-client.sh desktop ../desktop/src/open-api
./generate-client.sh mobile ../mobile/api
./generate-client.sh mcp <output-path>
```

**IMPORTANT**: Never manually edit generated client code in `desktop/src/open-api/` or `mobile/api/`. Changes will be overwritten.

## Component Development

### Backend Development (api/)
```bash
cd api
go run main.go                    # Run API server
go test -v ./...                  # Run tests
./set-up-dependencies.sh          # Install system deps (first time)
```

See `api/CLAUDE.md` for detailed backend architecture and testing requirements.

### Frontend Development (desktop/)
```bash
cd desktop
npm start                         # Dev server with API proxy (localhost:4200)
npm test                          # Run tests with coverage
npm run build                     # Production build
```

See `desktop/CLAUDE.md` for Angular architecture, NGXS state management, and component structure.

### Mobile Development (mobile/)
```bash
cd mobile
flutter run                       # Run on device/emulator
flutter test                      # Run tests
flutter build apk                 # Build Android APK
flutter build ios                 # Build iOS app
```

See `mobile/CLAUDE.md` for Flutter architecture, Provider state management, and navigation.

## Critical Cross-Component Considerations

### API Changes Workflow
1. Modify backend code in `api/internal/`
2. Update `api/swagger.yml` to reflect API changes
3. Regenerate clients: `cd api && ./generate-client.sh desktop ../desktop/src/open-api`
4. Update frontend code to use new client methods
5. Test integration between components

### Authentication Flow
- JWT-based authentication with refresh tokens
- Backend issues tokens in `api/internal/handlers/auth.go`
- Desktop stores tokens via NGXS persistent storage
- Mobile uses `flutter_secure_storage` for secure token storage
- All API endpoints except `/api/auth/login` and `/api/auth/signup` require authentication

### Authorization (Roles & Permissions)
- A configurable role system layers on top of auth: administrators define **app-level** and
  **group-level** roles from granular permission strings (e.g. `app.users.create`,
  `group.receipts.read`) and assign them to users / group members.
- **Backend source of truth** is the hardcoded permission registry plus role CRUD in `api/` —
  exposed via `GET /api/permission` and `/api/role`, and mirrored in `swagger.yml` (so regenerated
  clients carry the `Permission` enum and role types). See `api/CLAUDE.md` → "Roles & Permissions".
- The server **never trusts the JWT for authorization** — it re-checks a user's current permissions
  from the database on every request. The JWT no longer carries any role field.
- **Role rollout is complete across backend and desktop.** Handlers fully enforce the permission
  system, and the legacy `UserRole`/`GroupRole` enums have been **removed from the backend** (Go
  types, model fields, JWT role claim, and the `userRole`/`groupRole` API fields are all gone; only
  the physical `user_role`/`group_role` DB columns are retained for the one-time upgrade migration).
  The **desktop** has likewise dropped every legacy-role consumer: the user-list and group-member
  tables now resolve `appRoleId`/`groupRoleId` to a role **name** (via a shared `RoleNamePipe`), the
  group-form "must have an owner" rule is gone (the backend no longer enforces an owner concept), and
  the `AuthState.userRole`/`hasRole` selectors plus the group-member legacy-enum bridge are removed.
  See `api/CLAUDE.md` → "Roles & Permissions" and `desktop/CLAUDE.md`.

### State Management Patterns
- **Backend**: Service layer handles business logic, repositories handle data access
- **Desktop**: NGXS store with actions/selectors, persistent storage for auth/preferences
- **Mobile**: Provider pattern with ChangeNotifier models, models own their state

### Background Processing
- Backend uses Asynq for async jobs (OCR processing, email polling, cleanup)
- Long-running operations (OCR, AI extraction) run as background jobs
- Frontend polls for completion or uses WebSocket-like patterns where implemented

## Version Management

Each component has version tagging scripts:
- `api/tag-version.sh` - Tag API version
- `desktop/tag-version.sh` - Tag desktop version
- `mobile/tag-version.sh` - Tag mobile version

Version is embedded in Docker builds via `VERSION` and `BUILD_DATE` build args.

## Data Persistence

### Development
- API defaults to SQLite in `api/sqlite/`
- Desktop proxy config in `desktop/proxy.conf.json` routes to localhost:8081
- Mobile configures API base URL in app settings

### Production (Docker)
- Volumes for persistent data:
  - `/app/receipt-wrangler-api/data` - Receipt images and uploads
  - `/app/receipt-wrangler-api/sqlite` - SQLite database
  - `/app/receipt-wrangler-api/logs` - Application logs
- nginx serves frontend from `/usr/share/nginx/html`
- API runs on same container, proxied via nginx

## Common Pitfalls

1. **Forgot to regenerate clients**: After API changes, clients are out of sync → regenerate!
2. **Editing generated code**: Changes to `desktop/src/open-api/` or `mobile/api/` will be lost
3. **Missing system dependencies**: API requires Tesseract, ImageMagick → run `api/set-up-dependencies.sh`
4. **Test database cleanup**: Failed Go tests leave `app.db` in test dirs → remove before rerunning
5. **Port conflicts**: API (8081), desktop dev (4200), docker prod (80) must be available
6. **CORS in development**: Desktop proxy handles CORS, but mobile needs proper API base URL

## Project Structure Summary

```
receipt-wrangler-api/          # Monorepo root
├── api/                       # Go backend
│   ├── internal/              # Core application code
│   │   ├── handlers/          # HTTP handlers
│   │   ├── services/          # Business logic
│   │   ├── repositories/      # Database access
│   │   ├── models/            # Data models
│   │   └── wranglerasynq/     # Background jobs
│   ├── swagger.yml            # API specification (source of truth)
│   └── CLAUDE.md              # Backend-specific guidance
├── desktop/                   # Angular web app
│   ├── src/
│   │   ├── app/               # Application modules
│   │   ├── store/             # NGXS state management
│   │   ├── shared-ui/         # Reusable components
│   │   └── open-api/          # Generated API client (DO NOT EDIT)
│   └── CLAUDE.md              # Frontend-specific guidance
├── mobile/                    # Flutter mobile app
│   ├── lib/
│   │   ├── models/            # Provider state models
│   │   ├── groups/            # Group features
│   │   ├── receipts/          # Receipt features
│   │   └── shared/            # Shared widgets
│   ├── api/                   # Generated API client (DO NOT EDIT)
│   └── CLAUDE.md              # Mobile-specific guidance
└── docker/                    # Docker build configs
    ├── Dockerfile             # Production monolith
    └── dev/Dockerfile         # Development container
```

## Code Changes Philosophy

- Prefer minimal, targeted changes. Do not refactor or restructure code beyond what was explicitly requested.
- A primary focus of yours is overall code quality. Your focus should be on producing code that is stable, flexible when
  needed, readable and maintainable. You should not be writing code that is difficult to read, confusing, insecure or
  too long.
- Follow **DRY (Don't Repeat Yourself) pragmatically**. If two or more places share nearly identical logic that would
  need to be updated together, extract it into a shared utility, function, or component. This is not a dogmatic rule —
  three similar lines in a single file or minor template repetition is fine. Apply DRY when it meaningfully reduces
  maintenance burden, not for every tiny duplication.
- When the first approach fails, stop and ask the user for direction rather than trying multiple speculative approaches
  in sequence.
- After you have completed the planning phase, and you have your plan, please iterate over your plan at a maximum of 3
  times. During these iterations, your goals are to verify that your code makes sense, and solves the requested things,
  that your code is sound, secure and consistent with style across the codebase, and that your code is clean, and not a
  hacked together solution.

## Parallel Agent Execution

When a task spans multiple components (e.g., backend `api/` and frontend `desktop/` or `mobile/`), follow these rules:

- **Run backend and frontend agents in parallel** whenever possible. Do not serialize work across components unless
  there is a hard dependency.
- **Frontend agents should order their work to defer backend-dependent tasks.** If the frontend needs something from the
  backend (generated client, models, API endpoints), schedule that work last so independent frontend work happens first.
- **If the frontend agent is blocked on the backend agent** (e.g., waiting for a generated client, new API models, or
  endpoint changes), the frontend agent should:
    1. Continue planning its backend-dependent work (design the component, write the template, stub the types).
    2. **Wait** for the backend agent to finish before executing backend-dependent code. Do not guess at API shapes or
       generate placeholder clients.
    3. Resume execution once the backend deliverables are available.
- **The backend agent should signal completion clearly** — after finishing its work, the orchestrating agent should
  trigger any required client regeneration (e.g., `./generate-client.sh desktop ../desktop/src/open-api`) before
  unblocking the frontend agent.
- **Mobile (`mobile/`) changes** follow the same pattern: if a backend change requires a mobile update, run the mobile
  agent in parallel with the desktop agent after the backend agent completes.

### Example Task Ordering

For a feature that adds a new API endpoint and a corresponding UI:

1. **Phase 1 (parallel):**
    - Backend agent: handler → service → repository → route → tests → swagger update
    - Frontend agent: independent UI work (layout, styling, routing, non-API components)
2. **Phase 2 (sequential, after backend completes):**
    - Regenerate client (`cd api && ./generate-client.sh desktop ../desktop/src/open-api`)
    - Frontend agent: wire up API calls, integrate generated types, write dependent components
3. **Phase 3 (parallel):**
    - Backend agent: any follow-up fixes
    - Frontend agent: integration tests, final UI polish

## Testing

- After ANY code change, run the full relevant test suite before considering the task complete.
- When tests fail, fix both the code AND the tests — don't assume tests are correct or code is correct without
  verifying.

## Workflow Rules

- Always complete implementation AND verify (build + tests pass) before committing. Do not commit code that hasn't been
  validated.
- During your planning sessions, explicitly check if your planned code introduces regressions. We want to make sure that
  we do not break existing code, especially things that may not show themselves through build errors like scss changes,
  conflicting styles, and so on.
- During your planning sessions, take a moment to think if there are any edge cases, or possible regressions or any
  additional things for the user to test before considering the task complete.
- After implementing any full feature, always commit/push.

### Code Review Feedback Disposition

When addressing review feedback from CodeRabbit, human reviewers, or any other source, follow this protocol every time:

1. **Read every comment** before acting. Don't fix-by-fix.
2. **Build a disposition table for the user.** One row per comment: file/line, issue summary, decision (`ACCEPT` / `REJECT` / `ACCEPT + EXTEND` / `DEFER`), and a one-line justification. **Present this table inside the plan you write for the user** — it is for the user to review your reasoning before approving changes. Do NOT post this table as a bulk PR comment; it is a planning artifact, not a review response.
3. **Verify each "ACCEPT" against current code** — reviewers (especially bots) sometimes flag false positives, stale code, or generator output they don't recognize as such. If a flag turns out to be invalid on inspection, flip the decision to `REJECT` with the reason ("verified — generator output, matches existing repo convention" / "verified — code already handles this case at line X").
4. **Default to rejecting** comments that target auto-generated files (`desktop/src/open-api/`, `mobile/api/`) unless the comment identifies a real type/compile error the generator introduced. Hand-edits to generated files require an explicit justification and should match an established project precedent (search `git log` for "Fix build errors" / similar prior hand-patches).
5. **Reply on each individual review-comment thread** with the per-comment decision and justification. Every CodeRabbit (or human) comment must get its own reply explaining whether you accepted or rejected and why — that is the audit trail the reviewer sees. Use `gh api -X POST /repos/<owner>/<repo>/pulls/<num>/comments/<comment_id>/replies` (review comments live on the pulls endpoint, not the issues one) with a JSON body of `{"body": "..."}`. Fetch the review-comment IDs via `gh api /repos/<owner>/<repo>/pulls/<num>/comments`.
6. **Commit + push** only after every comment has an individual reply posted. The per-comment replies are the audit trail; the commit is the action.

## CLAUDE.md Maintenance

- After modifying files in any component, check whether the corresponding `CLAUDE.md` needs updating.
- Each component has its own documentation: `api/CLAUDE.md`, `desktop/CLAUDE.md`, `mobile/CLAUDE.md`.
- If a change alters behavior, configuration, architecture, commands, or conventions documented in a `CLAUDE.md` file,
  update that file to stay accurate before considering the task complete.
