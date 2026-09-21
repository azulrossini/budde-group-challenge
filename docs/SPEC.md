# SPEC — Bill Splitter

## 1. The brief (source requirements)

Technical task, **timebox 2 hours**. A tiny app to split a bill between people by percentage.

**Backend (Go)**
- `GET /bills/{id}/shares` — returns the people and their percentages for a bill.
- `PUT /bills/{id}/shares` — replaces the full set of shares for that bill.
- A schema or migration file for the tables.
- **Rule:** the percentages for a bill must sum to exactly 100.00. A request that breaks this must be rejected with a clear error.
- A bill has an id, a description and a total amount. Seed one bill so the app is usable on first run.

**Frontend (Vue 3)** — a single page that:
- lists the people on the bill and their percentages;
- lets you add and remove people, and edit percentages;
- shows the running total of percentages, and each person's amount in euros;
- prevents saving when the total isn't 100.
- uses a modern, dark UI with white/blue/grey accents and highlights (see §9a).

**Storage:** Postgres (chosen). Include the schema file.

**README** (half a page, no more): how to run it; decisions made and why; what we'd change with a week instead of 2 hours; what was left out on purpose.

**Submission requirements**
1. Starts with **one command**, stated in the README.
2. **At least one test** proving the sum-to-100 rule rejects a bad payload.
3. A git repo with history intact.

**Out of scope (do not implement):** authentication, user accounts, CI, deployment, Dockerfile polish, multiple bills in the UI, test coverage beyond the small tests listed in §11.

> **Amendment (post-brief):** visual design was originally out of scope for the 2-hour timebox, but the author later asked for a deliberate dark theme (see §9a). This is a scoped, intentional addition — not a reopening of the rest of the "out of scope" list.

---

## 2. Tech stack

| Area | Choice |
|---|---|
| Go | 1.27 (`go 1.27` in `go.mod`, `golang:1.27` image) |
| HTTP | Standard library `net/http` (method + path patterns, `r.PathValue`) |
| DB driver | `github.com/jackc/pgx/v5` (`pgxpool`) |
| Migrations | `github.com/pressly/goose/v3` with SQL files embedded via `embed` |
| Codegen (Go) | `github.com/oapi-codegen/oapi-codegen/v2` — models + std-http strict server |
| Frontend | Vue 3 + Vite + TypeScript (`npm create vue@latest`, with Vitest) |
| Codegen (TS) | `openapi-typescript` |
| Database | Postgres (official image, with healthcheck) |
| Run | Docker Compose + multi-stage Dockerfile |

Generated code is **committed** so reviewers don't need the generators installed.

---

## 3. Repository structure

```
split-bill/
  CLAUDE.md
  README.md
  Makefile
  docker-compose.yml
  Dockerfile
  .env.example
  contract/
    openapi.yaml                # single source of truth: endpoints, models, error format, error codes
    oapi-codegen.yaml           # Go codegen config
  docs/
    SPEC.md                     # this file
  backend/
    go.mod
    cmd/server/main.go          # wiring only: config, DB pool, migrations, router, graceful shutdown
    internal/
      config/                   # reads env vars (names are constants)
      models/                   # domain structs: Bill, Share, BillShares
      api/                      # GENERATED: DTOs + strict server interface (the API declaration)
      handlers/                 # implements the generated interface; DTO <-> model mapping; error -> HTTP mapping
      service/                  # business rules + validation (collects ALL field errors)
      repository/               # Postgres access only
      middleware/               # request ID, logging, panic recovery, auth placeholder
      apperr/                   # error type, codes, constructors, messages
      money/                    # basis points / cents helpers, largest-remainder split
      web/                      # go:embed of the built SPA + SPA handler
        dist/.gitkeep
    migrations/
      001_init.sql
  frontend/
    vite.config.ts              # dev proxy: /bills -> http://localhost:8080
    src/
      types/
        api.gen.ts              # GENERATED from the contract
        errors.ts               # ErrorCode (as const), ApiError class
        status.ts               # LoadStatus (as const)
      api/
        client.ts               # the ONLY place that calls fetch; timeout; error parsing
        billsApi.ts              # getShares(billId), replaceShares(billId, payload)
      composables/
        useShares.ts            # state, load/save, add/remove/edit, totals, validation, canSave
      utils/
        money.ts                # parse/format percentages and euros, largest-remainder split
        constants.ts            # BILL_ID, REQUEST_TIMEOUT_MS, LOCALE, CURRENCY, limits
        messages.ts             # user-facing strings
      styles/
        colors.css               # ALL color tokens as CSS custom properties (dark theme, white/blue/grey accents) — see §9a
        base.css                  # resets, typography, layout primitives; imports colors.css; no color literals
      components/
        SharesTable.vue
        ShareRow.vue
        TotalSummary.vue
        ErrorBanner.vue         # shows ApiError message + requestId; optional Retry action
        SharesSkeleton.vue      # loading placeholder rows
        LoadingSpinner.vue      # small inline spinner (used in the Save button)
      App.vue
```

---

## 4. Numbers: basis points and cents (no floats)

| Concept | Representation | Examples |
|---|---|---|
| Percentage | integer **basis points** | 33.33% → `3333`, 100% → `10000` |
| Money | integer **cents** (`int64` in Go) | €120.50 → `12050` |

Constants (Go `internal/money`, TS `src/utils/constants.ts`):
- `FULL_PERCENT_BASIS_POINTS = 10000`
- `MAX_DECIMAL_PLACES = 2`

Rules:
- The API and DB carry only integers. Field names say the unit: `percentageBasisPoints`, `totalCents`, `amountCents`.
- The frontend parses user input **as a string**: trim, accept `"33"`, `"33.3"`, `"33.33"` (and `,` as the decimal separator), reject more than 2 decimals, non-numeric input, and negatives. Never `parseFloat(x) * 100`.
- Display: `3333` → `33.33 %`; cents formatted with `Intl.NumberFormat(LOCALE, { style: 'currency', currency: 'EUR' })`.
- **Largest-remainder split** (same algorithm in Go and TS): for each share compute `total * bp / 10000` with integer floor division and keep the remainder; distribute the leftover cents one at a time to the shares with the largest remainders, ties broken by list order. The amounts always sum exactly to the total.
  - Example: total `1000` cents at `3333 / 3333 / 3334` → exact 333.3 / 333.3 / 333.4 → floors 333/333/333 (999) → leftover 1 cent to the largest remainder → **333 / 333 / 334**.

---

## 5. API contract (`contract/openapi.yaml`)

Routes (paths as in the brief; stored as constants):

**`GET /bills/{billId}/shares`** → `200 BillSharesResponse`

**`PUT /bills/{billId}/shares`** with body `ReplaceSharesRequest` → `200 BillSharesResponse` (returns the saved state so the UI can refresh from the server's truth).

`billId` is a positive integer (`int64`).

Schemas:

```yaml
Bill:                 { id: int64, description: string, totalCents: int64 }
Share:                { id: int64, name: string, percentageBasisPoints: int32, amountCents: int64 }
BillSharesResponse:   { bill: Bill, shares: Share[], totalBasisPoints: int32 }
ShareInput:           { name: string, percentageBasisPoints: int32 }
ReplaceSharesRequest: { shares: ShareInput[] }

ErrorCode (enum):     VALIDATION_ERROR | BAD_REQUEST | NOT_FOUND | UNAUTHORIZED | INTERNAL_ERROR
FieldError:           { field: string, message: string }         # e.g. field "shares[2].percentageBasisPoints"
ErrorResponse:        { code: ErrorCode, message: string, details?: FieldError[], requestId: string }
```

Every operation documents `400`, `404`, `422` and `500` responses with `ErrorResponse`.

---

## 6. Database (`backend/migrations/001_init.sql`, `002_seed.sql`)

Schema (`001_init.sql`):

```sql
CREATE TABLE bills (
  id           BIGSERIAL PRIMARY KEY,
  description  TEXT   NOT NULL CHECK (length(trim(description)) > 0),
  total_cents  BIGINT NOT NULL CHECK (total_cents >= 0),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE shares (
  id             BIGSERIAL PRIMARY KEY,
  bill_id        BIGINT  NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
  person_name    TEXT    NOT NULL CHECK (length(trim(person_name)) > 0),
  percentage_bp  INTEGER NOT NULL CHECK (percentage_bp > 0 AND percentage_bp <= 10000)
);

CREATE UNIQUE INDEX shares_bill_person_unique ON shares (bill_id, lower(person_name));
```

Seed data, in its own migration (`002_seed.sql`) so schema and data can evolve independently:

```sql
INSERT INTO bills (id, description, total_cents) VALUES (1, 'Team dinner', 12050);
INSERT INTO shares (bill_id, person_name, percentage_bp) VALUES
  (1, 'Alice', 5000), (1, 'Bob', 3000), (1, 'Carol', 2000);
SELECT setval('bills_id_seq', (SELECT max(id) FROM bills));
```

- Migrations run **on app startup** via goose (works the same locally and against a hosted DB).
- The sum-to-100 rule is enforced in the service layer, not the DB (a DB-level sum check would need a deferred constraint trigger; noted in the README as a decision).
- `ReplaceShares` runs in **one transaction**: `SELECT ... FROM bills WHERE id = $1 FOR UPDATE` (404 if missing, and it serializes concurrent PUTs on the same bill), `DELETE` the old shares, `INSERT` the new ones, commit.

---

## 7. Validation rules (identical in frontend and backend)

Limits are named constants (`MAX_NAME_LENGTH = 100`, `MAX_SHARES = 50`).

| Rule | Field error |
|---|---|
| At least one share | `shares` |
| At most `MAX_SHARES` shares | `shares` |
| Name required after trimming | `shares[i].name` |
| Name at most `MAX_NAME_LENGTH` characters | `shares[i].name` |
| Names unique (case-insensitive, trimmed) | `shares[i].name` |
| Percentage > 0 and ≤ 10000 basis points | `shares[i].percentageBasisPoints` |
| Sum of percentages == 10000 | `shares` — message includes the actual total, e.g. "Percentages must sum to 100.00 (got 99.50)" |

- The backend **collects all errors** and returns them together in `details` (not just the first one).
- Names are trimmed before saving.
- A person at 0% is rejected (decision: someone paying nothing shouldn't be on the bill).

---

## 8. Error handling

**Backend**
- `apperr.Error{ Code, Message, Details, Err }`; `Code` is a named type with constants matching the contract enum. Messages live in constants.
- `repository` wraps DB errors; `pgx.ErrNoRows` → `apperr.NotFound`.
- `service` returns `apperr.Validation` with field details.
- One mapping function in `handlers`: `VALIDATION_ERROR → 422`, `BAD_REQUEST → 400` (malformed JSON, invalid `billId`, unknown fields), `NOT_FOUND → 404`, `UNAUTHORIZED → 401`, anything else → `500` with a generic message. Wire it into the strict server's request/response error handlers so decoding errors use the same format.
- 500s: log the real error with the request ID via `log/slog`; the client only gets the generic message + `requestId`.
- Middleware:
  - `RequestID`: generates an ID, sets the `X-Request-ID` response header, stores it in the context.
  - `Logging`: method, path, status, duration, request ID.
  - `Recover`: turns panics into a `500 ErrorResponse`.
  - `Auth`: pass-through placeholder, intentionally empty:
    ```go
    // Auth is intentionally a pass-through: authentication is out of scope for this task.
    // TODO(auth): validate credentials here (e.g. Basic Auth or JWT) and return
    // apperr.Unauthorized on failure. Handlers need no changes when this is implemented.
    ```
- Graceful shutdown on SIGINT/SIGTERM.

**Frontend**
- `ApiError` class: `code: ErrorCode`, `message`, `details`, `requestId?`, `status?`.
- `client.ts` converts every failure into an `ApiError`:
  - network failure → `ErrorCode.Network` (frontend-only code);
  - timeout (`AbortController`, `REQUEST_TIMEOUT_MS`) → `ErrorCode.Timeout` (frontend-only code);
  - non-2xx with a JSON `ErrorResponse` body → use it as is;
  - non-2xx with a non-JSON body → `ErrorCode.Internal` with a generic message.
- `ErrorBanner.vue` shows the message (and request ID when present). Field `details` are shown next to the matching row.
- No `console.error` without also showing something to the user.

---

## 9. Loading states (frontend)

`LoadStatus` (as const): `Idle`, `Loading`, `Ready`, `Error`.

| Situation | UI |
|---|---|
| Initial load (`Loading`) | `SharesSkeleton` — 3 placeholder rows with a simple pulse animation, in place of the table |
| Load failed (`Error`) | `ErrorBanner` with a **Retry** button that calls `load()` again |
| Saving (`isSaving = true`) | Save button disabled, shows `LoadingSpinner` + "Saving…"; inputs and add/remove buttons disabled |
| Save succeeded | Replace local state with the server response; brief success message |
| Save failed | `ErrorBanner` + field errors; the user's edits are kept so nothing is lost |

Styling follows the dark theme in §9a; keep components simple and reuse the color tokens rather than introducing one-off styling per component.

---

## 9a. Visual design — dark theme

Scoped, deliberate addition to the brief (see amendment in §1). Goal: modern, dark UI with white/blue/grey accents and highlights.

- **All color values live in one file:** `frontend/src/styles/colors.css`, as CSS custom properties on `:root`. No hex/rgb/hsl literal is allowed in a component's `<style>` block or in TS/JS — components consume `var(--color-*)` tokens only.
- **Palette shape:**
  - Backgrounds: near-black / dark charcoal surfaces, layered (`--color-bg`, `--color-surface`, `--color-surface-raised`) for depth (page vs. card vs. row).
  - Text: white/off-white primary text, muted grey secondary text (`--color-text`, `--color-text-muted`).
  - Accent: a single blue accent family for interactive elements, focus rings, primary actions, and the running-total indicator when it's valid (`--color-accent`, `--color-accent-hover`, `--color-accent-muted`).
  - Borders/dividers: subtle grey, low-contrast against the dark surfaces (`--color-border`).
  - Status: keep the palette restricted to white/blue/grey plus the minimum needed for error/success feedback (e.g. a warning tone for "total ≠ 100" and a success tone for "saved") — these are the only non-blue/grey hues, and they're also tokens, not literals.
- **Structure:** `colors.css` defines tokens only (no component selectors). `base.css` imports it and sets global resets, typography, and the dark `background`/`color` on `body`. Component `<style scoped>` blocks reference tokens via `var(--color-*)`.
- Keep the rest of the UI simple: the dark theme is about color tokens and basic layout polish (spacing, borders, hover/focus states), not a component library or animation work beyond the existing skeleton/spinner.

---

## 10. Frontend behaviour

- Loads bill `BILL_ID = 1` on mount.
- Shows the bill description and total.
- Table: name input, percentage input (string, parsed with `utils/money.ts`), live amount in euros, remove button.
- "Add person" appends an empty row.
- `TotalSummary` shows the running total (e.g. `99.50 % / 100.00 %`) and a warning when it isn't exactly 100.
- `canSave` is true only when: status is `Ready`, not saving, no validation errors, total == 10000, and there are unsaved changes.
- Amounts are recalculated live with the TS largest-remainder split.

---

## 11. Tests (small, focused)

**Backend (Go, table-driven, standard `testing`)**
- `service`: rejects sum 9999 and 10001, a 0 or negative percentage, an empty name, a duplicate name (case-insensitive), an empty list; accepts a valid payload; returns **all** errors at once.
- `money`: the largest-remainder split sums to the total (the 1000-cent 3333/3333/3334 case, a 100% single share, a 0 total).
- `handlers` (**required by the brief**): `PUT` with percentages summing to 99.50 returns `422` with `code == VALIDATION_ERROR` and a sum error in `details`. Use `httptest` with a fake repository behind an interface, so no DB is needed.

**Frontend (Vitest)**
- `utils/money.ts`: `"33.33"` → `3333`, `"33,5"` → `3350`, `"100"` → `10000`; rejects `"33.333"`, `"abc"`, `"-5"`; formats cents as euros; same split cases as Go.
- `useShares`: `canSave` is false at a total of 9999 and true at 10000 (with changes).

---

## 12. Running

- `docker compose up --build` starts:
  - `db`: Postgres with a healthcheck (`pg_isready`);
  - `app`: waits for `db` to be healthy (`depends_on: condition: service_healthy`), runs migrations, serves the API and the SPA on port 8080.
- Multi-stage `Dockerfile`: Node stage builds the frontend → Go stage copies `frontend/dist` into `backend/internal/web/dist` and builds → small runtime image.
- `internal/web` embeds the SPA with `//go:embed all:dist` (the committed `.gitkeep` keeps local builds compiling) and serves `index.html` for `/`.
- Config via env vars only (names are constants): `DATABASE_URL`, `PORT` (default 8080). `.env.example` documents them.
- Makefile targets: `generate`, `test-backend`, `test-frontend`, `test`, `run`, `dev-db`, `dev-backend`, `dev-frontend`.

---

## 13. README outline (half a page max)

1. **Run:** `docker compose up --build` → http://localhost:8080. Tests: `make test`.
2. **Decisions:** integers for money/percentages (basis points + cents); largest-remainder rounding; OpenAPI contract as shared source of truth with generated Go/TS types; layered backend; sum rule in the service (not the DB); row lock for concurrent PUTs; Go serves the SPA (one container, no CORS); Postgres; dark theme with all color tokens centralized in `colors.css`.
3. **With a week:** real auth; optimistic concurrency (a `version` column + `409 Conflict`); multiple bills; e2e tests; CI; deployment (Render/Fly.io/Railway + managed Postgres, config already via env vars); a proper design system / component library.
4. **Left out on purpose:** auth (placeholder middleware left), multiple bills UI, CI, deployment.

---

## 14. Step-by-step plan

Work in this order. After each step: run tests, tick the box, commit, stop for review.

- [ ] **1. Skeleton** — repo layout, `go mod init`, `npm create vue@latest` (TypeScript + Vitest), Makefile, `.gitignore`, `.env.example`.
- [x] **2. Contract** — `contract/openapi.yaml`, codegen config, `make generate`, commit generated Go + TS files.
- [x] **3. Database** — `001_init.sql`, docker-compose `db` service, goose migrations on startup.
- [x] **4. Domain** — `models`, `apperr` (codes, messages), `money` + tests.
- [x] **5. Service** — validation + tests (the sum-to-100 rule).
- [x] **6. Repository** — pgx pool, `GetBillShares`, `ReplaceShares` (transaction + row lock).
- [ ] **7. HTTP** — handlers implementing the generated interface, error mapping, middleware (incl. auth placeholder), `main.go` wiring, graceful shutdown; handler test (422 on a bad sum). Check manually with curl.
- [ ] **8. Frontend foundations** — `constants.ts`, `messages.ts`, `errors.ts`, `status.ts`, `client.ts`, `billsApi.ts`, `utils/money.ts` + tests, `styles/colors.css` + `styles/base.css` (dark theme tokens, see §9a).
- [ ] **9. Frontend UI** — `useShares` + test, components (table, rows, total, error banner, skeleton, spinner) styled with the color tokens, `App.vue`, Vite dev proxy.
- [ ] **10. One command** — Dockerfile, `internal/web` embed, full `docker-compose.yml`; verify `docker compose up --build` from a fresh clone.
- [ ] **11. README** — following §13.
