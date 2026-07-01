DATABASE_URL=postgres://eduhub:eduhub_password@localhost:5432/eduhub?sslmode=disable

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

db-tables:
	docker exec -it eduhub-postgres psql -U eduhub -d eduhub -c "\dt"

db-version:
	docker exec -it eduhub-postgres psql -U eduhub -d eduhub -c "SELECT * FROM schema_migrations;"
