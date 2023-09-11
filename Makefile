fmt:
	go fmt ./...

lint: fmt
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run

up:
	docker-compose up -d 

down:
	docker-compose down
	
test:
	go test -count=1 -p 1 ./...

usqlBenchmarks:
	go test -benchmem -run=^# -bench ^BenchmarkApplyFilters github.com/carlosarismendi/utils/db/infrastructure/usql -count=1 -p=1 -benchtime=10s 