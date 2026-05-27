# Test Log

## Test-Driven Development: collections, environments, execution

All production code in these packages was deleted and rewritten test-first (strict TDD).

### execution package (`backend/internal/execution/`)

**`service_test.go`** — 28 tests for `httpService.Execute()`
- Basic GET/POST/DELETE requests return correct status, body, headers
- Query params appended, disabled params skipped
- Request headers sent, disabled/empty-key headers skipped
- Invalid URL returns error
- Redirects not followed
- Context timeout cancels
- Response has Size, ExecutionTarget, RequestID (prefixed `req_`)
- Existing URL query params preserved
- Connection refused handled
- All standard HTTP status codes mapped
- JSON response body parseable
- No-content response, duration bounds, custom methods, non-standard status codes

**`manager_test.go`** — 14 tests for `WSManager`
- `NewWSManager` creates empty manager
- Connect with invalid URL / connection refused returns error
- Send/Close on non-existent execution returns "not found" error
- Connect → Send → Close lifecycle works
- Broker events published on connect (`ws.connected`), close (`ws.closed`), send (`ws.message.out`)
- Double close returns error
- Multiple simultaneous connections isolated
- Works with `broker=nil` (no events)

### collection domain (`backend/internal/collection/`)

**`collection_test.go`** — 12 tests
- Error sentinels are distinct (`ErrFolderNotFound`, `ErrRequestNotFound`, `ErrFolderNameTaken`)
- `CreateFolderParams`, `UpdateFolderParams`, `CreateRequestParams`, `UpdateRequestParams` hold/omit values correctly
- `FolderWithRequests` aggregates correctly
- Zero-value defaults for `Folder` and `Request` structs

### environment domain (`backend/internal/environment/`)

**`environment_test.go`** — 11 tests
- Error sentinels are distinct (`ErrEnvironmentNotFound`, `ErrEnvironmentNameTaken`, `ErrVariableNotFound`, `ErrVariableKeyTaken`)
- All params structs hold/omit values correctly
- `EnvironmentWithVariables` aggregates correctly
- Zero-value defaults for `Environment` and `Variable` structs

### collection GORM repository (`backend/internal/storage/gorm/`)

**`collection_repository_test.go`** — 27 tests
- CreateFolder: returns correct fields, default icon (`PhGlobe`), duplicate name error, same name different workspace OK
- ListFolders: all folders for workspace, filters by workspace, empty workspace
- UpdateFolder: name, icon, not-found error
- DeleteFolder: removes folder, cascades to requests, not-found error
- CreateRequest: root-level (nil CollectionID), in-folder, default sort order
- ListRequests: root-level, in-folder, filters by workspace+collection, empty folder
- UpdateRequest: partial update (method only), full update (method+name+path), not-found error
- DeleteRequest: removes request, not-found error

### collection HTTP routes (`backend/internal/http/`)

**`collection_routes_test.go`** — 15 tests
- No auth → 401
- Missing workspaceId → 400
- Get empty collections → 200, empty arrays
- Get with folders and requests → 200, correct nesting
- Create folder missing fields → 400
- Create folder duplicate → 409
- Update folder name/icon → 200
- Update folder not found → 404
- Delete folder → 204, removes from list
- Delete folder not found → 404
- Create request missing fields → 400
- Create request defaults method to GET → 201
- Update request method/name/path → 200
- Update request not found → 404
- Delete request → 204
- Delete request not found → 404

### environment HTTP routes (`backend/internal/http/`)

**`environment_routes_test.go`** — 17 tests
- No auth → 401
- Missing workspaceId → 400
- Get empty → 200, empty array
- Get with data → 200, secret values masked as `••••••••••••••••`
- Create missing fields → 400
- Create duplicate name → 409
- Update name/visibility → 200
- Update not found → 404
- Delete → 204
- Delete not found → 404
- Create variable missing key → 400
- Create variable duplicate key → 409
- Update variable key/value/secret → 200
- Update variable not found → 404
- Delete variable → 204
- Delete variable not found → 404

### environment GORM repository (`backend/internal/storage/gorm/`)

**`environment_repository_test.go`** — 24 tests
- CreateEnvironment: correct fields, default visibility (`project`), duplicate name error, same name different workspace OK
- ListEnvironments: all for workspace, filters by workspace, empty workspace
- UpdateEnvironment: name+visibility, not-found error
- DeleteEnvironment: removes env, cascades variables, not-found error
- CreateVariable: correct fields, secret flag, duplicate key error, same key different env OK
- ListVariables: all for env, empty env
- UpdateVariable: full update (key+value+secret), partial update (value only), not-found error
- DeleteVariable: removes variable, not-found error
