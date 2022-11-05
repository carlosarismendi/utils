fmt:
 	go fmt .

lint:
	golangci-lint run

up:
	docker-compose up -d 

down:
	docker-compose down

enterpg:
	docker exec -it dddhelper_pg_1 psql --username "postgres"
	
test:
	go test ./...
	