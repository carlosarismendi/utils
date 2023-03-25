
fmt:
	go fmt ./...

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run

up:
	docker-compose up -d 

down:
	docker-compose down

enterpg:
	docker exec -it utils_pg_1 psql --username "postgres"
	
test:
	go test ./...
	