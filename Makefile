run-everything:
	cd misc && docker-compose -f everything-docker-compose.yml up

run-everything-rebuilt:
	cd misc && docker-compose -f everything-docker-compose.yml up --build --force-recreate

linters:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.28.0 golangci-lint run -v