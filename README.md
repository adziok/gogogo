# Feature Flag Service

A multi-tenant feature flag platform built with **Go**, **PostgreSQL**, **ClickHouse**, **Kafka** and **Auth0**, with a **SvelteKit** admin UI and a machine-to-machine (M2M) traffic generator.

Teams manage flags through a web dashboard. Backend services (the "flag consumers") fetch their tenant's current flag values using Auth0 machine-to-machine credentials, and every read is logged asynchronously and streamed into ClickHouse for analytics.

> This project started as a way to learn Go and backend architecture end-to-end: HTTP routing, authentication, SQL, event streaming, and a frontend on top. The Go code is human-written with LLM assistance; the `ui/` and `m2m/` folders are vibecoded.

## How it works

```
                    ┌──────────────────────┐
   Admin (SPA) ───▶ │  Management API      │ ──▶ PostgreSQL
   Auth0 user JWT   │  /feature-flag CRUD  │     (flags, tenant, audit fields)
                    └──────────────────────┘

                    ┌──────────────────────┐
   Service (M2M) ─▶ │  Consumer API        │ ──▶ PostgreSQL
   client_credentials│  GET /api/{name}    │ ──▶ Kafka "feature-flags-call-log"
                    └──────────────────────┘            │
                                                        ▼
                                            ┌───────────────────────┐
                                            │ Batch consumer        │
                                            │ (8k msgs / 60s)       │ ──▶ ClickHouse
                                            │                       │     usage_log
                                            └───────────────────────┘
```

- **PostgreSQL is the source of truth** for feature flag state.
- **ClickHouse stores analytics** — who requested which flag and when.
- **Kafka decouples analytics from reads**, so a slow or unavailable analytics pipeline never fails a flag read.

## Features

- Tenant-isolated feature flag CRUD (create, list, update, delete).
- Auth0 authentication: user login for the dashboard, `client_credentials` for machine consumers.
- Tenant and user identity derived from JWT claims (`org_id`, `sub`) — never trusted from request input.
- Dedicated consumer endpoint that returns only the requested flag's value.
- Asynchronous, batched Kafka → ClickHouse usage logging with graceful shutdown and manual offset commits.
- Input validation, structured JSON errors and `slog` JSON logging.
- SvelteKit admin UI with login/logout, flag table, create/edit/delete modals and toasts.
- Standalone Node.js M2M traffic generator to exercise the consumer API at load.

## Tech stack

| Layer | Technology |
|---|---|
| Backend | Go 1.26, `go-chi/chi`, `pgx`, `clickhouse-go`, `confluent-kafka-go` |
| Auth | Auth0 (`go-jwt-middleware`, JWKS, RS256) |
| Validation | `go-playground/validator` |
| State store | PostgreSQL 16 |
| Analytics store | ClickHouse |
| Event bus | Apache Kafka |
| Frontend | SvelteKit 2 / Svelte 5, Auth0 SPA JS, Vite |
| Tooling | Task, Air, golangci-lint, Goose migrations, Docker Compose |

## Repository layout

```
cmd/server/            Go entrypoint and HTTP wiring
internal/
  auth/                Auth0 config, JWT validation, claims, user context
  feature_flags/       Management API: handlers, models, repository
  external_api/        Consumer API + Kafka producer
  usage_log/           Kafka batch consumer + ClickHouse writer
  database/            PostgreSQL pool and ClickHouse client setup
  common/              Shared JSON error rendering
db/
  postgres/migrations/ PostgreSQL schema (Goose)
  clickhouse/migrations/ ClickHouse usage_log schema (Goose)
ui/                    SvelteKit admin dashboard
m2m/                   Node.js Auth0 M2M traffic generator
docker-compose.yml     PostgreSQL, ClickHouse and Kafka for local dev
```

## Getting started

Prerequisites: Go, Docker, [Task](https://taskfile.dev), and `air` for live reload.

```bash
cp .env.example .env      # fill in Auth0 and database values
task infra:up             # start PostgreSQL, ClickHouse and Kafka
task db:migrate           # apply migrations
task dev                  # run the server with live reload
```

The service listens on `:8080`; `GET /health` is an unauthenticated readiness check.

Admin UI:

```bash
cd ui
npm install
npm run dev
```

Traffic generator:

```bash
cd m2m
npm install
npm start
```

### Useful task commands

| Command | Description |
|---|---|
| `task dev` | Run the server with Air live reload |
| `task build` | Compile the server to `bin/server` |
| `task test` | Run all Go unit tests |
| `task lint` / `task format` | Run `golangci-lint` / apply formatting |
| `task infra:up` / `task infra:down` | Start / stop local infrastructure |
| `task db:migrate` | Apply PostgreSQL and ClickHouse migrations |

## Status

The core is working: management CRUD, Auth0 authorization, the M2M consumer endpoint, and the Kafka → ClickHouse analytics pipeline. Planned directions include stricter tenant-scoped authorization, automated test coverage for the critical paths, OpenTelemetry/Grafana observability, and caching for high-volume reads. See `PLAN_GPT.md` for the detailed roadmap.
