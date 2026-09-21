# CLAUDE.md

Bill splitter: a Go API + Vue 3 SPA that splits one bill between people by percentage.
The full specification and the step-by-step plan live in @docs/SPEC.md. Read it before doing anything.

## How to work

- Implement the plan in `docs/SPEC.md` **one step at a time**, in order. After each step: run the tests, tick the step's checkbox, commit, then **stop and summarize** what you did so I can review.
- Commit after every step with a short conventional message (`feat: ...`, `test: ...`, `chore: ...`). The git history is part of the submission.
- Do not add dependencies that are not listed in the spec without asking first.
- Do not implement anything listed under "Out of scope" in the spec.
- If something in the spec is ambiguous, ask instead of guessing.

## Non-negotiable conventions

1. **No floats for money or percentages.** Percentages are integer basis points (`10000` = 100.00%). Money is integer cents. This applies to the DB, the API, Go and TypeScript. Never use `parseFloat(x) * 100`; parse decimal strings manually.
2. **No literal strings or magic numbers in logic.** Error codes, messages, route paths, env var names, headers, content types, limits and statuses are named constants.
   - Go: named type + `const` block (e.g. `type Code string`).
   - TypeScript: `as const` objects + derived union type. Do not use the TS `enum` keyword.
3. **The contract is the source of truth.** `contract/openapi.yaml` defines all request/response models and the error-code enum. Go and TS types are generated from it (`make generate`). Never hand-edit generated files; change the contract and regenerate.
4. **Strict layering (backend):** `middleware → handlers → service → repository → Postgres`.
   - `handlers`: HTTP ↔ models only. No business rules, no SQL.
   - `service`: all business rules and validation.
   - `repository`: the only package that talks to the database.
   - `models`: all domain structs, in one place.
5. **Every error is handled and shown to the user.** One error shape everywhere (`ErrorResponse` in the contract). Never swallow an error, never `panic` for control flow, never send raw DB/internal errors to the client (log them with the request ID instead). On the frontend, only `src/api/client.ts` calls `fetch`; everything else receives a typed `ApiError`.
6. **Validate on both sides.** Same rules in the frontend (for instant feedback) and the backend (the real guard), plus DB `CHECK` constraints.
7. **Shared helpers are mirrored.** Money/percentage parsing, formatting and the largest-remainder split exist in Go (`internal/money`) and TS (`src/utils/money.ts`) with the **same test cases**.
8. Go packages are named by purpose (`apperr`, `money`), never `common`/`utils`.
9. **Theme tokens live in one place.** All colors are defined once in `frontend/src/styles/colors.css` as CSS custom properties. Components reference `var(--color-*)` tokens; never hardcode a hex/rgb value in a component or in `colors.ts`/JS.

## Commands

- Run everything: `docker compose up --build` → http://localhost:8080
- Generate types from the contract: `make generate`
- Backend tests: `make test-backend` (`cd backend && go test ./...`)
- Frontend tests: `make test-frontend` (`cd frontend && npm test`)
- Dev (hot reload): `make dev-db`, then `make dev-backend` and `make dev-frontend` in separate terminals (Vite proxies API calls to :8080)
