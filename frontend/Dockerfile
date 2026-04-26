# syntax=docker/dockerfile:1

FROM oven/bun:1.2-alpine AS deps

WORKDIR /src/frontend

COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile

FROM oven/bun:1.2-alpine AS build

WORKDIR /src/frontend

COPY --from=deps /src/frontend/node_modules ./node_modules
COPY frontend/ ./
RUN bun run build

FROM oven/bun:1.2-alpine AS runtime

WORKDIR /app

ENV HOST=0.0.0.0
ENV PORT=3000
ENV NUXT_API_BASE_URL=http://localhost:8080
ENV NUXT_PUBLIC_APP_NAME="Vue API Workbench"

COPY --from=build /src/frontend/.output ./.output

EXPOSE 3000

CMD ["bun", ".output/server/index.mjs"]
