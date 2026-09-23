# Split Bill

A tiny Go + Vue 3 app that splits one bill between people by percentage. Full spec and build history: [`docs/SPEC.md`](docs/SPEC.md).

## Stack

- **Backend:** Go, standard-library `net/http` (Go 1.22+ method/path routing), `pgx`/`pgxpool`, `goose` migrations, `oapi-codegen`.
- **Frontend:** Vue 3 + TypeScript + Vite, Vitest, Naive UI (dark theme).
- **Contract:** one OpenAPI 3 spec (`contract/openapi.yaml`) generates both the Go and TS types.
- **Database:** Postgres.
- **Infra:** Docker Compose + a multi-stage Dockerfile (Go binary embeds the built frontend).

## Run

Prerequisites: Docker Desktop (or Docker Engine + Compose) installed and running.

```
git clone https://github.com/azulrossini/budde-group-challenge.git
cd budde-group-challenge
docker compose up --build
```

Then open **http://localhost:8080**. Seeded with one bill ("Team dinner") so it's usable immediately. First run takes a minute or two (pulling `node`/`golang`/`postgres` base images and building both stages); later runs are cached and much faster.

- **Custom DB credentials:** `cp .env.example .env` and edit it — Compose picks it up automatically.
- **Port conflicts:** this binds host ports `5432` (Postgres) and `8080` (app). If either's already in use, edit the `ports:` mappings in `docker-compose.yml`.
- **Stop:** `Ctrl+C`, then `docker compose down` (add `-v` to also wipe the DB volume).
- **Without Docker** (hot-reload dev loop): `make dev-db`, then `make dev-backend` and `make dev-frontend` in separate terminals — Vite proxies API calls to the Go server.
- **Tests:** `make test` (or `make test-backend` / `make test-frontend` separately).

## Decisions

- **No floats.** Percentages are integer basis points (`10000` = 100%), money is integer cents — in the DB, the API, Go, and TypeScript alike. Avoids the usual rounding bugs entirely.
- **Largest-remainder rounding**, the same algorithm mirrored in Go and TS, so the UI's live preview always matches what the server would actually save.
- **The OpenAPI contract is the source of truth**; Go and TS types are generated from it (`make generate`) and committed, so reviewers don't need the generators installed.
- **Layered backend** (`handlers → service → repository`): validation lives once, in the service; handlers only do HTTP↔model mapping; the repository is the only thing that talks to Postgres.
- **The sum-to-100 rule lives in the service**, not a DB constraint — a DB-level check would need a deferred trigger, not worth it for this size of app.
- **`SELECT ... FOR UPDATE`** on `ReplaceShares` so concurrent `PUT`s on the same bill serialize instead of racing.
- **Go serves the built SPA** (`go:embed`) — one container, one port, no CORS.
- **Postgres**, per the brief.
- **Dark UI** (explicitly requested after the brief): all color tokens centralized in `colors.css`; Naive UI for components, themed from those same tokens at runtime rather than hardcoded per-component styles.
- **A few things are hardcoded for now** (the seeded bill ID, EUR currency, `en-IE` locale) rather than made user-configurable — reasonable for a single-bill demo scope, called out below as one of the first things to generalize.

## With a week instead of two hours

- Real auth (a pass-through middleware placeholder is already wired in, with a TODO marking where it plugs in)
- Optimistic concurrency (a `version` column + `409 Conflict`) instead of relying on the row lock alone
- Remove the remaining hardcoded assumptions (bill ID, currency, locale) in favor of user-configurable settings, alongside multiple bills in the UI (tabs), instead of the one seeded bill
- Click-to-edit UX (click a field to edit it in place, Enter to confirm, plain text otherwise) and general visual polish, instead of always-on inputs
- An editable bill total — currently fixed at seed time; needs a new `PATCH /bills/{id}` endpoint plus the backend/contract work behind it
- A more refined visual design — a stronger color palette and finer layout details beyond this pass's first dark theme
- E2E tests, CI, and a real deployment (Render/Fly.io/Railway + managed Postgres — config is already env-var-only, so this is mostly wiring)

## Left out on purpose

- Authentication / user accounts
- Multiple bills in the UI
- CI and deployment
- Test coverage beyond what's listed in `docs/SPEC.md` §11
