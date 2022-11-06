
fmt:
	go fmt .

lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run

up:
	docker-compose up -d 

down:
	docker-compose down

enterpg:
	docker exec -it dddhelper_pg_1 psql --username "postgres"
	
test:
	go test ./...
	