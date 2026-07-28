# Repository Guidelines

## Project Structure & Module Organization

The Go service starts in `cmd/server/main.go`. Keep application packages under `internal/`, grouped by responsibility: `auth/`, `feature_flags/`, `external_api/`, `usage_log/`, `database/`, and shared helpers in `common/`. Database schema changes belong in `db/postgres/migrations/` or `db/clickhouse/migrations/` and use ordered SQL filenames such as `001_init.sql`.

The SvelteKit frontend is in `ui/`: routes live in `ui/src/routes/`, reusable components in `ui/src/lib/components/`, and client state/API helpers in `ui/src/lib`. The `m2m/` directory contains the standalone Node.js Auth0 machine-to-machine traffic generator.

## Build, Test, and Development Commands

Use Task to run backend and infrastructure workflows:

- `task dev` — run the Go service with Air live reload.
- `task build` — compile the server to `bin/server`.
- `task test` or `go test -v ./...` — run all Go unit tests.
- `task lint` — run `golangci-lint run`.
- `task format` — apply `golangci-lint fmt`.
- `task infra:up` / `task infra:down` — start or stop PostgreSQL and ClickHouse via Docker Compose.
- `task db:migrate` — apply migrations using `DB_POSTGRES_URL` and `DB_CLICKHOUSE_URL`.

For the UI, run `npm install`, then `npm run dev`, `npm run check`, or `npm run build` from `ui/`. Run the traffic generator from `m2m/` with `npm start`.

## Coding Style & Naming Conventions

Format Go with `task format` before committing; use tabs and standard Go naming: exported identifiers use `PascalCase`, unexported identifiers use `camelCase`, and package names are short lowercase words. Keep transport handlers, repositories, models, and database wiring in their existing package boundaries. Format TypeScript/Svelte consistently with nearby code; use `PascalCase.svelte` for components and lowercase filenames for utility modules.

## Testing Guidelines

Write Go tests beside the code they cover as `*_test.go`, with test functions named `TestThing` (for example, `TestHealthHandler`). Prefer table-driven tests for multiple inputs and `httptest` for HTTP handlers. Run `task test` and `task lint` before opening a pull request. There is no configured coverage threshold; add focused regression tests for behavior changes and bug fixes.

## Commit & Pull Request Guidelines

Recent commits follow concise conventional prefixes, such as `feat: add simple ui` and `chore: init project`. Use `feat:`, `fix:`, `chore:`, or a similarly clear type, followed by an imperative summary. Pull requests should explain the change, list validation performed, link relevant issues when available, and include screenshots for visible UI changes. Never commit secrets: copy `.env.example` to `.env` locally and keep real Auth0 and database credentials out of Git.
