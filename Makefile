.PHONY: run stop test benchmark mocks

run:
	docker-compose up -d
	cd cmd/server && go run main.go	

stop:
	docker-compose down

test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem ./...

mocks:
	go run github.com/vektra/mockery/v2@latest

