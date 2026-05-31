# Changelog

# Backend

## New

- Added selected collection export by collection id with workspace validation.
- Added Postman collection import adapter for folders, root request items, query params, headers, raw bodies, and simple auth mapping.
- Added YAML import parsing support for OpenAPI and Swagger collection files.
- Added Swagger 2.0 JSON collection import adapter that normalizes legacy operations into persisted request rows.
- Added OpenAPI JSON collection import adapter that maps operations into request rows with query parameters and request body placeholders.
- Added backend import preview endpoint for Workbench, OpenAPI, Swagger, Postman, YAML, invalid JSON, and unknown JSON detection.
- Added native Workbench collection export backend endpoint with workspace permission checks and persisted request payload serialization.
- Added native Workbench collection import backend endpoint with workspace permission checks, duplicate collection renaming, request persistence, and order preservation.
- Added workspace deletion with admin/manager authorization and transactional cleanup of grants, memberships, collections, requests, environments, and variables.
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

- Updated collection export to download the active selected collection when one is selected, otherwise export the whole workspace.
- Added in-sheet collection import result summary with persisted collection/request counts and warning display.
- Updated collection import confirmation to submit raw file content so YAML imports work without browser-side parsing.
- Moved collection import preview detection from frontend-only logic to the backend preview endpoint.
- Added import/export feedback on the collections page with busy export state and success/failure toasts.
- Wired the collections export action to download the backend export payload instead of serializing stale local state.
- Wired the collections import confirmation flow to the backend import endpoint and reloads collections after successful ingestion.
- Added a dedicated Workspaces page for selecting, renaming, creating, and deleting workspaces; moved destructive workspace control out of the sidebar switcher.
- Added frontend workspace deletion state handling with API call, local workspace list update, and preferred-workspace fallback.
- Darkened `--muted-foreground` from oklch(0.552→0.45) in light mode and lightened from oklch(0.765→0.88) in dark mode for WCAG-compliant text contrast.
- Raised fractional opacity text classes globally: `/10→/40`, `/15→/50`, `/20→/50`, `/30→/60`, `/40→/65`, `/55→/70`.
- Raised `placeholder:text-muted-foreground` from `/20→/50` across all inputs.
- Restructured `index.vue` outer shell to `h-[calc(100dvh-5.5rem)] border bg-card` matching Access and Environments page layout.
- Standardized all tactile buttons from `active:scale-95` to `active:translate-x-0.5 active:translate-y-0.5 active:shadow-none`.
- Replaced indigo-500 command palette in `AppTopbar.vue` with primary color tokens.
- Replaced indigo-500 localhost badge in `WorkbenchTabBar.vue` with primary color tokens.
- Removed glowing `shadow-[0_-2px_8px_rgba(16,185,129,0.35)]` from active tab indicator.
- Replaced custom `tracking-[0.15em]`, `tracking-[0.2em]`, `tracking-[0.25em]` values with `tracking-widest` across Response panel, AuthPanel, and empty states.
- Reduced URL input font from `text-[13px]` to `text-[11px]` in `WorkbenchCommandBar.vue`.
- Standardized Send button letter-spacing from `tracking-[0.2em]` to `tracking-widest`.
- Fixed type errors in `mock-data/workbench.ts` (string sizes→number, removed stale `decoded` property).
- Fixed type errors in `useWorkbench.ts` (`appendWebSocketEvent` signature, error response type, added `moveHeader` for headers table drag-to-reorder).
- Fixed type error in `access.vue` (typed `accessOptions` as `AccessLevel[]`).
- Removed drag-drop glow shadows (`shadow-[0_0_10px_rgba(16,185,129,0.45)]`) from `WorkbenchHeadersTable` and `WorkbenchParamsTable`.
- Removed glow shadow and `animate-ping` from `WorkbenchResizer` handle.
- Fixed remaining `text-muted-foreground/20` on delete icons and empty states in `CollectionsWorkbench`, `WorkbenchAuthPanel`, `EnvironmentsPolicyPanel` (bumped to `/50`).
- Fixed remaining `active:scale-90` and `active:scale-95` patterns in `WorkbenchSidebar`, `CollectionsWorkbench`, `AccessPolicyPanel`, `AccessGrantEditor` (standardized to tactile translate pattern).
- Fixed `border-muted-foreground/20` on checkbox toggles in `HeadersTable` and `ParamsTable` (bumped to `/40`).
- Added backend CRUD for Collections (folders + requests) and Environments (environments + variables) with 4 new DB tables.
- Added `/v1/collections` and `/v1/environments` route groups with auth middleware.
- Wired frontend environments page to real backend API (`useEnvironments` composable).
- Wired frontend collections page to load from real backend API (`loadCollections` in workbench).
- Made `addFolder`, `addRequest`, `deleteItem` in workbench composable call backend API.
- Removed hardcoded 'Local'/'Staging' mock environment references from `useCollections` and `useAccess`; derived from real API environments instead.
- Added skeleton loading states and empty states for environments and collections pages.
- Added skeleton loading state for WorkbenchSidebar.
- Removed mock data fallback from workbench composable (useWorkbench no longer imports mock data).

## Old

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
