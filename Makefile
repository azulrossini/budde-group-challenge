.PHONY: generate test-backend test-frontend test run dev-db dev-backend dev-frontend

# Generates Go + TS types from contract/openapi.yaml. Wired up in step 2.
generate:
	cd backend && go generate ./...
	cd frontend && npx openapi-typescript ../contract/openapi.yaml -o src/types/api.gen.ts

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test

test: test-backend test-frontend

# Starts everything with a single command. Wired up in step 10.
run:
	docker compose up --build

# Starts only the Postgres container for local (non-Docker) development.
dev-db:
	docker compose up -d db

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev
