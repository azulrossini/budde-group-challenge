# Split Bill

A tiny Go + Vue 3 app that splits one bill between people by percentage. Full spec and build history: [`docs/SPEC.md`](docs/SPEC.md).

## Run

```
docker compose up --build
```

Then open **http://localhost:8080**. Seeded with one bill ("Team dinner") so it's usable immediately.

Tests: `make test` (or `make test-backend` / `make test-frontend` separately).

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

## With a week instead of two hours

- Real auth (a pass-through middleware placeholder is already wired in, with a TODO marking where it plugs in)
- Optimistic concurrency (a `version` column + `409 Conflict`) instead of relying on the row lock alone
- Multiple bills in the UI (tabs), instead of the one seeded bill
- Click-to-edit UX (click a field to edit it in place, Enter to confirm, plain text otherwise) and general visual polish, instead of always-on inputs
- An editable bill total — currently fixed at seed time; needs a new `PATCH /bills/{id}` endpoint plus the backend/contract work behind it
- E2E tests, CI, and a real deployment (Render/Fly.io/Railway + managed Postgres — config is already env-var-only, so this is mostly wiring)

## Left out on purpose

- Authentication / user accounts
- Multiple bills in the UI
- CI and deployment
- Test coverage beyond what's listed in `docs/SPEC.md` §11
