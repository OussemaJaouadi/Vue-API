# Vue API Workbench

Vue API Workbench is a Nuxt and Go API workbench for creating collections, configuring request data, proxying execution through the backend, and inspecting responses.

## Repository

- `frontend/`: Nuxt 4 application.
- `backend/`: Go API, auth, storage, migrations, and request execution proxy.
- `.tmp/`: local tool/cache output, ignored by git.

## Request Lifecycle

```mermaid
flowchart TD
  A[User creates or opens a request] --> B[Workbench loads request state]
  B --> C[User edits method, URL path, params, headers, body, and auth]
  C --> D[Save request state]
  D --> E[(Backend database)]
  C --> F[User clicks Execute]
  F --> G[Frontend builds execution payload]
  G --> H[Backend validates workbench bearer token]
  H --> I[Execution proxy sends outbound API request]
  I --> J[Target API returns response]
  J --> K[Backend normalizes status, headers, timing, size, body]
  K --> L[Frontend renders response inspector]
```

## Execution Translation

```mermaid
flowchart LR
  subgraph Request State
    M[Method]
    U[Base URL + Path]
    P[Enabled Params]
    H[Enabled Headers]
    B[Body + Body Language]
    A[Auth Config]
  end

  A --> AT{Auth Mode}
  AT -->|Bearer| AH[Authorization header]
  AT -->|API Key header| KH[Configured header]
  AT -->|API Key query| KQ[Configured query param]
  AT -->|Basic| BH[Basic Authorization header]
  AT -->|OAuth2/OIDC| OH[Bearer Authorization header]
  AT -->|None/Inherit scaffold| NO[No direct injection yet]

  P --> X[Execution Payload]
  H --> X
  AH --> X
  KH --> X
  KQ --> X
  BH --> X
  OH --> X
  M --> X
  U --> X
  B --> X

  X --> BP[Backend Proxy]
  BP --> TA[Target API]
  TA --> R[Response Inspector]
```

## Local Commands

```bash
task frontend:dev
task backend:dev
task db:plan
task db:migrate
```

The latest generated migration must be applied locally before persisted request editor fields are available in an existing database.
