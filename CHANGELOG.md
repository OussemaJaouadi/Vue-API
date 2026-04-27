# Changelog

# Backend

## Unreleased

- Added auth routes: register, login, refresh, logout, and me.
- Added Argon2id password hashing.
- Added access and refresh JWTs with separate secrets.
- Added HttpOnly refresh-token cookie handling.
- Added global roles: `manager`, `user`.
- Added scoped roles: `admin`, `developer`, `tester`.
- Added user records with unique email and unique mutable username.
- Switched user IDs to UUID v7.
- Added bootstrap manager seeding before server startup.
- Split env examples into app, HTTP, database, password, JWT, cookie, bootstrap, frontend, and Compose sections.
- Added Docker Compose env propagation for backend auth config.
- Added in-memory SSE broker.
- Added single-use SSE tickets.
- Added `POST /events/ticket`.
- Added `GET /events`.
- Added `user.registered` events for managers.
- Removed stale `backend/internal/auth/.keep`.
