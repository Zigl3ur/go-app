default: install_deps migrate seed run

install_deps: 
	go mod download

migrate: cmd/migrate/migrate.go
	go run cmd/migrate/migrate.go

seed: cmd/seed/seed.go
	go run cmd/seed/seed.go

run: main.go
	go run main.go

clean:
	go clean