run:
	go run ./cmd/server/main.go

build:
	go build -o bin/server ./cmd/server/main.go

migrate-up:
	migrate -path ./db/migrations -database "$env:DATABASE_URL" up

migrate-down:
	migrate -path ./db/migrations -database "$env:DATABASE_URL" down