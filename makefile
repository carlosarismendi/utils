up:
	docker-compose up -d 

down:
	docker-compose down

enterpg:
	docker exec -it ddd-hexa_pg_1 psql --username "postgres"
	
test:
	go test ./...
	