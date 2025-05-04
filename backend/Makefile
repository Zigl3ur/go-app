default: install_deps build

install_deps: 
	go mod download

build: cmd/main.go cmd/migrate/migrate.go
	go build -o bin/app cmd/main.go
	go build -o bin/migrate cmd/migrate/migrate.go

clean:
	rm -rf bin/
	go clean