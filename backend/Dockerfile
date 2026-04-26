# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build

WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/backend ./cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=build /out/backend /app/backend
COPY backend/migrations /app/migrations

RUN mkdir -p /data && chown -R app:app /app /data

USER app

EXPOSE 8080

ENTRYPOINT ["/app/backend"]
