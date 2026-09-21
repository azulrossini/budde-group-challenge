.PHONY: generate test-backend test-frontend test run dev-db dev-backend dev-frontend

generate:
	cd backend && go generate ./...
	cd frontend && npx openapi-typescript ../contract/openapi.yaml -o src/types/api.gen.ts

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test

test: test-backend test-frontend

run:
	docker compose up --build

dev-db:
	docker compose up -d db

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev
