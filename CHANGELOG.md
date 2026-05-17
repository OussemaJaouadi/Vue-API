# Changelog

# Backend

## New

- Extended login to accept either email or username through a single `login` field.
- Added username lookup to the user repository contract, memory adapter, and GORM adapter.
- Redacted sensitive request query values from logs: `ticket`, `token`, `access_token`, and `refresh_token`.
- Quieted expected GORM `record not found` logs.
- Prefer username lookup before email lookup when login input does not contain `@`.

## Old

- Added auth routes: register, login, refresh, logout, and me.
- Added Argon2id password hashing.
- Added access and refresh JWTs with separate secrets.
- Added HttpOnly refresh-token cookie handling.
- Added global roles: `manager`, `user`.
- Added scoped roles: `admin`, `developer`, `tester`.
- Added user records with unique email and unique mutable username.
- Switched user IDs to UUID v7.
- Added SQLite persistence through GORM with explicit approved indexes.
- Added explicit DB generate/plan/migrate commands without startup AutoMigrate.
- Added bootstrap manager seeding before server startup.
- Split env examples into app, HTTP, database, password, JWT, cookie, bootstrap, frontend, and Compose sections.
- Added Docker Compose env propagation for backend auth config.
- Added in-memory SSE broker.
- Added single-use SSE tickets.
- Added `POST /events/ticket`.
- Added `GET /events`.
- Added `user.registered` events for managers.
- Removed stale `backend/internal/auth/.keep`.

# Frontend

## New

- Added workbench shell with collapsible app navigation, project topbar, request tabs, split-pane editor/response layout, and persisted pane sizing.
- Added mock collection/request data consumed through `useWorkbench`.
- Added request controls for method selection, URL editing, params, headers, request body, and request auth configuration.
- Added CodeMirror-backed request/response surfaces with language selection, formatting, minify, copy, bracket closing, fold gutters, and tab handling.
- Added mock execution responses per endpoint.
- Added WebSocket display scaffold with `WS` requests, connect/disconnect state, message composer, mock send/receive events, and timeline display.
- Added backend-owned WebSocket execution notes to the phase docs.
- Added sidebar theme preference control for light, system, and dark modes.
- Improved dark-mode surface, border, input, and primary color contrast.

## Old

- Added shadcn-vue UI components from the configured preset.
- Added app logo asset usage and favicon wiring.
- Refactored the auth pages, app shell, and first dashboard state around the generated UI primitives.
- Added auth session state with memory-only access token storage.
- Added refresh-cookie lifecycle through backend auth routes.
- Added API client bearer injection and one-time refresh retry on `401`.
- Added login and register pages.
- Added auth-aware shell with user info and sign out.
- Added manager and pending-workspace landing states.
- Added SSE connection setup through backend event tickets.
