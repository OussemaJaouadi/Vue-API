# Changelog

# Backend

## New

- Added HTTP execution service with configurable timeout, TLS, TTFB tracing, and header/query injection.
- Added WebSocket execution manager with bidirectional relay and broker event publishing.
- Added authenticated /execute API route group with Bearer token middleware:
  - POST /execute — execute an HTTP request
  - POST /execute/ws — open a WebSocket connection
  - POST /execute/:id/ws/send — send a payload over an active socket
  - DELETE /execute/:id — close a WebSocket execution
- Integrated execution events into the real-time SSE broker.
- Extended login to support email or username lookup.
- Improved logging for security and database noise.
- Upgraded gorilla/websocket and gorm dependencies from indirect to direct.

## Old

- Added HTTP and WebSocket execution engine.
- Integrated execution events into the real-time SSE broker.
- Added authenticated /execute API route group.
- Extended login to support email or username lookup.
- Improved logging for security and database noise.
- Upgraded gorilla/websocket and gorm dependencies from indirect to direct.
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

- Added Access Control page with user roster, grant editor, policy panel, execution matrix, and invite sheet.
- Added Collections page with expand/collapse tree, environment policies, and multi-format import (OpenAPI, Swagger, Postman, Workbench).
- Added Environments page with variable table, secret masking, visibility tiers, and danger-zone delete.
- Added Settings page with theme preference selector.
- Added `useAccess`, `useCollections`, `useEnvironments` composables and `types/access.ts`.
- Integrated workbench with real backend execution replacing all mock responses.
- Added real-time WebSocket event timeline driven by SSE broker events.
- Added auth guard in default layout; moved `loadMe` check out of index page.
- Registered SSE listeners for all WS lifecycle events in useAuthSession.
- Refactored HeadersTable with drag-to-reorder, new column layout, and auto ghost rows.
- Refactored AuthPanel from table layout to labeled form with segmented API-key placement.
- Refactored Response panel with dynamic Content-Type language detection and `formatSize` helper.
- Refactored CodeSurface, WebSocket panel, Editor, Sidebar, RequestPanel, TabBar, Resizer, CommandBar with updated visual tokens and consistent tactile button styles.
- Added AppSidebar nav items with route targets and floating theme dropdown; added AppTopbar accent refinements.
- Added custom-scrollbar, workbench-shadow, wb-active-indicator, and btn-tactile CSS utilities.
- Fixed error handling in login, register, and API client to safely display non-string errors.

## Old

- Integrated workbench with real backend execution for HTTP/WS.
- Refactored page architecture into modular composables (`useAccess`, `useCollections`, `useEnvironments`).
- Added real-time WebSocket event timeline.
- Improved response viewer with dynamic formatting and language detection.
- Added workbench shell with request controls, editor surfaces, and theme support.
- Refactored workbench UI: sidebar navigation, command bar, editor layout, params/headers editors, auth panel, code surface, response viewer, WS panel, and tab bar.
- Added custom scrollbar, workbench shadow, and tactile button utility CSS.
- Fixed error handling in login, register, and API client to safely display non-string errors.
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
